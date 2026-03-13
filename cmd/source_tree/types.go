// SPDX-License-Identifier: Apache-2.0

package main

// fileInfo holds line count and relative path for a source file.
type fileInfo struct {
	lines  int
	path   string     // relative, e.g. "sdk/mpr/reader.go"
	dupPct int        // percentage of lines in duplicated blocks (0-100)
	fnData *fnMetrics // function metrics (nil if not computed)
	churn  int        // number of commits that touched this file (0 = not computed)
}

// fnMetrics holds per-file function metrics.
type fnMetrics struct {
	count  int // number of functions
	maxLen int // longest function in lines
	avgLen int // average function length
}

// covData holds per-file coverage data.
type covData struct {
	pct     int
	total   int
	covered int
}

// ANSI color codes.
const (
	codeBold   = "\033[1m"
	codeDim    = "\033[2m"
	codeCyan   = "\033[36m"
	codeGreen  = "\033[32m"
	codeYellow = "\033[33m"
	codeOrange = "\033[38;5;208m"
	codeRed    = "\033[31m"
	codeWhite  = "\033[37m"
	codeReset  = "\033[0m"
)

const barMax = 40
const covBarMax = 10
const dupBarMax = 10

var useColor bool

func c(code string) string {
	if useColor {
		return code
	}
	return ""
}
