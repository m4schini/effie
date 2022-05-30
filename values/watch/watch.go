package watch

import (
	"context"
	"effie3/cache"
	"effie3/logger"
	"effie3/state"
	"errors"
	"github.com/m4schini/exstate"
)

const (
	DataKey = "watch"
)

var log = logger.Get("watchlist").Sugar()

type Watchlist struct {
	add       exstate.SetAdd
	ls        exstate.SetGet
	rm        exstate.SetRemove
	contains  exstate.SetContains
	cancelMap map[string]context.CancelFunc

	isBlocked func(id string) bool
}

func New(isBlocked func(id string) bool) *Watchlist {
	w := new(Watchlist)
	w.add, w.ls, w.rm, w.contains = state.Set(DataKey)
	w.cancelMap = make(map[string]context.CancelFunc)
	w.isBlocked = isBlocked

	return w
}

func (w *Watchlist) AddName(summonerName string) error {
	id := cache.GetSummonerId(summonerName)
	if id == "" {
		return errors.New("summoner doesn't exist")
	}
	return w.Add(id)
}

func (w *Watchlist) Add(id string) error {
	blocked := w.isBlocked(id)
	if blocked {
		log.Debugw("user tried to add blocked summoner")
		return errors.New("this user is blocked and can't be tracked")
	}

	if w.Contains(id) {
		log.Warnw("summoner is already on watchlist",
			"summonerId", id)
		return errors.New("summoner is already on watchlist")
	}

	/*ctx*/
	_, cancel := context.WithCancel(context.Background())
	//TODO go StartTracking(ctx, summonerId)

	w.cancelMap[id] = cancel
	w.add(id)
	return nil
}

func (w *Watchlist) Contains(id string) bool {
	return w.contains(id)
}

func (w *Watchlist) Remove(id string) {
	cancel, ok := w.cancelMap[id]
	if ok {
		cancel()
	}
	w.rm(id)
}

func (w *Watchlist) List() []string {
	return w.ls()
}
