# FlowerPecker 🌸

A full-screen TUI tool written in Go that scaffolds Flutter projects with
production-ready boilerplate — architecture layers, state management wiring,
a complete design-system widget library, and a local SQLite storage service
— all generated before you write your first line of feature code.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the tool](#running-the-tool)
- [Screen-by-screen walkthrough](#screen-by-screen-walkthrough)
- [What gets generated](#what-gets-generated)
  - [Architecture scaffolding](#architecture-scaffolding)
  - [State management](#state-management)
  - [UI component library](#ui-component-library)
  - [Local storage service](#local-storage-service)
- [Creating a module after project init](#creating-a-module-after-project-init)
- [Keyboard reference](#keyboard-reference)
- [Supported platforms](#supported-platforms)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

| Requirement | Minimum version |
|---|---|
| Go | 1.24.3 |
| Flutter SDK | 3.x (any channel) |
| macOS / Linux | — |

Flutter must be reachable from your shell. FlowerPecker searches for the
binary in the following order:

1. Your shell's `PATH` (via `exec.LookPath`)
2. `$FLUTTER_ROOT/bin/flutter` and `$FLUTTER_HOME/bin/flutter`
3. Your login shell (`fish -l`, `zsh -l`, `bash -l`) — covers custom install
   paths configured in `~/.config/fish/config.fish`, `~/.zshrc`, etc.
4. A curated list of well-known locations (Homebrew, snap, asdf, fvm, …)

If you use **fish**, make sure your flutter path is exported via
`set -Ux FLUTTER_ROOT /path/to/flutter` or added to `fish_user_paths`.

---

## Installation

```bash
# Clone the repository
git clone https://github.com/xedeveloper/flowerpecker.git
cd flowerpecker

# Build a local binary
go build -o flowerpecker .

# (Optional) Install globally
go install .
```

---

## Running the tool

```bash
# From the project directory
go run .

# Or if installed globally
flowerpecker
```

The tool opens in full-screen (alt-screen) mode and takes over the terminal.

---

## Screen-by-screen walkthrough

```
Splash  ──►  Flutter Check  ──►  Project Name  ──►  Bundle ID
                                                         │
                                                         ▼
Done  ◄──  Generating  ◄──  State Management  ◄──  Architecture
 │
 └──►  Module Creation  (repeatable)
```

### 1. Splash screen
Displays the FlowerPecker ASCII logo and sparrow. Press any key or wait 3 s
to continue.

### 2. Flutter environment check
Runs `flutter doctor -v` and `flutter devices` in the background while
showing a spinner. Results are displayed as a checklist:

```
✓ Flutter SDK detected   v3.22.0
✓ Android toolchain
✗ Xcode  (warning — you can still continue)

Connected Devices:
  • iPhone 15 Pro
```

Press `Enter` to continue if Flutter is found. Press `Q` to quit.

### 3. Project name
Enter a Flutter-valid project name in `snake_case`.

```
Enter your Flutter project name:
> my_awesome_app
```

Validation rejects uppercase letters, hyphens, and spaces inline.

### 4. Bundle identifier
Enter the bundle ID / application ID.

```
Enter bundle identifier:
> com.example.myawesomeapp
```

### 5. Architecture selection
Choose the structural pattern for your Flutter project:

| Option | Layers generated |
|---|---|
| **CLEAN Architecture** | `data/` · `domain/` · `presentation/` per feature |
| **MVC** | `models/` · `views/` · `controllers/` |
| **MVVM** | `models/` · `views/` · `viewmodels/` |

### 6. State management selection
Choose how state is managed inside the generated screens:

| Option | Package added | Boilerplate generated |
|---|---|---|
| **Flutter BLoC** | `flutter_bloc: ^8.1.6` | `_bloc`, `_event`, `_state` files |
| **Riverpod** | `flutter_riverpod: ^2.5.1` | `_provider`, `_notifier` files |
| **Signals** | `signals_flutter: ^5.4.0` | `_signals` file |
| **Provider** | `provider: ^6.1.2` | `_provider` ChangeNotifier file |
| **GetX** | `get: ^4.6.6` | `_controller`, `_binding` files |

### 7. Generation
A progress log shows each step as it completes:

```
✓  Creating Flutter project...
✓  Scaffolding architecture...
✓  Adding state management...
✓  Generating UI components...
✓  Generating local storage service...
✓  Updating pubspec.yaml...
✓  Writing README...
✓  Running flutter pub get...

✓  Project generated successfully!
```

### 8. Done screen
Shows the generated project tree and two options:

- **`C`** — create a new feature module
- **`Q`** — quit

---

## What gets generated

### Architecture scaffolding

#### CLEAN Architecture

```
lib/
├── core/
│   ├── constants/
│   │   └── app_constants.dart
│   ├── errors/
│   │   └── failures.dart
│   ├── theme/
│   │   └── app_theme.dart          ← full design-system tokens
│   ├── utils/
│   │   └── extensions.dart
│   ├── services/
│   │   └── local_storage_service.dart
│   └── widgets/                    ← pre-built UI component library
│       ├── fp_button.dart
│       ├── fp_card.dart
│       ├── fp_divider.dart
│       ├── fp_footer.dart
│       ├── fp_hero.dart
│       ├── fp_input.dart
│       ├── fp_nav.dart
│       ├── fp_toast.dart
│       └── fp_typography.dart
└── features/
    └── <module>/                   ← added when you create a module
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
            ├── widgets/
            └── <state-mgmt-files>
```

#### MVC

```
lib/
├── core/          ← theme, constants, utils, widgets, services
├── models/
├── views/
└── controllers/
```

#### MVVM

```
lib/
├── core/          ← theme, constants, utils, widgets, services
├── models/
├── views/
└── viewmodels/
```

---

### State management

Every architecture option adds the chosen state management package to
`pubspec.yaml` and generates a working `home` feature as a starter example.

---

### UI component library

All components are placed in `lib/core/widgets/` and follow the
**Wired editorial design system** — strict black-and-white palette, square
corners, serif display type, and hairline dividers.

| File | Widgets |
|---|---|
| `fp_button.dart` | `FPButtonPrimary`, `FPButtonOutline`, `FPButtonIconCircular` |
| `fp_card.dart` | `FPStoryCardLarge`, `FPStoryCard`, `FPStoryRow` |
| `fp_input.dart` | `FPTextInput` |
| `fp_nav.dart` | `FPNavBar`, `FPNavLink` |
| `fp_footer.dart` | `FPFooter` |
| `fp_hero.dart` | `FPHeroBand` |
| `fp_divider.dart` | `FPHairlineDivider` |
| `fp_toast.dart` | `FPToast` |
| `fp_typography.dart` | `FPText` with named constructors for every type token |

All components include micro-animations using Flutter's built-in animation
system (`AnimatedContainer`, `AnimatedOpacity`, `TweenAnimationBuilder`).
No external animation packages required.

Design tokens live in `lib/core/theme/app_theme.dart`:

```dart
AppColors.primary       // #000000
AppColors.canvas        // #ffffff
AppColors.hairline      // #e0e0e0
AppColors.link          // #057dbc
AppTextStyles.displayHero  // Playfair Display 64 px
AppTextStyles.bodySm       // Inter 14 px
// … all tokens from the design system
```

---

### Local storage service

`lib/core/services/local_storage_service.dart` is a singleton SQLite-backed
key-value store. Store any JSON-serialisable Dart value with a single call:

```dart
// Save
await LocalStorageService.instance.save('user', {'id': 1, 'name': 'Ada'});

// Read
final user = await LocalStorageService.instance.get<Map<String, dynamic>>('user');

// Check existence
final exists = await LocalStorageService.instance.containsKey('user');

// Delete
await LocalStorageService.instance.delete('user');

// Wipe all
await LocalStorageService.instance.clear();
```

Dependencies added automatically: `sqflite: ^2.3.3` · `path: ^1.9.0`

---

## Creating a module after project init

On the Done screen press **`C`**, enter a module name (e.g. `auth`), and
FlowerPecker scaffolds all required files for the chosen architecture:

**CLEAN example — `auth` module:**

```
lib/features/auth/
├── data/
│   ├── datasources/
│   │   ├── auth_remote_datasource.dart
│   │   └── auth_local_datasource.dart
│   ├── models/
│   │   └── auth_model.dart
│   └── repositories/
│       └── auth_repository_impl.dart
├── domain/
│   ├── entities/
│   │   └── auth_entity.dart
│   ├── repositories/
│   │   └── auth_repository.dart
│   └── usecases/
│       └── auth_usecase.dart
└── presentation/
    ├── pages/
    │   └── auth_page.dart
    ├── widgets/
    │   └── auth_form_widget.dart
    └── <state-mgmt-files>
```

The project's `README.md` is updated automatically with a new module section
after each creation.

---

## Keyboard reference

| Key | Action |
|---|---|
| `↑` / `↓` | Navigate list items |
| `Enter` | Confirm / continue |
| `Esc` | Go back to previous screen |
| `C` | Create a module (Done screen) |
| `Q` / `Ctrl+C` | Quit |

---

## Supported platforms

| Platform | Status |
|---|---|
| macOS (Apple Silicon + Intel) | ✓ |
| Linux | ✓ |
| Windows | Not supported |

---

## Troubleshooting

### Flutter SDK not detected

FlowerPecker runs a 4-tier search for the flutter binary. If it still fails:

1. Verify flutter works in your shell: `which flutter`
2. Set the explicit env var: `set -Ux FLUTTER_ROOT /path/to/flutter` (fish)
   or `export FLUTTER_ROOT=/path/to/flutter` (bash/zsh)
3. Restart your terminal so the variable is exported to new processes.

### `flutter pub get` fails

The tool runs `pub get` from inside the project directory. If it fails,
check that your Flutter installation is healthy: `flutter doctor`.

### Text is invisible / hard to read

All UI text is rendered in white (`#ffffff`). If text is still not visible,
ensure your terminal emulator is not overriding foreground colours globally.
