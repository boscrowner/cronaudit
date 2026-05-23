package reporter_test

import (
	"fmt"

	"github.com/yourorg/cronaudit/internal/reporter"
)

func ExampleClusterByHour() {
	r := reporter.Report{
		Source: "example",
		Entries: []reporter.ReportEntry{
			{Line: 1, Schedule: "0 3 * * *", Command: "backup.sh", Valid: true},
			{Line: 2, Schedule: "0 3 * * *", Command: "cleanup.sh", Valid: true},
			{Line: 3, Schedule: "15 9 * * *", Command: "report.sh", Valid: true},
		},
	}

	clusters := reporter.ClusterByHour(r)
	for _, c := range clusters {
		fmt.Printf("pattern=%q size=%d\n", c.Pattern, c.Size)
	}

	// Output:
	// pattern="0 3" size=2
	// pattern="15 9" size=1
}

func ExampleClusterByDayOfWeek() {
	r := reporter.Report{
		Source: "example",
		Entries: []reporter.ReportEntry{
			{Line: 1, Schedule: "0 1 * * 0", Command: "weekly-a.sh", Valid: true},
			{Line: 2, Schedule: "0 2 * * 0", Command: "weekly-b.sh", Valid: true},
			{Line: 3, Schedule: "0 6 * * 1", Command: "monday.sh", Valid: true},
		},
	}

	clusters := reporter.ClusterByDayOfWeek(r)
	for _, c := range clusters {
		fmt.Printf("weekday=%q count=%d\n", c.Pattern, c.Size)
	}

	// Output:
	// weekday="0" count=2
	// weekday="1" count=1
}
