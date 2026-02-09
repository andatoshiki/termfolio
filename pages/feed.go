package pages

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/view"
)

const (
	feedLeftWidth  = 52
	feedRightWidth = 12
)

type FeedItem struct {
	Title string
	Link  string
	Date  string
}

func RenderFeed(styles view.ThemeStyles, items []FeedItem, cursor int, offset int, pageSize int, loading bool, errMsg string, themeLabel string) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("━━━ Feed ━━━"))
	b.WriteString("\n\n")

	if loading {
		b.WriteString(styles.Subtle.Render("Loading feed..."))
		b.WriteString("\n")
		b.WriteString(styles.Help.Render(themeLabel + " • esc: back to menu"))
		return b.String()
	}

	if errMsg != "" {
		b.WriteString(styles.Subtle.Render("Failed to load feed."))
		b.WriteString("\n")
		b.WriteString(styles.Subtle.Render(errMsg))
		b.WriteString("\n")
		b.WriteString(styles.Help.Render(themeLabel + " • esc: back to menu"))
		return b.String()
	}

	if len(items) == 0 {
		b.WriteString(styles.Subtle.Render("No posts found."))
		b.WriteString("\n")
		b.WriteString(styles.Help.Render(themeLabel + " • esc: back to menu"))
		return b.String()
	}

	if pageSize <= 0 {
		pageSize = len(items)
	}
	start := offset
	if start < 0 {
		start = 0
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	for i := start; i < end; i++ {
		item := items[i]
		cursorMark := "  "
		if i == cursor {
			cursorMark = "→ "
		}

		title := truncate(item.Title, feedLeftWidth-3)
		if item.Link != "" {
			title = view.ClickableLink(title, item.Link)
		}
		leftCell := lipgloss.NewStyle().Width(feedLeftWidth).Render(cursorMark + title)
		if i == cursor {
			leftCell = styles.Selected.Render(leftCell)
		} else {
			leftCell = styles.Content.Render(leftCell)
		}

		dateText := item.Date
		rightCell := lipgloss.NewStyle().
			Width(feedRightWidth).
			Align(lipgloss.Right).
			Render(dateText)
		if i == cursor {
			rightCell = styles.Selected.Copy().Bold(false).Faint(true).Render(rightCell)
		} else {
			rightCell = styles.Subtle.Render(rightCell)
		}

		b.WriteString(leftCell)
		b.WriteString(rightCell)
		b.WriteString("\n")
	}

	help := themeLabel + " • ↑/↓: browse • esc: back to menu"
	b.WriteString(styles.Help.Render("\n" + help))

	return b.String()
}

func truncate(value string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}
