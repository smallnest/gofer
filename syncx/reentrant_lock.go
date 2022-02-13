package syncx

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kortschak/goroutine"
)

// ReentrantLock is held by the goroutine last successfully locking, but not yet unlocking it.
type ReentrantMutex struct {
	sync.Mutex
	owner     int64
	recursion int32
}

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *ReentrantMutex) Lock() {
	gid := goroutine.ID()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// we are now inside the lock
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// This lock must be unlock by the held goroutine, not like Mutex in the standard lib.
func (m *ReentrantMutex) Unlock() {
	gid := goroutine.ID()
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}
