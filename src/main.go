package main

import (
	"github.com/oshosanya/go-dm/src/rpc"
	"net"

	"github.com/oshosanya/go-dm/src/websocket"

	"github.com/oshosanya/go-dm/protobuf"
	"github.com/oshosanya/go-dm/src/download"
	"github.com/oshosanya/go-dm/src/logger"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	println("Creating grpc interface")
	var downloads rpc.DownloadsServer
	protobuf.RegisterDownloadsServer(server, downloads)
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		logger.GetInstance().Fatalf("Could not bind to address: %v", err)
	}
	go download.RunDownloadPool()
	println("Creating websocket interface")
	go websocket.Start()
	println("Finished init")
	logger.GetInstance().Fatal(server.Serve(l))
}
