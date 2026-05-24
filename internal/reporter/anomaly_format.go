package reporter

import (
	"fmt"
	"strings"
)

// FormatAnomalyReport returns a human-readable summary of detected anomalies.
func FormatAnomalyReport(ar AnomalyReport) string {
	if len(ar.Anomalies) == 0 {
		return fmt.Sprintf("source: %s\nNo anomalies detected.\n", ar.Source)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "source: %s\n", ar.Source)
	fmt.Fprintf(&sb, "%d anomaly(ies) detected:\n\n", len(ar.Anomalies))

	for _, a := range ar.Anomalies {
		fmt.Fprintf(&sb, "  line %-4d [%-15s] %s\n", a.Line, a.Kind, a.Detail)
		fmt.Fprintf(&sb, "             schedule: %s\n", a.Schedule)
		fmt.Fprintf(&sb, "             command:  %s\n\n", a.Command)
	}

	return sb.String()
}
