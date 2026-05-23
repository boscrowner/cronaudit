// Package reporter provides tools for analyzing and transforming parsed crontab reports.
//
// # Cluster
//
// ClusterByHour and ClusterByDayOfWeek group valid crontab entries by shared
// schedule components, making it easy to spot scheduling patterns and potential
// conflicts in a crontab file.
//
// Example usage:
//
//	clusters := reporter.ClusterByHour(report)
//	for _, c := range clusters {
//		fmt.Printf("Pattern %q: %d entries\n", c.Pattern, c.Size)
//	}
package reporter
