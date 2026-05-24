package screens

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type ProjectNameModel struct {
	width      int
	height     int
	input      textinput.Model
	validErr   string
	advance    bool
	projectName string
}

func NewProjectNameModel() ProjectNameModel {
	ti := textinput.New()
	ti.Placeholder = "my_awesome_app"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40

	return ProjectNameModel{input: ti}
}

func (m ProjectNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ProjectNameModel) Update(msg tea.Msg) (ProjectNameModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			name := strings.TrimSpace(m.input.Value())
			if err := validateProjectName(name); err != "" {
				m.validErr = err
			} else {
				m.projectName = name
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

func (m ProjectNameModel) ShouldAdvance() bool  { return m.advance }
func (m ProjectNameModel) ProjectName() string   { return m.projectName }

func (m *ProjectNameModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m ProjectNameModel) View() string {
	header := renderScreenHeader("Project Configuration  ·  Step 1 of 2", m.width)

	prompt := style.TitleStyle.Render("Enter your Flutter project name:")
	note := style.SubtitleStyle.Render("Use snake_case (e.g. my_awesome_app)")

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

func validateProjectName(name string) string {
	if name == "" {
		return "Project name cannot be empty"
	}
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9_]*$`, name)
	if !matched {
		return "Only lowercase letters, numbers, and underscores allowed"
	}
	if strings.Contains(name, "__") {
		return "Cannot have consecutive underscores"
	}
	return ""
}
