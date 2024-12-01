package views

import (
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type search struct {
	quries      *database.Queries
	spinner     spinner.Model
	loading     bool
	list        list.Model
	searchInput textinput.Model
	err         error
	searchTerm  string
}

func initialiseSearch(q *database.Queries) search {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, feedItemDelegate(), width, height)
	l.Title = "SEARCH RESULTS"
	l.Styles.Title = highlightStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			enterBinding,
		}
	}

	ti := textinput.New()
	ti.Placeholder = "Enter search terms..."
	ti.Focus()
	ti.Width = width

	return search{
		quries:      q,
		spinner:     s,
		loading:     true,
		list:        l,
		searchInput: ti,
	}
}

func (s search) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return s, nil
}

func (s search) View() string {

	return ""
}
