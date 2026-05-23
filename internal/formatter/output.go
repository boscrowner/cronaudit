package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/cronaudit/internal/parser"
)

// Format controls the output style for crontab entries.
type Format int

const (
	FormatText Format = iota
	FormatJSON
	FormatMarkdown
)

// EntryOutput holds a parsed entry alongside its human-readable summary.
type EntryOutput struct {
	Raw     string
	Summary string
	Valid   bool
	Error   string
}

// Render writes formatted crontab entries to the given writer.
func Render(w io.Writer, entries []parser.Entry, format Format) error {
	outputs := buildOutputs(entries)
	switch format {
	case FormatJSON:
		return renderJSON(w, outputs)
	case FormatMarkdown:
		return renderMarkdown(w, outputs)
	default:
		return renderText(w, outputs)
	}
}

func buildOutputs(entries []parser.Entry) []EntryOutput {
	outs := make([]EntryOutput, 0, len(entries))
	for _, e := range entries {
		out := EntryOutput{
			Raw:   e.Raw,
			Valid: e.Valid,
		}
		if e.Valid {
			out.Summary = parser.Summarize(e)
		} else {
			out.Error = e.Error
		}
		outs = append(outs, out)
	}
	return outs
}

func renderText(w io.Writer, outputs []EntryOutput) error {
	for _, o := range outputs {
		if o.Valid {
			_, err := fmt.Fprintf(w, "[OK]  %s\n      => %s\n", o.Raw, o.Summary)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(w, "[ERR] %s\n      => %s\n", o.Raw, o.Error)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func renderMarkdown(w io.Writer, outputs []EntryOutput) error {
	_, err := fmt.Fprintln(w, "| Status | Entry | Summary / Error |")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "|--------|-------|-----------------|")
	if err != nil {
		return err
	}
	for _, o := range outputs {
		status := "✅"
		detail := o.Summary
		if !o.Valid {
			status = "❌"
			detail = o.Error
		}
		raw := strings.ReplaceAll(o.Raw, "|", "\\|")
		_, err := fmt.Fprintf(w, "| %s | `%s` | %s |\n", status, raw, detail)
		if err != nil {
			return err
		}
	}
	return nil
}
