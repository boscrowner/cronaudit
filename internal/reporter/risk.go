package reporter

import (
	"fmt"
	"strings"
)

// RiskLevel represents the severity of a cron entry's risk assessment.
type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// RiskEntry holds the risk assessment for a single cron entry.
type RiskEntry struct {
	Line     int
	Schedule string
	Command  string
	Level    RiskLevel
	Reasons  []string
}

// RiskReport holds the full risk assessment for a report.
type RiskReport struct {
	Source  string
	Entries []RiskEntry
}

// AssessRisk evaluates each valid cron entry for risk indicators and returns a RiskReport.
func AssessRisk(r Report) RiskReport {
	var entries []RiskEntry
	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		level, reasons := assessEntry(e)
		entries = append(entries, RiskEntry{
			Line:     e.Line,
			Schedule: e.Schedule,
			Command:  e.Command,
			Level:    level,
			Reasons:  reasons,
		})
	}
	return RiskReport{Source: r.Source, Entries: entries}
}

func assessEntry(e Entry) (RiskLevel, []string) {
	var reasons []string

	if containsRiskyCommand(e.Command) {
		reasons = append(reasons, "command contains potentially destructive operation")
	}
	if runsAsRoot(e.Command) {
		reasons = append(reasons, "command runs with elevated privileges")
	}
	if isHighFrequency(e.Schedule) {
		reasons = append(reasons, "schedule runs more frequently than every 5 minutes")
	}
	if writesToSensitivePath(e.Command) {
		reasons = append(reasons, "command writes to a sensitive system path")
	}

	switch {
	case len(reasons) >= 3:
		return RiskHigh, reasons
	case len(reasons) >= 1:
		return RiskMedium, reasons
	default:
		return RiskLow, reasons
	}
}

func containsRiskyCommand(cmd string) bool {
	risky := []string{"rm ", "rm -", "mkfs", "dd ", "truncate", "shred", "> /dev/"}
	lower := strings.ToLower(cmd)
	for _, r := range risky {
		if strings.Contains(lower, r) {
			return true
		}
	}
	return false
}

func runsAsRoot(cmd string) bool {
	lower := strings.ToLower(cmd)
	return strings.Contains(lower, "sudo ") || strings.HasPrefix(strings.TrimSpace(lower), "su ")
}

func writesToSensitivePath(cmd string) bool {
	paths := []string{"/etc/", "/boot/", "/sys/", "/proc/"}
	for _, p := range paths {
		if strings.Contains(cmd, p) {
			return true
		}
	}
	return false
}

func isHighFrequency(schedule string) bool {
	parts := strings.Fields(schedule)
	if len(parts) < 2 {
		return false
	}
	minute := parts[0]
	return strings.HasPrefix(minute, "*/") && minute != "*/5" && minute != "*/10" && minute != "*/15" && minute != "*/30"
}

// FormatRiskReport returns a human-readable summary of the risk report.
func FormatRiskReport(rr RiskReport) string {
	if len(rr.Entries) == 0 {
		return fmt.Sprintf("risk report for %s: no valid entries assessed\n", rr.Source)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("risk report for %s:\n", rr.Source))
	for _, e := range rr.Entries {
		sb.WriteString(fmt.Sprintf("  [%s] line %d: %s %s\n", e.Level, e.Line, e.Schedule, e.Command))
		for _, reason := range e.Reasons {
			sb.WriteString(fmt.Sprintf("    - %s\n", reason))
		}
	}
	return sb.String()
}
