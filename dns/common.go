package dns

import "errors"

const (
	// Header.QR Value
	qrQuery    uint16 = 0
	qrResponse uint16 = 1

	// Header.Opcode Value
	opcodeQuery  uint16 = 0
	opcodeIQuery uint16 = 1
	opcodeStatus uint16 = 2

	// Header.AA Value
	nonAuthoritative uint16 = 0
	authoritative    uint16 = 1

	// Header.TC Value
	nonTruncated uint16 = 0
	truncated    uint16 = 1

	// Header.RD Value (query)
	nonRecursiveDesired uint16 = 0
	recursiveDesired    uint16 = 1

	// Header.RA Value (response)
	nonRecursiveAvailable uint16 = 0
	recursiveAvailable    uint16 = 1

	// Header.Rcode Value
	rcodeSuccess        uint16 = 0
	rcodeFormatError    uint16 = 1
	rcodeServerFailure  uint16 = 2
	rcodeNameError      uint16 = 3
	rcodeNotImplemented uint16 = 4
	rcodeRefused        uint16 = 5
)

const (
	// Type Value
	typeNone  uint16 = 0
	typeA     uint16 = 1
	typeNS    uint16 = 2
	typeCNAME uint16 = 5
	typeSOA   uint16 = 6
	typeWKS   uint16 = 11
	typePTR   uint16 = 12
	typeMX    uint16 = 15
	typeSRV   uint16 = 33
	typeAAAA  uint16 = 28
	typeANY   uint16 = 255

	// Class Value
	classINET uint16 = 1
)

var (
	errDNSPacketTooShort = errors.New("DNS packet too short")
)

func uintToType(atype uint16) string {
	switch atype {
	case typeNone:
		return "None"
	case typeA:
		return "A"
	case typeNS:
		return "NS"
	case typeCNAME:
		return "CNAME"
	case typeSOA:
		return "SOA"
	case typeWKS:
		return "WKS"
	case typePTR:
		return "PTR"
	case typeMX:
		return "Mx"
	case typeSRV:
		return "SRV"
	case typeAAAA:
		return "AAAA"
	case typeANY:
		return "ANY"
	default:
		return "Unknown"
	}
}

func uintToClass(aclass uint16) string {
	switch aclass {
	case classINET:
		return "INET"
	default:
		return "Unknown"
	}
}
