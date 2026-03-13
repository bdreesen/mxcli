// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/mdl/executor"
	"github.com/mendixlabs/mxcli/mdl/visitor"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check <file>",
	Short: "Check an MDL script for errors without executing it",
	Long: `Check an MDL script file for syntax errors and optionally validate references.

By default, only checks syntax (parsing). Use --references to also validate
that all referenced modules, entities, etc. exist in the project.

Reference validation is smart: it automatically skips references to objects
that are created within the script itself. For example, if your script creates
a module "MyModule" and then creates entities in it, no error will be reported
for the module reference.

Examples:
  # Check syntax only (no project needed)
  mxcli check script.mdl

  # Check syntax and validate references against a project
  mxcli check script.mdl -p app.mpr --references
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		projectPath, _ := cmd.Flags().GetString("project")
		checkRefs, _ := cmd.Flags().GetBool("references")

		// Read the file
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		// Parse the script
		fmt.Printf("Checking syntax: %s\n", filePath)
		prog, errs := visitor.Build(string(content))
		if len(errs) > 0 {
			fmt.Fprintf(os.Stderr, "Syntax errors found:\n")
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "  - %v\n", err)
			}
			// Hint: if script contains IMPORT/QUERY with single $ but not $$, suggest dollar-quoting
			src := string(content)
			if (strings.Contains(src, "IMPORT") || strings.Contains(src, "import")) &&
				(strings.Contains(src, "QUERY") || strings.Contains(src, "query")) &&
				strings.Contains(src, "$") && !strings.Contains(src, "$$") {
				fmt.Fprintf(os.Stderr, "\nHint: SQL queries in IMPORT statements should use dollar-quoting ($$...$$) instead of single quotes.\n")
				fmt.Fprintf(os.Stderr, "  Example: IMPORT FROM alias QUERY $$SELECT * FROM table$$ INTO Module.Entity MAP (...)\n")
			}
			os.Exit(1)
		}
		fmt.Printf("✓ Syntax OK (%d statements)\n", len(prog.Statements))

		// Validate statements (doesn't require project connection)
		var oqlErrors []string
		for i, stmt := range prog.Statements {
			// Check enumeration values for reserved words
			if enumStmt, ok := stmt.(*ast.CreateEnumerationStmt); ok {
				if errs := executor.ValidateEnumeration(enumStmt); len(errs) > 0 {
					for _, e := range errs {
						oqlErrors = append(oqlErrors, fmt.Sprintf("statement %d (%s): %s", i+1, enumStmt.Name.String(), e))
					}
				}
			}
			// Check microflow body for common issues
			if mfStmt, ok := stmt.(*ast.CreateMicroflowStmt); ok {
				if warns := executor.ValidateMicroflow(mfStmt); len(warns) > 0 {
					for _, w := range warns {
						oqlErrors = append(oqlErrors, fmt.Sprintf("statement %d (%s): %s", i+1, mfStmt.Name.String(), w))
					}
				}
			}
			// Check view entity OQL
			if viewStmt, ok := stmt.(*ast.CreateViewEntityStmt); ok {
				if viewStmt.Query.RawQuery != "" {
					if errs := executor.ValidateOQLSyntax(viewStmt.Query.RawQuery); len(errs) > 0 {
						for _, e := range errs {
							oqlErrors = append(oqlErrors, fmt.Sprintf("statement %d (%s): %s", i+1, viewStmt.Name.String(), e))
						}
					}
					if errs := executor.ValidateOQLTypes(viewStmt.Query.RawQuery, viewStmt.Attributes); len(errs) > 0 {
						for _, e := range errs {
							oqlErrors = append(oqlErrors, fmt.Sprintf("statement %d (%s): %s", i+1, viewStmt.Name.String(), e))
						}
					}
				}
			}
		}
		if len(oqlErrors) > 0 {
			fmt.Fprintf(os.Stderr, "\nValidation errors found:\n")
			for _, e := range oqlErrors {
				fmt.Fprintf(os.Stderr, "  - %s\n", e)
			}
			os.Exit(1)
		}

		// If reference checking requested
		if checkRefs {
			if projectPath == "" {
				fmt.Fprintln(os.Stderr, "Error: --project (-p) is required for reference checking")
				os.Exit(1)
			}

			fmt.Printf("\nValidating references against: %s\n", projectPath)
			fmt.Printf("(Note: References to objects created within the script are skipped)\n")
			exec, logger := newLoggedExecutor("check")
			defer logger.Close()
			defer exec.Close()

			// Connect to project
			connectProg, _ := visitor.Build(fmt.Sprintf("CONNECT LOCAL '%s'", projectPath))
			for _, stmt := range connectProg.Statements {
				if err := exec.Execute(stmt); err != nil {
					fmt.Fprintf(os.Stderr, "Error connecting: %v\n", err)
					os.Exit(1)
				}
			}

			// Validate the program (considers objects defined within the script)
			validationErrors := exec.ValidateProgram(prog)
			if len(validationErrors) > 0 {
				fmt.Fprintf(os.Stderr, "Reference errors:\n")
				for _, err := range validationErrors {
					fmt.Fprintf(os.Stderr, "  %v\n", err)
				}
				fmt.Fprintf(os.Stderr, "\n✗ %d reference error(s) found\n", len(validationErrors))
				os.Exit(1)
			}
			fmt.Printf("✓ All references valid\n")
		}

		fmt.Println("\nCheck passed!")
	},
}
