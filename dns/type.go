package dns

const (
	// Valid Header.QR
	QrQuery    uint16 = 0
	QrResponse uint16 = 1

	// Message Header.Opcode
	OpcodeQuery  uint16 = 0
	OpcodeIQuery uint16 = 1
	OpcodeStatus uint16 = 2

	// Valid Header.AA
	NonAuthoritative uint16 = 0
	Authoritative    uint16 = 1

	// Valid Header.TC
	NonTruncated uint16 = 0
	Truncated    uint16 = 1

	// Valid Header.RD (query)
	NonRecursiveDesired uint16 = 0
	RecursiveDesired    uint16 = 1

	// Valid Header.RA (response)
	NonRecursiveAvailable uint16 = 0
	RecursiveAvailable    uint16 = 1

	// Message Header.Rcode
	RcodeSuccess        uint16 = 0
	RcodeFormatError    uint16 = 1
	RcodeServerFailure  uint16 = 2
	RcodeNameError      uint16 = 3
	RcodeNotImplemented uint16 = 4
	RcodeRefused        uint16 = 5

	// Valid Question.Qtype
	TypeNone uint16 = 0
	TypeA    uint16 = 1

	// Valid Question.Qclass
	ClassINET uint16 = 1
)

type Header struct {
	Id      uint16
	Bits    uint16
	Qdcount uint16
	Ancount uint16
	Nscount uint16
	Arcount uint16
}

const (
	// Header.Bits
	_QR     = 15
	_Opcode = 14
	_AA     = 10
	_TC     = 9
	_RD     = 8
	_RA     = 7
)

type Question struct {
	Qname  string
	Qtype  uint16
	Qclass uint16
}

type Query struct {
	header   Header
	question Question
}

type Answer struct {
	Aname    string
	Atype    uint16
	Aclass   uint16
	Attl     uint32
	RdLength uint16
}
