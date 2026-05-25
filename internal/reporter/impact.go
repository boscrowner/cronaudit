package reporter

import (
	"fmt"
	"sort"
	"strings"
)

// ImpactLevel represents the estimated blast radius of a cron entry failing.
type ImpactLevel string

const (
	ImpactCritical ImpactLevel = "critical"
	ImpactHigh     ImpactLevel = "high"
	ImpactMedium   ImpactLevel = "medium"
	ImpactLow      ImpactLevel = "low"
)

// ImpactEntry holds the assessed impact for a single cron entry.
type ImpactEntry struct {
	Line     int
	Schedule string
	Command  string
	Level    ImpactLevel
	Reasons  []string
}

// ImpactReport is the result of AssessImpact.
type ImpactReport struct {
	Source  string
	Entries []ImpactEntry
}

// AssessImpact evaluates each valid cron entry and assigns an impact level
// based on heuristics such as command keywords, frequency, and path targets.
func AssessImpact(r Report) ImpactReport {
	var entries []ImpactEntry
	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		ie := assessImpact(e)
		entries = append(entries, ie)
	}
	sort.Slice(entries, func(i, j int) bool {
		return impactWeight(entries[i].Level) > impactWeight(entries[j].Level)
	})
	return ImpactReport{Source: r.Source, Entries: entries}
}

func assessImpact(e Entry) ImpactEntry {
	cmd := strings.ToLower(e.Command)
	var reasons []string

	if containsAny(cmd, []string{"rm ", "drop ", "truncate", "delete", "format"}) {
		reasons = append(reasons, "destructive command")
	}
	if containsAny(cmd, []string{"/etc/", "/var/", "/usr/", "/bin/", "/sbin/"}) {
		reasons = append(reasons, "targets system path")
	}
	if containsAny(cmd, []string{"sudo", "su ", "chmod", "chown"}) {
		reasons = append(reasons, "privilege escalation")
	}
	if containsAny(cmd, []string{"curl", "wget", "nc ", "ncat", "ssh"}) {
		reasons = append(reasons, "network activity")
	}

	level := ImpactLow
	switch {
	case len(reasons) >= 3:
		level = ImpactCritical
	case len(reasons) == 2:
		level = ImpactHigh
	case len(reasons) == 1:
		level = ImpactMedium
	}

	return ImpactEntry{
		Line:     e.Line,
		Schedule: e.Schedule,
		Command:  e.Command,
		Level:    level,
		Reasons:  reasons,
	}
}

func impactWeight(l ImpactLevel) int {
	switch l {
	case ImpactCritical:
		return 4
	case ImpactHigh:
		return 3
	case ImpactMedium:
		return 2
	default:
		return 1
	}
}

func containsAny(s string, keywords []string) bool {
	for _, k := range keywords {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}

// FormatImpactReport returns a human-readable summary of the impact report.
func FormatImpactReport(ir ImpactReport) string {
	if len(ir.Entries) == 0 {
		return fmt.Sprintf("impact report for %s: no valid entries\n", ir.Source)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "impact report for %s (%d entries):\n", ir.Source, len(ir.Entries))
	for _, e := range ir.Entries {
		reasons := "none"
		if len(e.Reasons) > 0 {
			reasons = strings.Join(e.Reasons, ", ")
		}
		fmt.Fprintf(&sb, "  line %d [%s] %s — %s (%s)\n",
			e.Line, e.Level, e.Schedule, e.Command, reasons)
	}
	return sb.String()
}
