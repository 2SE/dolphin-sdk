package mongo

import (
	"context"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	//URI fomat https://docs.mongodb.com/manual/reference/connection-string/
	//https://docs.mongodb.com/manual/reference/connection-string/#connections-connection-options
	URI string
}

func NewMongo(c *Config) (*mgo.Client, error) {
	ctx := context.Background()
	cli, err := mgo.Connect(ctx, options.Client().ApplyURI(c.URI))
	if err != nil {
		return nil, err
	}
	ctx, _ = context.WithTimeout(ctx, time.Second*5)
	err = cli.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
