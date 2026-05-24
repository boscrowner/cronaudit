package reporter

import (
	"strings"
	"testing"
)

func buildAnomalyReport() Report {
	return Report{
		Source: "test",
		Entries: []Entry{
			{Line: 1, Raw: "* * * * * /bin/every-minute", Valid: true},
			{Line: 2, Raw: "0 3 * * * /bin/odd-hour", Valid: true},
			{Line: 3, Raw: "30 12 * * * /bin/normal", Valid: true},
			{Line: 4, Raw: "0 4 1 * * /bin/early", Valid: true},
			{Line: 5, Raw: "bad entry", Valid: false, Error: "invalid"},
		},
	}
}

func TestDetectAnomalies_HighFrequency(t *testing.T) {
	ar := DetectAnomalies(buildAnomalyReport())
	var found bool
	for _, a := range ar.Anomalies {
		if a.Kind == "high-frequency" && a.Line == 1 {
			found = true
		}
	}
	if !found {
		t.Error("expected high-frequency anomaly for line 1")
	}
}

func TestDetectAnomalies_OddHour(t *testing.T) {
	ar := DetectAnomalies(buildAnomalyReport())
	var count int
	for _, a := range ar.Anomalies {
		if a.Kind == "odd-hour" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 odd-hour anomalies, got %d", count)
	}
}

func TestDetectAnomalies_SkipsInvalidEntries(t *testing.T) {
	ar := DetectAnomalies(buildAnomalyReport())
	for _, a := range ar.Anomalies {
		if a.Line == 5 {
			t.Error("invalid entry should not produce anomaly")
		}
	}
}

func TestDetectAnomalies_NormalEntryClean(t *testing.T) {
	ar := DetectAnomalies(buildAnomalyReport())
	for _, a := range ar.Anomalies {
		if a.Line == 3 {
			t.Errorf("normal entry at line 3 should not be flagged, got kind=%s", a.Kind)
		}
	}
}

func TestDetectAnomalies_SourcePreserved(t *testing.T) {
	ar := DetectAnomalies(buildAnomalyReport())
	if ar.Source != "test" {
		t.Errorf("expected source 'test', got %q", ar.Source)
	}
}

func TestFormatAnomalyReport_NoAnomalies(t *testing.T) {
	ar := AnomalyReport{Source: "clean", Anomalies: nil}
	out := FormatAnomalyReport(ar)
	if !strings.Contains(out, "No anomalies") {
		t.Errorf("expected no-anomaly message, got: %s", out)
	}
}

func TestFormatAnomalyReport_ContainsKind(t *testing.T) {
	ar := DetectAnomalies(buildAnomalyReport())
	out := FormatAnomalyReport(ar)
	if !strings.Contains(out, "high-frequency") {
		t.Errorf("expected 'high-frequency' in output, got: %s", out)
	}
}
