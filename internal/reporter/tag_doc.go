// Package reporter provides utilities for analysing, transforming, and
// reporting on parsed crontab entries.
//
// # Tag
//
// TagEntries applies a slice of Tag rules to every entry in a Report and
// returns a TagReport whose Results carry the matched tag names alongside
// the original entry.
//
// Built-in rules are available via CommonTags, which recognises the
// following patterns:
//
//   - frequent  – runs every minute or every few minutes
//   - daily     – fixed minute/hour, runs every day
//   - weekly    – restricted to specific weekdays
//   - monthly   – restricted to specific days of the month
//
// FilterByTag narrows a TagReport to only those results carrying a
// particular tag name (case-insensitive).
package reporter
