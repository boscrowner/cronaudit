// Package formatter provides output rendering for parsed crontab entries.
//
// It supports three output formats:
//
//   - FormatText     — plain text with [OK]/[ERR] prefixes (default)
//   - FormatJSON     — structured JSON report with totals and per-entry detail
//   - FormatMarkdown — GitHub-flavoured Markdown table
//
// Usage:
//
//	entries, _ := parser.Parse(r)
//	formatter.Render(os.Stdout, entries, formatter.FormatText)
//
// The Render function delegates to the appropriate renderer based on the
// Format value and writes all output to the provided io.Writer.
package formatter
