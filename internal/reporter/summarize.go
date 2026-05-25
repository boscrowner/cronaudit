package reporter

import (
	"fmt"
	"strings"

	"github.com/example/cronaudit/internal/parser"
)

// SummaryReport holds a human-readable summary for each entry in a Report.
type SummaryReport struct {
	Source  string
	Entries []SummaryEntry
}

// SummaryEntry pairs a cron entry with its plain-English description.
type SummaryEntry struct {
	Line    int
	Raw     string
	Valid   bool
	Summary string
	Error   string
}

// SummarizeReport produces a SummaryReport from a Report, attaching
// human-readable descriptions to every valid entry and error messages
// to every invalid one.
func SummarizeReport(r Report) SummaryReport {
	entries := make([]SummaryEntry, 0, len(r.Entries))
	for _, e := range r.Entries {
		se := SummaryEntry{
			Line:  e.Line,
			Raw:   e.Raw,
			Valid: e.Valid,
		}
		if e.Valid {
			se.Summary = parser.Summarize(e.Entry)
		} else {
			se.Error = formatSummaryError(e.Error)
		}
		entries = append(entries, se)
	}
	return SummaryReport{
		Source:  r.Source,
		Entries: entries,
	}
}

// FormatSummaryReport renders a SummaryReport as a multi-line string
// suitable for display in a terminal.
func FormatSummaryReport(sr SummaryReport) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Source: %s\n", sr.Source)
	fmt.Fprintf(&sb, "Entries: %d\n\n", len(sr.Entries))
	for _, e := range sr.Entries {
		if e.Valid {
			fmt.Fprintf(&sb, "  [line %d] %s\n", e.Line, e.Raw)
			fmt.Fprintf(&sb, "           → %s\n", e.Summary)
		} else {
			fmt.Fprintf(&sb, "  [line %d] %s\n", e.Line, e.Raw)
			fmt.Fprintf(&sb, "           ✗ %s\n", e.Error)
		}
	}
	return sb.String()
}

// Stats returns the count of valid and invalid entries in the SummaryReport.
func (sr SummaryReport) Stats() (valid, invalid int) {
	for _, e := range sr.Entries {
		if e.Valid {
			valid++
		} else {
			invalid++
		}
	}
	return valid, invalid
}

// InvalidEntries returns a slice containing only the invalid SummaryEntries.
// This is useful when callers want to report or log only the problematic lines.
func (sr SummaryReport) InvalidEntries() []SummaryEntry {
	var invalid []SummaryEntry
	for _, e := range sr.Entries {
		if !e.Valid {
			invalid = append(invalid, e)
		}
	}
	return invalid
}

func formatSummaryError(err string) string {
	if err == "" {
		return "invalid entry"
	}
	return err
}
