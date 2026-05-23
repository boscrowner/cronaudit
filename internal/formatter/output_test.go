package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cronaudit/internal/formatter"
	"github.com/cronaudit/internal/parser"
)

func validEntry(raw string) parser.Entry {
	return parser.Entry{Raw: raw, Valid: true}
}

func invalidEntry(raw, errMsg string) parser.Entry {
	return parser.Entry{Raw: raw, Valid: false, Error: errMsg}
}

func TestRender_TextFormat(t *testing.T) {
	entries := []parser.Entry{
		validEntry("0 * * * * /usr/bin/backup"),
		invalidEntry("99 * * * * /bin/sh", "minute out of range"),
	}
	var buf bytes.Buffer
	if err := formatter.Render(&buf, entries, formatter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Error("expected [OK] tag in text output")
	}
	if !strings.Contains(out, "[ERR]") {
		t.Error("expected [ERR] tag in text output")
	}
	if !strings.Contains(out, "minute out of range") {
		t.Error("expected error message in text output")
	}
}

func TestRender_MarkdownFormat(t *testing.T) {
	entries := []parser.Entry{
		validEntry("*/5 * * * * /bin/check"),
	}
	var buf bytes.Buffer
	if err := formatter.Render(&buf, entries, formatter.FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Status |") {
		t.Error("expected markdown table header")
	}
	if !strings.Contains(out, "✅") {
		t.Error("expected checkmark for valid entry")
	}
}

func TestRender_JSONFormat(t *testing.T) {
	entries := []parser.Entry{
		validEntry("0 0 * * * /bin/daily"),
		invalidEntry("bad entry", "invalid format"),
	}
	var buf bytes.Buffer
	if err := formatter.Render(&buf, entries, formatter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"total": 2`) {
		t.Error("expected total count in JSON output")
	}
	if !strings.Contains(out, `"valid": 1`) {
		t.Error("expected valid count in JSON output")
	}
	if !strings.Contains(out, `"invalid": 1`) {
		t.Error("expected invalid count in JSON output")
	}
}

func TestRender_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Render(&buf, []parser.Entry{}, formatter.FormatText); err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected empty output for no entries")
	}
}
