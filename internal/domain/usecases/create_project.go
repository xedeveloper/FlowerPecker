package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xedeveloper/flowerpecker/internal/domain/models"
	"github.com/xedeveloper/flowerpecker/internal/infrastructure/flutter"
	"github.com/xedeveloper/flowerpecker/internal/infrastructure/generator/components"
	genproject "github.com/xedeveloper/flowerpecker/internal/infrastructure/generator/project"
	"github.com/xedeveloper/flowerpecker/internal/infrastructure/generator/readme"
	"github.com/xedeveloper/flowerpecker/internal/infrastructure/generator/state"
)

func CreateProject(config models.ProjectConfig, onStep func(string)) error {
	onStep("Creating Flutter project...")
	// Run `flutter create` from config.Path so the project lands in the right directory.
	// Using RunCommandInDir avoids path-quoting issues with spaces or brackets.
	_, stderr, err := flutter.RunCommandInDir(config.Path, "create", "--org", orgFromBundleID(config.BundleID), config.Name)
	if err != nil {
		return fmt.Errorf("flutter create failed: %s", stderr)
	}

	projectPath := filepath.Join(config.Path, config.Name)

	onStep("Scaffolding architecture...")
	switch config.Architecture {
	case models.ArchClean:
		if err := genproject.GenerateCleanArchitecture(projectPath, config.Name); err != nil {
			return err
		}
	case models.ArchMVC:
		if err := genproject.GenerateMVCArchitecture(projectPath, config.Name); err != nil {
			return err
		}
	case models.ArchMVVM:
		if err := genproject.GenerateMVVMArchitecture(projectPath, config.Name); err != nil {
			return err
		}
	}

	onStep("Adding state management...")
	if err := generateStateFiles(projectPath, config); err != nil {
		return err
	}

	onStep("Generating UI components...")
	if err := components.GenerateUIComponents(projectPath); err != nil {
		return err
	}

	onStep("Generating local storage service...")
	if err := genproject.GenerateLocalStorage(projectPath); err != nil {
		return err
	}

	onStep("Updating pubspec.yaml...")
	if err := injectPubspecDependencies(projectPath, config.StateManagement); err != nil {
		return err
	}

	onStep("Writing README...")
	updatedConfig := config
	updatedConfig.Path = projectPath
	if err := readme.GenerateReadme(updatedConfig); err != nil {
		return err
	}

	onStep("Running flutter pub get...")
	// Run pub get from inside the project directory — this is more reliable
	// than the global -C flag when the path contains spaces or brackets.
	if _, stderr, err := flutter.RunCommandInDir(projectPath, "pub", "get"); err != nil {
		return fmt.Errorf("flutter pub get failed: %s", strings.TrimSpace(stderr))
	}

	return nil
}

func orgFromBundleID(bundleID string) string {
	parts := strings.Split(bundleID, ".")
	if len(parts) >= 2 {
		return strings.Join(parts[:len(parts)-1], ".")
	}
	return bundleID
}

func generateStateFiles(projectPath string, config models.ProjectConfig) error {
	switch config.StateManagement {
	case models.StateBLoC:
		return state.GenerateBlocFiles(projectPath, "home")
	case models.StateRiverpod:
		return state.GenerateRiverpodFiles(projectPath, "home")
	case models.StateSignals:
		return state.GenerateSignalsFiles(projectPath, "home")
	case models.StateProvider:
		return state.GenerateProviderFiles(projectPath, "home")
	case models.StateGetX:
		return state.GenerateGetXFiles(projectPath, "home")
	}
	return nil
}

func injectPubspecDependencies(projectPath string, stateMgmt models.StateManagement) error {
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return err
	}

	var stateDep string
	switch stateMgmt {
	case models.StateBLoC:
		stateDep = state.PubspecDependency()
	case models.StateRiverpod:
		stateDep = state.RiverpodPubspecDependency()
	case models.StateSignals:
		stateDep = state.SignalsPubspecDependency()
	case models.StateProvider:
		stateDep = state.ProviderPubspecDependency()
	case models.StateGetX:
		stateDep = state.GetXPubspecDependency()
	}

	googleFontsDep := "  google_fonts: ^6.1.0"
	localStorageDeps := genproject.LocalStoragePubspecDeps()
	injection := googleFontsDep + "\n" + localStorageDeps + "\n" + stateDep + "\n"

	updated := strings.Replace(string(content), "dependencies:\n", "dependencies:\n"+injection, 1)
	return os.WriteFile(pubspecPath, []byte(updated), 0644)
}
