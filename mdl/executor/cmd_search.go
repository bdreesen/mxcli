// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"fmt"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/ast"
	mdlerrors "github.com/mendixlabs/mxcli/mdl/errors"
)

// execShowCallers handles SHOW CALLERS OF Module.Microflow [TRANSITIVE].
func (e *Executor) execShowCallers(s *ast.ShowStmt) error {
	if s.Name == nil {
		return mdlerrors.NewValidation("target name required for SHOW CALLERS")
	}

	// Ensure catalog is available with full mode for refs
	if err := e.ensureCatalog(true); err != nil {
		return err
	}

	targetName := s.Name.String()
	fmt.Fprintf(e.output, "\nCallers of %s", targetName)
	if s.Transitive {
		fmt.Fprintln(e.output, " (transitive)")
	} else {
		fmt.Fprintln(e.output, "")
	}

	var query string
	if s.Transitive {
		// Recursive CTE for transitive callers
		query = `
			WITH RECURSIVE callers_cte AS (
				SELECT SourceName as Caller, 1 as Depth
				FROM refs
				WHERE TargetName = ? AND RefKind = 'call'
				UNION ALL
				SELECT r.SourceName, c.Depth + 1
				FROM refs r
				JOIN callers_cte c ON r.TargetName = c.Caller
				WHERE r.RefKind = 'call' AND c.Depth < 10
			)
			SELECT DISTINCT Caller, MIN(Depth) as Depth
			FROM callers_cte
			GROUP BY Caller
			ORDER BY Depth, Caller
		`
	} else {
		// Direct callers only
		query = `
			SELECT DISTINCT SourceName as Caller, 1 as Depth
			FROM refs
			WHERE TargetName = ? AND RefKind = 'call'
			ORDER BY Caller
		`
	}

	result, err := e.catalog.Query(strings.Replace(query, "?", "'"+targetName+"'", 1))
	if err != nil {
		return mdlerrors.NewBackend("query callers", err)
	}

	if result.Count == 0 {
		fmt.Fprintln(e.output, "(no callers found)")
		return nil
	}

	fmt.Fprintf(e.output, "Found %d caller(s)\n", result.Count)
	e.outputCatalogResults(result)
	return nil
}

// execShowCallees handles SHOW CALLEES OF Module.Microflow [TRANSITIVE].
func (e *Executor) execShowCallees(s *ast.ShowStmt) error {
	if s.Name == nil {
		return mdlerrors.NewValidation("target name required for SHOW CALLEES")
	}

	// Ensure catalog is available with full mode for refs
	if err := e.ensureCatalog(true); err != nil {
		return err
	}

	sourceName := s.Name.String()
	fmt.Fprintf(e.output, "\nCallees of %s", sourceName)
	if s.Transitive {
		fmt.Fprintln(e.output, " (transitive)")
	} else {
		fmt.Fprintln(e.output, "")
	}

	var query string
	if s.Transitive {
		// Recursive CTE for transitive callees
		query = `
			WITH RECURSIVE callees_cte AS (
				SELECT TargetName as Callee, 1 as Depth
				FROM refs
				WHERE SourceName = ? AND RefKind = 'call'
				UNION ALL
				SELECT r.TargetName, c.Depth + 1
				FROM refs r
				JOIN callees_cte c ON r.SourceName = c.Callee
				WHERE r.RefKind = 'call' AND c.Depth < 10
			)
			SELECT DISTINCT Callee, MIN(Depth) as Depth
			FROM callees_cte
			GROUP BY Callee
			ORDER BY Depth, Callee
		`
	} else {
		// Direct callees only
		query = `
			SELECT DISTINCT TargetName as Callee, 1 as Depth
			FROM refs
			WHERE SourceName = ? AND RefKind = 'call'
			ORDER BY Callee
		`
	}

	result, err := e.catalog.Query(strings.Replace(query, "?", "'"+sourceName+"'", 1))
	if err != nil {
		return mdlerrors.NewBackend("query callees", err)
	}

	if result.Count == 0 {
		fmt.Fprintln(e.output, "(no callees found)")
		return nil
	}

	fmt.Fprintf(e.output, "Found %d callee(s)\n", result.Count)
	e.outputCatalogResults(result)
	return nil
}

// execShowReferences handles SHOW REFERENCES TO Module.Entity.
func (e *Executor) execShowReferences(s *ast.ShowStmt) error {
	if s.Name == nil {
		return mdlerrors.NewValidation("target name required for SHOW REFERENCES")
	}

	// Ensure catalog is available with full mode for refs
	if err := e.ensureCatalog(true); err != nil {
		return err
	}

	targetName := s.Name.String()
	fmt.Fprintf(e.output, "\nReferences to %s\n", targetName)

	// Find all references to this target
	query := `
		SELECT SourceType, SourceName, RefKind
		FROM refs
		WHERE TargetName = ?
		ORDER BY RefKind, SourceType, SourceName
	`

	result, err := e.catalog.Query(strings.Replace(query, "?", "'"+targetName+"'", 1))
	if err != nil {
		return mdlerrors.NewBackend("query references", err)
	}

	if result.Count == 0 {
		fmt.Fprintln(e.output, "(no references found)")
		return nil
	}

	fmt.Fprintf(e.output, "Found %d reference(s)\n", result.Count)
	e.outputCatalogResults(result)
	return nil
}

// execShowImpact handles SHOW IMPACT OF Module.Entity.
// This shows all elements that would be affected by changing the target.
func (e *Executor) execShowImpact(s *ast.ShowStmt) error {
	if s.Name == nil {
		return mdlerrors.NewValidation("target name required for SHOW IMPACT")
	}

	// Ensure catalog is available with full mode for refs
	if err := e.ensureCatalog(true); err != nil {
		return err
	}

	targetName := s.Name.String()
	fmt.Fprintf(e.output, "\nImpact analysis for %s\n", targetName)

	// Find all direct references to this target
	directQuery := `
		SELECT SourceType, SourceName, RefKind
		FROM refs
		WHERE TargetName = ?
		ORDER BY SourceType, SourceName
	`

	result, err := e.catalog.Query(strings.Replace(directQuery, "?", "'"+targetName+"'", 1))
	if err != nil {
		return mdlerrors.NewBackend("query impact", err)
	}

	if result.Count == 0 {
		fmt.Fprintln(e.output, "(no impact - element is not referenced)")
		return nil
	}

	// Group by type for summary
	typeCounts := make(map[string]int)
	for _, row := range result.Rows {
		if len(row) > 0 {
			if t, ok := row[0].(string); ok {
				typeCounts[t]++
			}
		}
	}

	fmt.Fprintf(e.output, "\nSummary:\n")
	for t, count := range typeCounts {
		fmt.Fprintf(e.output, "  %s: %d\n", t, count)
	}
	fmt.Fprintln(e.output)

	fmt.Fprintf(e.output, "Found %d affected element(s)\n", result.Count)
	e.outputCatalogResults(result)

	return nil
}
