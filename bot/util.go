package bot

import (
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
)

const (
	defaultLoadingText = "something is happening"
)

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

type Option struct {
	Name  string
	Value interface{}
}

func GetCommandOptions(data discordgo.InteractionData) ([]Option, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	interationOptions, ok := result["options"].([]interface{})
	if !ok {
		return nil, errors.New("something happened")
	}

	subCommand, ok := interationOptions[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("something happened")
	}

	scOptions := subCommand["options"].([]interface{})
	if !ok {
		return nil, errors.New("something happened")
	}

	options := make([]Option, 0, len(scOptions))
	for _, option := range scOptions {
		optionMap := option.(map[string]interface{})

		options = append(options, Option{optionMap["name"].(string), optionMap["value"]})
	}

	return options, nil
}

func addCatchAllHandler(session *discordgo.Session, handlers map[string]Command) {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Interaction == nil ||
			i.Interaction.ApplicationCommandData().Options == nil ||
			len(i.Interaction.ApplicationCommandData().Options) < 1 {
			log.Errorw("missing data",
				"command", i.ApplicationCommandData().Name)
		}

		var cmdName string
		if len(i.Interaction.ApplicationCommandData().Options) > 0 {
			cmdName = i.Interaction.ApplicationCommandData().Options[0].Name
		} else {
			cmdName = i.Interaction.ApplicationCommandData().Name
		}

		log.Infow("incoming discord interaction", "userID", UserId(i.Interaction), "command", cmdName)
		if h, ok := handlers[cmdName]; ok {
			if h.LoadingText != "" {
				RespondToCommand(s, i.Interaction, h.LoadingText)
			} else {
				RespondToCommand(s, i.Interaction, defaultLoadingText)
			}
			h.Run(s, i.Interaction, func(content string) {
				UpdateRespondToCommand(s, i.Interaction, content)
			})
		} else {
			if i.User != nil {
				log.Warnw("command invoked by direct messages",
					"invokedBy", i.User.ID,
					"invokedIn", i.ChannelID,
					"commandName", cmdName,
				)
			}
			log.Warnw("tried to handle command without handler",
				"invokedBy", i.Member.Mention(),
				"invokedIn", i.ChannelID,
				"commandName", cmdName,
			)
		}
	})
}

func RespondToCommand(s *discordgo.Session, i *discordgo.Interaction, content string) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
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

func UpdateRespondToCommand(s *discordgo.Session, i *discordgo.Interaction, content string) {
	_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
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
