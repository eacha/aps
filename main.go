package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"encoding/json"
	"sync"
	"time"

	"github.com/eacha/aps/dns"
	"github.com/eacha/aps/scan"
	"github.com/eacha/aps/tools/conn"
	"github.com/eacha/aps/tools/thread"
	"github.com/eacha/aps/util"
)

const (
	inputChannelBuffer  = 1
	outputChannelBuffer = 1
	end                 = 1
)

var (
	modulesList       = []string{"DNS"}
	protocolList      = []string{conn.UDP, conn.TCP}
	showModules       bool
	showProtocols     bool
	options           scan.Options
	connectionTimeout uint
	ioTimeout         uint
)

func init() {
	flag.BoolVar(&showModules, "module-list", false, "Print module list and exit")
	flag.BoolVar(&showProtocols, "protocol-list", false, "Print protocol list and exit")

	flag.StringVar(&options.InputFileName, "input-file", "-", "Input file name, use - for stdin")
	flag.StringVar(&options.OutputFileName, "output-file", "-", "Output file name, use - for stdout")
	flag.IntVar(&options.Port, "port", 0, "Port number to scan")
	flag.StringVar(&options.Module, "module", "", "Set module to scan")
	flag.StringVar(&options.Protocol, "protocol", conn.TCP, "Set protocol to scan")
	flag.UintVar(&options.Threads, "threads", 1, "Set the number of corutines")
	flag.UintVar(&connectionTimeout, "connection-timeout", 10, "Set connection timeout in seconds")
	flag.UintVar(&ioTimeout, "io-timeout", 10, "Set input output timeout in seconds")

	flag.StringVar(&options.DNSOptions.QuestionURL, "dns-question", "www.uchile.cl", "Question sent to DNS resolver")
	flag.StringVar(&options.DNSOptions.IpResponse, "dns-response", "172.29.1.16", "Expected response of the DNS resolver")

	flag.Parse()

	// Help arguments
	if showModules {
		printModules()
	}

	if showProtocols {
		printProtocols()
	}

	// Check the arguments
	if options.Port < 0 || options.Port > 65535 {
		log.Fatal("--port must be in the range [0, 65535]")
	}

	if options.Module == "" || !util.StringInSlice(options.Module, modulesList) {
		log.Fatal("--module must be in the --module-list")
	}

	if !util.StringInSlice(options.Protocol, protocolList) {
		log.Fatal("--protocol must be in the --protocol-list")
	}

	if connectionTimeout <= 0 && ioTimeout <= 0 {
		log.Fatal("--connection-timeout and  --io-timeout must be positive")
	}

	options.ConnectionTimeout = time.Duration(connectionTimeout)
	options.IOTimeout = time.Duration(ioTimeout)
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

func printProtocols() {
	fmt.Println("protocols:")
	for _, mod := range protocolList {
		fmt.Printf("\t- %s\n", mod)
	}
	os.Exit(0)
}

func main() {
	var (
		wg       sync.WaitGroup
		endWrite = make(chan int)
		ts       = make([]*thread.Statistic, int(options.Threads))
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
