package bot

import (
	"effie3/logger"
	"github.com/bwmarrin/discordgo"
	"os"
)

var log = logger.Get("discord", "api").Sugar()

type Api interface {
	SendMessage(channelId, content string) (Message, error)
	EditMessage(channelId, messageId, content string) (Message, error)
	DeleteMessage(channelId, messageId string) error
}

type discordApi struct {
	session *discordgo.Session

	commandHandlers map[string]Command
}

func New(token string, commands ...*Command) (*discordApi, error) {
	var err error

	if token == "" {
		log.Warn("discord token is empty")
	}

	api := new(discordApi)
	api.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	log.Infow("discord session created",
		"lastHeartbeatSent", api.session.LastHeartbeatSent,
		"lastHeartbeatAck", api.session.LastHeartbeatAck,
	)
	err = api.session.Open()
	if err != nil {
		return nil, err
	}

	if len(commands) > 0 {
		api.commandHandlers = make(map[string]Command)
		addCatchAllHandler(api.session, api.commandHandlers)
	}

	for _, command := range commands {
		_, err := api.session.ApplicationCommandCreate(
			api.session.State.User.ID,
			os.Getenv("DISCORD_GUILD"),
			command.ApplicationCommand(),
		)
		if err != nil {
			log.Errorw("error while registering command",
				"command", command.ApplicationCommand().Name,
				"err", err,
			)
			continue
		}
		log.Infow("registered command with discord", "command", command.Use)

		if !command.HasSubCommands() {
			api.commandHandlers[command.Use] = *command
		} else {
			for _, subCommand := range command._subCommands {
				api.commandHandlers[subCommand.Use] = *subCommand
				log.Infow("added subcommand handler", "command", subCommand.Use)
			}
		}
	}

	return api, nil
}

func (d *discordApi) SendMessage(channelId, content string) (Message, error) {
	dgoMsg, err := d.session.ChannelMessageSend(channelId, content)
	if err != nil {
		return nil, err
	}
	log.Debugw("send discord message", "channelId", channelId, "contentSize", len(content))

	return NewMessage(d, dgoMsg.Reference()), nil
}

func (d *discordApi) EditMessage(channelId, messageId, content string) (Message, error) {
	dgoMsg, err := d.session.ChannelMessageEdit(channelId, messageId, content)
	if err != nil {
		return nil, err
	}
	log.Debugw("edited discord message",
		"channelId", channelId,
		"messageId", messageId,
		"contentSize", len(content))

	return NewMessage(d, dgoMsg.Reference()), nil
}

func (d *discordApi) DeleteMessage(channelId, messageId string) error {
	err := d.session.ChannelMessageDelete(channelId, messageId)
	if err == nil {
		log.Debugw("deleted discord message",
			"channelId", channelId,
			"messageId", messageId)
	}
	return err
}
