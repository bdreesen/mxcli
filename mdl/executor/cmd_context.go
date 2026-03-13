// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"fmt"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/ast"
)

// execShowContext handles SHOW CONTEXT OF <name> [DEPTH n] command.
// It assembles relevant context information for LLM consumption.
func (e *Executor) execShowContext(s *ast.ShowStmt) error {
	if s.Name == nil {
		return fmt.Errorf("SHOW CONTEXT requires a qualified name")
	}

	// Ensure catalog is built with full mode for refs
	if err := e.ensureCatalog(true); err != nil {
		return fmt.Errorf("failed to build catalog: %w", err)
	}

	name := s.Name.String()
	depth := s.Depth
	if depth <= 0 {
		depth = 2
	}

	// Detect the type of the target element
	targetType, err := e.detectElementType(name)
	if err != nil {
		return err
	}

	// Assemble context based on type
	var output strings.Builder
	output.WriteString(fmt.Sprintf("## Context: %s\n\n", name))

	switch targetType {
	case "microflow", "nanoflow":
		e.assembleMicroflowContext(&output, name, depth)
	case "entity":
		e.assembleEntityContext(&output, name, depth)
	case "page":
		e.assemblePageContext(&output, name, depth)
	case "enumeration":
		e.assembleEnumerationContext(&output, name)
	case "workflow":
		e.assembleWorkflowContext(&output, name, depth)
	case "snippet":
		e.assembleSnippetContext(&output, name, depth)
	case "javaaction":
		e.assembleJavaActionContext(&output, name)
	case "odataclient":
		e.assembleODataClientContext(&output, name)
	case "odataservice":
		e.assembleODataServiceContext(&output, name)
	default:
		output.WriteString(fmt.Sprintf("Unknown element type for: %s\n", name))
	}

	fmt.Fprint(e.output, output.String())
	return nil
}

// detectElementType determines what kind of element the name refers to.
func (e *Executor) detectElementType(name string) (string, error) {
	// Check catalog tables for known element types
	catalogChecks := []struct {
		table    string
		elemType string
	}{
		{"microflows", "microflow"},
		{"entities", "entity"},
		{"pages", "page"},
		{"enumerations", "enumeration"},
		{"snippets", "snippet"},
		{"workflows", "workflow"},
		{"java_actions", "javaaction"},
		{"odata_clients", "odataclient"},
		{"odata_services", "odataservice"},
	}

	for _, check := range catalogChecks {
		result, err := e.catalog.Query(fmt.Sprintf(
			"SELECT 1 FROM %s WHERE QualifiedName = '%s' LIMIT 1", check.table, name))
		if err == nil && result.Count > 0 {
			return check.elemType, nil
		}
	}

	return "", fmt.Errorf("element not found: %s", name)
}

// assembleMicroflowContext assembles context for a microflow.
func (e *Executor) assembleMicroflowContext(out *strings.Builder, name string, depth int) {
	// Get microflow basic info
	out.WriteString("### Microflow Definition\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		"SELECT Name, ReturnType, ParameterCount, ActivityCount FROM microflows WHERE QualifiedName = '%s'", name))
	if err == nil && result.Count > 0 {
		row := result.Rows[0]
		out.WriteString(fmt.Sprintf("- **Name**: %v\n", row[0]))
		out.WriteString(fmt.Sprintf("- **Return Type**: %v\n", row[1]))
		out.WriteString(fmt.Sprintf("- **Parameters**: %v\n", row[2]))
		out.WriteString(fmt.Sprintf("- **Activities**: %v\n", row[3]))
	}
	out.WriteString("\n")

	// Entities used by this microflow
	out.WriteString("### Entities Used\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName, RefKind FROM refs
		 WHERE SourceName = '%s' AND TargetType = 'entity'
		 ORDER BY RefKind, TargetName`, name))
	if err == nil && result.Count > 0 {
		out.WriteString("| Entity | Usage |\n")
		out.WriteString("|--------|-------|\n")
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("| %v | %v |\n", row[0], row[1]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Pages shown by this microflow
	out.WriteString("### Pages Shown\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName FROM refs
		 WHERE SourceName = '%s' AND RefKind = 'show_page'
		 ORDER BY TargetName`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none)\n")
	}
	out.WriteString("\n")

	// Called microflows (with depth)
	out.WriteString(fmt.Sprintf("### Called Microflows (depth %d)\n\n", depth))
	if depth > 0 {
		e.addCallees(out, name, depth, 1)
	}
	out.WriteString("\n")

	// Direct callers
	out.WriteString("### Direct Callers\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT SourceName FROM refs
		 WHERE TargetName = '%s' AND RefKind = 'call'
		 ORDER BY SourceName LIMIT 10`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
		if result.Count == 10 {
			out.WriteString("- ... (more callers exist)\n")
		}
	} else {
		out.WriteString("(none)\n")
	}
}

// addCallees recursively adds callees up to the specified depth.
func (e *Executor) addCallees(out *strings.Builder, name string, maxDepth, currentDepth int) {
	if currentDepth > maxDepth {
		return
	}

	indent := strings.Repeat("  ", currentDepth-1)
	result, err := e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName FROM refs
		 WHERE SourceName = '%s' AND RefKind = 'call'
		 ORDER BY TargetName`, name))
	if err != nil || result.Count == 0 {
		return
	}

	for _, row := range result.Rows {
		callee := fmt.Sprintf("%v", row[0])
		out.WriteString(fmt.Sprintf("%s- %s\n", indent, callee))
		// Recurse for deeper levels
		if currentDepth < maxDepth {
			e.addCallees(out, callee, maxDepth, currentDepth+1)
		}
	}
}

// assembleEntityContext assembles context for an entity.
func (e *Executor) assembleEntityContext(out *strings.Builder, name string, depth int) {
	// Get entity basic info
	out.WriteString("### Entity Definition\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		"SELECT Name, EntityType, Generalization, AttributeCount, IndexCount FROM entities WHERE QualifiedName = '%s'", name))
	if err == nil && result.Count > 0 {
		row := result.Rows[0]
		out.WriteString(fmt.Sprintf("- **Name**: %v\n", row[0]))
		out.WriteString(fmt.Sprintf("- **Type**: %v\n", row[1]))
		if row[2] != nil && row[2] != "" {
			out.WriteString(fmt.Sprintf("- **Extends**: %v\n", row[2]))
		}
		out.WriteString(fmt.Sprintf("- **Attributes**: %v\n", row[3]))
		out.WriteString(fmt.Sprintf("- **Indexes**: %v\n", row[4]))
	}
	out.WriteString("\n")

	// Microflows that use this entity
	out.WriteString("### Microflows Using This Entity\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT SourceName, RefKind FROM refs
		 WHERE TargetName = '%s' AND SourceType = 'microflow'
		 ORDER BY RefKind, SourceName LIMIT 20`, name))
	if err == nil && result.Count > 0 {
		out.WriteString("| Microflow | Usage |\n")
		out.WriteString("|-----------|-------|\n")
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("| %v | %v |\n", row[0], row[1]))
		}
		if result.Count == 20 {
			out.WriteString("\n(limited to 20 results)\n")
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Pages displaying this entity
	out.WriteString("### Pages Displaying This Entity\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT SourceName FROM refs
		 WHERE TargetName = '%s' AND SourceType = 'page'
		 ORDER BY SourceName LIMIT 10`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Related entities (via associations or generalization)
	out.WriteString("### Related Entities\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName, RefKind FROM refs
		 WHERE SourceName = '%s' AND TargetType = 'entity'
		 UNION
		 SELECT DISTINCT SourceName, RefKind FROM refs
		 WHERE TargetName = '%s' AND SourceType = 'entity'
		 ORDER BY RefKind, TargetName LIMIT 10`, name, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v (%v)\n", row[0], row[1]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
}

// assemblePageContext assembles context for a page.
func (e *Executor) assemblePageContext(out *strings.Builder, name string, depth int) {
	// Get page basic info
	out.WriteString("### Page Definition\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		"SELECT Name, Title, URL, LayoutRef, WidgetCount FROM pages WHERE QualifiedName = '%s'", name))
	if err == nil && result.Count > 0 {
		row := result.Rows[0]
		out.WriteString(fmt.Sprintf("- **Name**: %v\n", row[0]))
		if row[1] != nil && row[1] != "" {
			out.WriteString(fmt.Sprintf("- **Title**: %v\n", row[1]))
		}
		if row[2] != nil && row[2] != "" {
			out.WriteString(fmt.Sprintf("- **URL**: %v\n", row[2]))
		}
		if row[3] != nil && row[3] != "" {
			out.WriteString(fmt.Sprintf("- **Layout**: %v\n", row[3]))
		}
		out.WriteString(fmt.Sprintf("- **Widgets**: %v\n", row[4]))
	}
	out.WriteString("\n")

	// Entities used on this page
	out.WriteString("### Entities Used\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName FROM refs
		 WHERE SourceName = '%s' AND TargetType = 'entity'
		 ORDER BY TargetName`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Microflows called from this page
	out.WriteString("### Microflows Called\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName FROM refs
		 WHERE SourceName = '%s' AND TargetType = 'microflow'
		 ORDER BY TargetName LIMIT 15`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Microflows that show this page
	out.WriteString("### Shown By\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT SourceName FROM refs
		 WHERE TargetName = '%s' AND RefKind = 'show_page'
		 ORDER BY SourceName LIMIT 10`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
}

// assembleEnumerationContext assembles context for an enumeration.
func (e *Executor) assembleEnumerationContext(out *strings.Builder, name string) {
	// Get enumeration basic info
	out.WriteString("### Enumeration Definition\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		"SELECT Name, ValueCount FROM enumerations WHERE QualifiedName = '%s'", name))
	if err == nil && result.Count > 0 {
		row := result.Rows[0]
		out.WriteString(fmt.Sprintf("- **Name**: %v\n", row[0]))
		out.WriteString(fmt.Sprintf("- **Values**: %v\n", row[1]))
	}
	out.WriteString("\n")

	// Entities with attributes of this enumeration type
	out.WriteString("### Used By Entities\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT SourceName FROM refs
		 WHERE TargetName = '%s' AND SourceType = 'entity'
		 ORDER BY SourceName LIMIT 15`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Microflows that use this enumeration
	out.WriteString("### Used By Microflows\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT SourceName FROM refs
		 WHERE TargetName = '%s' AND SourceType = 'microflow'
		 ORDER BY SourceName LIMIT 15`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
}

// assembleSnippetContext assembles context for a snippet.
func (e *Executor) assembleSnippetContext(out *strings.Builder, name string, depth int) {
	out.WriteString("### Snippet Definition\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		"SELECT Name, ParameterCount, WidgetCount FROM snippets WHERE QualifiedName = '%s'", name))
	if err == nil && result.Count > 0 {
		row := result.Rows[0]
		out.WriteString(fmt.Sprintf("- **Name**: %v\n", row[0]))
		out.WriteString(fmt.Sprintf("- **Parameters**: %v\n", row[1]))
		out.WriteString(fmt.Sprintf("- **Widgets**: %v\n", row[2]))
	}
	out.WriteString("\n")

	// MDL source via DESCRIBE
	out.WriteString("### MDL Source\n\n```sql\n")
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		descStmt := &ast.DescribeStmt{
			ObjectType: ast.DescribeSnippet,
			Name:       ast.QualifiedName{Module: parts[0], Name: parts[1]},
		}
		savedOutput := e.output
		e.output = out
		e.execDescribe(descStmt)
		e.output = savedOutput
	}
	out.WriteString("```\n\n")

	// Pages that use this snippet
	out.WriteString("### Used By Pages\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT SourceName FROM refs
		 WHERE TargetName = '%s' AND RefKind = 'snippet_call'
		 ORDER BY SourceName LIMIT 15`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
}

// assembleJavaActionContext assembles context for a java action.
func (e *Executor) assembleJavaActionContext(out *strings.Builder, name string) {
	out.WriteString("### Java Action Definition\n\n```sql\n")
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		descStmt := &ast.DescribeStmt{
			ObjectType: ast.DescribeJavaAction,
			Name:       ast.QualifiedName{Module: parts[0], Name: parts[1]},
		}
		savedOutput := e.output
		e.output = out
		e.execDescribe(descStmt)
		e.output = savedOutput
	}
	out.WriteString("```\n\n")

	// Microflows that call this java action
	out.WriteString("### Called By Microflows\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT SourceName FROM refs
		 WHERE TargetName = '%s' AND RefKind = 'call'
		 ORDER BY SourceName LIMIT 15`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
}

// assembleODataClientContext assembles context for a consumed OData service.
func (e *Executor) assembleODataClientContext(out *strings.Builder, name string) {
	out.WriteString("### Consumed OData Service\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		"SELECT Name, Version, ODataVersion, MetadataUrl FROM odata_clients WHERE QualifiedName = '%s'", name))
	if err == nil && result.Count > 0 {
		row := result.Rows[0]
		out.WriteString(fmt.Sprintf("- **Name**: %v\n", row[0]))
		out.WriteString(fmt.Sprintf("- **Version**: %v\n", row[1]))
		out.WriteString(fmt.Sprintf("- **OData Version**: %v\n", row[2]))
		out.WriteString(fmt.Sprintf("- **Metadata URL**: %v\n", row[3]))
	}
	out.WriteString("\n")

	// External entities from this service
	out.WriteString("### External Entities\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName FROM refs
		 WHERE SourceName = '%s' AND RefKind = 'odata_entity'
		 ORDER BY TargetName LIMIT 15`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
}

// assembleWorkflowContext assembles context for a workflow.
func (e *Executor) assembleWorkflowContext(out *strings.Builder, name string, depth int) {
	// Get workflow basic info
	out.WriteString("### Workflow Definition\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		"SELECT Name, ParameterEntity, ActivityCount, UserTaskCount, MicroflowCallCount, DecisionCount, Description FROM workflows WHERE QualifiedName = '%s'", name))
	if err == nil && result.Count > 0 {
		row := result.Rows[0]
		out.WriteString(fmt.Sprintf("- **Name**: %v\n", row[0]))
		if row[1] != nil && row[1] != "" {
			out.WriteString(fmt.Sprintf("- **Parameter Entity**: %v\n", row[1]))
		}
		out.WriteString(fmt.Sprintf("- **Activities**: %v\n", row[2]))
		out.WriteString(fmt.Sprintf("- **User Tasks**: %v\n", row[3]))
		out.WriteString(fmt.Sprintf("- **Microflow Calls**: %v\n", row[4]))
		out.WriteString(fmt.Sprintf("- **Decisions**: %v\n", row[5]))
		if row[6] != nil && row[6] != "" {
			out.WriteString(fmt.Sprintf("- **Description**: %v\n", row[6]))
		}
	}
	out.WriteString("\n")

	// MDL source via DESCRIBE
	out.WriteString("### MDL Source\n\n```sql\n")
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		descStmt := &ast.DescribeStmt{
			ObjectType: ast.DescribeWorkflow,
			Name:       ast.QualifiedName{Module: parts[0], Name: parts[1]},
		}
		savedOutput := e.output
		e.output = out
		e.execDescribe(descStmt)
		e.output = savedOutput
	}
	out.WriteString("```\n\n")

	// Microflows called by this workflow
	out.WriteString("### Microflows Called\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName, RefKind FROM refs
		 WHERE SourceName = '%s' AND TargetType = 'MICROFLOW'
		 ORDER BY RefKind, TargetName`, name))
	if err == nil && result.Count > 0 {
		out.WriteString("| Microflow | Usage |\n")
		out.WriteString("|-----------|-------|\n")
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("| %v | %v |\n", row[0], row[1]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Pages used by this workflow (user task pages, overview page)
	out.WriteString("### Pages Used\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName, RefKind FROM refs
		 WHERE SourceName = '%s' AND TargetType = 'PAGE'
		 ORDER BY TargetName`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v (%v)\n", row[0], row[1]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Entities referenced by this workflow
	out.WriteString("### Entities Used\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName, RefKind FROM refs
		 WHERE SourceName = '%s' AND TargetType = 'ENTITY'
		 ORDER BY TargetName`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v (%v)\n", row[0], row[1]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
	out.WriteString("\n")

	// Direct callers (what calls this workflow)
	out.WriteString("### Direct Callers\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT SourceName, SourceType FROM refs
		 WHERE TargetName = '%s'
		 ORDER BY SourceName LIMIT 15`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v (%v)\n", row[0], row[1]))
		}
		if result.Count == 15 {
			out.WriteString("- ... (more callers exist)\n")
		}
	} else {
		out.WriteString("(none found)\n")
	}
}

// assembleODataServiceContext assembles context for a published OData service.
func (e *Executor) assembleODataServiceContext(out *strings.Builder, name string) {
	out.WriteString("### Published OData Service\n\n")
	result, err := e.catalog.Query(fmt.Sprintf(
		"SELECT Name, Path, Version, ODataVersion, EntitySetCount FROM odata_services WHERE QualifiedName = '%s'", name))
	if err == nil && result.Count > 0 {
		row := result.Rows[0]
		out.WriteString(fmt.Sprintf("- **Name**: %v\n", row[0]))
		out.WriteString(fmt.Sprintf("- **Path**: %v\n", row[1]))
		out.WriteString(fmt.Sprintf("- **Version**: %v\n", row[2]))
		out.WriteString(fmt.Sprintf("- **OData Version**: %v\n", row[3]))
		out.WriteString(fmt.Sprintf("- **Entity Sets**: %v\n", row[4]))
	}
	out.WriteString("\n")

	// Published entities
	out.WriteString("### Published Entities\n\n")
	result, err = e.catalog.Query(fmt.Sprintf(
		`SELECT DISTINCT TargetName FROM refs
		 WHERE SourceName = '%s' AND RefKind = 'odata_publish'
		 ORDER BY TargetName LIMIT 15`, name))
	if err == nil && result.Count > 0 {
		for _, row := range result.Rows {
			out.WriteString(fmt.Sprintf("- %v\n", row[0]))
		}
	} else {
		out.WriteString("(none found)\n")
	}
}
