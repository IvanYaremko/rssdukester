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
