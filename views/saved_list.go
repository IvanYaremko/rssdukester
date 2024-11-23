package views

import (
	"time"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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
