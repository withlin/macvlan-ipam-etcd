package etcd

import (
	"testing"

	"github.com/withlin/macvlan-ipam-etcd/backend/allocator"
)

func Test_Lock(t *testing.T) {
	ipamConfig := &allocator.IPAMConfig{
		Range: &allocator.Range{
			EtcdConfig: &allocator.EtcdConfig{
				EtcdURL: "http://localhost:2379",
			},
		},
	}

	s1, err := New("", ipamConfig)
	err = s1.Unlock()
	if err != nil {
		t.Error(err)
	}
	err = s1.Lock()
	if err != nil {
		t.Error(err)
	}
	s2, err := New("", ipamConfig)
	err = s2.Lock()
	if err != nil {
		t.Error(err)
	}

}
