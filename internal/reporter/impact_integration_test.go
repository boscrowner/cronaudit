package reporter_test

import (
	"strings"
	"testing"

	"github.com/cronaudit/cronaudit/internal/reporter"
)

// TestAssessImpact_IntegrationWithBuild verifies that AssessImpact works
// correctly with a Report produced by reporter.Build.
func TestAssessImpact_IntegrationWithBuild(t *testing.T) {
	lines := []string{
		"# daily cleanup",
		"0 2 * * * sudo rm -rf /var/tmp && curl http://notify.internal",
		"*/10 * * * * echo ping",
		"bad line no schedule",
	}
	r := reporter.Build(lines, "integration-test")
	ir := reporter.AssessImpact(r)

	// Only valid entries should appear.
	for _, e := range ir.Entries {
		if !entryIsValid(r, e.Line) {
			t.Errorf("line %d is invalid but appeared in impact report", e.Line)
		}
	}
}

func entryIsValid(r reporter.Report, line int) bool {
	for _, e := range r.Entries {
		if e.Line == line {
			return e.Valid
		}
	}
	return false
}

func TestAssessImpact_FormatContainsAllEntries(t *testing.T) {
	lines := []string{
		"0 1 * * * /usr/bin/backup",
		"0 2 * * * sudo chown root /etc/passwd",
	}
	r := reporter.Build(lines, "format-test")
	ir := reporter.AssessImpact(r)
	out := reporter.FormatImpactReport(ir)

	if !strings.Contains(out, "format-test") {
		t.Errorf("output missing source name")
	}
	if len(ir.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ir.Entries))
	}
	for _, e := range ir.Entries {
		if !strings.Contains(out, e.Command) {
			t.Errorf("output missing command: %s", e.Command)
		}
	}
}
