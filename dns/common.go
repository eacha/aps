package dns

const (
	// Header.QR Value
	QrQuery    uint16 = 0
	QrResponse uint16 = 1

	// Header.Opcode Value
	OpcodeQuery  uint16 = 0
	OpcodeIQuery uint16 = 1
	OpcodeStatus uint16 = 2

	// Header.AA Value
	NonAuthoritative uint16 = 0
	Authoritative    uint16 = 1

	// Header.TC Value
	NonTruncated uint16 = 0
	Truncated    uint16 = 1

	// Header.RD Value (query)
	NonRecursiveDesired uint16 = 0
	RecursiveDesired    uint16 = 1

	// Header.RA Value (response)
	NonRecursiveAvailable uint16 = 0
	RecursiveAvailable    uint16 = 1

	// Header.Rcode Value
	RcodeSuccess        uint16 = 0
	RcodeFormatError    uint16 = 1
	RcodeServerFailure  uint16 = 2
	RcodeNameError      uint16 = 3
	RcodeNotImplemented uint16 = 4
	RcodeRefused        uint16 = 5
)

const (
	// Type Value
	TypeNone  uint16 = 0
	TypeA     uint16 = 1
	TypeNS    uint16 = 2
	TypeCNAME uint16 = 5
	TypeSOA   uint16 = 6
	TypeWKS   uint16 = 11
	TypePTR   uint16 = 12
	TypeMX    uint16 = 15
	TypeSRV   uint16 = 33
	TypeAAAA  uint16 = 28
	TypeANY   uint16 = 255

	// Class Value
	ClassINET uint16 = 1
)

func uintToType(atype uint16) string {
	switch atype {
	case TypeNone:
		return "None"
	case TypeA:
		return "A"
	case TypeNS:
		return "NS"
	case TypeCNAME:
		return "CNAME"
	case TypeSOA:
		return "SOA"
	case TypeWKS:
		return "WKS"
	case TypePTR:
		return "PTR"
	case TypeMX:
		return "Mx"
	case TypeSRV:
		return "SRV"
	case TypeAAAA:
		return "AAAA"
	case TypeANY:
		return "ANY"
	default:
		return "Unknown"
	}
}

func uintToClass(aclass uint16) string {
	switch aclass {
	case ClassINET:
		return "INET"
	default:
		return "Unknown"
	}
}
