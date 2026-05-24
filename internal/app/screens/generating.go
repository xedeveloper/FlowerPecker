package screens

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/flowerpecker/internal/domain/models"
	"github.com/xedeveloper/flowerpecker/internal/domain/usecases"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type generationStepMsg struct{ step string }
type generationDoneMsg struct{ err error }

type GeneratingModel struct {
	width         int
	height        int
	spinner       spinner.Model
	config        models.ProjectConfig
	currentStep   string
	completedSteps []string
	done          bool
	err           error
	advance       bool
}

func NewGeneratingModel(config models.ProjectConfig) GeneratingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle
	return GeneratingModel{
		spinner: s,
		config:  config,
	}
}

func (m GeneratingModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		runProjectGeneration(m.config),
	)
}

func runProjectGeneration(config models.ProjectConfig) tea.Cmd {
	return func() tea.Msg {
		err := usecases.CreateProject(config, func(step string) {})
		return generationDoneMsg{err: err}
	}
}

func (m GeneratingModel) Update(msg tea.Msg) (GeneratingModel, tea.Cmd) {
	switch msg := msg.(type) {
	case generationStepMsg:
		if m.currentStep != "" {
			m.completedSteps = append(m.completedSteps, m.currentStep)
		}
		m.currentStep = msg.step
		return m, nil

	case generationDoneMsg:
		m.done = true
		m.err = msg.err
		if m.currentStep != "" {
			m.completedSteps = append(m.completedSteps, m.currentStep)
		}
		m.currentStep = ""
		return m, nil

	case tea.KeyMsg:
		if m.done {
			m.advance = true
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m GeneratingModel) ShouldAdvance() bool { return m.advance }

func (m *GeneratingModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m GeneratingModel) View() string {
	header := renderScreenHeader("Generating Project", m.width)

	lines := []string{}
	for _, step := range m.completedSteps {
		lines = append(lines, style.SuccessStyle.Render("✓  "+step))
	}

	if m.currentStep != "" {
		lines = append(lines, m.spinner.View()+"  "+style.NormalStyle.Render(m.currentStep))
	}

	if m.done {
		if m.err != nil {
			lines = append(lines, "", style.ErrorStyle.Render("✗ Error: "+m.err.Error()))
		} else {
			lines = append(lines, "", style.SuccessStyle.Render("✓ Project generated successfully!"))
			lines = append(lines, "", style.SubtitleStyle.Render("Press any key to continue..."))
		}
	}

	body := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))

	footer := renderFooterHints(m.width, "Please wait...")
	if m.done {
		footer = renderFooterHints(m.width, "[enter] continue  [q] quit")
	}
	return renderFullScreen(header, body, footer, m.width, m.height)
}
