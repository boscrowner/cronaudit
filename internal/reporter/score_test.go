package reporter

import (
	"strings"
	"testing"
)

func buildScoreReport(entries []Entry) Report {
	return Report{
		Source:  "score_test",
		Entries: entries,
		Stats:   Stats{Total: len(entries)},
	}
}

func TestScoreReport_PerfectScore(t *testing.T) {
	r := buildScoreReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "echo a", Valid: true},
		{Line: 2, Schedule: "5 * * * *", Command: "echo b", Valid: true},
	})
	s := ScoreReport(r)
	if s.Value != 100 {
		t.Errorf("expected 100, got %d", s.Value)
	}
	if s.Grade != "A" {
		t.Errorf("expected grade A, got %s", s.Grade)
	}
}

func TestScoreReport_HalfValid(t *testing.T) {
	r := buildScoreReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "echo a", Valid: true},
		{Line: 2, Schedule: "bad", Command: "echo b", Valid: false, Error: "invalid"},
	})
	s := ScoreReport(r)
	if s.Value != 50 {
		t.Errorf("expected 50, got %d", s.Value)
	}
	if s.Grade != "F" {
		t.Errorf("expected grade F, got %s", s.Grade)
	}
}

func TestScoreReport_DuplicatePenalty(t *testing.T) {
	r := buildScoreReport([]Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "echo a", Valid: true},
		{Line: 2, Schedule: "0 * * * *", Command: "echo b", Valid: true},
	})
	s := ScoreReport(r)
	// base 100, penalty 5 for one duplicate key
	if s.Value != 95 {
		t.Errorf("expected 95, got %d", s.Value)
	}
}

func TestScoreReport_EmptyReport(t *testing.T) {
	r := buildScoreReport([]Entry{})
	s := ScoreReport(r)
	if s.Value != 100 {
		t.Errorf("expected 100 for empty report, got %d", s.Value)
	}
	if !strings.Contains(s.Summary, "No entries") {
		t.Errorf("expected 'No entries' in summary, got: %s", s.Summary)
	}
}

func TestScoreReport_ScoreNeverNegative(t *testing.T) {
	entries := make([]Entry, 20)
	for i := range entries {
		entries[i] = Entry{Line: i + 1, Schedule: "0 * * * *", Command: "echo x", Valid: true}
	}
	r := buildScoreReport(entries)
	s := ScoreReport(r)
	if s.Value < 0 {
		t.Errorf("score should not be negative, got %d", s.Value)
	}
}

func TestLetterGrade_Boundaries(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{100, "A"}, {90, "A"}, {89, "B"}, {75, "B"},
		{74, "C"}, {60, "C"}, {59, "D"}, {40, "D"},
		{39, "F"}, {0, "F"},
	}
	for _, tc := range cases {
		got := letterGrade(tc.score)
		if got != tc.want {
			t.Errorf("letterGrade(%d) = %s, want %s", tc.score, got, tc.want)
		}
	}
}
