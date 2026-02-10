package ui

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"

	"github.com/andatoshiki/termfolio/pages"
)

const (
	feedURL          = "https://note.toshiki.dev/feed.xml"
	feedFetchTimeout = 5 * time.Second
	feedCacheTTL     = 15 * time.Minute
	feedMaxItems     = 25
)

type feedMsg struct {
	items []pages.FeedItem
	err   error
}

func fetchFeedCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), feedFetchTimeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
		if err != nil {
			return feedMsg{err: err}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return feedMsg{err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return feedMsg{err: fmt.Errorf("feed request failed: %s", resp.Status)}
		}

		parser := gofeed.NewParser()
		feed, err := parser.Parse(resp.Body)
		if err != nil {
			return feedMsg{err: err}
		}

		items := make([]pages.FeedItem, 0, len(feed.Items))
		for _, item := range feed.Items {
			if item == nil {
				continue
			}
			title := normalizeFeedText(item.Title)
			if title == "" {
				title = "Untitled"
			}
			link := normalizeFeedText(item.Link)
			if link == "" {
				link = normalizeFeedText(feed.Link)
			}
			date := ""
			if item.PublishedParsed != nil {
				date = item.PublishedParsed.Format("01-02-2006")
			} else if item.UpdatedParsed != nil {
				date = item.UpdatedParsed.Format("01-02-2006")
			}
			items = append(items, pages.FeedItem{
				Title: title,
				Link:  link,
				Date:  date,
			})
			if len(items) >= feedMaxItems {
				break
			}
		}

		return feedMsg{items: items}
	}
}

func normalizeFeedText(value string) string {
	if value == "" {
		return value
	}
	value = unescapeHTMLEntities(value)
	value = strings.ReplaceAll(value, "\u00a0", " ")
	value = strings.ReplaceAll(value, "\u200b", "")
	return value
}

func unescapeHTMLEntities(value string) string {
	const maxPasses = 3
	for i := 0; i < maxPasses; i++ {
		unescaped := html.UnescapeString(value)
		if unescaped == value {
			break
		}
		value = unescaped
	}
	return value
}

func shouldFetchFeed(m model) bool {
	if m.feedLoading {
		return false
	}
	if len(m.feedItems) == 0 {
		return true
	}
	if m.feedFetchedAt.IsZero() {
		return true
	}
	return time.Since(m.feedFetchedAt) > feedCacheTTL
}
