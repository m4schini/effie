package block

import (
	"effie3/state"
	"errors"
)

const (
	DataKey = "blocks"
)

var addBlock, getBlocked, removeBlock, isBlocked = state.Set(DataKey)

func AddName(summonerName string) error {
	return errors.New("not implemented")
}

func Add(id string) {
	addBlock(id)
}

func IsBlocked(id string) bool {
	return isBlocked(id)
}

func Remove(id string) {
	removeBlock(id)
}
