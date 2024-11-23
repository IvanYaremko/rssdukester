package views

import (
	"context"
	"time"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func feedItemDelegate(q *database.Queries, rss item) list.DefaultDelegate {
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

			case key.Matches(msg, saveBinding):
				return savePostItem(q, item, rss.title)
			}
		}
		return nil
	}

	return d
}

func savePostItem(q *database.Queries, selected item, feed string) tea.Cmd {
	return func() tea.Msg {
		params := database.SavePostParams{
			Url:       selected.url,
			Title:     selected.title,
			Feed:      feed,
			CreatedAt: time.Now(),
		}
		err := q.SavePost(context.Background(), params)
		if err != nil {
			return fail{}
		}
		return success{}
	}
}
