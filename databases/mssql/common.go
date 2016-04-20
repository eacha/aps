package mssql

const DEFAULT_VERSION = "9.00.1399.00"

// Packet Type
const (
	Query            = 0x01
	Response         = 0x04
	Login            = 0x10
	NTAuthentication = 0x11
	PreLogin         = 0x12
)

// Pre-Login Option Types
const (
	Version    = 0x00
	Encryption = 0x01
	InstOpt    = 0x02
	ThreadId   = 0x03
	MARS       = 0x04
	Terminator = 0xFF
)

var OPTION_LENGTH_CLIENT = map[byte]uint16{
	Version:    6,
	Encryption: 1,
	//InstOpt: -1,
	ThreadId:   4,
	MARS:       1,
	Terminator: 0,
}
