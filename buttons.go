package telegram

import api "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CreateButtons(buttons ...Button) (a *Buttons) {
	a = new(Buttons)
	if buttons != nil {
		a.list = append(a.list, buttons)
	}
	return
}

type Buttons struct {
	list [][]Button
}

func (a *Buttons) AddRow(buttons ...Button) {
	a.list = append(a.list, buttons)
}

func (a *Buttons) Add(b Button) {
	if a.list == nil {
		a.list = [][]Button{
			[]Button{
				b,
			},
		}
		return
	}
	a.list[len(a.list)-1] = append(a.list[len(a.list)-1], b)
}

func (a *Buttons) Done() (res *api.InlineKeyboardMarkup) {
	var list [][]api.InlineKeyboardButton
	for _, x := range a.list {
		var row []api.InlineKeyboardButton
		for _, r := range x {
			row = append(row, r.button())
		}
		p := api.NewInlineKeyboardRow(row...)
		list = append(list, p)
	}
	rr := api.NewInlineKeyboardMarkup(list...)
	return &rr
}
