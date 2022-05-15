package twin

import (
	"github.com/KnutZuidema/golio/riot/lol"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

type GameManager struct {
	games      map[int]*Game
	gamesMutex sync.Mutex
}

func NewManager() *GameManager {
	gm := new(GameManager)
	gm.games = make(map[int]*Game)
	return gm
}

func (gm *GameManager) Add(game *Game) bool {
	defer gm.gamesMutex.Unlock()
	gm.gamesMutex.Lock()

	_, exists := gm.games[game.id]
	if exists {
		return false
	}

	gm.games[game.id] = game
	return true
}

type Game struct {
	id           int
	participants []string
	startTime    time.Time
	mu           sync.Mutex
}

func NewGame(gameId int) *Game {
	g := new(Game)
	g.id = gameId
	g.participants = make([]string, 0, 5)
	return g
}

func NewGameFromGameInfo(info *lol.GameInfo) *Game {
	g := NewGame(info.GameID)
	g.startTime = time.UnixMilli(int64(info.GameStartTime))
	return g
}

func (g *Game) GameLength() time.Duration {
	return time.Since(g.startTime)
}

func (g *Game) AddParticipant(summonerId string) bool {
	defer g.mu.Unlock()
	g.mu.Lock()

	if slices.Contains(g.participants, summonerId) {
		return false
	}

	g.participants = append(g.participants, summonerId)
	return true
}
