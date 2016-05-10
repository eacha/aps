package scan

import (
	"github.com/eacha/aps/tools/thread"
	"sync"
)

type Scannable interface {
	Scan(name int, options ScanOptions, statistic *thread.ThreadStatistic)
}

type DNSOptions struct {
	QuestionURL string
	IpResponse  string
}

type ScanOptions struct {
	// Basic Scan Setup
	WaitGroup *sync.WaitGroup

	InputFileName string
	OutputFileName string

	InputFile *thread.SyncRead
	OutputFile *thread.SyncWrite

	Port uint
	Module string
	//Protocol string
	Threads uint
	ConnectionTimeout uint
	IOTimeout         uint

	// DNS options
	DNSOptions DNSOptions

	// More options in the future
}


