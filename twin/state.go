package twin

import (
	"github.com/KnutZuidema/golio/riot/lol"
	"sync"
)

type GameState uint8

const (
	Idle GameState = iota
	Started
	InGame
	Stopped
)

const (
	startState = Idle
)

func nop(info *lol.GameInfo) {}

type GameStateMachine interface {
	OnFoundData(info *lol.GameInfo)
	OnNoData()
	GetState() GameState
	Reset()

	SetOnIdle(func(info *lol.GameInfo))
	SetOnStarted(func(info *lol.GameInfo))
	SetOnInGame(func(info *lol.GameInfo))
	SetOnStopped(func(info *lol.GameInfo))
}

type match struct {
	summonerId string
	state      GameState
	cachedInfo *lol.GameInfo
	mu         sync.Mutex

	__onIdle    func(info *lol.GameInfo)
	__onStarted func(info *lol.GameInfo)
	__onInGame  func(info *lol.GameInfo)
	__onStopped func(info *lol.GameInfo)
}

func NewState(summonerId string) *match {
	m := new(match)
	m.summonerId = summonerId
	m.state = startState
	m.__onIdle = nop
	m.__onStarted = nop
	m.__onInGame = nop
	m.__onStopped = nop
	return m
}

func (m *match) SetOnIdle(f func(info *lol.GameInfo)) {
	m.mu.Lock()
	m.__onIdle = f
	m.mu.Unlock()
}

func (m *match) SetOnStarted(f func(info *lol.GameInfo)) {
	m.mu.Lock()
	m.__onStarted = f
	m.mu.Unlock()
}

func (m *match) SetOnInGame(f func(info *lol.GameInfo)) {
	m.mu.Lock()
	m.__onInGame = f
	m.mu.Unlock()
}

func (m *match) SetOnStopped(f func(info *lol.GameInfo)) {
	m.mu.Lock()
	m.__onStopped = f
	m.mu.Unlock()
}

func (m *match) update(info *lol.GameInfo) {
	m.cachedInfo = info
}

func (m *match) setState(state GameState, info *lol.GameInfo) {
	m.mu.Lock()
	m.state = state
	m.update(info)
	m.mu.Unlock()

	getGameId := func() int {
		if info != nil {
			return info.GameID
		} else {
			return -1
		}
	}

	switch state {
	case Idle:
		log.Debugw("inGame status:    idle",
			"summonerId", m.summonerId,
		)
		m.__onIdle(info)
		break
	case Started:
		log.Debugw("inGame status: started",
			"summonerId", m.summonerId,
			"gameId", getGameId(),
		)
		m.__onStarted(info)
		break
	case InGame:
		log.Debugw("inGame status:  active",
			"summonerId", m.summonerId,
			"gameId", getGameId(),
		)
		m.__onInGame(info)
		break
	case Stopped:
		log.Debugw("inGame status: stopped",
			"summonerId", m.summonerId,
			"gameId", getGameId(),
		)
		m.__onStopped(info)
		break
	}
}

func (m *match) GetState() GameState {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.state
}

func (m *match) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cachedInfo = nil
	m.state = startState
}

func (m *match) OnFoundData(info *lol.GameInfo) {
	switch m.state {
	case Idle:
		// idle -> started
		m.setState(Started, info)
		break
	case Started:
		// started -> InGame
		m.setState(InGame, info)
		break
	case InGame:
		// InGame -> InGame
		m.setState(InGame, info)
		break
	case Stopped:
		// Stopped -> Started
		m.setState(Started, info)
		break
	}
}

func (m *match) OnNoData() {
	switch m.state {
	case Idle:
		// Idle -> Idle
		m.setState(Idle, m.cachedInfo)
		break
	case Started:
		// Started -> Stopped
		m.setState(Stopped, m.cachedInfo)
		break
	case InGame:
		// InGame -> Stopped
		m.setState(Stopped, m.cachedInfo)
		break
	case Stopped:
		// Stopped -> Idle
		m.setState(Idle, m.cachedInfo)
		break
	}
}
