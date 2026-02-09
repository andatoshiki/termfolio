package pages

import (
	"strings"

	"github.com/andatoshiki/termfolio/view"
)

type Project struct {
	Name string
	Desc string
	Tech string
	Link string
}

var projects = []Project{
	{
		Name: "SSH Portfolio",
		Desc: "This app",
		Tech: "Go, Bubble Tea, Wish",
		Link: "github.com/joe/ssh-portfolio",
	},
	{
		Name: "React From Scratch",
		Desc: "Built a toy React from scratch",
		Tech: "JavaScript",
		Link: "github.com/joe/react-0.5",
	},
	{
		Name: "HTTP Server From Scratch",
		Desc: "Build a HTTP server from scratch using TCP and HTTP/1.1",
		Tech: "Rust",
		Link: "github.com/joe/api-gateway",
	},
}

func Projects() []Project {
	return projects
}

func RenderProjects(styles view.ThemeStyles, projectCursor int, themeLabel string) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("━━━ Projects ━━━"))
	b.WriteString("\n\n")

	for i, p := range projects {
		cursor := "  "
		if projectCursor == i {
			cursor = "→ "
		}

		name := cursor + p.Name
		if projectCursor == i {
			b.WriteString(styles.ProjectName.Render(name))
		} else {
			b.WriteString(styles.Menu.Render(name))
		}
		b.WriteString("\n")

		// Expands project section
		if projectCursor == i {
			b.WriteString(styles.Subtle.Render("    " + p.Desc))
			b.WriteString("\n")
			b.WriteString("    ")
			b.WriteString(styles.Tech.Render(p.Tech))
			b.WriteString("\n")
			b.WriteString("    ")
			projectURL := p.Link
			if !strings.HasPrefix(projectURL, "http://") && !strings.HasPrefix(projectURL, "https://") {
				projectURL = "https://" + projectURL
			}
			b.WriteString(styles.Accent.Render(view.ClickableLink(projectURL, projectURL)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	help := themeLabel + " • ↑/↓: browse • esc: back to menu"
	b.WriteString(styles.Help.Render(help))

	return b.String()
}
