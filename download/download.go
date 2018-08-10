package download

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"

	"github.com/oshosanya/go-dm/counter"
)

//RoutineDefinition Struct defining a segment to be downloaded by goroutine
type RoutineDefinition struct {
	StartRange  int
	EndRange    int
	CurrentSize float32
	FileName    string
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

func DownloadRoutine(url string, downloadDef RoutineDefinition, wg *sync.WaitGroup, counter *counter.DataTransferred) {
	defer wg.Done()
	wg2 := sync.WaitGroup{}

	filePath := strings.Join([]string{DownloadsFolder(), downloadDef.FileName}, "")
	out, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
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
	go CopyResponseToFile(response, out, counter, &wg2)
	// wg2.Wait()
	// _, err = io.Copy(out, response.Body)
	wg2.Wait()
}

//MergeFiles Merge all downloaded segments into one file
func MergeFiles(filePath string, allDownloadDefs []RoutineDefinition) {
	out, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	fmt.Printf("Merging files into: %s", filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	for _, d := range allDownloadDefs {
		downloaded, err := ioutil.ReadFile(strings.Join([]string{DownloadsFolder(), d.FileName}, ""))
		if err != nil {
			panic(err)
		}
		out.Write(downloaded)
	}
}

func CopyResponseToFile(resp *http.Response, out *os.File, counter *counter.DataTransferred, wg2 *sync.WaitGroup) {
	defer wg2.Done()
	buf := make([]byte, 1024)
	for {
		_, err := io.ReadFull(resp.Body, buf)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				// println("End of file reached")
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
		counter.Inc(written)
		// wg2.Add(1)
		counter.PrintValue()
	}
}
