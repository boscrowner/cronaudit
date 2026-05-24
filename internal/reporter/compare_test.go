package reporter

import (
	"testing"
)

func buildCompareReport(entries []Entry) Report {
	return Report{
		Source:  "test",
		Entries: entries,
	}
}

func TestCompareReport_MatchingSchedules(t *testing.T) {
	a := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true, Summary: "every hour"},
		{Line: 2, Schedule: "0 0 * * *", Command: "/bin/bar", Valid: true, Summary: "daily"},
	})
	b := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true, Summary: "every hour"},
		{Line: 2, Schedule: "0 0 * * *", Command: "/bin/baz", Valid: true, Summary: "daily"},
	})

	results := CompareReport(a, b)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestCompareReport_MatchFlagTrue(t *testing.T) {
	a := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
	})
	b := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
	})

	results := CompareReport(a, b)
	if len(results) != 1 || !results[0].Match {
		t.Error("expected Match=true for identical commands")
	}
}

func TestCompareReport_MatchFlagFalse(t *testing.T) {
	a := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
	})
	b := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/bar", Valid: true},
	})

	results := CompareReport(a, b)
	if len(results) != 1 || results[0].Match {
		t.Error("expected Match=false for different commands")
	}
}

func TestCompareReport_SkipsInvalidEntries(t *testing.T) {
	a := buildCompareReport([]Entry{
		{Line: 1, Schedule: "bad", Command: "", Valid: false},
	})
	b := buildCompareReport([]Entry{
		{Line: 1, Schedule: "bad", Command: "", Valid: false},
	})

	results := CompareReport(a, b)
	if len(results) != 0 {
		t.Errorf("expected 0 results for invalid entries, got %d", len(results))
	}
}

func TestUnmatchedEntries_ReturnsEntriesOnlyInA(t *testing.T) {
	a := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
		{Line: 2, Schedule: "0 0 * * *", Command: "/bin/bar", Valid: true},
	})
	b := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
	})

	unmatched := UnmatchedEntries(a, b)
	if len(unmatched) != 1 {
		t.Fatalf("expected 1 unmatched entry, got %d", len(unmatched))
	}
	if unmatched[0].Schedule != "0 0 * * *" {
		t.Errorf("unexpected schedule: %s", unmatched[0].Schedule)
	}
}

func TestUnmatchedEntries_EmptyWhenAllMatch(t *testing.T) {
	a := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
	})
	b := buildCompareReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
	})

	unmatched := UnmatchedEntries(a, b)
	if len(unmatched) != 0 {
		t.Errorf("expected 0 unmatched entries, got %d", len(unmatched))
	}
}
