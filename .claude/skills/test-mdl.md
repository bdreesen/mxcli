# Test MDL Script

Use this skill to test MDL scripts against the ModelSDK Go implementation and verify they work correctly with Studio Pro.

## When to Use This Skill

- After implementing BSON serialization for a new document type
- To verify MDL scripts execute without errors
- To debug Studio Pro loading issues
- To compare BSON output with expected structure

## Test Workflow

### 1. Build the CLI

```bash
go build -o bin/mxcli ./cmd/mxcli
```

### 2. Run MDL Script

Execute MDL commands against the test project:

```bash
# Interactive REPL
./bin/mxcli

# Or direct execution
./bin/mxcli -p mx-test-projects/test1-go-app/test1-go.mpr -c "SHOW ENTITIES IN DmTest"
```

### 3. Execute Script File

```bash
./bin/mxcli -p mx-test-projects/test1-go-app/test1-go.mpr exec reference/mendix-repl/examples/doctype-tests/domain-model-examples.mdl
```

### 4. Verify in Studio Pro

After executing MDL:
1. Open the MPR file in Mendix Studio Pro
2. Check for errors in the Error pane
3. Navigate to the domain model to verify entities appear correctly

## Common Errors and Solutions

### TypeCacheUnknownTypeException

```
The type cache does not contain a type with qualified name DomainModels$Index
```

**Cause**: Using `qualifiedName` instead of `storageName` for `$Type` field.

**Fix**: Check `reference/mendixmodellib/reflection-data/<version>-structures.json` for the correct `storageName`.

### System.ArgumentNullException

```
System.ArgumentNullException: Value cannot be null. (Parameter 'AttributeId')
```

**Cause**: Wrong reference format. Using UUID for BY_NAME_REFERENCE or vice versa.

**Fix**: Check metamodel `typeInfo.kind` for the property:
- `BY_NAME_REFERENCE` → qualified name string (e.g., "Module.Entity.Attr")
- `BY_ID_REFERENCE` → binary UUID

### Enumeration Not Displayed

Enumeration attribute shows as just "Enumeration" without the reference.

**Fix**: Add `Enumeration` field with qualified name to the attribute type in BSON.

## Debug Tools

### Dump BSON for Comparison

Use the debug examples to inspect BSON:

```bash
go run ./examples/debug_bson/main.go mx-test-projects/test1-go-app/test1-go.mpr DmTest
```

### Compare with Studio Pro Output

1. Create entity manually in Studio Pro
2. Save project
3. Dump BSON using debug tool
4. Compare with SDK-generated BSON

### Check Metamodel

```bash
# Find type definition
grep -A 30 '"DomainModels\$Index"' reference/mendixmodellib/reflection-data/11.6.0-structures.json

# Find property reference kind
grep -B 5 -A 10 '"storageName" : "Attribute"' reference/mendixmodellib/reflection-data/11.6.0-structures.json
```

## Test Project

The test project is at:
```
mx-test-projects/test1-go-app/test1-go.mpr
```

Make a backup before testing destructive operations:
```bash
cp -r mx-test-projects/test1-go-app mx-test-projects/test1-go-app.bak
```

## Example Test Session

```bash
# Build
go build -o bin/mxcli ./cmd/mxcli

# Connect and show current state
./bin/mxcli -p mx-test-projects/test1-go-app/test1-go.mpr -c "SHOW ENTITIES IN DmTest"

# Execute test script
./bin/mxcli -p mx-test-projects/test1-go-app/test1-go.mpr exec reference/mendix-repl/examples/doctype-tests/domain-model-examples.mdl

# Verify created entity
./bin/mxcli -p mx-test-projects/test1-go-app/test1-go.mpr -c "DESCRIBE ENTITY DmTest.SalesOrder"
```

## Checklist Before Marking Complete

- [ ] MDL script executes without errors
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] Studio Pro opens project without errors
- [ ] Created elements appear correctly in Studio Pro
- [ ] BSON mapping documentation updated if new patterns discovered
