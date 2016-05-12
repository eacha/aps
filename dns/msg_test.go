package dns

import (
	"testing"

	. "gopkg.in/check.v1"
)

var (
	// Helper functions
	url              = "www.uchile.cl"
	queryBits uint16 = 0x0100

	// Packets
	questionName                = []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	packHeader                  = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	packQuestion                = []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01}
	packQuery                   = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x06, 0x75, 0x63, 0x68, 0x69, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01}
	packAnswerCname             = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	packAnswerA                 = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0xc0, 0x28, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x04, 0x42, 0x93, 0xf4, 0xc2}
	packAnswerOther             = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x03, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00}
	packResponse                = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0xc0, 0x28, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x04, 0x42, 0x93, 0xf4, 0xc2}
	unpackNameError             = []byte{0x12, 0x63, 0x32, 0x74, 0x54, 0x12, 0x63, 0x32, 0x74, 0x54, 0x12, 0x63, 0x32, 0x74, 0x54}
	unpackTypeAError            = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0xc0, 0x28, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x04}
	unpackTypeCnameError        = []byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x00, 0x00, 0x01, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x2e, 0xec, 0x00, 0x08, 0x03, 0x62, 0x6c, 0x65, 0x02, 0x63, 0x6c, 0x01}
	unpackResponseQuestionError = packHeader

	// JSON Marshal
	jsonQuestion = []byte{123, 34, 116, 121, 112, 101, 34, 58, 34, 65, 34, 44, 34, 99, 108, 97, 115, 115, 34, 58, 34, 73, 78, 69, 84, 34, 44, 34, 110, 97, 109, 101, 34, 58, 34, 116, 101, 115, 116, 46, 99, 108, 34, 125}
	jsonAnswer   = []byte{123, 34, 116, 121, 112, 101, 34, 58, 34, 65, 34, 44, 34, 99, 108, 97, 115, 115, 34, 58, 34, 73, 78, 69, 84, 34, 44, 34, 110, 97, 109, 101, 34, 58, 34, 116, 101, 115, 116, 46, 99, 108, 34, 44, 34, 116, 116, 108, 34, 58, 51, 48, 48, 48, 44, 34, 114, 100, 95, 108, 101, 110, 103, 116, 104, 34, 58, 49, 49, 44, 34, 114, 100, 95, 100, 97, 116, 97, 34, 58, 34, 119, 119, 119, 46, 116, 101, 115, 116, 46, 99, 108, 34, 125}
	// Header
	id      uint16 = 1
	bits    uint16 = 256
	qdCount uint16 = 1
	anCount uint16
	nsCount uint16
	arCount uint16
)

func TestMsg(t *testing.T) { TestingT(t) }

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

func (dns *DNSMsg) TestUnpackHeader(c *C) {
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

func (dns *DNSMsg) TestUnpackHeaderError(c *C) {
	var header Header
	pos, err := header.unpackBuffer(packHeader, 50)

	c.Assert(pos, Equals, -1)
	c.Assert(err, Equals, errDNSPacketTooShort)
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

func (dns *DNSMsg) TestUnpackQuestion(c *C) {
	var question Question
	question.unpackBuffer(packQuestion, 0)

	c.Assert(question.Qname, Equals, "www.uchile.cl.")
	c.Assert(question.Qtype, Equals, typeA)
	c.Assert(question.Qclass, Equals, classINET)
}

func (dns *DNSMsg) TestUnpackQuestionError(c *C) {
	var question Question
	pos, err := question.unpackBuffer(unpackNameError, 0)

	c.Assert(pos, Equals, -1)
	c.Assert(err, Equals, errDNSPacketTooShort)
}

func (dns *DNSMsg) TestMarshalQuestion(c *C) {
	question := newQuestion("test.cl", typeA, classINET)
	j, _ := question.MarshalJSON()

	c.Assert(j, DeepEquals, jsonQuestion)
}

func (dns *DNSMsg) TestPackQuery(c *C) {
	query := NewQuery(url, recursiveDesired)
	buffer := query.Pack()

	c.Assert(buffer, DeepEquals, packQuery)
}

func (dns *DNSMsg) TestUnpackQuery(c *C) {
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

func (dns *DNSMsg) TestUnpackQueryError(c *C) {
	var (
		query     Query
		emptyByte []byte
	)

	err := query.UnPack(emptyByte)
	c.Assert(err, Equals, errDNSPacketTooShort)

	err = query.UnPack(make([]byte, 12))
	c.Assert(err, Equals, errDNSPacketTooShort)

}

func (dns *DNSMsg) TestUnpackAnswerCname(c *C) {
	var answer Answer
	answer.unpackBuffer(packAnswerCname, 28)

	c.Assert(answer.Aname, Equals, "www.ble.cl.")
	c.Assert(answer.Atype, Equals, typeCNAME)
	c.Assert(answer.Aclass, Equals, classINET)
	c.Assert(answer.RdLength, Equals, uint16(8))
	c.Assert(answer.RdData, Equals, "ble.cl.")
}

func (dns *DNSMsg) TestUnpackAnswerA(c *C) {
	var answer Answer
	answer.unpackBuffer(packAnswerA, 48)

	c.Assert(answer.Aname, Equals, "ble.cl.")
	c.Assert(answer.Atype, Equals, typeA)
	c.Assert(answer.Aclass, Equals, classINET)
	c.Assert(answer.RdLength, Equals, uint16(4))
	c.Assert(answer.RdData, Equals, "66.147.244.194")
}

func (dns *DNSMsg) TestUnpackAnswerOther(c *C) {
	var answer Answer
	answer.unpackBuffer(packAnswerOther, 28)

	c.Assert(answer.Aname, Equals, "www.ble.cl.")
	c.Assert(answer.Atype, Equals, uint16(3))
	c.Assert(answer.Aclass, Equals, classINET)
	c.Assert(answer.RdLength, Equals, uint16(8))
	c.Assert(answer.RdData, Equals, "")
}

func (dns *DNSMsg) TestUnpackAnswerError(c *C) {
	var (
		answer    Answer
		emptyByte []byte
	)

	pos, err := answer.unpackBuffer(emptyByte, 0)
	c.Assert(pos, Equals, -1)
	c.Assert(err, Equals, errDNSPacketTooShort)

	pos, err = answer.unpackBuffer(unpackNameError, 0)
	c.Assert(pos, Equals, -1)
	c.Assert(err, Equals, errDNSPacketTooShort)

	pos, err = answer.unpackBuffer(unpackTypeAError, 48)
	c.Assert(pos, Equals, -1)
	c.Assert(err, Equals, errDNSPacketTooShort)

	pos, err = answer.unpackBuffer(unpackTypeCnameError, 28)
	c.Assert(pos, Equals, -1)
	c.Assert(err, Equals, errDNSPacketTooShort)

}

func (dns *DNSMsg) TestMarshalAnswer(c *C) {
	answer := &Answer{"test.cl", typeA, classINET, 3000, 11, "www.test.cl"}
	j, _ := answer.MarshalJSON()

	c.Assert(j, DeepEquals, jsonAnswer)
}

func (dns *DNSMsg) TestUnpackResponse(c *C) {
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

func (dns *DNSMsg) TestUnpackResponseError(c *C) {
	var (
		response  Response
		emptyByte []byte
	)

	err := response.UnPack(emptyByte)
	c.Assert(err, Equals, errDNSPacketTooShort)

	err = response.UnPack(unpackResponseQuestionError)
	c.Assert(err, Equals, errDNSPacketTooShort)

	err = response.UnPack(unpackTypeCnameError)
	c.Assert(err, Equals, errDNSPacketTooShort)
}
