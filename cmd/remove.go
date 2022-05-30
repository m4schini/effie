package cmd

import (
	"effie3/bot"
	"effie3/riot"
	"effie3/values"
	"github.com/bwmarrin/discordgo"
)

var removeCmd = &bot.Command{
	Use:         "remove",
	Description: "removes summoner from watchlist",
	LoadingText: "Losing interest...",
	Run: func(s *discordgo.Session, i *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		options, err := bot.GetCommandOptions(i.Data)
		if err != nil {
			log.Error(err)
			return
		}
		summonerName := options[0].Value.(string)

		log.Debugw("verifying user input", "summonerName", summonerName)

		summonerInfo, err := riot.GetSummonerByName(summonerName)
		if err != nil {
			log.Error(err)
			return
		}
		summonerId := summonerInfo.ID

		values.Targets.Remove(summonerId)
		respond("list interest in " + summonerName)
	},
}

func init() {
	err := removeCmd.AddParameter(
		discordgo.ApplicationCommandOptionString,
		"summoner",
		"account name of summoner",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddSubCommand(removeCmd)
}
