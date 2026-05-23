package reporter_test

import (
	"testing"

	"github.com/yourorg/cronaudit/internal/reporter"
)

var sampleLines = []string{
	"# daily backup",
	"",
	"0 2 * * * /usr/bin/backup.sh",
	"*/15 * * * * /usr/bin/healthcheck",
	"bad entry here",
}

func TestBuild_Stats(t *testing.T) {
	r := reporter.Build("test", sampleLines)

	if r.Stats.Total != len(sampleLines) {
		t.Errorf("expected Total=%d, got %d", len(sampleLines), r.Stats.Total)
	}
	if r.Stats.Valid != 4 {
		t.Errorf("expected Valid=4, got %d", r.Stats.Valid)
	}
	if r.Stats.Invalid != 1 {
		t.Errorf("expected Invalid=1, got %d", r.Stats.Invalid)
	}
}

func TestBuild_Source(t *testing.T) {
	r := reporter.Build("mycrontab", sampleLines)
	if r.Source != "mycrontab" {
		t.Errorf("expected source 'mycrontab', got %q", r.Source)
	}
}

func TestBuild_ValidEntryHasSummary(t *testing.T) {
	r := reporter.Build("test", []string{"0 2 * * * /usr/bin/backup.sh"})
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if !e.Valid {
		t.Errorf("expected entry to be valid")
	}
	if e.Summary == "" {
		t.Errorf("expected non-empty summary for valid entry")
	}
}

func TestBuild_InvalidEntryHasError(t *testing.T) {
	r := reporter.Build("test", []string{"bad entry here"})
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Valid {
		t.Errorf("expected entry to be invalid")
	}
	if e.Error == "" {
		t.Errorf("expected non-empty error for invalid entry")
	}
}

func TestBuild_EmptyLines(t *testing.T) {
	r := reporter.Build("empty", []string{})
	if r.Stats.Total != 0 {
		t.Errorf("expected Total=0 for empty input")
	}
	if r.Entries != nil && len(r.Entries) != 0 {
		t.Errorf("expected no entries for empty input")
	}
}
