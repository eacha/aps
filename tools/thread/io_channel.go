package thread

import (
	"bufio"
	"log"
	"os"
)

//type SyncRead struct {
//	mutex  sync.Mutex
//	file   *os.File
//	reader *bufio.Reader
//	finish bool
//}
//
//func NewSyncRead(fileName string) *SyncRead {
//	var sr SyncRead
//	var err error
//
//	switch fileName {
//	case "-":
//		sr.file = os.Stdin
//	default:
//		sr.file, err = os.Open(fileName)
//	}
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	sr.reader = bufio.NewReader(sr.file)
//	sr.finish = false
//
//	return &sr
//}
//
//func (sr *SyncRead) ReadLine() []byte {
//	sr.mutex.Lock()
//	if sr.finish {
//		return nil
//	}
//
//	line, _, err := sr.reader.ReadLine()
//	if err != nil {
//		sr.finish = true
//		return nil
//	}
//	sr.mutex.Unlock()
//
//	return line
//}
//
//func (sr *SyncRead) Close() error {
//	return sr.file.Close()
//}
//
//type SyncWrite struct {
//	mutex sync.Mutex
//	file  *os.File
//}
//
//func NewSyncWrite(fileName string) *SyncWrite {
//	var sw SyncWrite
//	var err error
//
//	switch fileName {
//	case "-":
//		sw.file = os.Stdout
//	default:
//		sw.file, err = os.Create(fileName)
//	}
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return &sw
//}
//
//func (sw *SyncWrite) WriteLine(b []byte) (int, error) {
//	sw.mutex.Lock()
//	writeBytes, err := sw.file.Write(b)
//	sw.mutex.Unlock()
//
//	return writeBytes, err
//}
//
//func (sw *SyncWrite) Close() error {
//	return sw.file.Close()
//}

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
			file.WriteString(line)
		case <- finish:
			return
		}
	}
}
