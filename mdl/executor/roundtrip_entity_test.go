// SPDX-License-Identifier: Apache-2.0

//go:build integration

package executor

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/mdl/visitor"
)

// --- Roundtrip Tests ---

func TestRoundtripEntity_Simple(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	entityName := testModule + ".TestEntitySimple"

	// Create entity (Boolean auto-defaults to false if no DEFAULT specified)
	createMDL := `CREATE OR MODIFY PERSISTENT ENTITY ` + entityName + ` (
		Name: String(100),
		Age: Integer,
		Active: Boolean DEFAULT false
	);`

	// Use diff-based helper to verify roundtrip
	env.assertContains(createMDL, []string{
		"PERSISTENT ENTITY",
		"Name:",
		"String(100)",
		"Age:",
		"Integer",
		"Active:",
		"Boolean",
	})
}

func TestRoundtripEntity_WithConstraints(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	entityName := testModule + ".TestEntityConstraints"

	// Create entity with constraints
	createMDL := `CREATE OR MODIFY PERSISTENT ENTITY ` + entityName + ` (
		Code: String(50) NOT NULL,
		Email: String(200) UNIQUE
	);`

	// Use diff-based helper - NOT NULL may be output as REQUIRED
	env.assertContains(createMDL, []string{
		"PERSISTENT ENTITY",
		"Code:",
		"String(50)",
		"Email:",
		"String(200)",
		"UNIQUE",
	})
}

func TestRoundtripEntity_WithIndex(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	entityName := testModule + ".TestEntityIndex"

	// Create entity with index
	createMDL := `CREATE OR MODIFY PERSISTENT ENTITY ` + entityName + ` (
		Code: String(50),
		Name: String(100)
	)
	INDEX (Code);`

	// Use diff-based helper to verify roundtrip
	env.assertContains(createMDL, []string{
		"PERSISTENT ENTITY",
		"Code:",
		"String(50)",
		"Name:",
		"String(100)",
		"INDEX",
	})
}

func TestRoundtripEntity_WithEventHandler(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	entityName := testModule + ".TestEntityEventHandler"
	mfName := testModule + ".ACT_ValidateTestEntity"

	// Create a microflow first (event handler references it)
	if err := env.executeMDL(`CREATE OR MODIFY MICROFLOW ` + mfName + ` ()
BEGIN
  LOG INFO 'validating';
END;`); err != nil {
		t.Fatalf("failed to create microflow: %v", err)
	}

	// Create entity with event handler
	createMDL := `CREATE OR MODIFY PERSISTENT ENTITY ` + entityName + ` (
		Name: String(100)
	)
	ON BEFORE COMMIT CALL ` + mfName + ` RAISE ERROR;`

	// Verify roundtrip preserves the event handler
	env.assertContains(createMDL, []string{
		"PERSISTENT ENTITY",
		"Name:",
		"String(100)",
		"ON BEFORE COMMIT CALL",
		mfName,
		"RAISE ERROR",
	})
}

func TestRoundtripEntity_AlterAddDropEventHandler(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	entityName := testModule + ".TestAlterEventHandler"
	mfName := testModule + ".ACT_AlterEventTest"

	// Create microflow
	if err := env.executeMDL(`CREATE OR MODIFY MICROFLOW ` + mfName + ` ()
BEGIN
  LOG INFO 'test';
END;`); err != nil {
		t.Fatalf("failed to create microflow: %v", err)
	}

	// Create entity without handlers
	if err := env.executeMDL(`CREATE OR MODIFY PERSISTENT ENTITY ` + entityName + ` (
		Code: String(50)
	);`); err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	// Add event handler via ALTER
	if err := env.executeMDL(`ALTER ENTITY ` + entityName + `
		ADD EVENT HANDLER ON AFTER CREATE CALL ` + mfName + `;`); err != nil {
		t.Fatalf("failed to add event handler: %v", err)
	}

	// Verify handler appears in DESCRIBE
	out, err := env.describeMDL(`DESCRIBE ENTITY ` + entityName + `;`)
	if err != nil {
		t.Fatalf("describe failed: %v", err)
	}
	if !strings.Contains(out, "ON AFTER CREATE CALL") {
		t.Errorf("expected ON AFTER CREATE CALL in DESCRIBE output, got:\n%s", out)
	}
	if !strings.Contains(out, mfName) {
		t.Errorf("expected microflow name %q in DESCRIBE output, got:\n%s", mfName, out)
	}

	// Drop the event handler
	if err := env.executeMDL(`ALTER ENTITY ` + entityName + `
		DROP EVENT HANDLER ON AFTER CREATE;`); err != nil {
		t.Fatalf("failed to drop event handler: %v", err)
	}

	// Verify handler is gone
	out, err = env.describeMDL(`DESCRIBE ENTITY ` + entityName + `;`)
	if err != nil {
		t.Fatalf("describe after drop failed: %v", err)
	}
	if strings.Contains(out, "ON AFTER CREATE CALL") {
		t.Errorf("event handler should be removed but still in DESCRIBE output:\n%s", out)
	}
}

func TestRoundtripEnumeration(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	enumName := testModule + ".TestStatus"

	// Create enumeration
	createMDL := `CREATE ENUMERATION ` + enumName + ` (
		Active 'Active',
		Inactive 'Inactive',
		Pending 'Pending Review'
	);`

	// Use diff-based helper to verify roundtrip
	env.assertContains(createMDL, []string{
		"ENUMERATION",
		"Active",
		"Inactive",
		"Pending",
	})
}

// --- Benchmark Tests ---

func BenchmarkRoundtripEntity(b *testing.B) {
	// Skip if source project doesn't exist
	srcDir, _ := filepath.Abs(sourceProject)
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		b.Skip("Source project not found")
	}

	// Copy source project for benchmark (read-only use, but keeps pattern consistent)
	destDir := b.TempDir()
	srcMPR := filepath.Join(srcDir, sourceProjectMPR)
	destMPR := filepath.Join(destDir, sourceProjectMPR)
	if err := copyFile(srcMPR, destMPR); err != nil {
		b.Fatalf("Failed to copy MPR: %v", err)
	}
	for _, dir := range []string{"mprcontents", "widgets", "themesource", "theme", "javascriptsource"} {
		srcSub := filepath.Join(srcDir, dir)
		if _, serr := os.Stat(srcSub); serr == nil {
			if err := copyDir(srcSub, filepath.Join(destDir, dir)); err != nil {
				b.Fatalf("Failed to copy %s: %v", dir, err)
			}
		}
	}

	output := &bytes.Buffer{}
	exec := New(output)

	// Connect once
	exec.Execute(&ast.ConnectStmt{
		Path: destMPR,
	})
	defer exec.Execute(&ast.DisconnectStmt{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output.Reset()

		// Create and describe entity
		prog, _ := visitor.Build(`DESCRIBE ENTITY MyFirstModule.MyEntity;`)
		for _, stmt := range prog.Statements {
			exec.Execute(stmt)
		}
	}
}
