package rss

import (
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var choices = []string{"Add", "Update", "View"}

type Feeds struct {
	dbQueries *database.Queries
	cursor    int
	choice    string
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
			switch f.choice {
			case "Add":
				return InitialiseAddFeedModel(f.dbQueries), nil
			case "Update":
				return f, nil
			case "View":
				return f, nil
			}

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
	PaddingLeft(10)

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Underline(true).
	Italic(true).
	Render("F E E D S\n\n")

func (f Feeds) View() string {
	s := strings.Builder{}
	s.WriteString("F E E D S\n\n")

	for i := 0; i < len(choices); i++ {
		if f.cursor == i {
			s.WriteString("[•] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	styledContent := s.String()
	return styledContent
}

func InitialiseFeedsModel(queries *database.Queries) Feeds {
	return Feeds{
		dbQueries: queries,
		cursor:    0,
		choice:    "",
	}
}
