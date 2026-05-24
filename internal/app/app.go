package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xedeveloper/flowerpecker/internal/app/screens"
	"github.com/xedeveloper/flowerpecker/internal/domain/models"
)

type Screen int

const (
	ScreenSplash Screen = iota
	ScreenFlutterCheck
	ScreenProjectName
	ScreenBundleID
	ScreenArchitecture
	ScreenStateManagement
	ScreenGenerating
	ScreenDone
	ScreenModule
)

type AppModel struct {
	currentScreen   Screen
	projectConfig   models.ProjectConfig
	width           int
	height          int
	splashModel     screens.SplashModel
	flutterCheck    screens.FlutterCheckModel
	projectName     screens.ProjectNameModel
	bundleID        screens.BundleIDModel
	architecture    screens.ArchitectureModel
	stateMgmt       screens.StateMgmtModel
	generating      screens.GeneratingModel
	done            screens.DoneModel
	module          screens.ModuleModel
}

func New() AppModel {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not determine working directory: %v; using '.'\n", err)
		cwd = "."
	}
	return AppModel{
		currentScreen: ScreenSplash,
		projectConfig: models.ProjectConfig{Path: cwd},
		splashModel:   screens.NewSplashModel(),
		flutterCheck:  screens.NewFlutterCheckModel(),
		projectName:   screens.NewProjectNameModel(),
		bundleID:      screens.NewBundleIDModel(),
		architecture:  screens.NewArchitectureModel(),
		stateMgmt:     screens.NewStateMgmtModel(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.splashModel.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.propagateSizeToScreens()
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case screens.NavigateBackMsg:
		return m.navigateBack()
	}

	return m.updateCurrentScreen(msg)
}

func (m *AppModel) propagateSizeToScreens() {
	m.splashModel.SetSize(m.width, m.height)
	m.flutterCheck.SetSize(m.width, m.height)
	m.projectName.SetSize(m.width, m.height)
	m.bundleID.SetSize(m.width, m.height)
	m.architecture.SetSize(m.width, m.height)
	m.stateMgmt.SetSize(m.width, m.height)
	m.generating.SetSize(m.width, m.height)
	m.done.SetSize(m.width, m.height)
	m.module.SetSize(m.width, m.height)
}

// navigateBack moves the app to the logically previous screen.
// Screens that cannot go back (splash, flutter check, generating, done) are no-ops.
func (m AppModel) navigateBack() (tea.Model, tea.Cmd) {
	switch m.currentScreen {
	case ScreenProjectName:
		m.currentScreen = ScreenFlutterCheck
		return m, m.flutterCheck.Init()

	case ScreenBundleID:
		// Reset project-name advance flag so the screen is re-enterable.
		m.projectName = screens.NewProjectNameModel()
		m.projectName.SetSize(m.width, m.height)
		m.currentScreen = ScreenProjectName
		return m, m.projectName.Init()

	case ScreenArchitecture:
		m.bundleID = screens.NewBundleIDModel()
		m.bundleID.SetSize(m.width, m.height)
		m.currentScreen = ScreenBundleID
		return m, m.bundleID.Init()

	case ScreenStateManagement:
		m.architecture = screens.NewArchitectureModel()
		m.architecture.SetSize(m.width, m.height)
		m.currentScreen = ScreenArchitecture
		return m, m.architecture.Init()

	case ScreenModule:
		// Back from module creation returns to the done screen.
		m.currentScreen = ScreenDone
		return m, m.done.Init()
	}

	// All other screens (splash, flutter check, generating, done) ignore back.
	return m, nil
}

func (m AppModel) updateCurrentScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentScreen {
	case ScreenSplash:
		updated, cmd := m.splashModel.Update(msg)
		m.splashModel = updated
		if m.splashModel.ShouldAdvance() {
			m.currentScreen = ScreenFlutterCheck
			return m, m.flutterCheck.Init()
		}
		return m, cmd

	case ScreenFlutterCheck:
		updated, cmd := m.flutterCheck.Update(msg)
		m.flutterCheck = updated
		if m.flutterCheck.ShouldQuit() {
			return m, tea.Quit
		}
		if m.flutterCheck.ShouldAdvance() {
			m.currentScreen = ScreenProjectName
			return m, m.projectName.Init()
		}
		return m, cmd

	case ScreenProjectName:
		updated, cmd := m.projectName.Update(msg)
		m.projectName = updated
		if m.projectName.ShouldAdvance() {
			m.projectConfig.Name = m.projectName.ProjectName()
			m.currentScreen = ScreenBundleID
			return m, m.bundleID.Init()
		}
		return m, cmd

	case ScreenBundleID:
		updated, cmd := m.bundleID.Update(msg)
		m.bundleID = updated
		if m.bundleID.ShouldAdvance() {
			m.projectConfig.BundleID = m.bundleID.BundleID()
			m.currentScreen = ScreenArchitecture
			return m, m.architecture.Init()
		}
		return m, cmd

	case ScreenArchitecture:
		updated, cmd := m.architecture.Update(msg)
		m.architecture = updated
		if m.architecture.ShouldAdvance() {
			m.projectConfig.Architecture = m.architecture.Selected()
			m.currentScreen = ScreenStateManagement
			return m, m.stateMgmt.Init()
		}
		return m, cmd

	case ScreenStateManagement:
		updated, cmd := m.stateMgmt.Update(msg)
		m.stateMgmt = updated
		if m.stateMgmt.ShouldAdvance() {
			m.projectConfig.StateManagement = m.stateMgmt.Selected()
			m.generating = screens.NewGeneratingModel(m.projectConfig)
			m.generating.SetSize(m.width, m.height)
			m.currentScreen = ScreenGenerating
			return m, m.generating.Init()
		}
		return m, cmd

	case ScreenGenerating:
		updated, cmd := m.generating.Update(msg)
		m.generating = updated
		if m.generating.ShouldAdvance() {
			m.done = screens.NewDoneModel(m.projectConfig)
			m.done.SetSize(m.width, m.height)
			m.currentScreen = ScreenDone
			return m, m.done.Init()
		}
		return m, cmd

	case ScreenDone:
		updated, cmd := m.done.Update(msg)
		m.done = updated
		if m.done.ShouldQuit() {
			return m, tea.Quit
		}
		if m.done.ShouldCreateModule() {
			m.module = screens.NewModuleModel(m.projectConfig)
			m.module.SetSize(m.width, m.height)
			m.currentScreen = ScreenModule
			return m, m.module.Init()
		}
		return m, cmd

	case ScreenModule:
		updated, cmd := m.module.Update(msg)
		m.module = updated
		if m.module.ShouldQuit() {
			return m, tea.Quit
		}
		if m.module.ShouldCreateAnother() {
			m.module = screens.NewModuleModel(m.projectConfig)
			m.module.SetSize(m.width, m.height)
			return m, m.module.Init()
		}
		return m, cmd
	}

	return m, nil
}

func (m AppModel) View() string {
	switch m.currentScreen {
	case ScreenSplash:
		return m.splashModel.View()
	case ScreenFlutterCheck:
		return m.flutterCheck.View()
	case ScreenProjectName:
		return m.projectName.View()
	case ScreenBundleID:
		return m.bundleID.View()
	case ScreenArchitecture:
		return m.architecture.View()
	case ScreenStateManagement:
		return m.stateMgmt.View()
	case ScreenGenerating:
		return m.generating.View()
	case ScreenDone:
		return m.done.View()
	case ScreenModule:
		return m.module.View()
	}
	return ""
}
