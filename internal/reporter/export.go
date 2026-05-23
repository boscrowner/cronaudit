package reporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// ExportCSV writes the report entries as CSV rows to the given writer.
// Columns: Line, Schedule, Command, Valid, Summary, Error
func ExportCSV(r Report, w io.Writer) error {
	cw := csv.NewWriter(w)

	header := []string{"Line", "Schedule", "Command", "Valid", "Summary", "Error"}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("export csv: write header: %w", err)
	}

	for _, e := range r.Entries {
		errMsg := ""
		if e.Err != nil {
			errMsg = e.Err.Error()
		}
		row := []string{
			strconv.Itoa(e.Line),
			e.Schedule,
			e.Command,
			strconv.FormatBool(e.Valid),
			e.Summary,
			errMsg,
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("export csv: write row %d: %w", e.Line, err)
		}
	}

	cw.Flush()
	return cw.Error()
}
