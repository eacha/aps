package dns

import (
	"github.com/eacha/aps/tools/connection"
	"time"
)

type DNSConn struct {
	conn *connection.ConnTimeout
}

func NewDNSConn(protocol, address string, port int, connectionTimeout, ioTimeout time.Duration) (*DNSConn, error) {
	var dnsConn DNSConn
	var err error

	dnsConn.conn, err = connection.NewConnTimeout(protocol, address, port, connectionTimeout, ioTimeout)
	if err != nil{
		return nil, err
	}

	return &dnsConn, nil
}

func (c *DNSConn)OpenResolver(question string) (string,  error){
	query := NewQuery(question, RecursiveDesired)
	buf := make([]byte, 1024)
	pack := query.Pack()

	_, err := c.conn.Write(pack)
	if err != nil {
		return "", err
	}

	b, err := c.conn.Read(buf)
	if err != nil {
		return "", err
	}

	var response Response
	response.UnPack(buf[:b])

	return "ip", nil
}