package twin

import (
	"context"
	"effie3/broker"
	"effie3/cache"
	"effie3/logger"
	"effie3/riot"
	"github.com/KnutZuidema/golio/riot/lol"
	"time"
)

var log = logger.Get("twin").Sugar()

const (
	RefreshInterval = 30 * time.Second
)

type SummonerTwin interface {
	ID() string
	InGame() (bool, *lol.GameInfo)
	Refresh() error
	AutoRefresh(ctx context.Context)
}

type summoner struct {
	id          string
	inGameState GameStateMachine
	IsInGame    bool
}

func NewSummonerTwin(id string) *summoner {
	s := new(summoner)
	s.id = id
	s.inGameState = NewState(id)

	s.inGameState.SetOnStarted(s.onGameStarted)
	s.inGameState.SetOnInGame(s.onGameUpdated)
	s.inGameState.SetOnStopped(s.onGameStopped)

	log.Debugw("created new summoner twin", "summonerId", id)
	return s
}

func (s *summoner) ID() string {
	return s.id
}

func (s *summoner) InGame() (bool, *lol.GameInfo) {
	info, err := riot.GetCurrentGame(s.id)
	if err != nil {
		if err.Error() != "not found" {
			log.Errorw("error while fetching summoner game status",
				"error", err)
			return false, nil
		}
	}

	log.Debugw("fetched summoner game status",
		"summonerId", s.id,
		"inGameState", info != nil)
	return info != nil, info
}

func (s *summoner) Refresh() error {
	inGame, info := s.InGame()

	if inGame {
		s.inGameState.OnFoundData(info)
	} else {
		s.inGameState.OnNoData()
	}

	log.Infow("twin updated",
		"summonerName", cache.GetSummonerName(s.id),
		"summonerId", s.id,
		"inGame", inGame)
	return nil
}

func (s *summoner) AutoRefresh(ctx context.Context) {
	ch := make(chan struct{})

	select {
	case <-ch:
		for {
			s.Refresh()
			log.Debugw("next twin refresh scheduled",
				"summonerId", s.id,
				"scheduledAt", time.Now().Add(RefreshInterval),
			)
			time.Sleep(RefreshInterval)
		}
	case <-ctx.Done():
		log.Warnw("Auto-Refresh was canceled",
			"summonerId", s.id,
			"reason", ctx.Err())
		return
	}
}

func (s *summoner) onGameStarted(info *lol.GameInfo) {
	if info == nil {
		log.Warnw("tried to trigger onGameStarted with nil gameInfo",
			"summonerId", s.id)
		return
	}

	msg := &broker.MatchUpdate{SummonerId: s.id, Game: info}
	switch info.GameQueueConfigID {
	case riot.QueueRankedSoloId:
		broker.Publish("match.ranked.solo.started", msg)
		break
	case riot.QueueRankedFlexId:
		broker.Publish("match.ranked.flex.started", msg)
		break
	default:
		broker.Publish("match.started", msg)
		break
	}

	leagues, err := riot.GetLeagues(s.id)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugw("retrieved leagues of summoner",
		"summonerId", s.id,
		"count", len(leagues),
	)

	for _, league := range leagues {
		if riot.ToQueueConfigId(league.QueueType) == info.GameQueueConfigID {
			log.Infow("found league",
				"summonerId", league.SummonerID,
				"summonerName", league.SummonerName,
				"tier", league.Tier,
				"rank", league.Rank,
				"inPromos", league.MiniSeries != nil)

			if league.MiniSeries == nil {
				continue
			}

			switch info.GameQueueConfigID {
			case riot.QueueRankedSoloId:
				broker.Publish("match.promo.solo.started", msg)
				break
			case riot.QueueRankedFlexId:
				broker.Publish("match.promo.flex.started", msg)
				break
			default:
				break
			}
		}
	}
}

func (s *summoner) onGameUpdated(info *lol.GameInfo) {
	if info == nil {
		log.Warnw("tried to trigger onGameUpdated with nil gameInfo",
			"summonerId", s.id)
		return
	}

	msg := &broker.MatchUpdate{SummonerId: s.id, Game: info}
	broker.Publish("match.updated", msg)
}

func (s *summoner) onGameStopped(info *lol.GameInfo) {
	if info == nil {
		log.Warnw("tried to trigger onGameStopped with nil gameInfo",
			"summonerId", s.id)
		return
	}

	msg := &broker.MatchUpdate{SummonerId: s.id, Game: info}
	broker.Publish("match.ended", msg)
}
