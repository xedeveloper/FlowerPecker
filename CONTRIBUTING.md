# Contributing to FlowerPecker

Thank you for your interest in contributing. This document explains how to
set up the project locally, understand the codebase, and submit changes.

---

## Table of Contents

- [Development setup](#development-setup)
- [Project structure at a glance](#project-structure-at-a-glance)
- [Running locally](#running-locally)
- [Code conventions](#code-conventions)
- [How to add a new architecture pattern](#how-to-add-a-new-architecture-pattern)
- [How to add a new state management option](#how-to-add-a-new-state-management-option)
- [How to add a new UI component](#how-to-add-a-new-ui-component)
- [How to add or modify a TUI screen](#how-to-add-or-modify-a-tui-screen)
- [Submitting a pull request](#submitting-a-pull-request)
- [Reporting a bug](#reporting-a-bug)

---

## Development setup

**Requirements**

| Tool | Version |
|---|---|
| Go | ≥ 1.24.3 |
| Flutter SDK | ≥ 3.x (for manual end-to-end testing) |
| macOS or Linux | — |

**Clone and build**

```bash
git clone https://github.com/xedeveloper/flowerpecker.git
cd flowerpecker
go mod download
go build ./...
go vet ./...
```

Both commands must exit with no output (no errors, no warnings) before any
change is considered ready.

---

## Project structure at a glance

```
flowerpecker/
├── main.go
├── assets/ascii/          ← splash screen art (plain text files)
├── pkg/style/theme.go     ← all Lip Gloss design tokens
└── internal/
    ├── app/               ← Bubble Tea screens (presentation layer)
    ├── domain/            ← models + use cases (business rules)
    └── infrastructure/    ← Flutter CLI runner + code generators
```

See `ARCHITECT.md` for the full file tree and layer-by-layer explanation.

---

## Running locally

```bash
# Run directly (recompiles on each invocation)
go run .

# Build a binary and run it
go build -o fp . && ./fp
```

The tool opens in full-screen (alt-screen) mode. Press `Ctrl+C` or `Q` to exit.

---

## Code conventions

### General

- Follow standard Go idioms (`gofmt`, `go vet` must pass).
- Package names are single lowercase words — no underscores, no capitals.
- Exported symbols must have doc comments.
- Error messages are lowercase and do not end with a period.

### Dependency rule (CLEAN architecture)

```
app  →  domain  →  infrastructure
```

- `internal/app` may import `internal/domain` but **not** `internal/infrastructure`.
- `internal/domain` may import `internal/infrastructure` (use cases orchestrate it).
- `internal/infrastructure` must **not** import `internal/app` or `internal/domain`.
- `pkg/style` is imported only by `internal/app`.

Any PR that breaks this rule will be rejected.

### Bubble Tea screens

- Each screen is a self-contained struct with `Init`, `Update`, `View`, `SetSize`, and a
  `ShouldAdvance() bool` method.
- **Never block inside `Update`**. Long-running work (CLI calls, file I/O)
  must be dispatched as a `tea.Cmd`.
- Use `NavigateBack()` from `screens/shared.go` to signal backward navigation.
- Every screen must render a consistent footer via `renderFooterHints()`.

### Dart code generation

- Dart template strings that contain backticks or null-aware operators (`??`,
  `?.`) **must** be built with `strings.Builder` — not raw Go string literals.
- Each generator function must be independently callable with only a
  `projectPath string` argument plus any feature-specific parameters.
- Generators must be idempotent: re-running them must not corrupt an existing
  project (use `os.MkdirAll` with appropriate flags, overwrite files explicitly).

---

## How to add a new architecture pattern

Let's say you want to add **Hexagonal (Ports & Adapters)** architecture.

### 1. Add the constant — `internal/domain/models/project.go`

```go
const (
    ArchClean Architecture = iota
    ArchMVC
    ArchMVVM
    ArchHexagonal          // ← new
)

func (a Architecture) String() string {
    // ...
    case ArchHexagonal:
        return "Hexagonal"
    // ...
}
```

### 2. Write the generator — `internal/infrastructure/generator/project/hexagonal.go`

```go
package project

import "os"
import "path/filepath"

func GenerateHexagonalArchitecture(projectPath, projectName string) error {
    dirs := []string{
        "lib/core/ports",
        "lib/core/adapters",
        "lib/features",
        // ...
    }
    for _, d := range dirs {
        if err := os.MkdirAll(filepath.Join(projectPath, d), 0755); err != nil {
            return err
        }
    }
    // write main.dart, app_theme.dart, etc.
    return nil
}
```

### 3. Wire into the use case — `internal/domain/usecases/create_project.go`

```go
case models.ArchHexagonal:
    if err := genproject.GenerateHexagonalArchitecture(projectPath, config.Name); err != nil {
        return err
    }
```

Do the same in `create_module.go` for module scaffolding.

### 4. Add to the selection list — `internal/app/screens/architecture.go`

```go
list.NewItem("Hexagonal", "Ports & Adapters — inward dependencies only"),
```

### 5. Verify

```bash
go build ./... && go vet ./...
```

---

## How to add a new state management option

Example: adding **Zustand-dart**.

### 1. Add the constant — `internal/domain/models/project.go`

```go
const (
    // existing...
    StateZustand       // ← new
)

func (s StateManagement) String() string {
    // ...
    case StateZustand:
        return "Zustand"
}
```

### 2. Write the generator — `internal/infrastructure/generator/state/zustand.go`

```go
package state

import (
    "os"
    "path/filepath"
)

func ZustandPubspecDependency() string {
    return "  zustand: ^1.0.0"
}

func GenerateZustandFiles(projectPath, feature string) error {
    // write <feature>_store.dart
    dir := filepath.Join(projectPath, "lib/features", feature, "presentation/stores")
    os.MkdirAll(dir, 0755)
    // ...
    return nil
}
```

### 3. Wire into use cases

In `create_project.go` — `generateStateFiles()`:
```go
case models.StateZustand:
    return state.GenerateZustandFiles(projectPath, "home")
```

In `create_project.go` — `injectPubspecDependencies()`:
```go
case models.StateZustand:
    stateDep = state.ZustandPubspecDependency()
```

In `create_module.go` — `generateModuleStateFiles()`:
```go
case models.StateZustand:
    return state.GenerateZustandFiles(config.ProjectPath, config.Name)
```

### 4. Add to the selection list — `internal/app/screens/state_mgmt.go`

```go
list.NewItem("Zustand", "Lightweight reactive store"),
```

### 5. Verify

```bash
go build ./... && go vet ./...
```

---

## How to add a new UI component

All Flutter widgets are generated by
`internal/infrastructure/generator/components/ui_components.go`.

### 1. Add a builder function

```go
func buildFPBadgeDart() string {
    var b strings.Builder
    b.WriteString("import 'package:flutter/material.dart';\n")
    // ... full widget source
    return b.String()
}
```

### 2. Register it in `GenerateUIComponents`

```go
func GenerateUIComponents(projectPath string) error {
    widgetsDir := filepath.Join(projectPath, "lib", "core", "widgets")
    os.MkdirAll(widgetsDir, 0755)

    files := map[string]string{
        // existing entries...
        "fp_badge.dart": buildFPBadgeDart(),  // ← new
    }
    for name, content := range files {
        if err := writeFile(filepath.Join(widgetsDir, name), content); err != nil {
            return err
        }
    }
    return nil
}
```

### Guidelines for widget content

- Use only Flutter's built-in animation primitives (`AnimatedContainer`,
  `AnimatedOpacity`, `TweenAnimationBuilder`, `SlideTransition`).
- No external animation packages.
- Honour the design tokens from `app_theme.dart` (colours, font sizes,
  spacing). Do not hardcode hex values inside widgets.
- Use `strings.Builder` if the Dart source contains backticks or `??`.

---

## How to add or modify a TUI screen

### Adding a new screen

1. Create `internal/app/screens/my_screen.go`:

```go
package screens

import (
    tea "github.com/charmbracelet/bubbletea"
    style "github.com/xedeveloper/flowerpecker/pkg/style"
)

type MyScreenModel struct {
    width, height int
    advance       bool
}

func NewMyScreenModel() MyScreenModel { return MyScreenModel{} }

func (m MyScreenModel) Init() tea.Cmd { return nil }

func (m MyScreenModel) Update(msg tea.Msg) (MyScreenModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            m.advance = true
        case "esc":
            return m, NavigateBack()
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m MyScreenModel) ShouldAdvance() bool { return m.advance }

func (m *MyScreenModel) SetSize(w, h int) { m.width = w; m.height = h }

func (m MyScreenModel) View() string {
    header := renderScreenHeader("My Screen", m.width)
    body   := style.NormalStyle.Render("Hello from MyScreen")
    footer := renderFooterHints(m.width, "[enter] continue  [esc] back  [q] quit")
    return renderFullScreen(header, body, footer, m.width, m.height)
}
```

2. Add the screen constant to `internal/app/app.go`:

```go
const (
    // existing...
    ScreenMyScreen
)
```

3. Add a sub-model field and wire it into `Update`, `navigateBack`, and the
   `windowSizeMsg` handler in `app.go`.

### Modifying an existing screen

- Keep `Init`, `Update`, `View` pure — no global state mutation.
- If you add a new `tea.Cmd` dispatch, make sure the corresponding message
  type is handled in `Update`.
- Run `go vet ./...` after every change.

---

## Submitting a pull request

1. **Fork** the repository and create a branch from `main`:
   ```bash
   git checkout -b feat/my-feature
   ```

2. **Make your changes** following the conventions above.

3. **Verify** the build is clean:
   ```bash
   go build ./... && go vet ./...
   ```

4. **Commit** with a concise, lowercase message:
   ```
   feat: add hexagonal architecture scaffold
   fix: version regex for flutter 3.22 user-branch output
   docs: update README with new state management option
   ```

5. **Open a pull request** against `main`. Include:
   - A clear description of what changed and why.
   - Steps to manually verify the change (which screens to navigate, what to
     look for in the generated Flutter project).
   - Any relevant context (issue number, flutter version tested against).

6. A maintainer will review within a few days. Address any requested changes
   in additional commits on the same branch — do not force-push during review.

---

## Reporting a bug

Open a GitHub issue and include:

- **Go version**: `go version`
- **OS and shell**: e.g. macOS 14, fish 3.7.1
- **Flutter version**: `flutter --version`
- **Flutter install path**: `which flutter`
- **What you expected** vs **what happened**
- **Reproduction steps** (be specific about which screen the error occurs on)

If the tool crashes with a Go panic, include the full stack trace.
