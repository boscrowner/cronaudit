package reporter_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"strings"
	"testing"

	"github.com/example/cronaudit/internal/reporter"
)

func buildExportReport() reporter.Report {
	return reporter.Report{
		Source: "test.crontab",
		Entries: []reporter.Entry{
			{
				Line:     1,
				Schedule: "0 * * * *",
				Command:  "/usr/bin/backup",
				Valid:    true,
				Summary:  "At minute 0 of every hour",
			},
			{
				Line:     2,
				Schedule: "99 * * * *",
				Command:  "/usr/bin/broken",
				Valid:    false,
				Err:      errors.New("invalid minute: 99"),
			},
		},
	}
}

func TestExportCSV_HeaderRow(t *testing.T) {
	var buf bytes.Buffer
	err := reporter.ExportCSV(buildExportReport(), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV output: %v", err)
	}

	expectedHeader := []string{"Line", "Schedule", "Command", "Valid", "Summary", "Error"}
	if len(records) == 0 {
		t.Fatal("expected at least a header row")
	}
	for i, col := range expectedHeader {
		if records[0][i] != col {
			t.Errorf("header col %d: got %q, want %q", i, records[0][i], col)
		}
	}
}

func TestExportCSV_RowCount(t *testing.T) {
	var buf bytes.Buffer
	_ = reporter.ExportCSV(buildExportReport(), &buf)

	r := csv.NewReader(&buf)
	records, _ := r.ReadAll()
	// 1 header + 2 entries
	if len(records) != 3 {
		t.Errorf("expected 3 rows (header + 2 entries), got %d", len(records))
	}
}

func TestExportCSV_ValidEntry(t *testing.T) {
	var buf bytes.Buffer
	_ = reporter.ExportCSV(buildExportReport(), &buf)

	r := csv.NewReader(&buf)
	records, _ := r.ReadAll()

	row := records[1]
	if row[0] != "1" {
		t.Errorf("line: got %q, want %q", row[0], "1")
	}
	if row[3] != "true" {
		t.Errorf("valid: got %q, want %q", row[3], "true")
	}
	if row[5] != "" {
		t.Errorf("error: got %q, want empty", row[5])
	}
}

func TestExportCSV_InvalidEntry_HasError(t *testing.T) {
	var buf bytes.Buffer
	_ = reporter.ExportCSV(buildExportReport(), &buf)

	r := csv.NewReader(&buf)
	records, _ := r.ReadAll()

	row := records[2]
	if row[3] != "false" {
		t.Errorf("valid: got %q, want %q", row[3], "false")
	}
	if !strings.Contains(row[5], "invalid minute") {
		t.Errorf("error field missing message, got %q", row[5])
	}
}
