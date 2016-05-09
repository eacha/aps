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
	PACK_ANSWER_CNAME = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	PACK_ANSWER_A =     []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0xc0, 0x28, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x04, 0x42, 0x93, 0xf4, 0xc2}
	PACK_ANSWER_OTHER = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x03, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	PACK_RESPONSE = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0xc0, 0x28, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x04, 0x42, 0x93, 0xf4, 0xc2}

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
	pos := qnameToBytes(question.Qname, buffer, 0)

	c.Assert(buffer[:pos], DeepEquals, QUESTION_NAME)
}

func (dns *DnsMsg) TestByteToDnsQueryName(c *C) {
	name, _ := bytesToQname(QUESTION_NAME, 0)

	c.Assert(name, Equals, "www.uchile.cl.")
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
	c.Assert(question.Qtype, Equals, TypeA)
	c.Assert(question.Qclass, Equals, ClassINET)
}

func (dns *DnsMsg) TestPackQuery(c *C) {
	query := NewQuery(URL, RecursiveDesired)
	buffer := query.Pack()

	c.Assert(buffer, DeepEquals, PACK_QUERY)
}

func (dns *DnsMsg) TestUnPackQuery(c *C) {
	var query Query
	query.UnPack(PACK_QUERY)

	c.Assert(query.Header.Id, Equals, ID)
	c.Assert(query.Header.Bits, Equals, BITS)
	c.Assert(query.Header.Qdcount, Equals, QD_COUNT)
	c.Assert(query.Header.Ancount, Equals, AN_COUNT)
	c.Assert(query.Header.Nscount, Equals, NS_COUNT)
	c.Assert(query.Header.Arcount, Equals, AR_COUNT)

	c.Assert(query.Question.Qname, Equals, "www.uchile.cl.")
	c.Assert(query.Question.Qtype, Equals, TypeA)
	c.Assert(query.Question.Qclass, Equals, ClassINET)
}

func (dns *DnsMsg) TestUnPackAnswerCname(c* C) {
	var answer Answer
	answer.unpackBuffer(PACK_ANSWER_CNAME, 28)

	c.Assert(answer.Aname, Equals, "www.ble.cl.")
	c.Assert(answer.Atype, Equals, TypeCNAME)
	c.Assert(answer.Aclass, Equals, ClassINET)
	c.Assert(answer.RdLength, Equals, uint16(8))
	c.Assert(answer.RdDataNS, Equals, "ble.cl.")
}

func (dns *DnsMsg) TestUnPackAnswerA(c* C) {
	var answer Answer
	answer.unpackBuffer(PACK_ANSWER_A, 48)

	c.Assert(answer.Aname, Equals, "ble.cl.")
	c.Assert(answer.Atype, Equals, TypeA)
	c.Assert(answer.Aclass, Equals, ClassINET)
	c.Assert(answer.RdLength, Equals, uint16(4))
	c.Assert(answer.RdDataA, Equals, uint32(1116992706))
}

func (dns *DnsMsg) TestUnPackAnswerOther(c* C) {
	var answer Answer
	answer.unpackBuffer(PACK_ANSWER_OTHER, 28)

	c.Assert(answer.Aname, Equals, "www.ble.cl.")
	c.Assert(answer.Atype, Equals, uint16(3))
	c.Assert(answer.Aclass, Equals, ClassINET)
	c.Assert(answer.RdLength, Equals, uint16(8))
	c.Assert(answer.RdDataNS, Equals, "")
	c.Assert(answer.RdDataA, Equals, uint32(0))
}

func (dns *DnsMsg) TestUnPackResponse(c *C){
	var response Response
	response.UnPack(PACK_RESPONSE)

	c.Assert(response.Header.Id, Equals, ID)
	c.Assert(response.Header.Bits, Equals, BITS)
	c.Assert(response.Header.Qdcount, Equals, QD_COUNT)
	c.Assert(response.Header.Ancount, Equals, uint16(2))
	c.Assert(response.Header.Nscount, Equals, NS_COUNT)
	c.Assert(response.Header.Arcount, Equals, AR_COUNT)

	c.Assert(response.Question[0].Qname, Equals, "www.ble.cl.")
	c.Assert(response.Question[0].Qtype, Equals, TypeA)
	c.Assert(response.Question[0].Qclass, Equals, ClassINET)

	c.Assert(response.Answer[0].Aname, Equals, "www.ble.cl.")
	c.Assert(response.Answer[0].Atype, Equals, TypeCNAME)
	c.Assert(response.Answer[0].Aclass, Equals, ClassINET)
	c.Assert(response.Answer[0].RdLength, Equals, uint16(8))
	c.Assert(response.Answer[0].RdDataNS, Equals, "ble.cl.")

	c.Assert(response.Answer[1].Aname, Equals, "ble.cl.")
	c.Assert(response.Answer[1].Atype, Equals, TypeA)
	c.Assert(response.Answer[1].Aclass, Equals, ClassINET)
	c.Assert(response.Answer[1].RdLength, Equals, uint16(4))
	c.Assert(response.Answer[1].RdDataA, Equals, uint32(1116992706))
}


