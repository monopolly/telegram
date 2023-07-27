package telegram

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
