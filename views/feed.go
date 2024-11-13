package views

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/IvanYaremko/rssdukester/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type feed struct {
	queries *database.Queries
	item    item
	spinner spinner.Model
	loading bool
	list    list.Model
}

func initialiseFeed(q *database.Queries, i item) feed {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.HighlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, list.NewDefaultDelegate(), 100, 40)
	l.Title = fmt.Sprintf("%s FEED", strings.ToUpper(i.title))
	l.Styles.Title = styles.HighlightStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			enterBinding,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			backBinding,
		}
	}

	return feed{
		queries: q,
		item:    i,
		spinner: s,
		loading: true,
		list:    l,
	}
}

type rssResponse struct {
	Channel struct {
		Title       string            `xml:"title"`
		Link        string            `xml:"link"`
		Description string            `xml:"description"`
		Item        []rssItemResponse `xml:"item"`
	} `xml:"channel"`
}

type rssItemResponse struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (f feed) fetchRssFeed() tea.Msg {
	req, err := http.NewRequestWithContext(context.Background(), "GET", f.item.url, nil)
	if err != nil {
		return failError{error: err}
	}

	req.Header.Add("User-Agent", "rssdukester")
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return failError{error: err}
	}
	defer response.Body.Close()

	decoder := xml.NewDecoder(response.Body)
	rss := rssResponse{}
	if err := decoder.Decode(&rss); err != nil {
		return failError{error: err}
	}

	items := make([]list.Item, len(rss.Channel.Item))
	for i, val := range rss.Channel.Item {
		items[i] = item{
			title:       val.Title,
			description: val.Link,
			url:         val.Link,
		}
	}

	return successItems{items: items}
}

func (f feed) Init() tea.Cmd {
	return tea.Batch(f.spinner.Tick, f.fetchRssFeed)
}

func (f feed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return f, tea.Quit
		case key.Matches(msg, backBinding):
			rssList := initialiseRssList(f.queries)
			return rssList, rssList.Init()
		}
	case successItems:
		cmd := f.list.SetItems(msg.items)
		f.loading = false
		return f, cmd
	}

	newList, cmd := f.list.Update(msg)
	cmds = append(cmds, cmd)
	f.list = newList

	f.spinner, cmd = f.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return f, tea.Batch(cmds...)
}

func (f feed) View() string {
	s := strings.Builder{}

	if f.loading {
		s.WriteString(fmt.Sprintf("%s loading %s...", f.spinner.View(), f.item.title))
	} else {
		s.WriteString(f.list.View())
	}

	return base.Render(s.String())
}
