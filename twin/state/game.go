package state

import "github.com/KnutZuidema/golio/riot/lol"

type GameData struct {
	Summoner *lol.Summoner
	Info     *lol.GameInfo
	Promo    *lol.MiniSeries
}

type GameState uint8

const (
	Started GameState = iota
	InGame
	Stopped
)

type GameStatus interface {
	OnData(gameInfo *GameData)
	OnNoData()
	State() GameState
	Data() *GameData
}
