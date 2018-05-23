// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	dl "github.com/oshosanya/go-dm/download"
	"github.com/oshosanya/go-dm/util"
	"github.com/spf13/cobra"
)

var numOfThreads int64
var wg sync.WaitGroup

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download [url]",
	Short: "Download the file from the given url",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.Join(args, ""))
		download(strings.Join(args, ""))
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	downloadCmd.Flags().Int64VarP(&numOfThreads, "threads", "t", 4, "Number of threads to use for connection")
}

func download(url string) {
	//url := "http://31.210.87.4/ringtones_new/fullmp3low/t/Timaya_feat_Phyno_feat_Olamide_Telli_Person.mp3?get=jjj"
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
		contentLengthPerRoutine := int(math.Ceil(float64(int64(contentLength) / numOfThreads)))
		fmt.Printf("Content Lenth per routine : %d", contentLengthPerRoutine)
		newStartRange := contentLengthPerRoutine + 1
		var allDownloadDefs []dl.RoutineDefinition
		var downloadDef dl.RoutineDefinition
		for i := int64(0); i < numOfThreads; i++ {
			if i == 0 {
				downloadDef = dl.RoutineDefinition{
					StartRange:  0,
					EndRange:    contentLengthPerRoutine,
					CurrentSize: 0,
					FileName:    strings.Join([]string{fileName, strconv.Itoa(int(i))}, ""),
				}
			} else if newStartRange+contentLengthPerRoutine > contentLength {
				downloadDef = dl.RoutineDefinition{
					StartRange:  newStartRange,
					EndRange:    contentLength,
					CurrentSize: 0,
					FileName:    strings.Join([]string{fileName, strconv.Itoa(int(i))}, ""),
				}
			} else {
				downloadDef = dl.RoutineDefinition{
					StartRange:  newStartRange,
					EndRange:    newStartRange + contentLengthPerRoutine,
					CurrentSize: 0,
					FileName:    strings.Join([]string{fileName, strconv.Itoa(int(i))}, ""),
				}
				newStartRange = downloadDef.EndRange + 1
			}
			allDownloadDefs = append(allDownloadDefs, downloadDef)
			wg.Add(1)
			go dl.DownloadRoutine(url, downloadDef, &wg)
		}
		wg.Wait()

		filePath := strings.Join([]string{dl.DownloadsFolder(), fileName}, "")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			dl.MergeFiles(filePath, allDownloadDefs)
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("File already exist, do you want to overwrite it? Y or N \n")
			text, _ := reader.ReadString('\n')
			if text == "Y" {
				dl.MergeFiles(filePath, allDownloadDefs)
			} else {

			}
		}

	} else {
		contentLength := "indefinite"
		fmt.Printf("Content-Length is: %s \n", contentLength)
		fileName := util.GetFileNameFromURL(url)
		dl.DownloadFile(url, fileName)
	}
}
