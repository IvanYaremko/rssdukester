package views

import (
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type search struct {
	quries      *database.Queries
	spinner     spinner.Model
	loading     bool
	list        list.Model
	searchInput textinput.Model
	err         error
	searchTerm  string
}

func (s search) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return s, nil
}

func (s search) View() string {

	return ""
}
