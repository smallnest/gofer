package syncx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestDSync_Barrier(t *testing.T) {
	endpoints := []string{"127.0.0.1:2379"}
	dsync, err := NewDSync(clientv3.Config{Endpoints: endpoints})
	assert.NoError(t, err)
	defer func() {
		err = dsync.Close()
		assert.NoError(t, err)
	}()

	barrier := dsync.NewBarrier("/defer/barrier1")

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			t.Logf("hold #%d", i)
			err := barrier.Hold()
			if err != nil {
				t.Logf("someone has created so #%d doesn't need to create", i)
			}

			err = barrier.Wait()
			assert.NoError(t, err)
		}()
	}

	time.Sleep(time.Second)

	err = barrier.Release()
	assert.NoError(t, err)

	// release again
	err = barrier.Release()
	assert.NoError(t, err)

	// wait on a released key
	err = barrier.Wait()
	assert.NoError(t, err)
}
