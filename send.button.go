package telegram

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *Bot) SendButton(chatID int64, text, title, link string) (res api.Message) {
	if text == "" {
		text = title
	}
	button := api.NewInlineKeyboardMarkup(
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonURL(title, link),
		),
	)
	msg := api.NewMessage(chatID, text)
	msg.ReplyMarkup = button
	res, _ = a.bot.Send(msg)
	return
}
