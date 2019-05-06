package redis

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestNewRedis(t *testing.T) {
	op := &redis.Options{
		Addr:               "127.0.0.1:6379",
		DB:                 15,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: 500 * time.Millisecond,
	}
	_, err := NewRedis(op)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}
