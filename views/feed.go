package views

import (
	"fmt"
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/IvanYaremko/rssdukester/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// xml object
type feedItem struct {
	name string
}

func (f feedItem) FilterValue() string { return f.name }
func (f feedItem) Description() string { return f.name }
func (f feedItem) Title() string       { return f.name }

type feed struct {
	queries *database.Queries
	item    rssItem
	spinner spinner.Model
	loading bool
	list    list.Model
}

func initialiseFeed(q *database.Queries, i rssItem) feed {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.HighlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, list.NewDefaultDelegate(), 30, 30)
	l.Title = fmt.Sprintf("%s FEED", strings.ToUpper(i.name))
	l.Styles.Title = styles.HighlightStyle
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

	return feed{
		queries: q,
		item:    i,
		spinner: s,
		loading: true,
		list:    l,
	}
}

func (f feed) Init() tea.Cmd {
	return f.spinner.Tick
}

func (f feed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return f, tea.Quit
		case key.Matches(msg, backBinding):
			rssList := initialiseRssList(f.queries)
			return rssList, rssList.Init()
		}
	}

	newList, cmd := f.list.Update(msg)
	cmds = append(cmds, cmd)
	f.list = newList

	f.spinner, cmd = f.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return f, tea.Batch(cmds...)
}

func (f feed) View() string {
	s := strings.Builder{}

	if f.loading {
		s.WriteString(fmt.Sprintf("%s loading %s...", f.spinner.View(), f.item.name))
	} else {
		s.WriteString(f.list.View())
	}

	return base.Render(s.String())
}
