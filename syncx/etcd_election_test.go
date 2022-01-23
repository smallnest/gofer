package syncx

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdSync_NewElection(t *testing.T) {
	electKey := "/defer/elect1"

	done := make(chan struct{}, 1)
	for i := 0; i < 3; i++ {
		i := i
		v := strconv.Itoa(i)

		go func() {
			endpoints := []string{"127.0.0.1:2379"}
			etcdSync, err := NewEtcdSync(clientv3.Config{Endpoints: endpoints})
			assert.NoError(t, err)
			defer func() {
				err = etcdSync.Close()
				assert.NoError(t, err)
			}()

			elect := etcdSync.NewElection(electKey)

			err = elect.Campaign(context.TODO(), v)
			assert.NoError(t, err)
			t.Logf("#%d becomes master", i)

			err = elect.Proclaim(context.TODO(), "g-"+v)
			assert.NoError(t, err)
			t.Logf("#%d updated its value", i)

			err = elect.Resign(context.TODO())
			assert.NoError(t, err)
			t.Logf("#%d becomes slave", i)

			done <- struct{}{}
		}()
	}

	<-done
	time.Sleep(time.Second)
}
