package main

import (
	"encoding/base64"
	"fmt"
	"github.com/eacha/aps/test"
	"net"
	"sync"
)

func main() {
	buffer := make([]byte, 1024)
	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", ":3306")
		if err != nil {
			fmt.Println(err)
		}
		defer conn.Close()

		n, _ := conn.Read(buffer)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(buffer[:n]))
	}()

	wg.Add(1)

	banner, _ := base64.StdEncoding.DecodeString("RQAAAP9qBEhvc3QgJzM1LjIuMTIwLjEzOCcgaXMgbm90IGFsbG93ZWQgdG8gY29ubmVjdCB0byB0aGlzIE15U1FMIHNlcnZlcg==")
	server := test.TestingBasicServer{Port: 3306, ToWrite: banner, WriteWait: 2}
	(&server).RunServer()

	wg.Wait()
	fmt.Println("hola")
}
