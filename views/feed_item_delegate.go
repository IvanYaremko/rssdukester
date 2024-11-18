package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func feedItemDelegate() list.DefaultDelegate {

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
			}
		}
		return nil
	}

	return d
}
