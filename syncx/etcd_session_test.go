package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestNewDSync(t *testing.T) {
	endpoints := []string{"127.0.0.1:2379"}
	dsync, err := NewDSync(clientv3.Config{Endpoints: endpoints})
	assert.NoError(t, err)

	err = dsync.Close()
	assert.NoError(t, err)
}
