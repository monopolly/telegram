package telegram

import api "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func createURLButtons(list ...Button) (res interface{}) {

	var row []api.InlineKeyboardButton

	for _, x := range list {
		row = append(row, x.button())
	}

	return api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(row...))
}

type Button struct {
	Title   string
	Link    string
	Handler func(m *Message)
}

func (a *Button) button() api.InlineKeyboardButton {
	return api.NewInlineKeyboardButtonURL(a.Title, a.Link)
}

// button :=
