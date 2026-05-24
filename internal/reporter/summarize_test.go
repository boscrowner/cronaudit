package reporter_test

import (
	"strings"
	"testing"

	"github.com/example/cronaudit/internal/reporter"
)

func buildSummarizeReport() reporter.Report {
	return reporter.Build([]string{
		"# daily backup",
		"0 2 * * * /usr/bin/backup",
		"*/5 * * * * /usr/bin/check",
		"bad entry here",
		"",
	})
}

func TestSummarizeReport_EntryCount(t *testing.T) {
	r := buildSummarizeReport()
	sr := reporter.SummarizeReport(r)
	// 2 valid + 1 invalid = 3 entries (comments and blanks skipped)
	if len(sr.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(sr.Entries))
	}
}

func TestSummarizeReport_ValidEntryHasSummary(t *testing.T) {
	r := buildSummarizeReport()
	sr := reporter.SummarizeReport(r)
	for _, e := range sr.Entries {
		if e.Valid && e.Summary == "" {
			t.Errorf("valid entry at line %d has empty summary", e.Line)
		}
	}
}

func TestSummarizeReport_InvalidEntryHasError(t *testing.T) {
	r := buildSummarizeReport()
	sr := reporter.SummarizeReport(r)
	for _, e := range sr.Entries {
		if !e.Valid && e.Error == "" {
			t.Errorf("invalid entry at line %d has empty error", e.Line)
		}
	}
}

func TestSummarizeReport_SourcePreserved(t *testing.T) {
	r := buildSummarizeReport()
	sr := reporter.SummarizeReport(r)
	if sr.Source != r.Source {
		t.Errorf("expected source %q, got %q", r.Source, sr.Source)
	}
}

func TestFormatSummaryReport_ContainsSource(t *testing.T) {
	r := buildSummarizeReport()
	sr := reporter.SummarizeReport(r)
	out := reporter.FormatSummaryReport(sr)
	if !strings.Contains(out, "Source:") {
		t.Error("formatted output missing 'Source:' header")
	}
}

func TestFormatSummaryReport_ContainsArrowForValid(t *testing.T) {
	r := buildSummarizeReport()
	sr := reporter.SummarizeReport(r)
	out := reporter.FormatSummaryReport(sr)
	if !strings.Contains(out, "→") {
		t.Error("formatted output missing '→' for valid entries")
	}
}

func TestFormatSummaryReport_ContainsCrossForInvalid(t *testing.T) {
	r := buildSummarizeReport()
	sr := reporter.SummarizeReport(r)
	out := reporter.FormatSummaryReport(sr)
	if !strings.Contains(out, "✗") {
		t.Error("formatted output missing '✗' for invalid entries")
	}
}
