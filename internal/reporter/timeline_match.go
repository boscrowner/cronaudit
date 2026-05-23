package reporter

import (
	"strconv"
	"strings"
	"time"
)

// splitFields splits a cron schedule string into its 5 fields.
func splitFields(schedule string) []string {
	return strings.Fields(schedule)
}

// matchesCron returns true if the given time matches all five cron fields.
func matchesCron(fields []string, t time.Time) bool {
	return matchField(fields[0], t.Minute(), 0, 59) &&
		matchField(fields[1], t.Hour(), 0, 23) &&
		matchField(fields[2], t.Day(), 1, 31) &&
		matchField(fields[3], int(t.Month()), 1, 12) &&
		matchField(fields[4], int(t.Weekday()), 0, 6)
}

// matchField checks whether a single cron field matches the given value.
func matchField(field string, value, min, max int) bool {
	if field == "*" {
		return true
	}
	// Handle step values: */n
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(field[2:])
		if err != nil || step <= 0 {
			return false
		}
		return (value-min)%step == 0
	}
	// Handle comma-separated lists
	for _, part := range strings.Split(field, ",") {
		if matchSinglePart(part, value, min, max) {
			return true
		}
	}
	return false
}

// matchSinglePart handles a range (a-b), step range (a-b/n), or literal value.
func matchSinglePart(part string, value, min, max int) bool {
	if strings.Contains(part, "-") {
		rangeParts := strings.SplitN(part, "-", 2)
		lo, err1 := strconv.Atoi(rangeParts[0])
		stepStr := rangeParts[1]
		step := 1
		if strings.Contains(stepStr, "/") {
			sp := strings.SplitN(stepStr, "/", 2)
			stepStr = sp[0]
			var err error
			step, err = strconv.Atoi(sp[1])
			if err != nil || step <= 0 {
				return false
			}
		}
		hi, err2 := strconv.Atoi(stepStr)
		if err1 != nil || err2 != nil {
			return false
		}
		if value < lo || value > hi {
			return false
		}
		return (value-lo)%step == 0
	}
	n, err := strconv.Atoi(part)
	return err == nil && n == value
}
