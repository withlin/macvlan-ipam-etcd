package etcd

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/withlin/macvlan-ipam-etcd/backend/allocator"
)

func connectStore(etcdConfig *allocator.EtcdConfig) (*clientv3.Client, error) {

	var etcdClient *clientv3.Client
	var err error
	if strings.HasPrefix(etcdConfig.EtcdURL, "https") {
		etcdClient, err = connectWithTLS(etcdConfig.EtcdURL, etcdConfig.EtcdCertFile, etcdConfig.EtcdKeyFile, etcdConfig.EtcdTrustedCAFileFile)
	} else {
		fmt.Println("--------------->进来了1")
		etcdClient, err = connectWithoutTLS(etcdConfig.EtcdURL)
	}

	return etcdClient, err
}

/*
	ETCD Related
*/
func connectWithoutTLS(url string) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: 5 * time.Second,
	})
	fmt.Println("--------------->进来了2")

	return cli, err
}

func connectWithTLS(url, cert, key, trusted string) (*clientv3.Client, error) {
	tlsInfo := transport.TLSInfo{
		CertFile:      cert,
		KeyFile:       key,
		TrustedCAFile: trusted,
	}

	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
	})

	return cli, err
}
