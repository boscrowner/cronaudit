package reporter

import "fmt"

// OverlapResult describes two entries whose schedules may fire at the same time.
type OverlapResult struct {
	A       Entry
	B       Entry
	Example string // human-readable description of when both fire
}

// OverlapReport holds all detected schedule overlaps within a report.
type OverlapReport struct {
	Source   string
	Overlaps []OverlapResult
}

// DetectOverlaps scans all valid entries in r and returns pairs whose cron
// schedules share at least one common firing minute within a standard week
// (10 080 minutes). Only valid entries are considered.
func DetectOverlaps(r Report) OverlapReport {
	valid := FilterValid(r).Entries
	var overlaps []OverlapResult

	for i := 0; i < len(valid); i++ {
		for j := i + 1; j < len(valid); j++ {
			a, b := valid[i], valid[j]
			if ex, ok := firstCommonMinute(a.Schedule, b.Schedule); ok {
				overlaps = append(overlaps, OverlapResult{
					A:       a,
					B:       b,
					Example: ex,
				})
			}
		}
	}

	return OverlapReport{
		Source:   r.Source,
		Overlaps: overlaps,
	}
}

// firstCommonMinute returns a descriptive string of the first minute-of-week
// where both schedules fire, or ("" , false) if none is found within one week.
func firstCommonMinute(schedA, schedB string) (string, bool) {
	const weekMinutes = 10080
	for m := 0; m < weekMinutes; m++ {
		day := (m / 1440) % 7
		hour := (m % 1440) / 60
		min := m % 60

		if matchesCron(schedA, min, hour, day) && matchesCron(schedB, min, hour, day) {
			return fmt.Sprintf("day=%d hour=%02d min=%02d", day, hour, min), true
		}
	}
	return "", false
}
