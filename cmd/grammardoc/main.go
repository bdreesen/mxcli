// SPDX-License-Identifier: Apache-2.0

// Command grammardoc generates documentation from ANTLR4 grammar files.
//
// It parses grammar files, extracts documentation comments and rules,
// and generates markdown documentation with examples and railroad diagrams.
//
// Usage:
//
//	go run ./cmd/grammardoc -grammar mdl/grammar/MDLParser.g4 -output docs/mdl-reference.md
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Rule represents a parsed grammar rule with its documentation
type Rule struct {
	Name        string
	Definition  string
	Description string
	Examples    []Example
	SeeAlso     []string
	Since       string
	Deprecated  string
}

// Example represents a code example from documentation
type Example struct {
	Title string
	Code  string
}

func main() {
	grammarFile := flag.String("grammar", "mdl/grammar/MDLParser.g4", "Path to ANTLR4 grammar file")
	lexerFile := flag.String("lexer", "mdl/grammar/MDLLexer.g4", "Path to ANTLR4 lexer file (optional)")
	outputFile := flag.String("output", "docs/06-mdl-reference/grammar-reference.md", "Output markdown file")
	title := flag.String("title", "MDL Grammar Reference", "Document title")
	flag.Parse()

	// Parse grammar file
	rules, err := parseGrammarFile(*grammarFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing grammar: %v\n", err)
		os.Exit(1)
	}

	// Parse lexer file for tokens (optional)
	var tokens []Rule
	if *lexerFile != "" {
		if _, err := os.Stat(*lexerFile); err == nil {
			tokens, _ = parseGrammarFile(*lexerFile)
		}
	}

	// Generate markdown
	markdown := generateMarkdown(*title, rules, tokens)

	// Write output
	if err := os.WriteFile(*outputFile, []byte(markdown), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s with %d rules\n", *outputFile, len(rules))
}

// parseGrammarFile parses an ANTLR4 grammar file and extracts rules with documentation
func parseGrammarFile(filename string) ([]Rule, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rules []Rule
	var currentDoc strings.Builder
	var currentRule *Rule
	var ruleDefBuilder strings.Builder
	inDocComment := false
	inRule := false

	scanner := bufio.NewScanner(file)
	// Increase buffer size for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Start of doc comment
		if strings.HasPrefix(trimmed, "/**") {
			inDocComment = true
			currentDoc.Reset()
			// Check if it's a single-line doc comment
			if strings.HasSuffix(trimmed, "*/") {
				content := strings.TrimPrefix(trimmed, "/**")
				content = strings.TrimSuffix(content, "*/")
				currentDoc.WriteString(strings.TrimSpace(content))
				inDocComment = false
			} else {
				content := strings.TrimPrefix(trimmed, "/**")
				if content != "" {
					currentDoc.WriteString(strings.TrimSpace(content))
					currentDoc.WriteString("\n")
				}
			}
			continue
		}

		// Inside doc comment
		if inDocComment {
			if before, ok := strings.CutSuffix(trimmed, "*/"); ok {
				content := before
				content = strings.TrimPrefix(content, "*")
				currentDoc.WriteString(strings.TrimSpace(content))
				inDocComment = false
			} else {
				content := strings.TrimPrefix(trimmed, "*")
				currentDoc.WriteString(strings.TrimSpace(content))
				currentDoc.WriteString("\n")
			}
			continue
		}

		// Skip regular comments
		if strings.HasPrefix(trimmed, "//") {
			continue
		}

		// Rule definition start (identifier followed by content or newline, then colon on same/next line)
		ruleMatch := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*$`).FindStringSubmatch(trimmed)
		if ruleMatch != nil {
			// This might be a rule name on its own line
			if currentRule != nil && inRule {
				// Save previous rule
				currentRule.Definition = strings.TrimSpace(ruleDefBuilder.String())
				rules = append(rules, *currentRule)
			}
			currentRule = &Rule{Name: ruleMatch[1]}
			parseDocComment(currentDoc.String(), currentRule)
			currentDoc.Reset()
			ruleDefBuilder.Reset()
			inRule = true
			continue
		}

		// Rule definition with colon on same line
		ruleMatch = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*:\s*(.*)$`).FindStringSubmatch(trimmed)
		if ruleMatch != nil {
			if currentRule != nil && inRule {
				// Save previous rule
				currentRule.Definition = strings.TrimSpace(ruleDefBuilder.String())
				rules = append(rules, *currentRule)
			}
			currentRule = &Rule{Name: ruleMatch[1]}
			parseDocComment(currentDoc.String(), currentRule)
			currentDoc.Reset()
			ruleDefBuilder.Reset()
			ruleDefBuilder.WriteString(ruleMatch[2])
			ruleDefBuilder.WriteString("\n")
			inRule = true
			continue
		}

		// Continue rule definition
		if inRule && currentRule != nil {
			if trimmed == ";" {
				// End of rule
				currentRule.Definition = strings.TrimSpace(ruleDefBuilder.String())
				rules = append(rules, *currentRule)
				currentRule = nil
				inRule = false
				ruleDefBuilder.Reset()
			} else if after, ok := strings.CutPrefix(trimmed, ":"); ok {
				// Colon on its own line (after rule name)
				ruleDefBuilder.WriteString(after)
				ruleDefBuilder.WriteString("\n")
			} else {
				ruleDefBuilder.WriteString(trimmed)
				ruleDefBuilder.WriteString("\n")
			}
		}
	}

	// Don't forget the last rule
	if currentRule != nil && inRule {
		currentRule.Definition = strings.TrimSpace(ruleDefBuilder.String())
		rules = append(rules, *currentRule)
	}

	return rules, scanner.Err()
}

// parseDocComment extracts structured information from a documentation comment
func parseDocComment(doc string, rule *Rule) {
	if doc == "" {
		return
	}

	lines := strings.Split(doc, "\n")
	var descLines []string
	var currentExample *Example
	inExample := false
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for @tags
		if strings.HasPrefix(trimmed, "@example") {
			if currentExample != nil {
				rule.Examples = append(rule.Examples, *currentExample)
			}
			title := strings.TrimSpace(strings.TrimPrefix(trimmed, "@example"))
			currentExample = &Example{Title: title}
			inExample = true
			continue
		}

		if after, ok := strings.CutPrefix(trimmed, "@see"); ok {
			ref := strings.TrimSpace(after)
			rule.SeeAlso = append(rule.SeeAlso, ref)
			continue
		}

		if after, ok := strings.CutPrefix(trimmed, "@since"); ok {
			rule.Since = strings.TrimSpace(after)
			continue
		}

		if after, ok := strings.CutPrefix(trimmed, "@deprecated"); ok {
			rule.Deprecated = strings.TrimSpace(after)
			continue
		}

		// Handle code blocks in examples
		if inExample && currentExample != nil {
			if strings.HasPrefix(trimmed, "```") {
				inCodeBlock = !inCodeBlock
				if !inCodeBlock {
					// End of code block
					continue
				}
				continue
			}
			if inCodeBlock {
				currentExample.Code += line + "\n"
			}
			continue
		}

		// Regular description line
		if !inExample {
			descLines = append(descLines, trimmed)
		}
	}

	// Save last example
	if currentExample != nil {
		rule.Examples = append(rule.Examples, *currentExample)
	}

	rule.Description = strings.TrimSpace(strings.Join(descLines, " "))
}

// generateMarkdown creates the markdown documentation
func generateMarkdown(title string, rules []Rule, _ []Rule) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("> Auto-generated from ANTLR4 grammar on %s\n\n", time.Now().Format("2006-01-02")))

	// Introduction
	sb.WriteString("This document provides a complete reference for the MDL (Mendix Definition Language) syntax.\n")
	sb.WriteString("Each grammar rule is documented with its syntax, description, and examples.\n\n")

	// Table of contents
	sb.WriteString("## Table of Contents\n\n")

	// Group rules by category
	categories := categorizeRules(rules)
	for _, cat := range categories {
		sb.WriteString(fmt.Sprintf("- [%s](#%s)\n", cat.Name, toAnchor(cat.Name)))
		for _, rule := range cat.Rules {
			if rule.Description != "" || len(rule.Examples) > 0 {
				sb.WriteString(fmt.Sprintf("  - [%s](#%s)\n", rule.Name, toAnchor(rule.Name)))
			}
		}
	}
	sb.WriteString("\n")

	// Generate rule documentation by category
	for _, cat := range categories {
		sb.WriteString(fmt.Sprintf("## %s\n\n", cat.Name))

		for _, rule := range cat.Rules {
			// Skip rules without documentation unless they have a definition
			if rule.Description == "" && len(rule.Examples) == 0 && rule.Definition == "" {
				continue
			}

			sb.WriteString(fmt.Sprintf("### %s\n\n", rule.Name))

			// Description
			if rule.Description != "" {
				sb.WriteString(rule.Description + "\n\n")
			}

			// Syntax (EBNF and railroad diagram)
			if rule.Definition != "" {
				sb.WriteString("**Syntax:**\n\n")
				sb.WriteString("```ebnf\n")
				sb.WriteString(formatDefinition(rule.Name, rule.Definition))
				sb.WriteString("\n```\n\n")

				// Embedded Mermaid state diagram as railroad visualization
				sb.WriteString("**Railroad Diagram:**\n\n")
				sb.WriteString(generateMermaidDiagram(rule.Name, rule.Definition))
				sb.WriteString("\n")
			}

			// Examples
			if len(rule.Examples) > 0 {
				sb.WriteString("**Examples:**\n\n")
				for _, ex := range rule.Examples {
					if ex.Title != "" {
						sb.WriteString(fmt.Sprintf("*%s:*\n\n", ex.Title))
					}
					sb.WriteString("```sql\n")
					sb.WriteString(strings.TrimSpace(ex.Code))
					sb.WriteString("\n```\n\n")
				}
			}

			// See also
			if len(rule.SeeAlso) > 0 {
				sb.WriteString("**See also:** ")
				refs := make([]string, len(rule.SeeAlso))
				for i, ref := range rule.SeeAlso {
					refs[i] = fmt.Sprintf("[%s](#%s)", ref, toAnchor(ref))
				}
				sb.WriteString(strings.Join(refs, ", ") + "\n\n")
			}

			// Since/Deprecated
			if rule.Since != "" {
				sb.WriteString(fmt.Sprintf("*Since: %s*\n\n", rule.Since))
			}
			if rule.Deprecated != "" {
				sb.WriteString(fmt.Sprintf("⚠️ **Deprecated:** %s\n\n", rule.Deprecated))
			}

			sb.WriteString("---\n\n")
		}
	}

	return sb.String()
}

// Category groups related rules
type Category struct {
	Name  string
	Rules []Rule
}

// categorizeRules groups rules into logical categories
func categorizeRules(rules []Rule) []Category {
	categories := []Category{
		{Name: "Statements", Rules: []Rule{}},
		{Name: "Entity Definitions", Rules: []Rule{}},
		{Name: "Microflow Statements", Rules: []Rule{}},
		{Name: "Page Definitions", Rules: []Rule{}},
		{Name: "OQL Queries", Rules: []Rule{}},
		{Name: "Expressions", Rules: []Rule{}},
		{Name: "Other Rules", Rules: []Rule{}},
	}

	for _, rule := range rules {
		name := strings.ToLower(rule.Name)
		switch {
		case strings.Contains(name, "statement") || strings.Contains(name, "program"):
			categories[0].Rules = append(categories[0].Rules, rule)
		case strings.Contains(name, "entity") || strings.Contains(name, "attribute") || strings.Contains(name, "association"):
			categories[1].Rules = append(categories[1].Rules, rule)
		case strings.Contains(name, "microflow") || strings.Contains(name, "action") || strings.Contains(name, "flow"):
			categories[2].Rules = append(categories[2].Rules, rule)
		case strings.Contains(name, "page") || strings.Contains(name, "widget") || strings.Contains(name, "layout"):
			categories[3].Rules = append(categories[3].Rules, rule)
		case strings.Contains(name, "oql") || strings.Contains(name, "select") || strings.Contains(name, "from") || strings.Contains(name, "where"):
			categories[4].Rules = append(categories[4].Rules, rule)
		case strings.Contains(name, "expr") || strings.Contains(name, "literal") || strings.Contains(name, "operator"):
			categories[5].Rules = append(categories[5].Rules, rule)
		default:
			categories[6].Rules = append(categories[6].Rules, rule)
		}
	}

	// Filter out empty categories
	var nonEmpty []Category
	for _, cat := range categories {
		if len(cat.Rules) > 0 {
			nonEmpty = append(nonEmpty, cat)
		}
	}

	return nonEmpty
}

// formatDefinition formats a rule definition for display
func formatDefinition(name, def string) string {
	// Clean up the definition
	lines := strings.Split(def, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && trimmed != ";" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return name + "\n    : " + strings.Join(cleaned, "\n    | ")
}

// generateMermaidDiagram creates a Mermaid state diagram for a grammar rule
func generateMermaidDiagram(name, def string) string {
	var sb strings.Builder
	sb.WriteString("```mermaid\n")
	sb.WriteString("stateDiagram-v2\n")
	sb.WriteString("    direction LR\n")

	// Parse the definition into alternatives
	lines := strings.Split(def, "\n")
	var alternatives []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && trimmed != ";" && !strings.HasPrefix(trimmed, "//") {
			// Remove leading | if present
			trimmed = strings.TrimPrefix(trimmed, "|")
			trimmed = strings.TrimSpace(trimmed)
			if trimmed != "" {
				alternatives = append(alternatives, trimmed)
			}
		}
	}

	// Generate unique state IDs
	stateID := 0
	nextID := func() string {
		stateID++
		return fmt.Sprintf("s%d", stateID)
	}

	if len(alternatives) == 0 {
		sb.WriteString("    [*] --> [*]\n")
	} else if len(alternatives) == 1 {
		// Single path
		elements := tokenizeRule(alternatives[0])
		writeStatePath(&sb, elements, nextID)
	} else {
		// Multiple alternatives - each starts from [*] and ends at [*]
		for _, alt := range alternatives {
			elements := tokenizeRule(alt)
			if len(elements) == 0 {
				sb.WriteString("    [*] --> [*]\n")
			} else {
				writeStatePath(&sb, elements, nextID)
			}
		}
	}

	sb.WriteString("```\n")
	return sb.String()
}

// tokenizeRule splits a rule alternative into tokens
func tokenizeRule(rule string) []string {
	// First strip any end-of-line comments
	if idx := strings.Index(rule, "//"); idx >= 0 {
		rule = strings.TrimSpace(rule[:idx])
	}

	var tokens []string
	// Simple tokenization - split on whitespace but keep parenthesized groups
	var current strings.Builder
	depth := 0
	inQuote := false

	for _, ch := range rule {
		switch ch {
		case '\'':
			inQuote = !inQuote
			current.WriteRune(ch)
		case '(':
			depth++
			current.WriteRune(ch)
		case ')':
			depth--
			current.WriteRune(ch)
		case ' ', '\t':
			if depth > 0 || inQuote {
				current.WriteRune(ch)
			} else if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// writeStatePath writes a single path in the state diagram
func writeStatePath(sb *strings.Builder, elements []string, nextID func() string) {
	if len(elements) == 0 {
		sb.WriteString("    [*] --> [*]\n")
		return
	}

	// First element connects from [*]
	firstID := nextID()
	sb.WriteString(fmt.Sprintf("    state \"%s\" as %s\n", escapeStateLabel(elements[0]), firstID))
	sb.WriteString(fmt.Sprintf("    [*] --> %s\n", firstID))

	prevID := firstID
	for i := 1; i < len(elements); i++ {
		currID := nextID()
		sb.WriteString(fmt.Sprintf("    state \"%s\" as %s\n", escapeStateLabel(elements[i]), currID))
		sb.WriteString(fmt.Sprintf("    %s --> %s\n", prevID, currID))
		prevID = currID
	}

	// Last element connects to [*]
	sb.WriteString(fmt.Sprintf("    %s --> [*]\n", prevID))
}

// escapeStateLabel escapes special characters for Mermaid state labels
func escapeStateLabel(s string) string {
	// Escape characters that break Mermaid state labels
	s = strings.ReplaceAll(s, "\"", "'")
	s = strings.ReplaceAll(s, "<", "‹")
	s = strings.ReplaceAll(s, ">", "›")
	return s
}

// toAnchor converts a string to a markdown anchor
func toAnchor(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(s, "")
	return s
}
