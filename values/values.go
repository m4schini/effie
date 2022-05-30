package values

import (
	"effie3/values/block"
	"effie3/values/volume"
	"effie3/values/watch"
)

var Volume *volume.Volume
var Targets *watch.Watchlist
var Forbidden *block.Blocklist

func init() {
	Volume = volume.New(volume.Ranked)
	Forbidden = block.New()
	Targets = watch.New(func(id string) bool {
		return Forbidden.IsBlocked(id)
	})
}
