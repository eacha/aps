package test

import (
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"time"
)

type TestingBasicServer struct {
	Port      int
	ToWrite   []byte
	WriteWait time.Duration
}

func (tbs *TestingBasicServer) RunServer() {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(tbs.Port))
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		time.Sleep(tbs.WriteWait * time.Second)
		_, err = conn.Write(tbs.ToWrite)
		if err != nil {
			fmt.Println(err)
		}

		return
	}
}

func ConnectionServer(protocol string, port *int, sleep time.Duration) {
	l, err := net.Listen(protocol, ":0")
	if err != nil {
		fmt.Println(err)
	}
	*port = parsePort(l.Addr().String())
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	time.Sleep(sleep * time.Millisecond)
}

func ServerWrite(port *int, buf []byte, sleep time.Duration) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println(err)
	}
	*port = parsePort(l.Addr().String())
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(sleep * time.Millisecond)
}

func ServerReadAndClose(port *int, sleep time.Duration) {
	buf := make([]byte, 10)
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println(err)
	}
	*port = parsePort(l.Addr().String())
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Read(buf)
	conn.Close()

	time.Sleep(sleep * time.Millisecond)
}

// Return a random port number between 10000 and 60000
func RandomPort() int {
	return rand.Intn(50000) + 10000
}

func parsePort(address string) int {
	re, _ := regexp.Compile(`.*?(\d+)$`)
	match := re.FindStringSubmatch(address)

	port, _ := strconv.Atoi(match[1])
	return port
}
