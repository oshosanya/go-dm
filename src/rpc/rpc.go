package rpc

import (
	"context"
	"fmt"
	"github.com/oshosanya/go-dm/protobuf"
	"github.com/oshosanya/go-dm/src/download"
)

type DownloadsServer struct{}

func (s DownloadsServer) Add(ctx context.Context, downloadItem *protobuf.Download) (*protobuf.Download, error) {
	download.AddDownload(downloadItem, 4)
	return downloadItem, nil
}

func (s DownloadsServer) List(ctx context.Context, void *protobuf.Void) (*protobuf.DownloadList, error) {
	return nil, fmt.Errorf("Method not implemented")
}
