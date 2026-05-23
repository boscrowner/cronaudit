// Package reporter builds structured reports from parsed crontab lines.
//
// A Report contains per-entry results (raw text, human-readable summary,
// and any validation error) together with aggregate statistics about the
// crontab as a whole.
//
// # Building a report
//
//	lines  := []string{"0 * * * * /usr/bin/backup", "# comment", "bad line"}
//	report := reporter.Build("mycrontab", lines)
//
// # Filtering entries
//
// Use FilterValid or FilterInvalid to obtain a focused view of the report
// without modifying the original:
//
//	valid   := reporter.FilterValid(report)
//	invalid := reporter.FilterInvalid(report)
//
// # Sorting entries
//
// Use SortEntries to reorder entries for display purposes. The original
// report is never mutated:
//
//	sorted := reporter.SortEntries(report, reporter.SortBySchedule)
package reporter
