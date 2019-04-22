package util

import (
	"log"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

//RoutineDefinition Struct defining a segment to be downloaded by goroutine
type RoutineDefinition struct {
	StartRange  int
	EndRange    int
	CurrentSize float32
	FileName    string
}

//GetFileNameFromURL Retrieve file name from url
func GetFileNameFromURL(url string) string {
	splitedURL := strings.Split(url, "/")
	if len(splitedURL) == 3 {
		fileName := "index"
		return fileName
	} else {
		splitIfQueryString := strings.Split(splitedURL[len(splitedURL)-1], "?")
		return splitIfQueryString[0]
	}
}

//ValidateURL Check if URL is valid
func ValidateURL(url string) (bool, error) {
	matched, err := regexp.MatchString(`^((http)s{0,1}(:\/\/){1}){0,1}[-a-zA-Z0-9@:%._\+~#=]{2,256}\.([a-z]{2,6}){0,1}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`, url)
	return matched, err
}

func BuildURL(url string) (string, error) {
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		appended := []string{"http://", url}
		url = strings.Join(appended, "")
	}
	return url, nil
}

func BuildRouteDefinition(threadIndex int64, startRange int, endRange int, fileName string) RoutineDefinition {
	routineDef := RoutineDefinition{
		StartRange:  startRange,
		EndRange:    endRange,
		CurrentSize: 0,
		FileName:    strings.Join([]string{fileName, strconv.Itoa(int(threadIndex))}, ""),
	}

	return routineDef
}

func BuildRoutineDownloadFileName(threadIndex int64, fileName string) string {
	return strings.Join([]string{fileName, strconv.Itoa(int(threadIndex))}, "")
}

func DownloadsFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	downloadsFolder := strings.Join([]string{usr.HomeDir, "/Downloads/"}, "")
	return downloadsFolder
}
