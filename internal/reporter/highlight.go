package reporter

import (
	"fmt"
	"strings"
)

// HighlightResult holds the original entry with matched fields annotated.
type HighlightResult struct {
	Line     int
	Schedule string
	Command  string
	Matched  bool
	Snippet  string
}

// HighlightReport searches entries in a Report for those whose schedule or
// command contains the given query string (case-insensitive) and returns
// annotated results with a context snippet.
func HighlightReport(r Report, query string) []HighlightResult {
	if query == "" {
		return nil
	}
	lower := strings.ToLower(query)
	var results []HighlightResult
	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		schedule := e.Raw
		command := e.Command
		matched := strings.Contains(strings.ToLower(schedule), lower) ||
			strings.Contains(strings.ToLower(command), lower)
		if !matched {
			continue
		}
		snippet := buildSnippet(schedule, command, query)
		results = append(results, HighlightResult{
			Line:     e.Line,
			Schedule: schedule,
			Command:  command,
			Matched:  true,
			Snippet:  snippet,
		})
	}
	return results
}

// buildSnippet constructs a human-readable snippet marking where the query
// was found within the schedule or command fields.
func buildSnippet(schedule, command, query string) string {
	lq := strings.ToLower(query)
	if idx := strings.Index(strings.ToLower(schedule), lq); idx >= 0 {
		return fmt.Sprintf("schedule: ...%s...", markMatch(schedule, idx, len(query)))
	}
	if idx := strings.Index(strings.ToLower(command), lq); idx >= 0 {
		return fmt.Sprintf("command: ...%s...", markMatch(command, idx, len(query)))
	}
	return ""
}

// markMatch wraps the matched portion of s with bracket markers.
func markMatch(s string, start, length int) string {
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	return s[:start] + "[" + s[start:end] + "]" + s[end:]
}
