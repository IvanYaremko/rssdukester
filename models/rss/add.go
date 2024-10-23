package rss

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e03434"))
)

type AddFeed struct {
	inputs []textinput.Model
	cursor int
	err    error
	errors []string
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
				// form submit
				return a, nil
			} else {
				a.nextInput()
			}
		}
	// catch errors from a CMD, maybe db error access error
	case errMsg:
		a.err = msg
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

	return builder.String()

	// for i := 0; i < len(a.errors); i++ {
	// 	if a.errors[i] != "" {
	// 		builder.WriteString(a.errors[i])
	// 	}
	// }
	//
	// s := fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s\n\n%s\n\n\ncursor: %v",
	// 	inputStyle.Render("Feed name"),
	// 	a.inputs[0].View(),
	// 	inputStyle.Render("Feed URL"),
	// 	a.inputs[1].View(),
	// 	builder.String(),
	// 	helpStyle.Render("tab - navigate\nenter - submit\nq - quit"),
	// 	a.cursor,
	// )
}

func textValidate(s string) error {
	if s == "" {
		return fmt.Errorf("can't be empty")
	}
	return nil
}

func InitialiseAddFeedModel() AddFeed {
	inputs := make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "hackernews"
	inputs[0].Validate = textValidate
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "https://hnrss.org/newest"
	inputs[1].Validate = textValidate

	return AddFeed{
		inputs: inputs,
		cursor: 0,
		err:    nil,
		errors: make([]string, len(inputs)),
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
