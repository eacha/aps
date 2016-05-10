package scan

import (
	"github.com/eacha/aps/tools/thread"
)

type ScanOptions struct {
	// Basic Scan Setup
	InputFileName string
	OutputFileName string

	InputFile *thread.SyncRead
	OutputFile *thread.SyncWrite

	Port uint
	Module string
	Threads uint
	ConnectionTimeout uint
	IOTimeout         uint

	// More options in the future
}

type Scannable interface {
	Scan(name int, options ScanOptions, statistic *thread.ThreadStatistic)
}
