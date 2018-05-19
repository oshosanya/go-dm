package main

import (
	"bufio"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/oshosanya/go-dm/download"
	"github.com/oshosanya/go-dm/util"
)

var wg sync.WaitGroup

func main() {
	// argsWithoutProgName := os.Args[1:]
	url := "http://31.210.87.4/ringtones_new/fullmp3low/t/Timaya_feat_Phyno_feat_Olamide_Telli_Person.mp3?get=jjj"
	numberOfRoutines := 4
	matched, err := util.ValidateURL(url)

	if matched == false {
		panic(fmt.Sprintf("URL %s is not valid", url))
	}

	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		appended := []string{"http://", url}
		url = strings.Join(appended, "")
	}
	fmt.Printf("Initiating connection to: %s \n", url)
	client := &http.Client{}
	request, err := http.NewRequest("HEAD", url, nil)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	fileName := util.GetFileNameFromURL(url)
	fmt.Println(fileName)
	if len(response.Header.Get("Content-Length")) > 0 {
		contentLength, _ := strconv.Atoi(response.Header.Get("Content-Length"))
		fmt.Printf("Content Lenth is : %d", contentLength)
		contentLengthPerRoutine := int(math.Ceil(float64(contentLength / numberOfRoutines)))
		fmt.Printf("Content Lenth per routine : %d", contentLengthPerRoutine)
		newStartRange := contentLengthPerRoutine + 1
		var allDownloadDefs []download.RoutineDefinition
		var downloadDef download.RoutineDefinition
		for i := 0; i < numberOfRoutines; i++ {
			if i == 0 {
				downloadDef = download.RoutineDefinition{
					StartRange:  0,
					EndRange:    contentLengthPerRoutine,
					CurrentSize: 0,
					FileName:    strings.Join([]string{fileName, strconv.Itoa(i)}, ""),
				}
			} else if newStartRange+contentLengthPerRoutine > contentLength {
				downloadDef = download.RoutineDefinition{
					StartRange:  newStartRange,
					EndRange:    contentLength,
					CurrentSize: 0,
					FileName:    strings.Join([]string{fileName, strconv.Itoa(i)}, ""),
				}
			} else {
				downloadDef = download.RoutineDefinition{
					StartRange:  newStartRange,
					EndRange:    newStartRange + contentLengthPerRoutine,
					CurrentSize: 0,
					FileName:    strings.Join([]string{fileName, strconv.Itoa(i)}, ""),
				}
				newStartRange = downloadDef.EndRange + 1
			}
			allDownloadDefs = append(allDownloadDefs, downloadDef)
			wg.Add(1)
			go download.DownloadRoutine(url, downloadDef, &wg)
		}
		wg.Wait()

		filePath := strings.Join([]string{download.DownloadsFolder(), fileName}, "")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			download.MergeFiles(filePath, allDownloadDefs)
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("File already exist, do you want to overwrite it? Y or N \n")
			text, _ := reader.ReadString('\n')
			if text == "Y" {
				download.MergeFiles(filePath, allDownloadDefs)
			} else {

			}
		}

	} else {
		contentLength := "indefinite"
		fmt.Printf("Content-Length is: %s", contentLength)
	}

	// if len(contentLength) > 0 {
	// 	fmt.Println("Content length exists")
	// } else {
	// 	fmt.Println("No Content-Length")
	// }
	// contentLengthPerRoutine := math.Ceil(contentLength / numberOfRoutines)
	// fmt.Printf("Content-Length per routine: %s", contentLengthPerRoutine)
	// downloadFile(url)

	// //Get file extension from Content-Type
	// extension, _ := mime.ExtensionsByType(response.Header.Get("Content-Type"))
	// fmt.Println(extension[0])

}

// func downloadFile(url string, extension string) {
// 	usr, err := user.Current()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	downloadsFolder := strings.Join([]string{usr.HomeDir, "/Downloads/"}, "")
// 	filePath := strings.Join([]string{downloadsFolder, getFileNameFromURL(url), extension[0]}, "")
// 	fmt.Printf("File Path: %s", filePath)
// 	out, err := os.Create(filePath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer out.Close()
// 	response, err = http.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer response.Body.Close()
// 	fmt.Println("Writing to file")
// 	_, err = io.Copy(out, response.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// }
