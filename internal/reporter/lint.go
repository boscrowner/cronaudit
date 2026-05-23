package reporter

import "fmt"

// LintSeverity represents the severity level of a lint warning.
type LintSeverity string

const (
	SeverityWarn  LintSeverity = "warn"
	SeverityError LintSeverity = "error"
	SeverityInfo  LintSeverity = "info"
)

// LintWarning describes a potential issue found in a cron entry.
type LintWarning struct {
	Line     int          `json:"line"`
	Schedule string       `json:"schedule"`
	Message  string       `json:"message"`
	Severity LintSeverity `json:"severity"`
}

// LintResult holds all warnings produced by LintReport.
type LintResult struct {
	Warnings []LintWarning `json:"warnings"`
	Total    int           `json:"total"`
}

// LintReport inspects a Report for suspicious or risky cron patterns
// and returns a LintResult with categorised warnings.
func LintReport(r Report) LintResult {
	var warnings []LintWarning

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		warnings = append(warnings, lintEntry(e)...)
	}

	return LintResult{
		Warnings: warnings,
		Total:    len(warnings),
	}
}

func lintEntry(e Entry) []LintWarning {
	var ws []LintWarning

	if e.Schedule == "* * * * *" {
		ws = append(ws, LintWarning{
			Line:     e.Line,
			Schedule: e.Schedule,
			Message:  "schedule runs every minute — ensure this is intentional",
			Severity: SeverityWarn,
		})
	}

	if containsWildcardHeavy(e.Schedule) {
		ws = append(ws, LintWarning{
			Line:     e.Line,
			Schedule: e.Schedule,
			Message:  "schedule uses wildcards in all time fields — may run more often than intended",
			Severity: SeverityInfo,
		})
	}

	if e.Command == "" {
		ws = append(ws, LintWarning{
			Line:     e.Line,
			Schedule: e.Schedule,
			Message:  "entry has an empty command",
			Severity: SeverityError,
		})
	}

	return ws
}

// containsWildcardHeavy returns true when the first four schedule fields are all '*'.
func containsWildcardHeavy(schedule string) bool {
	var fields [5]string
	n, _ := fmt.Sscanf(schedule, "%s %s %s %s %s",
		&fields[0], &fields[1], &fields[2], &fields[3], &fields[4])
	if n < 4 {
		return false
	}
	// Only flag when minute is NOT '*' but all others are — i.e. very broad
	return fields[1] == "*" && fields[2] == "*" && fields[3] == "*" && fields[4] == "*"
}
