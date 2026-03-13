// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"regexp"
	"strings"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// symbolPattern matches CREATE/ALTER/DROP statements and extracts the type and qualified name.
// Group 1: action (CREATE, ALTER, DROP)
// Group 2: optional entity modifier (PERSISTENT, NON-PERSISTENT, VIEW, EXTERNAL)
// Group 3: object type keyword (ENTITY, MICROFLOW, etc.)
// Group 4: qualified name
var symbolPattern = regexp.MustCompile(`(?i)^\s*(CREATE(?:\s+OR\s+(?:REPLACE|MODIFY))?|ALTER|DROP)\s+(?:(PERSISTENT|NON[_-]PERSISTENT|VIEW|EXTERNAL)\s+)?(ENTITY|MICROFLOW|NANOFLOW|PAGE|SNIPPET|LAYOUT|ENUMERATION|ASSOCIATION|CONSTANT|MODULE|JAVA\s+ACTION)\s+(\S+)`)

// symbolKindMap maps MDL object types to LSP symbol kinds.
var symbolKindMap = map[string]protocol.SymbolKind{
	"ENTITY":      protocol.SymbolKindClass,
	"MICROFLOW":   protocol.SymbolKindFunction,
	"NANOFLOW":    protocol.SymbolKindFunction,
	"PAGE":        protocol.SymbolKindFile,
	"SNIPPET":     protocol.SymbolKindFile,
	"LAYOUT":      protocol.SymbolKindFile,
	"ENUMERATION": protocol.SymbolKindEnum,
	"ASSOCIATION": protocol.SymbolKindInterface,
	"CONSTANT":    protocol.SymbolKindConstant,
	"MODULE":      protocol.SymbolKindModule,
	"JAVA ACTION": protocol.SymbolKindFunction,
}

// parsedSymbol holds extracted information about a CREATE/ALTER/DROP statement.
type parsedSymbol struct {
	action    string // CREATE, ALTER, DROP
	modifier  string // PERSISTENT, NON-PERSISTENT, etc. (optional)
	objType   string // ENTITY, MICROFLOW, PAGE, etc.
	module    string // Module part of qualified name
	shortName string // Name part of qualified name
	folder    string // Folder path (if specified)
	startLine int
	endLine   int
	nameStart int
}

// folderPatternMicroflow matches FOLDER 'path' for microflows/nanoflows (keyword syntax)
var folderPatternMicroflow = regexp.MustCompile(`(?i)\bFOLDER\s+'([^']*)'`)

// folderPatternPage matches Folder: 'path' for pages/snippets (property syntax)
var folderPatternPage = regexp.MustCompile(`(?i)\bFolder\s*:\s*'([^']*)'`)

// DocumentSymbol handles textDocument/documentSymbol requests.
func (s *mdlServer) DocumentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) ([]any, error) {
	docURI := uri.URI(params.TextDocument.URI)
	s.mu.Lock()
	text := s.docs[docURI]
	s.mu.Unlock()

	if text == "" {
		return nil, nil
	}

	symbols := extractDocumentSymbols(text)
	result := make([]any, len(symbols))
	for i := range symbols {
		result[i] = symbols[i]
	}
	return result, nil
}

// extractDocumentSymbols scans document text for CREATE/ALTER/DROP statements
// and returns them as DocumentSymbol objects organized like Mendix Studio Pro:
// - Module (top level)
//   - Domain Model (contains entities, associations, enumerations)
//   - Folders (contain microflows, pages, etc.)
//   - Documents without folders (directly under module)
func extractDocumentSymbols(text string) []protocol.DocumentSymbol {
	lines := strings.Split(text, "\n")

	// First pass: collect all symbols with their metadata
	var parsed []parsedSymbol
	for i, line := range lines {
		matches := symbolPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		action := strings.ToUpper(matches[1])
		modifier := strings.ToUpper(matches[2])
		objType := strings.ToUpper(strings.Join(strings.Fields(matches[3]), " "))
		qualName := strings.TrimRight(matches[4], ";(,{")

		if qualName == "" {
			continue
		}

		// Split qualified name into module and short name
		module, shortName := splitQualifiedName(qualName)

		// Find the end of this statement
		endLine := findStatementEnd(lines, i)

		// Extract folder path from the statement body
		folder := extractFolderPath(lines, i, endLine, objType)

		nameStart := max(strings.Index(line, qualName), 0)

		parsed = append(parsed, parsedSymbol{
			action:    action,
			modifier:  modifier,
			objType:   objType,
			module:    module,
			shortName: shortName,
			folder:    folder,
			startLine: i,
			endLine:   endLine,
			nameStart: nameStart,
		})
	}

	// Second pass: organize into hierarchical structure
	return buildSymbolHierarchy(parsed, lines)
}

// splitQualifiedName splits "Module.Name" into ("Module", "Name").
// If no dot, returns ("", name).
func splitQualifiedName(qualName string) (module, name string) {
	if idx := strings.Index(qualName, "."); idx > 0 {
		return qualName[:idx], qualName[idx+1:]
	}
	return "", qualName
}

// extractFolderPath finds the folder path in a statement's body.
func extractFolderPath(lines []string, startLine, endLine int, objType string) string {
	// Build the statement text
	var sb strings.Builder
	for i := startLine; i <= endLine && i < len(lines); i++ {
		sb.WriteString(lines[i])
		sb.WriteString("\n")
	}
	stmtText := sb.String()

	// Use different patterns based on object type
	switch objType {
	case "MICROFLOW", "NANOFLOW":
		// Keyword syntax: FOLDER 'path'
		if m := folderPatternMicroflow.FindStringSubmatch(stmtText); m != nil {
			return m[1]
		}
	case "PAGE", "SNIPPET", "LAYOUT":
		// Property syntax: Folder: 'path'
		if m := folderPatternPage.FindStringSubmatch(stmtText); m != nil {
			return m[1]
		}
	}
	return ""
}

// buildSymbolHierarchy creates a hierarchical symbol tree organized like Studio Pro.
func buildSymbolHierarchy(parsed []parsedSymbol, lines []string) []protocol.DocumentSymbol {
	// Group by module
	moduleMap := make(map[string][]parsedSymbol)
	var moduleOrder []string

	for _, p := range parsed {
		mod := p.module
		if mod == "" {
			mod = "(no module)"
		}
		if _, exists := moduleMap[mod]; !exists {
			moduleOrder = append(moduleOrder, mod)
		}
		moduleMap[mod] = append(moduleMap[mod], p)
	}

	var result []protocol.DocumentSymbol
	for _, mod := range moduleOrder {
		symbols := moduleMap[mod]
		moduleSymbol := buildModuleSymbol(mod, symbols, lines)
		result = append(result, moduleSymbol)
	}
	return result
}

// buildModuleSymbol creates a module symbol with its children organized by folder.
func buildModuleSymbol(moduleName string, symbols []parsedSymbol, lines []string) protocol.DocumentSymbol {
	// Separate domain model items from other documents
	var domainModelItems []parsedSymbol
	var otherItems []parsedSymbol

	for _, p := range symbols {
		switch p.objType {
		case "ENTITY", "ASSOCIATION", "ENUMERATION":
			domainModelItems = append(domainModelItems, p)
		default:
			otherItems = append(otherItems, p)
		}
	}

	// Build children
	var children []protocol.DocumentSymbol

	// Domain Model container (if any domain model items exist)
	if len(domainModelItems) > 0 {
		dmChildren := make([]protocol.DocumentSymbol, 0, len(domainModelItems))
		for _, p := range domainModelItems {
			dmChildren = append(dmChildren, createDocumentSymbol(p, lines))
		}

		// Find range spanning all domain model items
		startLine := domainModelItems[0].startLine
		endLine := domainModelItems[0].endLine
		for _, p := range domainModelItems[1:] {
			if p.startLine < startLine {
				startLine = p.startLine
			}
			if p.endLine > endLine {
				endLine = p.endLine
			}
		}

		children = append(children, protocol.DocumentSymbol{
			Name:     "Domain Model",
			Detail:   "",
			Kind:     protocol.SymbolKindNamespace,
			Children: dmChildren,
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(startLine), Character: 0},
				End:   protocol.Position{Line: uint32(endLine), Character: uint32(len(lines[endLine]))},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{Line: uint32(startLine), Character: 0},
				End:   protocol.Position{Line: uint32(startLine), Character: 12}, // "Domain Model"
			},
		})
	}

	// Group other items by folder
	folderMap := make(map[string][]parsedSymbol)
	var folderOrder []string
	var rootItems []parsedSymbol

	for _, p := range otherItems {
		if p.folder == "" {
			rootItems = append(rootItems, p)
		} else {
			// Get top-level folder
			topFolder := p.folder
			if idx := strings.Index(p.folder, "/"); idx > 0 {
				topFolder = p.folder[:idx]
			}
			if _, exists := folderMap[topFolder]; !exists {
				folderOrder = append(folderOrder, topFolder)
			}
			folderMap[topFolder] = append(folderMap[topFolder], p)
		}
	}

	// Add folder containers
	for _, folder := range folderOrder {
		folderItems := folderMap[folder]
		folderChildren := buildFolderChildren(folder, folderItems, lines)

		// Find range spanning all folder items
		startLine := folderItems[0].startLine
		endLine := folderItems[0].endLine
		for _, p := range folderItems[1:] {
			if p.startLine < startLine {
				startLine = p.startLine
			}
			if p.endLine > endLine {
				endLine = p.endLine
			}
		}

		children = append(children, protocol.DocumentSymbol{
			Name:     folder,
			Detail:   "folder",
			Kind:     protocol.SymbolKindPackage,
			Children: folderChildren,
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(startLine), Character: 0},
				End:   protocol.Position{Line: uint32(endLine), Character: uint32(len(lines[endLine]))},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{Line: uint32(startLine), Character: 0},
				End:   protocol.Position{Line: uint32(startLine), Character: uint32(len(folder))},
			},
		})
	}

	// Add root items (no folder)
	for _, p := range rootItems {
		children = append(children, createDocumentSymbol(p, lines))
	}

	// Find range spanning the entire module
	startLine := symbols[0].startLine
	endLine := symbols[0].endLine
	for _, p := range symbols[1:] {
		if p.startLine < startLine {
			startLine = p.startLine
		}
		if p.endLine > endLine {
			endLine = p.endLine
		}
	}

	return protocol.DocumentSymbol{
		Name:     moduleName,
		Detail:   "module",
		Kind:     protocol.SymbolKindModule,
		Children: children,
		Range: protocol.Range{
			Start: protocol.Position{Line: uint32(startLine), Character: 0},
			End:   protocol.Position{Line: uint32(endLine), Character: uint32(len(lines[endLine]))},
		},
		SelectionRange: protocol.Range{
			Start: protocol.Position{Line: uint32(startLine), Character: 0},
			End:   protocol.Position{Line: uint32(startLine), Character: uint32(len(moduleName))},
		},
	}
}

// buildFolderChildren creates symbols for items in a folder, handling nested folders.
func buildFolderChildren(parentFolder string, items []parsedSymbol, lines []string) []protocol.DocumentSymbol {
	// Group by subfolder
	subfolderMap := make(map[string][]parsedSymbol)
	var subfolderOrder []string
	var directItems []parsedSymbol

	for _, p := range items {
		// Check if this item is in a subfolder
		if len(p.folder) > len(parentFolder) && strings.HasPrefix(p.folder, parentFolder+"/") {
			remaining := p.folder[len(parentFolder)+1:]
			// Get next folder level
			nextFolder := remaining
			if idx := strings.Index(remaining, "/"); idx > 0 {
				nextFolder = remaining[:idx]
			}
			if _, exists := subfolderMap[nextFolder]; !exists {
				subfolderOrder = append(subfolderOrder, nextFolder)
			}
			subfolderMap[nextFolder] = append(subfolderMap[nextFolder], p)
		} else {
			directItems = append(directItems, p)
		}
	}

	var children []protocol.DocumentSymbol

	// Add subfolders recursively
	for _, subfolder := range subfolderOrder {
		subItems := subfolderMap[subfolder]
		fullPath := parentFolder + "/" + subfolder
		subChildren := buildFolderChildren(fullPath, subItems, lines)

		startLine := subItems[0].startLine
		endLine := subItems[0].endLine
		for _, p := range subItems[1:] {
			if p.startLine < startLine {
				startLine = p.startLine
			}
			if p.endLine > endLine {
				endLine = p.endLine
			}
		}

		children = append(children, protocol.DocumentSymbol{
			Name:     subfolder,
			Detail:   "folder",
			Kind:     protocol.SymbolKindPackage,
			Children: subChildren,
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(startLine), Character: 0},
				End:   protocol.Position{Line: uint32(endLine), Character: uint32(len(lines[endLine]))},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{Line: uint32(startLine), Character: 0},
				End:   protocol.Position{Line: uint32(startLine), Character: uint32(len(subfolder))},
			},
		})
	}

	// Add direct items
	for _, p := range directItems {
		children = append(children, createDocumentSymbol(p, lines))
	}

	return children
}

// createDocumentSymbol creates a leaf symbol for a document.
func createDocumentSymbol(p parsedSymbol, lines []string) protocol.DocumentSymbol {
	fullType := p.objType
	if p.modifier != "" {
		fullType = p.modifier + " " + p.objType
	}

	kind, ok := symbolKindMap[p.objType]
	if !ok {
		kind = protocol.SymbolKindVariable
	}

	endLineLen := 0
	if p.endLine < len(lines) {
		endLineLen = len(lines[p.endLine])
	}

	return protocol.DocumentSymbol{
		Name:   p.shortName,
		Detail: fullType,
		Kind:   kind,
		Range: protocol.Range{
			Start: protocol.Position{Line: uint32(p.startLine), Character: 0},
			End:   protocol.Position{Line: uint32(p.endLine), Character: uint32(endLineLen)},
		},
		SelectionRange: protocol.Range{
			Start: protocol.Position{Line: uint32(p.startLine), Character: uint32(p.nameStart)},
			End:   protocol.Position{Line: uint32(p.startLine), Character: uint32(p.nameStart + len(p.shortName))},
		},
	}
}

// findStatementEnd finds the last line of a statement starting at startLine.
// It looks for the next top-level statement or a line ending with a semicolon
// or END at the start of a line.
func findStatementEnd(lines []string, startLine int) int {
	endLine := startLine
	for j := startLine + 1; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])
		if trimmed == "" {
			continue
		}
		// Next top-level statement starts
		if symbolPattern.MatchString(lines[j]) {
			break
		}
		endLine = j
		// Statement terminators
		upper := strings.ToUpper(trimmed)
		if strings.HasSuffix(trimmed, ";") ||
			upper == "END;" || upper == "END" ||
			strings.HasPrefix(upper, "END;") ||
			strings.HasPrefix(upper, "END ") {
			break
		}
	}
	return endLine
}
