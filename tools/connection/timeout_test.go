package connection

import (
	"encoding/base64"
	"sync"
	"testing"

	"github.com/eacha/aps/tools/test"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ConnTimeoutSuite struct{}

var _ = Suite(&ConnTimeoutSuite{})

var (
	BUFFER = "JwAAAAo0L"
)

func (s *ConnTimeoutSuite) TestNewConnectionRefuse(c *C) {
	_, err := NewConnTimeout(TCP, "", 1, 10, 10)

	c.Assert(err, DeepEquals, &ConnError{ConnRefusedMsg, ""})
}

func (s *ConnTimeoutSuite) TestReadError(c *C) {
	var wc sync.WaitGroup
	port := test.RandomPort()
	sendData, _ := base64.StdEncoding.DecodeString(BUFFER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		buffer := make([]byte, 10)
		conn, _ := NewConnTimeout(TCP, "", port, 10, 10)

		conn.Close()
		_, err := conn.Read(buffer)

		c.Assert(err, DeepEquals, &IOError{ReadMsg, ""})

	}()

	// Server
	server := test.TestingBasicServer{Port: port, ToWrite: sendData, WriteWait: 0}
	(&server).RunServer()

	wc.Wait()
}

func (s *ConnTimeoutSuite) TestWriteError(c *C) {
	var wc sync.WaitGroup
	port := test.RandomPort()
	banner, _ := base64.StdEncoding.DecodeString(BUFFER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		conn, _ := NewConnTimeout(TCP, "", port, 10, 10)

		conn.Close()
		_, err := conn.Write(banner)

		c.Assert(err, DeepEquals, &IOError{WriteMsg, ""})

	}()

	// Server
	server := test.TestingBasicServer{Port: port, ToWrite: banner, WriteWait: 0}
	(&server).RunServer()

	wc.Wait()
}

func (s *ConnTimeoutSuite) TestReadTimeout(c *C) {
	var wc sync.WaitGroup
	port := test.RandomPort()
	banner, _ := base64.StdEncoding.DecodeString(BUFFER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		buffer := make([]byte, 10)
		conn, _ := NewConnTimeout(TCP, "", port, 1, 1)

		defer conn.Close()

		_, err := conn.Read(buffer)

		c.Assert(err, DeepEquals, &IOTimeoutError{ReadTimeoutMsg, ""})
	}()

	// Server
	server := test.TestingBasicServer{Port: port, ToWrite: banner, WriteWait: 2}
	(&server).RunServer()

	wc.Wait()
}

func (s *ConnTimeoutSuite) TestReadSuccess(c *C) {
	var wc sync.WaitGroup
	port := test.RandomPort()
	banner, _ := base64.StdEncoding.DecodeString(BUFFER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		buf := make([]byte, 10)

		conn, _ := NewConnTimeout(TCP, "", port, 10, 10)
		defer conn.Close()

		read, _ := conn.Read(buf)

		c.Assert(buf[:read], DeepEquals, banner)
	}()

	// Server
	server := test.TestingBasicServer{Port: port, ToWrite: banner, WriteWait: 0}
	(&server).RunServer()

	wc.Wait()
}

// Errors
func (s *ConnTimeoutSuite) TestConnError(c *C) {
	c.Assert((&ConnError{"test1", "test2"}).Error(), Equals, "test1, Host: test2")
	c.Assert((&IOError{"test1", "test2"}).Error(), Equals, "test1, Host: test2")
	c.Assert((&IOTimeoutError{"test1", "test2"}).Error(), Equals, "test1, Host: test2")
}
