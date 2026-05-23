package reporter_test

import (
	"testing"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/reporter"
)

func buildTagReport() reporter.Report {
	return reporter.Report{
		Source: "tag_test",
		Entries: []parser.Entry{
			{Valid: true, Minute: "*", Hour: "*", Day: "*", Month: "*", Weekday: "*", Command: "echo every-minute"},
			{Valid: true, Minute: "30", Hour: "6", Day: "*", Month: "*", Weekday: "*", Command: "echo daily"},
			{Valid: true, Minute: "0", Hour: "9", Day: "*", Month: "*", Weekday: "1", Command: "echo weekly"},
			{Valid: true, Minute: "0", Hour: "0", Day: "15", Month: "*", Weekday: "*", Command: "echo monthly"},
			{Valid: false, Raw: "bad entry", Error: "invalid"},
		},
	}
}

func TestTagEntries_AppliesCommonTags(t *testing.T) {
	report := buildTagReport()
	tr := reporter.TagEntries(report, reporter.CommonTags())

	if len(tr.Results) != len(report.Entries) {
		t.Fatalf("expected %d results, got %d", len(report.Entries), len(tr.Results))
	}

	expected := []string{"frequent", "daily", "weekly", "monthly", ""}
	for i, r := range tr.Results {
		if expected[i] == "" {
			if len(r.Tags) != 0 {
				t.Errorf("entry %d: expected no tags, got %v", i, r.Tags)
			}
			continue
		}
		found := false
		for _, tag := range r.Tags {
			if tag == expected[i] {
				found = true
			}
		}
		if !found {
			t.Errorf("entry %d: expected tag %q in %v", i, expected[i], r.Tags)
		}
	}
}

func TestTagEntries_InvalidEntryHasNoTags(t *testing.T) {
	report := buildTagReport()
	tr := reporter.TagEntries(report, reporter.CommonTags())
	last := tr.Results[len(tr.Results)-1]
	if len(last.Tags) != 0 {
		t.Errorf("invalid entry should have no tags, got %v", last.Tags)
	}
}

func TestTagEntries_SourcePreserved(t *testing.T) {
	report := buildTagReport()
	tr := reporter.TagEntries(report, reporter.CommonTags())
	if tr.Source != report.Source {
		t.Errorf("expected source %q, got %q", report.Source, tr.Source)
	}
}

func TestFilterByTag_ReturnsMatchingEntries(t *testing.T) {
	report := buildTagReport()
	tr := reporter.TagEntries(report, reporter.CommonTags())
	results := reporter.FilterByTag(tr, "daily")
	if len(results) != 1 {
		t.Fatalf("expected 1 daily entry, got %d", len(results))
	}
	if results[0].Entry.Command != "echo daily" {
		t.Errorf("unexpected command: %s", results[0].Entry.Command)
	}
}

func TestFilterByTag_EmptyWhenNoMatch(t *testing.T) {
	report := buildTagReport()
	tr := reporter.TagEntries(report, reporter.CommonTags())
	results := reporter.FilterByTag(tr, "nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestTagEntries_DoesNotMutateOriginal(t *testing.T) {
	report := buildTagReport()
	origLen := len(report.Entries)
	reporter.TagEntries(report, reporter.CommonTags())
	if len(report.Entries) != origLen {
		t.Errorf("original report mutated")
	}
}
