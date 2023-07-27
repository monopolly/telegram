package telegram

import (
	"fmt"
)

func (a *Bot) SendLink(chatID int64, title, url string) {
	if a == nil {
		return
	}
	a.SendMarkdown(chatID, fmt.Sprintf(`[%s](%s)`, title, url))
}
