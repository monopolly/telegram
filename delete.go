package telegram

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *Bot) Delete(user, message int) (err error) {
	if a == nil {
		return
	}

	p := api.DeleteMessageConfig{ChatID: int64(user), MessageID: message}
	_, err = a.bot.Send(p)
	return
}
