package telegram

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

func (a *Bot) SendMarkdown(chatID int64, text string, buttons ...Button) (err error) {
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
		for _, x := range buttons {
			if x.Handler != nil {
				a.callbacks.list[uuid.New().String()] = x.Handler
			}
		}
		msg.ReplyMarkup = createURLButtons(buttons...)
	}

	_, err = a.bot.Send(msg)
	return
}
