package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type downloadRoutineDefinition struct {
	startRange  int
	endRange    int
	currentSize float32
	fileName    string
}

func main() {
	// argsWithoutProgName := os.Args[1:]
	url := "https://images.pexels.com/photos/33109/fall-autumn-red-season.jpg?cs=srgb&dl=autumn-colorful-colourful-33109.jpg&fm=jpg"
	numberOfRoutines := 4
	matched, err := validateURL(url)
	if matched == false {
		panic("URL not valid")
	}
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		appended := []string{"http://", url}
		url = strings.Join(appended, "")
		fmt.Printf("New Url: %s \n", url)
	}
	fmt.Printf("Initiating connection to: %s \n", url)
	client := &http.Client{}
	request, err := http.NewRequest("HEAD", url, nil)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	fileName := getFileNameFromURL(url)
	extension, _ := mime.ExtensionsByType(response.Header.Get("Content-Type"))
	if len(response.Header.Get("Content-Length")) > 0 {
		contentLength, _ := strconv.Atoi(response.Header.Get("Content-Length"))
		fmt.Println("Content Lenth is : %d", contentLength)
		contentLengthPerRoutine := int(math.Ceil(float64(contentLength / numberOfRoutines)))
		fmt.Println("Content Lenth per routine : %d", contentLengthPerRoutine)
		newStartRange := contentLengthPerRoutine + 1
		var allDownloadDefs []downloadRoutineDefinition
		var downloadDef downloadRoutineDefinition
		for i := 0; i < numberOfRoutines; i++ {
			if i == 0 {
				downloadDef = downloadRoutineDefinition{
					startRange:  0,
					endRange:    contentLengthPerRoutine,
					currentSize: 0,
					fileName:    strings.Join([]string{fileName, strconv.Itoa(i)}, ""),
				}
			} else if newStartRange+contentLengthPerRoutine > contentLength {
				downloadDef = downloadRoutineDefinition{
					startRange:  newStartRange,
					endRange:    contentLength,
					currentSize: 0,
					fileName:    strings.Join([]string{fileName, strconv.Itoa(i)}, ""),
				}
			} else {
				downloadDef = downloadRoutineDefinition{
					startRange:  newStartRange,
					endRange:    newStartRange + contentLengthPerRoutine,
					currentSize: 0,
					fileName:    strings.Join([]string{fileName, strconv.Itoa(i)}, ""),
				}
				newStartRange = downloadDef.endRange + 1
			}
			allDownloadDefs = append(allDownloadDefs, downloadDef)
			wg.Add(1)
			go downloadFileRoutine(url, downloadDef, extension[0])
		}
		wg.Wait()
		fmt.Println(allDownloadDefs)
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		downloadsFolder := strings.Join([]string{usr.HomeDir, "/Downloads/"}, "")
		filePath := strings.Join([]string{downloadsFolder, fileName, extension[0]}, "")
		fmt.Printf("File Path: %s", filePath)
		out, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		fmt.Printf("Arranging file into: %s", filePath)
		if err != nil {
			panic(err)
		}
		defer out.Close()
		for _, d := range allDownloadDefs {
			downloaded, err := ioutil.ReadFile(strings.Join([]string{downloadsFolder, d.fileName, extension[0]}, ""))
			if err != nil {
				panic(err)
			}
			out.Write(downloaded)
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

func getFileNameFromURL(url string) string {
	splitedURL := strings.Split(url, "/")
	if len(splitedURL) == 3 {
		fileName := "index"
		return fileName
	} else {
		splitIfQueryString := strings.Split(splitedURL[4], "?")
		return splitIfQueryString[0]
	}
}

func validateURL(url string) (bool, error) {
	matched, err := regexp.MatchString(`^((http)s{0,1}(:\/\/){1}){0,1}[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`, url)
	return matched, err
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

func downloadFileRoutine(url string, downloadDef downloadRoutineDefinition, extension string) {
	fmt.Println(downloadDef)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	downloadsFolder := strings.Join([]string{usr.HomeDir, "/Downloads/"}, "")
	filePath := strings.Join([]string{downloadsFolder, downloadDef.fileName, extension}, "")
	fmt.Printf("File Path: %s", filePath)
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
	bytesRangeSlice := []string{"bytes=", strconv.Itoa(downloadDef.startRange), "-", strconv.Itoa(downloadDef.endRange)}
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
