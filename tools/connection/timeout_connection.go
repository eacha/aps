package connection

import (
	"net"
	"time"
)

type ConnTimeout struct  {
	conn 		  net.Conn
	address           string
	port              int
	connectionTimeout time.Duration
	ioTimeout         time.Duration
}

//func NewConnTimeout(address string, port int, connectionTimeout, ioTimeout time.Duration) (*ConnTimeout, error) {
//	var c ConnTimeout
//
//}