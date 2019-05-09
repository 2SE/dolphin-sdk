package server

import (
	"context"
	"fmt"
	"github.com/2se/dolphin-sdk/pb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

//grpc server start
//address: dolphin address  http://www.xxx.com:1111
//services: business service
func Start(c *Config, services ...interface{}) {
	newDolphinClient(c.DolphinGrpcAddr, c.RequestTimeout)
	registerManager.SetAppName(c.AppName)
	registerManager.SetAddress(c.Address)
	registerManager.SetTitle(c.AppName)
	err := parseServices(services...)
	if err != nil {
		panic(err)
	}
	go base.run(c)
	err = registerManager.RegisterServerOnDolpin(c.DolphinHttpAddr)
	if err != nil {
		panic(err)
	}
	select {}
}

//启动服务
func StartGrpcOnly(c *Config, services ...interface{}) {
	registerManager.SetAppName(c.AppName)
	registerManager.SetAddress(c.Address)
	registerManager.SetTitle(c.AppName)
	err := parseServices(services...)
	if err != nil {
		panic(err)
	}
	base.run(c)
}

//发送对其他GRPC服务的调用请求
func SendGrpcRequest(path *pb.MethodPath, message proto.Message) (*pb.ServerComResponse, error) {
	fmt.Println("TraceId ", t.GetTrace())
	object, err := ptypes.MarshalAny(message)
	if err != nil {
		return nil, err
	}
	req := &pb.ClientComRequest{
		TraceId:    t.GetTrace(),
		MethodPath: path,
		Params:     object,
	}
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	return dolphinClient.Request(ctx, req)
}
