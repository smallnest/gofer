package syncx

import (
	"go.etcd.io/etcd/client/v3/concurrency"
)

// NewElection 返回一个etcd Election 对象.
func (d *EtcdSync) NewElection(electKey string) *concurrency.Election {
	return concurrency.NewElection(d.session, electKey)
}
