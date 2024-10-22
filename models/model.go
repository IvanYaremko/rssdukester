package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}
