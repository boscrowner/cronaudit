package reporter

// FilterInvalid returns a new Report containing only the invalid entries
// from the original, with stats recalculated.
func FilterInvalid(r Report) Report {
	return filterBy(r, func(e Entry) bool { return !e.Valid })
}

// FilterValid returns a new Report containing only the valid, non-blank
// entries from the original, with stats recalculated.
func FilterValid(r Report) Report {
	return filterBy(r, func(e Entry) bool { return e.Valid && e.Fields != nil })
}

func filterBy(r Report, keep func(Entry) bool) Report {
	out := Report{Source: r.Source}
	for _, e := range r.Entries {
		if keep(e) {
			out.Entries = append(out.Entries, e)
			out.Stats.Total++
			if e.Valid {
				out.Stats.Valid++
			} else {
				out.Stats.Invalid++
			}
		}
	}
	return out
}
