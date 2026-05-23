package reporter

import (
	"sort"
	"strings"
)

// Cluster represents a group of report entries sharing a common schedule pattern.
type Cluster struct {
	Pattern  string
	Entries  []ReportEntry
	Size     int
}

// ClusterByHour groups entries that run at the same minute and hour fields,
// regardless of day/month/weekday. Clusters are sorted by size descending.
func ClusterByHour(r Report) []Cluster {
	return clusterBy(r, func(e ReportEntry) string {
		if !e.Valid {
			return ""
		}
		parts := strings.Fields(e.Schedule)
		if len(parts) < 2 {
			return ""
		}
		return parts[0] + " " + parts[1]
	})
}

// ClusterByDayOfWeek groups entries that share the same weekday field.
// Clusters are sorted by size descending.
func ClusterByDayOfWeek(r Report) []Cluster {
	return clusterBy(r, func(e ReportEntry) string {
		if !e.Valid {
			return ""
		}
		parts := strings.Fields(e.Schedule)
		if len(parts) < 5 {
			return ""
		}
		return parts[4]
	})
}

func clusterBy(r Report, keyFn func(ReportEntry) string) []Cluster {
	index := make(map[string][]ReportEntry)
	for _, e := range r.Entries {
		k := keyFn(e)
		if k == "" {
			continue
		}
		index[k] = append(index[k], e)
	}

	clusters := make([]Cluster, 0, len(index))
	for pattern, entries := range index {
		clusters = append(clusters, Cluster{
			Pattern: pattern,
			Entries: entries,
			Size:    len(entries),
		})
	}

	sort.Slice(clusters, func(i, j int) bool {
		if clusters[i].Size != clusters[j].Size {
			return clusters[i].Size > clusters[j].Size
		}
		return clusters[i].Pattern < clusters[j].Pattern
	})

	return clusters
}
