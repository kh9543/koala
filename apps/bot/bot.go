package bot

import (
	"fmt"
	"strings"

	"github.com/kh9543/bot/domain/bot"
	"github.com/kh9543/bot/domain/bot/discord"
	"github.com/kh9543/bot/domain/exchangerate"
)

func NewDiscordBot(prefix, token string) bot.Bot {
	bot := discord.NewDiscordBot(prefix, token)
	bot.AddHandler(pingHandler, currencyHandler)
	return bot
}

func pingHandler(msg string) (string, error) {
	if msg == "!ping" {
		return "pong", nil
	}
	return "", nil
}

func currencyHandler(msg string) (string, error) {
	argv := strings.Split(msg, " ")
	if argv[0] != "!currency" {
		return "", nil
	}

	if len(argv) != 2 {
		return "usage: !currency <nation>", nil
	}

	buyrate, sellRate, err := exchangerate.GetRate(argv[1])
	if err != nil {
		return err.Error(), nil
	}
	return fmt.Sprintf("現金買入: %s, 現金賣出: %s", buyrate, sellRate), nil
}
