package reporter

import (
	"strings"
	"testing"
)

func buildHighlightReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Raw: "0 * * * *", Command: "/usr/bin/backup.sh", Valid: true},
			{Line: 2, Raw: "*/5 * * * *", Command: "/usr/bin/cleanup", Valid: true},
			{Line: 3, Raw: "bad entry", Command: "", Valid: false, Error: "invalid"},
			{Line: 4, Raw: "0 0 * * 0", Command: "/opt/weekly.sh", Valid: true},
		},
	}
}

func TestHighlightReport_MatchesSchedule(t *testing.T) {
	r := buildHighlightReport()
	results := HighlightReport(r, "*/5")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Line != 2 {
		t.Errorf("expected line 2, got %d", results[0].Line)
	}
}

func TestHighlightReport_MatchesCommand(t *testing.T) {
	r := buildHighlightReport()
	results := HighlightReport(r, "backup")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0].Snippet, "[backup") {
		t.Errorf("expected snippet to mark match, got: %s", results[0].Snippet)
	}
}

func TestHighlightReport_SkipsInvalidEntries(t *testing.T) {
	r := buildHighlightReport()
	results := HighlightReport(r, "bad")
	if len(results) != 0 {
		t.Errorf("expected no results for invalid entries, got %d", len(results))
	}
}

func TestHighlightReport_EmptyQuery(t *testing.T) {
	r := buildHighlightReport()
	results := HighlightReport(r, "")
	if results != nil {
		t.Errorf("expected nil for empty query, got %v", results)
	}
}

func TestHighlightReport_CaseInsensitive(t *testing.T) {
	r := buildHighlightReport()
	results := HighlightReport(r, "WEEKLY")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Line != 4 {
		t.Errorf("expected line 4, got %d", results[0].Line)
	}
}

func TestHighlightReport_SnippetContainsField(t *testing.T) {
	r := buildHighlightReport()
	results := HighlightReport(r, "0 0")
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if !strings.HasPrefix(results[0].Snippet, "schedule:") {
		t.Errorf("expected snippet to start with 'schedule:', got: %s", results[0].Snippet)
	}
}
