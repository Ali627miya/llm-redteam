package detectors

import (
	"strings"

	"github.com/llm-redteam/redteam/pkg/config"
	"github.com/llm-redteam/redteam/pkg/models"
)

// Analyze runs enabled detectors on the model response for one attack.
func Analyze(cfg *config.Config, attack models.AttackCase, response string) []models.Finding {
	var out []models.Finding
	lower := strings.ToLower(response)

	if cfg.Detectors.IsPII() {
		out = append(out, detectPII(response)...)
	}
	if cfg.Detectors.IsToxicity() {
		out = append(out, detectToxicity(lower, response)...)
	}
	if cfg.Detectors.IsDisallowedTopics() {
		out = append(out, detectDisallowed(lower, response)...)
	}
	if cfg.Detectors.IsLeakage() && cfg.Context != "" {
		out = append(out, detectContextLeakage(cfg.Context, response)...)
	}
	if cfg.Detectors.IsJailbreakEcho() {
		out = append(out, jailbreakHeuristic(attack, lower, response)...)
	}
	return out
}
