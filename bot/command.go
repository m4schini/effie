package bot

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

func RespondToCommand(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		log.Errorw("failed to respond to command",
			"error", err,
			"channelId", i.ChannelID,
		)
	}
}

func UpdateRespondToCommand(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: content,
	})
	if err != nil {
		return
	}
	if err != nil {
		log.Errorw("failed to respond to command",
			"error", err,
			"channelId", i.ChannelID,
		)
	}
}

func IsAdmin(i *discordgo.Interaction) bool {
	if i == nil {
		return false
	}

	userId := UserId(i)
	if userId == "" {
		return false
	}

	return userId == os.Getenv("DISCORD_ADMIN")
}

func UserId(i *discordgo.Interaction) string {
	if i == nil {
		return ""
	}

	if i.User != nil {
		return i.User.ID
	} else if i.Member != nil {
		return i.Member.Mention()[3 : len(i.Member.Mention())-1] //"<@...>"
	} else {
		return ""
	}
}
