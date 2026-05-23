package reporter

import "sort"

// GroupedEntries holds report entries grouped by a common key.
type GroupedEntries struct {
	Key     string
	Entries []ReportEntry
}

// GroupByCommand groups report entries by their command field.
// Entries without a command (invalid) are grouped under "(invalid)".
func GroupByCommand(r Report) []GroupedEntries {
	return groupBy(r, func(e ReportEntry) string {
		if e.Command == "" {
			return "(invalid)"
		}
		return e.Command
	})
}

// GroupBySchedule groups report entries by their raw schedule expression.
func GroupBySchedule(r Report) []GroupedEntries {
	return groupBy(r, func(e ReportEntry) string {
		if e.Schedule == "" {
			return "(invalid)"
		}
		return e.Schedule
	})
}

// groupBy is the internal generic grouping function.
func groupBy(r Report, keyFn func(ReportEntry) string) []GroupedEntries {
	index := make(map[string][]ReportEntry)
	order := []string{}

	for _, entry := range r.Entries {
		key := keyFn(entry)
		if _, exists := index[key]; !exists {
			order = append(order, key)
		}
		index[key] = append(index[key], entry)
	}

	sort.Strings(order)

	result := make([]GroupedEntries, 0, len(order))
	for _, key := range order {
		result = append(result, GroupedEntries{
			Key:     key,
			Entries: index[key],
		})
	}
	return result
}
