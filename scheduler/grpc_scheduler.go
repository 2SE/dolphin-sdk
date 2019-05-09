package scheduler

import (
	"context"
	"fmt"
	"github.com/2se/dolphin-sdk/pb"
	"github.com/2se/dolphin-sdk/server"

	"google.golang.org/grpc"
	"log"
	"time"
)

var dolphinClient pb.AppServeClient

func Start(dolphinAddr string) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	conn, err := grpc.DialContext(ctx, dolphinAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect: %v", err)
		return
	}
	defer conn.Close()
	dolphinClient = pb.NewAppServeClient(conn)
}

func SendRequest(in *pb.ClientComRequest) (*pb.ServerComResponse, error) {
	fmt.Println("traceId=>", server.GetTrace())
	//ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	//return dolphinClient.Request(ctx, in)
	return nil, nil
}
