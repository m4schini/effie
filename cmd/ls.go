package cmd

import (
	"effie3/bot"
	"effie3/cache"
	"effie3/watch"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var lsCmd = &bot.Command{
	Use:         "list",
	Description: "lists of actively tracked summoners",
	LoadingText: "looking for the tribute list...",
	Run: func(s *discordgo.Session, i *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		var response string
		for _, id := range watch.List() {
			response = response + fmt.Sprintf("%v\n", cache.GetSummonerName(id))
		}

		if response == "" {
			respond("big brother is sleeping")
		} else {
			respond(response)
		}
	},
}

func init() {
	rootCmd.AddSubCommand(lsCmd)
}
