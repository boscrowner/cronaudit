// Package reporter provides tools for building, analyzing, and projecting cron reports.
//
// # Timeline
//
// BuildTimeline projects all valid cron entries in a Report onto a time window,
// computing the next N execution times for each entry within [from, to].
//
// Example usage:
//
//	from := time.Now()
//	to := from.Add(24 * time.Hour)
//	tl := reporter.BuildTimeline(report, from, to, 5)
//	for _, entry := range tl.Entries {
//		fmt.Printf("%s: next run at %s\n", entry.Command, entry.NextRuns[0])
//	}
//
// The timeline entries are sorted by their earliest next run time,
// making it easy to see which jobs fire soonest.
package reporter
