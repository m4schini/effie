package riot

import (
	"effie3/dependency"
	"effie3/values/volume"
	"errors"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/KnutZuidema/golio/static"
)

var apiSingleton = dependency.RiotApi()

func apiHealthy() bool {
	return apiSingleton != nil
}

func GetSummonerByName(summonerName string) (*lol.Summoner, error) {
	if !apiHealthy() {
		return nil, errors.New("api client not available")
	}

	return apiSingleton.GetSummonerByName(summonerName)
}

func GetSummonerById(summonerId string) (*lol.Summoner, error) {
	if !apiHealthy() {
		return nil, errors.New("api client not available")
	}

	return apiSingleton.GetSummonerById(summonerId)
}

func GetCurrentGame(summonerId string) (*lol.GameInfo, error) {
	if !apiHealthy() {
		return nil, errors.New("api client not available")
	}

	return apiSingleton.GetCurrentGame(summonerId)
}

func GetLeagues(summonerId string) ([]*lol.LeagueItem, error) {
	if !apiHealthy() {
		return nil, errors.New("api client not available")
	}

	return apiSingleton.GetLeagues(summonerId)
}

func GetMatch(matchId string, summonerId string) (*lol.Match, *lol.Participant, error) {
	if !apiHealthy() {
		return nil, nil, errors.New("api client not available")
	}

	return apiSingleton.GetMatch(matchId, summonerId)
}

func GetQueueInfo(queueId int) (static.Queue, error) {
	if !apiHealthy() {
		return static.Queue{}, errors.New("api client not available")
	}

	return apiSingleton.Static().GetQueue(queueId)
}

func GetGameLevel(summonerId string, info *lol.GameInfo) (volume.Level, error) {
	if !apiHealthy() {
		return volume.All, errors.New("api client not available")
	}

	return apiSingleton.GetGameLevel(summonerId, info)
}
