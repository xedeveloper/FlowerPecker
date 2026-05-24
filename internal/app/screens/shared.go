package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	style "github.com/xedeveloper/flowerpecker/pkg/style"
)

// NavigateBackMsg is sent by any screen when the user presses Esc,
// signalling the root AppModel to transition to the previous screen.
type NavigateBackMsg struct{}

func NavigateBack() tea.Cmd {
	return func() tea.Msg { return NavigateBackMsg{} }
}

type styledDelegate struct {
	list.DefaultDelegate
}

func newStyledDelegate() styledDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(style.ColorLink).
		BorderLeftForeground(style.ColorLink).
		Bold(true)

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(style.ColorBody).
		BorderLeftForeground(style.ColorLink)

	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Foreground(style.ColorCanvas)

	d.Styles.NormalDesc = d.Styles.NormalDesc.
		Foreground(style.ColorBody)

	d.ShowDescription = true

	return styledDelegate{d}
}

func renderPanelCentered(content string, width int) string {
	panel := style.PanelStyle.Render(content)
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(panel)
}
