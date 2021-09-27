package helper

func InterceptString(s string, size int, padding string) string {
	ss := []rune(s) // not []byte
	if len(ss) <= size {
		return s
	}
	return string(ss[:size]) + padding
}
