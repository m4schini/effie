package twin

import (
	"encoding/json"
	"github.com/KnutZuidema/golio/riot/lol"
)

// Input
// info == nil: no game data found -> summoner not ingame
// info != nil: game data found -> summoner in an active game

type State uint8

const (
	Idle    State = iota // no data
	Started              // no data -> data found
	InGame               // data found -> data found
	Ended                // data found -> no data
)

type TransitionRule func(in *lol.GameInfo) State

type Handler func(info *lol.GameInfo)

type InGameStatus interface {
	Input(in *lol.GameInfo)
}

type inGameState struct {
	state       State
	data        *lol.GameInfo
	transitions map[State]TransitionRule
	handlers    map[State]Handler
}

func NewInGameStatus(onIdle, onStarted, onUpdate, onEnded Handler) *inGameState {
	igs := new(inGameState)
	igs.state = Idle
	igs.transitions = make(map[State]TransitionRule)

	igs.transitions[Idle] = transitionFromIdle
	igs.transitions[Started] = transitionFromStarted
	igs.transitions[InGame] = transitionFromInGame
	igs.transitions[Ended] = transitionFromEnded

	//TODO handle nil parameter
	igs.handlers[Idle] = onIdle
	igs.handlers[Started] = onStarted
	igs.handlers[InGame] = onUpdate
	igs.handlers[Ended] = onEnded
	return igs
}

func (i *inGameState) Input(in *lol.GameInfo) {
	//TODO add game to known games
	//TODO add known game check for started

	// updating cached game data
	if in != nil {
		i.data = in
	}

	// get transition
	t, ok := i.transitions[i.state]
	if !ok {
		panic("invalid state")
	}
	//apply transition
	nextState := t(in)
	i.state = nextState

	// get handler
	h, ok := i.handlers[nextState]
	if !ok {
		panic("unhandled state")
	}
	//use handler
	h(i.data)
}

func transitionFromIdle(in *lol.GameInfo) State {
	if in == nil {
		return Idle
	} else {
		return Started
	}
}

func transitionFromStarted(in *lol.GameInfo) State {
	if in == nil {
		return Ended
	} else {
		return InGame
	}
}

func transitionFromInGame(in *lol.GameInfo) State {
	if in == nil {
		return Ended
	} else {
		return InGame
	}
}

func transitionFromEnded(in *lol.GameInfo) State {
	if in == nil {
		return Idle
	} else {
		return Started
	}
}

func (i *inGameState) publish(topic string, info *lol.GameInfo) {
	out, err := json.Marshal(i.data)
	log.Infow("ingame update",
		"topic", topic,
		"data", string(out),
		"err", err)
}
