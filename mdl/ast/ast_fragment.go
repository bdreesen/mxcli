// SPDX-License-Identifier: Apache-2.0

package ast

// DefineFragmentStmt represents: DEFINE FRAGMENT Name AS { widgets }
type DefineFragmentStmt struct {
	Name    string
	Widgets []*WidgetV3
}

func (s *DefineFragmentStmt) isStatement() {}

// DescribeFragmentFromStmt represents DESCRIBE FRAGMENT FROM PAGE/SNIPPET ... WIDGET ...
type DescribeFragmentFromStmt struct {
	ContainerType string        // "PAGE" or "SNIPPET"
	ContainerName QualifiedName // Module.PageName or Module.SnippetName
	WidgetName    string        // Target widget name
}

func (s *DescribeFragmentFromStmt) isStatement() {}
