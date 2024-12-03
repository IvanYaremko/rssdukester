package views

import (
	"fmt"
	"strings"

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
		"SEARCH RESULT FOR",
		specialStyle.Italic(true).Render(st),
	)
	l.Styles.Title = highlightStyle
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
	return tea.Batch(sl.spinner.Tick, sl.performSearch)
}

func (sl searchList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sl.list.SetSize(msg.Width-20, msg.Height-2)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return sl, tea.Quit
		case key.Matches(msg, quitBinding):
			return sl, tea.Quit
		case key.Matches(msg, backBinding):
			search := initialiseSearch(sl.queries, sl.searchTerm)
			return search, search.Init()
		}
	}
	sl.list, cmd = sl.list.Update(msg)
	cmds = append(cmds, cmd)

	sl.spinner, cmd = sl.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return sl, tea.Batch(cmds...)
}

func (sl searchList) View() string {
	sb := strings.Builder{}

	sb.WriteString(sl.list.View())
	return baseStyle.Render(sb.String())
}

func (sl searchList) performSearch() tea.Msg {

	return success{}
}
