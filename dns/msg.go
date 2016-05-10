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
	pos = qnameToBytes(q.Qname, buf, pos)
	binary.BigEndian.PutUint16(buf[pos:pos+2], q.Qtype)
	binary.BigEndian.PutUint16(buf[pos+2:pos+4], q.Qclass)

	return pos + 4
}

func (q *Question) unpackBuffer(buf []byte, pos int) int {
	q.Qname, pos = bytesToQname(buf, pos)
	q.Qtype = binary.BigEndian.Uint16(buf[pos : pos+2])
	q.Qclass = binary.BigEndian.Uint16(buf[pos+2 : pos+4])

	return pos + 4
}

func NewQuery(name string, recursive uint16) *Query {
	var q Query

	q.Header = *newHeader(QrQuery, OpcodeQuery, NonAuthoritative, NonTruncated, recursive, NonRecursiveAvailable, 1, 0, 0, 0)
	q.Question = *newQuestion(name, TypeA, ClassINET)

	return &q
}

func (q *Query) Pack() []byte {
	buf := make([]byte, 1024)

	offset := q.Header.packBuffer(buf, 0)
	offset = q.Question.packBuffer(buf, offset)

	return buf[:offset]
}

func (q *Query) UnPack(buf []byte) {
	pos := q.Header.unpackBuffer(buf, 0)
	pos = q.Question.unpackBuffer(buf, pos)
}

func (a *Answer) unpackBuffer(buf []byte, pos int) int {
	namePtr :=  binary.BigEndian.Uint16(buf[pos : pos+2]) & 0x3FFF

	a.Aname, _ = bytesToQname(buf, int(namePtr))
	a.Atype = binary.BigEndian.Uint16(buf[pos+2 : pos+4])
	a.Aclass = binary.BigEndian.Uint16(buf[pos+4 : pos+6])
	a.Attl = binary.BigEndian.Uint32(buf[pos+6 : pos+10])
	a.RdLength = binary.BigEndian.Uint16(buf[pos+10 : pos+12])

	if a.Atype == TypeA {
		a.RdDataA = binary.BigEndian.Uint32(buf[pos+12 : pos+16])
		pos += 16
	} else if a.Atype == TypeNS || a.Atype == TypeCNAME {
		a.RdDataNS, pos = bytesToQname(buf, pos+12)
	} else {
		pos += 12 + int(a.RdLength)
	}

	return pos
}

func (r *Response) UnPack(buf []byte) {
	pos := r.Header.unpackBuffer(buf, 0)

	r.Question = make([]Question, r.Header.Qdcount)
	r.Answer = make([]Answer, r.Header.Ancount)

	for i := 0; i < len(r.Question); i++ {
		pos = r.Question[i].unpackBuffer(buf, pos)
	}
	for i := 0; i < len(r.Answer); i++ {
		pos = r.Answer[i].unpackBuffer(buf, pos)
	}
}

func qnameToBytes(name string, buf []byte, pos int) int {
	for _, seg := range strings.Split(name, ".") {
		buf[pos] = uint8(len(seg))
		pos = util.CopySliceInto([]byte(seg), buf, pos+1)
	}

	return pos + 1
}

func bytesToQname(buf []byte, pos int) (string, int) {
	nullPos := util.ByteIndexOf(buf, 0x00, pos)
	qnameSlice, name := buf[pos:nullPos], ""

	for i := 0; i < len(qnameSlice); {
		stringLen := int(qnameSlice[i])
		name += string(qnameSlice[i+1 : i+stringLen+1])
		name += "."
		i += stringLen + 1
	}

	return name, nullPos + 1
}
