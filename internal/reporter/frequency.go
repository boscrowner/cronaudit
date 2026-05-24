package reporter

import (
	"fmt"
	"sort"
	"strings"
)

// FrequencyEntry holds a schedule pattern and the number of cron entries that use it.
type FrequencyEntry struct {
	Schedule string
	Count    int
	Commands []string
}

// FrequencyReport holds the result of a frequency analysis.
type FrequencyReport struct {
	Source  string
	Entries []FrequencyEntry
}

// FrequencyReport builds a frequency table of schedule patterns across all valid entries.
// Entries are sorted by count descending, then schedule ascending.
func BuildFrequencyReport(r Report) FrequencyReport {
	type bucket struct {
		count    int
		commands []string
	}

	index := make(map[string]*bucket)

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		parts := strings.Fields(e.Raw)
		if len(parts) < 6 {
			continue
		}
		schedule := strings.Join(parts[:5], " ")
		command := strings.Join(parts[5:], " ")

		if _, ok := index[schedule]; !ok {
			index[schedule] = &bucket{}
		}
		index[schedule].count++
		index[schedule].commands = append(index[schedule].commands, command)
	}

	result := make([]FrequencyEntry, 0, len(index))
	for sched, b := range index {
		result = append(result, FrequencyEntry{
			Schedule: sched,
			Count:    b.count,
			Commands: b.commands,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Schedule < result[j].Schedule
	})

	return FrequencyReport{
		Source:  r.Source,
		Entries: result,
	}
}

// FormatFrequencyReport returns a human-readable summary of the frequency report.
func FormatFrequencyReport(fr FrequencyReport) string {
	if len(fr.Entries) == 0 {
		return "No valid schedule entries found."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Schedule frequency report (%s):\n", fr.Source))
	for _, e := range fr.Entries {
		sb.WriteString(fmt.Sprintf("  [%d] %s\n", e.Count, e.Schedule))
	}
	return strings.TrimRight(sb.String(), "\n")
}
