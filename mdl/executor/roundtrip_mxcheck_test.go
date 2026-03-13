// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mendixlabs/mxcli/mdl/ast"
)

// --- MX Check Integration Tests ---
// These tests verify that created documents pass Mendix's validation.

// mxCheckPath is the path to the mx check command.
const mxCheckPath = "../../reference/mxbuild/modeler/mx"

// mxCheckAvailable checks if the mx command is available.
func mxCheckAvailable() bool {
	absPath, err := filepath.Abs(mxCheckPath)
	if err != nil {
		return false
	}
	_, err = os.Stat(absPath)
	return err == nil
}

// runMxCheck runs mx check on the given project and returns any errors.
func runMxCheck(t *testing.T, projectPath string) (string, error) {
	t.Helper()

	mxPath, err := filepath.Abs(mxCheckPath)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(mxPath, "check", projectPath)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// TestMxCheck_Entity creates an entity and verifies mx check passes.
func TestMxCheck_Entity(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	entityName := testModule + ".MxCheckEntity"
	env.registerCleanup("entity", entityName)

	// Create entity
	createMDL := `CREATE OR MODIFY PERSISTENT ENTITY ` + entityName + ` (
		Code: String(50) NOT NULL,
		Description: String(500),
		Count: Integer DEFAULT 0
	);`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Disconnect to flush changes
	env.executor.Execute(&ast.DisconnectStmt{})

	// Run mx check
	output, err := runMxCheck(t, env.projectPath)
	if err != nil {
		// mx check returns non-zero exit code if there are errors
		if strings.Contains(output, "error") || strings.Contains(output, "Error") {
			t.Errorf("mx check found errors:\n%s", output)
		} else {
			t.Logf("mx check output:\n%s", output)
		}
	} else {
		t.Logf("mx check passed:\n%s", output)
	}
}

// TestMxCheck_Enumeration creates an enumeration and verifies mx check passes.
func TestMxCheck_Enumeration(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	enumName := testModule + ".MxCheckPriority"
	env.registerCleanup("enumeration", enumName)

	// Create enumeration
	createMDL := `CREATE ENUMERATION ` + enumName + ` (
		Low 'Low Priority',
		Medium 'Medium Priority',
		High 'High Priority',
		Critical 'Critical'
	);`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create enumeration: %v", err)
	}

	// Disconnect to flush changes
	env.executor.Execute(&ast.DisconnectStmt{})

	// Run mx check
	output, err := runMxCheck(t, env.projectPath)
	if err != nil {
		if strings.Contains(output, "error") || strings.Contains(output, "Error") {
			t.Errorf("mx check found errors:\n%s", output)
		} else {
			t.Logf("mx check output:\n%s", output)
		}
	} else {
		t.Logf("mx check passed:\n%s", output)
	}
}

// TestMxCheck_RetrieveWithLimit validates that RETRIEVE with LIMIT produces
// BSON that passes mx check (Studio Pro validation).
func TestMxCheck_RetrieveWithLimit(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	if err := env.executeMDL(`CREATE OR MODIFY PERSISTENT ENTITY RoundtripTest.MxCheckItem (Name: String(100));`); err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	mfName := testModule + ".MxCheck_RetrieveLimit"
	env.registerCleanup("microflow", mfName)

	createMDL := `CREATE MICROFLOW ` + mfName + ` () RETURNS Boolean
BEGIN
  RETRIEVE $Item FROM RoundtripTest.MxCheckItem
    LIMIT 1;
  RETURN true;
END;`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create microflow: %v", err)
	}

	env.executor.Execute(&ast.DisconnectStmt{})

	output, err := runMxCheck(t, env.projectPath)
	if err != nil {
		if strings.Contains(output, "error") || strings.Contains(output, "Error") {
			t.Errorf("mx check found errors:\n%s", output)
		} else {
			t.Logf("mx check output:\n%s", output)
		}
	} else {
		t.Logf("mx check passed:\n%s", output)
	}
}

// TestMxCheck_RetrieveWithLimitOffset validates that RETRIEVE with LIMIT and OFFSET
// produces BSON that passes mx check. Regression guard for LimitExpression/OffsetExpression
// being stored in the correct BSON fields within Microflows$ConstantRange.
func TestMxCheck_RetrieveWithLimitOffset(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	if err := env.executeMDL(`CREATE OR MODIFY PERSISTENT ENTITY RoundtripTest.MxCheckItem (Name: String(100));`); err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	mfName := testModule + ".MxCheck_RetrieveLimOff"
	env.registerCleanup("microflow", mfName)

	createMDL := `CREATE MICROFLOW ` + mfName + ` () RETURNS Boolean
BEGIN
  RETRIEVE $Items FROM RoundtripTest.MxCheckItem
    LIMIT 2
    OFFSET 3;
  RETURN true;
END;`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create microflow: %v", err)
	}

	env.executor.Execute(&ast.DisconnectStmt{})

	output, err := runMxCheck(t, env.projectPath)
	if err != nil {
		if strings.Contains(output, "error") || strings.Contains(output, "Error") {
			t.Errorf("mx check found errors:\n%s", output)
		} else {
			t.Logf("mx check output:\n%s", output)
		}
	} else {
		t.Logf("mx check passed:\n%s", output)
	}
}

// TestMxCheck_RetrieveWithSortBy validates that RETRIEVE with SORT BY produces
// BSON that passes mx check. Regression guard for sort items being stored under
// the correct BSON key (NewSortings vs sortItemList).
func TestMxCheck_RetrieveWithSortBy(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	if err := env.executeMDL(`CREATE OR MODIFY PERSISTENT ENTITY RoundtripTest.MxCheckItem (Name: String(100));`); err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	mfName := testModule + ".MxCheck_RetrieveSort"
	env.registerCleanup("microflow", mfName)

	createMDL := `CREATE MICROFLOW ` + mfName + ` () RETURNS Boolean
BEGIN
  RETRIEVE $Items FROM RoundtripTest.MxCheckItem
    SORT BY RoundtripTest.MxCheckItem.Name ASC;
  RETURN true;
END;`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create microflow: %v", err)
	}

	env.executor.Execute(&ast.DisconnectStmt{})

	output, err := runMxCheck(t, env.projectPath)
	if err != nil {
		if strings.Contains(output, "error") || strings.Contains(output, "Error") {
			t.Errorf("mx check found errors:\n%s", output)
		} else {
			t.Logf("mx check output:\n%s", output)
		}
	} else {
		t.Logf("mx check passed:\n%s", output)
	}
}

// TestMxCheck_RetrieveWithWhereSortLimitOffset validates the full RETRIEVE pattern
// (WHERE + SORT BY + LIMIT + OFFSET) passes mx check. This matches the
// M028_DataForm_Getter microflow pattern.
func TestMxCheck_RetrieveWithWhereSortLimitOffset(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	if err := env.executeMDL(`CREATE OR MODIFY PERSISTENT ENTITY RoundtripTest.MxCheckItem (Name: String(100));`); err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	mfName := testModule + ".MxCheck_RetrieveFull"
	env.registerCleanup("microflow", mfName)

	createMDL := `CREATE MICROFLOW ` + mfName + ` () RETURNS Boolean
BEGIN
  RETRIEVE $Items FROM RoundtripTest.MxCheckItem
    WHERE (starts-with(Name, 'a'))
    SORT BY RoundtripTest.MxCheckItem.Name ASC
    LIMIT 2
    OFFSET 3;
  RETURN true;
END;`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create microflow: %v", err)
	}

	env.executor.Execute(&ast.DisconnectStmt{})

	output, err := runMxCheck(t, env.projectPath)
	if err != nil {
		if strings.Contains(output, "error") || strings.Contains(output, "Error") {
			t.Errorf("mx check found errors:\n%s", output)
		} else {
			t.Logf("mx check output:\n%s", output)
		}
	} else {
		t.Logf("mx check passed:\n%s", output)
	}
}

// TestMxCheck_MicroflowWithCallParams tests microflow with CALL unified param syntax.
func TestMxCheck_MicroflowWithCallParams(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	// Create a helper microflow
	helperName := testModule + ".MxCheckCallHelper"
	env.registerCleanup("microflow", helperName)

	createHelperMDL := `CREATE MICROFLOW ` + helperName + ` ($InputValue: String) RETURNS String
	BEGIN
		RETURN $InputValue;
	END;`

	if err := env.executeMDL(createHelperMDL); err != nil {
		t.Fatalf("Failed to create helper microflow: %v", err)
	}

	// Create caller microflow with unified param syntax
	callerName := testModule + ".MxCheckCallCaller"
	env.registerCleanup("microflow", callerName)

	createCallerMDL := `CREATE MICROFLOW ` + callerName + ` () RETURNS String
	BEGIN
		$Result = CALL MICROFLOW ` + helperName + ` (InputValue = 'TestValue');
		RETURN $Result;
	END;`

	if err := env.executeMDL(createCallerMDL); err != nil {
		t.Fatalf("Failed to create caller microflow: %v", err)
	}

	// Disconnect to flush changes
	env.executor.Execute(&ast.DisconnectStmt{})

	// Run mx check
	output, err := runMxCheck(t, env.projectPath)
	if err != nil {
		if strings.Contains(output, "error") || strings.Contains(output, "Error") {
			t.Errorf("mx check found errors:\n%s", output)
		} else {
			t.Logf("mx check output:\n%s", output)
		}
	} else {
		t.Logf("mx check passed:\n%s", output)
	}
}
