package reporter

// MergeReports combines two Reports into a single Report.
// Entries from both reports are concatenated in order (a first, then b).
// The Source field is taken from a; Stats are recomputed from the merged entries.
// Duplicate entries are NOT removed — call DedupeEntries afterwards if needed.
func MergeReports(a, b Report) Report {
	mergedEntries := make([]Entry, 0, len(a.Entries)+len(b.Entries))
	mergedEntries = append(mergedEntries, a.Entries...)
	mergedEntries = append(mergedEntries, b.Entries...)

	return Report{
		Source:  a.Source,
		Entries: mergedEntries,
		Stats:   ComputeStats(mergedEntries),
	}
}

// MergeMany combines multiple Reports into a single Report.
// Entries are appended in the order the reports are provided.
// The Source field is taken from the first non-empty report.
// Returns an empty Report if no reports are provided.
func MergeMany(reports []Report) Report {
	if len(reports) == 0 {
		return Report{}
	}

	var source string
	var total []Entry

	for _, r := range reports {
		if source == "" && r.Source != "" {
			source = r.Source
		}
		total = append(total, r.Entries...)
	}

	return Report{
		Source:  source,
		Entries: total,
		Stats:   ComputeStats(total),
	}
}
