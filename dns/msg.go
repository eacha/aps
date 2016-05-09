package dns

import (
	"encoding/binary"
	"github.com/eacha/aps/util"
	"strings"
)

func newHeader(QR, Opcode, AA, TC, RD, RA, Qdcount, Ancount, Nscount, Arcount uint16) *Header {
	var h Header

	h.Id = 1
	h.Bits = completeBits(QR, Opcode, AA, TC, RD, RA)
	h.Qdcount = Qdcount
	h.Ancount = Ancount
	h.Nscount = Nscount
	h.Arcount = Arcount

	return &h
}
func completeBits(qr, opcode, aa, tc, rd, ra uint16) uint16 {
	var bits uint16

	bits |= qr << _QR
	bits |= opcode << _Opcode
	bits |= aa << _AA
	bits |= tc << _TC
	bits |= rd << _RD
	bits |= ra << _RA

	return bits
}

func (h *Header) packBuffer(buf []byte, pos int) int {
	binary.BigEndian.PutUint16(buf[pos:pos+2], h.Id)
	binary.BigEndian.PutUint16(buf[pos+2:pos+4], h.Bits)
	binary.BigEndian.PutUint16(buf[pos+4:pos+6], h.Qdcount)
	binary.BigEndian.PutUint16(buf[pos+6:pos+8], h.Ancount)
	binary.BigEndian.PutUint16(buf[pos+8:pos+10], h.Nscount)
	binary.BigEndian.PutUint16(buf[pos+10:pos+12], h.Arcount)

	return pos + 12
}

func (h *Header) unpackBuffer(buf []byte, pos int) int {
	h.Id = binary.BigEndian.Uint16(buf[pos : pos+2])
	h.Bits = binary.BigEndian.Uint16(buf[pos+2 : pos+4])
	h.Qdcount = binary.BigEndian.Uint16(buf[pos+4 : pos+6])
	h.Ancount = binary.BigEndian.Uint16(buf[pos+6 : pos+8])
	h.Nscount = binary.BigEndian.Uint16(buf[pos+8 : pos+10])
	h.Arcount = binary.BigEndian.Uint16(buf[pos+10 : pos+12])

	return pos + 12
}

func newQuestion(qname string, qtype, qclass uint16) *Question {
	var q Question

	q.Qname = qname
	q.Qtype = qtype
	q.Qclass = qclass

	return &q
}

func (q *Question) packBuffer(buf []byte, pos int) int {
	pos = q.qnameToBytes(buf, pos)
	binary.BigEndian.PutUint16(buf[pos:pos+2], q.Qtype)
	binary.BigEndian.PutUint16(buf[pos+2:pos+4], q.Qclass)

	return pos + 4
}

func (q *Question) unpackBuffer(buf []byte, pos int) int {
	pos = q.bytesToQname(buf, pos)
	q.Qtype = binary.BigEndian.Uint16(buf[pos : pos+2])
	q.Qclass = binary.BigEndian.Uint16(buf[pos+2 : pos+4])

	return pos + 4
}

func (q *Question) qnameToBytes(buf []byte, pos int) int {
	for _, seg := range strings.Split(q.Qname, ".") {
		buf[pos] = uint8(len(seg))
		pos = util.CopySliceInto([]byte(seg), buf, pos+1)
	}

	return pos + 1
}

func (q *Question) bytesToQname(buf []byte, pos int) int {
	nullPos := util.ByteIndexOf(buf, 0x00, pos)
	qnameSlice := buf[pos:nullPos]

	for i := 0; i < len(qnameSlice); {
		stringLen := int(qnameSlice[i])
		q.Qname += string(qnameSlice[i+1 : i+stringLen+1])
		q.Qname += "."
		i += stringLen + 1
	}

	return nullPos + 1
}

func NewQuery(name string, recursive uint16) *Query {
	var q Query

	q.header = newHeader(QrQuery, OpcodeQuery, NonAuthoritative, NonTruncated, recursive, NonRecursiveAvailable, 1, 0, 0, 0)
	q.question = newQuestion(name, TypeA, ClassINET)

	return &q
}

func (q *Query) Pack() []byte {
	buf := make([]byte, 1024)

	offset := q.header.packBuffer(buf, 0)
	offset = q.question.packBuffer(buf, offset)

	return buf[:offset]
}
