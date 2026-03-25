// SPDX-License-Identifier: Apache-2.0

package ast

// ============================================================================
// REST Client Statements
// ============================================================================

// CreateRestClientStmt represents: CREATE REST CLIENT Module.Name BASE URL '...' AUTHENTICATION ... BEGIN ... END
type CreateRestClientStmt struct {
	Name           QualifiedName
	BaseUrl        string
	Authentication *RestAuthDef // nil = AUTHENTICATION NONE
	Operations     []*RestOperationDef
	Documentation  string
	Folder         string // Folder path within module
	CreateOrModify bool   // True if CREATE OR MODIFY was used
}

func (s *CreateRestClientStmt) isStatement() {}

// RestAuthDef represents authentication configuration in a CREATE REST CLIENT statement.
type RestAuthDef struct {
	Scheme   string // "BASIC"
	Username string // literal string or $variable name
	Password string // literal string or $variable name
}

// RestOperationDef represents a single operation in a CREATE REST CLIENT statement.
type RestOperationDef struct {
	Name             string
	Documentation    string
	Method           string // "GET", "POST", "PUT", "PATCH", "DELETE"
	Path             string
	Parameters       []RestParamDef  // path parameters
	QueryParameters  []RestParamDef  // query parameters
	Headers          []RestHeaderDef // HTTP headers
	BodyType         string          // "JSON", "FILE", "" (none)
	BodyVariable     string          // e.g. "$ItemData"
	ResponseType     string          // "JSON", "STRING", "FILE", "STATUS", "NONE"
	ResponseVariable string          // e.g. "$CreatedItem"
	Timeout          int
}

// RestParamDef represents a path or query parameter definition.
type RestParamDef struct {
	Name     string // includes $ prefix, e.g. "$userId"
	DataType string // "String", "Integer", "Boolean", "Decimal"
}

// RestHeaderDef represents an HTTP header definition.
type RestHeaderDef struct {
	Name     string // header name, e.g. "Accept"
	Value    string // static value, e.g. "application/json" (may be empty if Variable is set)
	Variable string // dynamic variable, e.g. "$Token" (may be empty if Value is set)
	Prefix   string // concatenation prefix, e.g. "Bearer " (used with Variable)
}

// DropRestClientStmt represents: DROP REST CLIENT Module.Name
type DropRestClientStmt struct {
	Name QualifiedName
}

func (s *DropRestClientStmt) isStatement() {}
