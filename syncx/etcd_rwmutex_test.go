package syncx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdSync_RWMutex(t *testing.T) {
	endpoints := []string{"127.0.0.1:2379"}
	etcdSync, err := NewEtcdSync(clientv3.Config{Endpoints: endpoints})
	assert.NoError(t, err)
	defer func() {
		err = etcdSync.Close()
		assert.NoError(t, err)
	}()

	m := etcdSync.NewRWMutex("/defer/rwmutex1")

	// 请求写锁
	t.Log("acquiring write lock")
	err = m.Lock()
	assert.NoError(t, err)
	t.Log("acquired write lock")

	// 等待一段时间
	time.Sleep(time.Second)

	// 释放写锁
	err = m.Unlock()
	assert.NoError(t, err)
	t.Log("released write lock")

	// 请求读锁
	t.Log("acquiring read lock")
	err = m.RLock()
	assert.NoError(t, err)
	t.Log("acquired read lock")

	// 等待一段时间
	time.Sleep(time.Second)

	// 释放写锁
	err = m.RUnlock()
	assert.NoError(t, err)
	t.Log("released read lock")
}
