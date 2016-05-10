package connection

import (
	"fmt"
)

var (
	ConnRefusedMsg  = "Connection refued by host"
	ConnTimeoutMsg  = "Connection timeout"
	SetTimeoutMsg   = "Can't set IO timeout"
	ReadTimeoutMsg  = "Read timeout"
	WriteTimeoutMsg = "Write timeout"
	ReadMsg         = "Can't do the read operation"
	WriteMsg        = "Can't do the write operation"
)

type ConnError struct {
	msg     string
	address string
}

func (e *ConnError) Error() string {
	return fmt.Sprintf("%s, Host: %s", e.msg, e.address)
}

type IOError struct {
	msg     string
	address string
}

func (e *IOError) Error() string {
	return fmt.Sprintf("%s, Host: %s", e.msg, e.address)
}

type IOTimeoutError struct {
	msg     string
	address string
}

func (e *IOTimeoutError) Error() string {
	return fmt.Sprintf("%s, Host: %s", e.msg, e.address)
}
