package scan

import (
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

type ScanSuite struct{}

var _ = Suite(&ScanSuite{})

func (s *ScanSuite) TestScanOptions(c *C) {
	so := NewScanOptions(25, 20, 50)

	c.Assert(so.Port, Equals, uint16(25))
	c.Assert(so.ConnectionTimeout, Equals, 20*time.Second)
	c.Assert(so.IOTimeout, Equals, 50*time.Second)
}
