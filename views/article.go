package views

import (
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/glamour"
)

type article struct {
	queries *database.Queries
	content string
	glammed string
	item    item
}

func InitialiseArticle(q *database.Queries, c string, i item) article {
	return article{
		queries: q,
		content: c,
		item:    i,
		glammed: "",
	}
}

func (a article) Init() tea.Cmd {
	return a.glamMarkdown
}

func (a article) glamMarkdown() tea.Msg {
	read, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())

	output, err := read.Render(a.content)
	if err != nil {
		return failError{error: err}
	}

	return successContent{content: output}
}

func (a article) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return a, tea.Quit
		case key.Matches(msg, backBinding):
			feed := initialiseFeed(a.queries, a.item)
			return feed, feed.Init()
		}

	case successContent:
		a.glammed = msg.content
		return a, nil
	}

	return a, nil
}

func (a article) View() string {
	if a.glammed == "" {
		return ""
	}
	return a.glammed
}
