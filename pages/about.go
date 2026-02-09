package pages

import (
	"math"
	"strings"

	"github.com/andatoshiki/termfolio/view"
)

var aboutContent = `
Hey, I'm Joe, a software developer interested in building entertaining or useful things.

Currently exploring React Internals and distributed systems.
`

var aboutRunes = []rune(aboutContent)
var scrambleRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:'\",.<>?/~")

const settleDurationTicks = 8

func aboutVisible(count int, scrambleTick int) string {
	if count <= 0 {
		return ""
	}
	total := len(aboutRunes)
	if count > total {
		count = total
	}

	if count >= total && scrambleTick >= count+AboutSettleTicks() {
		return aboutContent
	}

	out := make([]rune, count)
	i := 0
	for i < count {
		ch := aboutRunes[i]
		if isWhitespace(ch) {
			out[i] = ch
			i++
			continue
		}

		wordStart := i
		wordEnd := i
		for wordEnd < total && !isWhitespace(aboutRunes[wordEnd]) {
			wordEnd++
		}

		visibleEnd := wordEnd
		if visibleEnd > count {
			visibleEnd = count
		}

		if count < wordEnd {
			for j := wordStart; j < visibleEnd; j++ {
				out[j] = scrambleRune(scrambleTick, j)
			}
		} else {
			ticksSinceComplete := scrambleTick - wordEnd
			if ticksSinceComplete < 0 {
				ticksSinceComplete = 0
			}
			wordLen := wordEnd - wordStart
			settled := settledChars(wordLen, ticksSinceComplete)
			for j := wordStart; j < visibleEnd; j++ {
				if j-wordStart < settled {
					out[j] = aboutRunes[j]
					continue
				}
				out[j] = scrambleRune(scrambleTick, j)
			}
		}

		i = visibleEnd
	}

	return string(out)
}

func AboutRuneCount() int {
	return len(aboutRunes)
}

func AboutSettleTicks() int {
	return settleDurationForWord(lastWordLength())
}

func RenderAbout(styles view.ThemeStyles, revealCount int, scrambleTick int, themeLabel string) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("━━━ About Me ━━━"))
	b.WriteString("\n")
	b.WriteString(styles.Content.Render(aboutVisible(revealCount, scrambleTick)))
	b.WriteString("\n")
	b.WriteString(styles.Help.Render(themeLabel + " • esc: back to menu"))

	return b.String()
}

func scrambleRune(scrambleTick int, index int) rune {
	if len(scrambleRunes) == 0 {
		return ' '
	}
	value := uint64(scrambleTick)*1103515245 + uint64(index)*12345 + 12345
	return scrambleRunes[int(value%uint64(len(scrambleRunes)))]
}

func isWhitespace(ch rune) bool {
	switch ch {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}

func settledChars(wordLen int, ticksSinceComplete int) int {
	if wordLen <= 0 {
		return 0
	}
	durationTicks := settleDurationForWord(wordLen)
	if durationTicks <= 0 {
		return 0
	}
	if ticksSinceComplete >= durationTicks {
		return wordLen
	}
	if ticksSinceComplete <= 0 {
		return 0
	}
	step := float64(durationTicks) / float64(wordLen)
	settled := int(math.Floor(float64(ticksSinceComplete) / step))
	if settled > wordLen {
		return wordLen
	}
	if settled < 0 {
		return 0
	}
	return settled
}

func settleDurationForWord(wordLen int) int {
	if wordLen <= 0 {
		return 0
	}
	if settleDurationTicks < wordLen {
		return wordLen
	}
	return settleDurationTicks
}

func lastWordLength() int {
	if len(aboutRunes) == 0 {
		return 0
	}
	i := len(aboutRunes) - 1
	for i >= 0 && isWhitespace(aboutRunes[i]) {
		i--
	}
	length := 0
	for i >= 0 && !isWhitespace(aboutRunes[i]) {
		length++
		i--
	}
	return length
}
