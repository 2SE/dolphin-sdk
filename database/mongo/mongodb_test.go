package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestNewMongo(t *testing.T) {
	//db.createUser() https://docs.mongodb.com/manual/reference/method/db.createUser/index.html
	//https://docs.mongodb.com/manual/reference/connection-string/#connections-connection-options
	uri := "mongodb://rennbon1:123456@127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019/?replicaSet=rs0&maxPoolSize=512"
	cli, err := NewMongo(&Config{
		URI: uri,
	})
	if err != nil {
		t.Log("new err:", err)
		return
	}
	fmt.Println("aaa")
	//cli.Connect()
	demoCol := cli.Database("demo").Collection("post")
	//fmt.Println(demoCol.InsertOne(context.Background(), bson.M{"name": "Jack", "age": 24}))
	fmt.Println(demoCol.FindOne(context.Background(), bson.M{"name": "rennbon"}).DecodeBytes())
}
