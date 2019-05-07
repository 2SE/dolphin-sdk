package main

import (
	"github.com/2se/dolphin-sdk/mock"
	"github.com/2se/dolphin-sdk/server"
	"time"
)

func main() {
	c := &server.Config{
		AppName:      "userApp",
		Address:      "127.0.0.1:8848",
		WriteBufSize: 32 * 1024,
		ReadBufSize:  32 * 1024,
		ConnTimeout:  time.Second * 10,
	}
	//server.Start(c, mock.MkService)
	server.StartGrpcOnly(c, mock.MkService)
}