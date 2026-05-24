package reporter

import (
	"testing"
)

func buildSortTestReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Raw: "*/5 * * * * /bin/c", Summary: "every 5 minutes", Error: ""},
			{Raw: "0 2 * * * /bin/a", Summary: "at 02:00", Error: ""},
			{Raw: "bad entry", Summary: "", Error: "invalid format"},
			{Raw: "30 8 * * 1 /bin/b", Summary: "mondays at 08:30", Error: ""},
			{Raw: "also bad", Summary: "", Error: "invalid format"},
		},
		Stats: Stats{Total: 5, Valid: 3, Invalid: 2},
	}
}

func TestSortEntries_ByLine_PreservesOrder(t *testing.T) {
	r := buildSortTestReport()
	result := SortEntries(r, SortByLine)

	if len(result) != len(r.Entries) {
		t.Fatalf("expected %d entries, got %d", len(r.Entries), len(result))
	}
	for i, e := range result {
		if e.Raw != r.Entries[i].Raw {
			t.Errorf("entry %d: expected Raw %q, got %q", i, r.Entries[i].Raw, e.Raw)
		}
	}
}

func TestSortEntries_BySchedule_SortsLexicographically(t *testing.T) {
	r := buildSortTestReport()
	result := SortEntries(r, SortBySchedule)

	for i := 1; i < len(result); i++ {
		if result[i].Raw < result[i-1].Raw {
			t.Errorf("out of order at index %d: %q before %q", i, result[i-1].Raw, result[i].Raw)
		}
	}
}

func TestSortEntries_ByValid_InvalidFirst(t *testing.T) {
	r := buildSortTestReport()
	result := SortEntries(r, SortByValid)

	seenValid := false
	for _, e := range result {
		if e.Error == "" {
			seenValid = true
		}
		if seenValid && e.Error != "" {
			t.Errorf("invalid entry %q appeared after a valid entry", e.Raw)
		}
	}
}

func TestSortEntries_ByValid_CountsPreserved(t *testing.T) {
	r := buildSortTestReport()
	result := SortEntries(r, SortByValid)

	var validCount, invalidCount int
	for _, e := range result {
		if e.Error == "" {
			validCount++
		} else {
			invalidCount++
		}
	}

	if validCount != r.Stats.Valid {
		t.Errorf("expected %d valid entries, got %d", r.Stats.Valid, validCount)
	}
	if invalidCount != r.Stats.Invalid {
		t.Errorf("expected %d invalid entries, got %d", r.Stats.Invalid, invalidCount)
	}
}

func TestSortEntries_DoesNotMutateOriginal(t *testing.T) {
	r := buildSortTestReport()
	originalRaw := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		originalRaw[i] = e.Raw
	}

	_ = SortEntries(r, SortBySchedule)

	for i, e := range r.Entries {
		if e.Raw != originalRaw[i] {
			t.Errorf("original entry %d mutated: expected %q, got %q", i, originalRaw[i], e.Raw)
		}
	}
}
