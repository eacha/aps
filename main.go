package main

import (
	//"fmt"
	"github.com/eacha/aps/dns"
	"net"
	"fmt"
	"encoding/hex"
)

func main() {
	//conn, _ := net.Dial("udp", "200.89.70.3:53")
	conn, _ := net.Dial("udp", "62.133.85.107:53")
	query := dns.NewQuery("www.ble.cl", dns.RecursiveDesired)
	conn.Write(query.Pack())
	b := make([]byte, 1024)
	read, _ := conn.Read(b)
	fmt.Println(hex.Dump(b[:read]))
	var response dns.Response
	response.UnPack(b[:read])
	fmt.Println(response.Answer[0])
	fmt.Println(response.Answer[1])
}

//import (
//	"fmt"
//	"log"
//	"os/exec"
//	"sync"
//)
//
//func main() {
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go func (){
//		out, err := exec.Command("dig", "niclabs.cl" ).Output()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Printf("data: %s\n", out)
//		wg.Done()
//	}()
//	wg.Wait()
//}
