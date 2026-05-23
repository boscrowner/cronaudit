package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Snapshot captures a Report at a point in time for archival or comparison.
type Snapshot struct {
	CapturedAt time.Time `json:"captured_at"`
	Label      string    `json:"label,omitempty"`
	Report     Report    `json:"report"`
}

// TakeSnapshot wraps the given Report in a Snapshot with the current timestamp.
func TakeSnapshot(r Report, label string) Snapshot {
	return Snapshot{
		CapturedAt: time.Now().UTC(),
		Label:      label,
		Report:     r,
	}
}

// WriteSnapshot serialises the Snapshot as JSON to w.
func WriteSnapshot(w io.Writer, s Snapshot) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// ReadSnapshot deserialises a Snapshot from r.
func ReadSnapshot(r io.Reader) (Snapshot, error) {
	var s Snapshot
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return s, nil
}

// DiffSnapshots returns a DiffResult comparing the reports stored in two
// snapshots, treating a as the baseline and b as the newer version.
func DiffSnapshots(a, b Snapshot) DiffResult {
	return DiffReports(a.Report, b.Report)
}
