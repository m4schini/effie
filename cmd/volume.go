package cmd

import (
	"effie3/bot"
	"effie3/conf"
	"github.com/bwmarrin/discordgo"
)

var volumeCmd = &bot.Command{
	Use:         "volume",
	Description: "set volume level",
	LoadingText: "looking for the tv remote...",
	Run: func(s *discordgo.Session, i *discordgo.Interaction, respond bot.InteractionRespondFunc) {
		options, err := bot.GetCommandOptions(i.Data)
		if err != nil {
			log.Error(err)
			return
		}
		level := options[0].Value.(string)
		log.Infow("trying to set output level", "level", level)

		switch level {
		case "all":
			conf.SetVolume(conf.All)
			break
		case "ranked":
			conf.SetVolume(conf.Ranked)
			break
		case "promos":
			conf.SetVolume(conf.Promo)
			break
		}

		respond("Volume set to " + level)
	},
}

func init() {
	err := volumeCmd.AddParameter(
		discordgo.ApplicationCommandOptionString,
		"level",
		"volume level (what kind of games should effie pay attention to)",
		true,
		&discordgo.ApplicationCommandOptionChoice{
			Name:  "all",
			Value: conf.All,
		},
		&discordgo.ApplicationCommandOptionChoice{
			Name:  "ranked + promos",
			Value: conf.Ranked,
		},
		&discordgo.ApplicationCommandOptionChoice{
			Name:  "promos only",
			Value: conf.Promo,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddSubCommand(volumeCmd)
}
