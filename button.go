package telegram

import api "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func createButtons(list ...Button) (res *api.InlineKeyboardMarkup) {

	var row []api.InlineKeyboardButton

	for _, x := range list {
		row = append(row, x.button())
	}

	b := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(row...))
	return &b
}

type Button struct {
	Title string
	Link  string
	Data  string
}

func (a *Button) button() api.InlineKeyboardButton {
	switch a.Data != "" {
	case true:
		return api.NewInlineKeyboardButtonData(a.Title, a.Data)
	default:
		return api.NewInlineKeyboardButtonURL(a.Title, a.Link)
	}
}

func CreateCallbackButton(bot *Bot, title, data string, handler func(*Callback)) (r Button) {
	bot.callbacks.Lock()
	bot.callbacks.list[data] = handler
	bot.callbacks.Unlock()
	return Button{Title: title, Data: data}
}
