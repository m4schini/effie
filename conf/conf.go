package conf

import "effie3/logger"

var log = logger.Get("config").Sugar()

type VolumeLevel string

const (
	All    VolumeLevel = "all"
	Ranked VolumeLevel = "ranked"
	Promo  VolumeLevel = "promos"
)

var volumeLevel = Ranked

func SetVolume(level VolumeLevel) {
	log.Infow("changed volume level", "before", volumeLevel, "after", level)
	volumeLevel = level
}

func Volume() VolumeLevel {
	return volumeLevel
}

func VolumeOk(level VolumeLevel) bool {
	return volumeToInt(level) >= volumeToInt(volumeLevel)
}

func volumeToInt(level VolumeLevel) int {
	switch level {
	case Promo:
		return 3
	case Ranked:
		return 2
	case All:
		return 1
	default:
		return 0
	}
}
