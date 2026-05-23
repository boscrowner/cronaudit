// Package reporter provides tools for building, analysing, and reporting
// on parsed crontab entries.
//
// # Lint
//
// LintReport inspects a [Report] for suspicious or potentially problematic
// cron scheduling patterns and returns a [LintResult] containing zero or
// more [LintWarning] values.
//
// Each warning carries a [LintSeverity] (info, warn, or error) so callers
// can filter or format output accordingly.
//
// Example usage:
//
//	result := reporter.LintReport(report)
//	for _, w := range result.Warnings {
//		fmt.Printf("[%s] line %d: %s\n", w.Severity, w.Line, w.Message)
//	}
package reporter
