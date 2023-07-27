package telegram

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *Bot) SendMarkdown(chatID int64, text string, buttons ...*Buttons) (err error) {
	if a == nil {
		return
	}
	/*
		*bold text*
		_italic text_
		[inline URL](http://www.example.com/)
		[inline mention of a user](tg://user?id=123456789)
		pre-formatted fixed-width code block
	*/
	msg := api.NewMessage(chatID, text)
	msg.ParseMode = api.ModeMarkdown
	msg.DisableWebPagePreview = true

	if buttons != nil {
		msg.ReplyMarkup = buttons[0].Done()
	}

	_, err = a.bot.Send(msg)
	return
}
