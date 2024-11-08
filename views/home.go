package views

import (
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var homeNav = []string{
	"View",
	"Add",
	"Quit",
}

type Home struct {
	queries *database.Queries
	keys    keyMap
	help    help.Model
	cursor  int
}

func (h Home) Init() tea.Cmd {

	return tea.SetWindowTitle("RSSDUKESTER")
}

func (h Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.help.Width = msg.Width

	case tea.KeyMsg:
		{
			switch msg.String() {
			case "ctrl+c":
				return h, tea.Quit
			}
			switch {
			case key.Matches(msg, h.keys.Up):
				if h.cursor == 0 {
					h.cursor = len(homeNav) - 1
				} else {
					h.cursor--
				}
			case key.Matches(msg, h.keys.Down):
				if h.cursor == len(homeNav)-1 {
					h.cursor = 0
				} else {
					h.cursor++
				}
			}
		}
	}

	return h, tea.Batch(cmds...)
}

func (h Home) View() string {
	s := strings.Builder{}
	for i := 0; i < len(homeNav); i++ {
		if h.cursor == i {
			s.WriteString(highlightStyle.Render("[•] ", homeNav[i]))
		} else {
			s.WriteString("[ ] " + homeNav[i])
		}
		s.WriteString("\n")
		s.WriteString("\n")
	}
	s.WriteString("\n")
	s.WriteString(h.help.View(h.keys))
	return baseStyle.Render(s.String())
}

func InitHomeModel(q *database.Queries) Home {
	return Home{
		queries: q,
		keys:    keys,
		help:    help.New(),
		cursor:  0,
	}
}