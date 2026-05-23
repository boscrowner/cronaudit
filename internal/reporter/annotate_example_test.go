package reporter_test

import (
	"fmt"

	"github.com/cronaudit/cronaudit/internal/reporter"
)

func ExampleAnnotateEntries() {
	r := reporter.Report{
		Source: "crontab",
		Entries: []reporter.Entry{
			{Line: 1, Raw: "0 3 * * * /scripts/backup.sh", Valid: true, Summary: "At 03:00 daily"},
		},
	}
	tags := map[string]string{
		"0 3 * * * /scripts/backup.sh": "nightly-backup",
	}
	out := reporter.AnnotateEntries(r, tags)
	fmt.Println(out.Entries[0].Summary)
	// Output: [nightly-backup] At 03:00 daily
}

func ExampleAnnotateErrors() {
	r := reporter.Report{
		Source: "crontab",
		Entries: []reporter.Entry{
			{Line: 2, Raw: "bad line", Valid: false, Error: "invalid field count"},
		},
	}
	out := reporter.AnnotateErrors(r)
	fmt.Println(out.Entries[0].Summary)
	// Output: [ERROR: invalid field count]
}
