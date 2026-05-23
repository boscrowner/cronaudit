package reporter

import (
	"sort"

	"github.com/cronaudit/internal/parser"
)

// RankedEntry pairs a report entry with a computed rank score.
type RankedEntry struct {
	Entry parser.Entry
	Score float64
	Reason string
}

// RankReport ranks entries in a Report by a heuristic score.
// Valid entries are ranked higher than invalid ones. Among valid
// entries, those with more specific schedules (fewer wildcards)
// rank higher. The returned slice is sorted descending by score.
func RankReport(r Report) []RankedEntry {
	ranked := make([]RankedEntry, 0, len(r.Entries))
	for _, e := range r.Entries {
		ranked = append(ranked, scoreEntry(e))
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})
	return ranked
}

// TopN returns at most n highest-ranked entries from a Report.
func TopN(r Report, n int) []RankedEntry {
	all := RankReport(r)
	if n > len(all) {
		n = len(all)
	}
	return all[:n]
}

func scoreEntry(e parser.Entry) RankedEntry {
	if !e.Valid {
		return RankedEntry{Entry: e, Score: 0, Reason: "invalid entry"}
	}

	wildcards := countWildcards(e.Schedule)
	score := 100.0 - float64(wildcards)*15.0
	if score < 5 {
		score = 5
	}

	reason := "valid"
	switch {
	case wildcards == 0:
		reason = "fully specified schedule"
	case wildcards <= 2:
		reason = "mostly specific schedule"
	default:
		reason = "wildcard-heavy schedule"
	}

	return RankedEntry{Entry: e, Score: score, Reason: reason}
}

func countWildcards(schedule string) int {
	fields := splitFields(schedule)
	count := 0
	for _, f := range fields {
		if f == "*" {
			count++
		}
	}
	return count
}
