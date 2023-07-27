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
	ID   int
	bot  *Bot
	m    *Message
	user int
}

func (a *Play) Delete() error {
	return a.bot.Delete(a.user, a.ID)
}

// any
func (a *Play) Play(text string, args ...interface{}) (err error) {
	text = fmt.Sprintf(text, args...)
	switch a.m == nil {
	case true:
		m := api.NewMessage(int64(a.user), text)
		m.ParseMode = api.ModeMarkdown
		m.DisableWebPagePreview = true
		m.ReplyMarkup = nil
		out, er := a.bot.bot.Send(m)
		if er != nil {
			return fmt.Errorf(er.Error())
		}
		a.ID = out.MessageID
		a.m = &out
	default:
		msg := api.NewEditMessageText(int64(a.user), a.m.MessageID, text)
		_, err = a.bot.bot.Send(msg)
	}
	return
}

// any
func (a *Play) PlayButtons(text string, buttons ...*Buttons) (err error) {

	switch a.m == nil {
	case true:
		m := api.NewMessage(int64(a.user), text)
		m.ParseMode = api.ModeMarkdown
		m.DisableWebPagePreview = true

		if buttons != nil {
			m.ReplyMarkup = buttons[0].Done()
		}

		out, er := a.bot.bot.Send(m)
		if er != nil {
			return fmt.Errorf(er.Error())
		}
		a.ID = out.MessageID
		a.m = &out
	default:
		m := api.NewEditMessageText(int64(a.user), a.m.MessageID, text)
		if buttons != nil {
			m.ReplyMarkup = buttons[0].Done()
		}
		_, err = a.bot.bot.Send(m)
	}
	return
}
