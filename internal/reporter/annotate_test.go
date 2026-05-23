package reporter

import (
	"strings"
	"testing"
)

func buildAnnotateReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Raw: "0 2 * * * /backup.sh", Valid: true, Summary: "At 02:00 daily"},
			{Line: 2, Raw: "*/5 * * * * /poll.sh", Valid: true, Summary: "Every 5 minutes"},
			{Line: 3, Raw: "bad entry", Valid: false, Error: "invalid field count", Summary: ""},
		},
	}
}

func TestAnnotateEntries_AppliesLabel(t *testing.T) {
	r := buildAnnotateReport()
	tags := map[string]string{
		"0 2 * * * /backup.sh": "nightly-backup",
	}
	out := AnnotateEntries(r, tags)
	if !strings.Contains(out.Entries[0].Summary, "[nightly-backup]") {
		t.Errorf("expected annotation in summary, got: %s", out.Entries[0].Summary)
	}
}

func TestAnnotateEntries_PreservesUntagged(t *testing.T) {
	r := buildAnnotateReport()
	tags := map[string]string{}
	out := AnnotateEntries(r, tags)
	if out.Entries[1].Summary != r.Entries[1].Summary {
		t.Errorf("untagged entry should be unchanged")
	}
}

func TestAnnotateEntries_DoesNotMutateOriginal(t *testing.T) {
	r := buildAnnotateReport()
	orig := r.Entries[0].Summary
	tags := map[string]string{"0 2 * * * /backup.sh": "nightly"}
	AnnotateEntries(r, tags)
	if r.Entries[0].Summary != orig {
		t.Error("original report was mutated")
	}
}

func TestAnnotateErrors_PrefixesInvalidEntries(t *testing.T) {
	r := buildAnnotateReport()
	out := AnnotateErrors(r)
	if !strings.HasPrefix(out.Entries[2].Summary, "[ERROR:") {
		t.Errorf("expected ERROR prefix, got: %s", out.Entries[2].Summary)
	}
}

func TestAnnotateErrors_LeavesValidEntriesAlone(t *testing.T) {
	r := buildAnnotateReport()
	out := AnnotateErrors(r)
	if out.Entries[0].Summary != r.Entries[0].Summary {
		t.Errorf("valid entry summary should not change, got: %s", out.Entries[0].Summary)
	}
}

func TestAnnotateErrors_SourcePreserved(t *testing.T) {
	r := buildAnnotateReport()
	out := AnnotateErrors(r)
	if out.Source != r.Source {
		t.Errorf("source mismatch: got %s", out.Source)
	}
}
