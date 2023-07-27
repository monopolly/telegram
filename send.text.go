package telegram

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *Bot) SendText(chatID int64, text string) (res api.Message) {
	if a == nil {
		return
	}
	//max len 4000
	if len(text) > maxLen {
		for _, x := range split(text) {
			msg := api.NewMessage(chatID, x)
			res, _ = a.bot.Send(msg)
		}
		return
	}
	msg := api.NewMessage(chatID, text)
	res, _ = a.bot.Send(msg)
	return
}
