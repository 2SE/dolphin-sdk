package kafka

import (
	"context"
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
type WriteConfig struct {
	Brokers []string
	Topic   string
	//retry times
	MaxAttempts int
	//队列容量
	QueueCapacity int
}

func NewKafkaReader(cnf *Config) *kafka.Reader {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	for k, v := range cnf.Brokers {
		_, err := kafka.DialContext(ctx, "tcp", v)
		if err != nil {
			panic(fmt.Errorf("the broker dial err %v and the index is %d", err, k))
		}
	}

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

func NewKafkaWriter(c *kafka.WriterConfig) *kafka.Writer {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	for k, v := range c.Brokers {
		_, err := kafka.DialContext(ctx, "tcp", v)
		if err != nil {
			panic(fmt.Errorf("the broker dial err %v and the index is %d", err, k))
		}
	}
	writer := kafka.NewWriter(*c)
	return writer
}
