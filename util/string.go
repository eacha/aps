package util

import "strconv"

const nullTerminated = 0x00

func Atoi(number string) int {
	n, err := strconv.Atoi(number)

	if err != nil {
		return 0
	}

	return n
}

func ToByteArrayNullTerminated(str string) []byte {
	return append([]byte(str), []byte{nullTerminated}...)
}
