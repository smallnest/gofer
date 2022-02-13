package syncx

import (
	"context"

	"github.com/marusama/cyclicbarrier"
)

// CyclicDoubleBarrier is a synchronizer that allows a set of goroutines to wait for each other to enter,
// and wait for each other to leave.
type CyclicDoubleBarrier struct {
	b1 cyclicbarrier.CyclicBarrier
	b2 cyclicbarrier.CyclicBarrier
}

// NewCyclicDoubleBarrier 新建一个可重用的双次栅栏 CyclicDoubleBarrier.
func NewCyclicDoubleBarrier(parties int) *CyclicDoubleBarrier {
	if parties <= 0 {
		panic("parties must be positive number")
	}
	return &CyclicDoubleBarrier{
		b1: cyclicbarrier.New(parties),
		b2: cyclicbarrier.New(parties),
	}
}

// Enter 等待 parties goroutine进入.
func (b *CyclicDoubleBarrier) Enter(ctx context.Context) {
	b.b1.Await(ctx)
}

// Leave 等待 parties goroutine离开.
func (b *CyclicDoubleBarrier) Leave(ctx context.Context) {
	b.b2.Await(ctx)
}
