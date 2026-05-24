package reporter_test

import (
	"fmt"

	"github.com/cronaudit/cronaudit/internal/reporter"
)

func ExampleAssessRisk() {
	r := reporter.Report{
		Source: "crontab",
		Entries: []reporter.Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "/usr/bin/cleanup.sh", Valid: true},
			{Line: 2, Schedule: "*/2 * * * *", Command: "sudo rm -rf /tmp/old", Valid: true},
		},
	}
	rr := reporter.AssessRisk(r)
	for _, e := range rr.Entries {
		fmt.Printf("line %d: %s\n", e.Line, e.Level)
	}
	// Output:
	// line 1: low
	// line 2: high
}

func ExampleFormatRiskReport() {
	r := reporter.Report{
		Source: "crontab",
		Entries: []reporter.Entry{
			{Line: 1, Schedule: "0 6 * * *", Command: "/usr/bin/report.sh", Valid: true},
		},
	}
	rr := reporter.AssessRisk(r)
	out := reporter.FormatRiskReport(rr)
	if len(out) > 0 {
		fmt.Println("report generated")
	}
	// Output:
	// report generated
}
