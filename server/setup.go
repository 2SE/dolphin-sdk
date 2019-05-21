package server

import (
	"context"
	"github.com/2se/dolphin-sdk/log"
	"github.com/2se/dolphin-sdk/pb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"os"
)

//grpc server start
//address: dolphin address  http://www.xxx.com:1111
//services: business service
//需要有dolphin的启动
func Start(c *Config, services ...interface{}) {
	if c.LogCnf != nil {
		log.WithDB(c.LogCnf)

	}
	//设置格式
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	//设置控制台输出
	logrus.SetOutput(os.Stdout)
	//设置落库等级
	logrus.SetLevel(logrus.TraceLevel)
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

//单独启动服务
func StartGrpcOnly(c *Config, services ...interface{}) {
	if c.LogCnf != nil {
		log.WithDB(c.LogCnf)
	}
	//设置格式
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	//设置控制台输出
	logrus.SetOutput(os.Stdout)
	//设置落库等级
	logrus.SetLevel(logrus.TraceLevel)
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

func GetUserId() (string, error) {
	return t.GetUserId()
}
