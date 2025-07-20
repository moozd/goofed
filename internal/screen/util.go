package screen

func slice[T any](s []T, start, end int) []T {
	if start < 0 {
		start = 0
	}
	if start >= len(s) {
		return []T{}
	}
	if end > len(s) {
		end = len(s)
	}
	if end < start {
		return []T{}
	}
	return s[start:end]
}
