package syncx

import (
	"context"
	"strings"

	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// Mutex 代表一个排他锁对象，比sync.Locker更多的方法.
type Mutex interface {
	// TryLock locks the mutex if not already locked by another session.
	TryLock(ctx context.Context) error
	// Lock locks the mutex with a cancelable context.
	Lock(ctx context.Context) error
	// Unlock unlocks the mutex with a cancelable context.
	Unlock(ctx context.Context) error
	// IsOwner returns whether current session holds this mutex.
	IsOwner() bool
}

type mutex struct {
	*concurrency.Mutex
	cli *clientv3.Client
}

// NewMutex 返回一个分布式Mutex对象.
func (d *EtcdSync) NewMutex(mutexName string) Mutex {
	if !strings.HasSuffix(mutexName, "/") {
		mutexName += "/"
	}
	mu := concurrency.NewMutex(d.session, mutexName)

	return &mutex{Mutex: mu, cli: d.cli}
}

// IsOwner 返回当前session是否持有此锁.
func (m mutex) IsOwner() bool {
	cmp := m.Mutex.IsOwner()

	resp, err := m.cli.Txn(context.TODO()).If(cmp).Then(clientv3.OpGet(m.Key())).Commit()
	if err != nil || !resp.Succeeded {
		return false
	}

	return cmp.Result == pb.Compare_EQUAL
}
