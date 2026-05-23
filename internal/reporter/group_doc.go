// Package reporter provides tools for building, filtering, sorting,
// grouping, and exporting crontab audit reports.
//
// The group.go file provides grouping utilities:
//
//   - GroupByCommand groups report entries by their command string,
//     making it easy to spot commands scheduled multiple times.
//
//   - GroupBySchedule groups entries by their cron schedule expression,
//     useful for identifying schedules shared across multiple commands.
//
// Invalid entries (those that failed parsing) are grouped under the
// key "(invalid)" in both grouping modes.
package reporter
