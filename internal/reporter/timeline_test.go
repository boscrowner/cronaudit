package reporter

import (
	"testing"
	"time"
)

func buildTimelineReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Schedule: "* * * * *", Command: "every-minute", Valid: true},
			{Line: 2, Schedule: "0 * * * *", Command: "every-hour", Valid: true},
			{Line: 3, Schedule: "invalid", Command: "bad", Valid: false},
			{Line: 4, Schedule: "*/5 * * * *", Command: "every-5min", Valid: true},
		},
	}
}

func fixedWindow() (time.Time, time.Time) {
	// Use a fixed Monday 2024-01-01 00:00 UTC as base
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(30 * time.Minute)
	return from, to
}

func TestBuildTimeline_ExcludesInvalidEntries(t *testing.T) {
	r := buildTimelineReport()
	from, to := fixedWindow()
	tl := BuildTimeline(r, from, to, 10)
	for _, e := range tl.Entries {
		if e.Command == "bad" {
			t.Errorf("expected invalid entry to be excluded from timeline")
		}
	}
}

func TestBuildTimeline_EveryMinuteHasRuns(t *testing.T) {
	r := buildTimelineReport()
	from, to := fixedWindow()
	tl := BuildTimeline(r, from, to, 10)
	for _, e := range tl.Entries {
		if e.Command == "every-minute" {
			if len(e.NextRuns) != 10 {
				t.Errorf("expected 10 runs, got %d", len(e.NextRuns))
			}
			return
		}
	}
	t.Error("every-minute entry not found")
}

func TestBuildTimeline_SortedByFirstRun(t *testing.T) {
	r := buildTimelineReport()
	from, to := fixedWindow()
	tl := BuildTimeline(r, from, to, 5)
	for i := 1; i < len(tl.Entries); i++ {
		prev := tl.Entries[i-1].NextRuns[0]
		curr := tl.Entries[i].NextRuns[0]
		if curr.Before(prev) {
			t.Errorf("timeline not sorted: entry %d (%s) before entry %d (%s)",
				i, curr, i-1, prev)
		}
	}
}

func TestBuildTimeline_SourcePreserved(t *testing.T) {
	r := buildTimelineReport()
	from, to := fixedWindow()
	tl := BuildTimeline(r, from, to, 3)
	if tl.Source != "test" {
		t.Errorf("expected source 'test', got %q", tl.Source)
	}
}

func TestBuildTimeline_StepSchedule(t *testing.T) {
	r := buildTimelineReport()
	from, to := fixedWindow()
	tl := BuildTimeline(r, from, to, 10)
	for _, e := range tl.Entries {
		if e.Command == "every-5min" {
			for _, run := range e.NextRuns {
				if run.Minute()%5 != 0 {
					t.Errorf("expected run at minute divisible by 5, got %d", run.Minute())
				}
			}
			return
		}
	}
	t.Error("every-5min entry not found")
}
