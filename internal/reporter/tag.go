package reporter

import (
	"strings"

	"github.com/cronaudit/internal/parser"
)

// Tag represents a label applied to a cron entry based on matching rules.
type Tag struct {
	Name  string
	Match func(entry parser.Entry) bool
}

// TagResult holds an entry along with any tags applied to it.
type TagResult struct {
	Entry parser.Entry
	Tags  []string
}

// TagReport holds the source and tagged results.
type TagReport struct {
	Source  string
	Results []TagResult
}

// TagEntries applies a set of Tag rules to each valid entry in the report.
// Invalid entries are included with no tags applied.
func TagEntries(report Report, tags []Tag) TagReport {
	results := make([]TagResult, 0, len(report.Entries))
	for _, e := range report.Entries {
		var matched []string
		if e.Valid {
			for _, t := range tags {
				if t.Match(e) {
					matched = append(matched, t.Name)
				}
			}
		}
		results = append(results, TagResult{Entry: e, Tags: matched})
	}
	return TagReport{Source: report.Source, Results: results}
}

// FilterByTag returns only the TagResults that carry the given tag name.
func FilterByTag(tr TagReport, name string) []TagResult {
	var out []TagResult
	for _, r := range tr.Results {
		for _, t := range r.Tags {
			if strings.EqualFold(t, name) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

// CommonTags provides a set of built-in Tag rules for typical cron patterns.
func CommonTags() []Tag {
	return []Tag{
		{
			Name: "frequent",
			Match: func(e parser.Entry) bool {
				// runs every minute or every 5 minutes or less
				return e.Minute == "*" || strings.HasPrefix(e.Minute, "*/1") || e.Minute == "*/2" || e.Minute == "*/5"
			},
		},
		{
			Name: "daily",
			Match: func(e parser.Entry) bool {
				return e.Minute != "*" && e.Hour != "*" && e.Day == "*" && e.Month == "*"
			},
		},
		{
			Name: "weekly",
			Match: func(e parser.Entry) bool {
				return e.Weekday != "*" && e.Day == "*"
			},
		},
		{
			Name: "monthly",
			Match: func(e parser.Entry) bool {
				return e.Day != "*" && e.Month == "*"
			},
		},
	}
}
