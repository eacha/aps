package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eacha/aps/scan"
	"github.com/eacha/aps/tools/thread"
	"github.com/eacha/aps/util"
	"github.com/eacha/aps/dns"
	"sync"
	"encoding/json"
)

const (
	inputChannelBuffer  = 1
	outputChannelBuffer = 1
	end                 = 1
)

var (
	modulesList = []string{"DNS"}
	showModules bool
	options     scan.Options
)

func init() {
	flag.BoolVar(&showModules, "list-module", false, "Print module list and exit")

	flag.StringVar(&options.InputFileName, "input-file", "-", "Input file name, use - for stdin")
	flag.StringVar(&options.OutputFileName, "output-file", "-", "Output file name, use - for stdout")
	flag.UintVar(&options.Port, "port", 0, "Port number to scan")
	flag.StringVar(&options.Module, "module", "", "Set module to scan")
	flag.UintVar(&options.Threads, "threads", 1, "Set the number of corutines")
	flag.UintVar(&options.ConnectionTimeout, "connection-timeout", 10, "Set connection timeout in seconds")
	flag.UintVar(&options.IOTimeout, "io-timeout", 10, "Set input output timeout in seconds")

	flag.StringVar(&options.DNSOptions.QuestionURL, "dns-question", "www.uchile.cl", "Question sent to DNS resolver")
	flag.StringVar(&options.DNSOptions.IpResponse, "dns-response", "172.29.1.16", "Expected response of the DNS resolver")

	flag.Parse()

	// Help arguments
	if showModules {
		printModules()
	}

	// Check the arguments
	if options.Port > 65535 {
		log.Fatal("--port must be in the range [0, 65535]")
	}

	if options.Module == "" || !util.StringInSlice(options.Module, modulesList) {
		log.Fatal("--module must be in the --module-list")
	}

	options.InputChan = make(chan string, inputChannelBuffer)
	options.OutputChan = make(chan string, outputChannelBuffer)
}

func printModules() {
	fmt.Println("Modules:")
	for _, mod := range modulesList {
		fmt.Printf("\t- %s\n", mod)
	}
	os.Exit(0)
}

func main() {
	var (
		wg     sync.WaitGroup
		endWrite = make(chan int)
		ts     = make([]*thread.Statistic, int(options.Threads))
	)
	wg.Add(int(options.Threads))
	options.WaitGroup = &wg

	go thread.ReadChannel(options.InputFileName, options.InputChan)
	go thread.WriteChannel(options.OutputFileName, options.OutputChan, endWrite)

	switch options.Module {
	case "DNS":
		for i := 0; i < int(options.Threads); i++ {
			ts[i] = thread.NewStatistic(i)
			dns.Scan(&options, ts[i])
		}
	default:
		log.Fatal("")
	}

	wg.Wait()
	endWrite <- end

	for _, value := range ts {
		j, _ := json.Marshal(*value)
		fmt.Println(string(j))
	}
}
