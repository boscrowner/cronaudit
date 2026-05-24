// Package reporter provides tools for parsing, analyzing, and reporting
// on crontab entries.
//
// # Compare
//
// CompareReport compares two reports side by side, matching entries by
// their schedule string. For each schedule present in both reports it
// returns a CompareResult that includes both commands, both summaries,
// and a boolean Match flag indicating whether the commands are identical.
//
// UnmatchedEntries returns the valid entries from report A whose schedules
// do not appear in report B, making it easy to spot additions or removals
// between two crontab snapshots.
//
// Example:
//
//	results := reporter.CompareReport(reportA, reportB)
//	for _, r := range results {
//	    if !r.Match {
//	        fmt.Printf("schedule %s changed: %s -> %s\n", r.Schedule, r.CommandA, r.CommandB)
//	    }
//	}
package reporter
