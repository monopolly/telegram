package telegram

import (
	"log"
	"sync"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/telebot.v3"

	"github.com/monopolly/file"
)

type Callback = api.CallbackQuery

const maxLen = 4000 //telegram

func New(token string, admins ...int64) (a *Bot, err error) {
	a = new(Bot)

	a.bot, err = api.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return
	}

	a.admins = make(map[int64]bool)
	for _, x := range admins {
		a.admins[x] = true
	}

	a.filebot, _ = telebot.NewBot(telebot.Settings{
		Token: token,
	})

	a.callbacks.list = make(map[string]func(*Callback))
	return
}

type Bot struct {
	bot     *api.BotAPI
	router  func(*Context)
	admins  map[int64]bool
	m       sync.Mutex
	pass    string
	file    string
	filebot *telebot.Bot

	callbacks struct {
		sync.Mutex
		list map[string]func(*Callback)
	}
}

func (a *Bot) Name() string {
	return a.bot.Self.UserName
}

func (a *Bot) Player() *Player {
	return NewPlayer(a)
}

// admin password and admin json store
func (a *Bot) Pass(pass string) *Bot {
	a.pass = pass
	/* a.file = fmt.Sprintf("tg_%s.json", a.bot.Self.UserName)
	a.loadAdmins() */
	return a
}

func (a *Bot) Bot() *api.BotAPI {
	return a.bot
}

func (a *Bot) GetFileLink(id string) (link string) {
	if a == nil {
		return
	}
	f, err := a.bot.GetFile(api.FileConfig{id})
	if err != nil {
		return
	}
	return f.Link(a.bot.Token)
}

func (a *Bot) DownloadFile(id string) (body []byte) {
	if a == nil {
		return
	}
	link := a.GetFileLink(id)
	if link == "" {
		return
	}
	body, _ = file.Get(link)
	return
}

func (a *Bot) Start(router ...func(*Context)) {
	if a == nil {
		return
	}

	if len(router) > 0 {
		a.router = router[0]
	}
	//a.router = router
	if a.file != "" {
		a.loadAdmins()
	}
	u := api.NewUpdate(0)
	u.Timeout = 60

	updates := a.bot.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil && update.InlineQuery == nil { // ignore any non-Message Updates
			continue
		}
		// else if update.CallbackQuery != nil {
		// 	// Respond to the callback query, telling Telegram to show the user
		// 	// a message with the data received.
		// 	callback := api.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		// 	if _, err := a.bot.Request(callback); err != nil {
		// 		fmt.Println(err)
		// 		continue
		// 	}

		// 	// And finally, send a message containing the data received.
		// 	msg := api.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		// 	if _, err := a.bot.Send(msg); err != nil {
		// 		fmt.Println(err)
		// 		continue
		// 	}
		// }

		if a.pass != "" {
			switch update.Message.Text {
			case a.pass:
				a.admins[update.Message.From.ID] = true
				a.storeAdmins()
				msg := api.NewMessage(int64(update.Message.From.ID), "🔥 Admin logged")
				a.bot.Send(msg)
				continue
			default:
			}
		}

		if update.CallbackQuery != nil {
			go a.callback(update)
			continue
		}

		c := newContext(update)
		c.bot = a
		c.Admin = a.Admin(int(c.Message.From.ID))
		if a.router != nil {
			go a.router(c)
		}

	}
}

func (a *Bot) CreateDataButton(title, data string, handler func(*Callback)) (r Button) {
	a.callbacks.Lock()
	a.callbacks.list[data] = handler
	a.callbacks.Unlock()
	return Button{Title: title, Data: data}
}

func (a *Bot) callback(c api.Update) {
	defer a.callbackReply(c.CallbackQuery.ID)
	var f func(*Callback)
	a.callbacks.Lock()
	f = a.callbacks.list[c.CallbackQuery.Data]
	a.callbacks.Unlock()
	if f == nil {
		//have to reply
		return
	}
	f(c.CallbackQuery)

}

func (a *Bot) callbackReply(id string) {
	callback := api.NewCallback(id, "")
	a.bot.Request(callback)

	// msg := api.NewMessage(int64(chat), "")
	// msg.ReplyToMessageID = messageID
	// a.bot.Send(msg)
}

func (a *Bot) DeleteCallback(data string) {
	delete(a.callbacks.list, data)
}

func (a *Bot) Admin(id int) (res bool) {
	a.m.Lock()
	res = a.admins[int64(id)]
	a.m.Unlock()
	return
}
