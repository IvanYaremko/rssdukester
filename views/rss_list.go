package views

import (
	"context"
	"fmt"

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

	feedsList := list.New(items, delegate, width, height)
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
			ctrlcBinding,
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
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

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
		cmd = l.list.SetItems(msg.items)
		return l, cmd

	case selected:
		feed := initialiseFeed(l.queries, msg.selected)
		return feed, feed.Init()

	case success:
		if l.list.SelectedItem() == nil {
			message := errorStyle.Bold(true).Render("DELETED LAST FEED")
			cmd := l.list.NewStatusMessage(message)
			return l, cmd
		}
		selected := l.list.SelectedItem().(item)
		message := fmt.Sprintf("%s %s",
			errorStyle.Bold(true).Render("DELETED"),
			specialStyle.Italic(true).Render(selected.title),
		)
		cmd = l.list.NewStatusMessage(message)
		return l, cmd
	}

	l.list, cmd = l.list.Update(msg)
	cmds = append(cmds, cmd)

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
		hyperlink := fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", feeds[i].Url, feeds[i].Url)
		items[i] = item{
			title:       feeds[i].Name,
			description: specialStyle.Render(hyperlink),
			url:         feeds[i].Url,
			feedName:    feeds[i].Name,
			feedId:      feeds[i].ID,
		}
	}

	return successItems{
		items: items,
	}
}
