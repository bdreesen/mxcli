// SPDX-License-Identifier: Apache-2.0

package evalrunner

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ReportOptions configures report output.
type ReportOptions struct {
	// OutputDir is the directory for JSON/Markdown reports. Empty = no file output.
	OutputDir string

	// Color enables colored console output.
	Color bool
}

// PrintResult writes a human-readable summary of an eval result to the writer.
func PrintResult(w io.Writer, result *EvalResult, color bool) {
	fmt.Fprintf(w, "\nEval: %s (%s)", result.TestID, result.Category)
	if result.Title != "" {
		fmt.Fprintf(w, " — %s", result.Title)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("=", 60))

	// Initial phase
	fmt.Fprintln(w, "  Initial:")
	printPhaseResults(w, &result.Initial, color)

	// Iteration phase
	if result.Iteration != nil {
		fmt.Fprintln(w, "  Iteration:")
		printPhaseResults(w, result.Iteration, color)
	}

	// Overall
	fmt.Fprintln(w, strings.Repeat("-", 60))
	scoreStr := formatScore(result.OverallScore)
	fmt.Fprintf(w, "  Overall: %d/%d (%s)\n",
		result.TotalPassed(), result.TotalChecks(), scoreStr)

	if result.Duration > 0 {
		fmt.Fprintf(w, "  Duration: %s\n", result.Duration.Round(time.Millisecond))
	}

	fmt.Fprintln(w)
}

// printPhaseResults writes check results for a single phase.
func printPhaseResults(w io.Writer, phase *PhaseResult, color bool) {
	for _, cr := range phase.Checks {
		status := "PASS"
		if !cr.Passed {
			status = "FAIL"
		}

		if color {
			if cr.Passed {
				status = "\033[32mPASS\033[0m"
			} else {
				status = "\033[31mFAIL\033[0m"
			}
		}

		fmt.Fprintf(w, "    [%s] %s", status, cr.Check)
		if cr.Detail != "" && !cr.Passed {
			fmt.Fprintf(w, " — %s", cr.Detail)
		}
		fmt.Fprintln(w)
	}

	scoreStr := formatScore(phase.Score)
	fmt.Fprintf(w, "    Score: %d/%d (%s)\n", phase.Passed, phase.Total, scoreStr)
}

// formatScore formats a score as a percentage string.
func formatScore(score float64) string {
	return fmt.Sprintf("%.0f%%", score*100)
}

// PrintSummary writes a summary table for multiple eval results.
func PrintSummary(w io.Writer, summary *RunSummary, color bool) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Eval Run Summary — %s\n", summary.Timestamp.Format("2006-01-02 15:04"))
	fmt.Fprintln(w, strings.Repeat("=", 70))

	// Header
	fmt.Fprintf(w, "  %-12s %-15s %8s %8s %10s\n", "Test", "Category", "Score", "Checks", "Iteration")
	fmt.Fprintln(w, "  "+strings.Repeat("-", 58))

	for _, r := range summary.Results {
		iterScore := "—"
		if r.Iteration != nil {
			iterScore = formatScore(r.Iteration.Score)
		}

		fmt.Fprintf(w, "  %-12s %-15s %8s %5d/%d %10s\n",
			r.TestID,
			r.Category,
			formatScore(r.OverallScore),
			r.TotalPassed(), r.TotalChecks(),
			iterScore,
		)
	}

	fmt.Fprintln(w, "  "+strings.Repeat("-", 58))
	fmt.Fprintf(w, "  Average Score: %s\n", formatScore(summary.AverageScore()))
	if summary.Duration > 0 {
		fmt.Fprintf(w, "  Total Duration: %s\n", summary.Duration.Round(time.Second))
	}
	fmt.Fprintln(w)
}

// WriteJSONReport writes the eval result as JSON to the output directory.
func WriteJSONReport(result *EvalResult, outputDir string) error {
	dir := filepath.Join(outputDir, result.TestID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling result: %w", err)
	}

	path := filepath.Join(dir, "score.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	return nil
}

// WriteRunSummary writes the full run summary as JSON.
func WriteRunSummary(summary *RunSummary, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling summary: %w", err)
	}

	path := filepath.Join(outputDir, "summary.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	return nil
}

// WriteMarkdownReport writes a Markdown summary of the run.
func WriteMarkdownReport(summary *RunSummary, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("# Eval Run %s\n\n", summary.Timestamp.Format("2006-01-02 15:04")))
	buf.WriteString(fmt.Sprintf("Tests: %d | Duration: %s | Average Score: %s\n\n",
		len(summary.Results),
		summary.Duration.Round(time.Second),
		formatScore(summary.AverageScore()),
	))

	// Summary table
	buf.WriteString("| Test | Category | Score | Checks | Iteration |\n")
	buf.WriteString("|------|----------|-------|--------|-----------|\n")

	for _, r := range summary.Results {
		iterScore := "—"
		if r.Iteration != nil {
			iterScore = formatScore(r.Iteration.Score)
		}

		buf.WriteString(fmt.Sprintf("| %s | %s | %s | %d/%d | %s |\n",
			r.TestID,
			r.Category,
			formatScore(r.OverallScore),
			r.TotalPassed(), r.TotalChecks(),
			iterScore,
		))
	}

	buf.WriteString("\n")

	// Detailed results per test
	for _, r := range summary.Results {
		buf.WriteString(fmt.Sprintf("## %s: %s\n\n", r.TestID, r.Title))

		buf.WriteString("### Initial Checks\n\n")
		for _, cr := range r.Initial.Checks {
			status := "PASS"
			if !cr.Passed {
				status = "FAIL"
			}
			buf.WriteString(fmt.Sprintf("- [%s] `%s`", status, cr.Check))
			if cr.Detail != "" {
				buf.WriteString(fmt.Sprintf(" — %s", cr.Detail))
			}
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("\nScore: %s (%d/%d)\n\n",
			formatScore(r.Initial.Score), r.Initial.Passed, r.Initial.Total))

		if r.Iteration != nil {
			buf.WriteString("### Iteration Checks\n\n")
			for _, cr := range r.Iteration.Checks {
				status := "PASS"
				if !cr.Passed {
					status = "FAIL"
				}
				buf.WriteString(fmt.Sprintf("- [%s] `%s`", status, cr.Check))
				if cr.Detail != "" {
					buf.WriteString(fmt.Sprintf(" — %s", cr.Detail))
				}
				buf.WriteString("\n")
			}
			buf.WriteString(fmt.Sprintf("\nScore: %s (%d/%d)\n\n",
				formatScore(r.Iteration.Score), r.Iteration.Passed, r.Iteration.Total))
		}
	}

	path := filepath.Join(outputDir, "summary.md")
	if err := os.WriteFile(path, []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	return nil
}
