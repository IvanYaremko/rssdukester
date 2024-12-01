package views

import (
	"fmt"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type searchList struct {
	queries    *database.Queries
	searchTerm string
	list       list.Model
	spinner    spinner.Model
	loading    bool
}

func initialiseSearchList(q *database.Queries, st string) searchList {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, feedItemDelegate(), width, height)
	l.Title = fmt.Sprintf("%s %s",
		highlightStyle.Bold(true).Render("SEARCH RESULT FOR"),
		attentionStyle.Italic(true).Render(st),
	)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			enterBinding,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			backBinding,
		}
	}

	return searchList{
		queries:    q,
		list:       l,
		spinner:    s,
		searchTerm: st,
		loading:    true,
	}
}

func (sl searchList) Init() tea.Cmd {
	return sl.spinner.Tick
}

func (sl searchList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return sl, nil
}

func (sl searchList) View() string {
	return ""
}
