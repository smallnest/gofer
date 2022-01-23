package syncx

import (
	"strings"
	"sync"

	"go.etcd.io/etcd/client/v3/concurrency"
)

// NewLocker 返回一个实现了sync.Locker接口的锁对象.
func (d *EtcdSync) NewLocker(lockerName string) sync.Locker {
	if !strings.HasSuffix(lockerName, "/") {
		lockerName += "/"
	}

	return concurrency.NewLocker(d.session, lockerName)
}
