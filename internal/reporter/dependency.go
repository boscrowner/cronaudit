package reporter

import (
	"fmt"
	"sort"
	"strings"
)

// DependencyEntry represents a detected scheduling dependency between two cron entries.
type DependencyEntry struct {
	A       Entry
	B       Entry
	Reason  string
}

// DependencyReport holds the result of a dependency analysis between cron entries.
type DependencyReport struct {
	Source       string
	Dependencies []DependencyEntry
}

// DetectDependencies analyzes a Report and identifies pairs of valid entries
// that share the same command base or overlap on schedule, suggesting a
// potential ordering dependency.
func DetectDependencies(r Report) DependencyReport {
	valid := FilterValid(r)
	entries := valid.Entries

	var deps []DependencyEntry

	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			a := entries[i]
			b := entries[j]

			if reason, ok := detectDep(a, b); ok {
				deps = append(deps, DependencyEntry{A: a, B: b, Reason: reason})
			}
		}
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].A.Line < deps[j].A.Line
	})

	return DependencyReport{
		Source:       r.Source,
		Dependencies: deps,
	}
}

// FormatDependencyReport returns a human-readable summary of detected dependencies.
func FormatDependencyReport(dr DependencyReport) string {
	if len(dr.Dependencies) == 0 {
		return fmt.Sprintf("source: %s\nNo scheduling dependencies detected.\n", dr.Source)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "source: %s\n", dr.Source)
	fmt.Fprintf(&sb, "%d potential dependency pair(s) detected:\n", len(dr.Dependencies))
	for _, d := range dr.Dependencies {
		fmt.Fprintf(&sb, "  [line %d] %q  <->  [line %d] %q\n    reason: %s\n",
			d.A.Line, d.A.Raw, d.B.Line, d.B.Raw, d.Reason)
	}
	return sb.String()
}

func detectDep(a, b Entry) (string, bool) {
	cmdA := baseCommand(a.Command)
	cmdB := baseCommand(b.Command)

	if cmdA != "" && cmdA == cmdB {
		return "same base command: " + cmdA, true
	}

	if a.Schedule == b.Schedule {
		return "identical schedule: " + a.Schedule, true
	}

	return "", false
}

func baseCommand(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ""
	}
	base := parts[0]
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	return base
}
