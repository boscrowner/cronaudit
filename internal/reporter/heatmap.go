package reporter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// HeatmapCell represents a single cell in the heatmap grid.
type HeatmapCell struct {
	Hour  int
	Day   int // 0=Sunday .. 6=Saturday
	Count int
}

// HeatmapResult holds the full heatmap grid and a formatted ASCII representation.
type HeatmapResult struct {
	Cells     []HeatmapCell
	Formatted string
}

var dayLabels = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// BuildHeatmap produces a 7×24 frequency grid showing how many valid cron
// entries are scheduled to fire in each (day-of-week, hour) bucket.
func BuildHeatmap(r Report) HeatmapResult {
	counts := make(map[[2]int]int) // [day][hour] -> count

	for _, entry := range r.Entries {
		if !entry.Valid {
			continue
		}
		parts := strings.Fields(entry.Raw)
		if len(parts) < 6 {
			continue
		}
		hourField := parts[1]
		dowField := parts[4]

		hours := expandField(hourField, 0, 23)
		days := expandField(dowField, 0, 6)

		for _, d := range days {
			for _, h := range hours {
				counts[[2]int{d, h}]++
			}
		}
	}

	var cells []HeatmapCell
	for k, v := range counts {
		cells = append(cells, HeatmapCell{Day: k[0], Hour: k[1], Count: v})
	}
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].Day != cells[j].Day {
			return cells[i].Day < cells[j].Day
		}
		return cells[i].Hour < cells[j].Hour
	})

	return HeatmapResult{
		Cells:     cells,
		Formatted: renderHeatmap(counts),
	}
}

func renderHeatmap(counts map[[2]int]int) string {
	var sb strings.Builder
	// Header row: hours 0-23
	sb.WriteString("    ")
	for h := 0; h < 24; h++ {
		sb.WriteString(fmt.Sprintf("%2d ", h))
	}
	sb.WriteString("\n")

	for d := 0; d < 7; d++ {
		sb.WriteString(dayLabels[d] + " ")
		for h := 0; h < 24; h++ {
			c := counts[[2]int{d, h}]
			if c == 0 {
				sb.WriteString(" . ")
			} else {
				sb.WriteString(fmt.Sprintf("%2s ", strconv.Itoa(c)))
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// expandField returns the list of integers matched by a cron field token
// (supports "*", single values, ranges "a-b", and steps "*/n").
func expandField(field string, min, max int) []int {
	if field == "*" {
		return rangeInts(min, max)
	}
	var result []int
	for _, part := range strings.Split(field, ",") {
		if strings.HasPrefix(part, "*/") {
			step, err := strconv.Atoi(part[2:])
			if err != nil || step <= 0 {
				continue
			}
			for v := min; v <= max; v += step {
				result = append(result, v)
			}
		} else if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, e1 := strconv.Atoi(bounds[0])
			hi, e2 := strconv.Atoi(bounds[1])
			if e1 != nil || e2 != nil {
				continue
			}
			result = append(result, rangeInts(lo, hi)...)
		} else {
			v, err := strconv.Atoi(part)
			if err == nil {
				result = append(result, v)
			}
		}
	}
	return result
}

func rangeInts(lo, hi int) []int {
	out := make([]int, 0, hi-lo+1)
	for i := lo; i <= hi; i++ {
		out = append(out, i)
	}
	return out
}
