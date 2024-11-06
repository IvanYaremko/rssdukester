package rss

import (
	"context"
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	name, url string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.url }
func (i item) FilterValue() string { return i.name }

type ViewModel struct {
	dbQueries   *database.Queries
	list        list.Model
	err         error
	keys        *listKeyMap
	delgateKeys *delegateKeyMap
}

func (v ViewModel) Init() tea.Cmd {
	return tea.Batch(v.getFeeds, tea.EnterAltScreen)
}

func (v ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		x, y := docStyle.GetFrameSize()
		v.list.SetSize(msg.Width-x, msg.Height-y)

	case tea.KeyMsg:
		if v.list.FilterState() == list.Filtering {
			break
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return v, tea.Quit
		}

		switch msg.String() {
		case "r":
			return v, v.getFeeds
		}

		switch {
		case key.Matches(msg, v.keys.toggleTitleBar):
			flag := !v.list.ShowTitle()
			v.list.SetShowTitle(flag)
			return v, nil

		case key.Matches(msg, v.keys.toggleSpinner):
			cmd := v.list.ToggleSpinner()
			return v, cmd
		}

	case dbError:
		v.err = msg.dbErr
		return v, nil

	case dbSuccess:
		cmd := v.list.SetItems(msg.items)
		return v, cmd

	}

	newListModel, cmd := v.list.Update(msg)
	v.list = newListModel
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

func (v ViewModel) View() string {
	builder := strings.Builder{}
	builder.WriteString("VIEW feeds\n")
	if v.err != nil {
		builder.WriteString("error gettings feeds from db")
		builder.WriteString("\n")
		builder.WriteString(v.err.Error())
		return builder.String()
	}

	if len(v.list.Items()) == 0 {
		builder.WriteString("no feeds found in db")
		return builder.String()
	}

	return docStyle.Render(v.list.View())
}

func InitialiseViewModel(queries *database.Queries) ViewModel {
	var (
		listKeys     = newListKeyMap()
		delegateKeys = newDelegateKeyMap()
	)

	delagate := newItemDelegate(delegateKeys)

	l := list.New(make([]list.Item, 0), delagate, 30, 30)
	l.SetShowTitle(true)
	l.Title = "RSS FEEDS"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			delegateKeys.choose,
			delegateKeys.remove,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleHelpMenu,
			listKeys.togglePagination,
			listKeys.toggleStatusBar,
			listKeys.insertItem,
		}
	}

	return ViewModel{
		dbQueries:   queries,
		err:         nil,
		list:        l,
		delgateKeys: delegateKeys,
		keys:        listKeys,
	}
}

func (v *ViewModel) getFeeds() tea.Msg {
	feeds, err := v.dbQueries.GetFeeds(context.Background())
	if err != nil {
		return dbError{dbErr: err}
	}

	list := make([]list.Item, 0, len(feeds))
	for _, feed := range feeds {
		list = append(list, item{
			name: feed.Name,
			url:  feed.Url,
		})
	}

	return dbSuccess{items: list}
}

type dbSuccess struct {
	items []list.Item
}
