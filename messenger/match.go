package messenger

import (
	"effie/broker"
	"effie/output"
	"github.com/KnutZuidema/golio/riot/lol"
)

type State uint8

const (
	Started State = iota
	Loading
	Active
	Ended
	PostGame
)

type matchState struct {
	state        State
	level        output.Level
	participants []string
	info         *lol.GameInfo
	post         *broker.MatchPostData
}

func (m *matchState) SetState(state State, postgame *broker.MatchPostData, game *Game, participants []string) {
	if game != nil && game.GameInfo != nil {
		m.info = game.GameInfo
	}
	if postgame != nil {
		m.post = postgame
	}
	m.participants = participants
	m.level = game.Level

	if state < m.state {
		log.Warnw("tried to set invalid state", "current", m.state, "next", state)
		return
	}

	m.state = state
}

func (m *matchState) GetMessage() string {

	switch m.state {
	case Started:
		return GetStartedMessage(m.info, m.participants, m.level)
	case Loading:
		return GetLoadingMessage(m.info, m.participants, m.level)
	case Active:
		return GetActiveMessage(m.info, m.participants, m.level)
	case Ended:
		return GetEndedMessage(m.info, m.participants, m.level)
	case PostGame:
		return GetPostMessage(m.post, m.participants, m.level)
	}

	return "something went wrong :("
}
