package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type success struct{}

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
}

func (i item) FilterValue() string { return i.title }
func (i item) Description() string { return i.description }
func (i item) Title() string       { return i.title }

type selectedFeed struct {
	selected item
}

func transitionView(i item) tea.Cmd {
	return func() tea.Msg {
		return selectedFeed{selected: i}
	}
}
