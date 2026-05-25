// Package reporter provides tools for analysing and reporting on crontab entries.
//
// # Impact Assessment
//
// AssessImpact evaluates each valid cron entry and assigns one of four impact
// levels — critical, high, medium, or low — based on heuristics:
//
//   - Destructive commands (rm, drop, truncate, …)
//   - Sensitive system paths (/etc/, /var/, …)
//   - Privilege escalation (sudo, chmod, …)
//   - Network activity (curl, wget, ssh, …)
//
// The more heuristics that fire, the higher the impact level.
// Results are sorted from highest to lowest impact.
//
// Use FormatImpactReport to render results as a human-readable string.
package reporter
