package util

func CopySliceInto(source, target []byte, pos int) int {
	for i := 0; i < len(source); i++ {
		target[pos+i] = source[i]
	}

	return pos + len(source)
}
