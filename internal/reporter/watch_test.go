package reporter

import (
	"testing"
	"time"
)

func buildWatchSnapshots() (Snapshot, Snapshot) {
	base := Report{
		Source: "crontab",
		Entries: []Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
			{Line: 2, Schedule: "30 6 * * *", Command: "/bin/bar", Valid: true},
		},
	}
	before := TakeSnapshot(base, "before")

	modified := Report{
		Source: "crontab",
		Entries: []Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
			{Line: 3, Schedule: "15 3 * * 1", Command: "/bin/baz", Valid: true},
		},
	}
	after := TakeSnapshot(modified, "after")
	return before, after
}

func TestWatchSnapshots_ChangedTrue(t *testing.T) {
	before, after := buildWatchSnapshots()
	result := WatchSnapshots(before, after)
	if !result.Changed {
		t.Error("expected Changed to be true")
	}
}

func TestWatchSnapshots_DetectsAdded(t *testing.T) {
	before, after := buildWatchSnapshots()
	result := WatchSnapshots(before, after)
	if len(result.Diff.Added) != 1 {
		t.Errorf("expected 1 added entry, got %d", len(result.Diff.Added))
	}
}

func TestWatchSnapshots_DetectsRemoved(t *testing.T) {
	before, after := buildWatchSnapshots()
	result := WatchSnapshots(before, after)
	if len(result.Diff.Removed) != 1 {
		t.Errorf("expected 1 removed entry, got %d", len(result.Diff.Removed))
	}
}

func TestWatchSnapshots_NoChange(t *testing.T) {
	base := Report{
		Source: "crontab",
		Entries: []Entry{
			{Line: 1, Schedule: "0 * * * *", Command: "/bin/foo", Valid: true},
		},
	}
	before := TakeSnapshot(base, "v1")
	after := TakeSnapshot(base, "v2")
	result := WatchSnapshots(before, after)
	if result.Changed {
		t.Error("expected Changed to be false")
	}
}

func TestWatchSnapshots_CheckedAtIsRecent(t *testing.T) {
	before, after := buildWatchSnapshots()
	now := time.Now()
	result := WatchSnapshots(before, after)
	if result.CheckedAt.Before(now.Add(-time.Second)) {
		t.Error("CheckedAt should be recent")
	}
}

func TestWatchSummary_NoChange(t *testing.T) {
	w := WatchResult{Changed: false}
	got := WatchSummary(w)
	if got != "No changes detected between snapshots." {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestWatchSummary_WithChanges(t *testing.T) {
	before, after := buildWatchSnapshots()
	result := WatchSnapshots(before, after)
	summary := WatchSummary(result)
	if summary == "" {
		t.Error("expected non-empty summary for changed result")
	}
}
