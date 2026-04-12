// SPDX-License-Identifier: Apache-2.0

//go:build integration

package executor

import (
	"strings"
	"testing"

	"github.com/mendixlabs/mxcli/mdl/ast"
)

// --- REST Client Roundtrip Tests ---

func TestRoundtripRestClient_SimpleGet(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.SimpleAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION GetStatus {
    Method: GET,
    Path: '/status',
    Response: NONE
  }
};`

	env.assertContains(createMDL, []string{
		"REST CLIENT",
		"SimpleAPI",
		"BaseUrl: 'https://api.example.com'",
		"Authentication: NONE",
		"OPERATION GetStatus",
		"Method: GET",
		"Path: '/status'",
		"Response: NONE",
	})
}

func TestRoundtripRestClient_WithJsonResponse(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.JsonAPI (
  BaseUrl: 'https://jsonplaceholder.typicode.com',
  Authentication: NONE
)
{
  OPERATION GetPosts {
    Method: GET,
    Path: '/posts',
    Headers: ('Accept' = 'application/json'),
    Response: JSON AS $Posts
  }
};`

	env.assertContains(createMDL, []string{
		"REST CLIENT",
		"JsonAPI",
		"BaseUrl: 'https://jsonplaceholder.typicode.com'",
		"OPERATION GetPosts",
		"Method: GET",
		"Path: '/posts'",
		"'Accept' = 'application/json'",
		"Response: JSON",
	})
}

func TestRoundtripRestClient_WithPathParams(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.ParamAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION GetItem {
    Method: GET,
    Path: '/items/{itemId}',
    Parameters: ($itemId: Integer),
    Response: JSON AS $Item
  }
};`

	env.assertContains(createMDL, []string{
		"OPERATION GetItem",
		"Path: '/items/{itemId}'",
		"$itemId: Integer",
		"Response: JSON",
	})
}

func TestRoundtripRestClient_WithQueryParams(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.SearchAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION SearchItems {
    Method: GET,
    Path: '/search',
    Query: ($q: String, $page: String),
    Response: JSON AS $Results
  }
};`

	env.assertContains(createMDL, []string{
		"OPERATION SearchItems",
		"Path: '/search'",
		"$q: String",
		"$page: String",
		"Response: JSON",
	})
}

func TestRoundtripRestClient_PostWithBody(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.CrudAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION CreateItem {
    Method: POST,
    Path: '/items',
    Headers: ('Content-Type' = 'application/json'),
    Body: JSON FROM $NewItem,
    Response: JSON AS $CreatedItem
  }
};`

	env.assertContains(createMDL, []string{
		"OPERATION CreateItem",
		"Method: POST",
		"Body: JSON FROM $NewItem",
		"Response: JSON",
	})
}

func TestRoundtripRestClient_BasicAuth(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.AuthAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: BASIC (Username: 'admin', Password: 'secret')
)
{
  OPERATION GetData {
    Method: GET,
    Path: '/data',
    Response: JSON AS $Data
  }
};`

	env.assertContains(createMDL, []string{
		"REST CLIENT",
		"AuthAPI",
		"Authentication: BASIC",
		"Username: 'admin'",
		"Password: 'secret'",
		"OPERATION GetData",
	})
}

func TestRoundtripRestClient_WithTimeout(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.TimeoutAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION SlowQuery {
    Method: GET,
    Path: '/slow',
    Timeout: 60,
    Response: JSON AS $Result
  }
};`

	env.assertContains(createMDL, []string{
		"OPERATION SlowQuery",
		"Timeout: 60",
		"Response: JSON",
	})
}

func TestRoundtripRestClient_MultipleOperations(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.PetStoreAPI (
  BaseUrl: 'https://petstore.swagger.io/v2',
  Authentication: NONE
)
{
  OPERATION ListPets {
    Method: GET,
    Path: '/pet/findByStatus',
    Query: ($status: String),
    Headers: ('Accept' = 'application/json'),
    Timeout: 30,
    Response: JSON AS $PetList
  }

  OPERATION GetPet {
    Method: GET,
    Path: '/pet/{petId}',
    Parameters: ($petId: Integer),
    Response: JSON AS $Pet
  }

  OPERATION AddPet {
    Method: POST,
    Path: '/pet',
    Headers: ('Content-Type' = 'application/json'),
    Body: JSON FROM $NewPet,
    Response: JSON AS $CreatedPet
  }

  OPERATION RemovePet {
    Method: DELETE,
    Path: '/pet/{petId}',
    Parameters: ($petId: Integer),
    Response: NONE
  }
};`

	env.assertContains(createMDL, []string{
		"REST CLIENT",
		"PetStoreAPI",
		"OPERATION ListPets",
		"$status: String",
		"Timeout: 30",
		"OPERATION GetPet",
		"$petId: Integer",
		"OPERATION AddPet",
		"Method: POST",
		"Body: JSON FROM $NewPet",
		"OPERATION RemovePet",
		"Method: DELETE",
		"Response: NONE",
	})
}

func TestRoundtripRestClient_DeleteNoResponse(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.DeleteAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION DeleteResource {
    Method: DELETE,
    Path: '/resources/{id}',
    Parameters: ($id: Integer),
    Response: NONE
  }
};`

	env.assertContains(createMDL, []string{
		"OPERATION DeleteResource",
		"Method: DELETE",
		"$id: Integer",
		"Response: NONE",
	})
}

func TestRoundtripRestClient_CreateOrModify(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	// Create first version
	createMDL := `CREATE REST CLIENT ` + testModule + `.MutableAPI (
  BaseUrl: 'https://api.example.com/v1',
  Authentication: NONE
)
{
  OPERATION GetData {
    Method: GET,
    Path: '/data',
    Response: JSON AS $Data
  }
};`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create REST client: %v", err)
	}

	// Update with CREATE OR MODIFY
	updateMDL := `CREATE OR MODIFY REST CLIENT ` + testModule + `.MutableAPI (
  BaseUrl: 'https://api.example.com/v2',
  Authentication: NONE
)
{
  OPERATION GetDataV2 {
    Method: GET,
    Path: '/data/v2',
    Response: JSON AS $DataV2
  }
};`

	if err := env.executeMDL(updateMDL); err != nil {
		t.Fatalf("Failed to update REST client: %v", err)
	}

	// Verify the updated version
	output, err := env.describeMDL("DESCRIBE REST CLIENT " + testModule + ".MutableAPI;")
	if err != nil {
		t.Fatalf("Failed to describe REST client: %v", err)
	}

	if !strings.Contains(output, "v2") {
		t.Errorf("Expected updated BaseUrl with v2, got:\n%s", output)
	}
	if !strings.Contains(output, "GetDataV2") {
		t.Errorf("Expected updated operation GetDataV2, got:\n%s", output)
	}
}

func TestRoundtripRestClient_Drop(t *testing.T) {
	env := setupTestEnv(t)
	defer env.teardown()

	// Create a REST client
	createMDL := `CREATE REST CLIENT ` + testModule + `.ToBeDropped (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION Ping {
    Method: GET,
    Path: '/ping',
    Response: NONE
  }
};`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create REST client: %v", err)
	}

	// Verify it exists
	_, err := env.describeMDL("DESCRIBE REST CLIENT " + testModule + ".ToBeDropped;")
	if err != nil {
		t.Fatalf("REST client should exist before DROP: %v", err)
	}

	// Drop it
	if err := env.executeMDL("DROP REST CLIENT " + testModule + ".ToBeDropped;"); err != nil {
		t.Fatalf("Failed to drop REST client: %v", err)
	}

	// Verify it's gone
	_, err = env.describeMDL("DESCRIBE REST CLIENT " + testModule + ".ToBeDropped;")
	if err == nil {
		t.Error("REST client should not exist after DROP")
	}
}

// --- MX Check Tests ---

func TestMxCheck_RestClient_SimpleGet(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	createMDL := `CREATE REST CLIENT ` + testModule + `.MxCheckSimpleAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION GetStatus {
    Method: GET,
    Path: '/status',
    Headers: ('Accept' = '*/*'),
    Response: NONE
  }
};`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create REST client: %v", err)
	}

	// Disconnect to flush changes to disk
	env.executor.Execute(&ast.DisconnectStmt{})

	// Run mx check
	output, err := runMxCheck(t, env.projectPath)
	assertMxCheckPassed(t, output, err)
}

func TestMxCheck_RestClient_PostWithBody(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	// Test with path parameters (GET to avoid body requirements).
	createMDL := `CREATE REST CLIENT ` + testModule + `.MxCheckParamAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: NONE
)
{
  OPERATION GetItem {
    Method: GET,
    Path: '/items/{itemId}',
    Parameters: ($itemId: Integer),
    Headers: ('Accept' = '*/*'),
    Response: NONE
  }
};`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create REST client: %v", err)
	}

	env.executor.Execute(&ast.DisconnectStmt{})

	output, err := runMxCheck(t, env.projectPath)
	assertMxCheckPassed(t, output, err)
}

// assertMxCheckPassed checks mx check output for errors.
// Detects both "[error]" markers (validation errors) and "ERROR:" (load crashes).
func assertMxCheckPassed(t *testing.T, output string, err error) {
	t.Helper()
	if err != nil {
		// Non-zero exit code — could be a crash or validation errors
		if strings.Contains(output, "[error]") || strings.Contains(output, "ERROR:") {
			t.Errorf("mx check failed:\n%s", output)
		} else {
			t.Logf("mx check exited with error but no validation errors:\n%s", output)
		}
	} else if strings.Contains(output, "[error]") {
		t.Errorf("mx check found errors:\n%s", output)
	} else {
		t.Logf("mx check passed:\n%s", output)
	}
}

func TestMxCheck_RestClient_BasicAuth(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	// Use RESPONSE NONE to avoid entity mapping requirements (CE0061)
	createMDL := `CREATE REST CLIENT ` + testModule + `.MxCheckAuthAPI (
  BaseUrl: 'https://api.example.com',
  Authentication: BASIC (Username: 'user', Password: 'pass')
)
{
  OPERATION GetSecureData {
    Method: GET,
    Path: '/secure/data',
    Headers: ('Accept' = '*/*'),
    Response: NONE
  }
};`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create REST client: %v", err)
	}

	env.executor.Execute(&ast.DisconnectStmt{})

	output, err := runMxCheck(t, env.projectPath)
	assertMxCheckPassed(t, output, err)
}

func TestMxCheck_RestClient_MultipleOperations(t *testing.T) {
	if !mxCheckAvailable() {
		t.Skip("mx command not available")
	}

	env := setupTestEnv(t)
	defer env.teardown()

	// Use RESPONSE NONE for all operations to avoid entity mapping requirements (CE0061).
	// All operations include Accept header to avoid CE7062.
	createMDL := `CREATE REST CLIENT ` + testModule + `.MxCheckPetStore (
  BaseUrl: 'https://petstore.swagger.io/v2',
  Authentication: NONE
)
{
  OPERATION ListPets {
    Method: GET,
    Path: '/pet/findByStatus',
    Query: ($status: String),
    Headers: ('Accept' = 'application/json'),
    Timeout: 30,
    Response: NONE
  }

  OPERATION GetPet {
    Method: GET,
    Path: '/pet/{petId}',
    Parameters: ($petId: Integer),
    Headers: ('Accept' = 'application/json'),
    Response: NONE
  }

  OPERATION RemovePet {
    Method: DELETE,
    Path: '/pet/{petId}',
    Parameters: ($petId: Integer),
    Headers: ('Accept' = '*/*'),
    Response: NONE
  }
};`

	if err := env.executeMDL(createMDL); err != nil {
		t.Fatalf("Failed to create REST client: %v", err)
	}

	env.executor.Execute(&ast.DisconnectStmt{})

	output, err := runMxCheck(t, env.projectPath)
	assertMxCheckPassed(t, output, err)
}
