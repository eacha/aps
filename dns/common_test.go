package dns

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestCommon(t *testing.T) { TestingT(t) }

type DNSCommon struct{}

var _ = Suite(&DNSCommon{})

var (
	typeTest = []struct {
		value uint16
		name  string
	}{
		{typeNone, "None"},
		{typeA, "A"},
		{typeNS, "NS"},
		{typeCNAME, "CNAME"},
		{typeSOA, "SOA"},
		{typeWKS, "WKS"},
		{typePTR, "PTR"},
		{typeMX, "Mx"},
		{typeSRV, "SRV"},
		{typeAAAA, "AAAA"},
		{typeANY, "ANY"},
		{128, "Unknown"},
	}
	classTest = []struct {
		value uint16
		name  string
	}{
		{classINET, "INET"},
		{255, "Unknown"},
	}
)

func (dns *DNSCommon) TestUintToType(c *C) {
	for _, t := range typeTest {
		c.Assert(uintToType(t.value), Equals, t.name)
	}
}

func (dns *DNSCommon) TestUintToClass(c *C) {
	for _, t := range classTest {
		c.Assert(uintToClass(t.value), Equals, t.name)
	}
}
