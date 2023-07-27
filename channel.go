package telegram

import (
	"strings"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *Bot) Channel(ch, text string) (err error) {
	if a == nil {
		return
	}
	if !strings.HasPrefix(ch, "@") {
		ch = "@" + ch
	}
	b := api.NewMessageToChannel(ch, text)
	b.ParseMode = api.ModeMarkdown
	_, err = a.bot.Send(b)
	return
}
