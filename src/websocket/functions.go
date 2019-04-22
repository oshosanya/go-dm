package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/oshosanya/go-dm/protobuf"
	"github.com/oshosanya/go-dm/src/logger"
	"google.golang.org/grpc"

	"github.com/oshosanya/go-dm/src/data"
)

var funcMap = map[string]interface{}{
	"getDownloads": getDownloads,
	"addDownload":  addDownload,
}

type NewDownload struct {
	FileUrl string `json:"file_url"`
}

func callDynamically(message Message, args ...interface{}) {
	switch message.Action {
	case "getDownloads":
		funcMap["getDownloads"].(func())()
	case "addDownload":
		println(message.Args)
		funcMap["addDownload"].(func(string))(message.Args)
	}

}

func getDownloads() {
	downloads := data.GetAllDownloads()
	b, err := json.Marshal(downloads)
	if err != nil {
		logger.GetInstance().Printf("Malformed JSON: %e", err)
	}
	response := ResponseMessage{
		Action:  "all",
		Payload: string(b),
	}
	SendMessageToClient(response)
}

func addDownload(args string) {
	log.Println("Adding new download")
	download := NewDownload{}
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to backend: %v", err)
	}
	println(args)
	err = json.Unmarshal([]byte(args), &download)
	if err != nil {
		log.Fatalf("Could not parse download request: %v", err)
	}
	client := protobuf.NewDownloadsClient(conn)
	proto := &protobuf.Download{
		Url:           download.FileUrl,
		ContentLength: int64(0),
		Done:          false,
	}
	_, err = client.Add(context.Background(), proto)
	if err != nil {
		log.Fatalf("Could not add download: %v", err)
	}
	getDownloads()
}

func PublishDownloadProgress(download data.Download) {

}
