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
	fmt.Println("Writing to file")
	_, err = io.Copy(out, response.Body)
	if err != nil {
		panic(err)
	}
}

func DownloadRoutine(url string, downloadDef RoutineDefinition, wg *sync.WaitGroup) {
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
	fmt.Println("Writing to file")
	_, err = io.Copy(out, response.Body)
	if err != nil {
		panic(err)
	}
	wg.Done()
}

//MergeFiles Merge all downloaded segments into one file
func MergeFiles(filePath string, allDownloadDefs []RoutineDefinition) {
	fmt.Printf("File Path: %s", filePath)
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
