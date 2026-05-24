package reporter

import (
	"strings"
	"testing"
)

func buildDependencyReport() Report {
	return Report{
		Source: "testcron",
		Entries: []Entry{
			{Line: 1, Raw: "0 * * * * /usr/bin/backup", Schedule: "0 * * * *", Command: "/usr/bin/backup", Valid: true},
			{Line: 2, Raw: "30 * * * * /usr/bin/backup --verify", Schedule: "30 * * * *", Command: "/usr/bin/backup --verify", Valid: true},
			{Line: 3, Raw: "0 2 * * * /usr/bin/cleanup", Schedule: "0 2 * * *", Command: "/usr/bin/cleanup", Valid: true},
			{Line: 4, Raw: "0 * * * * /usr/bin/sync", Schedule: "0 * * * *", Command: "/usr/bin/sync", Valid: true},
			{Line: 5, Raw: "bad entry", Valid: false, Error: "invalid"},
		},
	}
}

func TestDetectDependencies_SameBaseCommand(t *testing.T) {
	r := buildDependencyReport()
	dr := DetectDependencies(r)

	for _, d := range dr.Dependencies {
		if strings.Contains(d.Reason, "backup") {
			return
		}
	}
	t.Error("expected dependency for shared base command 'backup'")
}

func TestDetectDependencies_SameSchedule(t *testing.T) {
	r := buildDependencyReport()
	dr := DetectDependencies(r)

	for _, d := range dr.Dependencies {
		if strings.Contains(d.Reason, "identical schedule") {
			return
		}
	}
	t.Error("expected dependency for identical schedule '0 * * * *'")
}

func TestDetectDependencies_SkipsInvalidEntries(t *testing.T) {
	r := buildDependencyReport()
	dr := DetectDependencies(r)

	for _, d := range dr.Dependencies {
		if !d.A.Valid || !d.B.Valid {
			t.Error("dependency includes an invalid entry")
		}
	}
}

func TestDetectDependencies_SourcePreserved(t *testing.T) {
	r := buildDependencyReport()
	dr := DetectDependencies(r)

	if dr.Source != r.Source {
		t.Errorf("expected source %q, got %q", r.Source, dr.Source)
	}
}

func TestFormatDependencyReport_NoDeps(t *testing.T) {
	dr := DependencyReport{Source: "empty", Dependencies: nil}
	out := FormatDependencyReport(dr)

	if !strings.Contains(out, "No scheduling dependencies") {
		t.Errorf("expected no-dependency message, got: %s", out)
	}
}

func TestFormatDependencyReport_ContainsLineNumbers(t *testing.T) {
	r := buildDependencyReport()
	dr := DetectDependencies(r)
	out := FormatDependencyReport(dr)

	if !strings.Contains(out, "line") {
		t.Errorf("expected line references in output, got: %s", out)
	}
}
