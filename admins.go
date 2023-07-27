package telegram

import (
	"github.com/monopolly/file"
)

func (a *Bot) AddAdmin(id int) *Bot {
	a.admins[int64(id)] = true
	return a
}

func (a *Bot) Admins() []int64 {
	var list []int64
	a.m.Lock()
	for x := range a.admins {
		list = append(list, x)
	}
	a.m.Unlock()
	return list
}

func (a *Bot) storeAdmins() {
	file.Json(a.file, a.admins)
}
func (a *Bot) loadAdmins() {
	file.LoadJson(a.file, &a.admins)
}

func (a *Bot) Filename(f string) {
	a.file = f
}
