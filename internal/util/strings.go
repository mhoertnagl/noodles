package util

func IndexOf(vals []string, val string) int {
	for pos := 0; pos < len(vals); pos++ {
		if vals[pos] == val {
			return pos
		}
	}
	return -1
}
