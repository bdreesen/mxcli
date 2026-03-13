// SPDX-License-Identifier: Apache-2.0

// Package executor - Snippet commands (SHOW/DESCRIBE SNIPPETS)
package executor

import (
	"fmt"
	"sort"
	"strings"
)

// showSnippets handles SHOW SNIPPETS command.
func (e *Executor) showSnippets(moduleName string) error {
	// Get hierarchy for module/folder resolution
	h, err := e.getHierarchy()
	if err != nil {
		return fmt.Errorf("failed to build hierarchy: %w", err)
	}

	// Get all snippets
	snippets, err := e.reader.ListSnippets()
	if err != nil {
		return fmt.Errorf("failed to list snippets: %w", err)
	}

	// Collect rows and calculate column widths
	type row struct {
		qualifiedName string
		module        string
		name          string
		folderPath    string
		params        int
	}
	var rows []row
	qnWidth := len("Qualified Name")
	modWidth := len("Module")
	nameWidth := len("Name")
	pathWidth := len("Folder")
	paramsWidth := len("Params")

	for _, s := range snippets {
		modID := h.FindModuleID(s.ContainerID)
		modName := h.GetModuleName(modID)
		if moduleName == "" || modName == moduleName {
			qualifiedName := modName + "." + s.Name
			folderPath := h.BuildFolderPath(s.ContainerID)

			rows = append(rows, row{qualifiedName, modName, s.Name, folderPath, len(s.Parameters)})
			if len(qualifiedName) > qnWidth {
				qnWidth = len(qualifiedName)
			}
			if len(modName) > modWidth {
				modWidth = len(modName)
			}
			if len(s.Name) > nameWidth {
				nameWidth = len(s.Name)
			}
			if len(folderPath) > pathWidth {
				pathWidth = len(folderPath)
			}
			paramsStr := fmt.Sprintf("%d", len(s.Parameters))
			if len(paramsStr) > paramsWidth {
				paramsWidth = len(paramsStr)
			}
		}
	}

	// Sort by qualified name
	sort.Slice(rows, func(i, j int) bool {
		return strings.ToLower(rows[i].qualifiedName) < strings.ToLower(rows[j].qualifiedName)
	})

	// Markdown table with aligned columns
	fmt.Fprintf(e.output, "| %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		qnWidth, "Qualified Name", modWidth, "Module", nameWidth, "Name",
		pathWidth, "Folder", paramsWidth, "Params")
	fmt.Fprintf(e.output, "|-%s-|-%s-|-%s-|-%s-|-%s-|\n",
		strings.Repeat("-", qnWidth), strings.Repeat("-", modWidth), strings.Repeat("-", nameWidth),
		strings.Repeat("-", pathWidth), strings.Repeat("-", paramsWidth))
	for _, r := range rows {
		fmt.Fprintf(e.output, "| %-*s | %-*s | %-*s | %-*s | %-*d |\n",
			qnWidth, r.qualifiedName, modWidth, r.module, nameWidth, r.name,
			pathWidth, r.folderPath, paramsWidth, r.params)
	}
	fmt.Fprintf(e.output, "\n(%d snippets)\n", len(rows))
	return nil
}
