package readme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xedeveloper/flowerpecker/internal/domain/models"
)

func GenerateReadme(config models.ProjectConfig) error {
	content := buildReadmeContent(config)
	path := filepath.Join(config.Path, "README.md")
	return os.WriteFile(path, []byte(content), 0644)
}

func UpdateReadmeWithModule(projectPath, moduleName string, config models.ProjectConfig) error {
	readmePath := filepath.Join(projectPath, "README.md")
	existing, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	moduleHeading := fmt.Sprintf("### Module: `%s`", moduleName)
	if strings.Contains(string(existing), moduleHeading) {
		// Section already present — skip to avoid duplication.
		return nil
	}

	moduleSection := buildModuleSection(moduleName, config)
	updated := string(existing) + "\n" + moduleSection
	return os.WriteFile(readmePath, []byte(updated), 0644)
}

func buildReadmeContent(config models.ProjectConfig) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", config.Name))
	sb.WriteString(fmt.Sprintf("A Flutter application scaffolded with **FlowerPecker**.\n\n"))
	sb.WriteString("## Project Configuration\n\n")
	sb.WriteString(fmt.Sprintf("| Property | Value |\n|---|---|\n"))
	sb.WriteString(fmt.Sprintf("| Bundle ID | `%s` |\n", config.BundleID))
	sb.WriteString(fmt.Sprintf("| Architecture | %s |\n", config.Architecture.String()))
	sb.WriteString(fmt.Sprintf("| State Management | %s |\n\n", config.StateManagement.String()))

	sb.WriteString("## Project Structure\n\n")
	sb.WriteString("```\n")
	sb.WriteString(buildProjectTree(config))
	sb.WriteString("```\n\n")

	sb.WriteString("## Getting Started\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("# Install dependencies\nflutter pub get\n\n")
	sb.WriteString("# Run the app\nflutter run\n\n")
	sb.WriteString("# Run tests\nflutter test\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Modules\n\n")
	sb.WriteString("_No modules created yet. Use FlowerPecker to scaffold new modules._\n")

	return sb.String()
}

func buildProjectTree(config models.ProjectConfig) string {
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
├── test/
└── pubspec.yaml
`, config.Name)
	case models.ArchMVC:
		return fmt.Sprintf(`%s/
├── lib/
│   ├── core/
│   ├── models/
│   ├── views/
│   ├── controllers/
│   └── main.dart
├── test/
└── pubspec.yaml
`, config.Name)
	default:
		return fmt.Sprintf(`%s/
├── lib/
│   ├── core/
│   ├── models/
│   ├── views/
│   ├── viewmodels/
│   └── main.dart
├── test/
└── pubspec.yaml
`, config.Name)
	}
}

func buildModuleSection(moduleName string, config models.ProjectConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### Module: `%s`\n\n", moduleName))
	sb.WriteString(fmt.Sprintf("Architecture: %s | State: %s\n\n", config.Architecture.String(), config.StateManagement.String()))
	sb.WriteString("```\n")

	switch config.Architecture {
	case models.ArchClean:
		sb.WriteString(fmt.Sprintf(`features/%s/
├── data/
│   ├── datasources/
│   ├── models/
│   └── repositories/
├── domain/
│   ├── entities/
│   ├── repositories/
│   └── usecases/
└── presentation/
    ├── pages/
    └── widgets/
`, moduleName))
	case models.ArchMVC:
		sb.WriteString(fmt.Sprintf("models/%s_model.dart\nviews/%s_view.dart\ncontrollers/%s_controller.dart\n", moduleName, moduleName, moduleName))
	default:
		sb.WriteString(fmt.Sprintf("models/%s_model.dart\nviews/%s_view.dart\nviewmodels/%s_viewmodel.dart\n", moduleName, moduleName, moduleName))
	}

	sb.WriteString("```\n")
	return sb.String()
}
