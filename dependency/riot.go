package dependency

import (
	"effie3/riot"
	"os"
	"sync"
)

// ====================
// RIOT API
// ====================

const (
	envNameRiotApiKey    = "RIOT_API_KEY"
	envNameRiotApiRegion = "RIOT_API_REGION"
)

var riotApi riot.Api
var riotApiInitialised = false
var riotInitLock sync.Mutex

func readyRiot() {
	riotInitLock.Lock()
	defer riotInitLock.Unlock()
	if !riotApiInitialised {
		initRiotApi()
	}
}

func initRiotApi() {
	if os.Getenv(envNameRiotApiRegion) == "" {
		log.Warn(envNameRiotApiRegion + " is missing")
		return
	}
	if os.Getenv(envNameRiotApiKey) == "" {
		log.Warn(envNameRiotApiKey + " is missing")
		return
	}

	apiSingleton, err := riot.NewApi(os.Getenv(envNameRiotApiRegion), os.Getenv(envNameRiotApiKey))
	if err != nil {
		log.Warnw("riot api connection could not be established")
	}
	log.Infow("initialized new riot-api connection")

	riotApi = apiSingleton
	riotApiInitialised = true
}

func RiotApi() riot.Api {
	readyRiot()

	return riotApi
}
