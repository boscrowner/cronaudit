package reporter

import (
	"fmt"
	"strings"
)

// RetentionEntry holds a cron entry paired with an estimated data retention
// window based on how frequently the job runs.
type RetentionEntry struct {
	Entry     Entry
	Frequency string // e.g. "every minute", "hourly", "daily"
	WindowDays int    // suggested retention window in days
	Rationale  string
}

// RetentionReport is the result of analysing a report for retention windows.
type RetentionReport struct {
	Source  string
	Entries []RetentionEntry
}

// BuildRetentionReport analyses each valid entry and suggests a log/data
// retention window based on its run frequency.
func BuildRetentionReport(r Report) RetentionReport {
	rr := RetentionReport{Source: r.Source}

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}

		freq, days, rationale := estimateRetention(e)
		rr.Entries = append(rr.Entries, RetentionEntry{
			Entry:      e,
			Frequency:  freq,
			WindowDays: days,
			Rationale:  rationale,
		})
	}

	return rr
}

// FormatRetentionReport returns a human-readable summary of the retention report.
func FormatRetentionReport(rr RetentionReport) string {
	if len(rr.Entries) == 0 {
		return "No valid entries found for retention analysis."
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Retention Report — %s\n", rr.Source)
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 60))

	for _, re := range rr.Entries {
		fmt.Fprintf(&sb, "Schedule : %s\n", re.Entry.Schedule)
		fmt.Fprintf(&sb, "Command  : %s\n", re.Entry.Command)
		fmt.Fprintf(&sb, "Frequency: %s\n", re.Frequency)
		fmt.Fprintf(&sb, "Suggested retention: %d day(s)\n", re.WindowDays)
		fmt.Fprintf(&sb, "Rationale: %s\n\n", re.Rationale)
	}

	return strings.TrimRight(sb.String(), "\n")
}

// estimateRetention returns a frequency label, window in days, and rationale
// for the given entry.
func estimateRetention(e Entry) (string, int, string) {
	parts := strings.Fields(e.Schedule)
	if len(parts) < 5 {
		return "unknown", 30, "Could not parse schedule fields."
	}

	minute, hour, dom, month, dow := parts[0], parts[1], parts[2], parts[3], parts[4]

	switch {
	case minute == "*" && hour == "*":
		return "every minute", 7, "High-frequency job; 7-day window prevents unbounded log growth."
	case minute != "*" && hour == "*":
		return "hourly", 14, "Runs every hour; 14-day window balances storage and auditability."
	case minute != "*" && hour != "*" && dom == "*" && month == "*" && dow == "*":
		return "daily", 90, "Daily job; 90-day window supports month-over-month comparison."
	case dow != "*" && dom == "*":
		return "weekly", 180, "Weekly job; 180-day window covers a 6-month audit horizon."
	case dom != "*" && month == "*":
		return "monthly", 365, "Monthly job; 1-year window supports annual reporting."
	case month != "*":
		return "yearly", 730, "Infrequent job; 2-year window ensures historical records are kept."
	default:
		return "custom", 90, "Custom schedule; defaulting to 90-day retention window."
	}
}
