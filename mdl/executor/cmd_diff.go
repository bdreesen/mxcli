// SPDX-License-Identifier: Apache-2.0

// Package executor - Diff command implementation for comparing MDL scripts against project state
package executor

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/ast"
)

// DiffFormat represents the output format for diff results
type DiffFormat string

const (
	DiffFormatUnified    DiffFormat = "unified"
	DiffFormatSideBySide DiffFormat = "side"
	DiffFormatStructural DiffFormat = "struct"
)

// DiffOptions configures diff output
type DiffOptions struct {
	Format   DiffFormat
	UseColor bool
	Width    int
}

// ChangeType represents the type of structural change
type ChangeType string

const (
	ChangeAdded    ChangeType = "+"
	ChangeRemoved  ChangeType = "-"
	ChangeModified ChangeType = "~"
)

// StructuralChange represents a single structural change within an object
type StructuralChange struct {
	ChangeType  ChangeType
	ElementType string // "Attribute", "Parameter", "Value", etc.
	ElementName string
	Details     string
}

// DiffResult represents the diff for a single object
type DiffResult struct {
	ObjectType string
	ObjectName ast.QualifiedName
	Current    string // MDL from MPR (empty if new)
	Proposed   string // MDL from script
	IsNew      bool
	IsDeleted  bool
	Changes    []StructuralChange
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
)

// DiffProgram compares an MDL program against the current project state
func (e *Executor) DiffProgram(prog *ast.Program, opts DiffOptions) error {
	if e.reader == nil {
		return fmt.Errorf("not connected to a project")
	}

	// Set defaults
	if opts.Format == "" {
		opts.Format = DiffFormatUnified
	}
	if opts.Width == 0 {
		opts.Width = 120
	}

	var results []DiffResult
	var newCount, modifiedCount, unchangedCount int

	// Track processed objects to avoid duplicates (script may have multiple statements for same object)
	processed := make(map[string]bool)

	// Process each statement
	for _, stmt := range prog.Statements {
		result, err := e.diffStatement(stmt)
		if err != nil {
			// Skip statements that can't be diffed (e.g., connection statements)
			continue
		}
		if result != nil {
			// Create unique key for deduplication
			key := result.ObjectType + ":" + result.ObjectName.String()
			if processed[key] {
				// Skip duplicate - already processed this object
				continue
			}
			processed[key] = true

			results = append(results, *result)
			if result.IsNew {
				newCount++
			} else if result.Current != result.Proposed {
				modifiedCount++
			} else {
				unchangedCount++
			}
		}
	}

	// Output results based on format
	for _, result := range results {
		if result.Current == result.Proposed && !result.IsNew {
			// Skip unchanged objects unless showing structural
			if opts.Format != DiffFormatStructural {
				continue
			}
		}

		switch opts.Format {
		case DiffFormatUnified:
			e.outputUnifiedDiff(result, opts.UseColor)
		case DiffFormatSideBySide:
			e.outputSideBySideDiff(result, opts.Width, opts.UseColor)
		case DiffFormatStructural:
			e.outputStructuralDiff(result, opts.UseColor)
		}
	}

	// Output summary
	fmt.Fprintf(e.output, "\nSummary: %d new, %d modified, %d unchanged\n",
		newCount, modifiedCount, unchangedCount)

	return nil
}

// diffStatement generates a diff result for a single statement
func (e *Executor) diffStatement(stmt ast.Statement) (*DiffResult, error) {
	switch s := stmt.(type) {
	case *ast.CreateEntityStmt:
		return e.diffEntity(s)
	case *ast.CreateViewEntityStmt:
		return e.diffViewEntity(s)
	case *ast.CreateEnumerationStmt:
		return e.diffEnumeration(s)
	case *ast.CreateAssociationStmt:
		return e.diffAssociation(s)
	case *ast.CreateMicroflowStmt:
		return e.diffMicroflow(s)
	default:
		return nil, nil // Skip unsupported statements
	}
}

// diffEntity compares a CREATE ENTITY statement against the project
func (e *Executor) diffEntity(s *ast.CreateEntityStmt) (*DiffResult, error) {
	result := &DiffResult{
		ObjectType: "Entity",
		ObjectName: s.Name,
		Proposed:   e.entityStmtToMDL(s),
	}

	// Try to find existing entity
	module, err := e.findModule(s.Name.Module)
	if err != nil {
		result.IsNew = true
		return result, nil
	}

	dm, err := e.reader.GetDomainModel(module.ID)
	if err != nil {
		result.IsNew = true
		return result, nil
	}

	for _, entity := range dm.Entities {
		if entity.Name == s.Name.Name {
			// Found existing entity - get its MDL representation
			result.Current = e.entityToMDL(module.Name, entity, dm)
			result.Changes = e.compareEntities(result.Current, result.Proposed)
			return result, nil
		}
	}

	result.IsNew = true
	return result, nil
}

// diffViewEntity compares a CREATE VIEW ENTITY statement against the project
func (e *Executor) diffViewEntity(s *ast.CreateViewEntityStmt) (*DiffResult, error) {
	result := &DiffResult{
		ObjectType: "View Entity",
		ObjectName: s.Name,
		Proposed:   e.viewEntityStmtToMDL(s),
	}

	module, err := e.findModule(s.Name.Module)
	if err != nil {
		result.IsNew = true
		return result, nil
	}

	dm, err := e.reader.GetDomainModel(module.ID)
	if err != nil {
		result.IsNew = true
		return result, nil
	}

	for _, entity := range dm.Entities {
		if entity.Name == s.Name.Name {
			result.Current = e.viewEntityFromProjectToMDL(module.Name, entity, dm)
			return result, nil
		}
	}

	result.IsNew = true
	return result, nil
}

// diffEnumeration compares a CREATE ENUMERATION statement against the project
func (e *Executor) diffEnumeration(s *ast.CreateEnumerationStmt) (*DiffResult, error) {
	result := &DiffResult{
		ObjectType: "Enumeration",
		ObjectName: s.Name,
		Proposed:   e.enumerationStmtToMDL(s),
	}

	// Try to find existing enumeration
	existingEnum := e.findEnumeration(s.Name.Module, s.Name.Name)
	if existingEnum == nil {
		result.IsNew = true
		return result, nil
	}

	h, _ := e.getHierarchy()
	modName := h.GetModuleName(existingEnum.ContainerID)
	result.Current = e.enumerationToMDL(modName, existingEnum)
	result.Changes = e.compareEnumerations(result.Current, result.Proposed)

	return result, nil
}

// diffAssociation compares a CREATE ASSOCIATION statement against the project
func (e *Executor) diffAssociation(s *ast.CreateAssociationStmt) (*DiffResult, error) {
	result := &DiffResult{
		ObjectType: "Association",
		ObjectName: s.Name,
		Proposed:   e.associationStmtToMDL(s),
	}

	module, err := e.findModule(s.Name.Module)
	if err != nil {
		result.IsNew = true
		return result, nil
	}

	dm, err := e.reader.GetDomainModel(module.ID)
	if err != nil {
		result.IsNew = true
		return result, nil
	}

	for _, assoc := range dm.Associations {
		if assoc.Name == s.Name.Name {
			result.Current = e.associationToMDL(module.Name, assoc, dm)
			return result, nil
		}
	}

	result.IsNew = true
	return result, nil
}

// diffMicroflow compares a CREATE MICROFLOW statement against the project
func (e *Executor) diffMicroflow(s *ast.CreateMicroflowStmt) (*DiffResult, error) {
	result := &DiffResult{
		ObjectType: "Microflow",
		ObjectName: s.Name,
		Proposed:   e.microflowStmtToMDL(s),
	}

	// Try to find existing microflow
	h, err := e.getHierarchy()
	if err != nil {
		result.IsNew = true
		return result, nil
	}

	mfs, err := e.reader.ListMicroflows()
	if err != nil {
		result.IsNew = true
		return result, nil
	}

	for _, mf := range mfs {
		modID := h.FindModuleID(mf.ContainerID)
		modName := h.GetModuleName(modID)
		if modName == s.Name.Module && mf.Name == s.Name.Name {
			// Capture current MDL representation
			var buf bytes.Buffer
			oldOutput := e.output
			e.output = &buf
			e.describeMicroflow(s.Name)
			e.output = oldOutput
			result.Current = strings.TrimSuffix(buf.String(), "\n")
			result.Changes = e.compareMicroflows(result.Current, result.Proposed)
			return result, nil
		}
	}

	result.IsNew = true
	return result, nil
}

// ============================================================================
// Structural Comparison Functions
// ============================================================================

// compareEntities extracts structural changes between two entity MDL representations
func (e *Executor) compareEntities(current, proposed string) []StructuralChange {
	var changes []StructuralChange

	// Simple line-based comparison for now
	currentLines := strings.Split(current, "\n")
	proposedLines := strings.Split(proposed, "\n")

	// Extract attributes from both
	currentAttrs := e.extractAttributes(currentLines)
	proposedAttrs := e.extractAttributes(proposedLines)

	// Find added attributes
	for name, proposed := range proposedAttrs {
		if _, exists := currentAttrs[name]; !exists {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeAdded,
				ElementType: "Attribute",
				ElementName: name,
				Details:     proposed,
			})
		}
	}

	// Find removed attributes
	for name := range currentAttrs {
		if _, exists := proposedAttrs[name]; !exists {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeRemoved,
				ElementType: "Attribute",
				ElementName: name,
			})
		}
	}

	// Find modified attributes
	for name, proposed := range proposedAttrs {
		if current, exists := currentAttrs[name]; exists && current != proposed {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeModified,
				ElementType: "Attribute",
				ElementName: name,
				Details:     "changed",
			})
		}
	}

	return changes
}

// compareEnumerations extracts structural changes between two enumeration MDL representations
func (e *Executor) compareEnumerations(current, proposed string) []StructuralChange {
	var changes []StructuralChange

	currentValues := e.extractEnumValues(strings.Split(current, "\n"))
	proposedValues := e.extractEnumValues(strings.Split(proposed, "\n"))

	for name := range proposedValues {
		if _, exists := currentValues[name]; !exists {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeAdded,
				ElementType: "Value",
				ElementName: name,
			})
		}
	}

	for name := range currentValues {
		if _, exists := proposedValues[name]; !exists {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeRemoved,
				ElementType: "Value",
				ElementName: name,
			})
		}
	}

	return changes
}

// compareMicroflows extracts structural changes between two microflow MDL representations
func (e *Executor) compareMicroflows(current, proposed string) []StructuralChange {
	var changes []StructuralChange

	currentParams := e.extractParameters(strings.Split(current, "\n"))
	proposedParams := e.extractParameters(strings.Split(proposed, "\n"))

	for name := range proposedParams {
		if _, exists := currentParams[name]; !exists {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeAdded,
				ElementType: "Parameter",
				ElementName: name,
			})
		}
	}

	for name := range currentParams {
		if _, exists := proposedParams[name]; !exists {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeRemoved,
				ElementType: "Parameter",
				ElementName: name,
			})
		}
	}

	// Count body statements
	currentStmts := e.countBodyStatements(current)
	proposedStmts := e.countBodyStatements(proposed)
	if currentStmts != proposedStmts {
		diff := proposedStmts - currentStmts
		if diff > 0 {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeAdded,
				ElementType: "Body",
				ElementName: "statements",
				Details:     fmt.Sprintf("%d statements added", diff),
			})
		} else {
			changes = append(changes, StructuralChange{
				ChangeType:  ChangeRemoved,
				ElementType: "Body",
				ElementName: "statements",
				Details:     fmt.Sprintf("%d statements removed", -diff),
			})
		}
	}

	return changes
}

// extractAttributes extracts attribute definitions from MDL lines
func (e *Executor) extractAttributes(lines []string) map[string]string {
	attrs := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "CREATE") && !strings.HasPrefix(line, "/**") && !strings.HasPrefix(line, "*") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				if !strings.HasPrefix(name, "$") && !strings.HasPrefix(name, "@") {
					attrs[name] = strings.TrimSuffix(strings.TrimSpace(parts[1]), ",")
				}
			}
		}
	}
	return attrs
}

// extractEnumValues extracts enumeration values from MDL lines
func (e *Executor) extractEnumValues(lines []string) map[string]bool {
	values := make(map[string]bool)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "'") && !strings.HasPrefix(line, "CREATE") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				name := strings.TrimSuffix(parts[0], ",")
				if name != "" && !strings.HasPrefix(name, "/") && !strings.HasPrefix(name, "*") {
					values[name] = true
				}
			}
		}
	}
	return values
}

// extractParameters extracts parameter names from MDL lines
func (e *Executor) extractParameters(lines []string) map[string]bool {
	params := make(map[string]bool)
	inParams := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CREATE MICROFLOW") || strings.HasPrefix(line, "CREATE NANOFLOW") {
			inParams = true
			continue
		}
		if inParams {
			if strings.HasPrefix(line, ")") {
				inParams = false
				continue
			}
			if strings.HasPrefix(line, "$") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) >= 1 {
					name := strings.TrimPrefix(parts[0], "$")
					name = strings.TrimSuffix(name, ",")
					params[strings.TrimSpace(name)] = true
				}
			}
		}
	}
	return params
}

// countBodyStatements counts statements in a microflow body
func (e *Executor) countBodyStatements(mdl string) int {
	count := 0
	inBody := false
	for line := range strings.SplitSeq(mdl, "\n") {
		line = strings.TrimSpace(line)
		if line == "BEGIN" {
			inBody = true
			continue
		}
		if line == "END;" {
			break
		}
		if inBody && line != "" && !strings.HasPrefix(line, "--") {
			count++
		}
	}
	return count
}
