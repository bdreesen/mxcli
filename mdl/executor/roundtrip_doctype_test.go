// SPDX-License-Identifier: Apache-2.0

//go:build integration

package executor

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/mdl/visitor"
)

// TestMxCheck_DoctypeScripts executes each doctype-tests/*.mdl example script
// in its own fresh Mendix project and validates the result with mx check.
//
// Each script runs in isolation so errors are cleanly attributed.
// Files matching *.test.mdl or *.tests.mdl are skipped (they require Docker).
func TestMxCheck_DoctypeScripts(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	// Locate doctype-tests directory
	doctypeDir, err := filepath.Abs("../../mdl-examples/doctype-tests")
	if err != nil {
		t.Fatalf("Failed to resolve doctype-tests path: %v", err)
	}
	if _, err := os.Stat(doctypeDir); err != nil {
		t.Skipf("doctype-tests directory not found at %s", doctypeDir)
	}

	// Collect eligible scripts (skip .test.mdl and .tests.mdl)
	entries, err := os.ReadDir(doctypeDir)
	if err != nil {
		t.Fatalf("Failed to read doctype-tests directory: %v", err)
	}

	var scripts []string
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".mdl") {
			continue
		}
		if strings.HasSuffix(name, ".test.mdl") || strings.HasSuffix(name, ".tests.mdl") {
			continue
		}
		scripts = append(scripts, name)
	}
	sort.Strings(scripts)

	if len(scripts) == 0 {
		t.Skip("no eligible MDL scripts found")
	}

	for _, name := range scripts {
		scriptPath := filepath.Join(doctypeDir, name)
		content, err := os.ReadFile(scriptPath)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", name, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Fresh project for each script
			env := setupTestEnv(t)
			defer env.teardown()

			// Execute the script
			prog, errs := visitor.Build(string(content))
			if len(errs) > 0 {
				t.Fatalf("Parse error: %v", errs[0])
			}

			if err := env.executor.ExecuteProgram(prog); err != nil {
				t.Errorf("Execution error: %v", err)
			}

			// Flush to disk
			env.executor.Execute(&ast.DisconnectStmt{})

			// Run mx check
			output, mxErr := runMxCheck(t, env.projectPath)
			if mxErr != nil {
				lowerOutput := strings.ToLower(output)
				if strings.Contains(lowerOutput, "error") {
					t.Errorf("mx check found errors:\n%s", output)
				} else {
					t.Logf("mx check output:\n%s", output)
				}
			} else {
				t.Logf("mx check passed: 0 errors")
			}
		})
	}
}
