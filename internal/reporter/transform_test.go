package reporter

import (
	"strings"
	"testing"

	"github.com/user/cronaudit/internal/parser"
)

func buildTransformReport() Report {
	return Report{
		Source: "test",
		Entries: []parser.Entry{
			{Line: 1, Raw: "*/5 * * * * /usr/bin/backup", Valid: true, Command: "  /usr/bin/backup  ", Schedule: "*/5 * * * *"},
			{Line: 2, Raw: "0 2 * * * /usr/bin/clean", Valid: true, Command: "/usr/bin/clean", Schedule: "0 2 * * *"},
			{Line: 3, Raw: "bad entry", Valid: false, Command: "", Error: "invalid format"},
		},
	}
}

func TestTransformEntries_DoesNotMutateOriginal(t *testing.T) {
	r := buildTransformReport()
	origCmd := r.Entries[0].Command
	TransformEntries(r, RedactCommands("[redacted]"))
	if r.Entries[0].Command != origCmd {
		t.Errorf("original entry mutated: got %q, want %q", r.Entries[0].Command, origCmd)
	}
}

func TestNormalizeCommands_TrimsWhitespace(t *testing.T) {
	r := buildTransformReport()
	result := TransformEntries(r, NormalizeCommands())
	for _, e := range result.Entries {
		if strings.TrimSpace(e.Command) != e.Command {
			t.Errorf("command not trimmed: %q", e.Command)
		}
	}
}

func TestRedactCommands_ReplacesValidCommands(t *testing.T) {
	r := buildTransformReport()
	result := TransformEntries(r, RedactCommands("[redacted]"))
	for _, e := range result.Entries {
		if e.Valid && e.Command != "[redacted]" {
			t.Errorf("valid command not redacted: got %q", e.Command)
		}
	}
}

func TestRedactCommands_PreservesInvalidEntries(t *testing.T) {
	r := buildTransformReport()
	result := TransformEntries(r, RedactCommands("[redacted]"))
	for i, e := range result.Entries {
		if !e.Valid && e.Command != r.Entries[i].Command {
			t.Errorf("invalid entry command changed: got %q, want %q", e.Command, r.Entries[i].Command)
		}
	}
}

func TestChainTransforms_AppliesInOrder(t *testing.T) {
	r := buildTransformReport()
	chained := ChainTransforms(NormalizeCommands(), RedactCommands("[redacted]"))
	result := TransformEntries(r, chained)
	for _, e := range result.Entries {
		if e.Valid && e.Command != "[redacted]" {
			t.Errorf("chained transform did not redact: got %q", e.Command)
		}
	}
}

func TestTransformEntries_StatsRecomputed(t *testing.T) {
	r := buildTransformReport()
	result := TransformEntries(r, NormalizeCommands())
	expected := ComputeStats(result.Entries)
	if result.Stats.Total != expected.Total {
		t.Errorf("stats not recomputed: got total %d, want %d", result.Stats.Total, expected.Total)
	}
}

func TestTransformEntries_SourcePreserved(t *testing.T) {
	r := buildTransformReport()
	result := TransformEntries(r, NormalizeCommands())
	if result.Source != r.Source {
		t.Errorf("source changed: got %q, want %q", result.Source, r.Source)
	}
}
