package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// CronField represents a single field in a cron expression.
type CronField struct {
	Raw   string
	Min   int
	Max   int
	Label string
}

// CronEntry represents a parsed crontab entry.
type CronEntry struct {
	Raw     string
	Minute  CronField
	Hour    CronField
	Day     CronField
	Month   CronField
	Weekday CronField
	Command string
}

var fieldDefs = []struct {
	label string
	min   int
	max   int
}{
	{"minute", 0, 59},
	{"hour", 0, 23},
	{"day", 1, 31},
	{"month", 1, 12},
	{"weekday", 0, 7},
}

// Parse parses a single cron line into a CronEntry.
func Parse(line string) (*CronEntry, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, nil
	}

	parts := strings.Fields(line)
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid cron entry: expected at least 6 fields, got %d", len(parts))
	}

	fields := []CronField{}
	for i, def := range fieldDefs {
		if err := validateField(parts[i], def.min, def.max); err != nil {
			return nil, fmt.Errorf("invalid %s field %q: %w", def.label, parts[i], err)
		}
		fields = append(fields, CronField{
			Raw:   parts[i],
			Min:   def.min,
			Max:   def.max,
			Label: def.label,
		})
	}

	return &CronEntry{
		Raw:     line,
		Minute:  fields[0],
		Hour:    fields[1],
		Day:     fields[2],
		Month:   fields[3],
		Weekday: fields[4],
		Command: strings.Join(parts[5:], " "),
	}, nil
}

func validateField(field string, min, max int) error {
	if field == "*" {
		return nil
	}
	if strings.Contains(field, "/") {
		parts := strings.SplitN(field, "/", 2)
		step, err := strconv.Atoi(parts[1])
		if err != nil || step < 1 {
			return fmt.Errorf("invalid step value")
		}
		if parts[0] != "*" {
			return validateRange(parts[0], min, max)
		}
		return nil
	}
	if strings.Contains(field, ",") {
		for _, v := range strings.Split(field, ",") {
			if err := validateRange(v, min, max); err != nil {
				return err
			}
		}
		return nil
	}
	return validateRange(field, min, max)
}

func validateRange(field string, min, max int) error {
	if strings.Contains(field, "-") {
		parts := strings.SplitN(field, "-", 2)
		lo, err1 := strconv.Atoi(parts[0])
		hi, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("non-numeric range")
		}
		if lo < min || hi > max || lo > hi {
			return fmt.Errorf("range %d-%d out of bounds [%d-%d]", lo, hi, min, max)
		}
		return nil
	}
	n, err := strconv.Atoi(field)
	if err != nil {
		return fmt.Errorf("non-numeric value %q", field)
	}
	if n < min || n > max {
		return fmt.Errorf("value %d out of bounds [%d-%d]", n, min, max)
	}
	return nil
}
