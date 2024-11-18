package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func feedItemDelegate() list.DefaultDelegate {

	widthStyle := lipgloss.NewStyle().Width(80).Foreground(text)
	highlight := highlightStyle.
		Width(80).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(highlight).
		Padding(0, 0, 0, 1)

	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = widthStyle
	d.Styles.NormalDesc = widthStyle
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
