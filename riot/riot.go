package riot

import (
	"effie/logger"
	"effie/output"
	"errors"
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	"os"
	"strings"
)

var Api *golio.Client
var log = logger.Get("riot").Sugar()

// k: summonerName v: summonerId
var knownSummoners map[string]string

var online bool

const (
	envNameRiotApiKey    = "RIOT_API_KEY"
	envNameRiotApiRegion = "RIOT_API_REGION"

	QueueRankedFlex   = "RANKED_FLEX_SR"
	QueueRankedFlexId = 440
	QueueRankedSolo   = "RANKED_SOLO_5x5"
	QueueRankedSoloId = 420
)

func init() {
	online = true
	defer func() {
		if !online {
			log.Warnw("riot api is not available!")
		} else {
			log.Debugf("%s: %s", "SELECTED REGION", os.Getenv("RIOT_API_REGION"))
		}
	}()

	if os.Getenv(envNameRiotApiKey) == "" {
		online = false
		log.Warnw(envNameRiotApiKey+" is empty", "value", os.Getenv(envNameRiotApiKey))
	}
	if os.Getenv(envNameRiotApiRegion) == "" {
		online = false
		log.Warnw(envNameRiotApiRegion+" is empty", "value", os.Getenv(envNameRiotApiRegion))
	}

	if online {
		client := golio.NewClient(
			os.Getenv("RIOT_API_KEY"),
			golio.WithRegion(api.Region(os.Getenv("RIOT_API_REGION"))),
		)

		Api = client
	}
	knownSummoners = make(map[string]string)
}

type PromoSeries struct {
	Queue    string
	From     string
	To       string
	Progress string
	Target   int
	Wins     int
	Losses   int
}

func newPromoSeries(queue, from, to string, info *lol.MiniSeries) *PromoSeries {
	p := new(PromoSeries)
	p.Queue = queue
	p.From = from
	p.To = to
	p.Progress = info.Progress
	p.Target = info.Target
	p.Wins = info.Wins
	p.Losses = info.Losses
	return p
}

func GetSummonerByName(summonerName string) (*lol.Summoner, error) {
	if !online {
		return nil, errors.New("riot api not available")
	}
	s, err := Api.Riot.LoL.Summoner.GetByName(summonerName)
	if err != nil {
		return nil, err
	}
	knownSummoners[s.Name] = s.ID
	return s, err
}

func GetSummonerById(summonerId string) (*lol.Summoner, error) {
	if !online {
		return nil, errors.New("riot api not available")
	}
	s, err := Api.Riot.LoL.Summoner.GetByID(summonerId)
	if err != nil {
		return nil, err
	}
	knownSummoners[s.Name] = s.ID
	return s, err
}

func IsPlaying(summonerId string) (bool, *lol.GameInfo, error) {
	if !online {
		return false, nil, errors.New("riot api not available")
	}

	gameInfo, err := Api.Riot.LoL.Spectator.GetCurrent(summonerId)
	if err != nil && err.Error() == "not found" {
		return false, nil, nil
	}
	if err != nil {
		log.Error(err)
		return false, nil, nil
	}

	return true, gameInfo, nil
}

func getTier(summonerId, queueType string) (*string, *string, error) {
	if !online {
		return nil, nil, errors.New("riot api not available")
	}

	leagues, err := Api.Riot.LoL.League.ListBySummoner(summonerId)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	for _, league := range leagues {
		if league.QueueType == queueType {

			next := NextTier(league.Tier)
			return &league.Tier, &next, nil
		}
	}

	return nil, nil, errors.New("queueType doesn't exist")
}

func LeaguesBySummonerName(summonerName string) (map[int]*lol.LeagueItem, error) {
	if !online {
		return nil, errors.New("riot api not available")
	}

	id, ok := knownSummoners[summonerName]
	if !ok {
		return nil, errors.New("unknown summoner name")
	}

	info, err := Api.Riot.LoL.League.ListBySummoner(id)
	if err != nil {
		return nil, err
	}

	leagues := make(map[int]*lol.LeagueItem)

	for _, item := range info {
		queueId := ToQueueConfigId(item.QueueType)
		leagues[queueId] = item
	}

	return leagues, nil
}

func GetPromos(summonerId string) ([]*PromoSeries, error) {
	if !online {
		return nil, errors.New("riot api not available")
	}

	info, err := Api.Riot.LoL.League.ListBySummoner(summonerId)
	if err != nil {
		return nil, err
	}

	promos := make([]*PromoSeries, 0)
	for _, item := range info {
		if item.MiniSeries != nil {
			from, to, err := getTier(summonerId, item.QueueType)
			if err != nil {
				*from = ""
				*to = ""
			}

			promos = append(promos, newPromoSeries(item.QueueType, *from, *to, item.MiniSeries))
		}
	}

	return promos, nil
}

func GetMatch(matchId string, summonerId string) (*lol.Match, *lol.Participant, error) {
	if !online {
		return nil, nil, errors.New("riot api not available")
	}

	game, err := Api.Riot.LoL.Match.Get(matchId)
	if err != nil {
		return nil, nil, err
	}

	var userData *lol.Participant
	for _, participant := range game.Info.Participants {
		if participant.SummonerID == summonerId {
			userData = participant
			break
		}
	}

	return game, userData, nil
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

func GetGameLevel(summonerId string, info *lol.GameInfo) (output.Level, error) {
	if !online {
		return output.All, errors.New("riot api not available")
	}
	if info == nil {
		return output.All, errors.New("no info supplied")
	}

	if info.GameQueueConfigID != QueueRankedSoloId && info.GameQueueConfigID != QueueRankedFlexId {
		return output.All, nil
	}

	promos, err := GetPromos(summonerId)
	if err != nil {
		return output.All, err
	}

	for _, promo := range promos {
		if ToQueueConfigId(promo.Queue) == info.GameQueueConfigID {
			return output.Promo, nil
		}
	}

	return output.Ranked, nil
}
