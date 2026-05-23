package reporter

import (
	"time"
)

// WatchResult holds the outcome of comparing two snapshots taken over time.
type WatchResult struct {
	Before    Snapshot
	After     Snapshot
	Diff      DiffResult
	Changed   bool
	CheckedAt time.Time
}

// WatchSnapshots compares two snapshots and returns a WatchResult describing
// what changed between them. Changed is true if any entries were added or removed.
func WatchSnapshots(before, after Snapshot) WatchResult {
	diff := DiffSnapshots(before, after)
	changed := len(diff.Added) > 0 || len(diff.Removed) > 0
	return WatchResult{
		Before:    before,
		After:     after,
		Diff:      diff,
		Changed:   changed,
		CheckedAt: time.Now(),
	}
}

// WatchSummary returns a human-readable string summarising the watch result.
func WatchSummary(w WatchResult) string {
	if !w.Changed {
		return "No changes detected between snapshots."
	}
	var msg string
	if len(w.Diff.Added) > 0 {
		msg += formatCount(len(w.Diff.Added), "entry", "entries") + " added. "
	}
	if len(w.Diff.Removed) > 0 {
		msg += formatCount(len(w.Diff.Removed), "entry", "entries") + " removed. "
	}
	return msg
}

func formatCount(n int, singular, plural string) string {
	if n == 1 {
		return "1 " + singular
	}
	return itoa(n) + " " + plural
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
