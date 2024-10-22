package rss

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var choices = []string{"Taro", "Coffee", "Lychee"}

type Feeds struct {
	cursor int
	choice string
}

func (f Feeds) Init() tea.Cmd {
	// get feeds from database?
	return nil
}

func (f Feeds) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return f, tea.Quit

		case "enter":
			f.choice = choices[f.cursor]
			// return some other Model to change view?
			return f, tea.Quit

		case "down", "j":
			f.cursor++
			if f.cursor >= len(choices) {
				f.cursor = 0
			}

		case "up", "k":
			f.cursor--
			if f.cursor < 0 {
				f.cursor = len(choices) - 1
			}
		}
	}

	return f, nil
}

var borderStyle = lipgloss.NewStyle().
	Height(50).
	Width(100).
	PaddingLeft(10).
	PaddingTop(2)

func (f Feeds) View() string {
	s := strings.Builder{}
	s.WriteString("What kind of Bubble Tea would you like to order?\n\n")

	for i := 0; i < len(choices); i++ {
		if f.cursor == i {
			s.WriteString("[•] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	styledContent := borderStyle.Render(s.String())
	return styledContent
}

func InitialiseFeedsModel() Feeds {
	return Feeds{
		cursor: 0,
		choice: "",
	}
}
