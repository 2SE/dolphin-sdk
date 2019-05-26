package main

import (
	"github.com/2se/dolphin-sdk/mock"
	"github.com/2se/dolphin-sdk/server"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

func main() {
	c := &server.Config{
		AppName:         "userApp",
		Address:         "192.168.1.10:8848",
		WriteBufSize:    32 * 1024,
		ReadBufSize:     32 * 1024,
		ConnTimeout:     time.Second * 10,
		DolphinHttpAddr: "http://192.168.1.10:9527",
		DolphinGrpcAddr: "192.168.1.10:9528",
		RequestTimeout:  time.Second * 30,
	}

	//启动并注册到dolphin
	//1. 启动dolphin
	//2. 启动server
	//3. 启动client
	server.Start(c, mock.MkService)
	//只启动grpc
	//1. 启动server
	//2. 启动client
	//go runPprof(time.Minute, "pprof")
	//server.StartGrpcOnly(c, mock.MkService)
}

func runPprof(period time.Duration, pprofFile string) {

	log.Infof("pprof enabled. and it is path: %s", pprofFile)
	var err error

	cpuf, err := os.Create(pprofFile + ".cpu")
	if err != nil {
		log.Fatal("Failed to create CPU pprof file: ", err)
	}
	defer cpuf.Close()

	memf, err := os.Create(pprofFile + ".mem")
	if err != nil {
		log.Fatal("Failed to create Mem pprof file: ", err)
	}
	defer memf.Close()
	tracef, err := os.Create(pprofFile + ".trace")
	if err != nil {
		log.Fatal("Failed to create trace pprof file: ", err)
	}
	defer tracef.Close()

	pprof.StartCPUProfile(cpuf)
	trace.Start(tracef)
	defer pprof.StopCPUProfile()
	defer pprof.WriteHeapProfile(memf)
	defer trace.Stop()
	time.Sleep(period)
	log.Infof("Profiling info saved to '%s.(cpu|mem)'", pprofFile)
}
