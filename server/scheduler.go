package server

import (
	"context"
	"fmt"
	"github.com/2se/dolphin-sdk/pb"
	"github.com/2se/dolphin-sdk/trace"
	"google.golang.org/grpc"
	"time"
)

var (
	dolphinClient  pb.AppServeClient
	requestTimeout time.Duration
	t              = trace.GetTracer()
)

func newDolphinClient(dolphinAddr string, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	conn, err := grpc.DialContext(ctx, dolphinAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		panic(fmt.Errorf("did not connect dolphin: %v", err))
	}
	dolphinClient = pb.NewAppServeClient(conn)
	requestTimeout = timeout
}
