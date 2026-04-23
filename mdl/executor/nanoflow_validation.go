// SPDX-License-Identifier: Apache-2.0

// Package executor - nanoflow-specific validation rules
package executor

import (
	"fmt"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/ast"
)

// nanoflowDisallowedActions lists AST statement types that are not allowed in
// nanoflow bodies. These correspond to microflow-only actions in the Mendix
// runtime: Java actions, REST/web service calls, workflow actions, import/export,
// external object operations, download, push-to-client, show home page, and
// JSON transformation.
var nanoflowDisallowedActions = map[string]string{
	"*ast.RaiseErrorStmt":                 "ErrorEvent is not allowed in nanoflows",
	"*ast.CallJavaActionStmt":             "Java actions cannot be called from nanoflows",
	"*ast.ExecuteDatabaseQueryStmt":       "database queries are not allowed in nanoflows",
	"*ast.CallExternalActionStmt":         "external action calls are not allowed in nanoflows",
	"*ast.ShowHomePageStmt":               "SHOW HOME PAGE is not allowed in nanoflows",
	"*ast.RestCallStmt":                   "REST calls are not allowed in nanoflows",
	"*ast.SendRestRequestStmt":            "REST requests are not allowed in nanoflows",
	"*ast.ImportFromMappingStmt":          "import mapping is not allowed in nanoflows",
	"*ast.ExportToMappingStmt":            "export mapping is not allowed in nanoflows",
	"*ast.TransformJsonStmt":              "JSON transformation is not allowed in nanoflows",
	"*ast.CallWorkflowStmt":               "workflow calls are not allowed in nanoflows",
	"*ast.GetWorkflowDataStmt":            "workflow actions are not allowed in nanoflows",
	"*ast.GetWorkflowsStmt":               "workflow actions are not allowed in nanoflows",
	"*ast.GetWorkflowActivityRecordsStmt": "workflow actions are not allowed in nanoflows",
	"*ast.WorkflowOperationStmt":          "workflow actions are not allowed in nanoflows",
	"*ast.SetTaskOutcomeStmt":             "workflow actions are not allowed in nanoflows",
	"*ast.OpenUserTaskStmt":               "workflow actions are not allowed in nanoflows",
	"*ast.NotifyWorkflowStmt":             "workflow actions are not allowed in nanoflows",
	"*ast.OpenWorkflowStmt":               "workflow actions are not allowed in nanoflows",
	"*ast.LockWorkflowStmt":               "workflow actions are not allowed in nanoflows",
	"*ast.UnlockWorkflowStmt":             "workflow actions are not allowed in nanoflows",
}

// validateNanoflowBody checks that a nanoflow body does not contain disallowed
// actions or flow objects. Returns a list of human-readable error messages.
func validateNanoflowBody(body []ast.MicroflowStatement) []string {
	var errors []string
	validateNanoflowStatements(body, &errors)
	return errors
}

func validateNanoflowStatements(stmts []ast.MicroflowStatement, errors *[]string) {
	for _, stmt := range stmts {
		typeName := fmt.Sprintf("%T", stmt)
		if reason, disallowed := nanoflowDisallowedActions[typeName]; disallowed {
			*errors = append(*errors, reason)
			continue
		}
		// Recurse into compound statements
		switch s := stmt.(type) {
		case *ast.IfStmt:
			validateNanoflowStatements(s.ThenBody, errors)
			validateNanoflowStatements(s.ElseBody, errors)
		case *ast.LoopStmt:
			validateNanoflowStatements(s.Body, errors)
		case *ast.WhileStmt:
			validateNanoflowStatements(s.Body, errors)
		}
		// Also recurse into error handling bodies
		if eh := getErrorHandling(stmt); eh != nil && eh.Body != nil {
			validateNanoflowStatements(eh.Body, errors)
		}
	}
}

// getErrorHandling extracts the ErrorHandlingClause from statements that have one.
func getErrorHandling(stmt ast.MicroflowStatement) *ast.ErrorHandlingClause {
	switch s := stmt.(type) {
	case *ast.CreateObjectStmt:
		return s.ErrorHandling
	case *ast.MfCommitStmt:
		return s.ErrorHandling
	case *ast.DeleteObjectStmt:
		return s.ErrorHandling
	case *ast.RetrieveStmt:
		return s.ErrorHandling
	case *ast.CallMicroflowStmt:
		return s.ErrorHandling
	case *ast.CallNanoflowStmt:
		return s.ErrorHandling
	case *ast.CallJavaActionStmt:
		return s.ErrorHandling
	case *ast.CallExternalActionStmt:
		return s.ErrorHandling
	case *ast.RestCallStmt:
		return s.ErrorHandling
	case *ast.SendRestRequestStmt:
		return s.ErrorHandling
	case *ast.ImportFromMappingStmt:
		return s.ErrorHandling
	case *ast.ExportToMappingStmt:
		return s.ErrorHandling
	case *ast.TransformJsonStmt:
		return s.ErrorHandling
	case *ast.ExecuteDatabaseQueryStmt:
		return s.ErrorHandling
	case *ast.ListOperationStmt:
		return nil
	}
	return nil
}

// validateNanoflowReturnType checks that the return type is allowed for nanoflows.
// Binary and Float return types are not supported.
func validateNanoflowReturnType(retType *ast.MicroflowReturnType) string {
	if retType == nil {
		return ""
	}
	switch retType.Type.Kind {
	case ast.TypeBinary:
		return "Binary return type is not allowed in nanoflows"
	}
	return ""
}

// validateNanoflow runs all nanoflow-specific validations and returns a combined
// error message, or empty string if valid.
func validateNanoflow(name string, body []ast.MicroflowStatement, retType *ast.MicroflowReturnType) string {
	var allErrors []string

	if msg := validateNanoflowReturnType(retType); msg != "" {
		allErrors = append(allErrors, msg)
	}

	allErrors = append(allErrors, validateNanoflowBody(body)...)

	if len(allErrors) == 0 {
		return ""
	}

	var errMsg strings.Builder
	errMsg.WriteString(fmt.Sprintf("nanoflow '%s' has validation errors:\n", name))
	for _, e := range allErrors {
		errMsg.WriteString(fmt.Sprintf("  - %s\n", e))
	}
	return errMsg.String()
}
