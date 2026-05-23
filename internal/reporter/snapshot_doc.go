// Package reporter provides tools for building, analysing, and comparing
// crontab reports.
//
// # Snapshot
//
// The Snapshot type captures a [Report] together with a timestamp and an
// optional human-readable label. Snapshots can be persisted as JSON and later
// reloaded for historical comparison.
//
// Typical usage:
//
//	// Capture and persist
//	snap := reporter.TakeSnapshot(report, "production")
//	reporter.WriteSnapshot(file, snap)
//
//	// Reload and compare
//	prev, _ := reporter.ReadSnapshot(oldFile)
//	curr, _ := reporter.ReadSnapshot(newFile)
//	diff := reporter.DiffSnapshots(prev, curr)
//	fmt.Println("added:", len(diff.Added))
package reporter
