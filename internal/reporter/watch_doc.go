// Package reporter provides tools for building, analysing, and comparing
// crontab reports.
//
// # Watch
//
// The watch sub-feature compares two [Snapshot] values produced at different
// points in time and surfaces a [WatchResult] that describes what changed.
//
// Typical usage:
//
//	before := reporter.TakeSnapshot(reportA, "baseline")
//	// ... time passes, crontab is re-parsed ...
//	after := reporter.TakeSnapshot(reportB, "current")
//
//	result := reporter.WatchSnapshots(before, after)
//	if result.Changed {
//		fmt.Println(reporter.WatchSummary(result))
//	}
package reporter
