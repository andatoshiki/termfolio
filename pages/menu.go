package pages

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/view"
)

var menuItems = []string{"About", "Projects", "Experience", "Contact"}
var menuDescriptions = []string{
	"Who I am",
	"Selected work",
	"Roles and timeline",
	"Get in touch",
}

const (
	menuLeftWidth  = 16
	menuRightWidth = 40
)

func MenuItems() []string {
	return menuItems
}

func RenderMenu(styles view.ThemeStyles, menuCursor int, logoSweepIndex int, themeLabel string) string {
	var b strings.Builder

	b.WriteString(view.RenderGradientLogo(60, logoSweepIndex, styles.LogoBase, styles.LogoSnake))
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

	help := "↑/↓: navigate • enter: select • esc/backspace: menu • q: quit • " + themeLabel
	b.WriteString(styles.Help.Render("\n" + help))

	return b.String()
}
