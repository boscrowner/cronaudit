package reporter

import (
	"strings"
	"testing"
)

func buildFrequencyReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Raw: "0 9 * * 1 /bin/backup", Valid: true},
			{Raw: "0 9 * * 1 /bin/sync", Valid: true},
			{Raw: "*/5 * * * * /bin/poll", Valid: true},
			{Raw: "0 0 * * * /bin/nightly", Valid: true},
			{Raw: "bad entry", Valid: false},
		},
	}
}

func TestBuildFrequencyReport_SkipsInvalidEntries(t *testing.T) {
	r := buildFrequencyReport()
	fr := BuildFrequencyReport(r)
	for _, e := range fr.Entries {
		if e.Count == 0 {
			t.Errorf("expected non-zero count for schedule %q", e.Schedule)
		}
	}
}

func TestBuildFrequencyReport_CountsCorrectly(t *testing.T) {
	r := buildFrequencyReport()
	fr := BuildFrequencyReport(r)

	for _, e := range fr.Entries {
		if e.Schedule == "0 9 * * 1" {
			if e.Count != 2 {
				t.Errorf("expected count 2 for '0 9 * * 1', got %d", e.Count)
			}
			return
		}
	}
	t.Error("expected to find schedule '0 9 * * 1'")
}

func TestBuildFrequencyReport_SortedByCountDescending(t *testing.T) {
	r := buildFrequencyReport()
	fr := BuildFrequencyReport(r)

	for i := 1; i < len(fr.Entries); i++ {
		if fr.Entries[i].Count > fr.Entries[i-1].Count {
			t.Errorf("entries not sorted by count descending at index %d", i)
		}
	}
}

func TestBuildFrequencyReport_CommandsPopulated(t *testing.T) {
	r := buildFrequencyReport()
	fr := BuildFrequencyReport(r)

	for _, e := range fr.Entries {
		if e.Schedule == "0 9 * * 1" {
			if len(e.Commands) != 2 {
				t.Errorf("expected 2 commands, got %d", len(e.Commands))
			}
			return
		}
	}
}

func TestBuildFrequencyReport_SourcePreserved(t *testing.T) {
	r := buildFrequencyReport()
	fr := BuildFrequencyReport(r)
	if fr.Source != "test" {
		t.Errorf("expected source 'test', got %q", fr.Source)
	}
}

func TestFormatFrequencyReport_ContainsSchedule(t *testing.T) {
	r := buildFrequencyReport()
	fr := BuildFrequencyReport(r)
	out := FormatFrequencyReport(fr)
	if !strings.Contains(out, "0 9 * * 1") {
		t.Errorf("expected output to contain schedule, got: %s", out)
	}
}

func TestFormatFrequencyReport_EmptyReport(t *testing.T) {
	fr := FrequencyReport{Source: "empty"}
	out := FormatFrequencyReport(fr)
	if !strings.Contains(out, "No valid") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
