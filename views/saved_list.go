package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type saved struct {
	queries *database.Queries
	spinner spinner.Model
	loading bool
	list    list.Model
	err     error
}

func initialiseSaved(q *database.Queries) saved {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, savedItemDelegate(q), width, height)
	l.Title = "SAVED POSTS"
	l.Styles.Title = highlightStyle

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			enterBinding,
			backBinding,
			removeBinding,
		}
	}

	l.StatusMessageLifetime = 3 * time.Second

	return saved{
		queries: q,
		spinner: s,
		list:    l,
		loading: true,
		err:     nil,
	}
}

func (s saved) Init() tea.Cmd {
	return tea.Batch(s.spinner.Tick, s.getSavedPosts)
}

func (s *saved) getSavedPosts() tea.Msg {
	saved, err := s.queries.GetSavedPosts(context.Background())
	if err != nil {
		return failError{error: err}
	}

	items := make([]list.Item, 0, len(saved))
	for _, post := range saved {
		hyperlink := fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", post.Url, "Article link →")

		t, _ := parseDateTime(post.CreatedAt.String())
		savedDate := subtleStyle.Italic(true).Render(t.Format(formateDate))

		desc := fmt.Sprintf("%s %s %s",
			specialStyle.Render(hyperlink),
			attentionStyle.Render(post.Feed),
			savedDate,
		)

		items = append(items, item{
			title:       post.Title,
			description: desc,
			url:         post.Url,
			feed:        post.Feed,
		})
	}
	return successItems{
		items: items,
	}
}

func (s saved) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width-20, msg.Height-2)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return s, tea.Quit
		case key.Matches(msg, quitBinding):
			return s, tea.Quit
		case key.Matches(msg, backBinding):
			home := InitHomeModel(s.queries)
			return home, home.Init()
		}
	case successItems:
		cmd = s.list.SetItems(msg.items)
		s.loading = false
		return s, cmd

	case failError:
		s.err = msg.error
		s.loading = false
		return s, nil

	case success:
		if s.list.SelectedItem() == nil {
			message := errorStyle.Bold(true).Render("DELETED LAST FEED")
			cmd := s.list.NewStatusMessage(message)
			return s, cmd
		}
		selected := s.list.SelectedItem().(item)
		message := fmt.Sprintf("%s %s",
			errorStyle.Bold(true).Render("DELETED"),
			specialStyle.Italic(true).Render(selected.title),
		)
		cmd = s.list.NewStatusMessage(message)
		return s, cmd

	case selected:
		article := InitialiseArticle(s.queries, item{title: msg.selected.feed}, msg.selected, true)
		return article, article.Init()
	}

	s.list, cmd = s.list.Update(msg)
	cmds = append(cmds, cmd)

	s.spinner, cmd = s.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s saved) View() string {
	sb := strings.Builder{}

	if s.err != nil {
		return fmt.Sprintf("%s\n%s",
			errorStyle.Render("something has gone wrong!"),
			errorStyle.Render(s.err.Error()),
		)
	}

	if s.loading {
		sb.WriteString(fmt.Sprintf("%s loading...", s.spinner.View()))
	} else {
		sb.WriteString(s.list.View())
	}

	return baseStyle.Render(sb.String())
}
