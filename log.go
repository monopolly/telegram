package telegram

import (
	"fmt"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// safe
func (a *Bot) Log(text string, arg ...interface{}) {
	if a.bot == nil {
		return
	}
	mm := fmt.Sprintf(text, arg...)
	a.m.Lock()
	for id := range a.admins {
		//max len 4000
		if len(mm) > maxLen {
			for _, x := range split(mm) {
				msg := api.NewMessage(int64(id), x)
				a.bot.Send(msg)
			}
			return
		} else {
			m := api.NewMessage(int64(id), mm)
			m.ParseMode = api.ModeMarkdown
			m.DisableWebPagePreview = true
			a.bot.Send(m)
		}

	}
	a.m.Unlock()
}
