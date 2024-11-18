package views

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type feed struct {
	queries *database.Queries
	rss     item
	spinner spinner.Model
	loading bool
	list    list.Model
}

func initialiseFeed(q *database.Queries, i item) feed {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, feedItemDelegate(), 100, 40)
	l.Title = fmt.Sprintf("%s FEED", strings.ToUpper(i.title))
	l.Styles.Title = highlightStyle
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
		rss:     i,
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
	Guid        string `xml:"guid"`
	Atom        string `xml:"atom,link"`
}

func (f feed) fetchRssFeed() tea.Msg {
	req, err := http.NewRequestWithContext(context.Background(), "GET", f.rss.url, nil)
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
		url := val.Link

		if url == "" {
			url = val.Guid
		}
		if url == "" {
			url = val.Atom
		}

		timestamp, _ := time.Parse(time.RFC1123Z, val.PubDate)
		date := timestamp.Format("06 Jan Mon 15:04")
		items[i] = item{
			title:       val.Title,
			description: date,
			url:         url,
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
	case tea.WindowSizeMsg:
		f.list.SetSize(msg.Width-20, msg.Height-2)

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

	case selected:
		article := InitialiseArticle(f.queries, f.rss, msg.selected)
		return article, article.Init()

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
		s.WriteString(fmt.Sprintf("%s loading %s...", f.spinner.View(), f.rss.title))
	} else {
		s.WriteString(f.list.View())
	}

	return baseStyle.Render(s.String())
}
