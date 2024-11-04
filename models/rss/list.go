package rss

import (
	"context"
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	name, url string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.url }
func (i item) FilterValue() string { return i.name }

type ViewModel struct {
	dbQueries *database.Queries
	list      list.Model
	err       error
}

func (v ViewModel) Init() tea.Cmd {
	return tea.Batch(v.getFeeds, tea.EnterAltScreen)
}

func (v ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return v, tea.Quit
		}
		switch msg.String() {
		case "r":
			return v, v.getFeeds
		}

	case dbError:
		v.err = msg.dbErr
		return v, nil

	case dbSuccess:
		cmd := v.list.SetItems(msg.items)
		return v, cmd
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
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
	l := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(true)
	l.Title = "RSS FEEDS"
	l.SetSize(30, 30)

	return ViewModel{
		dbQueries: queries,
		err:       nil,
		list:      l,
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
