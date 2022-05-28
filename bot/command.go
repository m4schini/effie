package bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandFunc func(s *discordgo.Session, i *discordgo.Interaction, respond InteractionRespondFunc)

type InteractionRespondFunc func(content string)

type Command struct {
	Use         string
	Description string
	LoadingText string
	Run         CommandFunc

	_parameter   []*Parameter
	_type        discordgo.ApplicationCommandOptionType
	_subCommands []*Command
}

//Todo deprecated?
func NewCommand(use, description string) (*Command, error) {
	if use == "" {
		return nil, errors.New("name can't be empty")
	}
	if description == "" {
		return nil, errors.New("description can't be empty")
	}

	return &Command{
		Use:         use,
		Description: description,
		Run: func(session *discordgo.Session, interaction *discordgo.Interaction, respondFunc InteractionRespondFunc) {
			respondFunc("nothing to see here")
		},
		_parameter:   nil,
		_type:        discordgo.ApplicationCommandOptionSubCommand,
		_subCommands: nil,
	}, nil
}

func (c *Command) ApplicationCommand() *discordgo.ApplicationCommand {
	c.Use = strings.ToLower(c.Use)
	dcmd := &discordgo.ApplicationCommand{
		Name:        c.Use,
		Description: c.Description,
		Options:     make([]*discordgo.ApplicationCommandOption, 0),
	}

	for _, command := range c._subCommands {
		command.Use = strings.ToLower(command.Use)

		subCmd := &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        command.Use,
			Description: command.Description,
			Options:     make([]*discordgo.ApplicationCommandOption, 0),
		}

		for _, parameter := range command._parameter {
			subCmd.Options = append(subCmd.Options, &discordgo.ApplicationCommandOption{
				Name:        parameter.Name,
				Description: parameter.Description,
				Type:        parameter.Type,
				Required:    parameter.Required,
				Choices:     parameter.Choices,
			})
		}

		// remove empty array
		if len(subCmd.Options) == 0 {
			subCmd.Options = nil
		}
		dcmd.Options = append(dcmd.Options, subCmd)
	}

	// remove empty array
	if len(dcmd.Options) == 0 {
		dcmd.Options = nil
	}

	return dcmd
}

func (c *Command) HasSubCommands() bool {
	return c._subCommands != nil && len(c._subCommands) > 0
}

func (c *Command) AddSubCommand(subCommand *Command) {
	if c == nil {
		return
	}
	log.Infow("adding subcommand", "root", c.Use, "use", subCommand.Use)

	if c._subCommands == nil {
		c._subCommands = make([]*Command, 0)
	}

	c._subCommands = append(c._subCommands, subCommand)
}

func (c *Command) AddParameter(
	t discordgo.ApplicationCommandOptionType,
	name, description string,
	required bool,
	choices ...*discordgo.ApplicationCommandOptionChoice,
) error {
	if t == discordgo.ApplicationCommandOptionSubCommand || t == discordgo.ApplicationCommandOptionSubCommandGroup {
		return errors.New("type can't be SubCommand or SubCommandGroup")
	}
	if name == "" {
		return errors.New("name can't be empty")
	}
	if description == "" {
		return errors.New("description can't be empty")
	}

	if name != strings.ToLower(name) {
		return errors.New("command use has to be lowercase")
	}

	if c._parameter == nil {
		c._parameter = make([]*Parameter, 0, 2)
	}

	c._parameter = append(c._parameter, &Parameter{
		Type:        t,
		Name:        name,
		Description: description,
		Required:    required,
		Choices:     choices,
	})
	return nil
}

type Parameter struct {
	Type        discordgo.ApplicationCommandOptionType
	Name        string
	Description string
	Required    bool
	Choices     []*discordgo.ApplicationCommandOptionChoice
}
