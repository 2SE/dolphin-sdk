package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type Config struct {
	Brokers   []string
	Topic     string
	Offset    int64
	GroupID   string
	Partition int
	MinBytes  int
	MaxBytes  int
	MaxWait   time.Duration
}

func ConsumersInit(cnf *Config) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   cnf.Brokers,
		GroupID:   cnf.GroupID,
		Topic:     cnf.Topic,
		Partition: cnf.Partition,
		MinBytes:  cnf.MinBytes,
		MaxBytes:  cnf.MaxBytes,
		MaxWait:   cnf.MaxWait,
	})
	err := reader.SetOffset(cnf.Offset)
	if err != nil {
		panic(fmt.Errorf("kafka topic %s consumer err:%s", cnf.Topic, err.Error()))
	}
	return reader
}
