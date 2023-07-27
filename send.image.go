package telegram

import (
	"bytes"
	"fmt"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/telebot.v3"
)

func (a *Bot) SendImage(chatID int64, path string, caption ...string) {
	if a == nil {
		return
	}
	msg := api.NewPhoto(chatID, api.FilePath(path))
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	a.bot.Send(msg)
}

// images, video
func (a *Bot) SendImagesFiles(chatID int64, filenames ...string) {
	var files []telebot.File
	for _, x := range filenames {
		files = append(files, telebot.FromDisk(x))
	}
	a.sendAlbum(chatID, files...)
}

// images, video
func (a *Bot) SendImageBytes(chatID int64, image []byte) (err error) {
	return a.sendImage(chatID, telebot.FromReader(bytes.NewBuffer(image)))
}

// images, video
func (a *Bot) SendImagesBytes(chatID int64, images ...[]byte) {
	if images == nil {
		return
	}
	var files []telebot.File
	for _, x := range images {
		files = append(files, telebot.FromReader(bytes.NewBuffer(x)))
	}
	a.sendAlbum(chatID, files...)
}

// images, video
func (a *Bot) SendImagesURLs(chatID int64, urls ...string) {
	var files []telebot.File
	for _, x := range urls {
		files = append(files, telebot.FromURL(x))
	}
	a.sendAlbum(chatID, files...)
}

// images, video
func (a *Bot) sendAlbum(chatID int64, list ...telebot.File) {
	if a == nil {
		return
	}
	user := &telebot.User{ID: chatID}
	var files telebot.Album
	for _, x := range list {
		files = append(files, &telebot.Photo{File: x})
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

// images, video
func (a *Bot) sendImage(chatID int64, image telebot.File) (err error) {
	if a == nil {
		return
	}
	user := &telebot.User{ID: chatID}
	_, err = a.filebot.Send(user, &telebot.Photo{File: image})
	return
}
