// SPDX-License-Identifier: Apache-2.0

// Package executor - Microflow flow graph builder: core types and helpers
package executor

import (
	"fmt"

	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/microflows"
	"github.com/mendixlabs/mxcli/sdk/mpr"
)

// flowBuilder helps construct the flow graph from AST statements.
type flowBuilder struct {
	objects             []microflows.MicroflowObject
	flows               []*microflows.SequenceFlow
	annotationFlows     []*microflows.AnnotationFlow
	posX                int
	posY                int
	baseY               int // Base Y position (for returning after ELSE branches)
	spacing             int
	returnValue         string                   // Return value expression for RETURN statement (used by buildFlowGraph final EndEvent)
	endsWithReturn      bool                     // True if the flow already ends with EndEvent(s) from RETURN statements
	varTypes            map[string]string        // Variable name -> entity qualified name (for CHANGE statements)
	declaredVars        map[string]string        // Declared primitive variables: name -> type (e.g., "$IsValid" -> "Boolean")
	errors              []string                 // Validation errors collected during build
	measurer            *layoutMeasurer          // For measuring statement dimensions
	nextConnectionPoint model.ID                 // For compound statements: the exit point differs from entry point
	reader              *mpr.Reader              // For looking up page/microflow references
	hierarchy           *ContainerHierarchy      // For resolving container IDs to module names
	pendingAnnotations  *ast.ActivityAnnotations // Pending annotations to attach to next activity
}

// addError records a validation error during flow building.
func (fb *flowBuilder) addError(format string, args ...any) {
	fb.errors = append(fb.errors, fmt.Sprintf(format, args...))
}

// addErrorWithExample records a validation error with a code example showing the fix.
func (fb *flowBuilder) addErrorWithExample(message, example string) {
	fb.errors = append(fb.errors, fmt.Sprintf("%s\n\n  Example:\n%s", message, example))
}

// GetErrors returns all validation errors collected during build.
func (fb *flowBuilder) GetErrors() []string {
	return fb.errors
}

// errorExampleDeclareVariable returns an example for declaring a variable.
func errorExampleDeclareVariable(varName string) string {
	// Remove $ prefix if present for cleaner display
	cleanName := varName
	if len(varName) > 0 && varName[0] == '$' {
		cleanName = varName[1:]
	}
	return fmt.Sprintf(`    DECLARE $%s Boolean = true;  -- or String, Integer, Decimal, DateTime
    ...
    SET $%s = false;`, cleanName, cleanName)
}

// isVariableDeclared checks if a variable has been declared (either as primitive or entity).
func (fb *flowBuilder) isVariableDeclared(varName string) bool {
	// Check entity variables (from parameters with entity types)
	if _, ok := fb.varTypes[varName]; ok {
		return true
	}
	// Check primitive variables (from DECLARE statements or primitive parameters)
	if _, ok := fb.declaredVars[varName]; ok {
		return true
	}
	return false
}
