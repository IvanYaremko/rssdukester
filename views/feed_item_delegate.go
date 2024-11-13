package views

import (
	"github.com/IvanYaremko/rssdukester/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func feedItemDelegate() list.DefaultDelegate {
	highlight := styles.HighlightStyle.
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Highlight).
		Padding(0, 0, 0, 1)

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = highlight
	d.Styles.SelectedDesc = highlight
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
