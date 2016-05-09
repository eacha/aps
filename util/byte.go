package util

func CopySliceInto(source, target []byte, pos int) int {
	for i := 0; i < len(source); i++ {
		target[pos+i] = source[i]
	}

	return pos + len(source)
}

func ByteIndexOf(slice []byte, b byte, fromIndex int) int {
	if len(slice) <= fromIndex {
		return -1
	}

	for i := fromIndex; i < len(slice); i++ {
		if slice[i] == b {
			return i
		}
	}

	return -1
}
