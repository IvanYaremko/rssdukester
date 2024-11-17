package views

import (
	"context"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type rssList struct {
	queries *database.Queries
	list    list.Model
}

func initialiseRssList(q *database.Queries) rssList {
	delegate := rssItemDelegate(q)

	items := make([]list.Item, 0)

	feedsList := list.New(items, delegate, 100, 40)
	feedsList.Title = "RSS FEEDS"
	feedsList.Styles.Title = highlightStyle
	feedsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			enterBinding,
			removeBinding,
		}
	}
	feedsList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			backBinding,
		}
	}
	l := rssList{
		queries: q,
		list:    feedsList,
	}
	return l
}

func (l rssList) Init() tea.Cmd {
	return l.getRssFeeds
}

func (l rssList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.list.SetSize(msg.Width-20, msg.Height-2)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return l, tea.Quit
		case key.Matches(msg, escBinding):
			if l.list.FilterState() == list.Filtering {
				l.list.ResetFilter()
				return l, nil
			}
			return l, nil
		case key.Matches(msg, backBinding):
			home := InitHomeModel(l.queries)
			return home, home.Init()
		}

	case successItems:
		cmd := l.list.SetItems(msg.items)
		return l, cmd

	case selected:
		feed := initialiseFeed(l.queries, msg.selected)
		return feed, feed.Init()
	}

	newList, cmd := l.list.Update(msg)
	cmds = append(cmds, cmd)
	l.list = newList

	return l, tea.Batch(cmds...)
}

func (l rssList) View() string {
	return baseStyle.Render(l.list.View())
}

func (l *rssList) getRssFeeds() tea.Msg {
	feeds, err := l.queries.GetFeeds(context.Background())
	if err != nil {
		return failError{error: err}
	}

	items := make([]list.Item, len(feeds))
	for i := range feeds {
		items[i] = item{
			title:       feeds[i].Name,
			url:         feeds[i].Url,
			description: feeds[i].Url,
		}
	}

	return successItems{
		items: items,
	}
}
