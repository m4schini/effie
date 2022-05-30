package util

func GetNumberPostfix(n int) string {
	switch n {
	case 3:
		return "rd"
	case 2:
		return "nd"
	case 1:
		return "st"
	default:
		return "th"
	}
}

func SliceContains[T any](s []T, p func(T) bool) bool {
	for _, t := range s {
		if p(t) {
			return true
		}
	}

	return false
}
