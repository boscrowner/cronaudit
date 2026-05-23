package parser

import (
	"testing"
)

func TestParse_ValidEntries(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		cmd     string
		minute  string
		hour    string
	}{
		{
			name:   "every minute",
			line:   "* * * * * /usr/bin/backup",
			cmd:    "/usr/bin/backup",
			minute: "*",
			hour:   "*",
		},
		{
			name:   "specific time",
			line:   "30 6 * * 1 /scripts/weekly.sh",
			cmd:    "/scripts/weekly.sh",
			minute: "30",
			hour:   "6",
		},
		{
			name:   "step expression",
			line:   "*/15 * * * * /check.sh",
			cmd:    "/check.sh",
			minute: "*/15",
			hour:   "*",
		},
		{
			name:   "range expression",
			line:   "0 9-17 * * 1-5 /work.sh",
			cmd:    "/work.sh",
			minute: "0",
			hour:   "9-17",
		},
		{
			name:   "list expression",
			line:   "0 0,12 * * * /twice-daily.sh",
			cmd:    "/twice-daily.sh",
			minute: "0",
			hour:   "0,12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := Parse(tt.line)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if entry == nil {
				t.Fatal("expected entry, got nil")
			}
			if entry.Command != tt.cmd {
				t.Errorf("command: got %q, want %q", entry.Command, tt.cmd)
			}
			if entry.Minute.Raw != tt.minute {
				t.Errorf("minute: got %q, want %q", entry.Minute.Raw, tt.minute)
			}
			if entry.Hour.Raw != tt.hour {
				t.Errorf("hour: got %q, want %q", entry.Hour.Raw, tt.hour)
			}
		})
	}
}

func TestParse_InvalidEntries(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{"too few fields", "* * * * /cmd"},
		{"minute out of range", "60 * * * * /cmd"},
		{"hour out of range", "* 24 * * * /cmd"},
		{"invalid step", "*/0 * * * * /cmd"},
		{"invalid range", "* 10-25 * * * /cmd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := Parse(tt.line)
			if err == nil {
				t.Errorf("expected error for %q, got entry: %+v", tt.line, entry)
			}
		})
	}
}

func TestParse_SkipsCommentAndBlank(t *testing.T) {
	for _, line := range []string{"", "  ", "# this is a comment"} {
		entry, err := Parse(line)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", line, err)
		}
		if entry != nil {
			t.Errorf("expected nil entry for %q", line)
		}
	}
}
