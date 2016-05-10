package thread

import "time"

type ThreadStatistic struct {
	ThreadId       int           `json:"thread_id"`
	ProcessedLines int           `json:"processed_lines"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	DeltaTime      time.Duration `json:"delta_time"`
}

func NewThreadStadistic(threadId int) *ThreadStatistic {
	var ts ThreadStatistic

	ts.ThreadId = threadId
	ts.ProcessedLines = 0
	ts.StartTime = time.Now()

	return &ts
}

func (ts *ThreadStatistic) IncreaseProcessedLines() {
	ts.ProcessedLines += 1
}

func (ts *ThreadStatistic) SetEndTime() {
	ts.EndTime = time.Now()
	ts.DeltaTime = ts.EndTime.Sub(ts.StartTime) / time.Millisecond
}
