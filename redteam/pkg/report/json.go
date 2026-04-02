package report

import (
	"encoding/json"
	"os"
	"time"

	"github.com/llm-redteam/redteam/pkg/models"
)

const schemaVersion = "1.0"

// WriteJSON writes a scan report to path with stable schema_version.
func WriteJSON(path string, target string, results []models.RunResult) error {
	rep := models.ScanReport{
		SchemaVersion: schemaVersion,
		GeneratedAt:   time.Now().UTC(),
		Target:        target,
		Summary:       BuildSummary(results),
		Results:       results,
	}
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
