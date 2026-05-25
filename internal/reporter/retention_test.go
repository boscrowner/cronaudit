package reporter

import (
	"strings"
	"testing"
)

func buildRetentionReport() Report {
	return Report{
		Source: "retention_test",
		Entries: []Entry{
			{Line: 1, Schedule: "* * * * *", Command: "/bin/every-minute", Valid: true},
			{Line: 2, Schedule: "30 * * * *", Command: "/bin/hourly", Valid: true},
			{Line: 3, Schedule: "0 2 * * *", Command: "/bin/daily", Valid: true},
			{Line: 4, Schedule: "0 3 * * 0", Command: "/bin/weekly", Valid: true},
			{Line: 5, Schedule: "0 4 1 * *", Command: "/bin/monthly", Valid: true},
			{Line: 6, Schedule: "0 5 1 1 *", Command: "/bin/yearly", Valid: true},
			{Line: 7, Schedule: "bad entry", Command: "", Valid: false, Error: "invalid"},
		},
	}
}

func TestBuildRetentionReport_SkipsInvalidEntries(t *testing.T) {
	r := buildRetentionReport()
	rr := BuildRetentionReport(r)

	for _, re := range rr.Entries {
		if !re.Entry.Valid {
			t.Errorf("expected only valid entries, got invalid entry: %v", re.Entry)
		}
	}
}

func TestBuildRetentionReport_SourcePreserved(t *testing.T) {
	r := buildRetentionReport()
	rr := BuildRetentionReport(r)

	if rr.Source != r.Source {
		t.Errorf("expected source %q, got %q", r.Source, rr.Source)
	}
}

func TestBuildRetentionReport_EveryMinute(t *testing.T) {
	r := buildRetentionReport()
	rr := BuildRetentionReport(r)

	var found *RetentionEntry
	for i := range rr.Entries {
		if rr.Entries[i].Entry.Schedule == "* * * * *" {
			found = &rr.Entries[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected entry for '* * * * *'")
	}
	if found.Frequency != "every minute" {
		t.Errorf("expected frequency 'every minute', got %q", found.Frequency)
	}
	if found.WindowDays != 7 {
		t.Errorf("expected window 7, got %d", found.WindowDays)
	}
}

func TestBuildRetentionReport_Daily(t *testing.T) {
	r := buildRetentionReport()
	rr := BuildRetentionReport(r)

	var found *RetentionEntry
	for i := range rr.Entries {
		if rr.Entries[i].Entry.Schedule == "0 2 * * *" {
			found = &rr.Entries[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected entry for '0 2 * * *'")
	}
	if found.Frequency != "daily" {
		t.Errorf("expected frequency 'daily', got %q", found.Frequency)
	}
	if found.WindowDays != 90 {
		t.Errorf("expected window 90, got %d", found.WindowDays)
	}
}

func TestFormatRetentionReport_EmptyReport(t *testing.T) {
	rr := RetentionReport{Source: "empty"}
	out := FormatRetentionReport(rr)
	if !strings.Contains(out, "No valid entries") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestFormatRetentionReport_ContainsSchedule(t *testing.T) {
	r := buildRetentionReport()
	rr := BuildRetentionReport(r)
	out := FormatRetentionReport(rr)

	if !strings.Contains(out, "* * * * *") {
		t.Errorf("expected schedule in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Suggested retention") {
		t.Errorf("expected retention label in output, got:\n%s", out)
	}
}
