package server

import (
	"context"
	"fmt"
	"github.com/2se/dolphin-sdk/log"
	"github.com/2se/dolphin-sdk/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

//Service config,如果需要更详细的配置，可以加
type Config struct {
	AppName         string //服务名称，所有服务名称不可重复
	DolphinHttpAddr string //dolphin注册服务的端口
	DolphinGrpcAddr string //dolphin Grpc调度的端口，用于grpc服务之间的互相调用
	Address         string //grpc服务启动监听端口
	WriteBufSize    int    //grpc 写容量控制
	ReadBufSize     int    //grpc 读容量控制
	ConnTimeout     time.Duration
	RequestTimeout  time.Duration //请求时间跨度限制
	LogCnf          *log.Config
	LogLevel        logrus.Level
}

var base = new(baseService)

//服务主体
type baseService struct {
	listen net.Listener
	svc    *grpc.Server
	ready  bool
}

//基础请求
func (b *baseService) Request(ctx context.Context, req *pb.ClientComRequest) (*pb.ServerComResponse, error) {
	logrus.WithFields(logrus.Fields{
		"resource": req.MethodPath.Resource,
		"version":  req.MethodPath.Revision,
		"action":   req.MethodPath.Action,
		"traceId":  req.TraceId,
	}).Trace(req)
	response := delegate.invoke(req)
	if response == nil {
		response = &pb.ServerComResponse{
			Code: 500,
			Text: panicStr,
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"resource": req.MethodPath.Resource,
			"version":  req.MethodPath.Revision,
			"action":   req.MethodPath.Action,
			"traceId":  req.TraceId,
		}).Trace(response)
	}
	return response, nil
}
func (b *baseService) readyGo() {
	b.ready = true
}
func (b *baseService) run(c *Config) {
	if !b.ready {
		panic("The service is not ready,please register your business services first.")
	}
	l, err := net.Listen("tcp", c.Address)
	if err != nil {
		panic(fmt.Errorf("tpc listen err:%v ", err))
	}
	defer l.Close()
	b.listen = l
	b.svc = grpc.NewServer(
		grpc.ConnectionTimeout(c.ConnTimeout),
		grpc.WriteBufferSize(c.WriteBufSize),
		grpc.ReadBufferSize(c.ReadBufSize))
	pb.RegisterAppServeServer(b.svc, b)
	logrus.Info("Grpc server start and listen on ", c.Address)
	if err := b.svc.Serve(l); err != nil {
		panic(fmt.Errorf("failed to serve: %v", err))
	}
}
func (b *baseService) stop() {
	if b.svc != nil {
		b.svc.Stop()
		b.listen.Close()
	}
}
