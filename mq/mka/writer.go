package mka

import (
	"context"
	"sync/atomic"

	"github.com/gammazero/workerpool"
	"github.com/segmentio/kafka-go"
	"go.uber.org/multierr"
)

// RWMode 支持多写还是主备模式.
type RWMode byte

const (
	// RWModeMultiRW 多写模式.
	RWModeMultiRW RWMode = iota
	// RWModeBackup 主从模式，正常情况下写入主，主有问题时随机选择一个从.
	RWModeBackup
)

// Writer 支持多写Kafka集群.
// 可以选择多写模式还是主备模式.
// - 多写模式下: 轮询选择一个kafka集群进行写入
// - 主备模式: 优先写入主, 主失败的情况下写入从
type Writer struct {
	rwmode  RWMode
	configs []kafka.WriterConfig

	idx     uint64
	n       int
	writers []*kafka.Writer

	wp *workerpool.WorkerPool
}

// NewReNewWriterader 返回一个支持多Kafka集群的writer.
func NewWriter(rwmode RWMode, configs []kafka.WriterConfig) *Writer {
	if len(configs) == 0 {
		panic("must set at least one kafka cluster")
	}

	var writers []*kafka.Writer
	for _, config := range configs {
		writers = append(writers, kafka.NewWriter(config))
	}

	n := len(configs)

	return &Writer{
		rwmode:  rwmode,
		configs: configs,

		idx:     0,
		n:       n,
		writers: writers,

		wp: workerpool.New(n),
	}
}

// Close flushes pending writes, and waits for all writes to complete before
// returning. Calling Close also prevents new writes from being submitted to
// the writer, further calls to WriteMessages and the like will fail with
// io.ErrClosedPipe.
func (w *Writer) Close() error {
	var err error
	for _, w := range w.writers {
		e := w.Close()
		if e != nil {
			err = multierr.Append(err, e)
		}
	}

	w.wp.Stop()

	return err
}

// Stats returns a snapshot of the selected writer stats since the last time the method
// was called, or since the writer was created if it is called for the first
// time.
//
// A typical use of this method is to spawn a goroutine that will periodically
// call Stats on a kafka writer and report the metrics to a stats collection
// system.
func (w *Writer) Stats(i int) kafka.WriterStats {
	if i >= w.n {
		return kafka.WriterStats{}
	}

	return w.writers[i].Stats()
}

/// WriteMessages writes a batch of messages to the kafka topic configured on this
// writers.  If write fails, it will try write another kafka cluster again.
//
// Unless the writer was configured to write messages asynchronously, the method
// blocks until all messages have been written, or until the maximum number of
// attempts was reached.
//
// When sending synchronously and the writer's batch size is configured to be
// greater than 1, this method blocks until either a full batch can be assembled
// or the batch timeout is reached.  The batch size and timeouts are evaluated
// per partition, so the choice of Balancer can also influence the flushing
// behavior.  For example, the Hash balancer will require on average N * batch
// size messages to trigger a flush where N is the number of partitions.  The
// best way to achieve good batching behavior is to share one Writer amongst
// multiple go routines.
//
// When the method returns an error, it may be of type kafka.WriteError to allow
// the caller to determine the status of each message.
//
// The context passed as first argument may also be used to asynchronously
// cancel the operation. Note that in this case there are no guarantees made on
// whether messages were written to kafka. The program should assume that the
// whole batch failed and re-write the messages later (which could then cause
// duplicates).
func (w *Writer) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	var writer *kafka.Writer
	var idx uint64

	if w.rwmode == RWModeBackup {
		writer = w.writers[0]
	} else {
		idx = atomic.AddUint64(&w.idx, 1) % uint64(w.n)
		writer = w.writers[idx]
	}

	err := writer.WriteMessages(ctx, msgs...)
	if err == nil {
		return nil
	}

	if w.n == 1 {
		return err
	}

	if w.rwmode == RWModeBackup {
		idx := atomic.AddUint64(&w.idx, 1) % uint64(w.n-1)
		writer = w.writers[1:][idx]
	} else {
		idx := (idx + 1) % uint64(w.n)
		writer = w.writers[idx]
	}

	return writer.WriteMessages(ctx, msgs...)
}
