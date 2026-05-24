// Package reporter provides tools for analysing and reporting on parsed crontab
// entries.
//
// # Heatmap
//
// BuildHeatmap produces a 7×24 frequency grid that shows how many valid cron
// entries are scheduled to fire in each (day-of-week, hour) bucket.  It is
// useful for spotting scheduling hot-spots at a glance.
//
// Usage:
//
//	result := reporter.BuildHeatmap(report)
//	fmt.Println(result.Formatted)
//
// The Formatted field contains a plain-text ASCII grid with abbreviated day
// labels on the left and hour numbers (0–23) along the top.  Each cell shows
// the number of jobs scheduled for that slot, or a dot (".") when the slot is
// empty.
//
// Individual cells are also available via result.Cells for programmatic use.
package reporter
