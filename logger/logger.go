package logger

import (
	"go.uber.org/zap"
	"os"
)

var logger *zap.Logger

func init() {
	defer func() {
		if logger != nil {
			logger = Get("effie")
		}
	}()
	defer func() {
		if logger == nil {
			logger = zap.NewNop()
		}
	}()

	env := os.Getenv("LOGGER")
	switch env {
	case "nop":
		logger = zap.NewNop()
		break
	case "dev":
		logger, _ = zap.NewDevelopment()
		break
	case "prod":
		logger, _ = zap.NewProduction()
		break
	default:
		logger = zap.NewNop()
		break
	}
}

func Get(names ...string) *zap.Logger {
	l := logger

	for _, name := range names {
		l = l.Named(name)
	}

	return l
}

func Sync() {
	if logger != nil {
		logger.Sync()
	}
}
