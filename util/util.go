package util

import (
	"regexp"
	"strings"
)

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
