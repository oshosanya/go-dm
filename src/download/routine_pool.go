package download

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/oshosanya/go-dm/src/counter"
	"github.com/oshosanya/go-dm/src/data"
	"github.com/oshosanya/go-dm/src/util"
)

var availableRoutines = 2
var runningRoutines []data.DownloadRoutine
var runningDownloads []data.Download
var routinesCache []data.DownloadRoutine
var wg sync.WaitGroup
var mu sync.Mutex

type ChannelReturn struct {
	Status      string
	RoutineData data.DownloadRoutine
}

func RunDownloadPool() {
	counter := counter.DataTransferred{}
	routineChannel := make(chan ChannelReturn, availableRoutines)
	downloadQueue := make(chan data.DownloadRoutine, 5)
	dr := data.GetLimitedRoutines(availableRoutines)
	log.Debug("Pre Populating download queue")
	for _, r := range dr {
		fmt.Printf("Adding %s to download queue \n", r.FileName)
		r.Status = "running"
		r.Update()
		time.Sleep(500 * time.Millisecond)
		downloadQueue <- r
	}

	// Handle finished downloads
	log.Debug("Starting routine to handle completed downloads")
	go watchForDownloads(downloadQueue, routineChannel, counter)
	go addDownloadToQueue(downloadQueue)
	go handleFinishedDownloads(routineChannel, downloadQueue)
}

func watchForDownloads(downloadQueue chan data.DownloadRoutine, routineChannel chan ChannelReturn, counter counter.DataTransferred) {
	for {
		// fmt.Println("Watching for available downloads")
		select {
		case dr := <-downloadQueue:
			fmt.Printf("Download taken from queue %s \n", dr.FileName)
			if availableRoutines < 1 {
				downloadQueue <- dr
				continue
			}
			fmt.Printf("Starting download for routine with filename: %s", dr.FileName)
			go DownloadRoutine(dr, routineChannel, &counter)
			fmt.Printf("Available routines before adding %s is %d \n", dr.FileName, availableRoutines)
			decrementAvailableRoutines()
			fmt.Printf("Available routines after adding %s is %d \n", dr.FileName, availableRoutines)
		default:
			continue
		}
	}
}

func handleFinishedDownloads(c chan ChannelReturn, downloadQueue chan data.DownloadRoutine) {
	for {
		time.Sleep(3 * time.Millisecond)
		select {
		case returned := <-c:
			fmt.Printf("Download routine for %s completed \n", returned.RoutineData.FileName)
			if returned.Status == "done" {
				d := data.GetDownloadByID(returned.RoutineData.DownloadID)
				dr := data.GetAllDownloadRoutines(&d)
				doneStat := false
				for _, r := range dr {
					if r.Status != "done" {
						doneStat = false
						break
					} else if r.Status == "done" {
						doneStat = true
					}
				}
				if doneStat == true {
					log.Debugf("Number of done routines for %s: %d", d.Url, len(dr))
					mergeIfComplete(dr)
				}
			}
			fmt.Printf("Available routines after completing %s is %d \n", returned.RoutineData.FileName, availableRoutines)
			incrementAvailableRoutines()
			fmt.Printf("Available routines after closing %s is %d \n", returned.RoutineData.FileName, availableRoutines)
		default:
			continue
		}
	}
}

func mergeIfComplete(downloadRoutines []data.DownloadRoutine) {

	fileName := util.GetFileNameFromURL(downloadRoutines[0].Url)

	filePath := strings.Join([]string{util.DownloadsFolder(), fileName}, "")
	log.Debugf("Merging downloads for filepath: %s", filePath)
	MergeFiles(filePath, downloadRoutines)
}

func addDownloadToQueue(downloadQueue chan data.DownloadRoutine) {
	for {
		if getAvailableRoutines() == 0 {
			continue
		}
		dr := data.GetLimitedRoutines(availableRoutines)
		for _, r := range dr {
			fmt.Printf("Adding %s to download queue \n", r.FileName)
			r.Status = "running"
			r.Update()
			log.Debugf("Adding routine with ID %d to queue", r.ID)
			downloadQueue <- r
		}
	}
}

func getAvailableRoutines() int {
	mu.Lock()
	defer mu.Unlock()
	return availableRoutines
}

func incrementAvailableRoutines() {
	mu.Lock()
	defer mu.Unlock()
	availableRoutines = availableRoutines + 1
}

func decrementAvailableRoutines() {
	mu.Lock()
	defer mu.Unlock()
	availableRoutines = availableRoutines - 1
}
