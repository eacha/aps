package mysql

import (
	"net"
	"time"
)

type MySQLConnection struct {
	conn              net.Conn
	ip                string
	port              int
	connectionTimeout time.Duration
	ioTimeout         time.Duration
	data              *MySQL
}

type MySQLError struct {
	Code    uint16 `json:"error_code,omitempty"`
	Message string `json:"message,omitempty"`
}

type MySQL struct {
	Ip           string     `json:"ip"`
	Port         int        `json:"port"`
	Error        string     `json:"error,omitempty"`
	Banner       []byte     `json:"banner,omitempty"`
	MySQLError   MySQLError `json:"mysql_error,omitempty"`
	Proto        byte       `json:"protocol,omitempty"`
	Version      string     `json:"version,omitempty"`
	ThreadId     uint32     `json:"thread_id,omitempty"`
	Capabilities []string   `json:"capabilies,omitempty"`
	Charset      string     `json:"charset,omitempty"`
	Status       uint16     `json:"status,omitempty"`
}
