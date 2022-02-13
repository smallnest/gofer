package syncx

import (
	"context"
	"errors"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type rmutex struct {
	*redsync.Mutex
}

func NewRedLock(opts []*redisv8.Options, mutexName string, options ...redsync.Option) Mutex {
	var pool []redis.Pool

	for _, opt := range opts {
		client := redisv8.NewClient(opt)
		p := goredis.NewPool(client)
		pool = append(pool, p)
	}

	rs := redsync.New(pool...)
	mu := rs.NewMutex(mutexName, options...)

	return &rmutex{Mutex: mu}
}

// TryLock locks the mutex if not already locked by another session.
func (m *rmutex) TryLock(ctx context.Context) error {
	return errors.New("not implemented")
}

// Lock locks the mutex with a cancelable context.
func (m *rmutex) Lock(ctx context.Context) error {
	err := m.Mutex.LockContext(ctx)

	return err
}

// Unlock unlocks the mutex with a cancelable context.
func (m *rmutex) Unlock(ctx context.Context) error {
	_, err := m.Mutex.UnlockContext(ctx)

	return err
}

// IsOwner returns whether current session holds this mutex.
func (m *rmutex) IsOwner() bool {
	ok, _ := m.Mutex.Valid()

	return ok
}
