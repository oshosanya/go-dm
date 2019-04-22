package download

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"

	"github.com/oshosanya/go-dm/protobuf"
	"github.com/oshosanya/go-dm/src/counter"
	"github.com/oshosanya/go-dm/src/data"
	"github.com/oshosanya/go-dm/src/logger"
	"github.com/oshosanya/go-dm/src/util"
)

//RoutineDefinition Struct defining a segment to be downloaded by goroutine
type RoutineDefinition struct {
	StartRange  int
	EndRange    int
	CurrentSize float32
	FileName    string
}

var log = logger.GetInstance()

func AddDownload(downloadItem *protobuf.Download, numOfThreads int64) *data.Download {
	url := downloadItem.Url
	matched, err := util.ValidateURL(url)

	if matched == false {
		panic(fmt.Sprintf("URL %s is not valid", url))
	}

	url, err = util.BuildURL(url)

	log.Info("Initiating connection to: %s \n", url)
	client := &http.Client{}
	request, err := http.NewRequest("HEAD", url, nil)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	fileName := util.GetFileNameFromURL(url)

	// filePath := strings.Join([]string{DownloadsFolder(), fileName}, "")

	if !(len(response.Header.Get("Content-Length")) > 0) {
		var downloadDef util.RoutineDefinition
		contentLength := "indefinite"
		fmt.Printf("Content-Length is: %s \n", contentLength)
		fileName := util.GetFileNameFromURL(url)
		routineDownloadFileName := util.BuildRoutineDownloadFileName(int64(0), fileName)
		downloadDef = util.BuildRouteDefinition(int64(0), 0, 0, routineDownloadFileName)
		downloadModel := data.CreateDownload(url, fileName, 0)
		data.CreateDownloadRoutine(&downloadModel, downloadDef)
		// DownloadFile(url, fileName)
		return &downloadModel
	}
	contentLength, _ := strconv.Atoi(response.Header.Get("Content-Length"))
	contentLengthPerRoutine := int(math.Ceil(float64(int64(contentLength) / numOfThreads)))
	newStartRange := contentLengthPerRoutine + 1
	downloadModel := data.CreateDownload(url, fileName, contentLength)
	log.Println("Last added download is: ", downloadModel.ID)
	var downloadDef util.RoutineDefinition
	for i := int64(0); i < numOfThreads; i++ {
		if i == 0 {
			routineDownloadFileName := util.BuildRoutineDownloadFileName(i, fileName)
			downloadDef = util.BuildRouteDefinition(i, 0, contentLengthPerRoutine, routineDownloadFileName)
			data.CreateDownloadRoutine(&downloadModel, downloadDef)
		} else if newStartRange+contentLengthPerRoutine > contentLength {
			routineDownloadFileName := util.BuildRoutineDownloadFileName(i, fileName)
			downloadDef = util.BuildRouteDefinition(i, newStartRange, contentLength, routineDownloadFileName)
			data.CreateDownloadRoutine(&downloadModel, downloadDef)
		} else {
			routineDownloadFileName := util.BuildRoutineDownloadFileName(i, fileName)
			downloadDef = util.BuildRouteDefinition(i, newStartRange, newStartRange+contentLengthPerRoutine, routineDownloadFileName)
			data.CreateDownloadRoutine(&downloadModel, downloadDef)
			newStartRange = downloadDef.EndRange + 1
		}
	}
	return &downloadModel
}

func DownloadsFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	downloadsFolder := strings.Join([]string{usr.HomeDir, "/Downloads/"}, "")
	return downloadsFolder
}

func DownloadFile(url string, fileName string) {
	filePath := strings.Join([]string{DownloadsFolder(), fileName}, "")
	fmt.Printf("Saving to File Path: %s \n", filePath)
	out, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	_, err = io.Copy(out, response.Body)
	if err != nil {
		panic(err)
	}
}

func DownloadRoutine(downloadDef data.DownloadRoutine, c chan ChannelReturn, counter *counter.DataTransferred) {
	wg2 := sync.WaitGroup{}

	filePath := strings.Join([]string{DownloadsFolder(), downloadDef.FileName}, "")
	_ = os.Remove(filePath)
	out, err := os.Create(filePath)
	log.Debugf("Created file in path: %s", filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	client := &http.Client{}
	request, err := http.NewRequest("GET", downloadDef.Url, nil)
	if err != nil {
		panic(err)
	}
	bytesRangeSlice := []string{"bytes=", strconv.Itoa(downloadDef.StartRange), "-", strconv.Itoa(downloadDef.EndRange)}
	bytesRange := strings.Join(bytesRangeSlice, "")
	request.Header.Add("Range", bytesRange)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	wg2.Add(1)
	go CopyResponseToFile(downloadDef, response, out, counter, &wg2, c)
	wg2.Wait()
}

//MergeFiles Merge all downloaded segments into one file
func MergeFiles(filePath string, allDownloadDefs []data.DownloadRoutine) {
	_ = os.Remove(filePath)
	out, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	log.Debugf("Merging files into: %s", filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	for _, d := range allDownloadDefs {
		log.Debugf("Merging %s into %s", d.FileName, filePath)
		downloaded, err := ioutil.ReadFile(strings.Join([]string{DownloadsFolder(), d.FileName}, ""))
		if err != nil {
			panic(err)
		} else {
			out.Write(downloaded)
			os.Remove(strings.Join([]string{DownloadsFolder(), d.FileName}, ""))
		}
	}
}

func CopyResponseToFile(downloadDef data.DownloadRoutine, resp *http.Response, out *os.File, counter *counter.DataTransferred, wg2 *sync.WaitGroup, c chan ChannelReturn) {
	defer wg2.Done()
	buf := make([]byte, 16384)
	for {
		bytesRead, err := io.ReadFull(resp.Body, buf)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				// println("End of file reached")
				written, err := out.Write(buf[:bytesRead])
				if err != nil {
					fmt.Printf("Error occurred while writing %s \n", err)
					panic("Die")
				}
				counter.Inc(&downloadDef, written)
				break
			}
			fmt.Printf("Error occurred while copying %s \n", err)
			panic("Die")
		}
		written, err := out.Write(buf)
		if err != nil {
			fmt.Printf("Error occurred while writing %s \n", err)
			panic("Die")
		}
		// wg2.Add(1)
		counter.Inc(&downloadDef, written)
		// wg2.Add(1
		// counter.PrintValue()
	}
	// downloadDef = data.GetDownloadRoutineById(downloadDef.ID)
	downloadDef.Status = "done"
	downloadDef.Update()
	c <- ChannelReturn{
		Status:      "done",
		RoutineData: downloadDef,
	}
}
