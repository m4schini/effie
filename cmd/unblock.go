package cmd

import (
	"effie3/bot"
	"github.com/bwmarrin/discordgo"
)

var unblockCmd = &bot.Command{
	Use:         "unblock",
	Description: "this removes a summoner from the 'block tracking list'",
	LoadingText: "looking for the keys...",
	Run: func(s *discordgo.Session, i *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		respond("not implemented")
	},
}

func init() {
	err := unblockCmd.AddParameter(
		discordgo.ApplicationCommandOptionString,
		"summoner",
		"account name of summoner",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddSubCommand(unblockCmd)
}
