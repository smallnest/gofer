package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestNewEtcdSync(t *testing.T) {
	endpoints := []string{"127.0.0.1:2379"}
	etcdSync, err := NewEtcdSync(clientv3.Config{Endpoints: endpoints})
	assert.NoError(t, err)

	err = etcdSync.Close()
	assert.NoError(t, err)
}
