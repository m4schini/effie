package cmd

import (
	"effie3/bot"
	"github.com/bwmarrin/discordgo"
	"os"
)

var offCmd = &bot.Command{
	Use:         "off",
	Description: "turns effie irreversibly off",
	LoadingText: "killing effie...",
	Run: func(s *discordgo.Session, i *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		adminUser, _ := s.User(os.Getenv("DISCORD_ADMIN"))
		respond("you just killed Effie. If you want to bring Effie back to life, contact " + adminUser.Mention())
		log.Warnw("killed by discord command", "issuedBy", bot.UserId(i))
		log.Fatal("bye")
	},
}

func init() {
	rootCmd.AddSubCommand(offCmd)
}
