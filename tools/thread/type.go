package thread

import "time"

type Statistic struct {
	ThreadId       int           `json:"thread_id"`
	ProcessedLines int           `json:"processed_lines"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	DeltaTime      time.Duration `json:"delta_time"`
}

func NewStatistic(threadId int) *Statistic {
	var ts Statistic

	ts.ThreadId = threadId
	ts.ProcessedLines = 0
	ts.StartTime = time.Now()

	return &ts
}

func (ts *Statistic) IncreaseProcessedLines() {
	ts.ProcessedLines += 1
}

func (ts *Statistic) SetEndTime() {
	ts.EndTime = time.Now()
	ts.DeltaTime = ts.EndTime.Sub(ts.StartTime) / time.Millisecond
}
