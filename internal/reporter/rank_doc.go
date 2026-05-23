// Package reporter provides tools for analysing, transforming, and
// reporting on parsed crontab entries.
//
// # Rank
//
// RankReport scores every entry in a Report using a lightweight
// heuristic and returns them sorted from highest to lowest score.
//
// Scoring rules:
//   - Invalid entries always receive a score of 0.
//   - Valid entries start at 100 and lose 15 points per wildcard "*"
//     field in the schedule (minimum score: 5).
//
// Use TopN to retrieve only the n best-ranked entries:
//
//	ranked := reporter.TopN(report, 5)
//	for _, re := range ranked {
//		fmt.Printf("%.0f  %s  (%s)\n", re.Score, re.Entry.Schedule, re.Reason)
//	}
package reporter
