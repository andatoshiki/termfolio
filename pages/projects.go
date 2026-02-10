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
		Name: "Toshiki's Homepage",
		Desc: "All-in-one home landing page/blog/portfolio with a Nuxt 3 rebuild.",
		Tech: "Nuxt 3, Vue, Vercel, Cloudflare Workers, Netlify",
		Link: "https://github.com/andatoshiki/toshiki-home-nuxt3",
	},
	{
		Name: "Toshiki's Notebook",
		Desc: "VitePress-powered web notebook and knowledge base.",
		Tech: "VitePress, Vercel",
		Link: "https://note.toshiki.dev",
	},
	{
		Name: "Toshiki's HTTP",
		Desc: "Playful HTTP status code illustrations with the Ukuku character.",
		Tech: "Web, API",
		Link: "https://http.toshiki.dev",
	},
	{
		Name: "Toshiki's Live2D Viewer",
		Desc: "Simple web Live2D viewer built on Pixi.",
		Tech: "PixiJS, Live2D",
		Link: "https://live2d.toshiki.dev",
	},
	{
		Name: "Toshiki's Temple Block",
		Desc: "Virtual temple block to collect merits with a tap.",
		Tech: "Web",
		Link: "https://merit.toshiki.dev",
	},
	{
		Name: "Toshiki's Gallery",
		Desc: "Self-hosted photo gallery built with Hugo.",
		Tech: "Hugo",
		Link: "https://github.com/andatoshiki/toshiki-gallery",
	},
	{
		Name: "Toshiki's Mahjong Calculator",
		Desc: "Mahjong score calculator for quick scorekeeping.",
		Tech: "Web",
		Link: "https://github.com/andatoshiki/toshiki-mahjong-calc",
	},
}

func Projects() []Project {
	return projects
}

func RenderProjects(styles view.ThemeStyles, projectCursor int, themeLabel string) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("━━━ Projects ━━━"))
	b.WriteString("\n\n")

	const pageSize = 3
	start, end := projectWindow(projectCursor, len(projects), pageSize)

	for i := start; i < end; i++ {
		p := projects[i]
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
	if end < len(projects) {
		moreStyle := styles.Accent.Copy().Faint(true)
		b.WriteString(moreStyle.Render("more below!"))
		b.WriteString(styles.Help.Render(" • " + help))
	} else {
		b.WriteString(styles.Help.Render(help))
	}

	return b.String()
}

func projectWindow(cursor, total, pageSize int) (start, end int) {
	if total <= 0 {
		return 0, 0
	}
	if pageSize <= 0 || total <= pageSize {
		return 0, total
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= total {
		cursor = total - 1
	}
	start = cursor - (pageSize - 1)
	if start < 0 {
		start = 0
	}
	end = start + pageSize
	if end > total {
		end = total
		start = end - pageSize
	}
	return start, end
}
