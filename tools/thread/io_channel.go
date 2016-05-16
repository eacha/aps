package thread

import (
	"bufio"
	"log"
	"os"
	"time"
)

func ReadChannel(fileName string, read chan string) {
	var (
		err    error
		file   *os.File
		reader *bufio.Reader
	)

	switch fileName {
	case "-":
		file = os.Stdin
	default:
		file, err = os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()

	reader = bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			close(read)
			return
		}
		read <- string(line)
	}

}

func WriteChannel(fileName string, write chan string, finish chan int) {
	var (
		err  error
		file *os.File
	)

	switch fileName {
	case "-":
		file = os.Stdout
	default:
		file, err = os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()

	for {
		select {
		case line := <- write:
			file.WriteString(line + "\n")
		case <- finish:
			for {
				select {
				case line := <- write:
					file.WriteString(line + "\n")
				case <-time.After(time.Second * 1):
					return
				}
			}
		}
	}
}
