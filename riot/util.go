package riot

const (
	envNameRiotApiKey    = "RIOT_API_KEY"
	envNameRiotApiRegion = "RIOT_API_REGION"

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
