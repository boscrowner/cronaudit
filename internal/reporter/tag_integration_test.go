package reporter_test

import (
	"testing"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/reporter"
)

// TestTagEntries_WithCustomTag verifies that user-defined Tag rules work
// alongside CommonTags without interfering with each other.
func TestTagEntries_WithCustomTag(t *testing.T) {
	customTag := reporter.Tag{
		Name: "root-job",
		Match: func(e parser.Entry) bool {
			return len(e.Command) > 0 && e.Command[0] == '/'
		},
	}

	report := reporter.Report{
		Source: "integration",
		Entries: []parser.Entry{
			{Valid: true, Minute: "0", Hour: "1", Day: "*", Month: "*", Weekday: "*", Command: "/usr/bin/backup"},
			{Valid: true, Minute: "0", Hour: "2", Day: "*", Month: "*", Weekday: "*", Command: "relative-job"},
		},
	}

	allTags := append(reporter.CommonTags(), customTag)
	tr := reporter.TagEntries(report, allTags)

	rootJobs := reporter.FilterByTag(tr, "root-job")
	if len(rootJobs) != 1 {
		t.Fatalf("expected 1 root-job, got %d", len(rootJobs))
	}
	if rootJobs[0].Entry.Command != "/usr/bin/backup" {
		t.Errorf("unexpected command: %s", rootJobs[0].Entry.Command)
	}

	// Both entries should also carry the "daily" tag.
	daily := reporter.FilterByTag(tr, "daily")
	if len(daily) != 2 {
		t.Errorf("expected 2 daily entries, got %d", len(daily))
	}
}

// TestTagEntries_EmptyReport returns an empty TagReport without panicking.
func TestTagEntries_EmptyReport(t *testing.T) {
	report := reporter.Report{Source: "empty"}
	tr := reporter.TagEntries(report, reporter.CommonTags())
	if len(tr.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(tr.Results))
	}
	if tr.Source != "empty" {
		t.Errorf("expected source 'empty', got %q", tr.Source)
	}
}
