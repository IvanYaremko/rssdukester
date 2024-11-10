package views

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/IvanYaremko/rssdukester/bindings"
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type addFeed struct {
	queries    *database.Queries
	inputs     []textinput.Model
	cursor     int
	keyMap     bindings.AddKeyMap
	help       help.Model
	inputError error
	entries    []string
	dbErr      error
}

func textValidate(s string) error {
	if s == "" {
		return errors.New("empty input(s)")
	}
	return nil
}

func initialiseAddFeed(q *database.Queries) addFeed {
	a := addFeed{
		queries:    q,
		inputs:     make([]textinput.Model, 2),
		keyMap:     bindings.AddKeys,
		help:       help.New(),
		inputError: nil,
		entries:    make([]string, 0),
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
		a.inputs[i].Validate = textValidate
	}

	return a
}

func (a addFeed) Init() tea.Cmd {

	return nil
}

func (a addFeed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	moveFocus := []key.Binding{
		a.keyMap.Down,
		a.keyMap.Tab,
		a.keyMap.Up,
		a.keyMap.ShiftTab,
		a.keyMap.Enter,
	}
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.help.Width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		}
		switch {
		case key.Matches(msg, a.keyMap.Back):
			home := InitHomeModel(a.queries)
			return home, home.Init()
		case key.Matches(msg, a.keyMap.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.keyMap.Help):
			a.help.ShowAll = !a.help.ShowAll
			return a, nil
		}

		switch {
		case key.Matches(msg, moveFocus...):
			if key.Matches(msg, a.keyMap.Enter) && a.cursor == len(a.inputs) {
				// check inputs valid
				for i := range a.inputs {
					input := a.inputs[i]
					err := input.Validate(input.Value())
					if err != nil {
						a.inputError = err
						a.cursor = 0
						cmd := a.inputs[0].Focus()
						cmds = append(cmds, cmd)
						return a, tea.Batch(cmds...)
					}
				}
				return a, a.createFeed
			}

			if key.Matches(msg, a.keyMap.Up) || key.Matches(msg, a.keyMap.ShiftTab) {
				a.cursor--
			} else {
				a.cursor++
			}

			if a.cursor > len(a.inputs) {
				a.cursor = 0
			} else if a.cursor < 0 {
				a.cursor = len(a.inputs)
			}

			for i := range a.inputs {
				if i == a.cursor {
					cmd := a.inputs[i].Focus()
					cmds = append(cmds, cmd)
					a.inputs[i].TextStyle = highlightStyle.Bold(true)
					continue
				}
				a.inputs[i].Blur()
				a.inputs[i].TextStyle = lipgloss.NewStyle()
			}
			return a, tea.Batch(cmds...)
		}

	case success:
		a.entries = append(a.entries,
			fmt.Sprintf("✓ added %s %s",
				highlightStyle.Bold(false).Render(a.inputs[0].Value()),
				highlightStyle.Bold(false).Render(a.inputs[1].Value())),
		)
		a.cursor = 0
		for i := range a.inputs {
			a.inputs[i].SetValue("")
		}
		cmd := a.inputs[0].Focus()
		cmds = append(cmds, cmd)
		a.inputError = nil
		a.dbErr = nil
		return a, tea.Batch(cmds...)

	case failError:
		a.dbErr = msg.err
		a.cursor = 0
		for i := range a.inputs {
			a.inputs[i].SetValue("")
		}
		cmd := a.inputs[0].Focus()
		cmds = append(cmds, cmd)
		a.inputError = nil
		return a, tea.Batch(cmds...)
	}

	cmd := a.updateInputs(msg)
	cmds = append(cmds, cmd)
	return a, tea.Batch(cmds...)
}

func (a addFeed) View() string {
	s := strings.Builder{}

	if a.inputError != nil {
		s.WriteString(errorStyle.Render(a.inputError.Error()))
		s.WriteString("\n\n\n")
	}

	if a.dbErr != nil {
		s.WriteString(errorStyle.Render("failed database insert"))
		s.WriteString("\n\n\n")
	}

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

	for i := range a.entries {
		s.WriteString(attentionStyle.Render(a.entries[i]))
		s.WriteString("\n")
	}

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

func (a *addFeed) createFeed() tea.Msg {
	args := database.CreateFeedParams{
		Name:      a.inputs[0].Value(),
		Url:       a.inputs[1].Value(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := a.queries.CreateFeed(context.Background(), args)
	if err != nil {
		return failError{err: err}
	}
	return success{}
}
