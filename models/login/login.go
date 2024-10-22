package login

import (
	"fmt"

	"github.com/IvanYaremko/rssdukester/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	textInput textinput.Model
	err       error
}

func (f LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (f LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return f, tea.Quit
		}
	case error:
		f.err = msg
		return f, nil
	}

	username := f.textInput.Value()
	cfg := config.Config{}
	cfg.SetUser(username)

	f.textInput, cmd = f.textInput.Update(msg)
	return f, cmd
}

func (f LoginModel) View() string {
	return fmt.Sprintf("Enter your username:\n\n%s\n", f.textInput.View()) + "\n"
}

func InitialLoginModel() LoginModel {
	ti := textinput.New()
	ti.Placeholder = "Duke"
	ti.Focus()
	ti.CharLimit = 124
	ti.Width = 20

	return LoginModel{
		textInput: ti,
		err:       nil,
	}
}
