package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"
)

func TestNewKafkaReader(t *testing.T) {
	reader := NewKafkaReader(&Config{
		Brokers: []string{"192.168.10.189:9092"},
		Topic:   "dolphinhub",
		MaxWait: time.Second * 2,
	})

	reader.ReadMessage(context.Background())
}

//{0 0 0 0 0 0 {0s 0s 0s} {0s 0s 0s} {0s 0s 0s} {0 0 0} {0 0 0} 10 100 1s 10s 10s 15s 0 false 0 100  topic01}
//{0 0 0 0 0 0 {0s 0s 0s} {0s 0s 0s} {0s 0s 0s} {0 0 0} {0 0 0} 10 100 1s 10s 10s 15s 0 false 0 100  dolphinhub}
func TestNewKafkaWriter(t *testing.T) {
	writer := NewKafkaWriter(&kafka.WriterConfig{
		Brokers: []string{"192.168.10.189:9092"},
		Topic:   "dolphinhub",
	})

	return
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("key"),
		Value: []byte("val"),
	})
	if err != nil {
		t.Error(err)
		return
	}
}
