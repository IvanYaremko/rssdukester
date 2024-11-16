package views

import (
	"fmt"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/IvanYaremko/rssdukester/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/glamour"
)

type article struct {
	queries      *database.Queries
	content      string
	contentTitle string
	glammed      string
	backItem     item
	spinner      spinner.Model
	loading      bool
	viewport     viewport.Model
}

func InitialiseArticle(q *database.Queries, c, u string, i item) article {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.HighlightStyle

	vp := viewport.New(100, 50)

	return article{
		queries:      q,
		content:      c,
		contentTitle: u,
		backItem:     i,
		glammed:      "",
		loading:      true,
		spinner:      s,
		viewport:     vp,
	}
}

func (a article) Init() tea.Cmd {
	return tea.Batch(a.glamMarkdown, a.spinner.Tick)
}

func (a article) glamMarkdown() tea.Msg {
	read, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(100))

	output, err := read.Render(a.content)
	if err != nil {
		return failError{error: err}
	}

	return successContent{content: output}
}

func (a article) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return a, tea.Quit
		case key.Matches(msg, backBinding):
			feed := initialiseFeed(a.queries, a.backItem)
			return feed, feed.Init()
		}

	case successContent:
		a.loading = false
		a.glammed = msg.content
		a.viewport.SetContent(msg.content)
		return a, nil
	}

	a.viewport, cmd = a.viewport.Update(msg)
	cmds = append(cmds, cmd)

	a.spinner, cmd = a.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return a, tea.Batch(cmds...)
}

func (a article) View() string {
	if a.loading {
		return fmt.Sprintf("%s loading %s...", a.spinner.View(), a.contentTitle)
	}

	return a.viewport.View()
}
