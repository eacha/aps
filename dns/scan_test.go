package dns

import (
	"testing"

	"github.com/eacha/aps/scan"
	"github.com/eacha/aps/tools/conn"
	"github.com/eacha/aps/tools/test"
	"github.com/eacha/aps/tools/thread"
	. "gopkg.in/check.v1"
	"sync"
	"time"
)

func TestScan(t *testing.T) { TestingT(t) }

type DNSScan struct{}

var _ = Suite(&DNSScan{})

var stringResponse = "{\"open_resolveer\":{\"questions\":[{\"type\":\"A\",\"class\":\"INET\",\"name\":\"www.ble.cl.\"}],\"answer\":[{\"type\":\"CNAME\",\"class\":\"INET\",\"name\":\"www.ble.cl.\",\"ttl\":12012,\"rd_length\":8,\"rd_data\":\"ble.cl.\"},{\"type\":\"A\",\"class\":\"INET\",\"name\":\"ble.cl.\",\"ttl\":12012,\"rd_length\":4,\"rd_data\":\"66.147.244.194\"}],\"recursion_available\":false,\"resolve_correctly\":true,\"raw_response\":\"AAEBAAABAAIAAAAAA3d3dwNibGUCY2wAAAEAAcAMAAUAAQAALuwACANibGUCY2wAwCgAAQABAAAu7AAEQpP0wg==\"}}"

func (dns *DNSScan) TestHostScanConnError(c *C) {
	options := scan.Options{
		Port:              53,
		ConnectionTimeout: 1,
		IOTimeout:         1,
	}

	data := hostScan(&options, "qwerty")

	c.Assert(data.IP, Equals, "qwerty")
	c.Assert(data.Error, Equals, "Connection refued by host, Host: qwerty")
}

func (dns *DNSScan) TestHostScanOpenResolverError(c *C) {
	var (
		wg      sync.WaitGroup
		options = scan.Options{
			Protocol:          conn.TCP,
			ConnectionTimeout: 1,
			IOTimeout:         1,
		}
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)

		data := hostScan(&options, "")

		c.Assert(data.Error, Equals, "DNS packet too short")
	}()

	test.ServerWrite(&options.Port, []byte{1}, 250)

	wg.Wait()
}

func (dns *DNSScan) TestHostScan(c *C) {
	var (
		wg      sync.WaitGroup
		options = scan.Options{
			Protocol:          conn.TCP,
			ConnectionTimeout: 1,
			IOTimeout:         1,
			DNSOptions: scan.DNSOptions{
				QuestionURL: "www.ble.cl",
				IpResponse:  "66.147.244.194",
			},
		}
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)

		data := hostScan(&options, "")

		c.Assert(data.OpenResolver.Questions[0].Qclass, Equals, classINET)
		c.Assert(data.OpenResolver.Questions[0].Qtype, Equals, typeA)
		c.Assert(data.OpenResolver.Questions[0].Qname, Equals, "www.ble.cl.")

		c.Assert(data.OpenResolver.Answers[1].Aclass, Equals, classINET)
		c.Assert(data.OpenResolver.Answers[1].Atype, Equals, typeA)
		c.Assert(data.OpenResolver.Answers[1].Aname, Equals, "ble.cl.")
		c.Assert(data.OpenResolver.Answers[1].RdLength, Equals, uint16(4))
		c.Assert(data.OpenResolver.Answers[1].RdData, Equals, "66.147.244.194")

		c.Assert(data.OpenResolver.RecursionAvailable, Equals, false)
		c.Assert(data.OpenResolver.ResolveCorrectly, Equals, true)
	}()

	test.ServerWrite(&options.Port, packResponse, 250)

	wg.Wait()
}

func (dns *DNSScan) TestScan(c *C) {
	var (
		wg      sync.WaitGroup
		options = scan.Options{
			InputChan:         make(chan string, 1),
			OutputChan:        make(chan string, 1),
			Protocol:          conn.TCP,
			ConnectionTimeout: 1,
			IOTimeout:         1,
			DNSOptions: scan.DNSOptions{
				QuestionURL: "www.ble.cl",
				IpResponse:  "66.147.244.194",
			},
		}
		stat = thread.NewStatistic(1)
	)
	wg.Add(1)
	options.WaitGroup = &wg
	go func() {
		time.Sleep(500 * time.Millisecond)
		options.InputChan <- ""

		Scan(&options, stat)

		c.Assert(<-options.OutputChan, Equals, stringResponse)
		c.Assert(stat.ProcessedLines, Equals, 1)
	}()

	test.ServerWrite(&options.Port, packResponse, 250)
	close(options.InputChan)

	wg.Wait()

}
