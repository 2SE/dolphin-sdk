package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"

	pb2 "github.com/2se/dolphin-sdk/example/client/pb"
	"github.com/2se/dolphin-sdk/pb"
	//"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc"
	"log"
	"time"
)

//直接对grpc-server通过grpc-client的方式完成整个请求
//在没dolphin的情况下测试grpc-server流程可以用此方法
func main() {
	//"192.168.10.169:8848"
	addr := "192.168.1.10:9528"
	oneMethod(addr)
	//addr := "192.168.10.169:8848"
	//parallelMethod(time.Minute, addr, 5000, 50000)
}
func parallelMethod(period time.Duration, address string, numCli, numTurns int) {
	clis := getClients(numCli, address)
	for cli := range clis {
		go func(cli pb.AppServeClient) {
			sendRequest(cli, numTurns)
		}(cli)
	}
	time.Sleep(period)
	return
}
func getClients(num int, address string) chan pb.AppServeClient {
	ctx1, _ := context.WithTimeout(context.Background(), time.Second*5)
	clis := make(chan pb.AppServeClient, 30)
	for i := 0; i < num; i++ {
		go func() {
			conn, err := grpc.DialContext(ctx1, address, grpc.WithBlock(), grpc.WithInsecure())
			if err != nil {
				return
			}
			c := pb.NewAppServeClient(conn)
			clis <- c
		}()
	}
	return clis
}
func sendRequest(cli pb.AppServeClient, num int) {
	for i := 0; i < num; i++ {
		p := &pb2.GetUserRequest{
			UserId: 1,
		}
		object, err := ptypes.MarshalAny(p)
		if err != nil {
			log.Println(err)
			return
		}
		req := &pb.ClientComRequest{
			TraceId: "traceId_2123",
			Id:      "userid123",
			MethodPath: &pb.MethodPath{
				Resource: "MockUser",
				Revision: "v3",
				Action:   "GetUser",
			},
			Params: object,
		}
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		rep, err := cli.Request(ctx, req)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(rep.Code)
	}
}

//单个方法模拟请求
func oneMethod(address string) {
	ctx1, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(ctx1, address, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewAppServeClient(conn)
	p := &pb2.GetUserRequest{
		UserId: 1,
	}
	object, err := ptypes.MarshalAny(p)
	if err != nil {
		log.Println(err)
		return
	}
	//traceId 为客户端生成的随机数
	//methodPath 在启动服务时会在当前目录下生成document.md，这里生成了接口路径和参数名，具体参数需要结合protobuf查看
	req := &pb.ClientComRequest{
		TraceId: "traceId_2123",
		Id:      "userid123",
		MethodPath: &pb.MethodPath{
			Resource: "MockUser",
			Revision: "v2",
			Action:   "GetUser",
		},
		Params: object,
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	res, err := c.Request(ctx, req)
	if err != nil {
		log.Println(err)
		return
	}

	if res.Code == 200 {
		pmu := &pb2.User{}
		err = ptypes.UnmarshalAny(res.Body, pmu)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(pmu)
	} else {
		log.Println(res.Code)
	}
}
