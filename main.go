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
	//buffer := make([]byte, 1024)
	//var wg sync.WaitGroup
	//
	//go func() {
	//	defer wg.Done()
	//	conn, err := net.Dial("tcp", ":3306")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	defer conn.Close()
	//
	//	n, _ := conn.Read(buffer)
	//
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	fmt.Println(string(buffer[:n]))
	//}()
	//
	//wg.Add(1)
	//
	//banner, _ := base64.StdEncoding.DecodeString("RQAAAP9qBEhvc3QgJzM1LjIuMTIwLjEzOCcgaXMgbm90IGFsbG93ZWQgdG8gY29ubmVjdCB0byB0aGlzIE15U1FMIHNlcnZlcg==")
	//server := test.TestingBasicServer{Port: 3306, ToWrite: banner, WriteWait: 2}
	//(&server).RunServer()
	//
	//wg.Wait()
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
