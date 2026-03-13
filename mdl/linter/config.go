// SPDX-License-Identifier: Apache-2.0

package linter

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the lint configuration.
type Config struct {
	// ExcludeModules lists modules to exclude from linting
	ExcludeModules []string `yaml:"excludeModules"`

	// Rules configures individual rules
	Rules map[string]RuleConfigYAML `yaml:"rules"`
}

// RuleConfigYAML is the YAML representation of rule configuration.
type RuleConfigYAML struct {
	Enabled  *bool          `yaml:"enabled"`
	Severity string         `yaml:"severity"`
	Options  map[string]any `yaml:"options"`
}

// DefaultConfig returns a default configuration.
func DefaultConfig() *Config {
	return &Config{
		ExcludeModules: []string{},
		Rules:          make(map[string]RuleConfigYAML),
	}
}

// LoadConfig loads lint configuration from a file.
// Returns default config if file doesn't exist.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// FindConfigFile searches for lint config in standard locations.
// Searches: .claude/lint-config.yaml, lint-config.yaml, .lint-config.yaml
func FindConfigFile(projectDir string) string {
	candidates := []string{
		filepath.Join(projectDir, ".claude", "lint-config.yaml"),
		filepath.Join(projectDir, "lint-config.yaml"),
		filepath.Join(projectDir, ".lint-config.yaml"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// ApplyConfig applies configuration to a linter.
func (cfg *Config) ApplyConfig(l *Linter) {
	for ruleID, ruleCfg := range cfg.Rules {
		config := RuleConfig{
			Enabled: true,
			Options: ruleCfg.Options,
		}

		if ruleCfg.Enabled != nil {
			config.Enabled = *ruleCfg.Enabled
		}

		switch ruleCfg.Severity {
		case "error":
			config.Severity = SeverityError
		case "warning":
			config.Severity = SeverityWarning
		case "info":
			config.Severity = SeverityInfo
		case "hint":
			config.Severity = SeverityHint
		default:
			// Keep default severity from rule
			for _, rule := range l.Rules() {
				if rule.ID() == ruleID {
					config.Severity = rule.DefaultSeverity()
					break
				}
			}
		}

		l.ConfigureRule(ruleID, config)
	}
}

// ParseSeverity parses a severity string.
func ParseSeverity(s string) Severity {
	switch s {
	case "error":
		return SeverityError
	case "warning":
		return SeverityWarning
	case "info":
		return SeverityInfo
	case "hint":
		return SeverityHint
	default:
		return SeverityWarning
	}
}
