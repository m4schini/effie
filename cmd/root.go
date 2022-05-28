package cmd

import (
	"effie3/bot"
	"effie3/logger"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

var log = logger.Get("discord", "bot").Sugar()

var rootCmd = &bot.Command{
	Use:         "effie",
	Description: "effie bot description",
	Run: func(session *discordgo.Session, interaction *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		log.Warnw("Tried to call empty root command")
		respond("Tried to call empty root command")
	},
}

func Execute() (bot.Api, error) {
	d, err := bot.New(os.Getenv("DISCORD_TOKEN"), rootCmd)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func init() {
	if rootCmd.Use != strings.ToLower(rootCmd.Use) {
		log.Fatal("command use has to be lowercase")
	}
}
