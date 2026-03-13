// SPDX-License-Identifier: Apache-2.0

package visitor

import (
	"strings"

	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/mdl/grammar/parser"
)

// ExitCreateJavaActionStatement handles CREATE JAVA ACTION statements.
func (b *Builder) ExitCreateJavaActionStatement(ctx *parser.CreateJavaActionStatementContext) {
	stmt := &ast.CreateJavaActionStmt{}

	// Get qualified name
	if qn := ctx.QualifiedName(); qn != nil {
		stmt.Name = buildQualifiedName(qn)
	}

	// Get parameters
	if paramList := ctx.JavaActionParameterList(); paramList != nil {
		for _, paramCtx := range paramList.AllJavaActionParameter() {
			param := ast.JavaActionParam{}
			if pn := paramCtx.ParameterName(); pn != nil {
				param.Name = parameterNameText(pn)
			}
			if dt := paramCtx.DataType(); dt != nil {
				param.Type = buildDataType(dt)
			}
			// Check for NOT NULL constraint
			if paramCtx.NOT_NULL() != nil {
				param.IsRequired = true
			}
			stmt.Parameters = append(stmt.Parameters, param)
		}
	}

	// Extract type parameters from ENTITY <pEntity> parameter declarations
	for _, param := range stmt.Parameters {
		if param.Type.Kind == ast.TypeEntityTypeParam && param.Type.TypeParamName != "" {
			found := false
			for _, existing := range stmt.TypeParameters {
				if existing == param.Type.TypeParamName {
					found = true
					break
				}
			}
			if !found {
				stmt.TypeParameters = append(stmt.TypeParameters, param.Type.TypeParamName)
			}
		}
	}

	// Get return type
	if retType := ctx.JavaActionReturnType(); retType != nil {
		if dt := retType.DataType(); dt != nil {
			stmt.ReturnType = buildDataType(dt)
		}
	}

	// Get exposed clause (EXPOSED AS 'caption' IN 'category')
	if exposed := ctx.JavaActionExposedClause(); exposed != nil {
		allStrings := exposed.AllSTRING_LITERAL()
		if len(allStrings) >= 2 {
			stmt.ExposedCaption = unquoteString(allStrings[0].GetText())
			stmt.ExposedCategory = unquoteString(allStrings[1].GetText())
		}
	}

	// Get Java code from dollar-quoted string
	if dollarStr := ctx.DOLLAR_STRING(); dollarStr != nil {
		code := dollarStr.GetText()
		// Remove the $$ delimiters
		if len(code) >= 4 && strings.HasPrefix(code, "$$") && strings.HasSuffix(code, "$$") {
			code = code[2 : len(code)-2]
		}
		// Trim leading/trailing whitespace but preserve internal formatting
		stmt.JavaCode = strings.TrimSpace(code)
	}

	// Check for documentation comment from parent createStatement
	if parent, ok := ctx.GetParent().(*parser.CreateStatementContext); ok {
		if docComment := parent.DocComment(); docComment != nil {
			stmt.Documentation = extractDocComment(docComment.GetText())
		}
	}

	// Also check for doc comment at statement level (grammar allows it at both levels)
	if stmt.Documentation == "" {
		if stmtCtx := findParentStatement(ctx); stmtCtx != nil {
			if docCtx := stmtCtx.DocComment(); docCtx != nil {
				stmt.Documentation = extractDocComment(docCtx.GetText())
			}
		}
	}

	b.statements = append(b.statements, stmt)
}
