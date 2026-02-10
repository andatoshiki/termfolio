package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/view"
)

var menuItems = []string{"About", "Projects", "Education", "Contact", "Feed", "Privacy"}
var menuDescriptions = []string{
	"Who I am",
	"Selected work",
	"Academic timeline",
	"Get in touch",
	"Latest posts",
	"Tracking control",
}

const (
	menuLeftWidth  = 16
	menuRightWidth = 40
)

func MenuItems() []string {
	return menuItems
}

func RenderMenu(styles view.ThemeStyles, menuCursor int, logoSweepIndex int, themeLabel string, visitorCount int, boxWidth int) string {
	var b strings.Builder

	logoWidth := 60
	b.WriteString(view.RenderGradientLogo(logoWidth, logoSweepIndex, styles.LogoBase, styles.LogoSnake))

	infoLine := ""
	if visitorCount > 0 {
		infoLine = fmt.Sprintf("Visits: %d", visitorCount)
	}

	b.WriteString("\n")
	if boxWidth <= 0 {
		boxWidth = logoWidth
	}
	b.WriteString(styles.Subtle.Copy().Width(boxWidth).Align(lipgloss.Left).Render(infoLine))
	b.WriteString("\n\n")

	for i, item := range menuItems {
		cursor := "  "
		if menuCursor == i {
			cursor = "→ "
		}

		leftText := cursor + item
		leftCell := lipgloss.NewStyle().Width(menuLeftWidth).Render(leftText)
		if menuCursor == i {
			leftCell = styles.Selected.Render(leftCell)
		} else {
			leftCell = styles.Menu.Render(leftCell)
		}

		desc := menuDescriptions[i]
		rightCell := lipgloss.NewStyle().
			Width(menuRightWidth).
			Align(lipgloss.Right).
			Render(desc)
		if menuCursor == i {
			rightCell = styles.Selected.Copy().Bold(false).Faint(true).Render(rightCell)
		} else {
			rightCell = styles.Subtle.Render(rightCell)
		}

		b.WriteString(leftCell + rightCell)
		b.WriteString("\n")
	}

	helpMain := "↑/↓: navigate • enter: select • esc/backspace: menu • q: quit"
	helpTheme := "t: " + themeLabel[len("t: "):]
	b.WriteString(styles.Help.Render("\n" + helpMain + "\n" + helpTheme))

	return b.String()
}
