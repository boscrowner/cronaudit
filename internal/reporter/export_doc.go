// Package reporter provides report building, filtering, sorting, statistics,
// and export functionality for cronaudit.
//
// ExportCSV serialises a [Report] into comma-separated values format, useful
// for piping results into spreadsheets or downstream tooling.
//
// CSV column layout:
//
//	Line      – source line number of the crontab entry
//	Schedule  – the five-field cron schedule expression
//	Command   – the command or script being scheduled
//	Valid     – "true" if the schedule parsed without errors, otherwise "false"
//	Summary   – human-readable description of the schedule (empty when invalid)
//	Error     – validation error message (empty when valid)
package reporter
