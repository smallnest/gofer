package mka

import (
	"context"
	"strconv"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestSingleCluster(t *testing.T) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test",
	})

	for i := 0; i < 10; i++ {
		err := w.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte("Key-" + strconv.Itoa(i)),
			Value: []byte("Hello World: " + strconv.Itoa(i)),
		})

		assert.NoError(t, err)
	}
	err := w.Close()
	assert.NoError(t, err)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "test-group1",
		Topic:   "test",
	})

	for i := 0; i < 10; i++ {
		msg, err := reader.ReadMessage(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, msg)
	}
}

func TestReadWrite(t *testing.T) {
	// write
	writer := NewWriter(RWModeMultiRW, []kafka.WriterConfig{
		{Brokers: []string{"localhost:9092"}, Topic: "test"},
		{Brokers: []string{"localhost:9092"}, Topic: "test"},
	})

	for i := 0; i < 10; i++ {
		err := writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte("Key-" + strconv.Itoa(i)),
			Value: []byte("Hello World: " + strconv.Itoa(i)),
		})

		assert.NoError(t, err)
	}

	err := writer.Close()
	assert.NoError(t, err)

	// read
	reader := NewReader([]kafka.ReaderConfig{
		{
			Brokers: []string{"localhost:9092"},
			GroupID: "test-group",
			Topic:   "test",
		},
		{
			Brokers: []string{"localhost:9092"},
			GroupID: "test-group",
			Topic:   "test",
		},
	})
	defer reader.Close()

	count := 0
	for {
		msgs, err := reader.ReadMessage(context.Background())
		assert.NoError(t, err)
		assert.Greater(t, len(msgs), 0)
		count += len(msgs)
		if count >= 10 {
			break
		}
	}
}
