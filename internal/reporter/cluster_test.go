package reporter

import (
	"testing"
)

func buildClusterReport() Report {
	return Report{
		Source: "test",
		Entries: []ReportEntry{
			{Line: 1, Schedule: "0 2 * * 1", Command: "cmd-a", Valid: true},
			{Line: 2, Schedule: "0 2 * * 2", Command: "cmd-b", Valid: true},
			{Line: 3, Schedule: "0 2 * * 3", Command: "cmd-c", Valid: true},
			{Line: 4, Schedule: "30 6 * * 1", Command: "cmd-d", Valid: true},
			{Line: 5, Schedule: "30 6 * * 1", Command: "cmd-e", Valid: true},
			{Line: 6, Schedule: "bad schedule", Command: "cmd-f", Valid: false},
		},
	}
}

func TestClusterByHour_GroupsOnMinuteAndHour(t *testing.T) {
	r := buildClusterReport()
	clusters := ClusterByHour(r)

	patterns := make(map[string]int)
	for _, c := range clusters {
		patterns[c.Pattern] = c.Size
	}

	if patterns["0 2"] != 3 {
		t.Errorf("expected pattern '0 2' to have 3 entries, got %d", patterns["0 2"])
	}
	if patterns["30 6"] != 2 {
		t.Errorf("expected pattern '30 6' to have 2 entries, got %d", patterns["30 6"])
	}
}

func TestClusterByHour_SkipsInvalidEntries(t *testing.T) {
	r := buildClusterReport()
	clusters := ClusterByHour(r)

	for _, c := range clusters {
		for _, e := range c.Entries {
			if !e.Valid {
				t.Errorf("cluster should not contain invalid entry at line %d", e.Line)
			}
		}
	}
}

func TestClusterByHour_SortedBySizeDescending(t *testing.T) {
	r := buildClusterReport()
	clusters := ClusterByHour(r)

	for i := 1; i < len(clusters); i++ {
		if clusters[i].Size > clusters[i-1].Size {
			t.Errorf("clusters not sorted by size: index %d (%d) > index %d (%d)",
				i, clusters[i].Size, i-1, clusters[i-1].Size)
		}
	}
}

func TestClusterByDayOfWeek_GroupsByWeekday(t *testing.T) {
	r := buildClusterReport()
	clusters := ClusterByDayOfWeek(r)

	patterns := make(map[string]int)
	for _, c := range clusters {
		patterns[c.Pattern] = c.Size
	}

	if patterns["1"] != 3 {
		t.Errorf("expected weekday '1' to have 3 entries, got %d", patterns["1"])
	}
	if patterns["2"] != 1 {
		t.Errorf("expected weekday '2' to have 1 entry, got %d", patterns["2"])
	}
}

func TestClusterByHour_DoesNotMutateOriginal(t *testing.T) {
	r := buildClusterReport()
	origLen := len(r.Entries)
	ClusterByHour(r)
	if len(r.Entries) != origLen {
		t.Errorf("original report mutated: expected %d entries, got %d", origLen, len(r.Entries))
	}
}
