package screens

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/flowerpecker/internal/domain/models"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type stateMgmtItem struct {
	title, desc string
	state       models.StateManagement
}

func (i stateMgmtItem) Title() string       { return i.title }
func (i stateMgmtItem) Description() string { return i.desc }
func (i stateMgmtItem) FilterValue() string { return i.title }

type StateMgmtModel struct {
	width           int
	height          int
	list            list.Model
	advance         bool
	stateManagement models.StateManagement
}

func NewStateMgmtModel() StateMgmtModel {
	items := []list.Item{
		stateMgmtItem{
			title: "Flutter BLoC",
			desc:  "Business Logic Component — reactive, event-driven (recommended)",
			state: models.StateBLoC,
		},
		stateMgmtItem{
			title: "Riverpod",
			desc:  "Next generation Provider — compile-safe, testable",
			state: models.StateRiverpod,
		},
		stateMgmtItem{
			title: "Signals",
			desc:  "Reactive signals-based state management",
			state: models.StateSignals,
		},
		stateMgmtItem{
			title: "Provider",
			desc:  "Simple, lightweight InheritedWidget wrapper",
			state: models.StateProvider,
		},
		stateMgmtItem{
			title: "GetX",
			desc:  "All-in-one: state, routing & dependency injection",
			state: models.StateGetX,
		},
	}

	delegate := newStyledDelegate()
	l := list.New(items, delegate, 60, 14)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return StateMgmtModel{list: l}
}

func (m StateMgmtModel) Init() tea.Cmd { return nil }

func (m StateMgmtModel) Update(msg tea.Msg) (StateMgmtModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, ok := m.list.SelectedItem().(stateMgmtItem); ok {
				m.stateManagement = item.state
				m.advance = true
			}
			return m, nil
		case "esc":
			return m, NavigateBack()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m StateMgmtModel) ShouldAdvance() bool              { return m.advance }
func (m StateMgmtModel) Selected() models.StateManagement { return m.stateManagement }

func (m *StateMgmtModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.list.SetWidth(min(w-8, 70))
	m.list.SetHeight(min(h-10, 16))
}

func (m StateMgmtModel) View() string {
	header := renderScreenHeader("Select State Management", m.width)

	subtitle := style.SubtitleStyle.Render("Choose how to manage state in your Flutter app")
	listView := m.list.View()

	body := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Center, subtitle, "", listView))

	footer := renderFooterHints(m.width, "[↑↓] navigate  [enter] confirm  [esc] back  [q] quit")
	return renderFullScreen(header, body, footer, m.width, m.height)
}
