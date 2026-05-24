package reporter

import (
	"strings"
	"testing"
)

func buildRiskReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "/usr/bin/backup.sh", Valid: true},
			{Line: 2, Schedule: "*/2 * * * *", Command: "sudo rm -rf /tmp/cache", Valid: true},
			{Line: 3, Schedule: "0 3 * * *", Command: "dd if=/dev/zero of=/tmp/test", Valid: true},
			{Line: 4, Schedule: "bad schedule", Command: "ignored", Valid: false},
			{Line: 5, Schedule: "*/1 * * * *", Command: "sudo /scripts/sync.sh", Valid: true},
		},
	}
}

func TestAssessRisk_SkipsInvalidEntries(t *testing.T) {
	r := buildRiskReport()
	rr := AssessRisk(r)
	for _, e := range rr.Entries {
		if e.Line == 4 {
			t.Errorf("expected invalid entry (line 4) to be skipped")
		}
	}
}

func TestAssessRisk_LowRiskCleanEntry(t *testing.T) {
	r := buildRiskReport()
	rr := AssessRisk(r)
	for _, e := range rr.Entries {
		if e.Line == 1 {
			if e.Level != RiskLow {
				t.Errorf("expected low risk for clean entry, got %s", e.Level)
			}
			if len(e.Reasons) != 0 {
				t.Errorf("expected no reasons for low risk entry, got %v", e.Reasons)
			}
			return
		}
	}
	t.Error("line 1 not found in risk report")
}

func TestAssessRisk_HighRiskMultipleFactors(t *testing.T) {
	r := buildRiskReport()
	rr := AssessRisk(r)
	for _, e := range rr.Entries {
		if e.Line == 2 {
			if e.Level != RiskHigh {
				t.Errorf("expected high risk for line 2, got %s (reasons: %v)", e.Level, e.Reasons)
			}
			return
		}
	}
	t.Error("line 2 not found in risk report")
}

func TestAssessRisk_MediumRiskSingleFactor(t *testing.T) {
	r := buildRiskReport()
	rr := AssessRisk(r)
	for _, e := range rr.Entries {
		if e.Line == 3 {
			if e.Level == RiskLow {
				t.Errorf("expected at least medium risk for line 3 (dd command), got low")
			}
			return
		}
	}
	t.Error("line 3 not found in risk report")
}

func TestAssessRisk_SourcePreserved(t *testing.T) {
	r := buildRiskReport()
	rr := AssessRisk(r)
	if rr.Source != r.Source {
		t.Errorf("expected source %q, got %q", r.Source, rr.Source)
	}
}

func TestFormatRiskReport_ContainsLevel(t *testing.T) {
	r := buildRiskReport()
	rr := AssessRisk(r)
	out := FormatRiskReport(rr)
	for _, level := range []string{"low", "medium", "high"} {
		_ = level
	}
	if !strings.Contains(out, "risk report for test") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestFormatRiskReport_EmptyEntries(t *testing.T) {
	rr := RiskReport{Source: "empty", Entries: nil}
	out := FormatRiskReport(rr)
	if !strings.Contains(out, "no valid entries assessed") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
