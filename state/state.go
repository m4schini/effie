package state

import "github.com/m4schini/exstate"
import exsr "github.com/m4schini/exstate/redis"

var source exstate.Source

func init() {
	var err error
	source, err = exsr.New("", "", 0)
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
