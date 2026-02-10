package pages

import (
	"math"
	"strings"

	"github.com/andatoshiki/termfolio/view"
)

var aboutContent = "Hey there, I'm Anda Toshiki -- call me kiki for short (like the protagonist from Kiki's Delivery Service by Hayao Miyazaki). I'm a Maho ShouJo (魔法少女) who loves anime, drinks monster, writes code, documents tutorials, takes photos, eats burgers, and stays up way too late. I'm into clean UI, useful tooling, and anything that makes dev life a little smoother.\n\nPronounced as AHN-dah TOH-shee-kee in case you find my name obscure to read out. I am currently an undergraduate student at ASU. I'm the founder of Toshiki Dev, my first organizational community focused on building small, functional utility components and webapps for real needs, both mine and others'. I maintain the community alongside my two cats 😺 (玉ちゃん & 桃ちゃん), who are definitely cute girls.\n"

var aboutRunes = []rune(aboutContent)
var scrambleRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:'\",.<>?/~")

const settleDurationTicks = 8

func aboutSettled(count int, scrambleTick int) bool {
	total := len(aboutRunes)
	if count > total {
		count = total
	}
	return count >= total && scrambleTick >= count+AboutSettleTicks()
}

func aboutStyled(styles view.ThemeStyles) string {
	var b strings.Builder

	normal := styles.Content
	bold := styles.Content.Copy().Bold(true)
	italic := styles.Content.Copy().Italic(true)
	link := styles.Accent.Copy().Underline(true)

	b.WriteString(normal.Render("Hey there, I'm "))
	b.WriteString(bold.Render("Anda Toshiki"))
	b.WriteString(normal.Render(" -- call me "))
	b.WriteString(italic.Render("kiki"))
	b.WriteString(normal.Render(" for short (like the protagonist from "))
	b.WriteString(view.ClickableLink(link.Render("Kiki's Delivery Service"), "https://en.wikipedia.org/wiki/Kiki%27s_Delivery_Service"))
	b.WriteString(normal.Render(" by Hayao Miyazaki). I'm a "))
	b.WriteString(italic.Render("Maho ShouJo"))
	b.WriteString(normal.Render(" (魔法少女) who loves anime, drinks monster, writes code, documents tutorials, takes photos, eats burgers, and stays up way too late. I'm into clean UI, useful tooling, and anything that makes dev life a little smoother. "))
	b.WriteString("\n\n")
	b.WriteString(normal.Render("Pronounced as "))
	b.WriteString(bold.Render("AHN-dah TOH-shee-kee"))
	b.WriteString(normal.Render(" in case you find my name obscure to read out. I am currently an undergraduate student at "))
	b.WriteString(view.ClickableLink(link.Render("ASU"), "https://asu.edu"))
	b.WriteString(normal.Render(". I'm the founder of "))
	b.WriteString(view.ClickableLink(link.Render("Toshiki Dev"), "https://github.com/toshikidev"))
	b.WriteString(normal.Render(", my first organizational community focused on building small, functional utility components and webapps for real needs, both mine and others'. I maintain the community alongside my two cats 😺 (玉ちゃん & 桃ちゃん), who are definitely cute girls."))

	return b.String()
}

func aboutVisible(count int, scrambleTick int) string {
	if count <= 0 {
		return ""
	}
	total := len(aboutRunes)
	if count > total {
		count = total
	}

	if aboutSettled(count, scrambleTick) {
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
	if aboutSettled(revealCount, scrambleTick) {
		b.WriteString(aboutStyled(styles))
	} else {
		b.WriteString(styles.Content.Render(aboutVisible(revealCount, scrambleTick)))
	}
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
