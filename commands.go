package main

import (
	"context"
	"effie/blocklist"
	"effie/bot"
	"effie/logger"
	"effie/output"
	"effie/riot"
	"effie/watchlist"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

var commandLogger = logger.Get("discord", "command").Sugar()

var CmdBlock = bot.Cmd(
	"block",
	"block a summoner from being tracked by effie",
	true, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options, err := bot.GetSubCommandOptions(i.Interaction.Data)
		if err != nil {
			commandLogger.Errorw(err.Error(), "data", fmt.Sprintf("%v", i.Interaction.Data))
		}
		summonerName := options[0].Value

		commandLogger.Debugw("verifying user input", "summonerName", summonerName)

		summonerInfo, err := riot.GetSummonerByName(summonerName)
		if err != nil {
			commandLogger.Error(err)
			return
		}

		watchlist.Remove(summonerInfo.ID)
		alreadyBlocked, err := blocklist.Contains(summonerInfo.ID)
		if err != nil {
			commandLogger.Error(err)
			alreadyBlocked = false
		}
		if !alreadyBlocked {
			blocklist.Append(summonerInfo.ID)
		}

		bot.RespondToCommand(s, i, "summoner was blocked from tracking")
	}, bot.Summoner)

var CmdUnBlock = bot.Cmd(
	"unblock",
	"unblock a summoner from being tracked by effie",
	true, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if !bot.IsAdmin(i.Interaction) {
			adminUser, _ := s.User(os.Getenv("DISCORD_ADMIN"))

			bot.RespondToCommand(s, i,
				"Only Admins can add blocked summoners. "+
					"Pls contact "+adminUser.Mention())
		}

		bot.RespondToCommand(s, i, "not implemented")
	}, bot.Summoner)

var CmdAdd = bot.Cmd(
	"add",
	"add a summoner to effies tracking",
	false, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		bot.RespondToCommand(s, i, "trying to add summoner to tracking...")
		commandLogger.Debugw("add command", "data", i.Interaction.Data)

		options, err := bot.GetSubCommandOptions(i.Interaction.Data)
		if err != nil {
			commandLogger.Errorw(err.Error(), "data", fmt.Sprintf("%v", i.Interaction.Data))
		}
		summonerName := options[0].Value

		commandLogger.Debugw("verifying user input", "summonerName", summonerName)
		summonerInfo, err := riot.GetSummonerByName(summonerName)
		if err != nil {
			commandLogger.Error(err)
			bot.UpdateRespondToCommand(s, i, "something went wrong")
			return
		}
		summonerId := summonerInfo.ID
		commandLogger.Debugw("retrieved summoner data",
			"summonerName", summonerName,
			"summonerId", summonerId,
			"data", summonerInfo,
		)

		blocked, err := blocklist.Contains(summonerId)
		if err != nil {
			commandLogger.Error(err)
			bot.UpdateRespondToCommand(s, i, "something went wrong")
			return
		}
		if blocked {
			commandLogger.Info("user tried to add blocked summoner")
			bot.UpdateRespondToCommand(s, i, "this user is blocked and can't tracked")
			return
		}

		if watchlist.Contains(summonerId) {
			commandLogger.Warnw("summoner is already on watchlist",
				"summonerId", summonerInfo)
			bot.UpdateRespondToCommand(s, i, "summoner is already on watchlist")
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		go StartTracking(ctx, summonerId)

		watchlist.Set(summonerId, cancel)
		bot.UpdateRespondToCommand(s, i, "added summoner to tracking")
	}, bot.Summoner)

var CmdRemove = bot.Cmd(
	"remove",
	"remove a summoner from effies tracking",
	false, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options, err := bot.GetSubCommandOptions(i.Interaction.Data)
		if err != nil {
			commandLogger.Errorw(err.Error(), "data", fmt.Sprintf("%v", i.Interaction.Data))
		}
		summonerName := options[0].Value

		commandLogger.Debugw("verifying user input", "summonerName", summonerName)

		summonerInfo, err := riot.GetSummonerByName(summonerName)
		if err != nil {
			commandLogger.Error(err)
			return
		}
		summonerId := summonerInfo.ID

		watchlist.Remove(summonerId)
		bot.RespondToCommand(s, i, "removed summoner from tracking")
	}, bot.Summoner)

var CmdOn = bot.Cmd(
	"on",
	"turns effies tracking and messages on",
	false, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		bot.RespondToCommand(s, i, "not implemented")
	}, bot.None)

var CmdOff = bot.Cmd(
	"off",
	"effie stops tracking summoners and sending messages",
	false, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		adminUser, _ := s.User(os.Getenv("DISCORD_ADMIN"))
		bot.RespondToCommand(s, i, "you just killed Effie. If you want to bring Effie back to life, contact "+adminUser.Mention())
		commandLogger.Warnw("killed by discord command",
			"issuedBy", bot.UserId(i.Interaction))
		sc <- os.Kill
	}, bot.None)

var CmdLs = bot.Cmd(
	"list",
	"lists all summoners that are currently tracked",
	false, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		bot.RespondToCommand(s, i, "loading...")

		var r string
		r = "Tracking:\n"
		for id, _ := range watchlist.Summoners {
			s, err := riot.Api.Riot.LoL.Summoner.GetByID(id)
			if err != nil {
				r = r + "- " + id + "(failed to load name)\n"
			}

			r = r + "- " + s.Name + "\n"
		}

		bot.UpdateRespondToCommand(s, i, r)
	}, bot.None)

var CmdOutputLevel = bot.Cmd(
	"output",
	"set output level",
	false, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		bot.RespondToCommand(s, i, "loading...")

		options, err := bot.GetSubCommandOptions(i.Interaction.Data)
		if err != nil {
			bot.UpdateRespondToCommand(s, i, err.Error())
		}

		level := options[0].Value
		commandLogger.Infow("trying to set output level", "level", level)

		switch level {
		case "all":
			output.SetLevel(output.All)
			break
		case "ranked":
			output.SetLevel(output.Ranked)
			break
		case "promos":
			output.SetLevel(output.Promo)
			break
		}

		bot.UpdateRespondToCommand(s, i, "Output level set to "+level)
	}, bot.OutputLevel)
