package views

import (
	"fmt"
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/IvanYaremko/rssdukester/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/glamour"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type article struct {
	queries      *database.Queries
	content      string
	contentTitle string
	glammed      string
	rss          item
	post         item
	spinner      spinner.Model
	loading      bool
	viewport     viewport.Model
	ready        bool
}

func InitialiseArticle(q *database.Queries, c string, rssItem, feedItem item) article {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.HighlightStyle

	vp := viewport.New(60, 30)
	vp.SetContent(c)
	//	vp.Style = lipgloss.NewStyle().
	//		Border(lipgloss.NormalBorder())

	return article{
		queries:  q,
		content:  c,
		rss:      rssItem,
		post:     feedItem,
		glammed:  "",
		loading:  true,
		spinner:  s,
		viewport: vp,
		ready:    false,
	}
}

func (a article) Init() tea.Cmd {
	return nil
	return tea.Batch(a.glamMarkdown, a.spinner.Tick)
}

func (a article) glamMarkdown() tea.Msg {
	read, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0), // Let viewport handle wrapping
		glamour.WithEmoji(),     // Handle emoji if present
	)
	output, err := read.Render(a.content)
	if err != nil {
		return failError{error: err}
	}
	return successContent{content: output}
}

func (a article) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(a.headerView())
		footerHeight := lipgloss.Height(a.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		//if !a.ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		a.viewport.Height = msg.Height - verticalMarginHeight
		a.viewport.Width = msg.Width
		//a.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		a.viewport.YPosition = headerHeight
		//	a.ready = true

		// This is only necessary for high performance rendering, which in
		// most cases you won't need.
		//
		// Render the viewport one line below the header.
		//	a.viewport.YPosition = headerHeight + 1
		//} else {
		//a.viewport.Width = msg.Width - 2
		//a.viewport.Height = msg.Height - verticalMarginHeight
		a.viewport.SetContent(a.glammed)
	//}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return a, tea.Quit
		case key.Matches(msg, backBinding):
			feed := initialiseFeed(a.queries, a.rss)
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
	return base.Render(fmt.Sprintf("%s\n%s\n%s", a.headerView(), a.viewport.View(), a.footerView()))
	if a.loading {
		return fmt.Sprintf("%s loading...\n%s", a.spinner.View(), a.post.title)
	}

	return base.Render(fmt.Sprintf("%s\n%s\n%s", a.headerView(), a.viewport.View(), a.footerView()))
}

func (a article) headerView() string {
	title := titleStyle.Render("Article")
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
