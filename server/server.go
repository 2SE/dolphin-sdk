package server

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.2se.com/hashhash/server-sdk/pb"
	"google.golang.org/grpc"
	"net"
	"time"
)

//Service config,如果需要更详细的配置，可以加
type Config struct {
	AppName      string
	Address      string
	WriteBufSize int
	ReadBufSize  int
	ConnTimeout  time.Duration
}

var base = new(baseService)

type method func(proto.Message) (proto.Message, error)

//服务主体
type baseService struct {
	methods map[string]method
}

//基础请求
func (b *baseService) Request(ctx context.Context, req *pb.ClientComRequest) (*pb.ServerComResponse, error) {
	response := delegate.invoke(req)
	return response, nil
}
func (b *baseService) registerMP(f method) {

}

func (b *baseService) run(c *Config) {
	l, err := net.Listen("tcp", c.Address)
	if err != nil {
		panic(fmt.Errorf("tpc listen err:%v ", err))
	}
	defer l.Close()
	svc := grpc.NewServer(
		grpc.ConnectionTimeout(c.ConnTimeout),
		grpc.WriteBufferSize(c.WriteBufSize),
		grpc.ReadBufferSize(c.ReadBufSize))
	pb.RegisterAppServeServer(svc, b)
	if err := svc.Serve(l); err != nil {
		panic(fmt.Errorf("failed to serve: %v", err))
	}
}
