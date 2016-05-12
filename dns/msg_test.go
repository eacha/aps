package dns

import (
	"testing"

	. "gopkg.in/check.v1"
)

var (
	// Packets
	queryBits       uint16 = 0x0100
	url                    = "www.uchile.cl"
	questionName           = []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	packHeader             = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	packQuestion           = []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01}
	packQuery              = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01}
	packAnswerCname        = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	packAnswerA            = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0xc0, 0x28, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x04, 0x42, 0x93, 0xf4, 0xc2}
	packAnswerOther        = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x03, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	packResponse           = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0xc0, 0x28, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x04, 0x42, 0x93, 0xf4, 0xc2}

	// Header
	id      uint16 = 1
	bits    uint16 = 256
	qdCount uint16 = 1
	anCount uint16
	nsCount uint16
	arCount uint16
)

func Test(t *testing.T) { TestingT(t) }

type DNSMsg struct{}

var _ = Suite(&DNSMsg{})

func (dns *DNSMsg) TestCompleteBits(c *C) {
	bits := completeBits(qrQuery, opcodeQuery, nonAuthoritative, nonTruncated, recursiveAvailable, nonRecursiveAvailable)
	c.Assert(bits, Equals, queryBits)
}

func (dns *DNSMsg) TestPackHeader(c *C) {
	buffer := make([]byte, 1024)
	header := newHeader(qrQuery, opcodeQuery, nonAuthoritative, nonTruncated, recursiveAvailable, nonRecursiveAvailable, 1, 0, 0, 0)
	pos := header.packBuffer(buffer, 0)

	c.Assert(buffer[:pos], DeepEquals, packHeader)
}

func (dns *DNSMsg) TestUnPackHeader(c *C) {
	var header Header
	pos, _ := header.unpackBuffer(packHeader, 0)

	c.Assert(header.ID, Equals, id)
	c.Assert(header.Bits, Equals, bits)
	c.Assert(header.Qdcount, Equals, qdCount)
	c.Assert(header.Ancount, Equals, anCount)
	c.Assert(header.Nscount, Equals, nsCount)
	c.Assert(header.Arcount, Equals, arCount)

	c.Assert(pos, Equals, 12)
}

func (dns *DNSMsg) TestDnsQueryNameToByte(c *C) {
	buffer := make([]byte, 1024)
	question := newQuestion(url, typeA, classINET)
	pos := qnameToBytes(question.Qname, buffer, 0)

	c.Assert(buffer[:pos], DeepEquals, questionName)
}

func (dns *DNSMsg) TestByteToDnsQueryName(c *C) {
	name, _, _ := bytesToQname(questionName, 0)

	c.Assert(name, Equals, "www.uchile.cl.")
}

func (dns *DNSMsg) TestPackQuestion(c *C) {
	buffer := make([]byte, 1024)
	question := newQuestion(url, typeA, classINET)
	pos := question.packBuffer(buffer, 0)

	c.Assert(buffer[:pos], DeepEquals, packQuestion)
}

func (dns *DNSMsg) TestUnPackQuestion(c *C) {
	var question Question
	question.unpackBuffer(packQuestion, 0)

	c.Assert(question.Qname, Equals, "www.uchile.cl.")
	c.Assert(question.Qtype, Equals, typeA)
	c.Assert(question.Qclass, Equals, classINET)
}

func (dns *DNSMsg) TestPackQuery(c *C) {
	query := NewQuery(url, recursiveDesired)
	buffer := query.Pack()

	c.Assert(buffer, DeepEquals, packQuery)
}

func (dns *DNSMsg) TestUnPackQuery(c *C) {
	var query Query
	query.UnPack(packQuery)

	c.Assert(query.Header.ID, Equals, id)
	c.Assert(query.Header.Bits, Equals, bits)
	c.Assert(query.Header.Qdcount, Equals, qdCount)
	c.Assert(query.Header.Ancount, Equals, anCount)
	c.Assert(query.Header.Nscount, Equals, nsCount)
	c.Assert(query.Header.Arcount, Equals, arCount)

	c.Assert(query.Question.Qname, Equals, "www.uchile.cl.")
	c.Assert(query.Question.Qtype, Equals, typeA)
	c.Assert(query.Question.Qclass, Equals, classINET)
}

func (dns *DNSMsg) TestUnPackAnswerCname(c *C) {
	var answer Answer
	answer.unpackBuffer(packAnswerCname, 28)

	c.Assert(answer.Aname, Equals, "www.ble.cl.")
	c.Assert(answer.Atype, Equals, typeCNAME)
	c.Assert(answer.Aclass, Equals, classINET)
	c.Assert(answer.RdLength, Equals, uint16(8))
	c.Assert(answer.RdData, Equals, "ble.cl.")
}

func (dns *DNSMsg) TestUnPackAnswerA(c *C) {
	var answer Answer
	answer.unpackBuffer(packAnswerA, 48)

	c.Assert(answer.Aname, Equals, "ble.cl.")
	c.Assert(answer.Atype, Equals, typeA)
	c.Assert(answer.Aclass, Equals, classINET)
	c.Assert(answer.RdLength, Equals, uint16(4))
	c.Assert(answer.RdData, Equals, "66.147.244.194")
}

func (dns *DNSMsg) TestUnPackAnswerOther(c *C) {
	var answer Answer
	answer.unpackBuffer(packAnswerOther, 28)

	c.Assert(answer.Aname, Equals, "www.ble.cl.")
	c.Assert(answer.Atype, Equals, uint16(3))
	c.Assert(answer.Aclass, Equals, classINET)
	c.Assert(answer.RdLength, Equals, uint16(8))
	c.Assert(answer.RdData, Equals, "")
}

func (dns *DNSMsg) TestUnPackResponse(c *C) {
	var response Response
	response.UnPack(packResponse)

	c.Assert(response.Header.ID, Equals, id)
	c.Assert(response.Header.Bits, Equals, bits)
	c.Assert(response.Header.Qdcount, Equals, qdCount)
	c.Assert(response.Header.Ancount, Equals, uint16(2))
	c.Assert(response.Header.Nscount, Equals, nsCount)
	c.Assert(response.Header.Arcount, Equals, arCount)

	c.Assert(response.Question[0].Qname, Equals, "www.ble.cl.")
	c.Assert(response.Question[0].Qtype, Equals, typeA)
	c.Assert(response.Question[0].Qclass, Equals, classINET)

	c.Assert(response.Answer[0].Aname, Equals, "www.ble.cl.")
	c.Assert(response.Answer[0].Atype, Equals, typeCNAME)
	c.Assert(response.Answer[0].Aclass, Equals, classINET)
	c.Assert(response.Answer[0].RdLength, Equals, uint16(8))
	c.Assert(response.Answer[0].RdData, Equals, "ble.cl.")

	c.Assert(response.Answer[1].Aname, Equals, "ble.cl.")
	c.Assert(response.Answer[1].Atype, Equals, typeA)
	c.Assert(response.Answer[1].Aclass, Equals, classINET)
	c.Assert(response.Answer[1].RdLength, Equals, uint16(4))
	c.Assert(response.Answer[1].RdData, Equals, "66.147.244.194")
}
