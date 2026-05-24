package flutter

import (
	"regexp"
	"strings"
)

type DoctorCheck struct {
	Name   string
	Status bool
	Detail string
}

type DoctorResult struct {
	FlutterInstalled bool
	FlutterVersion   string
	Checks           []DoctorCheck
	Devices          []string
}

// RunDoctorCheck runs `flutter doctor -v` and returns a structured result.
//
// Flutter is considered installed as long as the binary runs and produces
// output — regardless of exit code (non-zero exit is normal when doctor
// reports warnings such as "unknown channel" or missing toolchains).
func RunDoctorCheck() DoctorResult {
	stdout, stderr, err := RunCommand("doctor", "-v")

	// No output at all means the binary could not be found or executed.
	if stdout == "" && stderr == "" && err != nil {
		return DoctorResult{FlutterInstalled: false}
	}

	// Use whichever stream has content (doctor writes to stdout but some
	// older versions and wrappers mix stderr in).
	raw := stdout
	if raw == "" {
		raw = stderr
	}

	result := ParseDoctorOutput(raw)
	result.Devices = GetDevices()
	return result
}

// ParseDoctorOutput parses the output of `flutter doctor -v`.
//
// It handles the three version-line formats seen across Flutter releases:
//
//	Classic:  "Flutter 3.x.x • channel stable • ..."
//	Verbose:  "Flutter version 3.x.x on channel stable ..."
//	Parenthesised (current):
//	          "[!] Flutter (Channel [user-branch], 3.x.x, on macOS ...)"
//	          "• Flutter version 3.x.x on channel ..."
func ParseDoctorOutput(output string) DoctorResult {
	result := DoctorResult{}
	lines := strings.Split(output, "\n")

	// Version patterns tried in order; first match wins.
	versionPatterns := []*regexp.Regexp{
		// "Flutter version 3.22.0" — verbose / indented line
		regexp.MustCompile(`Flutter version ([\d]+\.[\d]+\.[\d]+)`),
		// "Channel [user-branch], 3.22.0," — parenthesised header
		regexp.MustCompile(`Channel [^,]+,\s*([\d]+\.[\d]+\.[\d]+)`),
		// "Flutter 3.22.0 •" — classic one-line format
		regexp.MustCompile(`Flutter ([\d]+\.[\d]+\.[\d]+)\s+[•·]`),
		// Broad fallback: any semver next to the word Flutter
		regexp.MustCompile(`Flutter[^0-9]+([\d]+\.[\d]+\.[\d]+)`),
	}

	for _, line := range lines {
		for _, re := range versionPatterns {
			if m := re.FindStringSubmatch(line); len(m) > 1 {
				result.FlutterVersion = m[1]
				result.FlutterInstalled = true
				goto versionFound
			}
		}
	}
versionFound:

	// If the binary ran and produced output but the version could not be
	// parsed, we still mark Flutter as installed — a working binary is
	// sufficient evidence.
	if !result.FlutterInstalled && output != "" {
		result.FlutterInstalled = true
	}

	// Parse [✓], [✗], [!], [ ] check lines.
	checkRe := regexp.MustCompile(`^\s*\[(✓|✗|!| )\]\s+(.+)`)
	for _, line := range lines {
		if m := checkRe.FindStringSubmatch(line); len(m) > 2 {
			// Treat [✓] as pass; everything else (!, ✗, space) as warn/fail.
			status := m[1] == "✓"
			name := strings.TrimSpace(m[2])
			// Trim redundant sub-detail like "(Channel stable, 3.22.0, ...)"
			// from the check label so the UI stays readable.
			if idx := strings.Index(name, " (Channel"); idx != -1 {
				name = strings.TrimSpace(name[:idx])
			}
			result.Checks = append(result.Checks, DoctorCheck{
				Name:   name,
				Status: status,
			})
		}
	}

	if len(result.Checks) == 0 && result.FlutterInstalled {
		result.Checks = []DoctorCheck{
			{Name: "Flutter SDK", Status: true, Detail: "v" + result.FlutterVersion},
		}
	}

	return result
}

// GetDevices returns a list of connected device names from `flutter devices`.
func GetDevices() []string {
	stdout, _, err := RunCommand("devices")
	if err != nil {
		return []string{}
	}

	var devices []string
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Device lines contain "•" as a separator: "iPhone 15 • iPhone • arm64 • ..."
		// Skip header / summary lines that start with known non-device prefixes.
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "no devices") ||
			strings.HasPrefix(lower, "found ") ||
			strings.HasPrefix(line, "•") {
			continue
		}
		if strings.Contains(line, "•") {
			name := strings.TrimSpace(strings.SplitN(line, "•", 2)[0])
			if name != "" {
				devices = append(devices, name)
			}
		}
	}
	return devices
}
