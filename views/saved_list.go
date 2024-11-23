package views

import (
	"context"
	"fmt"
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
}

func initialiseSaved(q *database.Queries) saved {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, list.NewDefaultDelegate(), width, height)
	l.Title = "SAVED POSTS"
	l.Styles.Title = highlightStyle

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			enterBinding,
			backBinding,
			removeBinding,
		}
	}

	l.StatusMessageLifetime = 5 * time.Second

	return saved{
		queries: q,
		spinner: s,
		list:    l,
		loading: true,
	}
}

func (s saved) Init() tea.Cmd {
	return nil
}

func (s *saved) getSavedPosts() tea.Msg {
	saved, err := s.queries.GetSavedPosts(context.Background())
	if err != nil {
		return failError{error: err}
	}

	items := make([]list.Item, len(saved))
	for _, post := range saved {
		hyperlink := fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", "Article link →", post.Url)

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
		})
	}
	return items
}

func (s saved) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s saved) View() string {
	return ""
}
