package thread

import (
	"bufio"
	"log"
	"os"
	"time"

	. "gopkg.in/check.v1"
)

type ChanSuite struct{}

var _ = Suite(&ChanSuite{})

var (
	inputName  = "read.txt"
	outputName = "write.txt"
)

func (s *ChanSuite) SetUpSuite(c *C) {
	file, err := os.Create(inputName)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString("1234\n")
	file.WriteString("4567\n")

	file.Close()
}

func (s *ChanSuite) TearDownSuite(c *C) {
	os.Remove(inputName)
	os.Remove(outputName)
}

func (s *ChanSuite) TestReadChannel(c *C) {
	read := make(chan string, 1)
	go ReadChannel(inputName, read)

	r, more := <-read
	c.Assert(r, Equals, "1234")
	c.Assert(more, Equals, true)

	r, more = <-read
	c.Assert(r, Equals, "4567")
	c.Assert(more, Equals, true)

	r, more = <-read
	c.Assert(more, Equals, false)
}

func (s *ChanSuite) TestWriteChannel(c *C) {
	write := make(chan string, 1)
	end := make(chan int, 1)
	go WriteChannel(outputName, write, end)

	write <- "1234"
	end <- 1
	write <- "4567"

	time.Sleep(time.Second * 2)

	file, err := os.Open(outputName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	line, _, err := reader.ReadLine()
	c.Assert(string(line), Equals, "1234")

	line, _, err = reader.ReadLine()
	c.Assert(string(line), Equals, "4567")
}
