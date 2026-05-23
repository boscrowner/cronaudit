package reporter

import (
	"sort"
)

// SortBy defines the field to sort report entries by.
type SortBy int

const (
	// SortByLine sorts entries by their original line number (default order).
	SortByLine SortBy = iota
	// SortBySchedule sorts entries lexicographically by their raw schedule expression.
	SortBySchedule
	// SortByValid groups invalid entries first, then valid ones.
	SortByValid
)

// SortEntries returns a new slice of entries sorted by the given field.
// The original report is not mutated.
func SortEntries(r Report, by SortBy) []Entry {
	entries := make([]Entry, len(r.Entries))
	copy(entries, r.Entries)

	switch by {
	case SortBySchedule:
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].Raw < entries[j].Raw
		})
	case SortByValid:
		sort.SliceStable(entries, func(i, j int) bool {
			// invalid entries (Error != "") sort before valid ones
			iInvalid := entries[i].Error != ""
			jInvalid := entries[j].Error != ""
			if iInvalid == jInvalid {
				return false
			}
			return iInvalid
		})
	default: // SortByLine — already in insertion order
	}

	return entries
}
