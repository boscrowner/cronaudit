// Package reporter assembles parsed cron entries and their summaries
// into a structured Report suitable for rendering or further processing.
package reporter

import (
	"github.com/yourorg/cronaudit/internal/parser"
)

// Entry holds a single cron line together with its parsed result and
// human-readable summary.
type Entry struct {
	Line    string
	Valid   bool
	Fields  *parser.CronEntry
	Summary string
	Error   string
}

// Report is the top-level structure produced for a crontab source.
type Report struct {
	Source  string
	Entries []Entry
	Stats   Stats
}

// Stats holds aggregate counts for the report.
type Stats struct {
	Total   int
	Valid   int
	Invalid int
}

// Build constructs a Report from a slice of raw crontab lines and a source
// label (e.g. a filename or "stdin").
func Build(source string, lines []string) Report {
	report := Report{Source: source}

	for _, line := range lines {
		entry := processLine(line)
		report.Entries = append(report.Entries, entry)
		report.Stats.Total++
		if entry.Valid {
			report.Stats.Valid++
		} else {
			report.Stats.Invalid++
		}
	}

	return report
}

func processLine(line string) Entry {
	cronEntry, err := parser.Parse(line)
	if err != nil {
		return Entry{
			Line:  line,
			Valid: false,
			Error: err.Error(),
		}
	}

	if cronEntry == nil {
		// blank line or comment — skip silently
		return Entry{Line: line, Valid: true}
	}

	return Entry{
		Line:    line,
		Valid:   true,
		Fields:  cronEntry,
		Summary: parser.Summarize(cronEntry),
	}
}
