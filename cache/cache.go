package cache

import (
	"effie3/logger"
	"effie3/riot"
	"github.com/m4schini/exstate"
	er "github.com/m4schini/exstate/redis"
	"time"
)

const (
	cacheExpiration = 24 * time.Hour
)

var src exstate.Cache

var log = logger.Get("cache", "exstate").Sugar()

func init() {
	r, _ := er.New("", "", 0)
	src = r
}

func GetSummonerId(name string) string {
	if src == nil {
		return ""
	}

	getId, _, err := src.New("cache.id."+name, cacheExpiration, makeIdRetriever(name))
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
	if src == nil {
		return ""
	}

	getId, _, err := src.New("cache.id."+id, cacheExpiration, makeNameRetriever(id))
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
