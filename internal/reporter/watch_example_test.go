package reporter_test

import (
	"fmt"

	"github.com/example/cronaudit/internal/reporter"
)

func ExampleWatchSnapshots() {
	before := reporter.TakeSnapshot(reporter.Report{
		Source: "crontab",
		Entries: []reporter.Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "/bin/backup", Valid: true},
		},
	}, "v1")

	after := reporter.TakeSnapshot(reporter.Report{
		Source: "crontab",
		Entries: []reporter.Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "/bin/backup", Valid: true},
			{Line: 2, Schedule: "30 6 * * 1", Command: "/bin/report", Valid: true},
		},
	}, "v2")

	result := reporter.WatchSnapshots(before, after)
	fmt.Println(result.Changed)
	fmt.Println(reporter.WatchSummary(result))
	// Output:
	// true
	// 1 entry added.
}

func ExampleWatchSummary_noChange() {
	base := reporter.Report{
		Source: "crontab",
		Entries: []reporter.Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
		},
	}
	before := reporter.TakeSnapshot(base, "a")
	after := reporter.TakeSnapshot(base, "b")
	result := reporter.WatchSnapshots(before, after)
	fmt.Println(reporter.WatchSummary(result))
	// Output:
	// No changes detected between snapshots.
}
