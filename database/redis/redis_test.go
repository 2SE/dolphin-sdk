package redis

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestNewRedis(t *testing.T) {

	op := &redis.Options{
		Addr:               "192.168.10.189:19000",
		DB:                 15,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: 500 * time.Millisecond,
	}
	cli, err := NewRedis(op)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	key := "TEST:12345"
	val := "this is a test"
	err = cli.Set(key, val, time.Minute).Err()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	val2, err := cli.Get(key).Result()
	if err == redis.Nil {
		t.Error("missing_key does not exist")
		t.Fail()
	} else if err != nil {
		t.Error(err)
		t.Fail()
	} else {
		t.Log("val is :", val2)
	}
}
