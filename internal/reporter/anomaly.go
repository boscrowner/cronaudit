package reporter

import (
	"fmt"
	"strings"
)

// AnomalyResult holds a detected anomaly for a single cron entry.
type AnomalyResult struct {
	Line     int
	Schedule string
	Command  string
	Kind     string
	Detail   string
}

// AnomalyReport is the result of anomaly detection over a Report.
type AnomalyReport struct {
	Source    string
	Anomalies []AnomalyResult
}

// DetectAnomalies inspects entries for suspicious or unusual scheduling patterns.
// It flags entries that run too frequently, at unusual hours, or have empty commands.
func DetectAnomalies(r Report) AnomalyReport {
	var anomalies []AnomalyResult

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		fields := strings.Fields(e.Raw)
		if len(fields) < 6 {
			continue
		}
		minute, hour, command := fields[0], fields[1], strings.Join(fields[5:], " ")

		if minute == "*" && hour == "*" {
			anomalies = append(anomalies, AnomalyResult{
				Line:     e.Line,
				Schedule: fmt.Sprintf("%s %s %s %s %s", fields[0], fields[1], fields[2], fields[3], fields[4]),
				Command:  command,
				Kind:     "high-frequency",
				Detail:   "job runs every minute of every hour",
			})
		}

		if hour == "3" || hour == "4" {
			anomalies = append(anomalies, AnomalyResult{
				Line:     e.Line,
				Schedule: fmt.Sprintf("%s %s %s %s %s", fields[0], fields[1], fields[2], fields[3], fields[4]),
				Command:  command,
				Kind:     "odd-hour",
				Detail:   fmt.Sprintf("job scheduled at hour %s (early morning)", hour),
			})
		}

		if strings.TrimSpace(command) == "" {
			anomalies = append(anomalies, AnomalyResult{
				Line:     e.Line,
				Schedule: fmt.Sprintf("%s %s %s %s %s", fields[0], fields[1], fields[2], fields[3], fields[4]),
				Command:  command,
				Kind:     "empty-command",
				Detail:   "no command specified",
			})
		}
	}

	return AnomalyReport{
		Source:    r.Source,
		Anomalies: anomalies,
	}
}
