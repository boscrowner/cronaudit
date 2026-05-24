package reporter

import (
	"testing"
)

func buildClassifyReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Schedule: "0 2 * * *", Command: "/usr/bin/backup.sh", Valid: true},
			{Line: 2, Schedule: "0 0 * * *", Command: "/usr/local/bin/log-rotate", Valid: true},
			{Line: 3, Schedule: "*/5 * * * *", Command: "/opt/clean_tmp.sh", Valid: true},
			{Line: 4, Schedule: "0 12 * * *", Command: "/scripts/deploy.sh", Valid: true},
			{Line: 5, Schedule: "0 8 * * 1", Command: "/bin/send_report.py", Valid: true},
			{Line: 6, Schedule: "bad schedule", Command: "", Valid: false, Error: "invalid"},
			{Line: 7, Schedule: "30 6 * * *", Command: "/usr/bin/do_something", Valid: true},
		},
	}
}

func TestClassifyEntries_SourcePreserved(t *testing.T) {
	r := buildClassifyReport()
	cr := ClassifyEntries(r)
	if cr.Source != "test" {
		t.Errorf("expected source 'test', got %q", cr.Source)
	}
}

func TestClassifyEntries_CountMatchesEntries(t *testing.T) {
	r := buildClassifyReport()
	cr := ClassifyEntries(r)
	if len(cr.Classifications) != len(r.Entries) {
		t.Errorf("expected %d classifications, got %d", len(r.Entries), len(cr.Classifications))
	}
}

func TestClassifyEntries_BackupCategory(t *testing.T) {
	r := buildClassifyReport()
	cr := ClassifyEntries(r)
	if cr.Classifications[0].Category != "Backup" {
		t.Errorf("expected 'Backup', got %q", cr.Classifications[0].Category)
	}
}

func TestClassifyEntries_InvalidEntryCategory(t *testing.T) {
	r := buildClassifyReport()
	cr := ClassifyEntries(r)
	// entry at index 5 is invalid
	if cr.Classifications[5].Category != "Invalid" {
		t.Errorf("expected 'Invalid', got %q", cr.Classifications[5].Category)
	}
}

func TestClassifyEntries_GeneralFallback(t *testing.T) {
	r := buildClassifyReport()
	cr := ClassifyEntries(r)
	// last valid entry has no matching keyword
	last := cr.Classifications[6]
	if last.Category != "General" {
		t.Errorf("expected 'General', got %q", last.Category)
	}
}

func TestClassifyEntries_DoesNotMutateOriginal(t *testing.T) {
	r := buildClassifyReport()
	orig := make([]Entry, len(r.Entries))
	copy(orig, r.Entries)
	ClassifyEntries(r)
	for i, e := range r.Entries {
		if e.Command != orig[i].Command {
			t.Errorf("entry %d mutated", i)
		}
	}
}

func TestGroupByCategory_KeysPresent(t *testing.T) {
	r := buildClassifyReport()
	cr := ClassifyEntries(r)
	groups := GroupByCategory(cr)
	expected := []string{"Backup", "Logging", "Cleanup", "Deployment", "Reporting", "Invalid", "General"}
	for _, cat := range expected {
		if _, ok := groups[cat]; !ok {
			t.Errorf("expected category %q in groups", cat)
		}
	}
}

func TestGroupByCategory_InvalidCount(t *testing.T) {
	r := buildClassifyReport()
	cr := ClassifyEntries(r)
	groups := GroupByCategory(cr)
	if len(groups["Invalid"]) != 1 {
		t.Errorf("expected 1 invalid entry, got %d", len(groups["Invalid"]))
	}
}
