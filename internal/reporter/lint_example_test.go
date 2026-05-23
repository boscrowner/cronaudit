package reporter_test

import (
	"fmt"

	"github.com/cronaudit/internal/reporter"
)

func ExampleLintReport() {
	r := reporter.Report{
		Source: "example-crontab",
		Entries: []reporter.Entry{
			{
				Line:     1,
				Schedule: "* * * * *",
				Command:  "/usr/bin/check",
				Valid:    true,
			},
			{
				Line:     2,
				Schedule: "0 9 * * 1-5",
				Command:  "/usr/bin/report",
				Valid:    true,
			},
		},
	}

	result := reporter.LintReport(r)
	for _, w := range result.Warnings {
		fmt.Printf("[%s] line %d: %s\n", w.Severity, w.Line, w.Message)
	}
	fmt.Printf("total warnings: %d\n", result.Total)

	// Output:
	// [warn] line 1: schedule runs every minute — ensure this is intentional
	// total warnings: 1
}
