package reporter

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildSnapshotReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Raw: "*/5 * * * * /usr/bin/check", Valid: true, Schedule: "*/5 * * * *", Command: "/usr/bin/check"},
			{Line: 2, Raw: "bad entry", Valid: false, Error: "invalid cron expression"},
		},
	}
}

func TestTakeSnapshot_SetsLabel(t *testing.T) {
	r := buildSnapshotReport()
	s := TakeSnapshot(r, "prod")
	if s.Label != "prod" {
		t.Errorf("expected label 'prod', got %q", s.Label)
	}
}

func TestTakeSnapshot_CapturedAtIsRecent(t *testing.T) {
	before := time.Now().UTC()
	s := TakeSnapshot(buildSnapshotReport(), "")
	after := time.Now().UTC()
	if s.CapturedAt.Before(before) || s.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v not within expected range", s.CapturedAt)
	}
}

func TestWriteReadSnapshot_RoundTrip(t *testing.T) {
	orig := TakeSnapshot(buildSnapshotReport(), "ci")

	var buf bytes.Buffer
	if err := WriteSnapshot(&buf, orig); err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}

	got, err := ReadSnapshot(&buf)
	if err != nil {
		t.Fatalf("ReadSnapshot: %v", err)
	}

	if got.Label != orig.Label {
		t.Errorf("label: want %q, got %q", orig.Label, got.Label)
	}
	if got.Report.Source != orig.Report.Source {
		t.Errorf("source: want %q, got %q", orig.Report.Source, got.Report.Source)
	}
	if len(got.Report.Entries) != len(orig.Report.Entries) {
		t.Errorf("entries: want %d, got %d", len(orig.Report.Entries), len(got.Report.Entries))
	}
}

func TestReadSnapshot_InvalidJSON(t *testing.T) {
	_, err := ReadSnapshot(strings.NewReader("{invalid"))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestDiffSnapshots_DetectsAdded(t *testing.T) {
	base := TakeSnapshot(Report{Source: "a", Entries: []Entry{
		{Line: 1, Raw: "0 * * * * job1", Valid: true, Schedule: "0 * * * *", Command: "job1"},
	}}, "base")

	newer := TakeSnapshot(Report{Source: "b", Entries: []Entry{
		{Line: 1, Raw: "0 * * * * job1", Valid: true, Schedule: "0 * * * *", Command: "job1"},
		{Line: 2, Raw: "5 * * * * job2", Valid: true, Schedule: "5 * * * *", Command: "job2"},
	}}, "newer")

	diff := DiffSnapshots(base, newer)
	if len(diff.Added) != 1 {
		t.Errorf("expected 1 added entry, got %d", len(diff.Added))
	}
}
