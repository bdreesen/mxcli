// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"testing"

	"github.com/mendixlabs/mxcli/mdl/ast"
)

func TestValidateMicroflow_EnumSplitAllBranchesReturn(t *testing.T) {
	stmt := &ast.CreateMicroflowStmt{
		Name: ast.QualifiedName{Module: "Sample", Name: "Route"},
		ReturnType: &ast.MicroflowReturnType{
			Type: ast.DataType{Kind: ast.TypeBoolean},
		},
		Body: []ast.MicroflowStatement{
			&ast.EnumSplitStmt{
				Variable: "Status",
				Cases: []ast.EnumSplitCase{
					{Values: []string{"Open"}, Body: []ast.MicroflowStatement{
						&ast.ReturnStmt{Value: &ast.LiteralExpr{Kind: ast.LiteralBoolean, Value: true}},
					}},
					{Values: []string{"Closed"}, Body: []ast.MicroflowStatement{
						&ast.ReturnStmt{Value: &ast.LiteralExpr{Kind: ast.LiteralBoolean, Value: false}},
					}},
				},
			},
		},
	}

	violations := ValidateMicroflow(stmt)
	for _, v := range violations {
		if v.RuleID == "MDL003" {
			t.Fatalf("enum split with all cases returning must not trigger MDL003: %#v", v)
		}
	}
}

func TestValidateMicroflow_EnumSplitElseForbidden(t *testing.T) {
	stmt := &ast.CreateMicroflowStmt{
		Name: ast.QualifiedName{Module: "Sample", Name: "Route"},
		Body: []ast.MicroflowStatement{
			&ast.EnumSplitStmt{
				Variable: "Status",
				Cases: []ast.EnumSplitCase{
					{Values: []string{"Open"}, Body: []ast.MicroflowStatement{
						&ast.ReturnStmt{Value: &ast.LiteralExpr{Kind: ast.LiteralBoolean, Value: true}},
					}},
				},
				ElseBody: []ast.MicroflowStatement{
					&ast.ReturnStmt{Value: &ast.LiteralExpr{Kind: ast.LiteralBoolean, Value: false}},
				},
			},
		},
	}

	violations := ValidateMicroflow(stmt)
	for _, v := range violations {
		if v.RuleID == "MDL008" {
			return
		}
	}
	t.Fatalf("expected MDL008 for enum split with else branch, got %#v", violations)
}

func TestValidateMicroflow_EnumSplitMultipleValuesForbidden(t *testing.T) {
	stmt := &ast.CreateMicroflowStmt{
		Name: ast.QualifiedName{Module: "Sample", Name: "Route"},
		Body: []ast.MicroflowStatement{
			&ast.EnumSplitStmt{
				Variable: "Status",
				Cases: []ast.EnumSplitCase{
					{Values: []string{"Open", "Pending"}, Body: []ast.MicroflowStatement{
						&ast.ReturnStmt{Value: &ast.LiteralExpr{Kind: ast.LiteralBoolean, Value: true}},
					}},
				},
			},
		},
	}

	violations := ValidateMicroflow(stmt)
	for _, v := range violations {
		if v.RuleID == "MDL009" {
			return
		}
	}
	t.Fatalf("expected MDL009 for enum split with multiple values per branch, got %#v", violations)
}

func TestValidateMicroflow_EnumSplitBranchScopedVariable(t *testing.T) {
	stmt := &ast.CreateMicroflowStmt{
		Name: ast.QualifiedName{Module: "Sample", Name: "Route"},
		Body: []ast.MicroflowStatement{
			&ast.EnumSplitStmt{
				Variable: "Status",
				Cases: []ast.EnumSplitCase{
					{Values: []string{"Open"}, Body: []ast.MicroflowStatement{
						&ast.DeclareStmt{Variable: "OnlyInsideCase", Type: ast.DataType{Kind: ast.TypeString}},
					}},
				},
			},
			&ast.MfSetStmt{
				Target: "OnlyInsideCase",
				Value:  &ast.LiteralExpr{Kind: ast.LiteralString, Value: "outside"},
			},
		},
	}

	violations := ValidateMicroflow(stmt)
	for _, v := range violations {
		if v.RuleID == "MDL005" {
			return
		}
	}
	t.Fatalf("expected MDL005 for variable declared inside ENUM split branch, got %#v", violations)
}
