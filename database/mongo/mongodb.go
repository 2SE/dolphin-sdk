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
	//connTimeout in connstring does not work
	URI         string
	ConnTimeout time.Duration
}

func NewMongo(c *Config) (*mgo.Client, error) {
	ctx := context.Background()
	cli, err := mgo.Connect(ctx, options.Client().ApplyURI(c.URI))
	if err != nil {
		return nil, err
	}
	ctx, _ = context.WithTimeout(ctx, c.ConnTimeout)
	err = cli.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
