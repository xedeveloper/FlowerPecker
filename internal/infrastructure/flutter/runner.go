package flutter

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// ── Binary resolution ─────────────────────────────────────────────────────────

var (
	flutterOnce   sync.Once
	flutterBinary string

	shellPathOnce  sync.Once
	shellPathValue string // PATH as reported by the user's login shell
)

// FindFlutter resolves the absolute path to the flutter binary.
//
// Resolution order — each tier is tried only when the previous one fails:
//
//  1. Process PATH via exec.LookPath  — fast path, works for most users
//  2. Explicit env vars: $FLUTTER_ROOT, $FLUTTER_HOME
//  3. User's login shell: `$SHELL -l -c "which flutter"` (fish / zsh / bash)
//     This is the critical tier for fish users and custom install paths,
//     because the shell sources its own config and knows the real PATH.
//  4. Curated static candidate list — last-resort filesystem scan covering
//     Homebrew, snap, asdf, fvm, home-dir installs, etc.
//  5. Bare "flutter" — lets the OS produce a clean error if all else fails.
//
// The resolved path is cached; the scan happens only once per run.
func FindFlutter() string {
	flutterOnce.Do(func() {
		// Tier 1 — standard PATH lookup.
		if p, err := exec.LookPath("flutter"); err == nil {
			flutterBinary = p
			return
		}

		home, _ := os.UserHomeDir()

		// Tier 2 — explicit env-var overrides.
		for _, key := range []string{"FLUTTER_ROOT", "FLUTTER_HOME"} {
			if root := os.Getenv(key); root != "" {
				p := filepath.Join(root, "bin", "flutter")
				if fileExists(p) {
					flutterBinary = p
					return
				}
			}
		}

		// Tier 3 — ask the user's login shell.
		// This handles fish's fish_user_paths, zsh's .zshrc, custom install
		// directories (e.g. ~/Development/[05] Framework/flutter/bin) and
		// any other non-standard location the user already has working.
		if p := findFlutterViaShell(); p != "" {
			flutterBinary = p
			return
		}

		// Tier 4 — static filesystem scan for well-known locations.
		candidates := []string{
			// Homebrew (Apple Silicon + Intel)
			"/opt/homebrew/bin/flutter",
			"/usr/local/bin/flutter",
			"/usr/bin/flutter",

			// Home-directory installs (official Flutter docs defaults)
			filepath.Join(home, "flutter", "bin", "flutter"),
			filepath.Join(home, "development", "flutter", "bin", "flutter"),
			filepath.Join(home, "Development", "flutter", "bin", "flutter"),
			filepath.Join(home, "Developer", "flutter", "bin", "flutter"),
			filepath.Join(home, ".flutter", "bin", "flutter"),

			// snap (Ubuntu / Debian)
			filepath.Join(home, "snap", "flutter", "current", "bin", "flutter"),
			"/snap/flutter/current/bin/flutter",

			// asdf version manager
			filepath.Join(home, ".asdf", "shims", "flutter"),

			// fvm (Flutter Version Manager)
			filepath.Join(home, "fvm", "default", "bin", "flutter"),
			filepath.Join(home, ".fvm", "default", "bin", "flutter"),

			// pub-cache global activation
			filepath.Join(home, ".pub-cache", "bin", "flutter"),

			// XDG / opt
			filepath.Join(home, ".local", "share", "flutter", "bin", "flutter"),
			"/opt/flutter/bin/flutter",
		}

		for _, p := range candidates {
			if fileExists(p) {
				flutterBinary = p
				return
			}
		}

		// Tier 5 — bare fallback; lets the OS report a useful error.
		flutterBinary = "flutter"
	})
	return flutterBinary
}

// findFlutterViaShell invokes the user's login shell and asks it to locate
// the flutter binary. This is reliable for any shell that correctly configures
// PATH in its rc/config files — including fish with fish_user_paths.
func findFlutterViaShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	// Build the command variants to try, in preference order.
	// We prefer the login shell (-l) because that sources the full user profile;
	// interactive mode (-i) is a fallback for shells where -l does not load PATH.
	type attempt struct{ args []string }
	var attempts []attempt

	switch filepath.Base(shell) {
	case "fish":
		// Fish: -l loads ~/.config/fish/config.fish and merges fish_user_paths.
		// `command -v` is more portable inside fish than `which`.
		attempts = []attempt{
			{[]string{shell, "-l", "-c", "command -v flutter"}},
			{[]string{shell, "-c", "command -v flutter"}},
		}
	case "zsh":
		attempts = []attempt{
			{[]string{shell, "-l", "-c", "command -v flutter"}},
			{[]string{shell, "-i", "-c", "command -v flutter 2>/dev/null"}},
		}
	default: // bash, sh, dash, etc.
		attempts = []attempt{
			{[]string{shell, "-l", "-c", "command -v flutter"}},
			{[]string{shell, "-c", "command -v flutter"}},
		}
	}

	for _, a := range attempts {
		out, err := exec.Command(a.args[0], a.args[1:]...).Output()
		if err != nil {
			continue
		}
		// Some shells print one path per line (aliases, etc.); use the first.
		raw := strings.TrimSpace(string(out))
		if raw == "" {
			continue
		}
		line := raw
		if idx := strings.IndexByte(raw, '\n'); idx != -1 {
			line = strings.TrimSpace(raw[:idx])
		}
		// Reject shell error messages that land on stdout.
		if strings.Contains(line, "not found") || strings.Contains(line, "no flutter") {
			continue
		}
		if fileExists(line) {
			return line
		}
	}
	return ""
}

// ── PATH enrichment ───────────────────────────────────────────────────────────

// loginShellPath returns the PATH string as reported by the user's login shell.
// This is fetched once and cached; the result is used to enrich the environment
// of every flutter child process so that dart and other SDK tools are reachable.
func loginShellPath() string {
	shellPathOnce.Do(func() {
		shell := os.Getenv("SHELL")
		if shell == "" {
			return
		}

		var cmd *exec.Cmd
		switch filepath.Base(shell) {
		case "fish":
			// `string join` converts fish's list-typed PATH to a colon-separated string.
			cmd = exec.Command(shell, "-l", "-c", "string join ':' $PATH")
		default:
			cmd = exec.Command(shell, "-l", "-c", "printf '%s' \"$PATH\"")
		}

		out, err := cmd.Output()
		if err != nil {
			return
		}
		shellPathValue = strings.TrimSpace(string(out))
	})
	return shellPathValue
}

// augmentedEnv returns an environment slice that merges:
//   - The current process environment (os.Environ)
//   - The PATH from the user's login shell  (covers fish_user_paths, etc.)
//   - The flutter binary's own bin/ directory (ensures dart is on PATH too)
//
// The most specific directories are prepended so they win over system defaults.
func augmentedEnv() []string {
	env := os.Environ()

	// Collect extra dirs to prepend to PATH.
	var extra []string

	// 1. Flutter's own bin/ dir (ensures dart, pub, etc. are reachable).
	if bin := FindFlutter(); bin != "flutter" {
		extra = append(extra, filepath.Dir(bin))
	}

	// 2. Dirs from the login shell PATH that are not already present.
	if shellPath := loginShellPath(); shellPath != "" {
		// Find current PATH value in the env slice.
		currentPath := os.Getenv("PATH")
		for _, dir := range filepath.SplitList(shellPath) {
			if dir == "" {
				continue
			}
			if !strings.Contains(currentPath, dir) {
				extra = append(extra, dir)
			}
		}
	}

	if len(extra) == 0 {
		return env
	}

	// Prepend the extra dirs to the PATH entry in the env slice.
	prefix := strings.Join(extra, string(os.PathListSeparator)) + string(os.PathListSeparator)
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			env[i] = "PATH=" + prefix + e[len("PATH="):]
			return env
		}
	}
	// No PATH entry found — add one.
	return append(env, "PATH="+prefix)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// fileExists reports whether path is an existing regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// ── Public command runners ────────────────────────────────────────────────────

// RunCommand runs a flutter command and returns stdout, stderr, and any error.
func RunCommand(args ...string) (string, string, error) {
	cmd := exec.Command(FindFlutter(), args...)
	cmd.Env = augmentedEnv()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RunCommandInDir runs a flutter command with its working directory set to dir.
// This is the reliable alternative to flutter's global -C flag, which can
// fail when the project path contains spaces or special characters (e.g. [05]).
func RunCommandInDir(dir string, args ...string) (string, string, error) {
	cmd := exec.Command(FindFlutter(), args...)
	cmd.Dir = dir
	cmd.Env = augmentedEnv()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RunCommandWithProgress runs a flutter command under the supplied context,
// streaming each output line to onLine. The context can be used to cancel or
// time-out long-running commands from the caller.
func RunCommandWithProgress(ctx context.Context, args []string, onLine func(string)) error {
	cmd := exec.CommandContext(ctx, FindFlutter(), args...)
	cmd.Env = augmentedEnv()
	pr, pw, err := os.Pipe()
	if err != nil {
		return err
	}
	cmd.Stdout = pw
	cmd.Stderr = pw

	if err := cmd.Start(); err != nil {
		pw.Close()
		pr.Close()
		return err
	}

	pw.Close()

	scanner := bufio.NewScanner(pr)
	for scanner.Scan() {
		onLine(scanner.Text())
	}
	pr.Close()

	return cmd.Wait()
}
