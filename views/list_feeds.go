package views

import (
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	Name, Url, UpdatedAt string
}

func (i Item) FilterValue() string { return i.Name }

type listFeeds struct {
	queries *database.Queries
	list    list.Model
}

func InitialiseListFeeds(q *database.Queries) listFeeds {
	items := make([]list.Item, 0)

	feedsList := list.New(items, list.NewDefaultDelegate(), 40, 50)
	l := listFeeds{
		queries: q,
		list:    feedsList,
	}

	return l
}

func (l listFeeds) Init() tea.Cmd {

	return nil
}
