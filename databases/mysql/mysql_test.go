package mysql

import (
	"encoding/base64"
	. "gopkg.in/check.v1"
	"sync"
	"testing"
	"github.com/eacha/aps/tools/test"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Test(t *testing.T) { TestingT(t) }

type MySQLSuite struct{}

var _ = Suite(&MySQLSuite{})

var (
	ERROR_BANNER   = "RQAAAP9qBEhvc3QgJzM1LjIuMTIwLjEzOCcgaXMgbm90IGFsbG93ZWQgdG8gY29ubmVjdCB0byB0aGlzIE15U1FMIHNlcnZlcg=="
	SUCCESS_BANNER = "JwAAAAo0LjAuMjUAZytSAEZ1PnA7c1IhACwgCAIAAAAAAAAAAAAAAAAAAA=="
	WRONG_BANNER   = "RwAAAAo0LjAuMjUAZytSAEZ1PnA7c1IhACwgCAIAAAAAAAAAAAAAAAAAAA=="
)

func (s *MySQLSuite) TestDecodeHeader(c *C) {
	header := []byte{69, 13, 2, 3}
	length, id := decodeHeader(header)

	c.Assert(length, Equals, uint32(0x20d45))
	c.Assert(id, Equals, uint32(3))
}

func (s *MySQLSuite) TestDecodeBannerError(c *C) {
	var data MySQL
	banner, _ := base64.StdEncoding.DecodeString(ERROR_BANNER)

	(&data).decodeBanner(banner)

	c.Assert(data.Banner, DeepEquals, banner)
	c.Assert(data.MySQLError.Code, Equals, uint16(1130))
	c.Assert(data.MySQLError.Message, Equals, "Host '35.2.120.138' is not allowed to connect to this MySQL server")
}

func (s *MySQLSuite) TestDecodeBannerSuccess(c *C) {
	var data MySQL
	banner, _ := base64.StdEncoding.DecodeString(SUCCESS_BANNER)

	(&data).decodeBanner(banner)

	c.Assert(data.Banner, DeepEquals, banner)
	c.Assert(data.Proto, Equals, byte(10))
	c.Assert(data.Version, Equals, "4.0.25")
	c.Assert(contains(data.Capabilities, "ConnectWithDatabase"), Equals, true)
	c.Assert(contains(data.Capabilities, "SupportsTransactions"), Equals, true)
	c.Assert(contains(data.Capabilities, "SupportsCompression"), Equals, true)
	c.Assert(contains(data.Capabilities, "LongColumnFlag"), Equals, true)
	c.Assert(data.Charset, Equals, "latin1_swedish_ci")
	c.Assert(data.Status, Equals, uint16(2))
}

func (s *MySQLSuite) TestDecodeBannerWrong(c *C) {
	var data MySQL
	banner, _ := base64.StdEncoding.DecodeString(WRONG_BANNER)

	err := (&data).decodeBanner(banner)

	c.Assert(err, Equals, MalFormedPacketError)
}

func (s *MySQLSuite) TestParseCharset(c *C) {
	c.Assert(parseCharset(0x08), Equals, "latin1_swedish_ci")
	c.Assert(parseCharset(0x21), Equals, "utf8_general_ci")
	c.Assert(parseCharset(0x63), Equals, "binary")
	c.Assert(parseCharset(0x50), Equals, "")
}

func (s *MySQLSuite) TestConnectionRefuse(c *C) {
	var data MySQL
	_, err := NewMySQLConnection("", 1, 10, 10, &data)

	c.Assert(err, Equals, ConnectionRefusedError)
}

func (s *MySQLSuite) TestConnectionTimeout(c *C) {
	var data MySQL
	_, err := NewMySQLConnection("10.255.255.1", 1, 3, 10, &data)

	c.Assert(err, Equals, ConnectionTimeoutError)
}

func (s *MySQLSuite) TestReadError(c *C) {
	var wc sync.WaitGroup
	banner, _ := base64.StdEncoding.DecodeString(SUCCESS_BANNER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		buffer := make([]byte, 10)
		data := MySQL{Ip: "127.0.0.1", Port: 12345}
		mysqlConn, err := NewMySQLConnection("", 12345, 10, 10, &data)

		if err != nil {
			c.Fatalf("Connection Error: %s", err.Error())
		}

		mysqlConn.Close()
		_, err = mysqlConn.Read(buffer)

		c.Assert(err, Equals, ReadError)

	}()

	// Server
	server := test.TestingBasicServer{Port: 12345, ToWrite: banner, WriteWait: 0}
	(&server).RunServer()

	wc.Wait()
}

func (s *MySQLSuite) TestWriteError(c *C) {
	var wc sync.WaitGroup
	banner, _ := base64.StdEncoding.DecodeString(SUCCESS_BANNER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		data := MySQL{Ip: "127.0.0.1", Port: 12346}
		mysqlConn, err := NewMySQLConnection("", 12346, 10, 10, &data)

		if err != nil {
			c.Fatalf("Connection Error: %s", err.Error())
		}

		mysqlConn.Close()
		_, err = mysqlConn.Write(banner)

		c.Assert(err, Equals, WriteError)

	}()

	// Server
	server := test.TestingBasicServer{Port: 12346, ToWrite: banner, WriteWait: 0}
	(&server).RunServer()

	wc.Wait()
}

func (s *MySQLSuite) TestReadTimeout(c *C) {
	var wc sync.WaitGroup
	banner, _ := base64.StdEncoding.DecodeString(SUCCESS_BANNER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		buffer := make([]byte, 10)
		data := MySQL{Ip: "127.0.0.1", Port: 12347}
		mysqlConn, err := NewMySQLConnection("", 12347, 1, 1, &data)

		if err != nil {
			c.Fatalf("Connection Error: %s", err.Error())
		}
		defer mysqlConn.Close()

		_, err = mysqlConn.Read(buffer)

		c.Assert(err, Equals, IOTimeoutError)

	}()

	// Server
	server := test.TestingBasicServer{Port: 12347, ToWrite: banner, WriteWait: 2}
	(&server).RunServer()

	wc.Wait()
}

func (s *MySQLSuite) TestRealSuccess(c *C) {
	var wc sync.WaitGroup
	banner, _ := base64.StdEncoding.DecodeString(SUCCESS_BANNER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		data := MySQL{Ip: "127.0.0.1", Port: 12348}
		mysqlConn, err := NewMySQLConnection("", 12348, 10, 10, &data)

		if err != nil {
			c.Fatalf("Connection Error: %s", err.Error())
		}
		defer mysqlConn.Close()
		err = mysqlConn.GetBanner()

		c.Assert(data.Banner, DeepEquals, banner)
		c.Assert(data.Proto, Equals, byte(10))
		c.Assert(data.Version, Equals, "4.0.25")
		c.Assert(contains(data.Capabilities, "ConnectWithDatabase"), Equals, true)
		c.Assert(contains(data.Capabilities, "SupportsTransactions"), Equals, true)
		c.Assert(contains(data.Capabilities, "SupportsCompression"), Equals, true)
		c.Assert(contains(data.Capabilities, "LongColumnFlag"), Equals, true)
		c.Assert(data.Charset, Equals, "latin1_swedish_ci")
		c.Assert(data.Status, Equals, uint16(2))
	}()

	// Server
	server := test.TestingBasicServer{Port: 12348, ToWrite: banner, WriteWait: 0}
	(&server).RunServer()

	wc.Wait()
}

func (s *MySQLSuite) TestRealError(c *C) {
	var wc sync.WaitGroup
	banner, _ := base64.StdEncoding.DecodeString(ERROR_BANNER)

	wc.Add(1)
	go func() { // Client
		defer wc.Done()
		data := MySQL{Ip: "127.0.0.1", Port: 12349}
		mysqlConn, err := NewMySQLConnection("", 12349, 10, 10, &data)
		if err != nil {
			c.Fatalf("Connection Error: %s", err.Error())
		}
		defer mysqlConn.Close()
		err = mysqlConn.GetBanner()

		c.Assert(data.Banner, DeepEquals, banner)
		c.Assert(data.MySQLError.Code, Equals, uint16(1130))
		c.Assert(data.MySQLError.Message, Equals, "Host '35.2.120.138' is not allowed to connect to this MySQL server")
	}()

	// Server
	server := test.TestingBasicServer{Port: 12349, ToWrite: banner, WriteWait: 0}
	(&server).RunServer()

	wc.Wait()
}
