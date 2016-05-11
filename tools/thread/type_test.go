package thread

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ThreadSuite struct{}

var _ = Suite(&ThreadSuite{})

func simpleThread(threadId, readLines, sleepTime int, c chan Statistic) {
	ts := NewStatistic(threadId)

	for i := 0; i < readLines; i++ {
		ts.IncreaseProcessedLines()
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	}

	ts.SetEndTime()
	c <- *ts

}

func (s *ThreadSuite) TestSimple(c *C) {
	channel := make(chan Statistic)

	go simpleThread(1, 1, 100, channel)
	ts := <-channel

	c.Assert(ts.ThreadId, Equals, 1)
	c.Assert(ts.ProcessedLines, Equals, 1)
}

func (s *ThreadSuite) TestMultiRead(c *C) {
	channel := make(chan Statistic)

	go simpleThread(1, 10, 100, channel)
	ts := <-channel

	c.Assert(ts.ThreadId, Equals, 1)
	c.Assert(ts.ProcessedLines, Equals, 10)
	c.Assert(1000 <= ts.DeltaTime, Equals, true)
	c.Assert(1010 >= ts.DeltaTime, Equals, true)
}

func (s *ThreadSuite) TestMultiThreads(c *C) {
	channel := make(chan Statistic)

	go simpleThread(1, 5, 100, channel)
	go simpleThread(2, 10, 100, channel)
	go simpleThread(3, 15, 100, channel)

	ts := <-channel

	c.Assert(ts.ThreadId, Equals, 1)
	c.Assert(ts.ProcessedLines, Equals, 5)
	c.Assert(500 <= ts.DeltaTime, Equals, true)
	c.Assert(510 >= ts.DeltaTime, Equals, true)

	ts = <-channel

	c.Assert(ts.ThreadId, Equals, 2)
	c.Assert(ts.ProcessedLines, Equals, 10)
	c.Assert(1000 <= ts.DeltaTime, Equals, true)
	c.Assert(1010 >= ts.DeltaTime, Equals, true)

	ts = <-channel

	c.Assert(ts.ThreadId, Equals, 3)
	c.Assert(ts.ProcessedLines, Equals, 15)
	c.Assert(1500 <= ts.DeltaTime, Equals, true)
	c.Assert(1510 >= ts.DeltaTime, Equals, true)
}
