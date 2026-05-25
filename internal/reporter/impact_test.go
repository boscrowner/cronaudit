package reporter

import (
	"strings"
	"testing"
)

func buildImpactReport(entries []Entry) Report {
	return Report{
		Source:  "test",
		Entries: entries,
	}
}

func TestAssessImpact_SkipsInvalidEntries(t *testing.T) {
	r := buildImpactReport([]Entry{
		{Line: 1, Valid: false, Command: "rm -rf /"},
	})
	ir := AssessImpact(r)
	if len(ir.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(ir.Entries))
	}
}

func TestAssessImpact_LowRiskCleanEntry(t *testing.T) {
	r := buildImpactReport([]Entry{
		{Line: 1, Valid: true, Schedule: "0 * * * *", Command: "/usr/local/bin/backup.sh"},
	})
	ir := AssessImpact(r)
	if len(ir.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ir.Entries))
	}
	// /usr/local/bin matches system path heuristic → medium
	if ir.Entries[0].Level != ImpactMedium {
		t.Errorf("expected medium, got %s", ir.Entries[0].Level)
	}
}

func TestAssessImpact_CriticalMultipleFactors(t *testing.T) {
	r := buildImpactReport([]Entry{
		{Line: 2, Valid: true, Schedule: "* * * * *",
			Command: "sudo rm -rf /var/log && curl http://example.com"},
	})
	ir := AssessImpact(r)
	if len(ir.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if ir.Entries[0].Level != ImpactCritical {
		t.Errorf("expected critical, got %s", ir.Entries[0].Level)
	}
	if len(ir.Entries[0].Reasons) < 3 {
		t.Errorf("expected >=3 reasons, got %d", len(ir.Entries[0].Reasons))
	}
}

func TestAssessImpact_SortedDescending(t *testing.T) {
	r := buildImpactReport([]Entry{
		{Line: 1, Valid: true, Schedule: "0 1 * * *", Command: "echo hello"},
		{Line: 2, Valid: true, Schedule: "0 2 * * *",
			Command: "sudo rm -rf /etc/config && curl http://x.com && wget http://y.com"},
	})
	ir := AssessImpact(r)
	if len(ir.Entries) < 2 {
		t.Fatalf("expected 2 entries")
	}
	if impactWeight(ir.Entries[0].Level) < impactWeight(ir.Entries[1].Level) {
		t.Errorf("entries not sorted descending by impact")
	}
}

func TestAssessImpact_SourcePreserved(t *testing.T) {
	r := buildImpactReport([]Entry{
		{Line: 1, Valid: true, Schedule: "@daily", Command: "echo ok"},
	})
	ir := AssessImpact(r)
	if ir.Source != "test" {
		t.Errorf("expected source 'test', got %q", ir.Source)
	}
}

func TestFormatImpactReport_ContainsLevel(t *testing.T) {
	r := buildImpactReport([]Entry{
		{Line: 1, Valid: true, Schedule: "0 * * * *", Command: "sudo wget http://example.com"},
	})
	ir := AssessImpact(r)
	out := FormatImpactReport(ir)
	if !strings.Contains(out, "high") && !strings.Contains(out, "critical") {
		t.Errorf("expected high or critical in output, got: %s", out)
	}
}

func TestFormatImpactReport_EmptyReport(t *testing.T) {
	ir := ImpactReport{Source: "empty"}
	out := FormatImpactReport(ir)
	if !strings.Contains(out, "no valid entries") {
		t.Errorf("expected 'no valid entries' message, got: %s", out)
	}
}
