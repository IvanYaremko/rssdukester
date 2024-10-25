package rss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e03434"))
)

type AddFeed struct {
	dbQueris  *database.Queries
	inputs    []textinput.Model
	feed      database.Feed
	cursor    int
	err       error
	errors    []string
	isSuccess bool
}

func (a AddFeed) Init() tea.Cmd {
	return nil
}

func (a AddFeed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		case tea.KeyTab:
			a.nextInput()
		case tea.KeyEnter:
			tarInput := a.inputs[a.cursor]
			validationErr := tarInput.Validate(tarInput.Value())
			if validationErr != nil {
				a.errors[a.cursor] = validationErr.Error()
				return a, nil
			} else {
				a.errors[a.cursor] = ""
			}

			if a.cursor == len(a.inputs)-1 {
				return a, a.createFeed
			} else {
				a.nextInput()
			}
		}
		return a, nil
	case success:
		a.isSuccess = true
		return a, nil

	case dbError:
		a.isSuccess = false
		a.err = msg.dbErr
		return a, nil

	}

	for i := 0; i < len(a.inputs); i++ {
		var cmd tea.Cmd
		a.inputs[i], cmd = a.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return a, tea.Batch(cmds...)
}

func (a AddFeed) View() string {
	builder := strings.Builder{}
	builder.WriteString("FEED 0 \n\n")
	builder.WriteString(a.feed.Url)
	builder.WriteString("\n\n\n\n")
	builder.WriteString(inputStyle.Render("feed name"))
	builder.WriteString("\n")
	builder.WriteString(a.inputs[0].View())
	builder.WriteString("\n")
	if a.errors[0] != "" {
		builder.WriteString(errorStyle.Render(a.errors[0]))
		builder.WriteString("\n")
	}
	builder.WriteString(inputStyle.Render("feed url"))
	builder.WriteString("\n")
	builder.WriteString(a.inputs[1].View())
	builder.WriteString("\n")
	if a.errors[1] != "" {
		builder.WriteString(errorStyle.Render(a.errors[1]))
	}
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString(helpStyle.Render("tab (switch input)"))
	builder.WriteString("\n")
	builder.WriteString(helpStyle.Render("enter (submit)"))
	builder.WriteString("\n")
	builder.WriteString(helpStyle.Render("ctrl+c (quit)"))
	builder.WriteString("\n")
	if a.err != nil {
		builder.WriteString("ERROR INSERTING INTO DB\n")
		builder.WriteString(a.err.Error())
		builder.WriteString("\n")
	}

	if a.isSuccess {
		builder.WriteString("SUCCESS")
	}

	return builder.String()
}

func textValidate(s string) error {
	if s == "" {
		return fmt.Errorf("can't be empty")
	}
	return nil
}

func InitialiseAddFeedModel(queries *database.Queries) AddFeed {
	inputs := make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "hackernews"
	inputs[0].Validate = textValidate
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "https://hnrss.org/newest"
	inputs[1].Validate = textValidate

	return AddFeed{
		dbQueris: queries,
		inputs:   inputs,
		cursor:   0,
		feed: database.Feed{
			Name:      "test",
			Url:       "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		err:       nil,
		errors:    make([]string, len(inputs)),
		isSuccess: false,
	}
}

func (a *AddFeed) nextInput() {
	a.cursor = (a.cursor + 1) % len(a.inputs)
	for i := 0; i < len(a.inputs); i++ {
		if i == a.cursor {
			a.inputs[i].Focus()
		} else {
			a.inputs[i].Blur()
		}
	}
}

func (a *AddFeed) prevInput() {
	a.cursor--
	if a.cursor < 0 {
		a.cursor = len(a.inputs) - 1
	}
}

func (a AddFeed) createFeed() tea.Msg {
	name := a.inputs[0].Value()
	url := a.inputs[1].Value()
	params := database.CreateFeedParams{
		Name:      name,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := a.dbQueris.CreateFeed(context.Background(), params)

	if err != nil {
		return dbError{dbErr: err}
	}
	return success{}
}

type success struct{}

type dbError struct {
	dbErr error
}
