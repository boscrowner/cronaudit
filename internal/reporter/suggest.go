package reporter

import (
	"fmt"
	"strings"
)

// Suggestion holds a recommended fix or improvement for a cron entry.
type Suggestion struct {
	Line     int
	Schedule string
	Command  string
	Message  string
}

// SuggestReport analyzes a Report and returns actionable suggestions
// for improving cron entries. It targets common scheduling anti-patterns.
func SuggestReport(r Report) []Suggestion {
	var suggestions []Suggestion

	for _, entry := range r.Entries {
		if !entry.Valid {
			continue
		}

		parts := strings.Fields(entry.Raw)
		if len(parts) < 6 {
			continue
		}

		minute := parts[0]
		hour := parts[1]
		dom := parts[2]
		month := parts[3]
		dow := parts[4]

		// Suggest using @daily instead of "0 0 * * *"
		if minute == "0" && hour == "0" && dom == "*" && month == "*" && dow == "*" {
			suggestions = append(suggestions, Suggestion{
				Line:     entry.Line,
				Schedule: entry.Schedule,
				Command:  entry.Command,
				Message:  "consider using @daily instead of '0 0 * * *'",
			})
			continue
		}

		// Suggest using @hourly instead of "0 * * * *"
		if minute == "0" && hour == "*" && dom == "*" && month == "*" && dow == "*" {
			suggestions = append(suggestions, Suggestion{
				Line:     entry.Line,
				Schedule: entry.Schedule,
				Command:  entry.Command,
				Message:  "consider using @hourly instead of '0 * * * *'",
			})
			continue
		}

		// Suggest using @weekly instead of "0 0 * * 0"
		if minute == "0" && hour == "0" && dom == "*" && month == "*" && dow == "0" {
			suggestions = append(suggestions, Suggestion{
				Line:     entry.Line,
				Schedule: entry.Schedule,
				Command:  entry.Command,
				Message:  "consider using @weekly instead of '0 0 * * 0'",
			})
			continue
		}

		// Warn about running every minute on a heavy step
		if minute == "*" && hour == "*" {
			suggestions = append(suggestions, Suggestion{
				Line:     entry.Line,
				Schedule: entry.Schedule,
				Command:  entry.Command,
				Message:  fmt.Sprintf("entry runs every minute; verify '%s' is lightweight", entry.Command),
			})
		}
	}

	return suggestions
}
