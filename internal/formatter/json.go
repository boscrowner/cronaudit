package formatter

import (
	"encoding/json"
	"io"
)

type jsonEntry struct {
	Raw     string `json:"raw"`
	Valid   bool   `json:"valid"`
	Summary string `json:"summary,omitempty"`
	Error   string `json:"error,omitempty"`
}

type jsonReport struct {
	Total   int         `json:"total"`
	Valid   int         `json:"valid"`
	Invalid int         `json:"invalid"`
	Entries []jsonEntry `json:"entries"`
}

func renderJSON(w io.Writer, outputs []EntryOutput) error {
	valid := 0
	entries := make([]jsonEntry, 0, len(outputs))
	for _, o := range outputs {
		je := jsonEntry{
			Raw:   o.Raw,
			Valid: o.Valid,
		}
		if o.Valid {
			je.Summary = o.Summary
			valid++
		} else {
			je.Error = o.Error
		}
		entries = append(entries, je)
	}
	report := jsonReport{
		Total:   len(outputs),
		Valid:   valid,
		Invalid: len(outputs) - valid,
		Entries: entries,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
