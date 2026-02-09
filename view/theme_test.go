package view

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestThemeAt(t *testing.T) {
	if len(themePalettes) == 0 {
		t.Skip("no palettes defined")
	}

	cases := []struct {
		name  string
		index int
		want  ThemePalette
	}{
		{name: "first", index: 0, want: themePalettes[0]},
		{name: "last", index: len(themePalettes) - 1, want: themePalettes[len(themePalettes)-1]},
		{name: "negative", index: -1, want: themePalettes[0]},
		{name: "overflow", index: len(themePalettes), want: themePalettes[0]},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ThemeAt(tc.index)
			if got != tc.want {
				t.Fatalf("ThemeAt(%d) = %#v, want %#v", tc.index, got, tc.want)
			}
		})
	}
}

func TestNextThemeIndex(t *testing.T) {
	if len(themePalettes) == 0 {
		t.Skip("no palettes defined")
	}

	n := len(themePalettes)
	cases := []struct {
		name    string
		current int
		want    int
	}{
		{name: "zero", current: 0, want: (0 + 1) % n},
		{name: "wrap", current: n - 1, want: 0},
		{name: "negative", current: -2, want: 0},
		{name: "large", current: n + 2, want: (n + 3) % n},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NextThemeIndex(tc.current)
			if got != tc.want {
				t.Fatalf("NextThemeIndex(%d) = %d, want %d", tc.current, got, tc.want)
			}
		})
	}
}

func TestNewThemeStylesForegrounds(t *testing.T) {
	if len(themePalettes) == 0 {
		t.Skip("no palettes defined")
	}

	palette := themePalettes[0]
	styles := NewThemeStyles(palette)

	cases := []struct {
		name string
		got  lipgloss.TerminalColor
		want lipgloss.Color
	}{
		{name: "title", got: styles.Title.GetForeground(), want: palette.Title},
		{name: "menu", got: styles.Menu.GetForeground(), want: palette.Primary},
		{name: "selected", got: styles.Selected.GetForeground(), want: palette.Highlight},
		{name: "help", got: styles.Help.GetForeground(), want: palette.Muted},
		{name: "accent", got: styles.Accent.GetForeground(), want: palette.Link},
		{name: "tech", got: styles.Tech.GetForeground(), want: palette.Tech},
		{name: "project", got: styles.ProjectName.GetForeground(), want: palette.Project},
		{name: "logo", got: styles.LogoSnake.GetForeground(), want: palette.LogoSnake},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertColor(t, tc.got, tc.want)
		})
	}
}

func assertColor(t *testing.T, got lipgloss.TerminalColor, want lipgloss.Color) {
	t.Helper()

	gotColor, ok := got.(lipgloss.Color)
	if !ok {
		t.Fatalf("expected lipgloss.Color, got %T", got)
	}

	if gotColor != want {
		t.Fatalf("got %q, want %q", gotColor, want)
	}
}
