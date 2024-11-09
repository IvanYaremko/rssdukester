package views

import (
	"fmt"
	"strings"

	"github.com/IvanYaremko/rssdukester/bindings"
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type addFeed struct {
	queries *database.Queries
	inputs  []textinput.Model
	cursor  int
	keyMap  bindings.AddKeyMap
	help    help.Model
}

func initialiseAddFeed(q *database.Queries) addFeed {
	a := addFeed{
		queries: q,
		inputs:  make([]textinput.Model, 2),
		keyMap:  bindings.AddKeys,
		help:    help.New(),
	}

	var t textinput.Model
	for i := range a.inputs {
		t = textinput.New()
		t.Cursor.Style = baseStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "hackernews"
			t.Focus()
			t.TextStyle = highlightStyle.Bold(false)
		case 1:
			t.Placeholder = "https://hnrss.org/frontpage"
		}
		a.inputs[i] = t
	}

	return a
}

func (a addFeed) Init() tea.Cmd {

	return nil
}

func (a addFeed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit

		case "up", "down", "tab", "shift+tab", "enter":

		}

	}

	cmd := a.updateInputs(msg)
	cmds = append(cmds, cmd)
	return a, tea.Batch(cmds...)
}

func (a addFeed) View() string {
	s := strings.Builder{}

	for i := range a.inputs {
		s.WriteString(a.inputs[i].View())
		if i < len(a.inputs)-1 {
			s.WriteString("\n\n")
		}
	}
	s.WriteString("\n\n")
	button := fmt.Sprintf("[ %s ]", subtleStyle.Render("submit"))

	if a.cursor == len(a.inputs) {
		button = highlightStyle.Render("[ Submit ]")
	}

	s.WriteString(button)
	s.WriteString("\n\n\n")
	s.WriteString(a.help.View(a.keyMap))
	return baseStyle.Render(s.String())
}

func (a *addFeed) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(a.inputs))

	for i := range a.inputs {
		a.inputs[i], cmds[i] = a.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
