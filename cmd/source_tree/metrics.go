// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// computeChurn counts how many commits have touched each file over its full git history.
// Uses a single `git log --name-only` call for efficiency.
func computeChurn(root string, files []fileInfo) {
	cmd := exec.Command("git", "log", "--format=", "--name-only")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return
	}

	// Count occurrences of each file path
	counts := map[string]int{}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			counts[line]++
		}
	}

	// Map counts to files
	for i := range files {
		if n, ok := counts[files[i].path]; ok {
			files[i].churn = n
		}
	}
}

// computeFnMetrics computes function count and max function length for each file.
// For Go files: tracks `func` declarations at brace depth 0.
// For TypeScript files: tracks `function`, `export function`, and arrow functions.
func computeFnMetrics(root string, files []fileInfo) {
	for i := range files {
		data, err := os.ReadFile(filepath.Join(root, files[i].path))
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		isTS := strings.HasSuffix(files[i].path, ".ts")

		var fnLengths []int
		braceDepth := 0
		fnStartDepth := -1 // brace depth when current function started
		fnStartLine := 0

		for j, line := range lines {
			trimmed := strings.TrimSpace(line)

			// Detect function start at top level (braceDepth == 0)
			isFnStart := false
			if braceDepth == 0 {
				if isTS {
					isFnStart = strings.HasPrefix(trimmed, "function ") ||
						strings.HasPrefix(trimmed, "export function ") ||
						strings.HasPrefix(trimmed, "export async function ") ||
						strings.HasPrefix(trimmed, "async function ")
				} else {
					isFnStart = strings.HasPrefix(trimmed, "func ")
				}
			}

			// Count braces (simple: not inside strings/comments, good enough for metrics)
			for _, ch := range trimmed {
				if ch == '{' {
					braceDepth++
				} else if ch == '}' {
					braceDepth--
				}
			}

			if isFnStart {
				fnStartDepth = 0
				fnStartLine = j
			}

			// Function ends when we return to the depth where it started
			if fnStartDepth >= 0 && braceDepth <= fnStartDepth {
				fnLen := j - fnStartLine + 1
				fnLengths = append(fnLengths, fnLen)
				fnStartDepth = -1
			}
		}

		if len(fnLengths) == 0 {
			continue
		}

		maxLen := 0
		totalLen := 0
		for _, l := range fnLengths {
			totalLen += l
			if l > maxLen {
				maxLen = l
			}
		}

		files[i].fnData = &fnMetrics{
			count:  len(fnLengths),
			maxLen: maxLen,
			avgLen: totalLen / len(fnLengths),
		}
	}
}

// computeDuplication detects intra-file code duplication for each file.
// It uses a sliding window of 4 normalized lines as fingerprints.
// Lines that participate in duplicated blocks are counted toward the dup percentage.
func computeDuplication(root string, files []fileInfo) {
	const windowSize = 4

	for i := range files {
		data, err := os.ReadFile(filepath.Join(root, files[i].path))
		if err != nil {
			continue
		}

		// Normalize lines: trim whitespace, collect non-empty/non-comment lines
		// with their original line indices
		rawLines := strings.Split(string(data), "\n")
		type indexedLine struct {
			idx  int
			text string
		}
		var normLines []indexedLine
		for j, line := range rawLines {
			trimmed := strings.TrimSpace(line)
			// Skip blank lines, single braces, and comment-only lines
			if trimmed == "" || trimmed == "{" || trimmed == "}" || trimmed == "}," ||
				trimmed == ")" || trimmed == ")," ||
				strings.HasPrefix(trimmed, "//") {
				continue
			}
			normLines = append(normLines, indexedLine{idx: j, text: trimmed})
		}

		if len(normLines) < windowSize {
			continue
		}

		// Build fingerprints from sliding windows of normalized lines
		type fingerprint string
		seen := map[fingerprint][]int{} // fingerprint -> list of starting indices in normLines

		for j := 0; j <= len(normLines)-windowSize; j++ {
			var buf strings.Builder
			for k := range windowSize {
				if k > 0 {
					buf.WriteByte('\n')
				}
				buf.WriteString(normLines[j+k].text)
			}
			fp := fingerprint(buf.String())
			seen[fp] = append(seen[fp], j)
		}

		// Mark lines that participate in any duplicated fingerprint
		dupLines := map[int]bool{}
		for _, positions := range seen {
			if len(positions) < 2 {
				continue
			}
			for _, pos := range positions {
				for k := range windowSize {
					dupLines[normLines[pos+k].idx] = true
				}
			}
		}

		total := len(rawLines)
		if total == 0 {
			continue
		}
		files[i].dupPct = len(dupLines) * 100 / total
	}
}
