package views

import (
	"fmt"
	"strings"

	"github.com/IvanYaremko/rssdukester/reader"
	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	width      = 80
	height     = 40
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).Foreground(highlight)
	}()
	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type article struct {
	queries  *database.Queries
	markdown string
	rss      item
	post     item
	spinner  spinner.Model
	loading  bool
	viewport viewport.Model
	ready    bool
	help     help.Model
	keys     articleKeyMap
}

func InitialiseArticle(q *database.Queries, rss, selected item) article {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	vp := viewport.New(width, height)

	return article{
		queries:  q,
		rss:      rss,
		post:     selected,
		markdown: "",
		loading:  true,
		spinner:  s,
		viewport: vp,
		ready:    false,
		help:     help.New(),
		keys: articleKeyMap{
			Up:   upBinding,
			Down: downBinding,
			Back: backBinding,
			Save: saveBinding,
			Help: helpBinding,
			Quit: quitBinding,
		},
	}
}

func (a article) Init() tea.Cmd {
	return tea.Batch(loadMarkdown(a.post.url), a.spinner.Tick)
}

func loadMarkdown(url string) tea.Cmd {
	return func() tea.Msg {
		markdown, err := reader.GetMarkdown(url)
		if err != nil {
			return failError{
				error: err,
			}
		}

		return successContent{content: markdown}
	}
}

func prettifyMarkdown(title, markdown string, w, h int) string {
	wrappedTitle := lipgloss.
		NewStyle().
		Italic(true).
		Width(w).
		PaddingTop(1).
		Foreground(highlight).
		Render(title)

	wrappedContent := lipgloss.
		NewStyle().
		Width(w).
		Height(h).
		Foreground(text).
		Render(markdown)

	return wrappedTitle + "\n\n" + wrappedContent
}

func (a article) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.viewport.Height = msg.Height - 20
		a.viewport.Width = msg.Width - 40
		a.viewport.SetContent(prettifyMarkdown(a.post.title, a.markdown, msg.Width-41, msg.Height-21))

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, quitBinding):
			return a, tea.Quit
		case key.Matches(msg, ctrlcBinding):
			return a, tea.Quit
		case key.Matches(msg, backBinding):
			feed := initialiseFeed(a.queries, a.rss)
			return feed, feed.Init()
		}

	case successContent:
		a.loading = false
		a.markdown = msg.content
		prettyContent := prettifyMarkdown(a.post.title, msg.content, width, height)
		a.viewport.SetContent(prettyContent)
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
		return baseStyle.Render(
			fmt.Sprintf("%s loading...\n%s", a.spinner.View(), a.post.title),
		)
	}

	s := strings.Builder{}

	s.WriteString(a.headerView() + "\n")
	s.WriteString(a.viewport.View() + "\n")
	s.WriteString(a.footerView() + "\n\n")
	s.WriteString(a.help.View(a.keys) + "\n")

	return baseStyle.Render(s.String())
}

func (a article) headerView() string {
	title := titleStyle.Render(a.rss.title)
	line := strings.Repeat("─", max(0, a.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (a article) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", a.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, a.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
