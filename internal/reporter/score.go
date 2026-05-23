package reporter

import "fmt"

// Score represents a quality score for a crontab report.
type Score struct {
	// Value is the numeric score from 0 to 100.
	Value int
	// Grade is the letter grade derived from Value.
	Grade string
	// Summary is a human-readable description of the score.
	Summary string
}

// ScoreReport computes a quality score for the given Report based on the
// ratio of valid entries to total entries, penalising duplicate schedules.
//
// Scoring breakdown:
//   - Base score: (valid / total) * 100
//   - Duplicate penalty: 5 points per duplicate schedule key
func ScoreReport(r Report) Score {
	stats := ComputeStats(r)

	if stats.Total == 0 {
		return Score{Value: 100, Grade: "A", Summary: "No entries to evaluate."}
	}

	base := int((float64(stats.Valid) / float64(stats.Total)) * 100)

	dupes := DuplicateSchedules(r)
	penalty := len(dupes) * 5

	value := base - penalty
	if value < 0 {
		value = 0
	}

	grade := letterGrade(value)
	summary := buildSummary(value, grade, stats.Valid, stats.Total, len(dupes))

	return Score{
		Value:   value,
		Grade:   grade,
		Summary: summary,
	}
}

func letterGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}

func buildSummary(value int, grade string, valid, total, dupes int) string {
	base := fmt.Sprintf("Score %d/100 (%s): %d/%d entries valid", value, grade, valid, total)
	if dupes > 0 {
		return fmt.Sprintf("%s, %d duplicate schedule(s) penalised.", base, dupes)
	}
	return base + "."
}
