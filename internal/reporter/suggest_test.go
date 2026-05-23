package reporter

import (
	"testing"
)

func buildSuggestReport(raws []string) Report {
	var entries []Entry
	for i, raw := range raws {
		parts := splitFields(raw)
		valid := len(parts) >= 6
		schedule := ""
		command := ""
		if valid {
			schedule = parts[0] + " " + parts[1] + " " + parts[2] + " " + parts[3] + " " + parts[4]
			command = parts[5]
		}
		entries = append(entries, Entry{
			Line:     i + 1,
			Raw:      raw,
			Schedule: schedule,
			Command:  command,
			Valid:    valid,
		})
	}
	return Report{Source: "test", Entries: entries}
}

func TestSuggestReport_DailyAlias(t *testing.T) {
	r := buildSuggestReport([]string{"0 0 * * * /usr/bin/backup"})
	suggestions := SuggestReport(r)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Message == "" {
		t.Error("expected non-empty suggestion message")
	}
	contains := false
	for _, s := range suggestions {
		if s.Message != "" && len(s.Message) > 0 {
			contains = true
		}
	}
	if !contains {
		t.Error("expected suggestion about @daily")
	}
}

func TestSuggestReport_HourlyAlias(t *testing.T) {
	r := buildSuggestReport([]string{"0 * * * * /usr/bin/cleanup"})
	suggestions := SuggestReport(r)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestSuggestReport_WeeklyAlias(t *testing.T) {
	r := buildSuggestReport([]string{"0 0 * * 0 /usr/bin/report"})
	suggestions := SuggestReport(r)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestSuggestReport_EveryMinuteWarning(t *testing.T) {
	r := buildSuggestReport([]string{"* * * * * /usr/bin/poll"})
	suggestions := SuggestReport(r)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestSuggestReport_SkipsInvalidEntries(t *testing.T) {
	r := buildSuggestReport([]string{"not a valid cron entry"})
	suggestions := SuggestReport(r)
	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions for invalid entry, got %d", len(suggestions))
	}
}

func TestSuggestReport_NoSuggestionForCleanEntry(t *testing.T) {
	r := buildSuggestReport([]string{"30 4 * * 1 /usr/bin/task"})
	suggestions := SuggestReport(r)
	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(suggestions))
	}
}
