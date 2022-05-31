package riot

import "strings"

const (
	QueueRankedFlex   = "RANKED_FLEX_SR"
	QueueRankedFlexId = 440
	QueueRankedSolo   = "RANKED_SOLO_5x5"
	QueueRankedSoloId = 420
)

func ToQueueConfigId(queueType string) int {
	switch queueType {
	case QueueRankedFlex:
		return QueueRankedFlexId
	case QueueRankedSolo:
		return QueueRankedSoloId
	default:
		return -1
	}
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
