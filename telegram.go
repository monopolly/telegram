package telegram

import (
	"fmt"
	"log"
	"strings"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tucnak/telebot"

	"github.com/monopolly/file"
)

func New(token string, admins ...int) (a *Bot, err error) {
	a = new(Bot)

	a.bot, err = api.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return
	}

	/* switch client != nil {
	case true:

	default:
		a.bot, err = api.NewBotAPIWithClient(token, api.APIEndpoint, client[0])
		if err != nil {
			log.Println(err)
			return
		}
	} */

	a.admins = make(map[int]bool)
	for _, x := range admins {
		a.admins[x] = true
	}
	//log.Println("telegram", a.bot.Self.UserName)
	a.filebot, _ = telebot.NewBot(telebot.Settings{
		Token: token,
	})
	return
}

type Bot struct {
	bot     *api.BotAPI
	router  func(*Context)
	admins  map[int]bool
	pass    string
	file    string
	filebot *telebot.Bot
}

func (a *Bot) Name() string {

	return a.bot.Self.UserName
}

//admin password and admin json store
func (a *Bot) Pass(pass string) *Bot {
	a.pass = pass
	/* a.file = fmt.Sprintf("tg_%s.json", a.bot.Self.UserName)
	a.loadAdmins() */
	return a
}

func (a *Bot) AddAdmin(id int) *Bot {
	a.admins[id] = true
	return a
}

func (a *Bot) Admins() []int {
	var list []int
	for x := range a.admins {
		list = append(list, x)
	}
	return list
}

func (a *Bot) storeAdmins() {
	file.Json(a.file, a.admins)
}
func (a *Bot) loadAdmins() {
	file.LoadJson(a.file, &a.admins)
}

func (a *Bot) Log(text string, arg ...interface{}) {
	if a.bot == nil {
		return
	}
	for id := range a.admins {
		m := api.NewMessage(int64(id), fmt.Sprintf(text, arg...))
		m.ParseMode = api.ModeMarkdown
		m.DisableWebPagePreview = true
		a.bot.Send(m)
	}
}

type play struct {
	bot  *Bot
	list []api.Message
}

func (a *Bot) LogPlay(text string, args ...interface{}) (p *play) {
	if a.bot == nil {
		return
	}
	p = new(play)
	p.bot = a

	for id := range a.admins {
		m := api.NewMessage(int64(id), fmt.Sprintf(text, args...))
		m.ParseMode = api.ModeMarkdown
		m.DisableWebPagePreview = true
		ms, err := a.bot.Send(m)
		if err != nil {
			continue
		}
		p.list = append(p.list, ms)
	}

	return
}

func (a *play) Send(text string, args ...interface{}) {
	for _, x := range a.list {
		msg := api.NewEditMessageText(x.Chat.ID, x.MessageID, fmt.Sprintf(text, args...))
		a.bot.bot.Send(msg)
	}
}

func (a *play) Store(filename string)   {}
func (a *play) Restore(filename string) {}

func (a *Bot) Filename(f string) {
	a.file = f
}

func (a *Bot) Bot() *api.BotAPI {
	return a.bot
}

func (a *Bot) Channel(ch, text string) (err error) {
	if a == nil {
		return
	}
	if !strings.HasPrefix(ch, "@") {
		ch = "@" + ch
	}
	b := api.NewMessageToChannel(ch, text)
	b.ParseMode = api.ModeMarkdown
	_, err = a.bot.Send(b)
	return
}

func (a *Bot) SendMarkdown(chatID int64, text string) (err error) {
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
	_, err = a.bot.Send(msg)

	return
}

func (a *Bot) SendHTML(chatID int64, text string) (err error) {
	if a == nil {
		return
	}
	/*
		<b>bold</b>,
		<strong>bold</strong>
		<i>italic</i>,
		<em>italic</em>
		<a href="http://www.example.com/">inline URL</a>
		<a href="tg://user?id=123456789">inline mention of a user</a>
		<code>inline fixed-width code</code>
		<pre>pre-formatted fixed-width code block</pre>
	*/
	msg := api.NewMessage(chatID, text)
	msg.ParseMode = api.ModeHTML
	_, err = a.bot.Send(msg)

	return
}

func (a *Bot) SendLink(chatID int64, title, url string) {
	if a == nil {
		return
	}
	a.SendMarkdown(chatID, fmt.Sprintf(`[%s](%s)`, title, url))
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
	body, _ = file.Downloads(link)
	return
}

/* func (a *Bot) SendButton(chatID int64, title, link string) {
	msg := api.NewInlineKeyboardButtonURL() .NewMessage(chatID, text)
	res, _ = a.bot.Send(msg)
} */

func (a *Bot) SendText(chatID int64, text string) (res api.Message) {
	if a == nil {
		return
	}
	msg := api.NewMessage(chatID, text)
	res, _ = a.bot.Send(msg)
	return
}

func (a *Bot) SendImage(chatID int64, path string, caption ...string) {
	if a == nil {
		return
	}
	msg := api.NewPhotoUpload(chatID, path)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	a.bot.Send(msg)
}

func (a *Bot) SendSticker(chatID int64, id string) {
	if a == nil {
		return
	}
	msg := api.NewStickerShare(chatID, id)
	a.bot.Send(msg)
}

//images, video
func (a *Bot) SendImages(chatID int64, url ...string) {
	if a == nil {
		return
	}
	user := &telebot.User{ID: int(chatID)}
	var files telebot.Album
	for _, x := range url {
		files = append(files, &telebot.Photo{File: telebot.FromDisk(x)})
		if len(files) >= 10 {
			_, err := a.filebot.SendAlbum(user, files)
			if err != nil {
				fmt.Println("send images", err)
				continue
			}
			files = nil
		}
	}
	_, err := a.filebot.SendAlbum(user, files)
	if err != nil {
		fmt.Println("send images", err)
	}

}

func (a *Bot) SendImageBytes(chatID int64, f []byte, caption ...string) {
	if a == nil {
		return
	}
	msg := api.NewPhotoUpload(chatID, f)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	a.bot.Send(msg)
}

func (a *Bot) Sendfile(chatID int64, path string, caption ...string) {
	if a == nil {
		return
	}
	msg := api.NewDocumentUpload(chatID, path)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	a.bot.Send(msg)
}

func (a *Bot) SendfileBytes(chatID int64, f interface{}, caption ...string) (err error) {
	if a == nil {
		return
	}
	msg := api.NewDocumentUpload(chatID, f)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	_, err = a.bot.Send(msg)
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

	updates, err := a.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

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

		c := newContext(update)
		c.bot = a.bot
		if a.admins[c.c.Message.From.ID] {
			c.Admin = true
		}

		if a.router != nil {
			go a.router(c)
		}

	}
}
