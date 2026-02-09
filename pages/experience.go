package pages

import (
	"fmt"
	"strings"

	"github.com/andatoshiki/termfolio/view"
)

type Experience struct {
	Role    string
	Company string
	Period  string
	Desc    string
}

var experiences = []Experience{
	{
		Role:    "Software Engineer",
		Company: "Microsoft",
		Period:  "2025 - Present",
		Desc:    "Azure SQL VM team",
	},
	{
		Role:    "Software Engineer Intern",
		Company: "Jenni AI",
		Period:  "2024 - 2025",
		Desc:    "Developed new product that reviews manuscripts for Jenni AI",
	},
	{
		Role:    "Software Engineer Intern",
		Company: "Blue Origin",
		Period:  "Fall 2023",
		Desc:    "New Glenn Rocket Software",
	},
}

func Experiences() []Experience {
	return experiences
}

func RenderExperience(styles view.ThemeStyles, expCursor int, themeLabel string) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("━━━ Experience ━━━"))
	b.WriteString("\n\n")

	for i, exp := range experiences {
		cursor := "  "
		if expCursor == i {
			cursor = "→ "
		}

		line := fmt.Sprintf("%s%s @ %s",
			cursor,
			styles.Role.Render(exp.Role),
			styles.Company.Render(exp.Company))
		b.WriteString(line)
		b.WriteString("\n")

		b.WriteString("    ")
		b.WriteString(styles.Period.Render(exp.Period))
		b.WriteString("\n")

		if expCursor == i {
			b.WriteString("    ")
			b.WriteString(styles.Content.Render(exp.Desc))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	help := themeLabel + " • ↑/↓: browse • esc: back to menu"
	b.WriteString(styles.Help.Render(help))

	return b.String()
}
