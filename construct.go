package telegram

import (
	"fmt"
	"strings"
)

type msg struct {
	lines []string
}

// markdown
func Text() *msg {
	return new(msg)
}

func (a *msg) Title(k string, v ...interface{}) *msg {
	a.lines = append(a.lines, fmt.Sprintf("*%s*", fmt.Sprintf(k, v...)))
	return a
}
func (a *msg) Text(k string, v ...interface{}) *msg {
	a.lines = append(a.lines, fmt.Sprintf(k, v...))
	return a
}

func (a *msg) Italic(k string, v ...interface{}) *msg {
	k = "_" + k + "_"
	a.lines = append(a.lines, fmt.Sprintf(k, v...))
	return a
}

func (a *msg) Div() *msg {
	a.lines = append(a.lines, "")
	return a
}
func (a *msg) String() string {
	return strings.Join(a.lines, "\n")
}
