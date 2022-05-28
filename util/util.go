package util

import "strings"

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

func NextTier(tier string) string {
	switch strings.ToUpper(tier) {
	case "IRON":
		return "BRONZE"
	case "BRONZE":
		return "SILVER"
	case "SILVER":
		return "GOLD"
	case "GOLD":
		return "PLATINUM"
	case "PLATINUM":
		return "DIAMOND"
	case "DIAMOND":
		return "MASTER"
	case "MASTER":
		return "GRANDMASTER"
	case "GRANDMASTER":
		return "CHALLENGER"
	default:
		return ""
	}
}
