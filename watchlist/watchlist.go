package watchlist

import (
	"context"
	"effie/logger"
	"fmt"
)

var log = logger.Get("watchlist").Sugar()

var Summoners map[string]context.CancelFunc

func init() {
	Summoners = make(map[string]context.CancelFunc)
}

func Remove(summonerId string) {
	log.Infow("removing summoner from watchlist", "summonerId", summonerId)
	cancel, exists := Summoners[summonerId]
	if !exists {
		log.Warnf("summoner '%s' doesn't exist", summonerId)
		return
	}

	cancel()
	delete(Summoners, summonerId)
}

func Contains(summonerId string) bool {
	_, exists := Summoners[summonerId]
	return exists
}

func Get(summonerId string) (context.CancelFunc, error) {
	cancelFunc, exists := Summoners[summonerId]
	if exists {
		return nil, fmt.Errorf("summoner '%s' is already being watched", summonerId)
	}

	return cancelFunc, nil
}

func Set(summonerId string, cancelFunc context.CancelFunc) {
	log.Infow("adding/updating summoner on watchlist", "summonerId", summonerId)
	Summoners[summonerId] = cancelFunc
}
