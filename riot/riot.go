package riot

import (
	"effie3/conf"
	"effie3/logger"
	"effie3/util"
	"errors"
	"github.com/KnutZuidema/golio"
	golioApi "github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	"go.uber.org/zap"
)

var log = logger.Get("riot").Sugar()

type PromoSeries struct {
	Queue    string
	From     string
	To       string
	Progress string
	Target   int
	Wins     int
	Losses   int
}

type Api interface {
	GetSummonerByName(summonerName string) (*lol.Summoner, error)
	GetSummonerById(summonerId string) (*lol.Summoner, error)
	GetCurrentGame(summonerId string) (*lol.GameInfo, error)
	GetLeagues(summonerId string) ([]*lol.LeagueItem, error)
	GetPromos(summonerId string) ([]*PromoSeries, error)
	GetMatch(matchId string, summonerId string) (*lol.Match, *lol.Participant, error)
	GetGameLevel(summonerId string, info *lol.GameInfo) (conf.VolumeLevel, error)
}

type api struct {
	client *golio.Client
	logger *zap.SugaredLogger
}

func NewApi(region, apiKey string) (*api, error) {
	a := new(api)

	log.Debugw("creating new riot api client", "region", region)
	a.client = golio.NewClient(apiKey, golio.WithRegion(golioApi.Region(region)))

	a.logger = log.Named(region)
	return a, nil
}

func (a *api) GetSummonerByName(summonerName string) (*lol.Summoner, error) {
	s, err := a.client.Riot.LoL.Summoner.GetByName(summonerName)
	if err != nil {
		return nil, err
	}

	a.logger.Debugw("retrieved summoner info", "summonerId", s.ID)
	return s, nil
}

func (a *api) GetSummonerById(summonerId string) (*lol.Summoner, error) {
	s, err := a.client.Riot.LoL.Summoner.GetByID(summonerId)
	if err != nil {
		return nil, err
	}

	a.logger.Debugw("retrieved summoner info", "summonerId", s.ID)
	return s, nil
}

func (a *api) GetCurrentGame(summonerId string) (*lol.GameInfo, error) {
	s, err := a.client.Riot.LoL.Spectator.GetCurrent(summonerId)
	if err != nil {
		return nil, err
	}

	a.logger.Debugw("retrieved current game info", "summonerId", summonerId, "gameId", s.GameID)
	return s, nil
}

func (a *api) GetLeagues(summonerId string) ([]*lol.LeagueItem, error) {
	s, err := a.client.Riot.LoL.League.ListBySummoner(summonerId)
	if err != nil {
		return nil, err
	}

	a.logger.Debugw("retrieved leagues for summoner", "summonerId", summonerId)
	return s, nil
}

func (a *api) GetPromos(summonerId string) ([]*PromoSeries, error) {
	info, err := a.client.Riot.LoL.League.ListBySummoner(summonerId)
	if err != nil {
		return nil, err
	}

	promos := make([]*PromoSeries, 0)
	for _, item := range info {
		if item.MiniSeries != nil {
			from, to, err := a.getTier(summonerId, item.QueueType)
			if err != nil {
				*from = ""
				*to = ""
			}

			promos = append(promos, &PromoSeries{
				Queue:    item.QueueType,
				From:     *from,
				To:       *to,
				Progress: item.MiniSeries.Progress,
				Target:   item.MiniSeries.Target,
				Wins:     item.MiniSeries.Wins,
				Losses:   item.MiniSeries.Losses,
			})
		}
	}

	return promos, nil
}

//TODO avoid/replace
func (a *api) getTier(summonerId, queueType string) (*string, *string, error) {

	leagues, err := a.client.Riot.LoL.League.ListBySummoner(summonerId)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	for _, league := range leagues {
		if league.QueueType == queueType {

			next := util.NextTier(league.Tier)
			return &league.Tier, &next, nil
		}
	}

	return nil, nil, errors.New("queueType doesn't exist")
}

func (a *api) GetMatch(matchId string, summonerId string) (*lol.Match, *lol.Participant, error) {
	game, err := a.client.Riot.LoL.Match.Get(matchId)
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

//TODO avoid/replace
func (a *api) GetGameLevel(summonerId string, info *lol.GameInfo) (conf.VolumeLevel, error) {
	if info == nil {
		return conf.All, errors.New("no info supplied")
	}

	if info.GameQueueConfigID != QueueRankedSoloId && info.GameQueueConfigID != QueueRankedFlexId {
		return conf.All, nil
	}

	promos, err := a.GetPromos(summonerId)
	if err != nil {
		return conf.All, err
	}

	for _, promo := range promos {
		if ToQueueConfigId(promo.Queue) == info.GameQueueConfigID {
			return conf.Promo, nil
		}
	}

	return conf.Ranked, nil
}
