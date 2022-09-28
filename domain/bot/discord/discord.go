package discord

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kh9543/koala/domain/bot"
)

type DiscordBot struct {
	token   string
	id      string
	session *discordgo.Session

	handlers           []bot.Handler
	handlersWithPrefix []bot.Handler
	botPrefix          string
}

func NewDiscordBot(botPrefix, token string) bot.Bot {
	return &DiscordBot{
		botPrefix: botPrefix,
		token:     token,
	}
}

func (b *DiscordBot) AddHandlerFuncs(usePrefix bool, fs ...bot.Handler) {
	if usePrefix {
		for i := range fs {
			b.handlersWithPrefix = append(b.handlersWithPrefix, fs[i])
		}
	} else {
		for i := range fs {
			b.handlers = append(b.handlers, fs[i])
		}
	}
}

func (b *DiscordBot) Start() error {
	if len(b.handlers) == 0 {
		return errors.New("empty handler")
	}

	session, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return err
	}

	user, err := session.User("@me")
	if err != nil {
		return err
	}

	b.id = user.ID
	b.session = session

	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == b.id {
			return
		}
		if m.Content[0:1] == b.botPrefix {
			for i := range b.handlersWithPrefix {
				if handle(b.handlersWithPrefix[i], s, m) {
					return
				}
			}
		} else {
			for i := range b.handlers {
				if handle(b.handlers[i], s, m) {
					return
				}
			}
		}

	})

	if err := session.Open(); err != nil {
		return err
	}
	return nil
}

func (b *DiscordBot) Send(channelID, message string) error {
	if _, err := b.session.ChannelMessageSend(channelID, message); err != nil {
		return err
	}
	return nil
}

func handle(h bot.Handler, s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if reply, err := h(m.Content); err != nil {
		err = fmt.Errorf("err: %s, handling %s in channel %s", err.Error(), m.Content, m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return true
	} else if reply != "" {
		fmt.Println(m.ChannelID, reply)
		s.ChannelMessageSend(m.ChannelID, reply)
		return true
	}
	return false
}
