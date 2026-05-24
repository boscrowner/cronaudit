package reporter

import (
	"strings"
	"testing"
)

func buildOverlapReport(entries []Entry) Report {
	return Report{
		Source:  "test",
		Entries: entries,
	}
}

func TestDetectOverlaps_FindsCommonSchedule(t *testing.T) {
	r := buildOverlapReport([]Entry{
		{Line: 1, Schedule: "0 9 * * *", Command: "cmd-a", Valid: true},
		{Line: 2, Schedule: "0 9 * * *", Command: "cmd-b", Valid: true},
	})
	out := DetectOverlaps(r)
	if len(out.Overlaps) != 1 {
		t.Fatalf("expected 1 overlap, got %d", len(out.Overlaps))
	}
	if out.Overlaps[0].A.Command != "cmd-a" || out.Overlaps[0].B.Command != "cmd-b" {
		t.Errorf("unexpected overlap pair: %+v", out.Overlaps[0])
	}
}

func TestDetectOverlaps_ExampleString(t *testing.T) {
	r := buildOverlapReport([]Entry{
		{Line: 1, Schedule: "30 8 * * *", Command: "a", Valid: true},
		{Line: 2, Schedule: "30 8 * * *", Command: "b", Valid: true},
	})
	out := DetectOverlaps(r)
	if len(out.Overlaps) == 0 {
		t.Fatal("expected overlap")
	}
	if !strings.Contains(out.Overlaps[0].Example, "hour=08") {
		t.Errorf("example should mention hour=08, got: %s", out.Overlaps[0].Example)
	}
}

func TestDetectOverlaps_NoOverlapForDifferentSchedules(t *testing.T) {
	r := buildOverlapReport([]Entry{
		{Line: 1, Schedule: "0 6 * * 1", Command: "mon-only", Valid: true},
		{Line: 2, Schedule: "0 6 * * 3", Command: "wed-only", Valid: true},
	})
	out := DetectOverlaps(r)
	if len(out.Overlaps) != 0 {
		t.Errorf("expected no overlaps, got %d", len(out.Overlaps))
	}
}

func TestDetectOverlaps_SkipsInvalidEntries(t *testing.T) {
	r := buildOverlapReport([]Entry{
		{Line: 1, Schedule: "0 9 * * *", Command: "valid", Valid: true},
		{Line: 2, Schedule: "bad schedule", Command: "invalid", Valid: false, Error: "parse error"},
	})
	out := DetectOverlaps(r)
	if len(out.Overlaps) != 0 {
		t.Errorf("expected no overlaps involving invalid entry, got %d", len(out.Overlaps))
	}
}

func TestDetectOverlaps_SourcePreserved(t *testing.T) {
	r := buildOverlapReport(nil)
	r.Source = "myfile.crontab"
	out := DetectOverlaps(r)
	if out.Source != "myfile.crontab" {
		t.Errorf("expected source %q, got %q", "myfile.crontab", out.Source)
	}
}

func TestDetectOverlaps_EmptyReport(t *testing.T) {
	out := DetectOverlaps(Report{})
	if len(out.Overlaps) != 0 {
		t.Errorf("expected no overlaps for empty report")
	}
}
