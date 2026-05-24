package reporter

// CompareResult holds the result of comparing two report entries side by side.
type CompareResult struct {
	Schedule string
	CommandA string
	CommandB string
	SummaryA string
	SummaryB string
	Match     bool
}

// CompareReport compares entries from two reports by schedule, returning
// a slice of CompareResult for schedules that appear in both reports.
func CompareReport(a, b Report) []CompareResult {
	indexA := indexEntriesBySchedule(a)
	indexB := indexEntriesBySchedule(b)

	var results []CompareResult
	for schedule, entryA := range indexA {
		if entryB, ok := indexB[schedule]; ok {
			results = append(results, CompareResult{
				Schedule: schedule,
				CommandA: entryA.Command,
				CommandB: entryB.Command,
				SummaryA: entryA.Summary,
				SummaryB: entryB.Summary,
				Match:     entryA.Command == entryB.Command,
			})
		}
	}
	return results
}

// UnmatchedEntries returns entries from report a whose schedules do not appear in report b.
func UnmatchedEntries(a, b Report) []Entry {
	indexB := indexEntriesBySchedule(b)
	var unmatched []Entry
	for _, entry := range a.Entries {
		if !entry.Valid {
			continue
		}
		if _, found := indexB[entry.Schedule]; !found {
			unmatched = append(unmatched, entry)
		}
	}
	return unmatched
}

func indexEntriesBySchedule(r Report) map[string]Entry {
	m := make(map[string]Entry, len(r.Entries))
	for _, e := range r.Entries {
		if e.Valid {
			m[e.Schedule] = e
		}
	}
	return m
}
