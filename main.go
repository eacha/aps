package main

import (
	"flag"
	"github.com/eacha/aps/scan"
	"github.com/eacha/aps/tools/thread"
	"log"
	"fmt"
	"os"
	"github.com/eacha/aps/util"
)

var (
	MODULE_LIST = []string{"DNS"}
	listModule bool
	options scan.ScanOptions
)

func init(){
	flag.BoolVar(&listModule, "list-module", false, "Print module list and exit")

	flag.StringVar(&options.InputFileName, "input-file", "-", "Input file name, use - for stdin")
	flag.StringVar(&options.OutputFileName, "output-file", "-", "Output file name, use - for stdout")
	flag.UintVar(&options.Port, "port", 0, "Port number to scan")
	flag.StringVar(&options.Module, "module", "", "Set module to scan")
	flag.UintVar(&options.Threads, "threads", 1000, "Set the number of corutines")
	flag.UintVar(&options.ConnectionTimeout, "connection-timeout", 10, "Set connection timeout in seconds")
	flag.UintVar(&options.IOTimeout, "io-timeout", 10, "Set input output timeout in seconds")

	flag.StringVar(&options.DNSOptions.QuestionURL, "dns-question", "www.uchile.cl", "Question sent to DNS resolver")
	flag.StringVar(&options.DNSOptions.IpResponse, "dns-response", "172.29.1.16", "Expected response of the DNS resolver")

	flag.Parse()

	// Help arguments
	if listModule {
		printModules()
	}

	// Check the arguments
	if options.Port > 65535 {
		log.Fatal("--port must be in the range [0, 65535]")
	}

	if options.Module == "" || !util.StringInSlice(options.Module, MODULE_LIST) {
		log.Fatal("--module must be in the --module-list")
	}

	options.InputFile = thread.NewSyncRead(options.InputFileName)
	options.OutputFile = thread.NewSyncWrite(options.OutputFileName)
}

func printModules() {
	fmt.Println("Modules:")
	for _, mod := range MODULE_LIST{
		fmt.Printf("\t- %s\n", mod)
	}
	os.Exit(0)
}

func main() {

}







//import (
//	//"fmt"
//	"github.com/eacha/aps/dns"
//	"net"
//	"fmt"
//	"encoding/hex"
//)
//
//func main() {
//	//conn, _ := net.Dial("udp", "200.89.70.3:53")
//	conn, _ := net.Dial("udp", "62.133.85.107:53")
//	query := dns.NewQuery("www.ble.cl", dns.RecursiveDesired)
//	conn.Write(query.Pack())
//	b := make([]byte, 1024)
//	read, _ := conn.Read(b)
//	fmt.Println(hex.Dump(b[:read]))
//	var response dns.Response
//	response.UnPack(b[:read])
//	fmt.Println(response.Answer[0])
//	fmt.Println(response.Answer[1])
//}
