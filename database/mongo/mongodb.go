package mongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

type Config struct {
	Addr string // for trace
	DSN  string // write data source name.

	ConnTimeout  time.Duration
	QueryTimeout time.Duration // query sql timeout
	ExecTimeout  time.Duration // execute sql timeout
}

func NewMongo(c *Config) {
	ctx, _ := context.WithTimeout(context.Background(), c.ConnTimeout)
	cli, err := mongo.Connect(ctx)

}
