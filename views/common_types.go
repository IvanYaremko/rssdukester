package views

import "github.com/charmbracelet/bubbles/list"

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
