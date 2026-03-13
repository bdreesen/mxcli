// SPDX-License-Identifier: Apache-2.0

package visitor

import (
	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/mdl/grammar/parser"
)

func (b *Builder) ExitCreateModuleStatement(ctx *parser.CreateModuleStatementContext) {
	name := ""
	if id := ctx.IDENTIFIER(); id != nil {
		name = id.GetText()
	}
	b.statements = append(b.statements, &ast.CreateModuleStmt{
		Name: name,
	})
}

// ----------------------------------------------------------------------------
// Enumeration Statements
// ----------------------------------------------------------------------------

// ExitCreateEnumerationStatement is called when exiting the createEnumerationStatement production.
