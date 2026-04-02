package report

import (
	"github.com/Ali627miya/llm-redteam/redteam/pkg/models"
)

// BuildSummary counts outcomes for a scan report.
func BuildSummary(results []models.RunResult) models.Summary {
	var s models.Summary
	s.TotalAttacks = len(results)
	for _, r := range results {
		if r.Error != "" {
			s.Errors++
			continue
		}
		if len(r.Findings) > 0 {
			s.Failed++
			s.TotalFindings += len(r.Findings)
			for _, f := range r.Findings {
				switch f.Severity {
				case models.SeverityCritical:
					s.CriticalCount++
				case models.SeverityHigh:
					s.HighCount++
				}
			}
		} else {
			s.Passed++
		}
	}
	return s
}
