package dns

import (
	"encoding/binary"
	"fmt"
	"strings"

	"encoding/json"

	"github.com/eacha/aps/util"
)

const (
	bufferSize   = 1024
	headerSize   = 12
	questionSize = 4
	answerSize   = 12
)

const (
	// Header.Bits
	_QR     = 15
	_Opcode = 14
	_AA     = 10
	_TC     = 9
	_RD     = 8
	_RA     = 7
)

type Header struct {
	ID      uint16 `json:"id,omitempty"`
	Bits    uint16 `json:"bits,omitempty"`
	Qdcount uint16 `json:"qd_count,omitempty"`
	Ancount uint16 `json:"an_count,omitempty"`
	Nscount uint16 `json:"ns_count,omitempty"`
	Arcount uint16 `json:"ar_count,omitempty"`
}

type Question struct {
	Qname  string `json:"name,omitempty"`
	Qtype  uint16 `json:"type,omitempty"`
	Qclass uint16 `json:"class,omitempty"`
}

type Answer struct {
	Aname    string `json:"name,omitempty"`
	Atype    uint16 `json:"type,omitempty"`
	Aclass   uint16 `json:"class,omitempty"`
	Attl     uint32 `json:"ttl,omitempty"`
	RdLength uint16 `json:"rd_length,omitempty"`
	RdData   string `json:"rd_data,omitempty"`
}

type Query struct {
	Header   Header
	Question Question
}

type Response struct {
	Header   Header
	Question []Question
	Answer   []Answer
}

func newHeader(QR, Opcode, AA, TC, RD, RA, Qdcount, Ancount, Nscount, Arcount uint16) *Header {
	var h Header

	h.ID = 1
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
	binary.BigEndian.PutUint16(buf[pos:pos+2], h.ID)
	binary.BigEndian.PutUint16(buf[pos+2:pos+4], h.Bits)
	binary.BigEndian.PutUint16(buf[pos+4:pos+6], h.Qdcount)
	binary.BigEndian.PutUint16(buf[pos+6:pos+8], h.Ancount)
	binary.BigEndian.PutUint16(buf[pos+8:pos+10], h.Nscount)
	binary.BigEndian.PutUint16(buf[pos+10:pos+12], h.Arcount)

	return pos + headerSize
}

func (h *Header) unpackBuffer(buf []byte, pos int) (int, error) {
	if pos+headerSize > len(buf) {
		return -1, errDNSPacketTooShort
	}

	h.ID = binary.BigEndian.Uint16(buf[pos : pos+2])
	h.Bits = binary.BigEndian.Uint16(buf[pos+2 : pos+4])
	h.Qdcount = binary.BigEndian.Uint16(buf[pos+4 : pos+6])
	h.Ancount = binary.BigEndian.Uint16(buf[pos+6 : pos+8])
	h.Nscount = binary.BigEndian.Uint16(buf[pos+8 : pos+10])
	h.Arcount = binary.BigEndian.Uint16(buf[pos+10 : pos+12])

	return pos + headerSize, nil
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

	return pos + questionSize
}

func qnameToBytes(name string, buf []byte, pos int) int {
	for _, seg := range strings.Split(name, ".") {
		buf[pos] = uint8(len(seg))
		pos = util.CopySliceInto([]byte(seg), buf, pos+1)
	}

	return pos + 1
}

func (q *Question) unpackBuffer(buf []byte, pos int) (int, error) {
	var err error

	q.Qname, pos, err = bytesToQname(buf, pos)
	if err != nil || pos+questionSize > len(buf) {
		return -1, errDNSPacketTooShort
	}

	q.Qtype = binary.BigEndian.Uint16(buf[pos : pos+2])
	q.Qclass = binary.BigEndian.Uint16(buf[pos+2 : pos+4])

	return pos + questionSize, nil
}

func bytesToQname(buf []byte, pos int) (string, int, error) {
	nullPos := util.ByteIndexOf(buf, 0x00, pos)
	if nullPos == -1 {
		return "", -1, errDNSPacketTooShort
	}

	qnameSlice, name := buf[pos:nullPos], ""

	for i := 0; i < len(qnameSlice); {
		stringLen := int(qnameSlice[i])
		name += string(qnameSlice[i+1 : i+stringLen+1])
		name += "."
		i += stringLen + 1
	}

	return name, nullPos + 1, nil
}

func (q Question) MarshalJSON() ([]byte, error) {
	type Alias Question
	return json.Marshal(&struct {
		Qtype  string `json:"type"`
		Qclass string `json:"class"`
		*Alias
	}{
		Qtype:  uintToType(q.Qtype),
		Qclass: uintToClass(q.Qclass),
		Alias:  (*Alias)(&q),
	})
}

func (a *Answer) unpackBuffer(buf []byte, pos int) (int, error) {
	var err error
	if pos+answerSize > len(buf) {
		return -1, errDNSPacketTooShort
	}

	namePtr := binary.BigEndian.Uint16(buf[pos:pos+2]) & 0x3FFF
	a.Aname, _, err = bytesToQname(buf, int(namePtr))
	if err != nil {
		return -1, err
	}

	a.Atype = binary.BigEndian.Uint16(buf[pos+2 : pos+4])
	a.Aclass = binary.BigEndian.Uint16(buf[pos+4 : pos+6])
	a.Attl = binary.BigEndian.Uint32(buf[pos+6 : pos+10])
	a.RdLength = binary.BigEndian.Uint16(buf[pos+10 : pos+12])

	if a.Atype == typeA {
		if pos+answerSize+4 > len(buf) {
			return -1, errDNSPacketTooShort
		}
		intIP := binary.BigEndian.Uint32(buf[pos+12 : pos+16])
		a.RdData = fmt.Sprintf("%d.%d.%d.%d", byte(intIP>>24), byte(intIP>>16), byte(intIP>>8), byte(intIP))
		pos += 16
	} else if a.Atype == typeNS || a.Atype == typeCNAME {
		a.RdData, pos, err = bytesToQname(buf, pos+12)
		if err != nil {
			return -1, err
		}
	} else {
		pos += 12 + int(a.RdLength)
	}

	return pos, nil
}

func (a Answer) MarshalJSON() ([]byte, error) {
	type Alias Answer
	return json.Marshal(&struct {
		Atype  string `json:"type"`
		Aclass string `json:"class"`
		*Alias
	}{
		Atype:  uintToType(a.Atype),
		Aclass: uintToClass(a.Aclass),
		Alias:  (*Alias)(&a),
	})
}

func NewQuery(name string, recursive uint16) *Query {
	var q Query

	q.Header = *newHeader(qrQuery, opcodeQuery, nonAuthoritative, nonTruncated, recursive, nonRecursiveAvailable, 1, 0, 0, 0)
	q.Question = *newQuestion(name, typeA, classINET)

	return &q
}

func (q *Query) Pack() []byte {
	buf := make([]byte, bufferSize)

	offset := q.Header.packBuffer(buf, 0)
	offset = q.Question.packBuffer(buf, offset)

	return buf[:offset]
}

func (q *Query) UnPack(buf []byte) error {
	pos, err := q.Header.unpackBuffer(buf, 0)
	if err != nil {
		return err
	}

	pos, err = q.Question.unpackBuffer(buf, pos)
	if err != nil {
		return err
	}

	return nil
}

func (r *Response) UnPack(buf []byte) error {
	pos, err := r.Header.unpackBuffer(buf, 0)
	if err != nil {
		return err
	}

	r.Question = make([]Question, r.Header.Qdcount)
	r.Answer = make([]Answer, r.Header.Ancount)

	for i := 0; i < len(r.Question); i++ {
		pos, err = r.Question[i].unpackBuffer(buf, pos)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(r.Answer); i++ {
		pos, err = r.Answer[i].unpackBuffer(buf, pos)
		if err != nil {
			return err
		}
	}

	return nil
}
