package reporter

import (
	"fmt"
	"strings"
)

// Annotation holds a label and optional note attached to a report entry.
type Annotation struct {
	Label string
	Note  string
}

// String returns a formatted representation of the annotation.
func (a Annotation) String() string {
	if a.Note == "" {
		return fmt.Sprintf("[%s]", a.Label)
	}
	return fmt.Sprintf("[%s: %s]", a.Label, a.Note)
}

// AnnotateEntries returns a copy of the report with annotations applied to
// entries whose raw schedule matches one of the provided tags map keys.
// The map key is the raw schedule string; the value is the annotation label.
func AnnotateEntries(r Report, tags map[string]string) Report {
	copy := Report{
		Source:  r.Source,
		Entries: make([]Entry, len(r.Entries)),
	}
	for i, e := range r.Entries {
		copy.Entries[i] = e
		if label, ok := tags[strings.TrimSpace(e.Raw)]; ok {
			anno := Annotation{Label: label}
			copy.Entries[i].Summary = strings.TrimSpace(
				fmt.Sprintf("%s %s", anno.String(), e.Summary),
			)
		}
	}
	return copy
}

// AnnotateErrors returns a copy of the report where every invalid entry has
// its summary prefixed with an [ERROR] annotation for quick visual scanning.
func AnnotateErrors(r Report) Report {
	copy := Report{
		Source:  r.Source,
		Entries: make([]Entry, len(r.Entries)),
	}
	for i, e := range r.Entries {
		copy.Entries[i] = e
		if !e.Valid {
			anno := Annotation{Label: "ERROR", Note: e.Error}
			copy.Entries[i].Summary = anno.String()
		}
	}
	return copy
}
