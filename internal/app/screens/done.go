package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/flowerpecker/internal/domain/models"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type DoneModel struct {
	width      int
	height     int
	config     models.ProjectConfig
	createModule bool
	quit       bool
}

func NewDoneModel(config models.ProjectConfig) DoneModel {
	return DoneModel{config: config}
}

func (m DoneModel) Init() tea.Cmd { return nil }

func (m DoneModel) Update(msg tea.Msg) (DoneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c", "C":
			m.createModule = true
		case "q", "Q", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DoneModel) ShouldCreateModule() bool { return m.createModule }
func (m DoneModel) ShouldQuit() bool         { return m.quit }

func (m *DoneModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m DoneModel) View() string {
	header := renderScreenHeader("Project Created", m.width)

	successIcon := lipgloss.NewStyle().
		Foreground(style.ColorSuccess).
		Bold(true).
		Render("✓")

	title := style.TitleStyle.Render("Project created successfully!")
	projectInfo := style.SubtitleStyle.Render(fmt.Sprintf("Project: %s", m.config.Name))
	pathInfo := style.SubtitleStyle.Render(fmt.Sprintf("Location: %s/%s", m.config.Path, m.config.Name))
	archInfo := style.SubtitleStyle.Render(fmt.Sprintf("Architecture: %s  |  State: %s",
		m.config.Architecture.String(), m.config.StateManagement.String()))

	tree := buildProjectTreeView(m.config)
	treePanel := style.PanelStyle.Render(tree)

	actions := lipgloss.JoinVertical(lipgloss.Left,
		style.ButtonPrimaryStyle.Render("[C] Create a new module"),
		"",
		style.ButtonOutlineStyle.Render("[Q] Quit"),
	)

	body := lipgloss.JoinVertical(lipgloss.Center,
		successIcon,
		"",
		title,
		"",
		projectInfo,
		pathInfo,
		archInfo,
		"",
		treePanel,
		"",
		actions,
	)

	bodyRendered := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(body)

	footer := renderFooterHints(m.width, "[c] create module  [q] quit")
	return renderFullScreen(header, bodyRendered, footer, m.width, m.height)
}

func buildProjectTreeView(config models.ProjectConfig) string {
	switch config.Architecture {
	case models.ArchClean:
		return fmt.Sprintf(`%s/
├── lib/
│   ├── core/
│   │   ├── constants/
│   │   ├── errors/
│   │   ├── theme/
│   │   ├── utils/
│   │   └── widgets/
│   ├── features/
│   └── main.dart
└── pubspec.yaml`, config.Name)
	case models.ArchMVC:
		return fmt.Sprintf(`%s/
├── lib/
│   ├── core/
│   ├── models/
│   ├── views/
│   ├── controllers/
│   └── main.dart
└── pubspec.yaml`, config.Name)
	default:
		return fmt.Sprintf(`%s/
├── lib/
│   ├── core/
│   ├── models/
│   ├── views/
│   ├── viewmodels/
│   └── main.dart
└── pubspec.yaml`, config.Name)
	}
}
