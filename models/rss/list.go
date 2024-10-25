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
	tea.Println("ViewModel init called")
	return getFeeds(v.dbQueries)
}

func (v ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return v, tea.Quit
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

func getFeeds(q *database.Queries) tea.Cmd {
	return func() tea.Msg {
		tea.Println("getfeeds invoked")
		data, err := q.GetFeeds(context.Background())
		tea.Println("getfeeds feeds", data)
		if err != nil {
			return dbError{dbErr: err}
		}
		return dbSuccess{feeds: data}
	}

}

func (v ViewModel) getFeed() tea.Msg {
	data, err := v.dbQueries.GetFeedById(context.Background(), 0)
	if err != nil {
		return dbError{dbErr: err}
	}
	return successFeed{feed: data}
}

type successFeed struct {
	feed database.Feed
}

type dbSuccess struct {
	feeds []database.Feed
}
