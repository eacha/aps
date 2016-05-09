package dns

import (
	. "gopkg.in/check.v1"
	"testing"
)

var (
	// Packets
	QUERY_BITS    uint16 = 0x0100
	URL                  = "www.uchile.cl"
	QUESTION_NAME        = []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	PACK_HEADER          = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	PACK_QUESTION        = []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01}
	PACK_QUERY           = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01}

	// Header
	ID       uint16 = 1
	BITS     uint16 = 256
	QD_COUNT uint16 = 1
	AN_COUNT uint16 = 0
	NS_COUNT uint16 = 0
	AR_COUNT uint16 = 0
)

func Test(t *testing.T) { TestingT(t) }

type DnsMsg struct{}

var _ = Suite(&DnsMsg{})

func (dns *DnsMsg) TestCompleteBits(c *C) {
	bits := completeBits(QrQuery, OpcodeQuery, NonAuthoritative, NonTruncated, RecursiveAvailable, NonRecursiveAvailable)
	c.Assert(bits, Equals, QUERY_BITS)
}

func (dns *DnsMsg) TestPackHeader(c *C) {
	buffer := make([]byte, 1024)
	header := newHeader(QrQuery, OpcodeQuery, NonAuthoritative, NonTruncated, RecursiveAvailable, NonRecursiveAvailable, 1, 0, 0, 0)
	pos := header.packBuffer(buffer, 0)

	c.Assert(buffer[:pos], DeepEquals, PACK_HEADER)
}

func (dns *DnsMsg) TestUnPackHeader(c *C) {
	var header Header
	pos := header.unpackBuffer(PACK_HEADER, 0)

	c.Assert(header.Id, Equals, ID)
	c.Assert(header.Bits, Equals, BITS)
	c.Assert(header.Qdcount, Equals, QD_COUNT)
	c.Assert(header.Ancount, Equals, AN_COUNT)
	c.Assert(header.Nscount, Equals, NS_COUNT)
	c.Assert(header.Arcount, Equals, AR_COUNT)

	c.Assert(pos, Equals, 12)
}

func (dns *DnsMsg) TestDnsQueryNameToByte(c *C) {
	buffer := make([]byte, 1024)
	question := newQuestion(URL, TypeA, ClassINET)
	pos := question.qnameToBytes(buffer, 0)

	c.Assert(buffer[:pos], DeepEquals, QUESTION_NAME)
}

func (dns *DnsMsg) TestByteToDnsQueryName(c *C) {
	var question Question
	question.bytesToQname(QUESTION_NAME, 0)

	c.Assert(question.Qname, Equals, "www.uchile.cl.")
}

func (dns *DnsMsg) TestPackQuestion(c *C) {
	buffer := make([]byte, 1024)
	question := newQuestion(URL, TypeA, ClassINET)
	pos := question.packBuffer(buffer, 0)

	c.Assert(buffer[:pos], DeepEquals, PACK_QUESTION)
}

func (dns *DnsMsg) TestUnPackQuestion(c *C) {
	var question Question
	question.unpackBuffer(PACK_QUESTION, 0)

	c.Assert(question.Qname, Equals, "www.uchile.cl.")
	c.Assert(question.Qtype, Equals, uint16(1))
	c.Assert(question.Qclass, Equals, uint16(1))
}

func (dns *DnsMsg) TestPackQuery(c *C) {
	query := NewQuery(URL, RecursiveDesired)
	buffer := query.Pack()

	c.Assert(buffer, DeepEquals, PACK_QUERY)
}
