package cmd

import (
	"effie3/bot"
	"effie3/riot"
	"effie3/values"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var blockCmd = &bot.Command{
	Use:         "block",
	Description: "blocks a summoner from being tracked",
	LoadingText: "loading...",
	Run: func(s *discordgo.Session, i *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		options, err := bot.GetCommandOptions(i.Data)
		if err != nil {
			log.Error(err)
			return
		}
		summonerName := options[0].Value.(string)
		summonerInfo, err := riot.GetSummonerByName(summonerName)
		if err != nil {
			log.Error(err)
			return
		}

		values.Forbidden.AddName(summonerName)
		respond(fmt.Sprintf("blocked %v", summonerInfo.Name))
	},
}

func init() {
	err := blockCmd.AddParameter(
		discordgo.ApplicationCommandOptionString,
		"summoner",
		"account name of summoner",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddSubCommand(blockCmd)
}
