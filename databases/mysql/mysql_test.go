package mysql

import (
	"encoding/base64"
	. "gopkg.in/check.v1"
	"testing"
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
	c.Assert(data.Charset, Equals, byte(8))
	c.Assert(data.Status, Equals, uint16(2))

}
