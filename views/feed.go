package views

import (
	"github.com/IvanYaremko/rssdukester/sql/database"
	tea "github.com/charmbracelet/bubbletea"
)

type feed struct {
	queries *database.Queries
	item    item
}

func InitialiseFeed(q *database.Queries, i item) feed {
	return feed{
		queries: q,
		item:    i,
	}
}

func (f feed) Init() tea.Cmd {
	return nil
}

func (f feed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return f, nil
}

func (f feed) View() string {
	return ""
}
