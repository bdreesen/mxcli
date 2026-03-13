// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"strings"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// FoldingRanges handles textDocument/foldingRange requests.
func (s *mdlServer) FoldingRanges(ctx context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	docURI := uri.URI(params.TextDocument.URI)
	s.mu.Lock()
	text := s.docs[docURI]
	s.mu.Unlock()

	if text == "" {
		return nil, nil
	}

	return extractFoldingRanges(text), nil
}

// extractFoldingRanges scans document text for foldable regions:
// - BEGIN ... END blocks (microflow/nanoflow bodies)
// - IF ... END IF blocks
// - LOOP ... END LOOP blocks
// - { ... } brace blocks (page bodies, widget containers)
// - ( ... ) paren blocks (entity definitions, parameter lists)
// - /* ... */ block comments
// - Consecutive -- line comments
func extractFoldingRanges(text string) []protocol.FoldingRange {
	lines := strings.Split(text, "\n")
	var ranges []protocol.FoldingRange

	// Stack-based tracking for nested blocks
	type stackEntry struct {
		kind      string // "begin", "if", "loop", "brace", "paren", "comment"
		startLine int
	}
	var stack []stackEntry

	// Track consecutive line comments
	lineCommentStart := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		upper := strings.ToUpper(trimmed)

		// Handle consecutive line comments
		isLineComment := strings.HasPrefix(trimmed, "--")
		if isLineComment {
			if lineCommentStart == -1 {
				lineCommentStart = i
			}
		} else {
			if lineCommentStart != -1 && i-1 > lineCommentStart {
				// At least 2 consecutive comment lines
				ranges = append(ranges, protocol.FoldingRange{
					StartLine: uint32(lineCommentStart),
					EndLine:   uint32(i - 1),
					Kind:      protocol.CommentFoldingRange,
				})
			}
			lineCommentStart = -1
		}

		if trimmed == "" {
			continue
		}

		// Block comments: /* ... */
		if strings.Contains(trimmed, "/*") && !strings.Contains(trimmed, "*/") {
			stack = append(stack, stackEntry{"comment", i})
		}
		if strings.Contains(trimmed, "*/") {
			for j := len(stack) - 1; j >= 0; j-- {
				if stack[j].kind == "comment" {
					if i > stack[j].startLine {
						ranges = append(ranges, protocol.FoldingRange{
							StartLine: uint32(stack[j].startLine),
							EndLine:   uint32(i),
							Kind:      protocol.CommentFoldingRange,
						})
					}
					stack = append(stack[:j], stack[j+1:]...)
					break
				}
			}
		}

		// Skip further keyword processing inside block comments
		inBlockComment := false
		for _, entry := range stack {
			if entry.kind == "comment" {
				inBlockComment = true
				break
			}
		}
		if inBlockComment || isLineComment {
			continue
		}

		// BEGIN ... END (microflow bodies)
		if upper == "BEGIN" || strings.HasSuffix(upper, " BEGIN") {
			stack = append(stack, stackEntry{"begin", i})
		}

		// IF ... THEN
		if strings.HasPrefix(upper, "IF ") && strings.HasSuffix(upper, "THEN") {
			stack = append(stack, stackEntry{"if", i})
		}

		// LOOP ... BEGIN
		if strings.HasPrefix(upper, "LOOP ") {
			stack = append(stack, stackEntry{"loop", i})
		}

		// END IF — also removes any nested entries inside the if block
		if upper == "END IF;" || upper == "END IF" {
			for j := len(stack) - 1; j >= 0; j-- {
				if stack[j].kind == "if" {
					if i > stack[j].startLine {
						ranges = append(ranges, protocol.FoldingRange{
							StartLine: uint32(stack[j].startLine),
							EndLine:   uint32(i),
							Kind:      protocol.RegionFoldingRange,
						})
					}
					stack = stack[:j]
					break
				}
			}
		}

		// END LOOP — also removes any nested begin entries inside the loop
		if upper == "END LOOP;" || upper == "END LOOP" {
			for j := len(stack) - 1; j >= 0; j-- {
				if stack[j].kind == "loop" {
					if i > stack[j].startLine {
						ranges = append(ranges, protocol.FoldingRange{
							StartLine: uint32(stack[j].startLine),
							EndLine:   uint32(i),
							Kind:      protocol.RegionFoldingRange,
						})
					}
					// Remove the loop entry and everything above it (nested begin, etc.)
					stack = stack[:j]
					break
				}
			}
		}

		// END; or END (for BEGIN blocks, but not END IF or END LOOP)
		if (upper == "END;" || upper == "END") &&
			!strings.HasPrefix(upper, "END IF") &&
			!strings.HasPrefix(upper, "END LOOP") {
			for j := len(stack) - 1; j >= 0; j-- {
				if stack[j].kind == "begin" {
					if i > stack[j].startLine {
						ranges = append(ranges, protocol.FoldingRange{
							StartLine: uint32(stack[j].startLine),
							EndLine:   uint32(i),
							Kind:      protocol.RegionFoldingRange,
						})
					}
					stack = append(stack[:j], stack[j+1:]...)
					break
				}
			}
		}

		// Track braces and parens character by character
		inString := false
		for _, ch := range trimmed {
			if ch == '\'' {
				inString = !inString
				continue
			}
			if inString {
				continue
			}
			switch ch {
			case '{':
				stack = append(stack, stackEntry{"brace", i})
			case '}':
				for j := len(stack) - 1; j >= 0; j-- {
					if stack[j].kind == "brace" {
						if i > stack[j].startLine {
							ranges = append(ranges, protocol.FoldingRange{
								StartLine: uint32(stack[j].startLine),
								EndLine:   uint32(i),
								Kind:      protocol.RegionFoldingRange,
							})
						}
						stack = append(stack[:j], stack[j+1:]...)
						break
					}
				}
			case '(':
				stack = append(stack, stackEntry{"paren", i})
			case ')':
				for j := len(stack) - 1; j >= 0; j-- {
					if stack[j].kind == "paren" {
						if i > stack[j].startLine {
							ranges = append(ranges, protocol.FoldingRange{
								StartLine: uint32(stack[j].startLine),
								EndLine:   uint32(i),
								Kind:      protocol.RegionFoldingRange,
							})
						}
						stack = append(stack[:j], stack[j+1:]...)
						break
					}
				}
			}
		}
	}

	// Close any remaining consecutive line comments at end of file
	if lineCommentStart != -1 && len(lines)-1 > lineCommentStart {
		ranges = append(ranges, protocol.FoldingRange{
			StartLine: uint32(lineCommentStart),
			EndLine:   uint32(len(lines) - 1),
			Kind:      protocol.CommentFoldingRange,
		})
	}

	return ranges
}
