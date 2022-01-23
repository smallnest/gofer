package mka

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/segmentio/kafka-go"
	"go.uber.org/multierr"
)

// Reader 代表一个支持多Kafka集群的reader.
// 它会从多个kafka集群同时读取消息.
type Reader struct {
	configs []kafka.ReaderConfig

	idx     uint64
	n       int
	readers []*kafka.Reader

	wp *workerpool.WorkerPool
}

// NewReader 返回一个支持多Kafka集群的reader.
func NewReader(configs []kafka.ReaderConfig) *Reader {
	if len(configs) == 0 {
		panic("must set at least one kafka cluster")
	}

	var readers []*kafka.Reader
	for _, config := range configs {
		readers = append(readers, kafka.NewReader(config))
	}

	n := len(configs)

	return &Reader{
		configs: configs,

		idx:     0,
		n:       n,
		readers: readers,

		wp: workerpool.New(n),
	}
}

// Close 关闭所有的reader, 阻止程序读取更多的kafka消息.
func (r *Reader) Close() error {
	var err error
	for _, r := range r.readers {
		e := r.Close()
		if e != nil {
			err = multierr.Append(err, e)
		}
	}

	r.wp.Stop()

	return err
}

// ReadMessage reads from all kafka clusters and return the next messages from the r. The method call
// blocks until at least a message becomes available, or all readers return errors. The program
// may also specify a context to asynchronously cancel the blocking operation.
//
// The method returns io.EOF to indicate that the reader has been closed.
//
// If consumer groups are used, ReadMessage will automatically commit the
// offset when called. Note that this could result in an offset being committed
// before the message is fully processed.
//
// If more fine grained control of when offsets are  committed is required, it
// is recommended to use FetchMessage with CommitMessages instead.
func (r *Reader) ReadMessage(ctx context.Context) ([]kafka.Message, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var mu sync.Mutex
	var msgs []kafka.Message
	var err error

	var wg sync.WaitGroup
	wg.Add(r.n)

	idx := atomic.AddUint64(&r.idx, 1) % uint64(r.n)
	for i := 0; i < r.n; i++ {
		reader := r.readers[(idx+uint64(i))%uint64(r.n)]

		r.wp.Submit(func() {
			defer wg.Done()

			msg, e := reader.ReadMessage(ctx)
			if e == nil {
				mu.Lock()
				msgs = append(msgs, msg)
				mu.Unlock()
				cancel()
			} else if !errors.Is(e, context.Canceled) {
				err = multierr.Append(err, e)
			}
		})
	}
	wg.Wait()

	if len(msgs) > 0 {
		return msgs, nil
	}

	return msgs, err
}

// Lag returns the lag of the last message returned by ReadMessage, or -1
// if r is backed by a consumer group.
func (r *Reader) Lag(i int) int64 {
	if i >= r.n {
		return 0
	}

	return r.readers[i].Lag()
}

// Offset returns the current absolute offset of the reader, or -1
// if r is backed by a consumer group.
func (r *Reader) Offset(i int) int64 {
	if i >= r.n {
		return 0
	}

	return r.readers[i].Offset()
}

// ReadLag returns the current lag of the reader by fetching the last offset of
// the topic and partition and computing the difference between that value and
// the offset of the last message returned by ReadMessage.
//
// This method is intended to be used in cases where a program may be unable to
// call ReadMessage to update the value returned by Lag, but still needs to get
// an up to date estimation of how far behind the reader is. For example when
// the consumer is not ready to process the next message.
//
// The function returns a lag of zero when the reader's current offset is
// negative.
func (r *Reader) ReadLag(ctx context.Context, i int) (lag int64, err error) {
	if i >= r.n {
		return 0, errors.New("wrong index")
	}

	return r.readers[i].ReadLag(ctx)
}

// SetOffset changes the offset from which the next batch of messages will be
// read. The method fails with io.ErrClosedPipe if the reader has already been closed.
//
// From version 0.2.0, FirstOffset and LastOffset can be used to indicate the first
// or last available offset in the partition. Please note while -1 and -2 were accepted
// to indicate the first or last offset in previous versions, the meanings of the numbers
// were swapped in 0.2.0 to match the meanings in other libraries and the Kafka protocol
// specification.
func (r *Reader) SetOffset(i int, offset int64) error {
	if i >= r.n {
		return errors.New("wrong index")
	}

	return r.readers[i].SetOffset(offset)
}

// SetOffsetAt changes the offset from which the next batch of messages will be
// read given the timestamp t.
//
// The method fails if the unable to connect partition leader, or unable to read the offset
// given the ts, or if the reader has been closed.
func (r *Reader) SetOffsetAt(ctx context.Context, i int, t time.Time) error {
	if i >= r.n {
		return errors.New("wrong index")
	}

	return r.readers[i].SetOffsetAt(ctx, t)
}

// Stats returns a snapshot of the reader stats since the last time the method
// was called, or since the reader was created if it is called for the first
// time.
//
// A typical use of this method is to spawn a goroutine that will periodically
// call Stats on a kafka reader and report the metrics to a stats collection
// system.
func (r *Reader) Stats(i int) kafka.ReaderStats {
	if i >= r.n {
		return kafka.ReaderStats{}
	}

	return r.readers[i].Stats()
}
