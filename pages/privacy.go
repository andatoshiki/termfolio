package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/view"
)

const (
	privacyLeftWidth  = 12
	privacyRightWidth = 40
)

func RenderPrivacy(styles view.ThemeStyles, cursor int, trackingEnabled bool, trackingAvailable bool, themeLabel string) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("━━━ Privacy ━━━"))
	b.WriteString("\n")

	if !trackingAvailable {
		b.WriteString(styles.Content.Render("Tracking is disabled on this server."))
		b.WriteString("\n")
		b.WriteString(styles.Help.Render(themeLabel + " • esc: back to menu"))
		return b.String()
	}

	options := []struct {
		label string
		desc  string
	}{
		{label: "Yes", desc: "Allow visit tracking"},
		{label: "No", desc: "Do not track"},
	}

	b.WriteString(styles.Content.Render("Choose whether your IP is tracked."))
	b.WriteString("\n\n")

	for i, opt := range options {
		cursorMark := "  "
		if cursor == i {
			cursorMark = "→ "
		}
		leftText := cursorMark + opt.label
		leftCell := lipgloss.NewStyle().Width(privacyLeftWidth).Render(leftText)
		if cursor == i {
			leftCell = styles.Selected.Render(leftCell)
		} else {
			leftCell = styles.Menu.Render(leftCell)
		}

		rightCell := lipgloss.NewStyle().
			Width(privacyRightWidth).
			Align(lipgloss.Right).
			Render(opt.desc)
		if cursor == i {
			rightCell = styles.Selected.Copy().Bold(false).Faint(true).Render(rightCell)
		} else {
			rightCell = styles.Subtle.Render(rightCell)
		}

		b.WriteString(leftCell + rightCell)
		b.WriteString("\n")
	}

	status := "No"
	if trackingEnabled {
		status = "Yes"
	}
	b.WriteString(styles.Subtle.Render(fmt.Sprintf("Current: %s", status)))

	help := "↑/↓: select • enter: confirm • esc/backspace: menu • q: quit • " + themeLabel
	b.WriteString("\n")
	b.WriteString(styles.Help.Render(help))

	return b.String()
}
