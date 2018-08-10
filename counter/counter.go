package counter

import (
	"fmt"
	"sync"

	humanize "github.com/dustin/go-humanize"
)

type DataTransferred struct {
	Count      int
	TotalCount int
	Mux        sync.Mutex
}

func (dt *DataTransferred) Inc(byteCount int) {
	dt.Mux.Lock()
	defer dt.Mux.Unlock()
	dt.Count += byteCount
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
