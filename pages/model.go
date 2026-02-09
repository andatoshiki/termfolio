package pages

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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
)

type model struct {
	currentPage    page
	menuCursor     int
	projectCursor  int
	expCursor      int
	aboutReveal    int
	aboutScramble  int
	width          int
	height         int
	logoSweepIndex int
	themeIndex     int
	styles         view.ThemeStyles
}

func initialModel() model {
	initialPalette := view.ThemeAt(0)
	return model{
		currentPage:    menuPage,
		menuCursor:     0,
		projectCursor:  0,
		expCursor:      0,
		aboutReveal:    0,
		aboutScramble:  0,
		width:          80,
		height:         24,
		logoSweepIndex: 0,
		themeIndex:     0,
		styles:         view.NewThemeStyles(initialPalette),
	}
}

func NewModel() tea.Model {
	return initialModel()
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
			if m.aboutReveal < AboutRuneCount() {
				m.aboutReveal++
				m.aboutScramble++
				return m, typewriterTickCmd()
			}
			if m.aboutScramble < m.aboutReveal+AboutSettleTicks() {
				m.aboutScramble++
				return m, typewriterTickCmd()
			}
		}
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
			}
			return m, nil

		case "down", "j":
			switch m.currentPage {
			case menuPage:
				if m.menuCursor < len(menuItems)-1 {
					m.menuCursor++
				}
			case projectsPage:
				if m.projectCursor < len(projects)-1 {
					m.projectCursor++
				}
			case experiencePage:
				if m.expCursor < len(experiences)-1 {
					m.expCursor++
				}
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
				case 2:
					m.currentPage = experiencePage
				case 3:
					m.currentPage = contactPage
				}
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

func typewriterTickCmd() tea.Cmd {
	return tea.Tick(typewriterTick, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

const typewriterTick = 40 * time.Millisecond

func (m model) View() string {
	themeLabel := m.themeLabel()
	var content string

	switch m.currentPage {
	case menuPage:
		content = RenderMenu(m.styles, m.menuCursor, m.logoSweepIndex, themeLabel)
	case aboutPage:
		content = RenderAbout(m.styles, m.aboutReveal, m.aboutScramble, themeLabel)
	case projectsPage:
		content = RenderProjects(m.styles, m.projectCursor, themeLabel)
	case experiencePage:
		content = RenderExperience(m.styles, m.expCursor, themeLabel)
	case contactPage:
		content = RenderContact(m.styles, themeLabel)
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
