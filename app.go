package main

import (
	"context"
	"effie/bot"
	"effie/logger"
	"effie/messenger"
	twinPkg "effie/twin"
	"effie/watchlist"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var log = logger.Get("main").Sugar()
var sc = make(chan os.Signal, 1)

func main() {
	defer logger.Sync()

	if os.Getenv("DISCORD_GUILD") == "" {
		panic(errors.New("DISCORD_GUILD is empty"))
	}

	disc := bot.NewBot("effie", os.Getenv("DISCORD_GUILD"),
		CmdOn,
		CmdOff,
		CmdAdd,
		CmdBlock,
		CmdRemove,
		CmdUnBlock,
		CmdLs,
		CmdOutputLevel,
	)
	defer disc.Close()

	if os.Getenv("DISCORD_CHANNEL") == "" {
		log.Warn("env DISCORD_CHANNEL is required")
		sc <- os.Kill
	} else {
		go messenger.Start(disc.Channel(os.Getenv("DISCORD_CHANNEL")))
	}

	summonerNames := os.Getenv("SUMMONER_LIST")
	if summonerNames != "" {
		names := strings.Split(summonerNames, "|")
		for _, name := range names {
			log.Warn("addSummoner is compared to string output")
			if msg := addSummoner(name); msg != "added summoner to tracking" {
				log.Errorw(msg, "summonerName", name)
			}
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Info("Received interrupt: Shutting down...")

	log.Debug("Trying to cancel summoner tracking")
	for s, cancelFunc := range watchlist.Summoners {
		log.Debugf("stopping tracking for summoner: %s", s)
		cancelFunc()
	}
}

func StartTracking(ctx context.Context, summonerId string) {
	log.Infow("started tracking summoner", "summonerId", summonerId)
	twin := twinPkg.NewSummonerTwin(summonerId)
	running := true
	go func() {
		for running {
			err := twin.Refresh()
			if err != nil {
				log.Errorw("error occurred while performing regular refresh of summoner",
					"summonerId", summonerId)
			}
			log.Debugw("Waiting for next scan interval",
				"summonerId", summonerId,
				"isRunning", running)
			time.Sleep(twinPkg.RefreshInterval)
		}
	}()

	select {
	case <-ctx.Done():
		running = false
		return
	}
}

func StopTracking(summonerId string) error {
	watchlist.Remove(summonerId)
	return nil
}
