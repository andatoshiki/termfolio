package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/view"
)

func RenderContact(styles view.ThemeStyles, themeLabel string) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("━━━ Contact ━━━"))
	b.WriteString("\n")

	b.WriteString(styles.Content.Render("Feel free to reach out!"))
	b.WriteString("\n\n")

	githubURL := "https://github.com/andatoshiki"
	mailtoURL := "mailto:hi@tosh1ki.de"
	telegramURL := "https://t.me/@andatoshiki"
	twitterURL := "https://x.com/andatoshiki"
	mastodonURL := "https://mastodon.social/@andatoshiki"

	leftLines := []string{
		styles.Accent.Render("Contacts"),
		styles.Content.Render(fmt.Sprintf("Email      %s", view.ClickableLink("hi@tosh1ki.de", mailtoURL))),
		styles.Content.Render(fmt.Sprintf("Telegram   %s", view.ClickableLink("@andatoshiki", telegramURL))),
		styles.Content.Render(fmt.Sprintf("GitHub     %s", view.ClickableLink("@andatoshiki", githubURL))),
	}

	rightLines := []string{
		styles.Accent.Render("Social"),
		styles.Content.Render(fmt.Sprintf("Twitter(X) %s", view.ClickableLink("andatoshiki", twitterURL))),
		styles.Content.Render(fmt.Sprintf("Mastodon   %s", view.ClickableLink("@andatoshiki", mastodonURL))),
		"",
	}

	leftColumn := lipgloss.NewStyle().Width(32).Render(strings.Join(leftLines, "\n"))
	rightColumn := lipgloss.NewStyle().Width(32).Render(strings.Join(rightLines, "\n"))
	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)

	b.WriteString(columns)

	b.WriteString("\n")
	b.WriteString(styles.Help.Render(themeLabel + " • esc: back to menu"))

	return b.String()
}
