package dns

import (
	"github.com/eacha/aps/tools/conn"
	. "gopkg.in/check.v1"

	"log"
	"sync"
	"time"

	"testing"

	"github.com/eacha/aps/tools/test"
)

func TestConn(t *testing.T) { TestingT(t) }

type DNSConnection struct{}

var _ = Suite(&DNSConnection{})

func (dns *DNSConnection) TestConnectionError(c *C) {
	conn, err := NewDNSConn(conn.UDP, "qwer", 53, 1, 1)

	c.Assert(conn, IsNil)
	c.Assert(err, NotNil)
}

func (dns *DNSConnection) TestConnectionSuccess(c *C) {
	var (
		wc   sync.WaitGroup
		port int
	)

	wc.Add(1)
	go func() {
		defer wc.Done()
		time.Sleep(500 * time.Millisecond)

		conn, err := NewDNSConn(conn.TCP, "", port, 1, 1)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer conn.Close()

		c.Assert(conn, NotNil)
		c.Assert(err, IsNil)
	}()

	test.ConnectionServer(conn.TCP, &port, 250)

	wc.Wait()
}

func (dns *DNSConnection) TestOpenResolverWriteError(c *C) {
	var (
		wc   sync.WaitGroup
		port int
	)

	wc.Add(1)
	go func() {
		defer wc.Done()
		time.Sleep(500 * time.Millisecond)

		conn, err := NewDNSConn(conn.TCP, "", port, 1, 1)
		if err != nil {
			log.Fatal(err.Error())
		}
		conn.Close()

		data, err := conn.OpenResolver("test", "test")
		c.Assert(data, IsNil)
		c.Assert(err, NotNil)
	}()

	test.ConnectionServer(conn.TCP, &port, 250)

	wc.Wait()
}

func (dns *DNSConnection) TestOpenResolverReadError(c *C) {
	var (
		wc   sync.WaitGroup
		port int
	)

	wc.Add(1)
	go func() {
		defer wc.Done()
		time.Sleep(500 * time.Millisecond)

		conn, err := NewDNSConn(conn.TCP, "", port, 1, 1)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer conn.Close()

		data, err := conn.OpenResolver("test", "test")
		c.Assert(data, IsNil)
		c.Assert(err, NotNil)
	}()

	test.ServerReadAndClose(&port, 0)

	wc.Wait()
}

func (dns *DNSConnection) TestOpenResolverBufferError(c *C) {
	var (
		wc   sync.WaitGroup
		port int
	)

	wc.Add(1)
	go func() {
		defer wc.Done()
		time.Sleep(500 * time.Millisecond)

		conn, err := NewDNSConn(conn.TCP, "", port, 1, 1)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer conn.Close()

		data, err := conn.OpenResolver("test", "test")
		c.Assert(data, IsNil)
		c.Assert(err, NotNil)
	}()

	test.ServerWrite(&port, []byte{1}, 250)

	wc.Wait()
}

func (dns *DNSConnection) TestOpenResolver(c *C) {
	var (
		wc   sync.WaitGroup
		port int
	)

	wc.Add(1)
	go func() {
		defer wc.Done()
		time.Sleep(500 * time.Millisecond)

		conn, err := NewDNSConn(conn.TCP, "", port, 1, 1)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer conn.Close()

		data, _ := conn.OpenResolver("ble.cl", "66.147.244.194")

		c.Assert(data.Questions[0].Qclass, Equals, classINET)
		c.Assert(data.Questions[0].Qtype, Equals, typeA)
		c.Assert(data.Questions[0].Qname, Equals, "www.ble.cl.")

		c.Assert(data.Answers[1].Aclass, Equals, classINET)
		c.Assert(data.Answers[1].Atype, Equals, typeA)
		c.Assert(data.Answers[1].Aname, Equals, "ble.cl.")
		c.Assert(data.Answers[1].RdLength, Equals, uint16(4))
		c.Assert(data.Answers[1].RdData, Equals, "66.147.244.194")

		c.Assert(data.RecursionAvailable, Equals, false)
		c.Assert(data.ResolveCorrectly, Equals, true)
	}()

	test.ServerWrite(&port, packResponse, 250)

	wc.Wait()
}

func (dns *DNSConnection) TestResolveCorrectly(c *C) {
	answerCorrect := Answer{RdData: "192.0.0.1"}
	answerFail := Answer{RdData: "qwerty"}

	c.Assert(resolveCorrectly([]Answer{answerCorrect}, "192.0.0.1"), Equals, true)
	c.Assert(resolveCorrectly([]Answer{answerFail}, "192.0.0.1"), Equals, false)

}
