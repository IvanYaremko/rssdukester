package views

import (
	"fmt"
	"strings"

	"github.com/IvanYaremko/rssdukester/bindings"
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
	item    item
	spinner spinner.Model
	loading bool
	list    list.Model
	keys    bindings.ListKeysMap
}

func initialiseFeed(q *database.Queries, i item) feed {
	s := spinner.New()
	s.Spinner = spinner.Dot

	items := make([]list.Item, 0)
	l := list.New(items, list.NewDefaultDelegate(), 30, 30)
	l.Title = fmt.Sprintf("%s FEED", strings.ToUpper(i.name))
	l.Styles.Title = styles.HighlightStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			bindings.ListItemDelegateKeys.Choose,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			bindings.ListKeys.Back,
		}
	}

	return feed{
		queries: q,
		item:    i,
		spinner: s,
		loading: true,
		keys:    bindings.ListKeys,
		list:    l,
	}
}

func (f feed) Init() tea.Cmd {
	return nil
}

func (f feed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keys.Ctrlc):
			return f, tea.Quit
		case key.Matches(msg, f.keys.Back):
			rssList := initialiseRssList(f.queries)
			return rssList, rssList.Init()
		}
	}

	newList, cmd := f.list.Update(msg)
	cmds = append(cmds, cmd)
	f.list = newList

	return f, tea.Batch(cmds...)
}

func (f feed) View() string {
	s := strings.Builder{}

	if f.loading {
		s.WriteString(fmt.Sprintf("%s loading...", f.spinner.View()))
	} else {
		s.WriteString(f.list.View())
	}

	return base.Render(s.String())
}
