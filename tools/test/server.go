package test

import (
	"fmt"
	"net"
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
