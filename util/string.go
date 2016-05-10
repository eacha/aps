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

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
