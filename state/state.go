package state

import (
	"effie3/logger"
	"github.com/m4schini/exstate"
	"os"
)
import exsr "github.com/m4schini/exstate/redis"

var source exstate.Source

var log = logger.Get("state", "deprecated")

func init() {
	var err error
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		log.Warn("REDIS_ADDR is missing")
	}
	pass := os.Getenv("REDIS_PASS")
	if pass == "" {
		log.Warn("REDIS_PASS is missing")
	}

	source, err = exsr.New(addr, pass, 0)
	if err != nil {
		panic(err)
	}
}

func String(path ...string) (exstate.GetString, exstate.Setter[string]) {
	return source.String(path...)
}
func Int(path ...string) (exstate.GetInt, exstate.Setter[int]) {
	return source.Int(path...)
}
func Float(path ...string) (exstate.GetFloat, exstate.Setter[float64]) {
	return source.Float(path...)
}
func Bool(path ...string) (exstate.GetBool, exstate.Setter[bool]) {
	return source.Bool(path...)
}
func Set(path ...string) (exstate.SetAdd, exstate.SetGet, exstate.SetRemove, exstate.SetContains) {
	return source.Set(path...)
}
