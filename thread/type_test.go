package thread

import (
	"testing"
	"time"
)

func simpleThread(threadId, readLines, sleepTime int, c chan ThreadStatistic) {
	ts := NewThreadStadistic(threadId)

	for i := 0; i < readLines; i++ {
		ts.IncreaseProcessedLines()
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	}

	ts.SetEndTime()
	c <- *ts

}

func TestSimple(t *testing.T) {
	c := make(chan ThreadStatistic)

	go simpleThread(1, 1, 100, c)
	ts := <-c

	if ts.ThreadId != 1 {
		t.Error("Expected Thread Id 1, got", ts.ThreadId)
	}

	if ts.ProcessedLines != 1 {
		t.Error("Expected Read Lines 1, got", ts.ProcessedLines)
	}
}

func TestMultiRead(t *testing.T) {
	c := make(chan ThreadStatistic)

	go simpleThread(1, 10, 100, c)
	ts := <-c

	if ts.ThreadId != 1 {
		t.Error("Expected Thread Id 1, got", ts.ThreadId)
	}

	if ts.ProcessedLines != 10 {
		t.Error("Expected Read Lines 10, got", ts.ProcessedLines)
	}

	if 1000 > ts.DeltaTime && 1010 < ts.DeltaTime {
		t.Error("Expected Delta Time 1000, got", ts.ThreadId)
	}
}

func TestMultiThreads(t *testing.T) {
	c := make(chan ThreadStatistic)

	go simpleThread(1, 5, 100, c)
	go simpleThread(2, 10, 100, c)
	go simpleThread(3, 15, 100, c)

	ts := <-c

	if ts.ThreadId != 1 {
		t.Error("Expected Thread Id 1, got", ts.ThreadId)
	}

	if ts.ProcessedLines != 5 {
		t.Error("Expected Read Lines 5, got", ts.ProcessedLines)
	}

	if 500 > ts.DeltaTime && 510 < ts.DeltaTime {
		t.Error("Expected Delta Time 500, got", ts.ThreadId)
	}

	ts = <-c

	if ts.ThreadId != 2 {
		t.Error("Expected Thread Id 2, got", ts.ThreadId)
	}

	if ts.ProcessedLines != 10 {
		t.Error("Expected Read Lines 10, got", ts.ProcessedLines)
	}

	if 1000 > ts.DeltaTime && 1010 < ts.DeltaTime {
		t.Error("Expected Delta Time 1000, got", int(ts.DeltaTime))
	}

	ts = <-c

	if ts.ThreadId != 3 {
		t.Error("Expected Thread Id 3, got", ts.ThreadId)
	}

	if ts.ProcessedLines != 15 {
		t.Error("Expected Read Lines 15, got", ts.ProcessedLines)
	}

	if 1500 > ts.DeltaTime && 1510 < ts.DeltaTime {
		t.Error("Expected Delta Time 1500, got", int(ts.DeltaTime))
	}
}
