package views

import (
	"fmt"
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type search struct {
	queries     *database.Queries
	searchInput textinput.Model
	err         error
	searchTerm  string
	cursor      int
	help        help.Model
	keys        addKeyMap
}

func initialiseSearch(q *database.Queries, st string) search {
	ti := textinput.New()
	ti.Placeholder = "Enter search terms..."
	ti.TextStyle = highlightStyle
	ti.Focus()
	ti.Width = width

	if st != "" {
		ti.SetValue(st)
	}

	return search{
		queries:     q,
		searchInput: ti,
		err:         nil,
		searchTerm:  st,
		cursor:      0,
		help:        help.New(),
		keys: addKeyMap{
			Up:       strictUpBinding,
			Down:     strictDownBinding,
			Tab:      tabBinding,
			ShiftTab: shiftTabBinding,
			Enter:    enterBinding,
			Help:     helpBinding,
			Back:     escBinding,
			Ctrlc:    ctrlcBinding,
		},
	}
}

func (s search) Init() tea.Cmd {
	return nil
}

func (s search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	if s.cursor == 0 {
		s.searchInput, cmd = s.searchInput.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	moveFocus := []key.Binding{
		s.keys.Down,
		s.keys.Tab,
		s.keys.Up,
		s.keys.ShiftTab,
		s.keys.Enter,
	}
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		s.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Ctrlc):
			return s, tea.Quit
		case key.Matches(msg, s.keys.Back):
			home := InitHomeModel(s.queries)
			return home, home.Init()
		case key.Matches(msg, s.keys.Help):
			s.help.ShowAll = !s.help.ShowAll
			return s, nil
		}

		switch {
		case key.Matches(msg, moveFocus...):
			if key.Matches(msg, s.keys.Enter) && s.cursor == 1 {
				// transition view submit
			}
			if s.cursor == 0 {
				s.cursor = 1
				s.searchInput.Blur()
				s.searchInput.TextStyle = lipgloss.NewStyle()
			} else {
				s.cursor = 0
				cmd = s.searchInput.Focus()
				cmds = append(cmds, cmd)
				s.searchInput.TextStyle = highlightStyle.Bold(true)
			}
		}
		return s, tea.Batch(cmds...)
	}

	return s, tea.Batch(cmds...)
}

func (s search) View() string {
	sb := strings.Builder{}

	if s.err != nil {
		return baseStyle.Render(
			fmt.Sprintf("%s \n\n %s",
				errorStyle.Render("something has gone wrong"),
				errorStyle.Render(s.err.Error()),
			),
		)
	}

	sb.WriteString(s.searchInput.View())

	sb.WriteString("\n\n")
	button := fmt.Sprintf("[ %s ]", subtleStyle.Render("submit"))

	if s.cursor == 1 {
		button = highlightStyle.Render("[ Submit ]")
	}

	sb.WriteString(button)
	sb.WriteString("\n\n\n")
	sb.WriteString(s.help.View(s.keys))
	return baseStyle.Render(sb.String())
}
