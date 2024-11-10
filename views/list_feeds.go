package views

import (
	"context"
	"time"

	"github.com/IvanYaremko/rssdukester/bindings"
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/IvanYaremko/rssdukester/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var base = styles.BaseStyle

type item struct {
	name      string
	url       string
	updatedAt time.Time
}

func (i item) FilterValue() string { return i.name }
func (i item) Description() string { return i.url }
func (i item) Title() string       { return i.name }

type listFeeds struct {
	queries *database.Queries
	list    list.Model
}

func initialiseListFeeds(q *database.Queries) listFeeds {
	var (
		delegateKeys = bindings.NewListItemDelegateKeyMap()
	)

	delegate := listItemDelegate(delegateKeys, q)

	items := make([]list.Item, 0)

	feedsList := list.New(items, delegate, 30, 30)
	feedsList.Title = "rss feeds"
	feedsList.Styles.Title = styles.HighlightStyle
	feedsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			delegateKeys.Choose,
			delegateKeys.Remove,
		}
	}
	l := listFeeds{
		queries: q,
		list:    feedsList,
	}
	return l
}

func (l listFeeds) Init() tea.Cmd {
	return l.getFeeds
}

func (l listFeeds) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		x, y := base.GetFrameSize()
		l.list.SetSize(msg.Width-x, msg.Height-y)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return l, tea.Quit
		case tea.KeyEsc:
			home := InitHomeModel(l.queries)
			return home, home.Init()
		}

	case successItems:
		cmd := l.list.SetItems(msg.items)
		return l, cmd
	}

	newList, cmd := l.list.Update(msg)
	cmds = append(cmds, cmd)
	l.list = newList

	return l, tea.Batch(cmds...)
}

func (l listFeeds) View() string {
	return base.Render(l.list.View())
}

func (l *listFeeds) getFeeds() tea.Msg {
	feeds, err := l.queries.GetFeeds(context.Background())
	if err != nil {
		return failError{error: err}
	}

	items := make([]list.Item, len(feeds))
	for i := range feeds {
		items[i] = item{
			name:      feeds[i].Name,
			url:       feeds[i].Url,
			updatedAt: feeds[i].UpdatedAt,
		}
	}

	return successItems{
		items: items,
	}
}
