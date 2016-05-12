package dns

import (
	"time"

	"github.com/eacha/aps/tools/connection"
)

type DNSConn struct {
	conn *connection.ConnTimeout
}

func NewDNSConn(protocol, address string, port int, connectionTimeout, ioTimeout time.Duration) (*DNSConn, error) {
	var dnsConn DNSConn
	var err error

	dnsConn.conn, err = connection.NewConnTimeout(protocol, address, port, connectionTimeout, ioTimeout)
	if err != nil {
		return nil, err
	}

	return &dnsConn, nil
}

func (c *DNSConn) OpenResolver(question string) (*Response, error) {
	query := NewQuery(question, recursiveDesired)
	buf := make([]byte, 1024)
	pack := query.Pack()

	_, err := c.conn.Write(pack)
	if err != nil {
		return nil, err
	}

	b, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	var response Response
	response.UnPack(buf[:b])

	return &response, nil
}
