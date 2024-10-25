package rss

import (
	"context"
	"fmt"
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewModel struct {
	dbQueries *database.Queries
	feeds     []database.Feed
	err       error
}

func (v ViewModel) Init() tea.Cmd {
	return v.getFeeds
}

func (v ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return v, tea.Quit
		}
		switch msg.String() {
		case "r":
			return v, v.getFeeds
		}
	case successFeed:
		v.feeds[0] = msg.feed
		return v, nil

	case dbError:
		v.err = msg.dbErr
		return v, nil

	case dbSuccess:
		v.feeds = msg.feeds
		return v, nil
	}
	return v, nil
}

func (v ViewModel) View() string {
	builder := strings.Builder{}
	builder.WriteString("VIEW feeds\n")
	if v.err != nil {
		builder.WriteString("error gettings feeds from db")
		builder.WriteString("\n")
		builder.WriteString(v.err.Error())
	}

	if len(v.feeds) == 0 {
		builder.WriteString("no feeds found in db")
	}

	for _, feed := range v.feeds {
		s := fmt.Sprintf("name: %s\n\nurl: %s\n", feed.Name, feed.Url)

		builder.WriteString(s)
	}

	return builder.String()
}

func InitialiseViewModel(queries *database.Queries) ViewModel {
	return ViewModel{
		dbQueries: queries,
		feeds:     make([]database.Feed, 0),
		err:       nil,
	}
}

func (v *ViewModel) getFeeds() tea.Msg {
	data, err := v.dbQueries.GetFeeds(context.Background())
	if err != nil {
		return dbError{dbErr: err}
	}
	return dbSuccess{feeds: data}
}

type successFeed struct {
	feed database.Feed
}

type dbSuccess struct {
	feeds []database.Feed
}
