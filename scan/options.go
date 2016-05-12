package scan

import (
	"sync"

	"github.com/eacha/aps/tools/thread"
)

type Scannable interface {
	Scan(options *Options, statistic *thread.Statistic)
}

type DNSOptions struct {
	QuestionURL string
	IpResponse  string
}

type Options struct {
	// Basic Scan Setup
	WaitGroup *sync.WaitGroup

	InputFileName  string
	OutputFileName string

	InputChan chan string
	OutputChan chan string

	Port   uint
	Module string
	//Protocol string
	Threads           uint
	ConnectionTimeout uint
	IOTimeout         uint

	// DNS options
	DNSOptions DNSOptions

	// More options in the future
}
