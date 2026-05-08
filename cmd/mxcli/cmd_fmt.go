// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/formatter"
	"github.com/mendixlabs/mxcli/mdl/visitor"
	"github.com/spf13/cobra"
)

var fmtCmd = &cobra.Command{
	Use:   "fmt <file.mdl>",
	Short: "Format an MDL file",
	Long: `Format an MDL script file with consistent styling:
  - Uppercase MDL keywords
  - Normalize indentation (2-space units)
  - Remove trailing whitespace
  - Normalize blank lines

Examples:
  # Format to stdout
  mxcli fmt script.mdl

  # Format in-place
  mxcli fmt script.mdl -w
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		writeInPlace, _ := cmd.Flags().GetBool("write")

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Reject unparseable input so automation scripts can detect failures.
		// Two failure modes:
		//   1. ANTLR reports explicit parse errors (structural violations).
		//   2. ANTLR silently skips unrecognised tokens — detected when no
		//      statements were produced from non-blank, non-comment content.
		prog, errs := visitor.Build(string(data))
		if len(errs) > 0 {
			var msgs []string
			for _, e := range errs {
				msgs = append(msgs, e.Error())
			}
			return fmt.Errorf("syntax errors in %s:\n%s", filePath, strings.Join(msgs, "\n"))
		}
		if prog != nil && len(prog.Statements) == 0 && hasSubstantiveContent(string(data)) {
			return fmt.Errorf("no valid MDL statements found in %s", filePath)
		}

		formatted := formatter.Format(string(data))

		if writeInPlace {
			if err := os.WriteFile(filePath, []byte(formatted), 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Formatted %s\n", filePath)
		} else {
			fmt.Print(formatted)
		}

		return nil
	},
}

func init() {
	fmtCmd.Flags().BoolP("write", "w", false, "Write result to source file instead of stdout")
}

// hasSubstantiveContent reports whether s contains at least one non-blank,
// non-comment line — used to distinguish empty/comment-only files (which
// produce zero statements legitimately) from garbage input.
func hasSubstantiveContent(s string) bool {
	for _, line := range strings.Split(s, "\n") {
		t := strings.TrimSpace(line)
		if t != "" && !strings.HasPrefix(t, "--") {
			return true
		}
	}
	return false
}
