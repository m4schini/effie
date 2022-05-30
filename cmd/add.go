package cmd

import (
	"effie3/bot"
	"effie3/values"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var addCmd = &bot.Command{
	Use:         "add",
	Description: "adds summoner to watchlist",
	LoadingText: "looking speechless at volunteer...",
	Run: func(s *discordgo.Session, i *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		log.Debugw("add command", "data", i.Data)

		options, err := bot.GetCommandOptions(i.Data)
		if err != nil {
			log.Error(err)
			return
		}
		summonerName := options[0].Value.(string)

		if err := values.Targets.AddName(summonerName); err == nil {
			respond(fmt.Sprintf("watching %v with great interest", summonerName))
		} else {
			log.Error(err)
			respond(fmt.Sprintf("I don't want to have anything to do with this... *%v*", summonerName))
		}
	},
}

func init() {
	err := addCmd.AddParameter(
		discordgo.ApplicationCommandOptionString,
		"summoner",
		"account name of summoner",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddSubCommand(addCmd)
}
