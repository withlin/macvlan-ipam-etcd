package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"github.com/withlin/macvlan-ipam-etcd/backend"
	"github.com/withlin/macvlan-ipam-etcd/backend/allocator"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

const ETCDPrefix string = "/etcd-macvlan-cni/networks"

// Store is a simple disk-backed store that creates one file per IP
// address in a given directory. The contents of the file are the container ID.
type Store struct {
	EtcdClient    *clientv3.Client
	EtcdKeyPrefix string
	session       *concurrency.Session
}

// Store implements the Store interface
var _ backend.Store = &Store{}

func New(name string, ipamConfig *allocator.IPAMConfig) (*Store, error) {
	etcdClient, err := connectStore(ipamConfig.EtcdConfig)
	if err != nil {
		return nil, err
	}
	netConfig, err := netConfigJson(ipamConfig)
	if err != nil {
		return nil, err
	}
	etcdKeyPrefix, err := initStore(name, netConfig, etcdClient)
	if err != nil {
		return nil, err
	}
	session, err := concurrency.NewSession(etcdClient, concurrency.WithTTL(3))
	if err != nil {
		return nil, err
	}
	// write values in Store object
	store := &Store{
		EtcdClient:    etcdClient,
		EtcdKeyPrefix: etcdKeyPrefix,
		session:       session,
	}
	return store, nil
}

func initStore(name string, netConfig string, etcdClient *clientv3.Client) (string, error) {
	key := ETCDPrefix + name

	_, err := etcdClient.Put(context.TODO(), key, netConfig)
	if err != nil {
		panic(err)
	}
	return key, nil
}

func netConfigJson(ipamConfig *allocator.IPAMConfig) (string, error) {
	conf, err := json.Marshal(ipamConfig.Ranges)
	return string(conf), err
}

func (s *Store) Reserve(id string, ip net.IP, rangeID string) (bool, error) {
	usedIPPrefix := s.EtcdKeyPrefix + "/used/"

	key := usedIPPrefix + ip.String()
	resp, err := s.EtcdClient.Get(context.TODO(), key)
	if err != nil {
		return false, err
	}
	if len(resp.Kvs) > 0 {
		return false, nil
	}
	value := id
	_, err = s.EtcdClient.Put(context.TODO(), key, value)
	if err != nil {
		return false, err
	}

	key = s.EtcdKeyPrefix + "/lastReserved/" + rangeID
	_, err = s.EtcdClient.Put(context.TODO(), key, ip.String())
	if err != nil {
		return false, err
	}
	return true, nil
}

// LastReservedIP returns the last reserved IP if exists
func (s *Store) LastReservedIP(rangeID string) (net.IP, error) {
	key := s.EtcdKeyPrefix + "/lastReserved/" + rangeID
	resp, err := s.EtcdClient.Get(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, errors.New("Can not find last reserved ip")
	}
	data := string(resp.Kvs[0].Value)
	return net.ParseIP(string(data)), nil
}

func (s *Store) Release(ip net.IP) error {
	key := s.EtcdKeyPrefix + "/used/" + ip.String()
	_, err := s.EtcdClient.Delete(context.TODO(), key)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// N.B. This function eats errors to be tolerant and
// release as much as possible
func (s *Store) ReleaseByID(id string) error {
	key := s.EtcdKeyPrefix + "/used/"
	resp, err := s.EtcdClient.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if len(resp.Kvs) > 0 {
		for _, kv := range resp.Kvs {
			if string(kv.Value) == id {
				_, err = s.EtcdClient.Delete(context.TODO(), string(kv.Key))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Store) Close() error {
	// stub we don't need close anything
	return nil
}

func (s *Store) Lock() error {
	key := s.EtcdKeyPrefix + "/lock/"

	m := concurrency.NewMutex(s.session, key)
	if err := m.Lock(context.TODO()); err != nil {
		return err
	}

	return nil
}

func (s *Store) Unlock() error {
	key := s.EtcdKeyPrefix + "/lock/"
	m := concurrency.NewMutex(s.session, key)
	if err := m.Unlock(context.TODO()); err != nil {
		return err
	}
	return nil
}
