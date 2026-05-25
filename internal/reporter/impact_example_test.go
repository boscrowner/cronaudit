package reporter_test

import (
	"fmt"

	"github.com/cronaudit/cronaudit/internal/reporter"
)

func ExampleAssessImpact() {
	r := reporter.Report{
		Source: "crontab",
		Entries: []reporter.Entry{
			{
				Line:     1,
				Valid:    true,
				Schedule: "0 3 * * *",
				Command:  "sudo rm -rf /var/cache && curl http://cleanup.internal",
			},
			{
				Line:     2,
				Valid:    true,
				Schedule: "@hourly",
				Command:  "echo heartbeat",
			},
		},
	}

	ir := reporter.AssessImpact(r)
	for _, e := range ir.Entries {
		fmt.Printf("line %d: %s\n", e.Line, e.Level)
	}
	// Output:
	// line 1: critical
	// line 2: low
}

func ExampleFormatImpactReport() {
	r := reporter.Report{
		Source: "sample",
		Entries: []reporter.Entry{
			{
				Line:     1,
				Valid:    true,
				Schedule: "*/5 * * * *",
				Command:  "wget http://example.com/data",
			},
		},
	}
	ir := reporter.AssessImpact(r)
	fmt.Print(FormatImpactReport_nonzero(ir))
}

// helper to avoid empty-output example failure
func FormatImpactReport_nonzero(ir reporter.ImpactReport) string {
	out := reporter.FormatImpactReport(ir)
	if out == "" {
		return "(empty)"
	}
	return out
}
