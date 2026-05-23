package reporter

import (
	"strings"

	"github.com/user/cronaudit/internal/parser"
)

// TransformFunc is a function that transforms a single cron entry.
type TransformFunc func(entry parser.Entry) parser.Entry

// TransformEntries applies the given TransformFunc to each entry in the report
// and returns a new Report with the transformed entries. The original report is
// not mutated.
func TransformEntries(r Report, fn TransformFunc) Report {
	transformed := make([]parser.Entry, len(r.Entries))
	for i, e := range r.Entries {
		transformed[i] = fn(e)
	}
	return Report{
		Source:  r.Source,
		Entries: transformed,
		Stats:   ComputeStats(transformed),
	}
}

// NormalizeCommands returns a TransformFunc that trims leading and trailing
// whitespace from each entry's Command field.
func NormalizeCommands() TransformFunc {
	return func(e parser.Entry) parser.Entry {
		e.Command = strings.TrimSpace(e.Command)
		return e
	}
}

// RedactCommands returns a TransformFunc that replaces the Command field of
// each valid entry with the provided placeholder string. Useful for producing
// sanitised output that omits sensitive command details.
func RedactCommands(placeholder string) TransformFunc {
	return func(e parser.Entry) parser.Entry {
		if e.Valid {
			e.Command = placeholder
		}
		return e
	}
}

// ChainTransforms combines multiple TransformFuncs into a single TransformFunc
// that applies them in order.
func ChainTransforms(fns ...TransformFunc) TransformFunc {
	return func(e parser.Entry) parser.Entry {
		for _, fn := range fns {
			e = fn(e)
		}
		return e
	}
}
