package telegram

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *Bot) SendFileBytes(chatID int64, data []byte, caption ...string) (err error) {

	f := api.FileBytes{Bytes: data}
	if len(caption) > 0 {
		f.Name = caption[0]
	} else {
		f.Name = "File"
	}
	msg := api.NewDocument(chatID, f)
	_, err = a.bot.Send(msg)
	return
}
