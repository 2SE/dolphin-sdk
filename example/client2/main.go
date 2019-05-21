package main

import (
	"context"
	"github.com/2se/dolphin-sdk/pb"

	"google.golang.org/grpc"
	"log"
	"time"
)

//直接对grpc-server通过grpc-client的方式完成整个请求
//在没dolphin的情况下测试grpc-server流程可以用此方法
func main() {
	//resource: MockUser
	//action: GetUser version:v2
	//input param:GetUserRequest
	//output param:User

	ctx1, _ := context.WithTimeout(context.Background(), time.Second*5)
	//defer cel()
	//conn, err := grpc.DialContext(ctx1, address, grpc.WithBlock(), grpc.WithInsecure())
	conn, err := grpc.DialContext(ctx1, "192.168.10.169:8848", grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewAppServeClient(conn)
	/*p := &pb2.GetUserRequest{
		UserId: 1,
	}*/

	/*p := &pb2.FindUserByTelReq{
		Tel: "13111111111",
	}*/
	/*object, err := ptypes.MarshalAny(p)
	if err != nil {
		log.Println(err)
		return
	}*/
	//traceId 为客户端生成的随机数
	//methodPath 在启动服务时会在当前目录下生成document.md，这里生成了接口路径和参数名，具体参数需要结合protobuf查看
	req := &pb.ClientComRequest{
		TraceId: "traceId_2123",
		Id:      "userid123",
		MethodPath: &pb.MethodPath{
			Resource: "MockUser", //"MockUser",
			Revision: "v1",       //"v2",
			Action:   "NotParam", //"GetUser",  //"FindUserByTel",
		},
		//Params: object,
	}
	req2 := &pb.ClientComRequest{
		TraceId: "traceId_2123",
		Id:      "userid321",
		MethodPath: &pb.MethodPath{
			Resource: "MockUser", //"MockUser",
			Revision: "v1",       //"v2",
			Action:   "NotParam", //"GetUser",  //"FindUserByTel",
		},
		//Params: object,
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	//res, err := c.Request(ctx, req)
	go c.Request(ctx, req2)
	go c.Request(ctx, req)
	if err != nil {
		log.Println(err)
		return
	}

	/*for {
		fmt.Println(conn.GetState())
	}*/
	//if res.Code == 200 {
	//pmu := &pb2.User{}
	/*pmu := &pb2.FindUserByTelRes{}
	err = ptypes.UnmarshalAny(res.Body, pmu)
	if err != nil {
		log.Println(err)
		return
	}*/
	//log.Println(pmu)
	log.Println("success")
	//}
	time.Sleep(time.Second * 2)
}