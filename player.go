package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message = api.Message

func NewPlayer(b *Bot) (a *Player) {
	a = new(Player)
	a.bot = b
	return
}

type Player struct {
	bot  *Bot
	list map[int64]*Message
	m    sync.Mutex
}

// admins only
func (a *Player) Log(text string, args ...interface{}) (err error) {
	if a.bot == nil {
		return errors.New("connection")
	}

	// first init
	if a.list == nil {
		a.list = make(map[int64]*Message)
		for id := range a.bot.admins {
			m := api.NewMessage(int64(id), fmt.Sprintf(text, args...))
			m.ParseMode = api.ModeMarkdown
			m.DisableWebPagePreview = true
			messageID, err := a.bot.bot.Send(m)
			if err != nil {
				continue
			}
			a.list[id] = &messageID
		}
		return
	}

	// send
	a.log(text, args...)
	return
}

// admins only
func (a *Player) log(text string, args ...interface{}) {
	a.m.Lock()
	for _, x := range a.list {
		msg := api.NewEditMessageText(x.Chat.ID, x.MessageID, fmt.Sprintf(text, args...))
		a.bot.bot.Send(msg)
	}
	a.m.Unlock()
}

// any
func (a *Player) Send(user int64, text string, args ...interface{}) (err error) {
	var x api.Message
	a.m.Lock()
	if a.list[user] == nil {
		m := api.NewMessage(int64(user), fmt.Sprintf(text, args...))
		m.ParseMode = api.ModeMarkdown
		m.DisableWebPagePreview = true
		x, err = a.bot.bot.Send(m)
		if err != nil {
			return
		}
		a.list[user] = &x
	} else {
		x = *a.list[user]
	}
	a.m.Unlock()

	msg := api.NewEditMessageText(x.Chat.ID, x.MessageID, fmt.Sprintf(text, args...))
	_, err = a.bot.bot.Send(msg)
	return
}

func LoadPlayer(b *Bot, p []byte) (a *Player, err error) {
	a = NewPlayer(b)
	err = json.Unmarshal(p, &a.list)
	return
}

func (a *Player) Pack() []byte {
	b, _ := json.Marshal(a.list)
	return b
}
