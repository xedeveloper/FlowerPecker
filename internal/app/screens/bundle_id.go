package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type BundleIDModel struct {
	width    int
	height   int
	input    textinput.Model
	validErr string
	advance  bool
	bundleID string
}

func NewBundleIDModel() BundleIDModel {
	ti := textinput.New()
	ti.Placeholder = "com.company.appname"
	ti.Focus()
	ti.CharLimit = 128
	ti.Width = 40

	return BundleIDModel{input: ti}
}

func (m BundleIDModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m BundleIDModel) Update(msg tea.Msg) (BundleIDModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			id := strings.TrimSpace(m.input.Value())
			if err := validateBundleID(id); err != "" {
				m.validErr = err
			} else {
				m.bundleID = id
				m.advance = true
			}
			return m, nil
		case "esc":
			return m, NavigateBack()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.validErr != "" {
		m.validErr = ""
	}
	return m, cmd
}

func (m BundleIDModel) ShouldAdvance() bool { return m.advance }
func (m BundleIDModel) BundleID() string    { return m.bundleID }

func (m *BundleIDModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m BundleIDModel) View() string {
	header := renderScreenHeader("Project Configuration  ·  Step 2 of 2", m.width)

	prompt := style.TitleStyle.Render("Enter bundle identifier:")
	note := style.SubtitleStyle.Render("e.g. com.company.appname")

	inputBox := style.InputStyle.Copy().Width(44).Render(m.input.View())

	var errLine string
	if m.validErr != "" {
		errLine = style.ErrorStyle.Render("⚠  " + m.validErr)
	}

	bodyParts := []string{prompt, "", note, "", inputBox}
	if errLine != "" {
		bodyParts = append(bodyParts, "", errLine)
	}

	body := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Left, bodyParts...))

	footer := renderFooterHints(m.width, "[enter] confirm  [esc] back  [q] quit")
	return renderFullScreen(header, body, footer, m.width, m.height)
}

func validateBundleID(id string) string {
	if id == "" {
		return "Bundle ID cannot be empty"
	}
	if !strings.Contains(id, ".") {
		return "Bundle ID must contain at least one dot (e.g. com.company.app)"
	}
	if strings.Contains(id, " ") {
		return "Bundle ID cannot contain spaces"
	}
	return ""
}
