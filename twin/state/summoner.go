package state

import "github.com/KnutZuidema/golio/riot/lol"

type SummonerStatus interface {
	OnFoundData(gameInfo *lol.GameInfo)
	OnNoData()
	State() bool
}
