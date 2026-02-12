package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/counter"
	"github.com/andatoshiki/termfolio/version"
	"github.com/andatoshiki/termfolio/view"
)

const (
	privacyOptionLabelWidth = 10
	privacyOptionDescWidth  = 30
	privacyCountryNameWidth = 18
)

func RenderPrivacy(
	styles view.ThemeStyles,
	cursor int,
	trackingEnabled bool,
	trackingAvailable bool,
	themeLabel string,
	statsEnabled bool,
	statsTotal int,
	statsTopCountries []counter.CountryCount,
	statsError string,
) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("━━━ Privacy ━━━"))
	b.WriteString("\n")

	b.WriteString(styles.Content.Render(version.AppDesc))
	b.WriteString("\n")
	b.WriteString(styles.Subtle.Render("Version: v" + version.VersionString()))
	b.WriteString("\n\n")

	if !trackingAvailable {
		b.WriteString(styles.Content.Render("Tracking is disabled on this server."))
	} else {
		options := []struct {
			label string
			desc  string
		}{
			{label: "Yes", desc: "Allow visit tracking"},
			{label: "No", desc: "Do not track"},
		}

		b.WriteString(styles.Content.Render("Choose whether your IP is tracked."))
		b.WriteString("\n\n")

		var optionRows []string
		for i, opt := range options {
			optionRows = append(optionRows, renderPrivacyOptionRow(styles, cursor == i, opt.label, opt.desc))
		}
		b.WriteString(lipgloss.NewStyle().PaddingLeft(2).Render(strings.Join(optionRows, "\n")))
		b.WriteString("\n")

		status := "No"
		if trackingEnabled {
			status = "Yes"
		}
		b.WriteString(styles.Subtle.Render(fmt.Sprintf("Current: %s", status)))
	}

	if statsEnabled {
		b.WriteString("\n\n")
		b.WriteString(styles.Title.Render("━━━ Stats ━━━"))
		b.WriteString("\n")

		if statsError != "" {
			b.WriteString(styles.Subtle.Render("Unavailable: " + statsError))
		} else {
			statsLines := []string{
				styles.Content.Render(fmt.Sprintf("Total unique visitors: %d", statsTotal)),
			}
			if len(statsTopCountries) == 0 {
				statsLines = append(statsLines, styles.Subtle.Render("Top 5 countries: N/A"))
			} else {
				statsLines = append(statsLines, styles.Content.Render("Top 5 countries:"))
				for i, country := range statsTopCountries {
					label := fmt.Sprintf("%d. %-*s %d", i+1, privacyCountryNameWidth, country.Name, country.Visitors)
					statsLines = append(statsLines, styles.Content.Render(label))
				}
			}

			b.WriteString(strings.Join(statsLines, "\n"))
		}
	}

	help := "↑/↓: select • enter: confirm • esc/backspace: menu • q: quit • " + themeLabel
	b.WriteString("\n")
	b.WriteString(styles.Help.Render(help))

	return b.String()
}

func renderPrivacyOptionRow(styles view.ThemeStyles, selected bool, label string, description string) string {
	cursorMark := "  "
	leftStyle := styles.Menu
	rightStyle := styles.Subtle
	if selected {
		cursorMark = "→ "
		leftStyle = styles.Selected
		rightStyle = styles.Selected.Copy().Bold(false).Faint(true)
	}

	leftCell := lipgloss.NewStyle().
		Width(privacyOptionLabelWidth).
		Align(lipgloss.Left).
		Render(cursorMark + label)
	rightCell := lipgloss.NewStyle().
		Width(privacyOptionDescWidth).
		Align(lipgloss.Left).
		Render(description)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftCell),
		rightStyle.Render(rightCell),
	)
}
