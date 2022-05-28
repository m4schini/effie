package bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

type Message interface {
	Edit(content string) error
	EditCh() (chan<- string, error)
	Delete() error
}

type message struct {
	ref *discordgo.MessageReference
	api Api
}

func NewMessage(api Api, ref *discordgo.MessageReference) *message {
	msg := new(message)
	msg.api = api
	msg.ref = ref
	return msg
}

func (m *message) Edit(content string) error {
	if m == nil || m.api == nil {
		return errors.New("something is missing")
	}
	_, err := m.api.EditMessage(m.ref.ChannelID, m.ref.MessageID, content)
	return err
}

func (m *message) EditCh() (chan<- string, error) {
	ch := make(chan string)

	go func() {
		for s := range ch {
			err := m.Edit(s)
			if err != nil {
				log.Error(err)
			}
		}
	}()

	return nil, errors.New("not implemented")
}

func (m *message) Delete() error {
	err := m.api.DeleteMessage(m.ref.ChannelID, m.ref.MessageID)
	return err
}
