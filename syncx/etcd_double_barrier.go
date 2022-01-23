package syncx

import (
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

// Barrier 代表一个分布式双次栅栏对象.
// 所有的分布式goroutine都会阻塞在Enter方法调用上，直到指定数量的goroutine调用了;
// 然后又会阻塞在Leave方法调用上，直到所有的goroutine都调用了Leave方法.
type DoubleBarrier interface {
	// Enter 等待"count"个goroutine进入才返回.
	Enter() error
	// Leave 等待"count"个goroutine调用Leave才返回.
	Leave() error
}

type doubleBarrier struct {
	*recipe.DoubleBarrier
}

// NewBarrier 返回一个双次栅栏对象.
func (d *EtcdSync) NewDoubleBarrier(key string, count int) DoubleBarrier {
	b := recipe.NewDoubleBarrier(d.session, key, count)

	return &doubleBarrier{b}
}
