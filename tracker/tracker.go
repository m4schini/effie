package tracker

import (
	"context"
	"effie3/logger"
	"effie3/twin"
	"time"
)

const (
	RefreshInterval = 20 * time.Second
)

var log = logger.Get("tracker").Sugar()

var twins = make(map[string]twin.SummonerTwin)

func Start(ctx context.Context, getTargets func() []string) error {
	for {
		done := make(chan struct{}, 0)

		go func() {
			summonerIds := getTargets()

			log.Infow("starting scan run", "count", len(summonerIds))
			for _, id := range summonerIds {
				t, ok := twins[id]
				if !ok {
					t = twin.NewSummonerTwin(id)
					twins[id] = t
				} else {
					log.Debugw("updating existing summoner twin")
				}

				err := t.Refresh()
				if err != nil {
					log.Errorw("error while updating summoner twin", "err", err)
				}
			}

			done <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			log.Infof("next summoner update scheduled for %v",
				time.Now().Add(RefreshInterval).Format(time.RFC3339),
			)
			time.Sleep(RefreshInterval)
		}
	}
}
