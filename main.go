/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"context"
	"effie3/bot"
	"effie3/cmd"
	"effie3/logger"
	"effie3/messenger"
	"effie3/riot"
	"effie3/tracker"
	"effie3/watch"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var sc = make(chan os.Signal, 1)

var botChannel = make(chan bot.Api)

var log = logger.Get("app").Sugar()

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		api, err := cmd.Execute()
		if err != nil {
			log.Error(err)
			sc <- os.Kill
		}

		botChannel <- api
		wg.Done()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if os.Getenv("DISCORD_CHANNEL") == "" {
			log.Warnw("DISCORD_CHANNEL is empty")
			sc <- os.Kill
			return
		}
		messenger.Start(<-botChannel, os.Getenv("DISCORD_CHANNEL"))
		wg.Done()
	}()

	go func() {
		wg.Wait()
		err := tracker.Start(ctx, watch.List)
		if err != nil {
			log.Error(err)
			sc <- os.Kill
		}
	}()

	for _, s := range strings.Split(os.Getenv("SUMMONER_LIST"), ";") {
		if s == "" {
			log.Fatal("empty name not allowed")
		}
		log.Debugw("adding summoner from env list", "name", s)
		info, err := riot.GetSummonerByName(s)
		if err != nil {
			log.Errorw("error while adding predefined summoner", "err", err, "summonerName", s)
			continue
		}
		err = watch.Add(info.ID)
		if err != nil {
			log.Warn(err)
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	//fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Warnw("shutting down...")
}
