// Package reporter provides utilities for analysing and transforming cron
// report data produced by the cronaudit parser.
//
// # Baseline
//
// BuildBaseline compares a live [Report] against a pre-approved list of
// schedule+command strings and categorises each entry as matched, unknown,
// or missing.
//
// Example usage:
//
//	approved := []string{
//		"0 9 * * 1-5 /usr/bin/backup",
//		"*/5 * * * * /usr/bin/health",
//	}
//	b := reporter.BuildBaseline(report, approved)
//	if !b.IsClean() {
//		fmt.Printf("unknown: %v\n", b.Unknown)
//		fmt.Printf("missing: %v\n", b.Missing)
//	}
package reporter
