package pages

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/view"
)

const splashIntroPrefix = "Hi! Welcome to "
const splashIntroName = "Toshiki's"
const splashIntroSuffix = " termfolio, say hi to me!"
const splashIntroLink = "https://toshiki.dev"
const splashIntroText = splashIntroPrefix + splashIntroName + splashIntroSuffix
const splashOpenSourcePrefix = "open sourced on "
const splashOpenSourceLabel = "github"
const splashOpenSourceLink = "https://github.com/andatoshiki/termfolio"
const splashCommandBar = "enter: continue"
const splashTickMillis = 45
const splashBlinkIntervalMillis = 500

var splashRunes = []rune(splashIntroText)

func SplashRuneCount() int {
	return len(splashRunes)
}

func RenderSplash(styles view.ThemeStyles, revealCount int, blinkStep int, boxWidth int) string {
	total := len(splashRunes)
	if revealCount < 0 {
		revealCount = 0
	}
	if revealCount > total {
		revealCount = total
	}

	contentWidth := splashContentWidth(boxWidth)
	text := renderSplashText(styles, revealCount, total)
	cursor := ""
	if revealCount >= total {
		cursor = " " + renderSplashCursor(styles, blinkStep)
	}

	var b strings.Builder
	line := text + cursor
	b.WriteString(lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(line))
	b.WriteString("\n")
	openSourceLine := splashOpenSourcePrefix + view.ClickableLink(splashOpenSourceLabel, splashOpenSourceLink)
	b.WriteString(styles.Accent.Copy().Bold(false).Faint(true).Width(contentWidth).Align(lipgloss.Center).Render(openSourceLine))
	b.WriteString("\n")
	b.WriteString(styles.Help.Copy().Width(contentWidth).Align(lipgloss.Center).Render(splashCommandBar))

	return b.String()
}

func renderSplashText(styles view.ThemeStyles, revealCount, total int) string {
	if revealCount < total {
		return styles.Content.Render(string(splashRunes[:revealCount]))
	}
	name := view.ClickableLink(styles.Accent.Copy().Bold(false).Underline(true).Render(splashIntroName), splashIntroLink)
	return styles.Content.Render(splashIntroPrefix) + name + styles.Content.Render(splashIntroSuffix)
}

func renderSplashCursor(styles view.ThemeStyles, blinkStep int) string {
	if blinkStep < 0 {
		blinkStep = 0
	}
	phase := (blinkStep * splashTickMillis) / splashBlinkIntervalMillis
	if phase%2 == 0 {
		return styles.Accent.Copy().Bold(true).Render("█")
	}
	return " "
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
