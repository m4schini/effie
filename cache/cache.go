package cache

import (
	"effie3/logger"
	"effie3/riot"
	"github.com/m4schini/exstate"
	er "github.com/m4schini/exstate/redis"
	"os"
	"time"
)

const (
	cacheExpiration = 24 * time.Hour
)

var cache exstate.Cache
var source exstate.Source

var log = logger.Get("cache", "exstate").Sugar()

func init() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		log.Warn("REDIS_ADDR is missing")
	}
	pass := os.Getenv("REDIS_PASS")
	if pass == "" {
		log.Warn("REDIS_PASS is missing")
	}

	r, err := er.New(addr, pass, 0)
	if err != nil {
		log.Error(err)
	}
	cache = r
	source = r
}

func GetHungerGamesCount() (exstate.GetInt, exstate.Setter[int]) {
	if source == nil {
		return func() int {
			return 0
		}, func(v int) {}
	}

	return source.Int("count.ranked")
}

func GetPromoGamesCount() (exstate.GetInt, exstate.Setter[int]) {
	if source == nil {
		return func() int {
			return 0
		}, func(v int) {}
	}

	return source.Int("count.promo")
}

func GetSummonerId(name string) string {
	if cache == nil {
		return ""
	}

	getId, _, err := cache.New("cache.id."+name, cacheExpiration, makeIdRetriever(name))
	if err != nil {
		log.Error(err)
		return ""
	}

	res, err := getId()
	if err != nil {
		log.Error(err)
		return ""
	}

	if res == nil {
		log.Warnw("no cached data available", "summonerName", name)
		return ""
	}

	return res.(string)
}

func GetSummonerName(id string) string {
	if cache == nil {
		return ""
	}

	getId, _, err := cache.New("cache.id."+id, cacheExpiration, makeNameRetriever(id))
	if err != nil {
		log.Error(err)
		return ""
	}

	res, err := getId()
	if err != nil {
		log.Error(err)
		return ""
	}

	if res == nil {
		log.Warnw("no cached data available", "summonerId", id)
		return ""
	}

	return res.(string)
}

func makeNameRetriever(id string) func() interface{} {
	return func() interface{} {
		log.Debugw("cached summoner name expired, retrieving update", "summonerId", id)
		info, err := riot.GetSummonerById(id)
		if err != nil {
			log.Errorw("cache update failed", "err", err, "id", id)
			return nil
		}

		return info.Name
	}
}

func makeIdRetriever(name string) func() interface{} {
	return func() interface{} {
		log.Debugw("cached summoner id expired, retrieving update", "summonerName", name)
		info, err := riot.GetSummonerByName(name)
		if err != nil {
			log.Errorw("cache update failed", "err", err, "name", name)
			return nil
		}

		return info.ID
	}
}
