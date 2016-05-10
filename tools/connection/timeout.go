package connection

import (
	"net"
	"strconv"
	"time"
)

const (
	UDP       = "udp"
	TCP       = "tcp"
	SEPARATOR = ":"
)

type ConnTimeout struct {
	conn              net.Conn
	Protocol          string
	Address           string
	Port              int
	connectionTimeout time.Duration
	ioTimeout         time.Duration
}

func NewConnTimeout(protocol, address string, port int, connectionTimeout, ioTimeout time.Duration) (*ConnTimeout, error) {
	var c ConnTimeout

	c.Address = address
	c.Port = port
	c.connectionTimeout = connectionTimeout * time.Second
	c.ioTimeout = ioTimeout * time.Second

	conn, err := net.DialTimeout(protocol, c.Address+SEPARATOR+strconv.Itoa(c.Port), c.connectionTimeout)
	if errTime, ok := err.(net.Error); ok && errTime.Timeout() {
		return nil, &ConnError{ConnTimeoutMsg, c.Address}
	} else if err != nil {
		return nil, &ConnError{ConnRefusedMsg, c.Address}
	}
	c.conn = conn

	if err = c.conn.SetDeadline(time.Now().Add(c.ioTimeout)); err != nil {
		c.conn.Close()
		return nil, &IOTimeoutError{SetTimeoutMsg, c.Address}
	}

	return &c, nil
}

func (c *ConnTimeout) Read(b []byte) (int, error) {
	readBytes, err := c.conn.Read(b)

	if errTime, ok := err.(net.Error); ok && errTime.Timeout() {
		return 0, &IOTimeoutError{ReadTimeoutMsg, c.Address}
	} else if err != nil {
		return 0, &IOError{ReadMsg, c.Address}
	}

	if err = c.conn.SetDeadline(time.Now().Add(c.ioTimeout)); err != nil {
		return 0, &IOTimeoutError{SetTimeoutMsg, c.Address}
	}

	return readBytes, nil
}

func (c *ConnTimeout) Write(b []byte) (int, error) {
	writeBytes, err := c.conn.Write(b)

	if errTime, ok := err.(net.Error); ok && errTime.Timeout() {
		return 0, &IOTimeoutError{WriteTimeoutMsg, c.Address}
	} else if err != nil {
		return 0, &IOError{WriteMsg, c.Address}
	}

	if err = c.conn.SetDeadline(time.Now().Add(c.ioTimeout)); err != nil {
		return 0, &IOTimeoutError{SetTimeoutMsg, c.Address}
	}

	return writeBytes, nil
}

func (c *ConnTimeout) Close() error {
	return c.conn.Close()
}
