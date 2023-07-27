package telegram

import "gopkg.in/telebot.v3"

// func NewMessage(b *Bot) (a *message) {
// 	a = new(message)
// 	a.mode = "text"
// 	return
// }

type message struct {
	mode    string
	id      string
	text    string
	link    string
	images  []telebot.File
	album   []telebot.Album
	files   []telebot.File
	video   []telebot.Video
	audio   []telebot.Audio
	buttons [][]button
}

type button struct {
	title   string
	link    string
	handler func()
}
