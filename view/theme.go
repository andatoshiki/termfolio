package view

import "github.com/charmbracelet/lipgloss"

type ThemePalette struct {
	Name      string
	Primary   lipgloss.Color
	Muted     lipgloss.Color
	Title     lipgloss.Color
	Link      lipgloss.Color
	Highlight lipgloss.Color
	Tech      lipgloss.Color
	Project   lipgloss.Color
	LogoSnake lipgloss.Color
}

type ThemeStyles struct {
	Title       lipgloss.Style
	Menu        lipgloss.Style
	Selected    lipgloss.Style
	Help        lipgloss.Style
	Content     lipgloss.Style
	Accent      lipgloss.Style
	Subtle      lipgloss.Style
	ProjectName lipgloss.Style
	Tech        lipgloss.Style
	Role        lipgloss.Style
	Company     lipgloss.Style
	Period      lipgloss.Style
	LogoBase    lipgloss.Style
	LogoSnake   lipgloss.Style
}

var themePalettes = []ThemePalette{
	{
		Name:      "Tokyo Night",
		Primary:   lipgloss.Color("#C0CAF5"),
		Muted:     lipgloss.Color("#565F89"),
		Title:     lipgloss.Color("#7AA2F7"),
		Link:      lipgloss.Color("#7DCFFF"),
		Highlight: lipgloss.Color("#BB9AF7"),
		Tech:      lipgloss.Color("#9ECE6A"),
		Project:   lipgloss.Color("#7AA2F7"),
		LogoSnake: lipgloss.Color("#9EC5FF"),
	},
	{
		Name:      "Nord",
		Primary:   lipgloss.Color("#D8DEE9"),
		Muted:     lipgloss.Color("#4C566A"),
		Title:     lipgloss.Color("#81A1C1"),
		Link:      lipgloss.Color("#88C0D0"),
		Highlight: lipgloss.Color("#B48EAD"),
		Tech:      lipgloss.Color("#A3BE8C"),
		Project:   lipgloss.Color("#EBCB8B"),
		LogoSnake: lipgloss.Color("#8FBCBB"),
	},
	{
		Name:      "Gruvbox",
		Primary:   lipgloss.Color("#EBDBB2"),
		Muted:     lipgloss.Color("#928374"),
		Title:     lipgloss.Color("#83A598"),
		Link:      lipgloss.Color("#8EC07C"),
		Highlight: lipgloss.Color("#D3869B"),
		Tech:      lipgloss.Color("#B8BB26"),
		Project:   lipgloss.Color("#FABD2F"),
		LogoSnake: lipgloss.Color("#FE8019"),
	},
	{
		Name:      "Catppuccin Mocha",
		Primary:   lipgloss.Color("#CDD6F4"),
		Muted:     lipgloss.Color("#6C7086"),
		Title:     lipgloss.Color("#89B4FA"),
		Link:      lipgloss.Color("#89DCEB"),
		Highlight: lipgloss.Color("#CBA6F7"),
		Tech:      lipgloss.Color("#A6E3A1"),
		Project:   lipgloss.Color("#F9E2AF"),
		LogoSnake: lipgloss.Color("#F5C2E7"),
	},
	{
		Name:      "Rose Pine",
		Primary:   lipgloss.Color("#E0DEF4"),
		Muted:     lipgloss.Color("#6E6A86"),
		Title:     lipgloss.Color("#9CCFD8"),
		Link:      lipgloss.Color("#C4A7E7"),
		Highlight: lipgloss.Color("#EBBCBA"),
		Tech:      lipgloss.Color("#31748F"),
		Project:   lipgloss.Color("#F6C177"),
		LogoSnake: lipgloss.Color("#EB6F92"),
	},
}

func ThemeAt(index int) ThemePalette {
	if len(themePalettes) == 0 {
		return ThemePalette{}
	}
	if index < 0 || index >= len(themePalettes) {
		return themePalettes[0]
	}
	return themePalettes[index]
}

func NextThemeIndex(current int) int {
	if len(themePalettes) == 0 {
		return 0
	}
	if current < 0 {
		return 0
	}
	return (current + 1) % len(themePalettes)
}

func NewThemeStyles(p ThemePalette) ThemeStyles {
	return ThemeStyles{
		Title: lipgloss.NewStyle().
			Foreground(p.Title).
			Bold(true).
			MarginBottom(1),
		Menu: lipgloss.NewStyle().
			Foreground(p.Primary),
		Selected: lipgloss.NewStyle().
			Foreground(p.Highlight).
			Bold(true),
		Help: lipgloss.NewStyle().
			Foreground(p.Muted).
			MarginTop(1),
		Content: lipgloss.NewStyle().
			Foreground(p.Primary),
		Accent: lipgloss.NewStyle().
			Foreground(p.Link).
			Bold(true),
		Subtle: lipgloss.NewStyle().
			Foreground(p.Muted),
		ProjectName: lipgloss.NewStyle().
			Foreground(p.Project).
			Bold(true),
		Tech: lipgloss.NewStyle().
			Foreground(p.Tech),
		Role: lipgloss.NewStyle().
			Foreground(p.Highlight).
			Bold(true),
		Company: lipgloss.NewStyle().
			Foreground(p.Link),
		Period: lipgloss.NewStyle().
			Foreground(p.Muted).
			Italic(true),
		LogoBase: lipgloss.NewStyle().
			Foreground(p.Title).
			Bold(true),
		LogoSnake: lipgloss.NewStyle().
			Foreground(p.LogoSnake).
			Bold(true),
	}
}
