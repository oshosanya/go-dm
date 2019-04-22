package counter

import (
	"fmt"
	"sync"

	humanize "github.com/dustin/go-humanize"
	"github.com/oshosanya/go-dm/src/data"
	"github.com/oshosanya/go-dm/src/websocket"
)

type DataTransferred struct {
	Count      int
	TotalCount int
	Mux        sync.Mutex
}

//Increment the bytes downloaded column of the download
func (dt *DataTransferred) Inc(downloadDef *data.DownloadRoutine, byteCount int) {
	dt.Mux.Lock()
	defer dt.Mux.Unlock()
	dt.Count += byteCount
	currentByteCount := downloadDef.BytesDownloaded
	downloadDef.BytesDownloaded = currentByteCount + byteCount
	if downloadDef.BytesDownloaded == downloadDef.EndRange {
		downloadDef.Done = true
	}
	downloadDef.Update()

	// fmt.Println("Parent bytes download is ", currentParentBytesCount)
	// fmt.Printf("Parent id is %s \n", downloadDef.DownloadID)
	download := data.GetDownloadByID(downloadDef.DownloadID)
	currentParentBytesCount := download.BytesDownloaded
	download.BytesDownloaded = currentParentBytesCount + byteCount
	download.Update()
	// websocket.PublishDownloadProgress(download)
	websocket.SendMessageToClient(download)
	// downloadDef.Download.Update()
	// downloadDef.Update()
	// getDefAgain := data.GetDownloadRoutineById(downloadDef.ID)
	// newBytesTransferred := getDefAgain.BytesDownloaded + byteCount
	// getDefAgain.BytesDownloaded = newBytesTransferred
	// getDefAgain.Update()
}

func (dt *DataTransferred) PrintValue() {
	dt.Mux.Lock()
	defer dt.Mux.Unlock()
	print("\033[1A")
	print("\033[K")
	print("\033[1A")
	print("\033[K")
	fmt.Printf("Total amount of data transferred: %s \n", humanize.Bytes(uint64(dt.Count)))
	percentageTransferred := int((dt.Count * 100) / dt.TotalCount)
	var downloadBarLength int
	if percentageTransferred == 99 {
		downloadBarLength = 10
	} else {
		downloadBarLength = int(percentageTransferred / 10)
	}
	// fmt.Println(percentageTransferred)
	fmt.Print("Progress:  [")
	for i := 0; i < downloadBarLength; i++ {
		fmt.Print("====")
	}
	fmt.Print(">")
	barsLeft := 10 - downloadBarLength
	for i := 0; i < barsLeft; i++ {
		fmt.Print("    ")
	}
	fmt.Print("]\n")
}
