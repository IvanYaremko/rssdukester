package views

import (
	"context"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func rssItemDelegate(q *database.Queries) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = itemNormalTitle
	d.Styles.NormalDesc = itemNormalDesc
	d.Styles.SelectedTitle = itemSelectedTitle
	d.Styles.SelectedDesc = itemSelectedDesc

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		selectedItem := m.SelectedItem()
		if selectedItem == nil {
			return nil
		}

		item, ok := selectedItem.(item)
		if !ok {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, enterBinding):
				return transitionView(item)

			case key.Matches(msg, removeBinding):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					removeBinding.SetEnabled(false)
				}
				return tea.Batch(removeFeed(item, q),
					m.NewStatusMessage("You deleted "+item.title))
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
