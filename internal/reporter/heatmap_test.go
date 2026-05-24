package reporter

import (
	"strings"
	"testing"
)

func buildHeatmapReport(entries []Entry) Report {
	return Report{
		Source:  "test",
		Entries: entries,
	}
}

func TestBuildHeatmap_SkipsInvalidEntries(t *testing.T) {
	r := buildHeatmapReport([]Entry{
		{Raw: "0 9 * * 1 /backup.sh", Valid: false},
	})
	result := BuildHeatmap(r)
	if len(result.Cells) != 0 {
		t.Errorf("expected 0 cells for invalid entries, got %d", len(result.Cells))
	}
}

func TestBuildHeatmap_SingleEntry(t *testing.T) {
	// Runs at 09:00 on Monday (dow=1)
	r := buildHeatmapReport([]Entry{
		{Raw: "0 9 * * 1 /backup.sh", Valid: true},
	})
	result := BuildHeatmap(r)
	if len(result.Cells) != 1 {
		t.Fatalf("expected 1 cell, got %d", len(result.Cells))
	}
	c := result.Cells[0]
	if c.Day != 1 || c.Hour != 9 || c.Count != 1 {
		t.Errorf("unexpected cell: day=%d hour=%d count=%d", c.Day, c.Hour, c.Count)
	}
}

func TestBuildHeatmap_WildcardExpands(t *testing.T) {
	// Runs every hour on Sunday (dow=0)
	r := buildHeatmapReport([]Entry{
		{Raw: "0 * * * 0 /hourly.sh", Valid: true},
	})
	result := BuildHeatmap(r)
	if len(result.Cells) != 24 {
		t.Errorf("expected 24 cells for wildcard hour, got %d", len(result.Cells))
	}
}

func TestBuildHeatmap_StepField(t *testing.T) {
	// Runs every 6 hours on Friday (dow=5)
	r := buildHeatmapReport([]Entry{
		{Raw: "0 */6 * * 5 /every6h.sh", Valid: true},
	})
	result := BuildHeatmap(r)
	// hours 0, 6, 12, 18 => 4 cells
	if len(result.Cells) != 4 {
		t.Errorf("expected 4 cells for */6 step, got %d", len(result.Cells))
	}
}

func TestBuildHeatmap_FormattedContainsDayLabels(t *testing.T) {
	r := buildHeatmapReport([]Entry{
		{Raw: "0 9 * * 1 /backup.sh", Valid: true},
	})
	result := BuildHeatmap(r)
	for _, label := range []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"} {
		if !strings.Contains(result.Formatted, label) {
			t.Errorf("formatted heatmap missing day label %q", label)
		}
	}
}

func TestBuildHeatmap_MultipleCellsAccumulate(t *testing.T) {
	// Two entries both fire at 08:00 on Wednesday (dow=3)
	r := buildHeatmapReport([]Entry{
		{Raw: "0 8 * * 3 /job1.sh", Valid: true},
		{Raw: "0 8 * * 3 /job2.sh", Valid: true},
	})
	result := BuildHeatmap(r)
	if len(result.Cells) != 1 {
		t.Fatalf("expected 1 merged cell, got %d", len(result.Cells))
	}
	if result.Cells[0].Count != 2 {
		t.Errorf("expected count=2, got %d", result.Cells[0].Count)
	}
}
