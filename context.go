package telegram

import (
	"fmt"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/monopolly/file"
)

func newContext(u api.Update) (a *Context) {
	a = new(Context)
	a.c = u
	a.User = int(u.Message.From.ID)
	a.Body = u.Message.Text
	a.Login = u.Message.From.UserName
	a.Command = u.Message.Command()
	a.Args = u.Message.CommandArguments()
	a.Message = u.Message
	a.keys = make(map[string]interface{})
	a.Callback = u.CallbackQuery
	return
}

type Context struct {
	User     int
	Login    string
	Body     string
	Command  string
	Args     string
	Admin    bool
	Working  bool
	Callback *api.CallbackQuery
	Message  *api.Message
	c        api.Update
	bot      *Bot
	keys     map[string]interface{}
}

func (a *Context) Get(key string) interface{} {
	return a.keys[key]
}

func (a *Context) CreateDataButton(title, data string, handler func(*Callback)) (r Button) {
	a.bot.callbacks.Lock()
	a.bot.callbacks.list[data] = handler
	a.bot.callbacks.Unlock()
	return Button{Title: title, Data: data}
}

func (a *Context) Set(key string, value interface{}) {
	a.keys[key] = value
}

func (a *Context) Text(text string, v ...interface{}) (res api.Message) {
	msg := api.NewMessage(a.c.Message.Chat.ID, fmt.Sprintf(text, v...))
	res, _ = a.bot.bot.Send(msg)
	return
}

func (a *Context) SendfileBytes(data []byte, caption ...string) (err error) {
	f := api.FileBytes{Bytes: data}
	if len(caption) > 0 {
		f.Name = caption[0]
	} else {
		f.Name = "upfile"
	}
	msg := api.NewDocument(a.c.Message.Chat.ID, f)
	_, err = a.bot.bot.Send(msg)
	return
}

/*
func (a *Context) SendImage(path string, caption ...string) {
	msg := api.NewPhotoUpload(a.c.Message.Chat.ID, path)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	a.bot.Send(msg)
}

func (a *Context) SendImageBytes(f []byte, caption ...string) {
	msg := api.NewPhotoUpload(a.c.Message.Chat.ID, f)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	a.bot.Send(msg)
}

func (a *Context) Sendfile(path string, caption ...string) (err error) {
	msg := api.NewDocumentUpload(a.c.Message.Chat.ID, path)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	_, err = a.bot.Send(msg)
	return
}

func (a *Context) SendfileBytes(f interface{}, caption ...string) (err error) {
	msg := api.NewDocumentUpload(a.c.Message.Chat.ID, f)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	_, err = a.bot.Send(msg)
	return
} */

func (a *Context) Download(id string) (body []byte, err error) {
	c := api.FileConfig{FileID: id}
	f, err := a.bot.bot.GetFile(c)
	if err != nil {
		return
	}
	link := f.Link(a.bot.bot.Token)
	body, _ = file.Get(link)
	return
}

func (a *Context) Markdown(text string, preview ...bool) {
	/*
		*bold text*
		_italic text_
		[inline URL](http://www.example.com/)
		[inline mention of a user](tg://user?id=123456789)
		pre-formatted fixed-width code block

		messageEntityBold => <b>bold</b>, <strong>bold</strong>, **bold**
		messageEntityItalic => <i>italic</i>, <em>italic</em> *italic*
		messageEntityCode => <code>code</code>, `code`
		messageEntityStrike => <s>strike</s>, <strike>strike</strike>, <del>strike</del>, ~~strike~~
		messageEntityUnderline => <u>underline</u>
		messageEntityPre => <pre language="c++">code</pre>,
		```c++
		code
		```
	*/
	msg := api.NewMessage(a.c.Message.Chat.ID, text)
	msg.ParseMode = api.ModeMarkdown
	if preview != nil {
		msg.DisableWebPagePreview = preview[0]
	}
	a.bot.bot.Send(msg)
}

func (a *Context) HTML(text string) {
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
	msg := api.NewMessage(a.c.Message.Chat.ID, text)
	msg.ParseMode = api.ModeHTML
	a.bot.bot.Send(msg)
}

func (a *Context) Link(title, url string, preview ...bool) {
	a.Markdown(fmt.Sprintf(`[%s](%s)`, title, url), preview...)
}

// Button with markdown
func (a *Context) Button(message, title, link string) {
	//msg := api.NewInlineKeyboardButtonURL(title, link) //api.NewEditMessageReplyMarkup(a.c.Message.Chat.ID, c.Message.MessageID, text)
	//msg := api.NewInlineKeyboardButtonData(title, link)
	msg := api.NewMessage(a.c.Message.Chat.ID, message)
	msg.ReplyMarkup = api.NewInlineKeyboardMarkup(
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonURL(title, link),
		),
	)
	msg.ParseMode = api.ModeMarkdown
	_, err := a.bot.bot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(err)
}

func (a *Context) Channel(name, text string) {
	msg := api.NewMessageToChannel(name, text)
	a.bot.bot.Send(msg)
}

func (a *Context) Reply(text string) (res api.Message) {
	msg := api.NewMessage(a.c.Message.Chat.ID, text)
	msg.ReplyToMessageID = a.c.Message.MessageID
	res, _ = a.bot.bot.Send(msg)
	return
}

func (a *Context) Play(res api.Message, newtext string) {
	msg := api.NewEditMessageText(res.Chat.ID, res.MessageID, newtext)
	a.bot.bot.Send(msg)
}

/* func (a *Context) Buttons(title, link string) {
	msg := api.NewInlineKeyboardButtonURL(title, link) //api.NewEditMessageReplyMarkup(a.c.Message.Chat.ID, c.Message.MessageID, text)
	_, err := a.bot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
} */

/* func (c *Context) Button(title, link string) {
	msg := api.NewEditMessageReplyMarkup(a.c.Message.Chat.ID, c.Message.MessageID text)

	res, _ = a.bot.Send(msg)
	return
}

type message struct {
	c        *Context
	msg      api.MessageConfig
	text     string
	title    string
	link     string
	inline   []api.InlineKeyboardButton
	keyboard []api.KeyboardButton
}

func (a *message) Text(text string) *message {
	a.text = text
	return a
}

func (a *message) Link(title, link string) *message {
	a.title = title
	a.link = link
	return a
}

func (a *message) Button(title, link string) *message {
	a.list = append(list.list, api.NewInlineKeyboardButtonURL(title, link))
	return a
}

func (a *message) Send() api.MessageConfig {

	return a
}
*/
