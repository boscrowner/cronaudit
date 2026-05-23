package reporter

// DedupeEntries returns a new Report containing only the first occurrence
// of each unique schedule expression. Duplicate schedules are tracked by
// their raw cron string (the Schedule field of each Entry).
//
// The original Report is not modified. Source and Stats are preserved from
// the input; Stats reflect the deduplicated entry set.
func DedupeEntries(r Report) Report {
	seen := make(map[string]bool)
	deduped := make([]Entry, 0, len(r.Entries))

	for _, e := range r.Entries {
		if seen[e.Schedule] {
			continue
		}
		seen[e.Schedule] = true
		deduped = append(deduped, e)
	}

	return Report{
		Source:  r.Source,
		Entries: deduped,
		Stats:   ComputeStats(deduped),
	}
}

// DuplicateSchedules returns a slice of schedule strings that appear more
// than once across the entries in the given Report.
func DuplicateSchedules(r Report) []string {
	counts := make(map[string]int)
	for _, e := range r.Entries {
		counts[e.Schedule]++
	}

	duplicates := make([]string, 0)
	for sched, count := range counts {
		if count > 1 {
			duplicates = append(duplicates, sched)
		}
	}
	return duplicates
}
