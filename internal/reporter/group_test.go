package reporter

import (
	"testing"
)

func buildGroupReport() Report {
	return Report{
		Source: "test",
		Entries: []ReportEntry{
			{Line: 1, Schedule: "0 * * * *", Command: "/bin/alpha", Valid: true},
			{Line: 2, Schedule: "0 * * * *", Command: "/bin/beta", Valid: true},
			{Line: 3, Schedule: "30 6 * * *", Command: "/bin/alpha", Valid: true},
			{Line: 4, Schedule: "@daily", Command: "/bin/gamma", Valid: true},
			{Line: 5, Schedule: "not-valid", Command: "", Valid: false, Error: "invalid"},
		},
	}
}

func TestGroupByCommand_GroupsCorrectly(t *testing.T) {
	r := buildGroupReport()
	groups := GroupByCommand(r)

	groupMap := make(map[string]int)
	for _, g := range groups {
		groupMap[g.Key] = len(g.Entries)
	}

	if groupMap["/bin/alpha"] != 2 {
		t.Errorf("expected 2 entries for /bin/alpha, got %d", groupMap["/bin/alpha"])
	}
	if groupMap["/bin/beta"] != 1 {
		t.Errorf("expected 1 entry for /bin/beta, got %d", groupMap["/bin/beta"])
	}
	if groupMap["(invalid)"] != 1 {
		t.Errorf("expected 1 entry for (invalid), got %d", groupMap["(invalid)"])
	}
}

func TestGroupBySchedule_GroupsCorrectly(t *testing.T) {
	r := buildGroupReport()
	groups := GroupBySchedule(r)

	groupMap := make(map[string]int)
	for _, g := range groups {
		groupMap[g.Key] = len(g.Entries)
	}

	if groupMap["0 * * * *"] != 2 {
		t.Errorf("expected 2 entries for '0 * * * *', got %d", groupMap["0 * * * *"])
	}
	if groupMap["30 6 * * *"] != 1 {
		t.Errorf("expected 1 entry for '30 6 * * *', got %d", groupMap["30 6 * * *"])
	}
}

func TestGroupByCommand_SortedKeys(t *testing.T) {
	r := buildGroupReport()
	groups := GroupByCommand(r)

	for i := 1; i < len(groups); i++ {
		if groups[i-1].Key > groups[i].Key {
			t.Errorf("keys not sorted: %q > %q", groups[i-1].Key, groups[i].Key)
		}
	}
}

func TestGroupByCommand_DoesNotMutateOriginal(t *testing.T) {
	r := buildGroupReport()
	origLen := len(r.Entries)
	GroupByCommand(r)
	if len(r.Entries) != origLen {
		t.Errorf("original report mutated: expected %d entries, got %d", origLen, len(r.Entries))
	}
}
