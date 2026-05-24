# FlowerPecker — Architecture

This document describes the internal architecture of the FlowerPecker tool
itself — its layer responsibilities, complete file tree, data flow, and key
design decisions.

---

## Table of Contents

- [Architectural philosophy](#architectural-philosophy)
- [Complete file tree](#complete-file-tree)
- [Layer responsibilities](#layer-responsibilities)
  - [Presentation — `internal/app`](#presentation--internalapp)
  - [Domain — `internal/domain`](#domain--internaldomain)
  - [Infrastructure — `internal/infrastructure`](#infrastructure--internalinfrastructure)
  - [Shared style — `pkg/style`](#shared-style--pkgstyle)
  - [Assets — `assets/`](#assets--assets)
- [Screen flow & state transitions](#screen-flow--state-transitions)
- [Flutter binary resolution](#flutter-binary-resolution)
- [Code generation pipeline](#code-generation-pipeline)
- [Dependency graph](#dependency-graph)
- [Key design decisions](#key-design-decisions)

---

## Architectural philosophy

FlowerPecker is itself built with **CLEAN Architecture** — the same pattern
it can scaffold for Flutter projects. The codebase is divided into three
concentric layers:

```
┌─────────────────────────────────────────┐
│  Presentation  (internal/app)           │  Bubble Tea TUI screens
├─────────────────────────────────────────┤
│  Domain  (internal/domain)              │  Business rules, models, use-cases
├─────────────────────────────────────────┤
│  Infrastructure  (internal/infrastructure) │  Flutter CLI, file generators
└─────────────────────────────────────────┘
```

Dependency rule: every arrow points inward. Infrastructure knows nothing
about screens; domain knows nothing about the CLI or file system specifics.

---

## Complete file tree

```
flowerpecker/
│
├── main.go                              Entry point — creates and runs the
│                                        Bubble Tea program with AltScreen.
│
├── go.mod / go.sum                      Module: github.com/xedeveloper/flowerpecker
│
├── assets/
│   └── ascii/
│       ├── sparrow.txt                  ASCII sparrow bird rendered on the
│       │                                splash screen.
│       └── logo.txt                     ASCII "FLOWERPECKER" logotype.
│
├── pkg/
│   └── style/
│       └── theme.go                     Lip Gloss design tokens — colours,
│                                        typography sizes, reusable styles.
│                                        Single source of truth for all UI
│                                        styling across every screen.
│
└── internal/
    │
    ├── app/
    │   ├── app.go                       Root Bubble Tea model.  Owns the
    │   │                                active screen enum, the shared
    │   │                                ProjectConfig, window dimensions,
    │   │                                and all screen sub-models.
    │   │                                Handles screen transitions and back
    │   │                                navigation (NavigateBackMsg).
    │   │
    │   └── screens/
    │       ├── shared.go                Cross-screen helpers:
    │       │                            • NavigateBackMsg / NavigateBack()
    │       │                            • newStyledDelegate() for lists
    │       │                            • renderPanelCentered()
    │       │                            Also hosts renderScreenHeader(),
    │       │                            renderFooterHints(), renderFullScreen(),
    │       │                            and countLines() — used by every screen.
    │       │
    │       ├── splash.go                Splash screen.  Loads ASCII assets,
    │       │                            renders white-on-black layout,
    │       │                            auto-advances after 3 s via tea.Tick.
    │       │
    │       ├── flutter_check.go         Environment check screen.  Dispatches
    │       │                            usecases.CheckFlutter() as a tea.Cmd,
    │       │                            shows spinner while running, renders
    │       │                            checklist of doctor results and devices.
    │       │
    │       ├── project_name.go          snake_case text-input screen (step 1/2).
    │       │                            Validates input inline; advances on Enter.
    │       │
    │       ├── bundle_id.go             Bundle-identifier input (step 2/2).
    │       │                            Validates presence of at least one dot.
    │       │
    │       ├── architecture.go          Three-item Bubble Tea list:
    │       │                            CLEAN · MVC · MVVM.
    │       │
    │       ├── state_mgmt.go            Five-item Bubble Tea list:
    │       │                            BLoC · Riverpod · Signals · Provider · GetX.
    │       │
    │       ├── generating.go            Progress screen.  Starts generation as
    │       │                            a tea.Cmd (never blocks the event loop).
    │       │                            Shows spinner + completed-step log.
    │       │                            Auto-advances to Done on success.
    │       │
    │       ├── done.go                  Success screen.  Shows project tree,
    │       │                            offers [C] create module / [Q] quit.
    │       │
    │       └── module.go                Module creation screen.  Text input →
    │                                    async generation → success, with option
    │                                    to create another.
    │
    ├── domain/
    │   ├── models/
    │   │   ├── project.go               Architecture enum (CLEAN/MVC/MVVM),
    │   │   │                            StateManagement enum (BLoC/…/GetX),
    │   │   │                            ProjectConfig struct.
    │   │   │
    │   │   └── module.go                ModuleConfig struct (name, project path,
    │   │                                arch, state management).
    │   │
    │   └── usecases/
    │       ├── flutter_check.go         CheckFlutter() → FlutterCheckResult.
    │       │                            Thin facade over flutter.RunDoctorCheck().
    │       │
    │       ├── create_project.go        Full 8-step project creation pipeline:
    │       │                            1. flutter create (RunCommandInDir)
    │       │                            2. Scaffold architecture
    │       │                            3. Generate state management files
    │       │                            4. Generate UI components
    │       │                            5. Generate local storage service
    │       │                            6. Inject pubspec.yaml dependencies
    │       │                            7. Write README.md
    │       │                            8. flutter pub get (RunCommandInDir)
    │       │
    │       └── create_module.go         Module scaffolding (CLEAN/MVC/MVVM),
    │                                    state file generation, README update.
    │                                    Contains all Dart template functions.
    │
    └── infrastructure/
        │
        ├── flutter/
        │   ├── runner.go                Flutter binary resolution (4-tier search),
        │   │                            login-shell PATH enrichment, augmentedEnv().
        │   │                            Public API:
        │   │                            • RunCommand(args…)
        │   │                            • RunCommandInDir(dir, args…)
        │   │                            • RunCommandWithProgress(ctx, args, onLine)
        │   │
        │   └── doctor.go                RunDoctorCheck(), ParseDoctorOutput(),
        │                                GetDevices().
        │                                Multi-pattern version extraction handles
        │                                all known flutter doctor output formats.
        │
        └── generator/
            ├── project/
            │   ├── shared.go            writeFile() helper + Dart templates
            │   │                        shared by all three arch generators:
            │   │                        app_theme.dart, failures.dart,
            │   │                        extensions.dart, app_constants.dart.
            │   │
            │   ├── clean.go             Generates CLEAN arch skeleton:
            │   │                        core/ + features/ + main.dart.
            │   │
            │   ├── mvc.go               Generates MVC skeleton:
            │   │                        models/ + views/ + controllers/ + main.dart.
            │   │
            │   ├── mvvm.go              Generates MVVM skeleton:
            │   │                        models/ + views/ + viewmodels/ + main.dart.
            │   │
            │   └── local_storage.go     Generates lib/core/services/
            │                            local_storage_service.dart.
            │                            LocalStoragePubspecDeps() returns the
            │                            sqflite + path dependency lines.
            │
            ├── state/
            │   ├── bloc.go              BLoC boilerplate + PubspecDependency().
            │   ├── riverpod.go          Riverpod boilerplate + RiverpodPubspecDependency().
            │   ├── signals.go           Signals boilerplate + SignalsPubspecDependency().
            │   ├── provider.go          Provider boilerplate + ProviderPubspecDependency().
            │   └── getx.go              GetX boilerplate + GetXPubspecDependency().
            │
            ├── components/
            │   └── ui_components.go     Writes all 9 Flutter widget files into
            │                            lib/core/widgets/. Each widget uses
            │                            Flutter's built-in animation system.
            │
            └── readme/
                └── generator.go         GenerateReadme() — initial project README.
                                         UpdateReadmeWithModule() — appends module
                                         section; deduplicates on repeated calls.
```

---

## Layer responsibilities

### Presentation — `internal/app`

- **Only layer** that imports Bubble Tea and Lip Gloss.
- Each screen is an independent model implementing `Init() tea.Cmd`,
  `Update(tea.Msg) (Model, tea.Cmd)`, and `View() string`.
- Screens communicate upward via typed messages (`NavigateBackMsg`,
  `generationDoneMsg`, etc.) — never by calling each other directly.
- `app.go` is the single orchestrator: it holds the screen enum, routes
  messages to the active sub-model, and owns the `ProjectConfig` that
  accumulates user inputs across screens.
- All async operations (flutter doctor, project generation) are dispatched as
  `tea.Cmd` so the event loop is never blocked.

### Domain — `internal/domain`

- **No imports from `internal/app` or `internal/infrastructure`** (dependency
  rule).
- `models/` defines plain Go structs and enums — no framework coupling.
- `usecases/` orchestrates the infrastructure layer; each use case is a
  single exported function with an `onStep func(string)` progress callback.

### Infrastructure — `internal/infrastructure`

- **Flutter runner** (`flutter/`): resolves the binary path once (cached with
  `sync.Once`), enriches the child-process environment with the login shell's
  PATH so `dart` and other SDK tools are always reachable, and exposes three
  command runners.
- **Generators** (`generator/`): pure file-writing code. Each generator
  receives a `projectPath string` and writes Dart source files. No network,
  no state, no side effects beyond disk writes.

### Shared style — `pkg/style`

Design tokens from `DESIGN.md` translated to Lip Gloss:

| Token group | Styles |
|---|---|
| Text | `TitleStyle` `SubtitleStyle` `NormalStyle` `EyebrowStyle` |
| Interactive | `ButtonPrimaryStyle` `ButtonOutlineStyle` `InputStyle` |
| Feedback | `SuccessStyle` `ErrorStyle` `SpinnerStyle` |
| Layout | `FooterStyle` `PanelStyle` `BorderStyle` `HairlineStyle` |
| Selection | `SelectedStyle` |

All text foregrounds are white (`#ffffff`) for visibility on dark terminals.
Borders use `ColorHairline` (`#e0e0e0`).

### Assets — `assets/`

Plain text files read at runtime by `screens/splash.go` via `os.ReadFile`.
No embed; the binary must be run from the project root or the assets directory
must be accessible from the working directory.

---

## Screen flow & state transitions

```
AppModel.currentScreen
         │
         ▼
  ScreenSplash ──(ShouldAdvance)──► ScreenFlutterCheck
                                            │
                                    (ShouldAdvance)
                                            │
                                            ▼
                                    ScreenProjectName ──(Esc)──► back
                                            │
                                    (ShouldAdvance)
                                            │
                                            ▼
                                    ScreenBundleID ──(Esc)──► back
                                            │
                                    (ShouldAdvance)
                                            │
                                            ▼
                                    ScreenArchitecture ──(Esc)──► back
                                            │
                                    (ShouldAdvance)
                                            │
                                            ▼
                                    ScreenStateManagement ──(Esc)──► back
                                            │
                                    (ShouldAdvance)
                                            │
                                            ▼
                                    ScreenGenerating
                                            │
                                    (generationDoneMsg)
                                            │
                                            ▼
                                    ScreenDone ──(C)──► ScreenModule
                                            │                  │
                                           (Q)          (repeatable)
                                            │
                                          Quit
```

Back navigation is handled via `NavigateBackMsg` — any screen can emit it;
`app.go`'s `navigateBack()` maps each screen to its predecessor.

---

## Flutter binary resolution

`internal/infrastructure/flutter/runner.go` — `FindFlutter()` — runs once
per process (guarded by `sync.Once`) and searches in four tiers:

```
Tier 1  exec.LookPath("flutter")
        Fast path — works when flutter is in the process PATH.

Tier 2  $FLUTTER_ROOT/bin/flutter
        $FLUTTER_HOME/bin/flutter
        Explicit environment variable overrides.

Tier 3  $SHELL -l -c "command -v flutter"
        Invokes the user's login shell (fish / zsh / bash).
        Reliable for fish users who configure PATH via fish_user_paths
        or config.fish, and for any custom install path.

Tier 4  Static filesystem scan
        /opt/homebrew/bin, ~/flutter/bin, ~/development/flutter/bin,
        snap, asdf, fvm, pub-cache, XDG dirs, /opt/flutter/bin.

Tier 5  Bare "flutter" — lets the OS report a clear error.
```

`augmentedEnv()` merges the login-shell PATH into every child process so
`dart` and other SDK tools that live beside the `flutter` binary are
reachable even if not in the original process PATH.

---

## Code generation pipeline

`usecases.CreateProject` drives the pipeline; each step calls one or more
infrastructure functions:

```
1. RunCommandInDir(config.Path, "create", "--org", …, name)
         └── flutter create — project skeleton

2. generator/project/{clean,mvc,mvvm}.go
         └── architecture folder structure + main.dart + core/ templates

3. generator/state/{bloc,riverpod,signals,provider,getx}.go
         └── state management boilerplate for "home" feature

4. generator/components/ui_components.go
         └── 9 widget files in lib/core/widgets/

5. generator/project/local_storage.go
         └── lib/core/services/local_storage_service.dart

6. injectPubspecDependencies()
         └── reads pubspec.yaml, inserts google_fonts + sqflite +
             path + state-mgmt dep under `dependencies:`

7. generator/readme/generator.go
         └── README.md with project tree

8. RunCommandInDir(projectPath, "pub", "get")
         └── resolves all pubspec dependencies
```

Module creation follows the same pattern but scoped to a single feature
directory and without steps 4–8.

---

## Dependency graph

```
main.go
  └── internal/app
        ├── internal/domain/models          (read-only)
        └── internal/domain/usecases
              ├── internal/domain/models
              └── internal/infrastructure
                    ├── flutter/runner.go
                    ├── flutter/doctor.go
                    └── generator/*
```

The `pkg/style` package is imported only by `internal/app` (presentation).
Infrastructure packages never import `pkg/style` or `internal/app`.

---

## Key design decisions

| Decision | Rationale |
|---|---|
| Bubble Tea with AltScreen | Full-screen takeover; consistent experience regardless of terminal scroll buffer. |
| `sync.Once` for binary resolution | Flutter path search involves shell subprocess calls — doing it once and caching prevents repeated latency on every command. |
| `RunCommandInDir` over `-C` flag | Flutter's `-C` flag fails on paths with spaces or special characters. Setting `cmd.Dir` at the OS level is unconditional. |
| Multi-pattern version regex | Flutter doctor output format differs across channels (stable vs user-branch vs debug build). Regex cascade handles all known variants and degrades gracefully. |
| `onStep` progress callback | Decouples long-running use cases from the TUI; screens receive step strings via `tea.Cmd` without directly depending on infrastructure. |
| `strings.Builder` for Dart templates | Dart source contains backticks and null-aware `??` operators that would prematurely close Go raw-string literals. `strings.Builder` avoids all delimiter escaping issues. |
| README deduplication | `UpdateReadmeWithModule` checks for an existing section heading before appending, making module creation idempotent on repeated runs. |
