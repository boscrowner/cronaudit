// Package reporter provides utilities for building, filtering, sorting,
// grouping, deduplicating, exporting, and annotating crontab reports.
//
// # Annotations
//
// AnnotateEntries applies user-defined labels to entries based on their raw
// schedule string. This is useful for tagging known schedules (e.g. marking
// "0 2 * * *" as "nightly-backup") so that rendered output is self-documenting.
//
// AnnotateErrors prefixes every invalid entry's summary with an [ERROR] tag
// containing the parse error message, making problems immediately visible in
// text or markdown output without requiring a separate filter step.
package reporter
