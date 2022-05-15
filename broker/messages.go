package broker

import (
	"github.com/KnutZuidema/golio/riot/lol"
)

type MatchUpdate struct {
	SummonerId string
	Game       *lol.GameInfo
}

type MatchPostData struct {
	GameId   int
	Summoner *lol.Participant
	Match    *lol.Match
}
