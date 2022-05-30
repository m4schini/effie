package volume

import "effie3/logger"

var log = logger.Get("volume").Sugar()

type Level string

const (
	All    Level = "all"
	Ranked Level = "ranked"
	Promo  Level = "promos"
)

type Volume struct {
	level Level
}

func New(def Level) *Volume {
	return &Volume{level: def}
}

func (v *Volume) Set(level Level) {
	log.Infow("changed volume level", "before", v.level, "after", level)
	v.level = level
}

func (v *Volume) Get() Level {
	return v.level
}

func (v *Volume) Ok(level Level) bool {
	return volumeToInt(level) >= volumeToInt(v.level)
}

func volumeToInt(level Level) int {
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
