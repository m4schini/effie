package block

import (
	"effie3/cache"
	"effie3/state"
	"github.com/m4schini/exstate"
)

const (
	DataKey = "blocks"
)

type Blocklist struct {
	add      exstate.SetAdd
	get      exstate.SetGet
	rm       exstate.SetRemove
	contains exstate.SetContains
}

func New() *Blocklist {
	b := new(Blocklist)
	b.add, b.get, b.rm, b.contains = state.Set(DataKey)

	return b
}

func (b *Blocklist) AddName(summonerName string) {
	id := cache.GetSummonerId(summonerName)
	b.Add(id)
}

func (b *Blocklist) Add(id string) {
	b.add(id)
}

func (b *Blocklist) IsBlocked(id string) bool {
	return b.contains(id)
}

func (b *Blocklist) Remove(id string) {
	b.rm(id)
}
