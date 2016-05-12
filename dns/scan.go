package dns

import (
	"github.com/eacha/aps/scan"
	"github.com/eacha/aps/tools/thread"
	"fmt"
)

func Scan(options *scan.Options, statistic *thread.Statistic) {
	defer options.WaitGroup.Done()
	for {
		address, more := <- options.InputChan
		if !more {
			break
		}
		statistic.IncreaseProcessedLines()

		fmt.Println(address)
	}
	statistic.SetEndTime()
}
