package reporter

import (
	"testing"

	"github.com/cronaudit/internal/parser"
)

func buildRankReport() Report {
	return Report{
		Source: "test",
		Entries: []parser.Entry{
			{Line: 1, Schedule: "* * * * *", Command: "echo every-minute", Valid: true},
			{Line: 2, Schedule: "0 9 * * 1", Command: "echo weekly", Valid: true},
			{Line: 3, Schedule: "0 6 1 1 *", Command: "echo yearly", Valid: true},
			{Line: 4, Schedule: "bad schedule", Command: "echo bad", Valid: false, Error: "parse error"},
			{Line: 5, Schedule: "30 8 15 * *", Command: "echo monthly", Valid: true},
		},
	}
}

func TestRankReport_InvalidScoresZero(t *testing.T) {
	r := buildRankReport()
	ranked := RankReport(r)
	for _, re := range ranked {
		if !re.Entry.Valid && re.Score != 0 {
			t.Errorf("expected invalid entry to have score 0, got %f", re.Score)
		}
	}
}

func TestRankReport_SortedDescending(t *testing.T) {
	r := buildRankReport()
	ranked := RankReport(r)
	for i := 1; i < len(ranked); i++ {
		if ranked[i].Score > ranked[i-1].Score {
			t.Errorf("expected descending order at index %d: %f > %f", i, ranked[i].Score, ranked[i-1].Score)
		}
	}
}

func TestRankReport_FullySpecifiedScoresHighest(t *testing.T) {
	r := Report{
		Entries: []parser.Entry{
			{Schedule: "* * * * *", Command: "a", Valid: true},
			{Schedule: "0 6 1 1 0", Command: "b", Valid: true},
		},
	}
	ranked := RankReport(r)
	if ranked[0].Entry.Command != "b" {
		t.Errorf("expected fully specified entry to rank first, got %q", ranked[0].Entry.Command)
	}
}

func TestTopN_ReturnsAtMostN(t *testing.T) {
	r := buildRankReport()
	top := TopN(r, 2)
	if len(top) != 2 {
		t.Errorf("expected 2 entries, got %d", len(top))
	}
}

func TestTopN_NLargerThanEntries(t *testing.T) {
	r := buildRankReport()
	top := TopN(r, 100)
	if len(top) != len(r.Entries) {
		t.Errorf("expected %d entries, got %d", len(r.Entries), len(top))
	}
}

func TestRankReport_ReasonNotEmpty(t *testing.T) {
	r := buildRankReport()
	ranked := RankReport(r)
	for _, re := range ranked {
		if re.Reason == "" {
			t.Errorf("expected non-empty reason for entry %q", re.Entry.Command)
		}
	}
}
