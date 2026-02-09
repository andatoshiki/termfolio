package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/andatoshiki/termfolio/counter"
	"github.com/andatoshiki/termfolio/pages"
	"github.com/andatoshiki/termfolio/view"
)

type tickMsg time.Time

type page int

const (
	menuPage page = iota
	aboutPage
	projectsPage
	experiencePage
	contactPage
	privacyPage
	feedPage
)

type model struct {
	currentPage     page
	menuCursor      int
	projectCursor   int
	expCursor       int
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
		currentPage:     menuPage,
		menuCursor:      0,
		projectCursor:   0,
		expCursor:       0,
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
	return tickCmd()
}

// Controls
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
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
			if m.currentPage == menuPage {
				return m, tea.Quit
			}
			m.currentPage = menuPage
			return m, tickCmd()

		case "esc", "backspace":
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
			case experiencePage:
				if m.expCursor > 0 {
					m.expCursor--
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
			case experiencePage:
				if m.expCursor < len(pages.Experiences())-1 {
					m.expCursor++
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
					m.currentPage = experiencePage
					return m, nil
				case 3:
					m.currentPage = contactPage
					return m, nil
				case 4:
					m.currentPage = privacyPage
					m.privacyCursor = 0
					return m, nil
				case 5:
					m.currentPage = feedPage
					m.feedCursor = 0
					m.feedOffset = 0
					if shouldFetchFeed(m) {
						m.feedLoading = true
						m.feedError = ""
						return m, fetchFeedCmd()
					}
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

const typewriterTick = 40 * time.Millisecond

const feedPageSize = 8

func typewriterTickCmd() tea.Cmd {
	return tea.Tick(typewriterTick, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) View() string {
	var content string

	switch m.currentPage {
	case menuPage:
		content = pages.RenderMenu(m.styles, m.menuCursor, m.logoSweepIndex, m.themeLabel(), m.visitorCount)
	case aboutPage:
		content = pages.RenderAbout(m.styles, m.aboutReveal, m.aboutScramble, m.themeLabel())
	case projectsPage:
		content = pages.RenderProjects(m.styles, m.projectCursor, m.themeLabel())
	case experiencePage:
		content = pages.RenderExperience(m.styles, m.expCursor, m.themeLabel())
	case contactPage:
		content = pages.RenderContact(m.styles, m.themeLabel())
	case privacyPage:
		content = pages.RenderPrivacy(m.styles, m.privacyCursor, m.trackingEnabled, m.counterStore != nil && m.remoteIP != "", m.themeLabel())
	case feedPage:
		content = pages.RenderFeed(m.styles, m.feedItems, m.feedCursor, m.feedOffset, feedPageSize, m.feedLoading, m.feedError, m.themeLabel())
	}

	boxWidth := min(m.width-4, 70)
	boxedContent := lipgloss.NewStyle().
		Padding(1, 2).
		Width(boxWidth).
		Render(content)

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
