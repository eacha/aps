package dns

import (
	"time"

	"github.com/eacha/aps/tools/conn"
)

// Conn is a specific connection to perform dns queries
type Conn struct {
	conn *conn.ConnTimeout
}

// NewDNSConn connect to the address and port on the named network. Take a connectionTimeout to perform a connection
// and set ioTimeout to read and write operations.
// It return the connection and any connection error encountered.
func NewDNSConn(protocol, address string, port int, connectionTimeout, ioTimeout time.Duration) (*Conn, error) {
	var dnsConn Conn
	var err error

	dnsConn.conn, err = conn.NewConnTimeout(protocol, address, port, connectionTimeout, ioTimeout)
	if err != nil {
		return nil, err
	}

	return &dnsConn, nil
}

// OpenResolver send a dns query who contains the specify question to resolve and check if the response is equals to
// expected. if positive the server has a open dns resolver.
// It return a Open Resolver struct with the data and any error encountered.
func (c *Conn) OpenResolver(question, expected string) (*OpenResolver, error) {
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
	if err = response.Unpack(buf[:b]); err != nil {
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

// Close closes the connection.
func (c *Conn) Close() {
	c.conn.Close()
}
