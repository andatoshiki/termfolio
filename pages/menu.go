package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/view"
)

var menuItems = []string{"About", "Projects", "Experience", "Contact", "Privacy", "Feed"}
var menuDescriptions = []string{
	"Who I am",
	"Selected work",
	"Roles and timeline",
	"Get in touch",
	"Tracking control",
	"Latest posts",
}

const (
	menuLeftWidth  = 16
	menuRightWidth = 40
)

func MenuItems() []string {
	return menuItems
}

func RenderMenu(styles view.ThemeStyles, menuCursor int, logoSweepIndex int, themeLabel string, visitorCount int) string {
	var b strings.Builder

	b.WriteString(view.RenderGradientLogo(60, logoSweepIndex, styles.LogoBase, styles.LogoSnake))
	if visitorCount > 0 {
		b.WriteString("\n")
		b.WriteString(styles.Subtle.Render(fmt.Sprintf("Visits: %d", visitorCount)))
		b.WriteString("\n\n")
	} else {
		b.WriteString("\n\n")
	}

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

	help := "↑/↓: navigate • enter: select • esc/backspace: menu • q: quit • " + themeLabel
	b.WriteString(styles.Help.Render("\n" + help))

	return b.String()
}
