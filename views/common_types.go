package views

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type success struct{}

type successContent struct {
	content string
}

type successItems struct {
	items []list.Item
}

type fail struct{}

type failError struct {
	error error
}

type item struct {
	title       string
	description string
	url         string
	feedName    string
	feedId      int64
	pubDate     time.Time
}

func (i item) FilterValue() string { return i.title }
func (i item) Description() string {
	return i.description
}

func (i item) Title() string {
	return i.title
}

type selected struct {
	selected item
}

func transitionView(i item) tea.Cmd {
	return func() tea.Msg {
		return selected{selected: i}
	}
}
