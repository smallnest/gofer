package syncx

import (
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

// Barrier 代表一个分布式栅栏对象,类似标准库中的sync.WaitGroup,可以实现分布式的任务编排.
type Barrier interface {
	// Hold 创建栅栏的key,让当前goroutine阻塞在Wait方法调用上.
	Hold() error
	// Release 解锁阻塞的goroutine.
	Release() error
	// Wait 阻塞在栅栏key上，直到它被删除.
	Wait() error
}

type barrier struct {
	*recipe.Barrier
}

// NewBarrier 返回一个栅栏对象.
// Barrier creates a key in etcd to block processes, then deletes the key to release all blocked processes.
func (d *EtcdSync) NewBarrier(key string) Barrier {
	b := recipe.NewBarrier(d.cli, key)

	return &barrier{b}
}
