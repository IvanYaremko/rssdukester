package views

import (
	"context"

	"github.com/IvanYaremko/rssdukester/bindings"
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func listItemDelegate(keys bindings.ListItemDelegateKeyMap, q *database.Queries) list.DefaultDelegate {
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
			case key.Matches(msg, keys.Choose):
				return m.NewStatusMessage("You chose " + title)

			case key.Matches(msg, keys.Remove):
				index := m.Index()
				item := m.SelectedItem().(item)
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.Remove.SetEnabled(false)
				}
				return tea.Batch(removeFeed(item, q),
					m.NewStatusMessage("You deleted "+title))
			}
		}
		return nil
	}

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
