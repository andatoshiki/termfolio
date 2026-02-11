package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"

	"github.com/andatoshiki/termfolio/counter"
	"github.com/andatoshiki/termfolio/pages"
	"github.com/andatoshiki/termfolio/view"
)

type tickMsg time.Time

type page int

const (
	splashPage page = iota
	menuPage
	aboutPage
	projectsPage
	educationPage
	contactPage
	privacyPage
	feedPage
)

type model struct {
	currentPage     page
	splashReveal    int
	splashBlinkStep int
	menuCursor      int
	projectCursor   int
	eduCursor       int
	aboutReveal     int
	aboutScramble   int
	visitorCount    int
	remoteIP        string
	trackingEnabled bool
	privacyCursor   int
	feedItems       []pages.FeedItem
	feedCursor      int
	feedOffset      int
	feedLoading     bool
	feedError       string
	feedFetchedAt   time.Time
	counterStore    *counter.Store
	width           int
	height          int
	logoSweepIndex  int
	themeIndex      int
	styles          view.ThemeStyles
}

func initialModel() model {
	initialPalette := view.ThemeAt(0)
	return model{
		currentPage:     splashPage,
		splashReveal:    0,
		splashBlinkStep: 0,
		menuCursor:      0,
		projectCursor:   0,
		eduCursor:       0,
		aboutReveal:     0,
		aboutScramble:   0,
		visitorCount:    0,
		remoteIP:        "",
		trackingEnabled: false,
		privacyCursor:   0,
		feedItems:       nil,
		feedCursor:      0,
		feedOffset:      0,
		feedLoading:     false,
		feedError:       "",
		feedFetchedAt:   time.Time{},
		counterStore:    nil,
		width:           80,
		height:          24,
		logoSweepIndex:  0,
		themeIndex:      0,
		styles:          view.NewThemeStyles(initialPalette),
	}
}

func NewModel() tea.Model {
	return initialModel()
}

func NewModelWithVisitorCount(visitorCount int) tea.Model {
	m := initialModel()
	m.visitorCount = visitorCount
	return m
}

func NewModelWithCounter(store *counter.Store, visitorCount int, remoteIP string, trackingEnabled bool) tea.Model {
	m := initialModel()
	m.counterStore = store
	m.visitorCount = visitorCount
	m.remoteIP = remoteIP
	m.trackingEnabled = trackingEnabled
	return m
}

func (m model) Init() tea.Cmd {
	return splashTickCmd()
}

// Controls
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.currentPage == splashPage {
			m.splashBlinkStep++
			if m.splashBlinkStep >= 1_000_000 {
				m.splashBlinkStep = 0
			}
			logoTotal := pages.SplashLogoRuneCount()
			total := pages.SplashRuneCount()
			if m.splashReveal < logoTotal {
				m.splashReveal += splashLogoRevealStep
				if m.splashReveal > logoTotal {
					m.splashReveal = logoTotal
				}
			} else if m.splashReveal < total {
				if m.splashBlinkStep%splashTextRevealTickDivisor == 0 {
					m.splashReveal++
				}
			}
			return m, splashTickCmd()
		}
		if m.currentPage == menuPage {
			m.logoSweepIndex++
			return m, tickCmd()
		}
		if m.currentPage == aboutPage {
			if m.aboutReveal < pages.AboutRuneCount() {
				m.aboutReveal++
				m.aboutScramble++
				return m, typewriterTickCmd()
			}
			if m.aboutScramble < m.aboutReveal+pages.AboutSettleTicks() {
				m.aboutScramble++
				return m, typewriterTickCmd()
			}
		}
		return m, nil

	case feedMsg:
		m.feedLoading = false
		if msg.err != nil {
			m.feedError = msg.err.Error()
			return m, nil
		}
		m.feedError = ""
		m.feedItems = msg.items
		m.feedFetchedAt = time.Now()
		if len(m.feedItems) == 0 {
			m.feedCursor = 0
			m.feedOffset = 0
			return m, nil
		}
		if m.feedCursor >= len(m.feedItems) {
			m.feedCursor = 0
		}
		m = m.adjustFeedWindow()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.currentPage == menuPage || m.currentPage == splashPage {
				return m, tea.Quit
			}
			m.currentPage = menuPage
			return m, tickCmd()

		case "esc", "backspace":
			if m.currentPage == splashPage {
				return m, nil
			}
			if m.currentPage != menuPage {
				m.currentPage = menuPage
			}
			return m, tickCmd()

		case "up", "k":
			switch m.currentPage {
			case menuPage:
				if m.menuCursor > 0 {
					m.menuCursor--
				}
			case projectsPage:
				if m.projectCursor > 0 {
					m.projectCursor--
				}
			case educationPage:
				if m.eduCursor > 0 {
					m.eduCursor--
				}
			case privacyPage:
				if m.privacyCursor > 0 {
					m.privacyCursor--
				}
			case feedPage:
				if m.feedCursor > 0 {
					m.feedCursor--
				}
				m = m.adjustFeedWindow()
			}
			return m, nil

		case "down", "j":
			switch m.currentPage {
			case menuPage:
				if m.menuCursor < len(pages.MenuItems())-1 {
					m.menuCursor++
				}
			case projectsPage:
				if m.projectCursor < len(pages.Projects())-1 {
					m.projectCursor++
				}
			case educationPage:
				if m.eduCursor < len(pages.Educations())-1 {
					m.eduCursor++
				}
			case privacyPage:
				if m.privacyCursor < 1 {
					m.privacyCursor++
				}
			case feedPage:
				if m.feedCursor < len(m.feedItems)-1 {
					m.feedCursor++
				}
				m = m.adjustFeedWindow()
			}
			return m, nil

		case "enter", " ":
			if m.currentPage == splashPage {
				m.currentPage = menuPage
				m.logoSweepIndex = 0
				return m, tickCmd()
			}
			if m.currentPage == menuPage {
				switch m.menuCursor {
				case 0:
					m.currentPage = aboutPage
					m.aboutReveal = 0
					m.aboutScramble = 0
					return m, typewriterTickCmd()
				case 1:
					m.currentPage = projectsPage
					return m, nil
				case 2:
					m.currentPage = educationPage
					return m, nil
				case 3:
					m.currentPage = contactPage
					return m, nil
				case 4:
					m.currentPage = feedPage
					m.feedCursor = 0
					m.feedOffset = 0
					if shouldFetchFeed(m) {
						m.feedLoading = true
						m.feedError = ""
						return m, fetchFeedCmd()
					}
					return m, nil
				case 5:
					m.currentPage = privacyPage
					m.privacyCursor = 0
					return m, nil
				}
			}
			if m.currentPage == privacyPage {
				var cmd tea.Cmd
				m, cmd = m.setTracking(m.privacyCursor == 0)
				m.currentPage = menuPage
				return m, cmd
			}
			return m, nil

		case "t", "T":
			m.themeIndex = view.NextThemeIndex(m.themeIndex)
			m.styles = view.NewThemeStyles(view.ThemeAt(m.themeIndex))
			return m, nil

		}
	}
	return m, nil
}

func (m model) themeLabel() string {
	name := view.ThemeAt(m.themeIndex).Name
	if name == "" {
		return "t: change theme"
	}
	return "t: theme (" + name + ")"
}

func tickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

const splashTick = 40 * time.Millisecond

const splashLogoRevealStep = 3

const splashTextRevealTickDivisor = 1

const splashRainFrameDivisor = 1

const splashRainLeadPadding = 12

const typewriterTick = 40 * time.Millisecond

const feedPageSize = 8

func splashTickCmd() tea.Cmd {
	return tea.Tick(splashTick, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func typewriterTickCmd() tea.Cmd {
	return tea.Tick(typewriterTick, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) View() string {
	boxWidth := min(m.width-4, 70)
	var content string

	switch m.currentPage {
	case splashPage:
		content = pages.RenderSplash(m.styles, m.splashReveal, m.splashBlinkStep, boxWidth)
	case menuPage:
		content = pages.RenderMenu(m.styles, m.menuCursor, m.logoSweepIndex, m.themeLabel(), m.visitorCount, boxWidth)
	case aboutPage:
		content = pages.RenderAbout(m.styles, m.aboutReveal, m.aboutScramble, m.themeLabel(), boxWidth)
	case projectsPage:
		content = pages.RenderProjects(m.styles, m.projectCursor, m.themeLabel())
	case educationPage:
		content = pages.RenderEducation(m.styles, m.eduCursor, m.themeLabel())
	case contactPage:
		content = pages.RenderContact(m.styles, m.themeLabel())
	case privacyPage:
		content = pages.RenderPrivacy(m.styles, m.privacyCursor, m.trackingEnabled, m.counterStore != nil && m.remoteIP != "", m.themeLabel())
	case feedPage:
		content = pages.RenderFeed(m.styles, m.feedItems, m.feedCursor, m.feedOffset, feedPageSize, m.feedLoading, m.feedError, m.themeLabel())
	}

	boxedContent := lipgloss.NewStyle().
		Padding(1, 2).
		Width(boxWidth).
		Render(content)

	if m.currentPage == splashPage {
		left := max(0, (m.width-lipgloss.Width(boxedContent))/2)
		top := max(0, (m.height-lipgloss.Height(boxedContent))/2)
		return renderSplashWithRain(m.styles, m.splashBlinkStep, m.width, m.height, boxedContent, left, top)
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		boxedContent)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) setTracking(enabled bool) (model, tea.Cmd) {
	if m.counterStore == nil || m.remoteIP == "" {
		return m, nil
	}
	if enabled == m.trackingEnabled {
		return m, nil
	}
	count, err := m.counterStore.SetOptOut(m.remoteIP, !enabled)
	if err != nil {
		return m, nil
	}
	m.trackingEnabled = enabled
	m.visitorCount = count
	return m, nil
}

func (m model) adjustFeedWindow() model {
	if feedPageSize <= 0 || len(m.feedItems) == 0 {
		m.feedOffset = 0
		return m
	}
	if m.feedCursor < m.feedOffset {
		m.feedOffset = m.feedCursor
	} else if m.feedCursor >= m.feedOffset+feedPageSize {
		m.feedOffset = m.feedCursor - feedPageSize + 1
	}
	if m.feedOffset < 0 {
		m.feedOffset = 0
	}
	maxOffset := max(0, len(m.feedItems)-feedPageSize)
	if m.feedOffset > maxOffset {
		m.feedOffset = maxOffset
	}
	return m
}

func renderSplashWithRain(styles view.ThemeStyles, frame int, width int, height int, overlay string, left int, top int) string {
	if width <= 0 || height <= 0 {
		return overlay
	}

	rainFrame := frame
	if rainFrame < 0 {
		rainFrame = 0
	}
	rainFrame /= splashRainFrameDivisor

	tokens := buildSplashRainTokens(styles)
	seeds, speeds, trails := buildSplashRainColumns(width)
	overlayLines := strings.Split(overlay, "\n")
	contentLeft := max(0, left)
	if contentLeft > width {
		contentLeft = width
	}

	var out strings.Builder
	for y := 0; y < height; y++ {
		if y > 0 {
			out.WriteByte('\n')
		}

		overlayRow := y - top
		if overlayRow < 0 || overlayRow >= len(overlayLines) {
			writeSplashRainRange(&out, tokens, rainFrame, width, height, y, 0, width, false, seeds, speeds, trails)
			continue
		}

		overlayLine := overlayLines[overlayRow]
		plainOverlayLine := xansi.Strip(overlayLine)
		spanStart, spanEnd, hasVisible := nonSpaceSpan(plainOverlayLine)
		if !hasVisible {
			writeSplashRainRange(&out, tokens, rainFrame, width, height, y, 0, width, true, seeds, speeds, trails)
			continue
		}

		lineWidth := xansi.StringWidth(plainOverlayLine)
		if spanStart < 0 {
			spanStart = 0
		}
		if spanEnd > lineWidth {
			spanEnd = lineWidth
		}

		overlayStartCol := contentLeft + spanStart
		overlayEndCol := contentLeft + spanEnd
		if overlayStartCol < 0 {
			overlayStartCol = 0
		}
		if overlayEndCol > width {
			overlayEndCol = width
		}

		if overlayStartCol > 0 {
			writeSplashRainRange(&out, tokens, rainFrame, width, height, y, 0, overlayStartCol, true, seeds, speeds, trails)
		}

		if overlayEndCol > overlayStartCol {
			out.WriteString(xansi.Cut(overlayLine, spanStart, spanEnd))
		}

		if overlayEndCol < width {
			writeSplashRainRange(&out, tokens, rainFrame, width, height, y, overlayEndCol, width, true, seeds, speeds, trails)
		}
	}

	return out.String()
}

type splashRainTokens struct {
	headZero  string
	headOne   string
	trailZero string
	trailOne  string
	dimZero   string
	dimOne    string
}

func buildSplashRainTokens(styles view.ThemeStyles) splashRainTokens {
	headStyle := styles.Accent.Copy().Bold(false)
	trailStyle := styles.Accent.Copy().Bold(false).Faint(true)
	dimStyle := styles.Subtle.Copy().Bold(false).Faint(true)

	return splashRainTokens{
		headZero:  headStyle.Render("0"),
		headOne:   headStyle.Render("1"),
		trailZero: trailStyle.Render("0"),
		trailOne:  trailStyle.Render("1"),
		dimZero:   dimStyle.Render("0"),
		dimOne:    dimStyle.Render("1"),
	}
}

func buildSplashRainColumns(width int) ([]int, []int, []int) {
	seeds := make([]int, width)
	speeds := make([]int, width)
	trails := make([]int, width)

	for x := 0; x < width; x++ {
		seed := splashRainHash((x + 1) * 7919)
		seeds[x] = seed
		speeds[x] = 1 + seed%3
		trails[x] = 4 + seed%6
	}

	return seeds, speeds, trails
}

func writeSplashRainRange(
	out *strings.Builder,
	tokens splashRainTokens,
	rainFrame int,
	width int,
	height int,
	y int,
	start int,
	end int,
	soft bool,
	seeds []int,
	speeds []int,
	trails []int,
) {
	if start < 0 {
		start = 0
	}
	if end > width {
		end = width
	}
	if start >= end {
		return
	}

	for x := start; x < end; x++ {
		seed := seeds[x]
		speed := speeds[x]
		trail := trails[x]
		headRange := height + trail + splashRainLeadPadding
		head := (rainFrame/speed + seed) % headRange

		switch {
		case head < height && y <= head && y > head-trail:
			one := splashBinaryOne(rainFrame, x, y, seed)
			if soft {
				if one {
					out.WriteString(tokens.dimOne)
				} else {
					out.WriteString(tokens.dimZero)
				}
				continue
			}

			dist := head - y
			switch {
			case dist == 0:
				if one {
					out.WriteString(tokens.headOne)
				} else {
					out.WriteString(tokens.headZero)
				}
			case dist <= trail/2:
				if one {
					out.WriteString(tokens.trailOne)
				} else {
					out.WriteString(tokens.trailZero)
				}
			default:
				if one {
					out.WriteString(tokens.dimOne)
				} else {
					out.WriteString(tokens.dimZero)
				}
			}
		case splashRainHash(rainFrame*29+x*97+y*53+seed)%41 == 0:
			if splashBinaryOne(rainFrame/2, x, y, seed+17) {
				out.WriteString(tokens.dimOne)
			} else {
				out.WriteString(tokens.dimZero)
			}
		default:
			out.WriteByte(' ')
		}
	}
}

func nonSpaceSpan(s string) (start int, end int, ok bool) {
	runes := []rune(s)
	start = 0
	for start < len(runes) && runes[start] == ' ' {
		start++
	}
	if start >= len(runes) {
		return 0, 0, false
	}

	last := len(runes) - 1
	for last >= 0 && runes[last] == ' ' {
		last--
	}

	return start, last + 1, true
}

func splashRainHash(n int) int {
	n ^= n << 13
	n ^= n >> 17
	n ^= n << 5
	if n < 0 {
		n = -n
	}
	return n
}

func splashBinaryOne(frame int, x int, y int, seed int) bool {
	return splashRainHash(frame*131+x*17+y*23+seed)%2 == 0
}
