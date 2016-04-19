package scan

import "time"

type ScanOptions struct {
	Port              uint16
	ConnectionTimeout time.Duration
	IOTimeout         time.Duration
	// More options in the future
}

func NewScanOptions(port uint16, connectionTimeout, ioTimeout time.Duration) ScanOptions {
	var so ScanOptions

	so.Port = port
	so.ConnectionTimeout = connectionTimeout * time.Second
	so.IOTimeout = ioTimeout * time.Second

	return so
}

type Scannable interface {
	Scan(option ScanOptions, input chan string, output chan []byte)
}
