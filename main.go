package main

import (
	"fmt"
	"gitlab.2se.com/hashhash/server-sdk/mock"
	"gitlab.2se.com/hashhash/server-sdk/server"
	"time"
)

func main() {

	server.RegisterService(mock.MockService)
	fmt.Println("mock service init ok")
	c := &server.Config{
		AppName:      "userApp",
		Address:      "127.0.0.1:8848",
		WriteBufSize: 32 * 1024,
		ReadBufSize:  32 * 1024,
		ConnTimeout:  time.Second * 10,
	}
	server.Run(c)

	//server.RegisterServerOnDolpin("")
}
