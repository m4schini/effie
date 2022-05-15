package output

import "sync"

var outputlevel Level
var olMu sync.Mutex

type Level string

func init() {
	SetLevel(All)
}

const (
	All    Level = "all"
	Ranked Level = "ranked"
	Promo  Level = "promos"
)

func SetLevel(level Level) {
	olMu.Lock()
	outputlevel = level
	olMu.Unlock()
}

func GetLevel() Level {
	olMu.Lock()
	defer olMu.Unlock()
	return outputlevel
}

func Ok(level Level) bool {
	return toInt(level) >= toInt(outputlevel)
}

func toInt(level Level) int {
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
