// SPDX-License-Identifier: Apache-2.0

// Package executor - Layout commands (SHOW LAYOUTS)
package executor

import (
	"fmt"
	"sort"
	"strings"
)

// showLayouts handles SHOW LAYOUTS command.
func (e *Executor) showLayouts(moduleName string) error {
	// Get hierarchy for module/folder resolution
	h, err := e.getHierarchy()
	if err != nil {
		return fmt.Errorf("failed to build hierarchy: %w", err)
	}

	// Get all layouts
	layouts, err := e.reader.ListLayouts()
	if err != nil {
		return fmt.Errorf("failed to list layouts: %w", err)
	}

	// Collect rows and calculate column widths
	type row struct {
		qualifiedName string
		module        string
		name          string
		folderPath    string
		layoutType    string
	}
	var rows []row
	qnWidth := len("Qualified Name")
	modWidth := len("Module")
	nameWidth := len("Name")
	pathWidth := len("Folder")
	typeWidth := len("Type")

	for _, l := range layouts {
		modID := h.FindModuleID(l.ContainerID)
		modName := h.GetModuleName(modID)
		if moduleName == "" || modName == moduleName {
			qualifiedName := modName + "." + l.Name
			folderPath := h.BuildFolderPath(l.ContainerID)
			layoutType := string(l.LayoutType)

			rows = append(rows, row{qualifiedName, modName, l.Name, folderPath, layoutType})
			if len(qualifiedName) > qnWidth {
				qnWidth = len(qualifiedName)
			}
			if len(modName) > modWidth {
				modWidth = len(modName)
			}
			if len(l.Name) > nameWidth {
				nameWidth = len(l.Name)
			}
			if len(folderPath) > pathWidth {
				pathWidth = len(folderPath)
			}
			if len(layoutType) > typeWidth {
				typeWidth = len(layoutType)
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
		pathWidth, "Folder", typeWidth, "Type")
	fmt.Fprintf(e.output, "|-%s-|-%s-|-%s-|-%s-|-%s-|\n",
		strings.Repeat("-", qnWidth), strings.Repeat("-", modWidth), strings.Repeat("-", nameWidth),
		strings.Repeat("-", pathWidth), strings.Repeat("-", typeWidth))
	for _, r := range rows {
		fmt.Fprintf(e.output, "| %-*s | %-*s | %-*s | %-*s | %-*s |\n",
			qnWidth, r.qualifiedName, modWidth, r.module, nameWidth, r.name,
			pathWidth, r.folderPath, typeWidth, r.layoutType)
	}
	fmt.Fprintf(e.output, "\n(%d layouts)\n", len(rows))
	return nil
}
