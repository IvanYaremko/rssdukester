package views

import (
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	homeNav = []string{
		"feeds",
		"add",
		"saved",
		"search",
		"recommended",
	}
	width  = 100
	height = 40
)

type Home struct {
	queries *database.Queries
	keys    homeKeyMap
	help    help.Model
	cursor  int
}

func InitHomeModel(q *database.Queries) Home {
	return Home{
		queries: q,
		keys: homeKeyMap{
			Up:       upBinding,
			Down:     downBinding,
			Tab:      tabBinding,
			ShiftTab: shiftTabBinding,
			Enter:    enterBinding,
			Add:      addViewBiding,
			Help:     helpBinding,
			Quit:     quitBinding,
			Ctrlc:    ctrlcBinding,
		},
		help:   help.New(),
		cursor: 0,
	}
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
					rssFeeds := initialiseRssList(h.queries)
					return rssFeeds, rssFeeds.Init()
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
	highlightItalic := highlightStyle.Italic(true)
	s.WriteString(highlightItalic.Render("R S S D U K E S T E R"))
	s.WriteString("\n\n\n")
	for i := 0; i < len(homeNav); i++ {
		if h.cursor == i {
			s.WriteString(highlightStyle.Render("[•] ", strings.ToUpper(homeNav[i])))
		} else {
			s.WriteString("[ ] " + homeNav[i])
		}
		s.WriteString("\n")
		s.WriteString("\n")
	}
	s.WriteString("\n\n\n")
	s.WriteString(h.help.View(h.keys))
	return baseStyle.Render(s.String())
}
