package reporter

// Stats holds aggregate counts derived from a Report.
type Stats struct {
	Total   int
	Valid   int
	Invalid int
}

// ComputeStats returns a Stats summary for the given Report.
func ComputeStats(r Report) Stats {
	s := Stats{
		Total: r.Stats.Total,
	}
	for _, e := range r.Entries {
		if e.Valid {
			s.Valid++
		} else {
			s.Invalid++
		}
	}
	return s
}

// PassRate returns the percentage of valid entries as a value between 0 and 1.
// Returns 0 if there are no entries.
func (s Stats) PassRate() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Valid) / float64(s.Total)
}

// HasErrors reports whether any invalid entries exist.
func (s Stats) HasErrors() bool {
	return s.Invalid > 0
}
