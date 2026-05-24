// Package reporter provides tools for analyzing and reporting on parsed crontab entries.
//
// # Risk Assessment
//
// AssessRisk evaluates each valid cron entry for potential security or reliability
// risks. Each entry is assigned a RiskLevel (low, medium, or high) based on
// heuristics such as:
//
//   - Use of destructive commands (rm, dd, mkfs, etc.)
//   - Elevated privilege execution (sudo, su)
//   - High-frequency scheduling (more often than every 5 minutes)
//   - Writes to sensitive system paths (/etc/, /boot/, etc.)
//
// Example usage:
//
//	rr := reporter.AssessRisk(report)
//	fmt.Print(reporter.FormatRiskReport(rr))
package reporter
