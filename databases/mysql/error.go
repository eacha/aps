package mysql

import "errors"

var (
	ConnectionRefusedError = errors.New("Error: Connection refued by host")
	ConnectionTimeoutError = errors.New("Error: Connection timeout")
	SetTimeoutError        = errors.New("Error: Can't set IO timeout")
	IOTimeoutError         = errors.New("Error: Input-Output timeout")
	ReadError              = errors.New("Error: Can't do the read operation")
	WriteError             = errors.New("Error: Can't do the write operation")
	MalFormedPacketError   = errors.New("Error: MySQL packet is malformed")
)
