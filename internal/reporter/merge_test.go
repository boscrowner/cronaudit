package reporter

import (
	"testing"
)

func buildMergeReport(source string, entries []Entry) Report {
	return Report{
		Source:  source,
		Entries: entries,
		Stats:   ComputeStats(entries),
	}
}

func TestMergeReports_CombinesEntries(t *testing.T) {
	a := buildMergeReport("a.cron", []Entry{
		{Line: 1, Raw: "* * * * * echo a", Valid: true},
	})
	b := buildMergeReport("b.cron", []Entry{
		{Line: 1, Raw: "0 9 * * * echo b", Valid: true},
	})

	result := MergeReports(a, b)

	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestMergeReports_SourceFromA(t *testing.T) {
	a := buildMergeReport("primary.cron", []Entry{{Line: 1, Raw: "* * * * * echo a", Valid: true}})
	b := buildMergeReport("secondary.cron", []Entry{{Line: 1, Raw: "* * * * * echo b", Valid: true}})

	result := MergeReports(a, b)

	if result.Source != "primary.cron" {
		t.Errorf("expected source %q, got %q", "primary.cron", result.Source)
	}
}

func TestMergeReports_StatsRecomputed(t *testing.T) {
	a := buildMergeReport("a.cron", []Entry{
		{Line: 1, Raw: "* * * * * echo a", Valid: true},
	})
	b := buildMergeReport("b.cron", []Entry{
		{Line: 1, Raw: "bad entry", Valid: false, Error: "invalid"},
	})

	result := MergeReports(a, b)

	if result.Stats.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Stats.Total)
	}
	if result.Stats.Valid != 1 {
		t.Errorf("expected valid 1, got %d", result.Stats.Valid)
	}
	if result.Stats.Invalid != 1 {
		t.Errorf("expected invalid 1, got %d", result.Stats.Invalid)
	}
}

func TestMergeReports_DoesNotMutateOriginals(t *testing.T) {
	a := buildMergeReport("a.cron", []Entry{{Line: 1, Raw: "* * * * * echo a", Valid: true}})
	b := buildMergeReport("b.cron", []Entry{{Line: 2, Raw: "* * * * * echo b", Valid: true}})

	_ = MergeReports(a, b)

	if len(a.Entries) != 1 {
		t.Errorf("original report a was mutated")
	}
	if len(b.Entries) != 1 {
		t.Errorf("original report b was mutated")
	}
}

func TestMergeMany_EmptySlice(t *testing.T) {
	result := MergeMany([]Report{})

	if result.Source != "" || len(result.Entries) != 0 {
		t.Errorf("expected empty report, got %+v", result)
	}
}

func TestMergeMany_PicksFirstSource(t *testing.T) {
	reports := []Report{
		buildMergeReport("", []Entry{}),
		buildMergeReport("second.cron", []Entry{{Line: 1, Raw: "* * * * * echo x", Valid: true}}),
		buildMergeReport("third.cron", []Entry{{Line: 2, Raw: "* * * * * echo y", Valid: true}}),
	}

	result := MergeMany(reports)

	if result.Source != "second.cron" {
		t.Errorf("expected source %q, got %q", "second.cron", result.Source)
	}
	if len(result.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Entries))
	}
}
