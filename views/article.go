package views

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

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
	queries    *database.Queries
	markdown   string
	rss        item
	post       item
	spinner    spinner.Model
	loading    bool
	viewport   viewport.Model
	ready      bool
	help       help.Model
	keys       articleKeyMap
	navToSaved bool
	saved      bool
}

func InitialiseArticle(q *database.Queries, rss, selected item, fromSaved bool) article {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	vp := viewport.New(width, height)

	k := articleKeyMap{
		Up:   upBinding,
		Down: downBinding,
		Back: backBinding,
		Help: helpBinding,
		Quit: quitBinding,
	}

	if !fromSaved {
		k.Save = saveBinding
	}

	return article{
		queries:    q,
		rss:        rss,
		post:       selected,
		markdown:   "",
		loading:    true,
		spinner:    s,
		viewport:   vp,
		ready:      false,
		navToSaved: fromSaved,
		help:       help.New(),
		keys:       k,
		saved:      false,
	}
}

func (a article) Init() tea.Cmd {
	return tea.Batch(a.loadMarkdown(a.post.url), a.checkIfSaved(a.post.url), a.spinner.Tick)
}

func (a article) loadMarkdown(url string) tea.Cmd {
	return func() tea.Msg {
		post, err := a.queries.GetPostByUrl(context.Background(), url)
		if err == nil && post.Content.Valid {
			return successContent{content: post.Content.String}
		}

		markdown, err := reader.GetMarkdown(url)
		if err != nil {
			return failError{
				error: err,
			}
		}

		err = a.queries.CreatePost(context.Background(), database.CreatePostParams{
			FeedID:      a.rss.feedId,
			Title:       a.post.title,
			Url:         url,
			Content:     sql.NullString{String: markdown, Valid: true},
			PublishedAt: time.Now(), // TODO: bring in actual post through list.Item
			LastViewed:  time.Now(),
		})
		if err != nil && !strings.Contains(err.Error(), "UNIQUE constraint") {
			return failError{error: err}
		}

		return successContent{content: markdown}
	}
}

func (a article) checkIfSaved(url string) tea.Cmd {
	return func() tea.Msg {
		saved, err := a.queries.IsPostSaved(context.Background(), url)
		if err != nil {
			return failError{}
		}
		if saved == 1 {
			return success{}
		}
		return item{}
	}
}

func prettifyMarkdown(title, markdown string, w, h int) string {
	s := strings.Builder{}

	wrappedTitle := lipgloss.
		NewStyle().
		Italic(true).
		Width(w).
		PaddingTop(1).
		Foreground(highlight).
		Render(title)
	s.WriteString(wrappedTitle + "\n\n")

	wrappedContent := lipgloss.
		NewStyle().
		Width(w).
		Height(h).
		Foreground(text).
		Render(markdown)

	s.WriteString(wrappedContent)
	return s.String()
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
			if a.navToSaved {
				saved := initialiseSaved(a.queries)
				return saved, saved.Init()
			}
			feed := initialiseFeed(a.queries, a.rss)
			return feed, feed.Init()
		case key.Matches(msg, saveBinding):
			if a.navToSaved {
				return a, nil
			}
			return a, savePostItem(a.queries, a.post)
		}

	case successContent:
		a.loading = false
		a.markdown = msg.content
		prettyContent := prettifyMarkdown(a.post.title, msg.content, width, height)
		a.viewport.SetContent(prettyContent)
		return a, nil

	case success:
		a.saved = true
		a.keys = articleKeyMap{
			Up:   upBinding,
			Down: downBinding,
			Back: backBinding,
			Help: helpBinding,
			Quit: quitBinding,
		}
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
	sb := strings.Builder{}
	if a.saved {
		sb.WriteString(
			titleStyle.Render(
				fmt.Sprintf("%s %s",
					attentionStyle.Bold(true).Render("SAVED"),
					highlightStyle.Render(a.rss.title),
				),
			),
		)
	} else {
		sb.WriteString(titleStyle.Render(a.rss.title))
	}
	line := strings.Repeat("─", max(0, a.viewport.Width-lipgloss.Width(sb.String())))
	return lipgloss.JoinHorizontal(lipgloss.Center, sb.String(), line)
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

func savePostItem(q *database.Queries, selected item) tea.Cmd {
	return func() tea.Msg {
		post, err := q.GetPostByUrl(context.Background(), selected.url)
		if err != nil {
			return failError{error: err}
		}

		params := database.CreateSavedPostParams{
			PostID:    post.ID,
			CreatedAt: time.Now(),
		}
		err = q.CreateSavedPost(context.Background(), params)
		if err != nil {
			return fail{}
		}
		return success{}
	}
}
