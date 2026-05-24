package reporter_test

import (
	"fmt"

	"github.com/example/cronaudit/internal/reporter"
)

func ExampleDetectAnomalies() {
	r := reporter.Report{
		Source: "example",
		Entries: []reporter.Entry{
			{Line: 1, Raw: "* * * * * /bin/every-minute", Valid: true},
			{Line: 2, Raw: "0 12 * * * /bin/noon-job", Valid: true},
		},
	}

	ar := reporter.DetectAnomalies(r)
	fmt.Println(len(ar.Anomalies) > 0)
	fmt.Println(ar.Anomalies[0].Kind)
	// Output:
	// true
	// high-frequency
}

func ExampleFormatAnomalyReport_clean() {
	ar := reporter.AnomalyReport{
		Source:    "crontab",
		Anomalies: nil,
	}
	out := reporter.FormatAnomalyReport(ar)
	fmt.Print(out)
	// Output:
	// source: crontab
	// No anomalies detected.
}
