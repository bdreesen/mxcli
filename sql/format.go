// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// FormatTable writes a QueryResult as a pipe-delimited table to w.
func FormatTable(w io.Writer, result *QueryResult) {
	if len(result.Columns) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(result.Columns))
	for i, col := range result.Columns {
		widths[i] = len(col)
	}
	for _, row := range result.Rows {
		for i, val := range row {
			s := formatValue(val)
			if len(s) > widths[i] {
				widths[i] = len(s)
			}
		}
	}

	// Cap column widths at 50 characters
	for i := range widths {
		if widths[i] > 50 {
			widths[i] = 50
		}
	}

	// Print header
	fmt.Fprint(w, "|")
	for i, col := range result.Columns {
		fmt.Fprintf(w, " %-*s |", widths[i], truncate(col, widths[i]))
	}
	fmt.Fprintln(w)

	// Print separator
	fmt.Fprint(w, "|")
	for _, wid := range widths {
		fmt.Fprintf(w, "-%s-|", strings.Repeat("-", wid))
	}
	fmt.Fprintln(w)

	// Print rows
	for _, row := range result.Rows {
		fmt.Fprint(w, "|")
		for i, val := range row {
			s := formatValue(val)
			fmt.Fprintf(w, " %-*s |", widths[i], truncate(s, widths[i]))
		}
		fmt.Fprintln(w)
	}
}

// FormatJSON writes a QueryResult as a JSON array of objects to w.
func FormatJSON(w io.Writer, result *QueryResult) error {
	objects := make([]map[string]any, 0, len(result.Rows))
	for _, row := range result.Rows {
		obj := make(map[string]any, len(result.Columns))
		for i, col := range result.Columns {
			obj[col] = row[i]
		}
		objects = append(objects, obj)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(objects)
}

// formatValue formats a value for table display.
func formatValue(val any) string {
	if val == nil {
		return "NULL"
	}
	s := fmt.Sprintf("%v", val)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}

// truncate truncates a string to max length with ellipsis.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
