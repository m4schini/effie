package messenger

import (
	"effie3/bot"
	"effie3/broker"
	"effie3/riot"
	"effie3/values/volume"
	"fmt"
	"github.com/KnutZuidema/golio/riot/lol"
	"golang.org/x/exp/slices"
)

var games = make(map[int]*Game)

var closedGames = make([]int, 0)

func isClosed(gameId int) bool {
	return slices.Contains(closedGames, gameId)
}

type Game struct {
	State        *matchState
	Level        volume.Level
	GameId       int
	Participants []string
	GameInfo     *lol.GameInfo
	Message      bot.Message
}

func AddGame(discovered *broker.MatchUpdate, level volume.Level) error {
	gameId := discovered.Game.GameID
	summoner, err := riot.GetSummonerById(discovered.SummonerId)
	if err != nil {
		return err
	}

	g, exists := games[gameId]
	if exists || isClosed(gameId) {
		if exists && !slices.Contains(g.Participants, summoner.Name) {
			g.Participants = append(g.Participants, summoner.Name)
		}
		return fmt.Errorf("game (%d) can't be added", gameId)
	}

	game := &Game{
		State:        &matchState{state: Started},
		Level:        level,
		GameId:       gameId,
		Participants: make([]string, 0),
		GameInfo:     discovered.Game,
		Message:      nil,
	}

	if !slices.Contains(game.Participants, summoner.Name) {
		game.Participants = append(game.Participants, summoner.Name)
	}

	games[gameId] = game

	return nil
}

func GetGame(gameId int) (*Game, bool) {
	game, ok := games[gameId]
	return game, ok
}

func RemoveGame(gameId int) {
	delete(games, gameId)
}

//CloseGame removes and closes (blocking from getting added again) the game
func CloseGame(gameId int) {
	delete(games, gameId)

	closedGames = append(closedGames, gameId)
}

func SetState(gameId int, state State, info *broker.MatchUpdate, post *broker.MatchPostData) {
	game, ok := GetGame(gameId)
	if !ok {
		return
	}
	if info != nil {
		game.GameInfo = info.Game
	}

	game.State.SetState(state, post, game, game.Participants)
}

func SetDMessage(gameId int, message bot.Message) {
	game, ok := games[gameId]
	if !ok {
		return
	}

	game.Message = message
}

func GetDMessage(gameId int) bot.Message {
	game, ok := games[gameId]
	if !ok {
		return nil
	}

	return game.Message
}

func GetMessageString(gameId int) string {
	game, ok := games[gameId]
	if !ok {
		return ""
	}

	return game.State.GetMessage()
}
