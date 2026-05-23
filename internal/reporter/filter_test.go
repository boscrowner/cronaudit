package reporter_test

import (
	"testing"

	"github.com/yourorg/cronaudit/internal/reporter"
)

func buildTestReport() reporter.Report {
	return reporter.Build("test", sampleLines)
}

func TestFilterInvalid_OnlyInvalidEntries(t *testing.T) {
	r := buildTestReport()
	filtered := reporter.FilterInvalid(r)

	for _, e := range filtered.Entries {
		if e.Valid {
			t.Errorf("FilterInvalid returned a valid entry: %q", e.Line)
		}
	}
	if filtered.Stats.Invalid != filtered.Stats.Total {
		t.Errorf("expected all filtered entries to be invalid")
	}
}

func TestFilterValid_OnlyValidEntries(t *testing.T) {
	r := buildTestReport()
	filtered := reporter.FilterValid(r)

	for _, e := range filtered.Entries {
		if !e.Valid {
			t.Errorf("FilterValid returned an invalid entry: %q", e.Line)
		}
		if e.Fields == nil {
			t.Errorf("FilterValid returned an entry with nil Fields: %q", e.Line)
		}
	}
	if filtered.Stats.Valid != filtered.Stats.Total {
		t.Errorf("expected all filtered entries to be valid")
	}
}

func TestFilterInvalid_CountMatchesOriginal(t *testing.T) {
	r := buildTestReport()
	filtered := reporter.FilterInvalid(r)

	if filtered.Stats.Total != r.Stats.Invalid {
		t.Errorf("expected %d invalid entries, got %d", r.Stats.Invalid, filtered.Stats.Total)
	}
}

func TestFilterValid_SourcePreserved(t *testing.T) {
	r := buildTestReport()
	filtered := reporter.FilterValid(r)
	if filtered.Source != r.Source {
		t.Errorf("expected source %q, got %q", r.Source, filtered.Source)
	}
}
