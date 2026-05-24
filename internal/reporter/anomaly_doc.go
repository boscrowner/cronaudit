// Package reporter provides report building, filtering, and analysis
// utilities for parsed crontab entries.
//
// # Anomaly Detection
//
// DetectAnomalies scans a Report for entries that exhibit suspicious or
// unusual scheduling patterns. The following anomaly kinds are recognised:
//
//   - high-frequency: the job runs every minute of every hour (* * ...)
//   - odd-hour:       the job is scheduled during early-morning hours (3–4 AM)
//   - empty-command:  no command string is present after the schedule fields
//
// Use FormatAnomalyReport to produce a human-readable summary suitable for
// CLI output or log files.
package reporter
