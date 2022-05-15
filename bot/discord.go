package bot

import (
	"effie/logger"
	"errors"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

var log = logger.Get("discordV2").Sugar()

type Builder interface {
	OnCommand(cmd string, f func()) error
	Build() Bot
}

type Bot interface {
	Channel(channelId string) Messenger
	Close()
}

type Messenger interface {
	SendMessage(content string) (*discordgo.Message, error)
	UpdateMessage(messageId, content string) (*discordgo.Message, error)
	DeleteMessage(messageId string) error
}

type bot struct {
	session  *discordgo.Session
	commands []*discordgo.ApplicationCommand
	handlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	channel  string
	guild    string
}

func (b *bot) Channel(channelId string) Messenger {
	return &channel{
		id:  channelId,
		bot: b,
	}
}

func (b *bot) ClearAllCommands(guildId string) error {
	if b.session == nil {
		return errors.New("missing discordgo session")
	}

	commands, err := b.session.ApplicationCommands(b.session.State.User.ID, guildId)
	if err != nil {
		return err
	}

	for _, command := range commands {
		b.session.ApplicationCommandDelete(b.session.State.User.ID, guildId, command.ID)
	}
	return nil
}

func (b *bot) Close() {
	if b == nil {
		return
	}
	if b.session == nil {
		return
	}
	defer func(session *discordgo.Session) {
		log.Warnw("closing discord session", "sessionId", b.session.State.SessionID)
		err := session.Close()
		if err != nil {
			log.Error("error while trying to close discord session", err)
		}
	}(b.session)

	if b.commands != nil {
		for _, v := range b.commands {
			err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.guild, v.ID)
			if err != nil {
				log.Errorf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	err := b.ClearAllCommands(b.guild)
	if err != nil {
		log.Error(err)
	}
	err = b.ClearAllCommands("")
	if err != nil {
		log.Error(err)
	}
}

type OptionType uint8

const (
	None        OptionType = iota
	Summoner               = 3
	OutputLevel            = 4
)

type commandOption struct {
	Type     OptionType
	Name     string
	Desc     string
	Required bool
}

type botCommand struct {
	Name        string
	Description string
	Global      bool
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Option      OptionType
}

func Cmd(
	name, description string,
	global bool,
	f func(s *discordgo.Session, i *discordgo.InteractionCreate),
	options OptionType,
) *botCommand {
	cmd := new(botCommand)
	cmd.Name = name
	cmd.Description = description
	cmd.Global = global
	cmd.Handler = f
	cmd.Option = options
	return cmd
}

func NewBot(botName, guildId string, commands ...*botCommand) *bot {
	var err error

	b := new(bot)
	b.guild = guildId

	log.Info("creating discord bot using token: ", os.Getenv("DISCORD_TOKEN"))
	b.session, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Debug("discord session created")
	}
	err = b.session.Open()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Debug("discord session opened")
	}

	// handle commands
	b.handlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Interaction == nil ||
			i.Interaction.ApplicationCommandData().Options == nil ||
			len(i.Interaction.ApplicationCommandData().Options) < 1 {
			log.Errorw("missing data",
				"guildId", guildId,
				"command", i.ApplicationCommandData().Name)
		}

		subCmdName := i.Interaction.ApplicationCommandData().Options[0].Name

		log.Infow("incoming discord interaction", "userID", UserId(i.Interaction), "command", subCmdName)
		if h, ok := b.handlers[subCmdName]; ok {
			h(s, i)
		} else {
			if i.User != nil {
				log.Warnw("command invoked by direct messages",
					"invokedBy", i.User.ID,
					"invokedIn", i.ChannelID,
					"guildId", guildId,
					"commandName", i.ApplicationCommandData().Name,
				)
			}
			log.Warnw("tried to handle command without handler",
				"invokedBy", i.Member.Mention(),
				"invokedIn", i.ChannelID,
				"guildId", guildId,
				"commandName", i.ApplicationCommandData().Name,
			)
		}
	})

	subCommands := make([]*discordgo.ApplicationCommandOption, 0)
	subCommandsGlobal := make([]*discordgo.ApplicationCommandOption, 0)
	for _, command := range commands {

		//b.commands = append(b.commands, cmd)
		cmd := &discordgo.ApplicationCommandOption{
			Name:        command.Name,
			Description: command.Description,
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		}

		if command.Option == Summoner {
			//TODO this is hardcoded shit
			cmd.Options = []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionType(Summoner),
				Name:        "summoner",
				Description: "Name of the summoner you want to modify",
				Required:    true,
			}}
		}
		if command.Option == OutputLevel {
			//TODO this is hardcoded shit
			cmd.Options = []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionType(Summoner),
				Name:        "level",
				Description: "Output level",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "all",
						Value: "all",
					},
					{
						Name:  "ranked + promos",
						Value: "ranked",
					},
					{
						Name:  "promos only",
						Value: "promos",
					},
				},
			}}
		}

		if command.Global {
			subCommandsGlobal = append(subCommandsGlobal, cmd)
		} else {
			subCommands = append(subCommands, cmd)
		}
		log.Infow("Registering handler", "commandName", command.Name)
		f := command.Handler
		n := command.Name
		b.handlers[command.Name] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Infow("valid command issued", "issuedBy", UserId(i.Interaction), "cmd", n)
			f(s, i)
		}
	}
	handlerNames := make([]string, 0, len(b.handlers))
	for s, _ := range b.handlers {
		handlerNames = append(handlerNames, s)
	}
	log.Infow("registered handlers", "available", "["+strings.Join(handlerNames, ",")+"]")

	// register guild commands
	b.commands = make([]*discordgo.ApplicationCommand, 0, len(commands))
	_, err = b.session.ApplicationCommandCreate(b.session.State.User.ID, guildId, &discordgo.ApplicationCommand{
		Name:        strings.ToLower(botName),
		Description: botName + " is here to annoy everyone with promo game announcements",
		Options:     subCommands,
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// register commands
	b.commands = make([]*discordgo.ApplicationCommand, 0, len(commands))
	_, err = b.session.ApplicationCommandCreate(b.session.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        strings.ToLower(botName),
		Description: botName + " is here to annoy everyone with promo game announcements",
		Options:     subCommandsGlobal,
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return b
}

type channel struct {
	id  string
	bot *bot
}

func (c *channel) SendMessage(content string) (*discordgo.Message, error) {
	return c.bot.session.ChannelMessageSend(c.id, content)
}

func (c *channel) UpdateMessage(messageId, content string) (*discordgo.Message, error) {
	return c.bot.session.ChannelMessageEdit(c.id, messageId, content)
}

func (c *channel) DeleteMessage(messageId string) error {
	return c.bot.session.ChannelMessageDelete(c.id, messageId)
}
