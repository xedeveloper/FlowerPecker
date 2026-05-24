package style

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary    = lipgloss.Color("#000000")
	ColorOnPrimary  = lipgloss.Color("#BFDDF0")
	ColorInk        = lipgloss.Color("#000000")
	ColorInkSoft    = lipgloss.Color("#1a1a1a")
	ColorBody       = lipgloss.Color("#757575")
	ColorHairline   = lipgloss.Color("#e0e0e0")
	ColorCanvas     = lipgloss.Color("#ffffff")
	ColorCanvasSoft = lipgloss.Color("#f5f5f5")
	ColorLink       = lipgloss.Color("#057dbc")
	ColorError      = lipgloss.Color("#d32f2f")
	ColorSuccess    = lipgloss.Color("#2e7d32")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCanvas)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorBody)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorLink).
			Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(ColorCanvas)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorHairline)

	HairlineStyle = lipgloss.NewStyle().
			Border(lipgloss.Border{Bottom: "─"}, false, false, true, false).
			BorderForeground(ColorHairline)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(ColorLink)

	ButtonPrimaryStyle = lipgloss.NewStyle().
				Background(ColorPrimary).
				Foreground(ColorOnPrimary).
				Padding(0, 1)

	ButtonOutlineStyle = lipgloss.NewStyle().
				Foreground(ColorCanvas).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorHairline).
				Padding(0, 1)

	FooterStyle = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorOnPrimary).
			Padding(0, 2)

	InputStyle = lipgloss.NewStyle().
			Foreground(ColorCanvas).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorHairline).
			Padding(0, 1)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorHairline).
			Padding(1, 2)

	EyebrowStyle = lipgloss.NewStyle().
			Foreground(ColorCanvas).
			Bold(true)
)
