// SPDX-License-Identifier: Apache-2.0

package ast

// AlterSettingsStmt represents ALTER SETTINGS commands.
type AlterSettingsStmt struct {
	Section    string         // "MODEL", "CONFIGURATION", "CONSTANT", "LANGUAGE", "WORKFLOWS"
	ConfigName string         // For CONFIGURATION section: the configuration name (e.g., "Default")
	Properties map[string]any // Key-value pairs to set
	// For CONSTANT section:
	ConstantId string // Qualified constant name
	Value      string // Constant value
}

func (s *AlterSettingsStmt) isStatement() {}
