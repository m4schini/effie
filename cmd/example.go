package cmd

import (
	"effie3/bot"
	"github.com/bwmarrin/discordgo"
)

var exampleCmd = &bot.Command{
	Use:         "exampleName",
	Description: "example description",
	LoadingText: "example loading text",
	Run: func(s *discordgo.Session, i *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		respond("this command should actually only be used as an example")
	},
}

func init() {
	//TODO change cmd
	err := exampleCmd.AddParameter(
		discordgo.ApplicationCommandOptionString,
		"summoner",
		"Name of summoner/account",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	//TODO comment this in if you want to actually add the command
	//rootCmd.AddSubCommand(blockCmd)
}
