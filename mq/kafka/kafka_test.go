package kafka

import (
	"context"
	"testing"
)

func TestNewKafka(t *testing.T) {
	reader := NewKafka(&Config{
		Brokers: []string{"192.168.10.189:9092"},
		Topic:   "topic01",
	})
	reader.ReadMessage(context.Background())
}
