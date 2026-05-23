// Package reporter provides types and functions for building, filtering,
// sorting, and summarising crontab audit reports.
//
// The Stats type and its helpers offer a lightweight way to obtain aggregate
// counts and derived metrics from a Report without re-processing raw entries:
//
//	s := reporter.ComputeStats(report)
//	fmt.Printf("Pass rate: %.0f%%\n", s.PassRate()*100)
//	if s.HasErrors() {
//		os.Exit(1)
//	}
package reporter
