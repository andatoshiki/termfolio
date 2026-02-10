package view

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

var asciiLogoLines = []string{
	`                         d8b        d8, d8b         d8,`,
	`   d8P                   ?88       ` + "`" + `8P  ?88        ` + "`" + `8P `,
	`d888888P                  88b            88b           `,
	`  ?88'   d8888b  .d888b,  888888b   88b  888  d88'  88b`,
	`  88P   d8P' ?88 ?8b,     88P ` + "`" + `?8b  88P  888bd8P'   88P`,
	`  88b   88b  d88   ` + "`" + `?8b  d88   88P d88  d88888b    d88 `,
	`  ` + "`" + `?8b  ` + "`" + `?8888P'` + "`" + `?888P' d88'   88bd88' d88' ` + "`" + `?88b,d88' `,
}

func RenderGradientLogo(width int, sweepIndex int, baseStyle, snakeStyle lipgloss.Style) string {
	var result strings.Builder

	linesToShow := len(asciiLogoLines)

	maxLineLen := 0
	for i := 0; i < linesToShow; i++ {
		lineLen := utf8.RuneCountInString(asciiLogoLines[i])
		if lineLen > maxLineLen {
			maxLineLen = lineLen
		}
	}

	if linesToShow == 0 || maxLineLen == 0 {
		return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render("")
	}

	padY := 1
	padX := 2
	gridW := maxLineLen + padX*2
	gridH := linesToShow + padY*2

	baseGrid := make([][]rune, gridH)
	for y := 0; y < gridH; y++ {
		row := make([]rune, gridW)
		for x := 0; x < gridW; x++ {
			row[x] = ' '
		}
		baseGrid[y] = row
	}

	for i := 0; i < linesToShow; i++ {
		lineRunes := []rune(asciiLogoLines[i])
		for j, r := range lineRunes {
			baseGrid[padY+i][padX+j] = r
		}
	}

	type pt struct{ x, y int }
	path := make([]pt, 0, gridW*2+gridH*2)

	for x := 0; x < gridW; x++ {
		path = append(path, pt{x: x, y: 0})
	}
	for y := 1; y < gridH-1; y++ {
		path = append(path, pt{x: gridW - 1, y: y})
	}
	if gridH > 1 {
		for x := gridW - 1; x >= 0; x-- {
			path = append(path, pt{x: x, y: gridH - 1})
		}
	}
	if gridW > 1 {
		for y := gridH - 2; y >= 1; y-- {
			path = append(path, pt{x: 0, y: y})
		}
	}

	pathLen := len(path)
	if pathLen == 0 {
		return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render("")
	}

	snakeLen := pathLen / 8
	if snakeLen < 6 {
		snakeLen = 6
	}
	if snakeLen > pathLen {
		snakeLen = pathLen
	}
	start := sweepIndex % pathLen

	snakeGrid := make([][]bool, gridH)
	for y := 0; y < gridH; y++ {
		snakeGrid[y] = make([]bool, gridW)
	}
	for i := 0; i < snakeLen; i++ {
		idx := (start + i) % pathLen
		p := path[idx]
		snakeGrid[p.y][p.x] = true
	}

	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			if snakeGrid[y][x] {
				result.WriteString(snakeStyle.Render("•"))
				continue
			}
			ch := string(baseGrid[y][x])
			if baseGrid[y][x] == ' ' {
				result.WriteString(ch)
			} else {
				result.WriteString(baseStyle.Render(ch))
			}
		}
		if y < gridH-1 {
			result.WriteString("\n")
		}
	}

	logoBlock := result.String()
	centered := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(logoBlock)
	return centered
}
