package dns

import (
	"time"

	"github.com/eacha/aps/tools/conn"
)

type DNSConn struct {
	conn *conn.ConnTimeout
}

func NewDNSConn(protocol, address string, port int, connectionTimeout, ioTimeout time.Duration) (*DNSConn, error) {
	var dnsConn DNSConn
	var err error

	dnsConn.conn, err = conn.NewConnTimeout(protocol, address, port, connectionTimeout, ioTimeout)
	if err != nil {
		return nil, err
	}

	return &dnsConn, nil
}

func (c *DNSConn) OpenResolver(question, expected string) (*OpenResolver, error) {
	var data OpenResolver
	buf := make([]byte, bufferSize)

	_, err := c.conn.Write(NewQuery(question, recursiveDesired).Pack())
	if err != nil {
		return nil, err
	}

	b, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	data.RawResponse = buf[:b]

	var response Response
	if err = response.UnPack(buf[:b]); err != nil {
		return nil, err
	}

	data.Questions = response.Question
	data.Answers = response.Answer
	data.RecursionAvailable = ((response.Header.Bits >> _RA) & 0x1) == 1
	data.ResolveCorrectly = resolveCorrectly(data.Answers, expected)

	return &data, nil
}

func resolveCorrectly(ans []Answer, expected string) bool {
	for _, value := range ans {
		if value.RdData == expected {
			return true
		}
	}
	return false
}

func (c *DNSConn) Close() {
	c.conn.Close()
}
