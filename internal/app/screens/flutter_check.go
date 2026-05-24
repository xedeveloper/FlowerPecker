package screens

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/flowerpecker/internal/domain/usecases"
	"github.com/xedeveloper/flowerpecker/internal/infrastructure/flutter"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type flutterCheckDoneMsg struct {
	result usecases.FlutterCheckResult
}

type FlutterCheckModel struct {
	width   int
	height  int
	spinner spinner.Model
	result  *usecases.FlutterCheckResult
	loading bool
	advance bool
	quit    bool
}

func NewFlutterCheckModel() FlutterCheckModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle
	return FlutterCheckModel{
		spinner: s,
		loading: true,
	}
}

func (m FlutterCheckModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		runFlutterCheck(),
	)
}

func runFlutterCheck() tea.Cmd {
	return func() tea.Msg {
		result := usecases.CheckFlutter()
		return flutterCheckDoneMsg{result: result}
	}
}

func (m FlutterCheckModel) Update(msg tea.Msg) (FlutterCheckModel, tea.Cmd) {
	switch msg := msg.(type) {
	case flutterCheckDoneMsg:
		m.loading = false
		m.result = &msg.result
		return m, nil

	case tea.KeyMsg:
		if !m.loading && m.result != nil {
			switch msg.String() {
			case "enter":
				if m.result.Available {
					m.advance = true
				}
			case "q", "ctrl+c":
				m.quit = true
			}
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m FlutterCheckModel) ShouldAdvance() bool { return m.advance }
func (m FlutterCheckModel) ShouldQuit() bool    { return m.quit }

func (m *FlutterCheckModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m FlutterCheckModel) View() string {
	header := renderScreenHeader("Environment Check", m.width)

	var body string
	if m.loading {
		body = lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(m.spinner.View() + "  Checking Flutter environment...")
	} else if m.result != nil {
		body = m.renderResults()
	}

	footer := renderFooterHints(m.width, "[enter] continue  [q] quit")
	return renderFullScreen(header, body, footer, m.width, m.height)
}

func (m FlutterCheckModel) renderResults() string {
	if !m.result.Available {
		errMsg := style.ErrorStyle.Render("✗ Flutter SDK not found")
		hint := style.SubtitleStyle.Render("Please install Flutter from https://flutter.dev/docs/get-started/install")
		return lipgloss.JoinVertical(lipgloss.Left, errMsg, "", hint)
	}

	lines := []string{}
	versionLine := style.SuccessStyle.Render("✓ Flutter SDK detected") + "  " +
		style.SubtitleStyle.Render("v"+m.result.Version)
	lines = append(lines, versionLine)

	for _, check := range m.result.Doctor.Checks {
		var line string
		if check.Status {
			line = style.SuccessStyle.Render("✓ " + check.Name)
		} else {
			line = style.ErrorStyle.Render("✗ " + check.Name)
		}
		lines = append(lines, line)
	}

	if len(m.result.Doctor.Devices) > 0 {
		lines = append(lines, "")
		lines = append(lines, style.TitleStyle.Render("Connected Devices:"))
		for _, d := range m.result.Doctor.Devices {
			lines = append(lines, "  "+style.NormalStyle.Render("• "+d))
		}
	}

	panel := style.PanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(panel)
}

func renderScreenHeader(title string, width int) string {
	titleRendered := style.TitleStyle.Copy().
		Width(width).
		Align(lipgloss.Center).
		Padding(1, 0).
		Render(title)
	divider := style.HairlineStyle.Copy().Width(width).Render("")
	return lipgloss.JoinVertical(lipgloss.Left, titleRendered, divider)
}

func renderFooterHints(width int, hints string) string {
	return style.FooterStyle.Copy().
		Width(width).
		Align(lipgloss.Center).
		Render(hints)
}

func renderFullScreen(header, body, footer string, width, height int) string {
	headerLines := countLines(header)
	footerLines := countLines(footer)
	bodyLines := height - headerLines - footerLines - 2
	if bodyLines < 1 {
		bodyLines = 1
	}

	bodyContainer := lipgloss.NewStyle().
		Width(width).
		Height(bodyLines).
		Align(lipgloss.Center, lipgloss.Center).
		Render(body)

	return lipgloss.JoinVertical(lipgloss.Left, header, bodyContainer, footer)
}

func countLines(s string) int {
	count := 1
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	return count
}

// Ensure flutter package is used
var _ = flutter.DoctorResult{}
