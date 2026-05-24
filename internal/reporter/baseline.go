package reporter

import (
	"sort"
	"strings"
)

// BaselineReport holds a reference set of expected cron entries used to
// validate a live report against an approved baseline.
type BaselineReport struct {
	Expected []string // canonical schedule+command strings
	Missing  []string // in baseline but not in report
	Unknown  []string // in report but not in baseline
	Matched  []string // present in both
}

// BuildBaseline compares a Report against a slice of approved schedule strings
// (formatted as "<schedule> <command>") and returns a BaselineReport.
func BuildBaseline(r Report, approved []string) BaselineReport {
	approvedSet := make(map[string]struct{}, len(approved))
	for _, a := range approved {
		approvedSet[strings.TrimSpace(a)] = struct{}{}
	}

	activeSet := make(map[string]struct{}, len(r.Entries))
	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		key := e.Schedule + " " + e.Command
		activeSet[key] = struct{}{}
	}

	var matched, unknown, missing []string

	for key := range activeSet {
		if _, ok := approvedSet[key]; ok {
			matched = append(matched, key)
		} else {
			unknown = append(unknown, key)
		}
	}

	for key := range approvedSet {
		if _, ok := activeSet[key]; !ok {
			missing = append(missing, key)
		}
	}

	sort.Strings(matched)
	sort.Strings(unknown)
	sort.Strings(missing)

	return BaselineReport{
		Expected: approved,
		Missing:  missing,
		Unknown:  unknown,
		Matched:  matched,
	}
}

// IsClean returns true when the live report matches the baseline exactly.
func (b BaselineReport) IsClean() bool {
	return len(b.Missing) == 0 && len(b.Unknown) == 0
}
