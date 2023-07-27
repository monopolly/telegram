package telegram

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *Bot) SendSticker(chatID int64, id string) {
	if a == nil {
		return
	}
	msg := api.NewSticker(chatID, api.FileID(id))
	a.bot.Send(msg)
}
