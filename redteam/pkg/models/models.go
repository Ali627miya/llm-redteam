package models

import "time"

// Severity ranks how serious a finding is for triage.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// AttackCase is one prompt from the built-in or custom library.
type AttackCase struct {
	ID          string   `yaml:"id" json:"id"`
	Category    string   `yaml:"category" json:"category"`
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Prompt      string   `yaml:"prompt" json:"prompt"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// RunResult is the outcome of sending one attack to the target.
type RunResult struct {
	AttackID    string            `json:"attack_id"`
	Category    string            `json:"category"`
	Name        string            `json:"name"`
	Prompt      string            `json:"prompt"`
	Response    string            `json:"response"`
	LatencyMS   int64             `json:"latency_ms"`
	Error       string            `json:"error,omitempty"`
	Findings    []Finding         `json:"findings"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	StartedAt   time.Time         `json:"started_at"`
	CompletedAt time.Time         `json:"completed_at"`
}

// Finding is a single security or safety signal on a response.
type Finding struct {
	Detector string   `json:"detector"`
	Rule     string   `json:"rule"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Evidence string   `json:"evidence,omitempty"`
}

// Summary aggregates a full scan.
type Summary struct {
	TotalAttacks   int `json:"total_attacks"`
	Passed         int `json:"passed"`
	Failed         int `json:"failed"`
	Errors         int `json:"errors"`
	TotalFindings  int `json:"total_findings"`
	CriticalCount  int `json:"critical_count"`
	HighCount      int `json:"high_count"`
}

// ScanReport is the top-level report object.
type ScanReport struct {
	SchemaVersion string       `json:"schema_version"`
	GeneratedAt   time.Time    `json:"generated_at"`
	Target        string       `json:"target"`
	Summary       Summary      `json:"summary"`
	Results       []RunResult  `json:"results"`
}
