# Proposal: Extend UPDATE WIDGETS to Built-in Widgets via Schema Registry

**Status:** Draft
**Author:** Claude
**Date:** 2026-02-18

## Problem Statement

`UPDATE WIDGETS` currently only works for **pluggable widgets** (ComboBox, DataGrid2, etc.) — it cannot modify properties on **built-in widgets** like `Forms$LayoutGrid`, `Forms$DataView`, `Forms$TextBox`, etc.

A real user scenario: 51 pages in a module needed `Width` changed from `"FixedWidth"` to `"FullWidth"` on `Forms$LayoutGrid`. The only workaround was `dump-bson`.

### Root Cause: Lossy Round-Trip

Built-in widget properties are stored as direct BSON fields (e.g., `"Width": "FullWidth"`), but the Go structs don't capture most of them. `LayoutGrid` has no `Width` field — it's hardcoded to `"FullWidth"` during serialization:

```
Read from MPR (BSON) → Parse to Go struct → Modify struct → Serialize back to BSON
                          ↑                                         ↑
                    Drops Width field                    Hardcodes Width="FullWidth"
```

## Design Principles

1. **What you can see = What you can change** — Any property that UPDATE WIDGETS can set MUST also be visible in DESCRIBE PAGE output. Neither ships without the other.

2. **Schema-driven, not hand-curated** — Property metadata (types, allowed values, defaults) comes from the Mendix reflection data, not from a manually maintained registry. This aligns with the [BSON Schema Registry proposal](BSON_SCHEMA_REGISTRY_PROPOSAL.md) and avoids creating a parallel metadata system.

3. **This is Phase 1 of the Schema Registry** — scoped to the built-in widget property read/write use case, but built on the same foundation that later phases will use for field completeness, validation, and version migration.

## Proposed Solution: Reflection-Data-Driven Property Metadata

### How the reflection data solves this

The reflection data (`reference/mendixmodellib/reflection-data/{version}-structures.json`) already contains exact property metadata for every built-in widget type. For example, `Pages$LayoutGrid`:

```json
{
  "properties": {
    "width": {
      "name": "width",
      "storageName": "Width",
      "typeInfo": {
        "type": "ENUMERATION",
        "values": ["FixedWidth", "FullWidth"]
      }
    }
  },
  "defaultSettings": {
    "width": "FullWidth"
  }
}
```

This gives us everything a hand-curated registry would — BSON field name, allowed values, default value — but for **every property of every type**, automatically, per Mendix version. No manual maintenance.

### Which properties are settable?

Not all BSON properties make sense for `UPDATE WIDGETS SET`. Only **scalar** properties can be assigned a simple value:

| `typeInfo.type` | Settable? | Example |
|----------------|-----------|---------|
| `ENUMERATION` | Yes | `Width` = `"FullWidth"`, `Editability` = `"Never"` |
| `PRIMITIVE` (STRING) | Yes | `Name` = `"grid1"` |
| `PRIMITIVE` (INTEGER) | Yes | `TabIndex` = `5` |
| `PRIMITIVE` (BOOLEAN) | Yes | `ShowFooter` = `true` |
| `ELEMENT` (PART) | No | `Appearance`, `DataSource` — complex nested objects |
| `ELEMENT` (BY_*_REFERENCE) | No | References to other elements |

The schema naturally provides this filter — only properties with `typeInfo.type` of `ENUMERATION` or `PRIMITIVE` are eligible for the `SET` clause.

### Architecture

```
Reflection Data ({version}-structures.json)
         │
         ▼
   SchemaRegistry (loaded at project open, keyed by storage name)
         │
    ┌────┴─────┐
    ▼          ▼
DESCRIBE    UPDATE WIDGETS
  PAGE        (write)
  (read)
    │          │
    ▼          ▼
  Display    Validate value against typeInfo.values
  non-default   Set BSON field directly in raw map
  scalar props   Write raw BSON back to MPR
```

Both read and write use the **same schema lookup**, so they are always in sync.

## Implementation Plan

### Step 1: Create `sdk/schema/` runtime package

Move the schema types out of `internal/codegen/schema/` (build-time only) into a new `sdk/schema/` package usable at runtime.

**File:** `sdk/schema/registry.go` (new)

```go
package schema

// Registry holds type schemas for a specific Mendix version.
// This is a minimal Phase 1 slice of the BSON Schema Registry.
type Registry struct {
    version      string
    byQualified  map[string]*TypeDefinition  // "Pages$LayoutGrid" → def
    byStorage    map[string]*TypeDefinition  // "Forms$LayoutGrid" → def
}

// LoadRegistry loads reflection data for the given Mendix version.
// Falls back to nearest available version if exact match not found.
func LoadRegistry(version string) (*Registry, error)

// LookupByStorage finds a type by its BSON $Type name.
func (r *Registry) LookupByStorage(storageName string) *TypeDefinition

// ScalarProperties returns only ENUMERATION and PRIMITIVE properties
// for a type — the ones that can be displayed and set via UPDATE WIDGETS.
func (r *Registry) ScalarProperties(storageName string) []*PropertyDef
```

The `TypeDefinition`, `PropertyDef`, `TypeInfo` types are the same structs already defined in `internal/codegen/schema/types.go`. We either:
- (a) Move them to `sdk/schema/` and have `internal/codegen/` import from there, or
- (b) Copy the subset needed (types + JSON parsing) to avoid breaking the codegen package

Option (a) is cleaner.

### Step 2: Embed reflection data for key versions

Following the BSON Schema Registry proposal's guidance on version selection, embed ~15 versions:

```go
//go:embed reflection-data/6.10.4-structures.json
//go:embed reflection-data/7.23.0-structures.json
//go:embed reflection-data/8.18.0-structures.json
//go:embed reflection-data/9.24.0-structures.json
//go:embed reflection-data/10.6.0-structures.json
//go:embed reflection-data/10.12.0-structures.json
//go:embed reflection-data/10.18.0-structures.json
//go:embed reflection-data/11.0.0-structures.json
// ... through 11.6.0
var reflectionFS embed.FS
```

**Size**: ~1.2 MB per version × ~15 = ~18 MB uncompressed. With Go's embed compression this is manageable, and we only parse one version at runtime.

**Version matching**: Detect the project's Mendix version from MPR metadata (already available via `reader.Version()`), find the nearest available embedded version.

### Step 3: Lazy-load registry in Reader

**File:** `sdk/mpr/reader_documents.go` (modify)

```go
// SchemaRegistry returns the type schema registry for this project's Mendix version.
// Loaded lazily on first access.
func (r *Reader) SchemaRegistry() *schema.Registry
```

The executor passes this registry to the DESCRIBE and UPDATE code paths.

### Step 4: Extend DESCRIBE PAGE to show scalar properties

**File:** `mdl/executor/cmd_pages_describe_parse.go` (modify)

When parsing any raw widget, extract scalar properties using the schema:

```go
func (e *Executor) extractScalarProperties(w map[string]interface{}) map[string]string {
    registry := e.reader.SchemaRegistry()
    if registry == nil {
        return nil
    }

    widgetType, _ := w["$Type"].(string)
    props := registry.ScalarProperties(widgetType)
    if len(props) == 0 {
        return nil
    }

    result := make(map[string]string)
    typeDef := registry.LookupByStorage(widgetType)

    for _, prop := range props {
        val, ok := w[prop.StorageName]
        if !ok {
            continue
        }
        // Skip if value equals default
        if typeDef != nil {
            if defaultVal, hasDefault := typeDef.DefaultSettings[prop.Name]; hasDefault {
                if fmt.Sprintf("%v", val) == fmt.Sprintf("%v", defaultVal) {
                    continue
                }
            }
        }
        result[prop.StorageName] = fmt.Sprintf("%v", val)
    }
    return result
}
```

This is called generically for **every** widget type in `parseRawWidget`, not just LayoutGrid. Any widget with non-default scalar properties gets them displayed.

**File:** `mdl/executor/cmd_pages_describe_output.go` (modify)

In the output for each widget type, include scalar properties:

```go
// For LayoutGrid, DataView, Container, etc. — all widget types
if len(w.BuiltinProps) > 0 {
    for key, val := range w.BuiltinProps {
        props = append(props, fmt.Sprintf("%s: %s", key, val))
    }
}
```

After this, `DESCRIBE PAGE` automatically shows:
```sql
LAYOUTGRID layoutgrid1 (Width: FixedWidth) {
  ROW row1 { ... }
}
```

Width only appears when non-default. If a DataView has `Editability: Never`, that also shows. No per-widget-type code needed.

### Step 5: Raw BSON write path for UPDATE WIDGETS

**File:** `mdl/executor/widget_property_raw.go` (new)

```go
// walkRawWidgets walks the raw BSON widget tree of a page/snippet,
// calling visitor for each widget map.
func walkRawWidgets(rawDoc map[string]interface{}, visitor func(widget map[string]interface{}) error) error

// setRawWidgetProperty sets a property on a raw widget using schema validation.
func setRawWidgetProperty(registry *schema.Registry, widget map[string]interface{}, propertyName string, value interface{}) error
```

The setter validates against the schema:

```go
func setRawWidgetProperty(registry *schema.Registry, widget map[string]interface{}, propertyName string, value interface{}) error {
    widgetType, _ := widget["$Type"].(string)

    // Pluggable widgets: existing Object.Properties[] path (no schema needed)
    if widgetType == "CustomWidgets$CustomWidget" {
        return setRawPluggableWidgetProperty(widget, propertyName, value)
    }

    // Built-in widgets: validate against schema
    typeDef := registry.LookupByStorage(widgetType)
    if typeDef == nil {
        return fmt.Errorf("unknown widget type: %s", widgetType)
    }

    // Find property by storage name (case-insensitive)
    var prop *schema.PropertyDef
    for _, p := range typeDef.Properties {
        if strings.EqualFold(p.StorageName, propertyName) {
            prop = p
            break
        }
    }
    if prop == nil {
        // Suggest available scalar properties
        scalars := registry.ScalarProperties(widgetType)
        var names []string
        for _, s := range scalars {
            names = append(names, s.StorageName)
        }
        return fmt.Errorf("unknown property '%s' on %s (settable: %s)",
            propertyName, widgetType, strings.Join(names, ", "))
    }

    // Only allow scalar types
    if prop.TypeInfo.Type != schema.TypeInfoEnumeration && prop.TypeInfo.Type != schema.TypeInfoPrimitive {
        return fmt.Errorf("property '%s' is a complex type (%s) and cannot be set via UPDATE WIDGETS",
            propertyName, prop.TypeInfo.Type)
    }

    // Validate enum values
    if prop.TypeInfo.Type == schema.TypeInfoEnumeration && len(prop.TypeInfo.Values) > 0 {
        strVal := fmt.Sprintf("%v", value)
        matched := false
        for _, allowed := range prop.TypeInfo.Values {
            if strings.EqualFold(allowed, strVal) {
                value = allowed // normalize casing
                matched = true
                break
            }
        }
        if !matched {
            return fmt.Errorf("invalid value '%v' for %s (allowed: %s)",
                value, prop.StorageName, strings.Join(prop.TypeInfo.Values, ", "))
        }
    }

    // Coerce type to match existing BSON value
    if existing, ok := widget[prop.StorageName]; ok {
        widget[prop.StorageName] = coerceToExistingType(value, existing)
    } else {
        widget[prop.StorageName] = value
    }
    return nil
}
```

### Step 6: Switch UPDATE WIDGETS to raw BSON path

**File:** `mdl/executor/cmd_widgets.go` (modify `updateWidgetsInPage`, `updateWidgetsInSnippet`)

```go
func (e *Executor) updateWidgetsInPage(...) (int, error) {
    // Read raw BSON (preserves ALL fields)
    rawBytes, err := e.reader.GetRawUnitBytes(model.ID(containerID))
    var rawDoc map[string]interface{}
    bson.Unmarshal(rawBytes, &rawDoc)

    registry := e.reader.SchemaRegistry()

    // Walk raw widgets, update matches
    walkRawWidgets(rawDoc, func(widget map[string]interface{}) error {
        widgetID := extractBinaryID(widget["$ID"])
        if ref, ok := widgetIDs[widgetID]; ok {
            for _, a := range assignments {
                if dryRun {
                    fmt.Fprintf(e.output, "  Would set '%s' = %v on %s in %s\n", ...)
                } else {
                    err := setRawWidgetProperty(registry, widget, a.PropertyPath, a.Value)
                    if err != nil { ... }
                }
            }
            updated++
        }
        return nil
    })

    // Write back
    if !dryRun && updated > 0 {
        data, _ := bson.Marshal(rawDoc)
        e.writer.UpdateRawUnit(model.ID(containerID), data)
    }
}
```

**File:** `sdk/mpr/writer_page.go` (add)

```go
// UpdateRawUnit writes raw BSON bytes back to a unit.
func (w *Writer) UpdateRawUnit(id model.ID, data []byte) error {
    return w.updateUnit(string(id), data)
}
```

## User's Scenario After Implementation

```sql
-- 1. See current values — Width: FixedWidth appears because it's non-default
DESCRIBE PAGE Main.AddressType_Overview;
-- LAYOUTGRID layoutgrid1 (Width: FixedWidth) { ... }

-- 2. Preview bulk change
UPDATE WIDGETS SET 'Width' = 'FullWidth'
WHERE WidgetType LIKE '%LayoutGrid%' IN Main DRY RUN;
-- Found 51 widget(s)... Would set 'Width' = FullWidth on layoutgrid1 in ...

-- 3. Apply
UPDATE WIDGETS SET 'Width' = 'FullWidth'
WHERE WidgetType LIKE '%LayoutGrid%' IN Main;
-- Updated 51 widget(s)

-- 4. Verify — Width disappears (it's now the default)
DESCRIBE PAGE Main.AddressType_Overview;
-- LAYOUTGRID layoutgrid1 { ... }
```

Validation is automatic:
```sql
UPDATE WIDGETS SET 'Width' = 'Narrow'
WHERE WidgetType LIKE '%LayoutGrid%' IN Main;
-- Error: invalid value 'Narrow' for Width (allowed: FixedWidth, FullWidth)

UPDATE WIDGETS SET 'Appearance' = 'something'
WHERE WidgetType LIKE '%LayoutGrid%' IN Main;
-- Error: property 'Appearance' is a complex type (ELEMENT) and cannot be set via UPDATE WIDGETS
```

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `sdk/schema/registry.go` | **Create** | Runtime schema registry (Phase 1 of BSON Schema Registry) |
| `sdk/schema/registry_test.go` | **Create** | Tests for schema loading and property lookup |
| `sdk/schema/reflection-data/` | **Create** | Embedded reflection data for ~15 key versions |
| `internal/codegen/schema/types.go` | **Move** | Move types to `sdk/schema/types.go`, import from codegen |
| `mdl/executor/widget_property_raw.go` | **Create** | Raw BSON widget walker + schema-validated setter |
| `mdl/executor/widget_property_raw_test.go` | **Create** | Tests for raw BSON operations |
| `mdl/executor/cmd_widgets.go` | **Modify** | Switch to raw BSON path for all widget updates |
| `mdl/executor/cmd_pages_describe.go` | **Modify** | Add `BuiltinProps` to `rawWidget` struct |
| `mdl/executor/cmd_pages_describe_parse.go` | **Modify** | Generic scalar property extraction via schema |
| `mdl/executor/cmd_pages_describe_output.go` | **Modify** | Display scalar properties for all widget types |
| `sdk/mpr/reader_documents.go` | **Modify** | Add lazy `SchemaRegistry()` accessor |
| `sdk/mpr/writer_page.go` | **Modify** | Add `UpdateRawUnit` method |

No grammar or AST changes needed.

## Alignment with BSON Schema Registry Proposal

This proposal implements a **narrow vertical slice of Phase 1** of the [BSON Schema Registry](BSON_SCHEMA_REGISTRY_PROPOSAL.md):

| Schema Registry Phase | This Proposal |
|---|---|
| **Phase 1**: Registry core — load reflection data, resolve storage names, type lookup | **Yes** — `sdk/schema/Registry` with `LoadRegistry(version)`, `LookupByStorage()`, `ScalarProperties()` |
| **Phase 2**: Write-side completion (fill missing defaults) | Not yet — but the registry has `DefaultSettings` ready |
| **Phase 3**: Widget schema resolution (from .mpk) | Not yet — pluggable widgets continue using existing path |
| **Phase 4**: Generic parser | Not yet — but `ScalarProperties()` is a simplified form |
| **Phase 5**: Version migration | Not yet — but version-aware loading is in place |

The `sdk/schema/` package created here becomes the foundation for later phases. Key design decisions (embedding strategy, version matching, storage name indexing) are made once and reused.

### What we're NOT doing (intentionally)

- **Full field completeness checking** — this proposal only reads/writes individual scalar fields, not entire documents
- **Array encoding / reference kinds** — not needed for scalar property updates
- **Generic parsing of all types** — only extracting scalar values from known BSON fields
- **Version migration** — properties are read and written within the same project version

These are future phases of the Schema Registry, not needed for the widget property use case.

## Risks

1. **Embedding ~18 MB of reflection data** — increases binary size. Mitigation: only embed key versions (~15); compress or lazy-load if needed.
2. **Version mismatch** — project may use a Mendix version not exactly matching any embedded schema. Mitigation: fall back to nearest version; scalar properties rarely change across patches.
3. **Schema loading performance** — parsing 1.2 MB JSON on first use. Mitigation: lazy-load, cache parsed registry; only ~50ms for JSON parse.
4. **Switching UPDATE to raw BSON** — changes existing pluggable widget behavior. Mitigation: pluggable widget raw BSON path already exists and is tested.
