package screens

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/flowerpecker/internal/domain/models"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type architectureItem struct {
	title, desc string
	arch        models.Architecture
}

func (i architectureItem) Title() string       { return i.title }
func (i architectureItem) Description() string { return i.desc }
func (i architectureItem) FilterValue() string { return i.title }

type ArchitectureModel struct {
	width        int
	height       int
	list         list.Model
	advance      bool
	architecture models.Architecture
}

func NewArchitectureModel() ArchitectureModel {
	items := []list.Item{
		architectureItem{
			title: "CLEAN Architecture",
			desc:  "Data, Domain & Presentation layers (recommended)",
			arch:  models.ArchClean,
		},
		architectureItem{
			title: "MVC",
			desc:  "Model, View & Controller pattern",
			arch:  models.ArchMVC,
		},
		architectureItem{
			title: "MVVM",
			desc:  "Model, View & ViewModel pattern",
			arch:  models.ArchMVVM,
		},
	}

	delegate := newStyledDelegate()
	l := list.New(items, delegate, 60, 10)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return ArchitectureModel{list: l}
}

func (m ArchitectureModel) Init() tea.Cmd { return nil }

func (m ArchitectureModel) Update(msg tea.Msg) (ArchitectureModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, ok := m.list.SelectedItem().(architectureItem); ok {
				m.architecture = item.arch
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

func (m ArchitectureModel) ShouldAdvance() bool        { return m.advance }
func (m ArchitectureModel) Selected() models.Architecture { return m.architecture }

func (m *ArchitectureModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.list.SetWidth(min(w-8, 70))
	m.list.SetHeight(min(h-10, 12))
}

func (m ArchitectureModel) View() string {
	header := renderScreenHeader("Select Architecture", m.width)

	subtitle := style.SubtitleStyle.Render("Choose the architecture pattern for your Flutter project")
	listView := m.list.View()

	body := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Center, subtitle, "", listView))

	footer := renderFooterHints(m.width, "[↑↓] navigate  [enter] confirm  [esc] back  [q] quit")
	return renderFullScreen(header, body, footer, m.width, m.height)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
