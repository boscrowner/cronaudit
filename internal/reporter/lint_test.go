package reporter

import (
	"testing"
)

func buildLintReport(entries []Entry) Report {
	return Report{
		Source:  "lint-test",
		Entries: entries,
	}
}

func TestLintReport_NoWarningsForCleanEntries(t *testing.T) {
	r := buildLintReport([]Entry{
		{Line: 1, Schedule: "0 9 * * 1-5", Command: "/usr/bin/backup", Valid: true},
		{Line: 2, Schedule: "30 6 * * *", Command: "/usr/bin/sync", Valid: true},
	})
	result := LintReport(r)
	if result.Total != 0 {
		t.Errorf("expected 0 warnings, got %d", result.Total)
	}
}

func TestLintReport_WarnOnEveryMinute(t *testing.T) {
	r := buildLintReport([]Entry{
		{Line: 3, Schedule: "* * * * *", Command: "/usr/bin/poll", Valid: true},
	})
	result := LintReport(r)
	if result.Total == 0 {
		t.Fatal("expected at least one warning for every-minute schedule")
	}
	found := false
	for _, w := range result.Warnings {
		if w.Severity == SeverityWarn && w.Line == 3 {
			found = true
		}
	}
	if !found {
		t.Error("expected a SeverityWarn warning on line 3")
	}
}

func TestLintReport_SkipsInvalidEntries(t *testing.T) {
	r := buildLintReport([]Entry{
		{Line: 5, Schedule: "bad schedule", Command: "/usr/bin/x", Valid: false},
	})
	result := LintReport(r)
	if result.Total != 0 {
		t.Errorf("expected 0 warnings for invalid entry, got %d", result.Total)
	}
}

func TestLintReport_ErrorOnEmptyCommand(t *testing.T) {
	r := buildLintReport([]Entry{
		{Line: 7, Schedule: "0 12 * * *", Command: "", Valid: true},
	})
	result := LintReport(r)
	var found bool
	for _, w := range result.Warnings {
		if w.Severity == SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected a SeverityError warning for empty command")
	}
}

func TestLintReport_TotalMatchesWarningsLen(t *testing.T) {
	r := buildLintReport([]Entry{
		{Line: 1, Schedule: "* * * * *", Command: "", Valid: true},
	})
	result := LintReport(r)
	if result.Total != len(result.Warnings) {
		t.Errorf("Total %d does not match len(Warnings) %d", result.Total, len(result.Warnings))
	}
}
