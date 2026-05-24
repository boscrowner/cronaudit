package reporter

import "strings"

// Classification represents a category assigned to a cron entry based on its command.
type Classification struct {
	Line     int
	Schedule string
	Command  string
	Category string
	Valid    bool
}

// ClassifyReport holds the classification results for all entries in a report.
type ClassifyReport struct {
	Source          string
	Classifications []Classification
}

// categoryRules maps keyword substrings to category labels.
var categoryRules = []struct {
	keyword  string
	category string
}{
	{"backup", "Backup"},
	{"restore", "Backup"},
	{"dump", "Backup"},
	{"log", "Logging"},
	{"rotate", "Logging"},
	{"clean", "Cleanup"},
	{"purge", "Cleanup"},
	{"prune", "Cleanup"},
	{"rm", "Cleanup"},
	{"deploy", "Deployment"},
	{"release", "Deployment"},
	{"sync", "Sync"},
	{"rsync", "Sync"},
	{"fetch", "Sync"},
	{"pull", "Sync"},
	{"report", "Reporting"},
	{"notify", "Notification"},
	{"alert", "Notification"},
	{"mail", "Notification"},
	{"monitor", "Monitoring"},
	{"check", "Monitoring"},
	{"health", "Monitoring"},
}

// ClassifyEntries assigns a category to each valid entry based on its command.
// Invalid entries are included with Category set to "Invalid".
func ClassifyEntries(r Report) ClassifyReport {
	result := ClassifyReport{Source: r.Source}
	for _, e := range r.Entries {
		c := Classification{
			Line:     e.Line,
			Schedule: e.Schedule,
			Command:  e.Command,
			Valid:    e.Valid,
		}
		if !e.Valid {
			c.Category = "Invalid"
		} else {
			c.Category = classifyCommand(e.Command)
		}
		result.Classifications = append(result.Classifications, c)
	}
	return result
}

// classifyCommand returns a category label for the given command string.
func classifyCommand(cmd string) string {
	lower := strings.ToLower(cmd)
	for _, rule := range categoryRules {
		if strings.Contains(lower, rule.keyword) {
			return rule.category
		}
	}
	return "General"
}

// GroupByCategory returns a map of category name to slice of Classifications.
func GroupByCategory(cr ClassifyReport) map[string][]Classification {
	groups := make(map[string][]Classification)
	for _, c := range cr.Classifications {
		groups[c.Category] = append(groups[c.Category], c)
	}
	return groups
}
