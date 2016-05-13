package dns

import (
	"encoding/json"
	"github.com/eacha/aps/scan"
	"github.com/eacha/aps/tools/connection"
	"github.com/eacha/aps/tools/thread"
)

func Scan(options *scan.Options, statistic *thread.Statistic) {
	defer options.WaitGroup.Done()
	for {
		address, more := <-options.InputChan
		if !more {
			break
		}
		statistic.IncreaseProcessedLines()

		data := hostScan(options, address)
		j, _ := json.Marshal(data)

		options.OutputChan <- string(j)
	}
	statistic.SetEndTime()
}

func hostScan(options *scan.Options, address string) DNSData {
	var dnsData DNSData

	dnsData.ip = address
	conn, err := NewDNSConn(connection.UDP, address, options.Port, options.ConnectionTimeout, options.IOTimeout)
	if err != nil {
		dnsData.error = err
		return dnsData
	}
	defer conn.Close()

	dnsData.openResolver, err = conn.OpenResolver(options.DNSOptions.QuestionURL, options.DNSOptions.IpResponse)
	if err != nil {
		dnsData.error = err
		return dnsData
	}

	return dnsData
}
