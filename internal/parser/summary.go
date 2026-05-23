package parser

import (
	"fmt"
	"strings"
)

// Summarize returns a human-readable description of a CronEntry.
func Summarize(e *CronEntry) string {
	parts := []string{}

	parts = append(parts, summarizeMinute(e.Minute.Raw))
	parts = append(parts, summarizeHour(e.Hour.Raw))
	parts = append(parts, summarizeDay(e.Day.Raw))
	parts = append(parts, summarizeMonth(e.Month.Raw))
	parts = append(parts, summarizeWeekday(e.Weekday.Raw))

	filtered := []string{}
	for _, p := range parts {
		if p != "" {
			filtered = append(filtered, p)
		}
	}

	schedule := strings.Join(filtered, ", ")
	if schedule == "" {
		schedule = "every minute"
	}
	return fmt.Sprintf("Runs %s — command: %s", schedule, e.Command)
}

func summarizeMinute(raw string) string {
	switch {
	case raw == "*":
		return ""
	case strings.HasPrefix(raw, "*/"):
		return fmt.Sprintf("every %s minutes", raw[2:])
	case strings.Contains(raw, "-"):
		return fmt.Sprintf("at minutes %s", raw)
	case strings.Contains(raw, ","):
		return fmt.Sprintf("at minutes %s", raw)
	default:
		return fmt.Sprintf("at minute %s", raw)
	}
}

func summarizeHour(raw string) string {
	switch {
	case raw == "*":
		return ""
	case strings.HasPrefix(raw, "*/"):
		return fmt.Sprintf("every %s hours", raw[2:])
	case strings.Contains(raw, "-"):
		return fmt.Sprintf("between hours %s", raw)
	case strings.Contains(raw, ","):
		return fmt.Sprintf("at hours %s", raw)
	default:
		return fmt.Sprintf("at hour %s", raw)
	}
}

func summarizeDay(raw string) string {
	if raw == "*" {
		return ""
	}
	return fmt.Sprintf("on day %s of month", raw)
}

func summarizeMonth(raw string) string {
	monthNames := map[string]string{
		"1": "January", "2": "February", "3": "March", "4": "April",
		"5": "May", "6": "June", "7": "July", "8": "August",
		"9": "September", "10": "October", "11": "November", "12": "December",
	}
	if raw == "*" {
		return ""
	}
	if name, ok := monthNames[raw]; ok {
		return fmt.Sprintf("in %s", name)
	}
	return fmt.Sprintf("in month %s", raw)
}

func summarizeWeekday(raw string) string {
	dayNames := map[string]string{
		"0": "Sunday", "1": "Monday", "2": "Tuesday", "3": "Wednesday",
		"4": "Thursday", "5": "Friday", "6": "Saturday", "7": "Sunday",
	}
	if raw == "*" {
		return ""
	}
	if name, ok := dayNames[raw]; ok {
		return fmt.Sprintf("on %s", name)
	}
	if strings.Contains(raw, "-") {
		return fmt.Sprintf("on weekdays %s", raw)
	}
	return fmt.Sprintf("on weekday %s", raw)
}
