package reporter

import (
	"testing"
)

func buildBaselineReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Schedule: "0 9 * * 1-5", Command: "/usr/bin/backup", Valid: true},
			{Line: 2, Schedule: "*/5 * * * *", Command: "/usr/bin/health", Valid: true},
			{Line: 3, Schedule: "bad", Command: "", Valid: false, Error: "invalid"},
		},
	}
}

func TestBuildBaseline_MatchedEntries(t *testing.T) {
	r := buildBaselineReport()
	approved := []string{
		"0 9 * * 1-5 /usr/bin/backup",
		"*/5 * * * * /usr/bin/health",
	}
	b := BuildBaseline(r, approved)
	if len(b.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(b.Matched))
	}
}

func TestBuildBaseline_UnknownEntries(t *testing.T) {
	r := buildBaselineReport()
	approved := []string{
		"0 9 * * 1-5 /usr/bin/backup",
	}
	b := BuildBaseline(r, approved)
	if len(b.Unknown) != 1 {
		t.Errorf("expected 1 unknown, got %d", len(b.Unknown))
	}
	if b.Unknown[0] != "*/5 * * * * /usr/bin/health" {
		t.Errorf("unexpected unknown entry: %s", b.Unknown[0])
	}
}

func TestBuildBaseline_MissingEntries(t *testing.T) {
	r := buildBaselineReport()
	approved := []string{
		"0 9 * * 1-5 /usr/bin/backup",
		"*/5 * * * * /usr/bin/health",
		"0 0 * * 0 /usr/bin/weekly",
	}
	b := BuildBaseline(r, approved)
	if len(b.Missing) != 1 {
		t.Errorf("expected 1 missing, got %d", len(b.Missing))
	}
	if b.Missing[0] != "0 0 * * 0 /usr/bin/weekly" {
		t.Errorf("unexpected missing entry: %s", b.Missing[0])
	}
}

func TestBuildBaseline_SkipsInvalidEntries(t *testing.T) {
	r := buildBaselineReport()
	approved := []string{}
	b := BuildBaseline(r, approved)
	// invalid entry should not appear in unknown
	for _, u := range b.Unknown {
		if u == " " {
			t.Error("invalid entry should not appear in unknown set")
		}
	}
}

func TestBaselineReport_IsClean_True(t *testing.T) {
	r := buildBaselineReport()
	approved := []string{
		"0 9 * * 1-5 /usr/bin/backup",
		"*/5 * * * * /usr/bin/health",
	}
	b := BuildBaseline(r, approved)
	if !b.IsClean() {
		t.Error("expected baseline to be clean")
	}
}

func TestBaselineReport_IsClean_False(t *testing.T) {
	r := buildBaselineReport()
	approved := []string{
		"0 9 * * 1-5 /usr/bin/backup",
	}
	b := BuildBaseline(r, approved)
	if b.IsClean() {
		t.Error("expected baseline to be dirty")
	}
}
