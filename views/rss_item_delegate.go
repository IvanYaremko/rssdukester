package views

import (
	"context"

	"github.com/IvanYaremko/rssdukester/bindings"
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/IvanYaremko/rssdukester/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func rssItemDelegate(q *database.Queries) list.DefaultDelegate {
	highlight := styles.HighlightStyle.
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Highlight).
		Padding(0, 0, 0, 1)

	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.name
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, bindings.ListItemDelegateKeys.Choose):
				item := m.SelectedItem().(item)
				return transitionView(item)

			case key.Matches(msg, bindings.ListItemDelegateKeys.Remove):
				index := m.Index()
				item := m.SelectedItem().(item)
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					bindings.ListItemDelegateKeys.Remove.SetEnabled(false)
				}
				return tea.Batch(removeFeed(item, q),
					m.NewStatusMessage("You deleted "+title))
			}
		}
		return nil
	}

	d.Styles.SelectedTitle = highlight
	d.Styles.SelectedDesc = highlight

	return d
}

func removeFeed(item item, q *database.Queries) tea.Cmd {
	return func() tea.Msg {
		err := q.DeleteFeed(context.Background(), item.url)
		if err != nil {
			return failError{error: err}
		}
		return success{}
	}
}

type selectedFeed struct {
	rssFeed item
}

func transitionView(i item) tea.Cmd {
	return func() tea.Msg {
		return selectedFeed{rssFeed: i}
	}
}
