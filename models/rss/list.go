package rss

import (
	"github.com/IvanYaremko/rssdukester/sql/database"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewModel struct {
	dbQueries database.Queries
}

func (v ViewModel) Init() tea.Cmd {
	return nil
}

func (v ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return v, nil
}

func (v ViewModel) View() string {
	return ""
}
