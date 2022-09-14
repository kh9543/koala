package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kh9543/koala/domain/bot"
	"github.com/kh9543/koala/domain/bot/discord"
	"github.com/kh9543/koala/domain/exchangerate"
	"github.com/kh9543/koala/domain/kv"
)

type Bot struct {
	core bot.Bot
	kv   kv.Kv
}

func NewDiscordBot(prefix, token string, kv kv.Kv) error {
	bot := Bot{
		core: discord.NewDiscordBot(prefix, token),
		kv:   kv,
	}

	bot.core.AddHandler(
		true,
		bot.pingHandler,
		bot.currencyHandler,
		bot.koalaBrainHandler,
	)

	bot.core.AddHandler(
		false,
		bot.koalaHandler,
	)

	if err := bot.core.Start(); err != nil {
		return err
	}
	return nil
}

func (b *Bot) pingHandler(msg string) (string, error) {
	if msg == "!ping" {
		return "pong", nil
	}
	return "", nil
}

func (b *Bot) currencyHandler(msg string) (string, error) {
	argv := strings.Split(msg, " ")
	if argv[0] != "!currency" {
		return "", nil
	}

	if len(argv) != 2 {
		return "usage: !currency <nation>", nil
	}

	buyrate, sellRate, err := exchangerate.GetRate(argv[1])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("現金買入: %s, 現金賣出: %s", buyrate, sellRate), nil
}

func (b *Bot) koalaBrainHandler(msg string) (string, error) {
	argv := strings.Split(msg, " ")
	if argv[0] != "!koala" {
		return "", nil
	}
	if len(argv) != 3 {
		return "usage: !koala <pattern> <response>", nil
	}

	if argv[1] == "!koala" {
		return "", nil
	}

	if err := b.kv.Add("koala", argv[1], argv[2]); err != nil {
		return "", err
	}
	return "👍", nil
}

func (b *Bot) koalaHandler(msg string) (string, error) {
	mp, err := b.kv.GetAll("koala")
	if err != nil {
		return "", err
	}

	ks := make([]string, 0, len(mp))
	for k, _ := range mp {
		ks = append(ks, k)
	}
	sort.Slice(ks, func(i, j int) bool {
		return len(ks[i]) > len(ks[j])
	})

	for _, k := range ks {
		if strings.Contains(msg, k) {
			return mp[k].(string), nil
		}
	}

	return "", nil
}
