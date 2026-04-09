// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"fmt"
	"sort"
	"strings"
)

// showPages handles SHOW PAGES command.
func (e *Executor) showPages(moduleName string) error {
	// Get hierarchy for module/folder resolution
	h, err := e.getHierarchy()
	if err != nil {
		return fmt.Errorf("failed to build hierarchy: %w", err)
	}

	// Get all pages
	pages, err := e.reader.ListPages()
	if err != nil {
		return fmt.Errorf("failed to list pages: %w", err)
	}

	// Collect rows and calculate column widths
	type row struct {
		qualifiedName string
		module        string
		name          string
		folderPath    string
		title         string
		url           string
		params        int
	}
	var rows []row
	qnWidth := len("Qualified Name")
	modWidth := len("Module")
	nameWidth := len("Name")
	pathWidth := len("Folder")
	titleWidth := len("Title")
	urlWidth := len("URL")
	paramsWidth := len("Params")

	for _, p := range pages {
		modID := h.FindModuleID(p.ContainerID)
		modName := h.GetModuleName(modID)
		if moduleName == "" || modName == moduleName {
			qualifiedName := modName + "." + p.Name
			if p.Excluded {
				qualifiedName += " [EXCLUDED]"
			}
			folderPath := h.BuildFolderPath(p.ContainerID)
			title := ""
			if p.Title != nil {
				// Try to get English title first, then any available translation
				title = p.Title.GetTranslation("en_US")
				if title == "" {
					for _, t := range p.Title.Translations {
						title = t
						break
					}
				}
			}
			url := p.URL

			rows = append(rows, row{qualifiedName, modName, p.Name, folderPath, title, url, len(p.Parameters)})
			if len(qualifiedName) > qnWidth {
				qnWidth = len(qualifiedName)
			}
			if len(modName) > modWidth {
				modWidth = len(modName)
			}
			if len(p.Name) > nameWidth {
				nameWidth = len(p.Name)
			}
			if len(folderPath) > pathWidth {
				pathWidth = len(folderPath)
			}
			if len(title) > titleWidth {
				titleWidth = len(title)
			}
			if len(url) > urlWidth {
				urlWidth = len(url)
			}
			paramsStr := fmt.Sprintf("%d", len(p.Parameters))
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
	fmt.Fprintf(e.output, "| %-*s | %-*s | %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		qnWidth, "Qualified Name", modWidth, "Module", nameWidth, "Name",
		pathWidth, "Folder", titleWidth, "Title", urlWidth, "URL", paramsWidth, "Params")
	fmt.Fprintf(e.output, "|-%s-|-%s-|-%s-|-%s-|-%s-|-%s-|-%s-|\n",
		strings.Repeat("-", qnWidth), strings.Repeat("-", modWidth), strings.Repeat("-", nameWidth),
		strings.Repeat("-", pathWidth), strings.Repeat("-", titleWidth), strings.Repeat("-", urlWidth),
		strings.Repeat("-", paramsWidth))
	for _, r := range rows {
		fmt.Fprintf(e.output, "| %-*s | %-*s | %-*s | %-*s | %-*s | %-*s | %-*d |\n",
			qnWidth, r.qualifiedName, modWidth, r.module, nameWidth, r.name,
			pathWidth, r.folderPath, titleWidth, r.title, urlWidth, r.url, paramsWidth, r.params)
	}
	fmt.Fprintf(e.output, "\n(%d pages)\n", len(rows))
	return nil
}
