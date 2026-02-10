package pages

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/view"
)

const splashIntroText = "Hi, welcome to termfolio"
const splashCommandBar = "enter: continue"

var splashRunes = []rune(splashIntroText)

func SplashRuneCount() int {
	return len(splashRunes)
}

func RenderSplash(styles view.ThemeStyles, revealCount int, cursorVisible bool, boxWidth int) string {
	total := len(splashRunes)
	if revealCount < 0 {
		revealCount = 0
	}
	if revealCount > total {
		revealCount = total
	}

	contentWidth := splashContentWidth(boxWidth)
	text := string(splashRunes[:revealCount])

	cursor := " "
	if cursorVisible {
		cursor = styles.Accent.Copy().Bold(true).Render("▌")
	}

	var b strings.Builder
	line := styles.Content.Render(text) + cursor
	b.WriteString(lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(line))
	b.WriteString("\n\n")
	b.WriteString(styles.Help.Copy().Width(contentWidth).Align(lipgloss.Center).Render(splashCommandBar))

	return b.String()
}

func splashContentWidth(boxWidth int) int {
	if boxWidth <= 0 {
		return 60
	}
	if boxWidth > 4 {
		return boxWidth - 4
	}
	return boxWidth
}
