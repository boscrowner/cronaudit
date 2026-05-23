package reporter

// DiffResult holds the comparison between two reports.
type DiffResult struct {
	Added   []Entry
	Removed []Entry
	Kept    []Entry
}

// DiffReports compares two reports and returns entries that were added,
// removed, or kept between the "before" and "after" report.
// Comparison is based on the raw schedule+command string.
func DiffReports(before, after Report) DiffResult {
	beforeIndex := indexEntries(before.Entries)
	afterIndex := indexEntries(after.Entries)

	var result DiffResult

	for key, entry := range afterIndex {
		if _, exists := beforeIndex[key]; exists {
			result.Kept = append(result.Kept, entry)
		} else {
			result.Added = append(result.Added, entry)
		}
	}

	for key, entry := range beforeIndex {
		if _, exists := afterIndex[key]; !exists {
			result.Removed = append(result.Removed, entry)
		}
	}

	return result
}

// indexEntries builds a map keyed by "schedule|command" for fast lookup.
func indexEntries(entries []Entry) map[string]Entry {
	idx := make(map[string]Entry, len(entries))
	for _, e := range entries {
		key := e.Schedule + "|" + e.Command
		idx[key] = e
	}
	return idx
}
