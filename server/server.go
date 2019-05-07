package server

import (
	"context"
	"fmt"
	"github.com/2se/dolphin-sdk/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

//Service config,如果需要更详细的配置，可以加
type Config struct {
	AppName      string
	DolphinAddr  string
	Address      string //grpc addr
	WriteBufSize int
	ReadBufSize  int
	ConnTimeout  time.Duration
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
	response := delegate.invoke(req)
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
