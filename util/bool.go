package util

func BoolToByte(val bool) byte {
	if val {
		return 1
	} else {
		return 0
	}
}
