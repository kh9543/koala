package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/kh9543/koala/domain/bot"
	"github.com/kh9543/koala/domain/bot/discord"
	"github.com/kh9543/koala/domain/chatgpt"
	"github.com/kh9543/koala/domain/constant"
	"github.com/kh9543/koala/domain/exchangerate"
	"github.com/kh9543/koala/domain/kv"
	"github.com/kh9543/koala/domain/stringmatch"
)

type Bot struct {
	bot bot.Bot
	kv  kv.Kv
}

func NewDiscordBot(prefix, token string, kv kv.Kv) bot.Bot {
	b := &Bot{
		bot: discord.NewDiscordBot(prefix, token),
		kv:  kv,
	}

	b.bot.AddHandlerFuncs(
		true,
		b.pingHandler,
		b.currencyHandler,
		b.koalaBrainHandler,
		b.koalaDeleteHandler,
		b.chatgptControlHandler,
	)

	b.bot.AddHandlerFuncs(
		false,
		b.koalaHandler,
	)

	b.bot.AddReactionHandlerFuncs(
		b.pollHandler,
	)

	b.bot.AddChannelMsgHandlerFuncs(
		b.chatgptHandler,
	)

	return b.bot
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

func (b *Bot) koalaDeleteHandler(msg string) (string, error) {
	argv := strings.Split(msg, " ")
	if argv[0] != "!delete" {
		return "", nil
	}
	if len(argv) != 2 {
		return "usage: !delete <key>", nil
	}

	if err := b.kv.Delete("koala", argv[1]); err != nil {
		return "", err
	}
	return "👍", nil
}

func (b *Bot) pollHandler(msg string) (string, int, error) {
	argv := strings.Split(msg, " ")
	if argv[0] != "!poll" {
		return "", -1, nil
	}
	if len(argv) < 2 {
		return "usage: !poll <title> [<choice>...]", -1, nil
	} else if len(argv) > 2 {
		response := fmt.Sprintf("📊 %s\n", argv[1])
		for i, option := range argv[2:] {
			response += fmt.Sprintf("%s %s\n", string(rune(i+0x1F1E6)), option)
		}
		return response, len(argv) - 2, nil
	}

	return fmt.Sprintf("📊 %s", argv[1]), 0, nil
}

func (b *Bot) koalaHandler(msg string) (string, error) {
	mp, err := b.kv.GetAll("koala")
	if err != nil {
		return "", err
	}

	ks := make([]string, 0, len(mp))
	for k := range mp {
		ks = append(ks, k)
	}

	matchedKeys := stringmatch.MatchString(msg, ks)
	if len(matchedKeys) == 0 {
		return "", nil
	}

	response := make([]string, 0, len(matchedKeys))
	for _, key := range matchedKeys {
		result, _ := mp[key].(string)
		// If there are more and more replaced words, add a new collection in mongodb
		replaced_result := strings.ReplaceAll(result, "time.Now()", fmt.Sprintf(" %s ", time.Now().Format("2006-01-02 15:04:05")))
		response = append(response, replaced_result)
	}

	return strings.Join(response, " "), nil
}

func (b *Bot) chatgptHandler(msg, channelID, userID string) (string, error) {
	if channelID != string(constant.ChatChannel) {
		return "", nil
	}

	ans, err := chatgpt.SendQuestion(msg)
	if err != nil {
		return "", err
	}
	return ans, nil
}

func (b *Bot) chatgptControlHandler(msg string) (string, error) {
	argv := strings.Split(msg, " ")
	if argv[0] != "!chatgpt" {
		return "", nil
	}
	if len(argv) != 2 {
		return "usage: !chatgpt <system_msg|delete>", nil
	}

	if msg == "delete" {
		msg = ""
	}

	chatgpt.SystemMessage = msg
	return "👍", nil
}
