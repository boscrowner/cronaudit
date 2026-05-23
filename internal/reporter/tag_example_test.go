package reporter_test

import (
	"fmt"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/reporter"
)

func ExampleTagEntries() {
	report := reporter.Report{
		Source: "example",
		Entries: []parser.Entry{
			{Valid: true, Minute: "*", Hour: "*", Day: "*", Month: "*", Weekday: "*", Command: "echo hi"},
			{Valid: true, Minute: "0", Hour: "3", Day: "*", Month: "*", Weekday: "*", Command: "backup.sh"},
		},
	}

	tr := reporter.TagEntries(report, reporter.CommonTags())
	for _, r := range tr.Results {
		fmt.Printf("%s -> %v\n", r.Entry.Command, r.Tags)
	}
	// Output:
	// echo hi -> [frequent]
	// backup.sh -> [daily]
}

func ExampleFilterByTag() {
	report := reporter.Report{
		Source: "example",
		Entries: []parser.Entry{
			{Valid: true, Minute: "*", Hour: "*", Day: "*", Month: "*", Weekday: "*", Command: "echo hi"},
			{Valid: true, Minute: "0", Hour: "3", Day: "*", Month: "*", Weekday: "*", Command: "backup.sh"},
		},
	}

	tr := reporter.TagEntries(report, reporter.CommonTags())
	daily := reporter.FilterByTag(tr, "daily")
	for _, r := range daily {
		fmt.Println(r.Entry.Command)
	}
	// Output:
	// backup.sh
}
