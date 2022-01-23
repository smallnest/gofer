package syncx

import (
	"strings"

	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

// RWMutex is a reader/writer mutual exclusion lock.
type RWMutex interface {
	// Lock locks rw for writing.
	Lock() error
	// Unlock unlocks rw for writing.
	Unlock() error
	// RLock locks rw for reading.
	RLock() error
	// RUnlock undoes a single RLock call;
	RUnlock() error
}

// NewBarrier 返回一个栅栏对象.
// Barrier creates a key in etcd to block processes, then deletes the key to release all blocked processes.
func (d *EtcdSync) NewRWMutex(rwMutexName string) RWMutex {
	if !strings.HasSuffix(rwMutexName, "/") {
		rwMutexName += "/"
	}

	return recipe.NewRWMutex(d.session, rwMutexName)
}
