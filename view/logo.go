package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var asciiLogoLines = []string{
	`             __       __           __      `,
	`            /\ \     /\ \         /\ \    `,
	`            \ \ \   /  \ \       /  \ \   `,
	`            /\ \_\ / /\ \ \     / /\ \ \  `,
	`           / /\/_// / /\ \ \   / / /\ \_\ `,
	`  __      / / /  / / /  \ \_\ / /_/_ \/_/ `,
	` /\ \    / / /  / / /   / / // /____/\    `,
	` \ \_\  / / /  / / /   / / // /\____\/    `,
	` / / /_/ / /  / / /___/ / // / /______    `,
	`/ / /__\/ /  / / /____\/ // / /_______\   `,
	`\/_______/   \/_________/ \/__________/   `,
}

func RenderGradientLogo(width int, sweepIndex int, baseStyle, snakeStyle lipgloss.Style) string {
	var result strings.Builder

	linesToShow := len(asciiLogoLines)

	maxLineLen := 0
	for i := 0; i < linesToShow; i++ {
		if len(asciiLogoLines[i]) > maxLineLen {
			maxLineLen = len(asciiLogoLines[i])
		}
	}

	if linesToShow == 0 || maxLineLen == 0 {
		return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render("")
	}

	pad := 1
	gridW := maxLineLen + pad*2
	gridH := linesToShow + pad*2

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
			baseGrid[pad+i][pad+j] = r
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

	snakeLen := 14
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
