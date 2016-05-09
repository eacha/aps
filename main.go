package main

import (
	"encoding/hex"
	"fmt"
	"github.com/eacha/aps/dns"
	"net"
)

func main() {
	conn, _ := net.Dial("udp", "200.89.70.3:53")
	query := dns.NewQuery("www.uchile.cl", dns.RecursiveDesired)
	conn.Write(query.Pack())
	b := make([]byte, 1024)
	read, err := conn.Read(b)
	fmt.Println(err)
	fmt.Println(hex.Dump(b[:read]))
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
