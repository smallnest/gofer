package syncx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestDSync_NewMutex(t *testing.T) {
	endpoints := []string{"127.0.0.1:2379"}
	dsync, err := NewDSync(clientv3.Config{Endpoints: endpoints})
	assert.NoError(t, err)
	defer func() {
		err = dsync.Close()
		assert.NoError(t, err)
	}()

	mutexName := "/defer/mutex1"
	mutex := dsync.NewMutex(mutexName)
	assert.False(t, mutex.IsOwner())

	err = mutex.Lock(context.Background())
	assert.NoError(t, err)

	assert.True(t, mutex.IsOwner())

	err = mutex.TryLock(context.Background())
	assert.NoError(t, err)

	err = mutex.Unlock(context.Background())
	assert.NoError(t, err)
	assert.False(t, mutex.IsOwner())
}
