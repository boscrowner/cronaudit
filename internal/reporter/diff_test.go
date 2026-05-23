package reporter

import (
	"testing"
)

func buildDiffReport(entries []Entry) Report {
	return Report{
		Source:  "test",
		Entries: entries,
	}
}

func TestDiffReports_AddedEntries(t *testing.T) {
	before := buildDiffReport([]Entry{
		{Schedule: "* * * * *", Command: "echo hello", Valid: true},
	})
	after := buildDiffReport([]Entry{
		{Schedule: "* * * * *", Command: "echo hello", Valid: true},
		{Schedule: "0 * * * *", Command: "echo new", Valid: true},
	})

	result := DiffReports(before, after)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added entry, got %d", len(result.Added))
	}
	if result.Added[0].Command != "echo new" {
		t.Errorf("unexpected added entry command: %s", result.Added[0].Command)
	}
}

func TestDiffReports_RemovedEntries(t *testing.T) {
	before := buildDiffReport([]Entry{
		{Schedule: "* * * * *", Command: "echo hello", Valid: true},
		{Schedule: "0 0 * * *", Command: "echo gone", Valid: true},
	})
	after := buildDiffReport([]Entry{
		{Schedule: "* * * * *", Command: "echo hello", Valid: true},
	})

	result := DiffReports(before, after)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed entry, got %d", len(result.Removed))
	}
	if result.Removed[0].Command != "echo gone" {
		t.Errorf("unexpected removed entry command: %s", result.Removed[0].Command)
	}
}

func TestDiffReports_KeptEntries(t *testing.T) {
	entry := Entry{Schedule: "* * * * *", Command: "echo hello", Valid: true}
	before := buildDiffReport([]Entry{entry})
	after := buildDiffReport([]Entry{entry})

	result := DiffReports(before, after)

	if len(result.Kept) != 1 {
		t.Fatalf("expected 1 kept entry, got %d", len(result.Kept))
	}
	if len(result.Added) != 0 {
		t.Errorf("expected 0 added, got %d", len(result.Added))
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
}

func TestDiffReports_EmptyReports(t *testing.T) {
	before := buildDiffReport([]Entry{})
	after := buildDiffReport([]Entry{})

	result := DiffReports(before, after)

	if len(result.Added) != 0 || len(result.Removed) != 0 || len(result.Kept) != 0 {
		t.Errorf("expected empty diff for empty reports")
	}
}
