package reporter

import (
	"sort"
	"time"
)

// TimelineEntry represents a single scheduled entry projected onto a timeline window.
type TimelineEntry struct {
	Line     int
	Schedule string
	Command  string
	NextRuns []time.Time
}

// Timeline holds projected run times for all valid cron entries within a window.
type Timeline struct {
	Source  string
	From    time.Time
	To      time.Time
	Entries []TimelineEntry
}

// BuildTimeline projects valid cron entries onto a time window [from, to],
// computing up to maxRuns next execution times per entry.
func BuildTimeline(r Report, from, to time.Time, maxRuns int) Timeline {
	tl := Timeline{
		Source: r.Source,
		From:   from,
		To:     to,
	}

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		runs := projectRuns(e.Schedule, from, to, maxRuns)
		if len(runs) == 0 {
			continue
		}
		tl.Entries = append(tl.Entries, TimelineEntry{
			Line:     e.Line,
			Schedule: e.Schedule,
			Command:  e.Command,
			NextRuns: runs,
		})
	}

	sort.Slice(tl.Entries, func(i, j int) bool {
		if len(tl.Entries[i].NextRuns) == 0 {
			return false
		}
		if len(tl.Entries[j].NextRuns) == 0 {
			return true
		}
		return tl.Entries[i].NextRuns[0].Before(tl.Entries[j].NextRuns[0])
	})

	return tl
}

// projectRuns returns up to maxRuns times within [from, to] that match the cron schedule.
// It uses a simplified minute-tick approach.
func projectRuns(schedule string, from, to time.Time, maxRuns int) []time.Time {
	fields := splitFields(schedule)
	if len(fields) != 5 {
		return nil
	}

	var runs []time.Time
	current := from.Truncate(time.Minute).Add(time.Minute)

	for !current.After(to) && len(runs) < maxRuns {
		if matchesCron(fields, current) {
			runs = append(runs, current)
		}
		current = current.Add(time.Minute)
	}
	return runs
}
