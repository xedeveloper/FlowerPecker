package screens

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/flowerpecker/internal/domain/models"
	"github.com/xedeveloper/flowerpecker/internal/domain/usecases"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type moduleGenerationDoneMsg struct{ err error }

type ModuleModel struct {
	width       int
	height      int
	input       textinput.Model
	spinner     spinner.Model
	config      models.ProjectConfig
	validErr    string
	generating  bool
	done        bool
	genErr      error
	moduleName  string
	createAnother bool
	quit        bool
}

func NewModuleModel(config models.ProjectConfig) ModuleModel {
	ti := textinput.New()
	ti.Placeholder = "auth"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	return ModuleModel{
		input:   ti,
		spinner: s,
		config:  config,
	}
}

func (m ModuleModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ModuleModel) Update(msg tea.Msg) (ModuleModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.done {
			switch msg.String() {
			case "c", "C":
				m.createAnother = true
			case "q", "Q", "ctrl+c":
				m.quit = true
				return m, tea.Quit
			}
			return m, nil
		}

		if !m.generating {
			switch msg.String() {
			case "enter":
				name := strings.TrimSpace(m.input.Value())
				if err := validateModuleName(name); err != "" {
					m.validErr = err
				} else {
					m.moduleName = name
					m.generating = true
					return m, tea.Batch(m.spinner.Tick, runModuleGeneration(m.config, name))
				}
				return m, nil
			case "esc":
				return m, NavigateBack()
			case "q", "ctrl+c":
				m.quit = true
				return m, tea.Quit
			}
		}

	case moduleGenerationDoneMsg:
		m.generating = false
		m.done = true
		m.genErr = msg.err
		return m, nil
	}

	if m.generating {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.validErr != "" {
		m.validErr = ""
	}
	return m, cmd
}

func runModuleGeneration(config models.ProjectConfig, moduleName string) tea.Cmd {
	return func() tea.Msg {
		moduleConfig := models.ModuleConfig{
			Name:            moduleName,
			ProjectPath:     filepath.Join(config.Path, config.Name),
			Architecture:    config.Architecture,
			StateManagement: config.StateManagement,
		}
		err := usecases.CreateModule(moduleConfig, func(step string) {})
		return moduleGenerationDoneMsg{err: err}
	}
}

func (m ModuleModel) ShouldCreateAnother() bool { return m.createAnother }
func (m ModuleModel) ShouldQuit() bool          { return m.quit }

func (m *ModuleModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m ModuleModel) View() string {
	header := renderScreenHeader("Create Module", m.width)

	var body string
	if m.generating {
		body = lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(m.spinner.View() + "  Generating module...")
	} else if m.done {
		body = m.renderDoneView()
	} else {
		body = m.renderInputView()
	}

	var footerHints string
	if m.done {
		footerHints = "[c] create another  [q] quit"
	} else {
		footerHints = "[enter] confirm  [esc] back  [q] quit"
	}

	footer := renderFooterHints(m.width, footerHints)
	return renderFullScreen(header, body, footer, m.width, m.height)
}

func (m ModuleModel) renderInputView() string {
	prompt := style.TitleStyle.Render("Enter module name:")
	note := style.SubtitleStyle.Render("e.g. auth, home, profile (use snake_case)")

	inputBox := style.InputStyle.Copy().Width(44).Render(m.input.View())

	var errLine string
	if m.validErr != "" {
		errLine = style.ErrorStyle.Render("⚠  " + m.validErr)
	}

	bodyParts := []string{prompt, "", note, "", inputBox}
	if errLine != "" {
		bodyParts = append(bodyParts, "", errLine)
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Left, bodyParts...))
}

func (m ModuleModel) renderDoneView() string {
	if m.genErr != nil {
		return lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(style.ErrorStyle.Render("✗ Error: " + m.genErr.Error()))
	}

	success := style.SuccessStyle.Render("✓ Module '" + m.moduleName + "' created successfully!")
	actions := lipgloss.JoinVertical(lipgloss.Left,
		style.ButtonPrimaryStyle.Render("[C] Create another module"),
		"",
		style.ButtonOutlineStyle.Render("[Q] Quit"),
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Center, success, "", actions))
}

func validateModuleName(name string) string {
	if name == "" {
		return "Module name cannot be empty"
	}
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9_]*$`, name)
	if !matched {
		return "Only lowercase letters, numbers, and underscores allowed"
	}
	return ""
}
