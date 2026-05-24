package screens

import (
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type SplashModel struct {
	width      int
	height     int
	sparrowArt string
	logoArt    string
	ready      bool
}

type splashTickMsg struct{}

func NewSplashModel() SplashModel {
	sparrow := readAsset("assets/ascii/sparrow.txt")
	logo := readAsset("assets/ascii/logo.txt")
	return SplashModel{
		sparrowArt: sparrow,
		logoArt:    logo,
	}
}

func readAsset(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func (m SplashModel) Init() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return splashTickMsg{}
	})
}

func (m SplashModel) Update(msg tea.Msg) (SplashModel, tea.Cmd) {
	switch msg.(type) {
	case splashTickMsg:
		m.ready = true
	case tea.KeyMsg:
		m.ready = true
	}
	return m, nil
}

func (m SplashModel) ShouldAdvance() bool {
	return m.ready
}

func (m *SplashModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m SplashModel) View() string {
	if m.width == 0 {
		return ""
	}

	// The splash always renders white text on a black canvas so the ASCII art
	// is legible regardless of whether the user's terminal theme is light or dark.
	splashBg := style.ColorPrimary   // #000000
	splashFg := style.ColorOnPrimary // #ffffff

	screenStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(splashBg).
		Align(lipgloss.Center, lipgloss.Center)

	logoBlock := lipgloss.NewStyle().
		Bold(true).
		Foreground(splashFg).
		Background(splashBg).
		Render(m.logoArt)

	sparrowBlock := lipgloss.NewStyle().
		Foreground(splashFg).
		Background(splashBg).
		Render(m.sparrowArt)

	tagline := lipgloss.NewStyle().
		Foreground(style.ColorHairline). // soft grey — readable on black
		Background(splashBg).
		Render("Flutter project scaffolding tool")

	hint := lipgloss.NewStyle().
		Foreground(style.ColorHairline).
		Background(splashBg).
		Render("Press any key to continue...")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		sparrowBlock,
		"",
		logoBlock,
		"",
		tagline,
		"",
		hint,
	)

	verticalPad := (m.height - strings.Count(content, "\n") - 4) / 2
	if verticalPad < 0 {
		verticalPad = 0
	}

	padding := strings.Repeat("\n", verticalPad)
	return screenStyle.Render(padding + content)
}
