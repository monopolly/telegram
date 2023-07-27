package telegram

import (
	"fmt"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewSimplePlay(b *Bot, user int) (a *Play) {
	a = new(Play)
	a.bot = b
	a.user = user
	return
}

type Play struct {
	bot  *Bot
	m    *Message
	user int
}

// any
func (a *Play) Play(text string, args ...interface{}) (err error) {
	text = fmt.Sprintf(text, args...)
	switch a.m == nil {
	case true:
		m := api.NewMessage(int64(a.user), text)
		m.ParseMode = api.ModeMarkdown
		m.DisableWebPagePreview = true
		out, er := a.bot.bot.Send(m)
		if er != nil {
			return fmt.Errorf(er.Error())
		}
		a.m = &out
	default:
		msg := api.NewEditMessageText(int64(a.user), a.m.MessageID, fmt.Sprintf(text, args...))
		_, err = a.bot.bot.Send(msg)
	}
	return
}
