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

var formateDate = "2006/01/02 15:04"

type feed struct {
	queries *database.Queries
	rss     item
	spinner spinner.Model
	loading bool
	list    list.Model
	err     error
}

func initialiseFeed(q *database.Queries, rss item) feed {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, feedItemDelegate(q, rss), width, height)
	l.Title = fmt.Sprintf("%s FEED", strings.ToUpper(rss.title))
	l.Styles.Title = highlightStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			enterBinding,
			saveBinding,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			backBinding,
			saveBinding,
		}
	}
	l.StatusMessageLifetime = 3 * time.Second

	return feed{
		queries: q,
		rss:     rss,
		spinner: s,
		loading: true,
		list:    l,
		err:     nil,
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

		hyperLink := specialStyle.Render(
			fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", url, "Article link →"),
		)

		t, _ := parseDateTime(val.PubDate)
		pubDate := attentionStyle.Italic(true).Render(t.Format(formateDate))

		items[i] = item{
			title:       val.Title,
			description: fmt.Sprintf("%s • %s", hyperLink, pubDate),
			url:         url,
		}
	}

	return successItems{items: items}
}

func parseDateTime(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		"2006-01-02T15:04:05-0700",
		"2006-01-02T15:04:05Z",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func (f feed) Init() tea.Cmd {
	return tea.Batch(f.spinner.Tick, f.fetchRssFeed)
}

func (f feed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.list.SetSize(msg.Width-20, msg.Height-2)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return f, tea.Quit
		case key.Matches(msg, quitBinding):
			return f, tea.Quit
		case key.Matches(msg, backBinding):
			rssList := initialiseRssList(f.queries)
			return rssList, rssList.Init()
		}
	case successItems:
		cmd = f.list.SetItems(msg.items)
		f.loading = false
		return f, cmd

	case selected:
		article := InitialiseArticle(f.queries, f.rss, msg.selected, false)
		return article, article.Init()

	case success:
		selected := f.list.SelectedItem().(item)
		message := fmt.Sprintf("%s %s",
			attentionStyle.Bold(true).Render("SAVED"),
			specialStyle.Italic(true).Render(selected.title),
		)
		cmd = f.list.NewStatusMessage(message)
		return f, cmd

	case fail:
		cmd = f.list.NewStatusMessage(
			errorStyle.Render("already saved!"),
		)
		return f, cmd

	case failError:
		f.err = msg.error
		f.loading = false
		return f, nil
	}

	f.list, cmd = f.list.Update(msg)
	cmds = append(cmds, cmd)

	f.spinner, cmd = f.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return f, tea.Batch(cmds...)
}

func (f feed) View() string {
	s := strings.Builder{}
	if f.err != nil {
		return fmt.Sprintf("%s\n%s",
			errorStyle.Render("something has gone wrong!"),
			errorStyle.Render(f.err.Error()),
		)
	}

	if f.loading {
		s.WriteString(fmt.Sprintf("%s loading %s...", f.spinner.View(), f.rss.title))
	} else {
		s.WriteString(f.list.View())
	}

	return baseStyle.Render(s.String())
}
