package views

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type searchResult struct {
	items []list.Item
	err   error
}

type searchList struct {
	queries    *database.Queries
	searchTerm string
	list       list.Model
	spinner    spinner.Model
	loading    bool
}

func initialiseSearchList(q *database.Queries, st string) searchList {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	items := make([]list.Item, 0)
	l := list.New(items, feedItemDelegate(), width, height)
	l.Title = fmt.Sprintf("%s %s",
		"SEARCH RESULT FOR",
		specialStyle.Italic(true).Render(st),
	)
	l.Styles.Title = highlightStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			enterBinding,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			escBinding,
		}
	}

	return searchList{
		queries:    q,
		list:       l,
		spinner:    s,
		searchTerm: st,
		loading:    true,
	}
}

func (sl searchList) Init() tea.Cmd {
	return tea.Batch(sl.spinner.Tick, sl.performSearch)
}

func (sl searchList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sl.list.SetSize(msg.Width-20, msg.Height-2)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ctrlcBinding):
			return sl, tea.Quit
		case key.Matches(msg, quitBinding):
			return sl, tea.Quit
		case key.Matches(msg, escBinding):
			search := initialiseSearch(sl.queries, sl.searchTerm)
			return search, search.Init()
		}

	case successItems:
		cmd = sl.list.SetItems(msg.items)
		sl.loading = false
		return sl, cmd

	case selected:
		article := InitialiseArticle(sl.queries, item{title: msg.selected.feedName}, msg.selected, backToSearch, sl.searchTerm)
		return article, article.Init()
	}
	sl.list, cmd = sl.list.Update(msg)
	cmds = append(cmds, cmd)

	sl.spinner, cmd = sl.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return sl, tea.Batch(cmds...)
}

func (sl searchList) View() string {
	sb := strings.Builder{}

	if sl.loading {
		return baseStyle.Render(fmt.Sprintf("%s %s",
			sl.spinner.View(),
			highlightStyle.Render("Searching...")),
		)
	}

	sb.WriteString(sl.list.View())
	return baseStyle.Render(sb.String())
}

func (sl searchList) performSearch() tea.Msg {

	feeds, err := sl.queries.GetFeeds(context.Background())
	if err != nil {
		return failError{error: err}
	}

	results := make(chan searchResult)
	wg := sync.WaitGroup{}

	for _, feed := range feeds {
		wg.Add(1)
		go func(feed database.Feed) {
			defer wg.Done()

			msg := fetchRssFeed(feed.Url)()
			if successItems, ok := msg.(successItems); ok {
				filteredItems := []list.Item{}
				searchTermLower := strings.ToLower(sl.searchTerm)

				for _, sI := range successItems.items {
					if feedItem, ok := sI.(item); ok {
						if strings.Contains(strings.ToLower(feedItem.title), searchTermLower) {
							filteredItems = append(filteredItems, item{
								title:       feedItem.title,
								url:         feedItem.url,
								pubDate:     feedItem.pubDate,
								description: fmt.Sprintf("%s • %s", feedItem.description, attentionStyle.Render(feed.Name)),
								feedId:      feed.ID,
								feedName:    feed.Name,
							})
						}
					}
				}
				results <- searchResult{items: filteredItems}
			} else if errMsg, ok := msg.(failError); ok {
				results <- searchResult{err: errMsg.error}
			}
		}(feed)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	allItems := []list.Item{}
	errors := []error{}

	for result := range results {
		if result.err != nil {
			errors = append(errors, result.err)
		}
		allItems = append(allItems, result.items...)
	}

	if len(errors) > 0 {
		return failError{error: errors[0]}
	}

	return successItems{items: allItems}
}
