package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config drives how the runner invokes the application under test.
type Config struct {
	Version int `yaml:"version"`

	Target TargetConfig `yaml:"target"`

	Context string `yaml:"context,omitempty"`

	Run RunConfig `yaml:"run"`

	Categories []string `yaml:"categories,omitempty"`

	Detectors DetectorConfig `yaml:"detectors"`
}

type TargetConfig struct {
	URL string `yaml:"url"`

	Method string `yaml:"method"`

	Headers map[string]string `yaml:"headers,omitempty"`

	// BodyTemplate is JSON with Go template; use {{toJSON .Prompt}} and {{toJSON .Context}}.
	BodyTemplate string `yaml:"body_template"`

	ResponsePath string `yaml:"response_path"`
}

type RunConfig struct {
	TimeoutSeconds int `yaml:"timeout_seconds"`
	Concurrency    int `yaml:"concurrency"`
	MaxAttacks     int `yaml:"max_attacks"`
}

// DetectorConfig uses pointers so omitted YAML keys default to "on".
type DetectorConfig struct {
	PII           *bool `yaml:"pii,omitempty"`
	Toxicity      *bool `yaml:"toxicity,omitempty"`
	Disallowed    *bool `yaml:"disallowed_topics,omitempty"`
	Leakage       *bool `yaml:"context_leakage,omitempty"`
	JailbreakEcho *bool `yaml:"jailbreak_success_heuristic,omitempty"`
}

func (d DetectorConfig) IsPII() bool             { return ptrBool(d.PII, true) }
func (d DetectorConfig) IsToxicity() bool        { return ptrBool(d.Toxicity, true) }
func (d DetectorConfig) IsDisallowedTopics() bool { return ptrBool(d.Disallowed, true) }
func (d DetectorConfig) IsLeakage() bool       { return ptrBool(d.Leakage, true) }
func (d DetectorConfig) IsJailbreakEcho() bool { return ptrBool(d.JailbreakEcho, true) }

func ptrBool(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}

// Default applies defaults to the config.
func (c *Config) Default() {
	if c.Version == 0 {
		c.Version = 1
	}
	if c.Target.Method == "" {
		c.Target.Method = "POST"
	}
	if c.Target.BodyTemplate == "" {
		c.Target.BodyTemplate = `{"model":"gpt-4o-mini","messages":[{"role":"user","content":{{toJSON .Prompt}}}]}`
	}
	if c.Target.ResponsePath == "" {
		c.Target.ResponsePath = "choices.0.message.content"
	}
	if c.Run.TimeoutSeconds == 0 {
		c.Run.TimeoutSeconds = 120
	}
	if c.Run.Concurrency == 0 {
		c.Run.Concurrency = 4
	}
}

// Load reads and validates a YAML config file.
func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	c.Default()
	c.Target.URL = os.ExpandEnv(c.Target.URL)
	c.Target.BodyTemplate = os.ExpandEnv(c.Target.BodyTemplate)
	c.Context = os.ExpandEnv(c.Context)
	for k, v := range c.Target.Headers {
		c.Target.Headers[k] = os.ExpandEnv(v)
	}
	if c.Target.URL == "" {
		return nil, fmt.Errorf("target.url is required")
	}
	return &c, nil
}
