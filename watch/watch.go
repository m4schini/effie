package watch

import (
	"context"
	"effie3/block"
	"effie3/cache"
	"effie3/logger"
	"effie3/state"
	"errors"
)

const (
	DataKey = "watch"
)

var log = logger.Get("watchlist").Sugar()
var add, ls, rm, contains = state.Set(DataKey)
var cancelMap = make(map[string]context.CancelFunc)

func AddName(summonerName string) error {
	id := cache.GetSummonerId(summonerName)
	if id == "" {
		return errors.New("summoner doesn't exist")
	}
	return Add(id)
}

func Add(id string) error {
	blocked := block.IsBlocked(id)
	if blocked {
		log.Debugw("user tried to add blocked summoner")
		return errors.New("this user is blocked and can't be tracked")
	}

	if Contains(id) {
		log.Warnw("summoner is already on watchlist",
			"summonerId", id)
		return errors.New("summoner is already on watchlist")
	}

	/*ctx*/
	_, cancel := context.WithCancel(context.Background())
	//TODO go StartTracking(ctx, summonerId)

	cancelMap[id] = cancel
	add(id)
	return nil
}

func Contains(id string) bool {
	return contains(id)
}

func Remove(id string) {
	cancel, ok := cancelMap[id]
	if ok {
		cancel()
	}
	rm(id)
}

func List() []string {
	return ls()
}
