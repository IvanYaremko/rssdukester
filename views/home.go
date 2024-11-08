package views

import (
	"github.com/IvanYaremko/rssdukester/sql/database"
	tea "github.com/charmbracelet/bubbletea"
)

type Home struct {
	queries *database.Queries
}

func (h Home) Init() tea.Cmd {

	return tea.SetWindowTitle("RSSDUKESTER")
}

func (h Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	return h, tea.Batch(cmds...)
}

func (h Home) View() string {

	return ""
}

func InitHomeModel(q *database.Queries) Home {

	return Home{queries: q}
}
