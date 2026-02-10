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

type logoCell struct {
	row int
	col int
	ch  rune
}

type logoTypewriterData struct {
	width int
	lines [][]rune
	order []logoCell
}

var typewriterLogoData = buildTypewriterLogoData()

func buildTypewriterLogoData() logoTypewriterData {
	lines := make([][]rune, len(asciiLogoLines))
	width := 0
	for i := 0; i < len(asciiLogoLines); i++ {
		lines[i] = []rune(asciiLogoLines[i])
		if len(lines[i]) > width {
			width = len(lines[i])
		}
	}

	perRow := make([][]logoCell, len(lines))
	maxRowCells := 0
	for row := 0; row < len(lines); row++ {
		for col, ch := range lines[row] {
			if ch == ' ' {
				continue
			}
			perRow[row] = append(perRow[row], logoCell{row: row, col: col, ch: ch})
		}
		if len(perRow[row]) > maxRowCells {
			maxRowCells = len(perRow[row])
		}
	}

	order := make([]logoCell, 0, maxRowCells*len(lines))
	for i := 0; i < maxRowCells; i++ {
		for row := 0; row < len(perRow); row++ {
			if i < len(perRow[row]) {
				order = append(order, perRow[row][i])
			}
		}
	}

	return logoTypewriterData{
		width: width,
		lines: lines,
		order: order,
	}
}

func LogoTypewriterRuneCount() int {
	return len(typewriterLogoData.order)
}

func RenderTypewriterLogo(width int, revealCount int, logoStyle lipgloss.Style) string {
	if revealCount < 0 {
		revealCount = 0
	}
	if revealCount > len(typewriterLogoData.order) {
		revealCount = len(typewriterLogoData.order)
	}

	var b strings.Builder
	visible := make([][]bool, len(typewriterLogoData.lines))
	for row := 0; row < len(visible); row++ {
		visible[row] = make([]bool, typewriterLogoData.width)
	}

	for i := 0; i < revealCount; i++ {
		cell := typewriterLogoData.order[i]
		visible[cell.row][cell.col] = true
	}

	for row := 0; row < len(typewriterLogoData.lines); row++ {
		line := typewriterLogoData.lines[row]
		for col := 0; col < typewriterLogoData.width; col++ {
			if col >= len(line) {
				b.WriteRune(' ')
				continue
			}
			ch := line[col]
			if ch == ' ' {
				b.WriteRune(' ')
				continue
			}
			if visible[row][col] {
				b.WriteString(logoStyle.Render(string(ch)))
				continue
			}
			b.WriteRune(' ')
		}
		if row < len(typewriterLogoData.lines)-1 {
			b.WriteRune('\n')
		}
	}

	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(b.String())
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
