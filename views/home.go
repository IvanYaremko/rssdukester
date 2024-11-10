package views

import (
	"strings"

	"github.com/IvanYaremko/rssdukester/bindings"
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/IvanYaremko/rssdukester/styles"
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
	keys    bindings.HomeKeyMap
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
			switch {
			case key.Matches(msg, h.keys.Ctrlc):
				return h, tea.Quit
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
			case key.Matches(msg, h.keys.Quit):
				return h, tea.Quit
			case key.Matches(msg, h.keys.Help):
				h.help.ShowAll = !h.help.ShowAll
			case key.Matches(msg, h.keys.Enter):
				switch h.cursor {
				case 0:
					listFeeds := initialiseListFeeds(h.queries)
					return listFeeds, listFeeds.Init()
				case 1:
					//
					addFeed := initialiseAddFeed(h.queries)
					return addFeed, addFeed.Init()
				case 2:
					//
					return h, nil
				}
			}
		}
	}

	return h, tea.Batch(cmds...)
}

func (h Home) View() string {
	s := strings.Builder{}
	highlightItalic := styles.HighlightStyle.Italic(true)
	s.WriteString(highlightItalic.Render("R S S D U K E S T E R"))
	s.WriteString("\n\n\n")
	for i := 0; i < len(homeNav); i++ {
		if h.cursor == i {
			s.WriteString(styles.HighlightStyle.Render("[•] ", homeNav[i]))
		} else {
			s.WriteString("[ ] " + homeNav[i])
		}
		s.WriteString("\n")
		s.WriteString("\n")
	}
	s.WriteString("\n\n\n")
	s.WriteString(h.help.View(h.keys))
	return styles.BaseStyle.Render(s.String())
}

func InitHomeModel(q *database.Queries) Home {
	return Home{
		queries: q,
		keys:    bindings.HomeKeys,
		help:    help.New(),
		cursor:  0,
	}
}
