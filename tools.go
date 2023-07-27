package telegram

import (
	"bytes"
	"strings"
)

// разбивает длинные сообщения на части из за лимита телеграма
func split(body string) (lines []string) {
	if len(body) < maxLen {
		return []string{body}
	}

	var b bytes.Buffer
	for _, x := range strings.Split(body, "\n") {
		switch b.Len()+len(x) > maxLen {
		case true:
			lines = append(lines, b.String())
			b.Reset()
			b.WriteString(x + "\n")
		default:
			b.WriteString(x + "\n")
		}
	}
	return
}
