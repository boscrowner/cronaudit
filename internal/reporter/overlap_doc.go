// Package reporter provides tools for analysing and reporting on parsed
// crontab entries.
//
// # Overlap Detection
//
// DetectOverlaps scans all valid entries in a Report and identifies pairs
// whose cron schedules fire at the same minute within a standard week.
// This is useful for spotting unintentional scheduling conflicts where two
// jobs compete for the same resources at the same time.
//
// Example:
//
//	report := reporter.Build(lines)
//	overlaps := reporter.DetectOverlaps(report)
//	for _, o := range overlaps.Overlaps {
//		fmt.Printf("conflict: %q and %q both fire at %s\n",
//			o.A.Command, o.B.Command, o.Example)
//	}
package reporter
