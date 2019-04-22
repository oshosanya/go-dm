package main

import (
	"context"
	"log"

	// "github.com/oshosanya/go-dm/master/cmd"
	"github.com/oshosanya/go-dm/protobuf"
	"google.golang.org/grpc"
)

func list(ctx context.Context, client protobuf.DownloadsClient) {
	l, err := client.List(ctx, &protobuf.Void{})
	print(err, l)
	// if err != nil {
	// 	print(fmt.Errorf("Could not fetch downloads list: %v", err))
	// }

	// for _, t := range l.DownloadItem {
	// 	fmt.Println(t.Url)
	// }
}

func add(ctx context.Context, client protobuf.DownloadsClient) {
	download := &protobuf.Download{
		Url:           "https://preview.redd.it/b9g7pq71ncw11.jpg?width=960&crop=smart&auto=webp&s=e44d55ca28c0ebb63fe938368c86f5d63634d7c8",
		ContentLength: int64(676673),
		Done:          false,
	}
	i, err := client.Add(ctx, download)
	if err != nil {
		log.Fatalf("Could not add download: %v", err)
	}
	print(i.Url)
}

func main() {
	// cmd.Execute()
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to backend: %v", err)
	}
	client := protobuf.NewDownloadsClient(conn)
	add(context.Background(), client)
}
