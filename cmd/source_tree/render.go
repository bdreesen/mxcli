// SPDX-License-Identifier: Apache-2.0

package main

import (
	"path/filepath"
	"sort"
	"strings"
)

// filesInDir returns files directly in the given directory, sorted by lines descending.
func filesInDir(files []fileInfo, dir string) []fileInfo {
	var result []fileInfo
	for _, f := range files {
		fDir := filepath.ToSlash(filepath.Dir(f.path))
		if dir == "." {
			// Root: files with no directory separator
			if !strings.Contains(f.path, "/") {
				result = append(result, f)
			}
		} else {
			if fDir == dir {
				result = append(result, f)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].lines != result[j].lines {
			return result[i].lines > result[j].lines
		}
		return result[i].path < result[j].path
	})
	return result
}

func makeBar(lines, maxLines int) string {
	length := lines * barMax / maxLines
	var buf strings.Builder
	displayLen := 0

	for range length {
		buf.WriteString("█")
	}
	displayLen = length

	// Add half block if we'd round up
	remainder := (lines*barMax*2/maxLines - length*2)
	if remainder > 0 && length < barMax {
		buf.WriteString("▌")
		displayLen++
	}

	// Pad to fixed width
	for i := displayLen; i < barMax; i++ {
		buf.WriteByte(' ')
	}
	return buf.String()
}

func barColor(lines int) string {
	switch {
	case lines >= 1500:
		return codeRed
	case lines >= 1000:
		return codeOrange
	case lines >= 500:
		return codeYellow
	default:
		return codeGreen
	}
}

func covBarColor(pct int) string {
	switch {
	case pct >= 80:
		return codeGreen
	case pct >= 60:
		return codeYellow
	case pct >= 30:
		return codeOrange
	default:
		return codeRed
	}
}

func maxFnColor(lines int) string {
	switch {
	case lines >= 100:
		return codeRed
	case lines >= 60:
		return codeOrange
	case lines >= 30:
		return codeYellow
	default:
		return codeGreen
	}
}

func churnColor(commits int) string {
	switch {
	case commits >= 51:
		return codeRed
	case commits >= 26:
		return codeOrange
	case commits >= 11:
		return codeYellow
	default:
		return codeGreen
	}
}

func dupBarColor(pct int) string {
	switch {
	case pct >= 40:
		return codeRed
	case pct >= 25:
		return codeOrange
	case pct >= 10:
		return codeYellow
	default:
		return codeGreen
	}
}

func commitColor(changed int) string {
	switch {
	case changed >= 500:
		return codeRed
	case changed >= 200:
		return codeOrange
	case changed >= 50:
		return codeYellow
	default:
		return codeGreen
	}
}
