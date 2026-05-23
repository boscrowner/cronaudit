package reporter_test

import (
	"fmt"
	"time"

	"github.com/example/cronaudit/internal/reporter"
)

func ExampleBuildTimeline() {
	r := reporter.Report{
		Source: "example",
		Entries: []reporter.Entry{
			{Line: 1, Schedule: "0 9 * * 1-5", Command: "/usr/bin/backup", Valid: true},
			{Line: 2, Schedule: "*/15 * * * *", Command: "/usr/bin/healthcheck", Valid: true},
		},
	}

	// Project over a 1-hour window starting at a known Monday 09:00
	from := time.Date(2024, 1, 8, 8, 45, 0, 0, time.UTC) // Monday
	to := from.Add(time.Hour)

	tl := reporter.BuildTimeline(r, from, to, 3)

	fmt.Printf("Timeline for: %s\n", tl.Source)
	for _, e := range tl.Entries {
		if len(e.NextRuns) > 0 {
			fmt.Printf("  %s → first run at %s\n", e.Command, e.NextRuns[0].Format("15:04"))
		}
	}

	// Output:
	// Timeline for: example
	//   /usr/bin/healthcheck → first run at 09:00
	//   /usr/bin/backup → first run at 09:00
}
