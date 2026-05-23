package reporter

import (
	"testing"
)

func buildDedupeReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "echo a", Valid: true},
			{Line: 2, Schedule: "0 * * * *", Command: "echo b", Valid: true},
			{Line: 3, Schedule: "30 6 * * *", Command: "echo c", Valid: true},
			{Line: 4, Schedule: "30 6 * * *", Command: "echo d", Valid: true},
			{Line: 5, Schedule: "0 0 * * 0", Command: "echo e", Valid: false, Error: "bad field"},
		},
	}
}

func TestDedupeEntries_RemovesDuplicates(t *testing.T) {
	r := buildDedupeReport()
	result := DedupeEntries(r)

	if len(result.Entries) != 3 {
		t.Errorf("expected 3 deduplicated entries, got %d", len(result.Entries))
	}
}

func TestDedupeEntries_KeepsFirstOccurrence(t *testing.T) {
	r := buildDedupeReport()
	result := DedupeEntries(r)

	if result.Entries[0].Command != "echo a" {
		t.Errorf("expected first entry command 'echo a', got %q", result.Entries[0].Command)
	}
	if result.Entries[1].Command != "echo c" {
		t.Errorf("expected second entry command 'echo c', got %q", result.Entries[1].Command)
	}
}

func TestDedupeEntries_DoesNotMutateOriginal(t *testing.T) {
	r := buildDedupeReport()
	origLen := len(r.Entries)
	DedupeEntries(r)

	if len(r.Entries) != origLen {
		t.Errorf("original report mutated: expected %d entries, got %d", origLen, len(r.Entries))
	}
}

func TestDedupeEntries_SourcePreserved(t *testing.T) {
	r := buildDedupeReport()
	result := DedupeEntries(r)

	if result.Source != r.Source {
		t.Errorf("expected source %q, got %q", r.Source, result.Source)
	}
}

func TestDedupeEntries_StatsUpdated(t *testing.T) {
	r := buildDedupeReport()
	result := DedupeEntries(r)

	if result.Stats.Total != len(result.Entries) {
		t.Errorf("stats total %d does not match entry count %d", result.Stats.Total, len(result.Entries))
	}
}

func TestDuplicateSchedules_FindsDuplicates(t *testing.T) {
	r := buildDedupeReport()
	dupes := DuplicateSchedules(r)

	if len(dupes) != 2 {
		t.Errorf("expected 2 duplicate schedules, got %d", len(dupes))
	}
}

func TestDuplicateSchedules_NoDuplicates(t *testing.T) {
	r := Report{
		Entries: []Entry{
			{Schedule: "0 * * * *"},
			{Schedule: "30 6 * * *"},
		},
	}
	dupes := DuplicateSchedules(r)

	if len(dupes) != 0 {
		t.Errorf("expected no duplicates, got %d", len(dupes))
	}
}
