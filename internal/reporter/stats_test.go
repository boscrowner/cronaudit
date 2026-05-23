package reporter

import (
	"testing"
)

func buildStatsReport() Report {
	return Report{
		Source: "test",
		Stats: struct {
			Total   int
			Valid   int
			Invalid int
		}{Total: 4, Valid: 3, Invalid: 1},
		Entries: []Entry{
			{Line: 1, Raw: "* * * * * cmd1", Valid: true},
			{Line: 2, Raw: "* * * * * cmd2", Valid: true},
			{Line: 3, Raw: "* * * * * cmd3", Valid: true},
			{Line: 4, Raw: "bad entry", Valid: false, Error: "invalid format"},
		},
	}
}

func TestComputeStats_Counts(t *testing.T) {
	r := buildStatsReport()
	s := ComputeStats(r)

	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Valid != 3 {
		t.Errorf("expected Valid=3, got %d", s.Valid)
	}
	if s.Invalid != 1 {
		t.Errorf("expected Invalid=1, got %d", s.Invalid)
	}
}

func TestStats_PassRate(t *testing.T) {
	s := Stats{Total: 4, Valid: 3, Invalid: 1}
	got := s.PassRate()
	want := 0.75
	if got != want {
		t.Errorf("expected PassRate=%v, got %v", want, got)
	}
}

func TestStats_PassRate_ZeroTotal(t *testing.T) {
	s := Stats{}
	if s.PassRate() != 0 {
		t.Errorf("expected PassRate=0 for empty stats")
	}
}

func TestStats_HasErrors_True(t *testing.T) {
	s := Stats{Total: 2, Valid: 1, Invalid: 1}
	if !s.HasErrors() {
		t.Error("expected HasErrors=true")
	}
}

func TestStats_HasErrors_False(t *testing.T) {
	s := Stats{Total: 2, Valid: 2, Invalid: 0}
	if s.HasErrors() {
		t.Error("expected HasErrors=false")
	}
}

func TestComputeStats_EmptyReport(t *testing.T) {
	r := Report{Source: "empty"}
	s := ComputeStats(r)
	if s.Valid != 0 || s.Invalid != 0 {
		t.Errorf("expected zero counts for empty report, got %+v", s)
	}
}
