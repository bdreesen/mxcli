# Widget Templates

This directory contains JSON templates for Mendix pluggable widgets. These templates are extracted from a reference Mendix project and embedded into the mxcli binary.

## Structure

```
templates/
├── mendix-11.6/                    # Templates for Mendix 11.6.x
│   ├── combobox.json               # com.mendix.widget.web.combobox.Combobox
│   ├── datagrid.json               # com.mendix.widget.web.datagrid.Datagrid
│   ├── datagrid-text-filter.json   # DatagridTextFilter
│   ├── datagrid-date-filter.json   # DatagridDateFilter
│   ├── datagrid-dropdown-filter.json
│   └── datagrid-number-filter.json
├── mendix-10.x/                    # Templates for older versions (if needed)
└── README.md
```

## Template Format

Each template is a JSON file containing **both** the `CustomWidgetType` and `WidgetObject` structures:

```json
{
  "widgetId": "com.mendix.widget.web.combobox.Combobox",
  "name": "Combo box",
  "version": "11.6.0",
  "extractedFrom": "PageTemplates.Customer_NewEdit",
  "type": { ... },   // The full CustomWidgetType BSON converted to JSON
  "object": { ... }  // The default WidgetObject with all property values
}
```

### Why Both `type` AND `object` Are Required

The `type` field defines the widget's PropertyTypes (schema), while the `object` field contains the actual property values with correct defaults. Studio Pro expects:

1. **Consistent IDs**: `object.Properties[].TypePointer` must reference valid `type.ObjectType.PropertyTypes[].$ID` values
2. **All properties present**: Every PropertyType must have a corresponding WidgetProperty in the object
3. **Correct default values**: Properties like `TextTemplate` need proper `Forms$ClientTemplate` structures, not null

Without the `object` field, mxcli must build the WidgetObject from scratch, which is error-prone and often triggers "widget definition has changed" warnings in Studio Pro.

## Extracting New Templates

### Important: Use Studio Pro-Created Widgets

When extracting templates, **always use widgets that have been created or "fixed" by Studio Pro**. This ensures the WidgetObject contains correct default values. If you programmatically create a widget and extract it, you'll just get the same incorrect structure back.

### Extraction Process

1. **Create the widget in Studio Pro** - Add the widget to a page in Studio Pro and configure it with default settings

2. **If updating an existing template** - If Studio Pro shows "widget definition has changed", right-click and select "Update widget" to let Studio Pro fix it

3. **Extract using mxcli**:
```bash
# Extract BSON template + skeleton .def.json from .mpk widget package
mxcli widget extract --mpk path/to/widget.mpk

# Generates:
#   .mxcli/widgets/<widget-name>.json      (template with type + object)
#   .mxcli/widgets/<widget-name>.def.json  (skeleton definition)
```

4. **Manual extraction** (current method):
```go
// Use reader.GetRawUnit() to get the page, then extract CustomWidget.Type and CustomWidget.Object
// Convert BSON binary IDs to hex strings for JSON storage
```

### Verifying Templates

After updating a template, verify it works:

```bash
# Create a test page with the widget
mxcli -p test.mpr -c "CREATE PAGE Test.TestPage ... DATAGRID ..."

# Check for errors (should have no CE0463 errors)
mx check test.mpr
```

## Usage

Templates are automatically used when creating pluggable widgets via MDL:

```sql
COMBOBOX myCombo ATTRIBUTE Country;
```

### Priority Chain (3-Tier Widget Registry)

When creating a pluggable widget, mxcli resolves definitions and templates using a 3-tier registry:

1. **Embedded** (`sdk/widgets/definitions/*.def.json` + `sdk/widgets/templates/`) — Built-in definitions, compiled into the binary
2. **Global** (`~/.mxcli/widgets/*.def.json` + `*.json`) — User-defined global overrides
3. **Project** (`<project>/.mxcli/widgets/*.def.json` + `*.json`) — Per-project overrides (highest priority)

Each `.def.json` declares property mappings and child slots; the engine applies them to the BSON template at build time. See `docs/plans/2026-03-25-pluggable-widget-engine-design.md` for the full architecture.

### Why Templates Are Needed

Mendix pluggable widgets (like ComboBox, DataGrid2) require a full `CustomWidgetType` definition with 50+ PropertyTypes. These definitions are embedded in each widget instance in the MPR file. Without the complete definition, Mendix will show "widget definition has changed" warnings.

By embedding templates extracted from a known-good project, mxcli can create widgets that are fully compatible with Mendix Studio Pro.

### Known Limitation: Widget Version Drift

Static templates are tied to the widget version they were extracted from. If the target project has a **newer** version of the widget `.mpk` (in `widgets/`), Studio Pro will detect that the serialized Type definition doesn't match the installed widget and report CE0463.

For example, the ComboBox template was extracted from a Mendix 11.6.0 project, but a 11.6.3 project may ship ComboBox v2.5.0 which added 3 new properties (`staticDataSourceCaption`, `staticDataSourceCustomContent`, `staticDataSourceValue`). Our template lacks these → CE0463.

**The correct long-term fix**: read the widget definition from the project's actual `widgets/*.mpk` file at runtime instead of relying on static templates. The `.mpk` is a ZIP containing an XML schema (e.g., `Combobox.xml`) that defines all property keys, types, and defaults. Two approaches:

1. **Parse `.mpk` XML, generate full BSON** — map each XML property type (`attribute`, `expression`, `widgets`, `textTemplate`, etc.) to the BSON structure with correct defaults. Eliminates version drift entirely.
2. **Augment static template from `.mpk` at runtime** — keep the current template for BSON structure patterns, but read the `.mpk` XML to discover which properties should exist, adding missing ones and removing stale ones.

Either way, the `.mpk` in the project's `widgets/` folder is the **source of truth** for what properties a widget should have.

## TextTemplate Property Requirements

Properties with `"Type": "TextTemplate"` in the Type definition require special handling. They cannot be `null` in the Object section.

### Problem: CE0463 "widget definition has changed"

If a TextTemplate property is `null` in the Object section, Studio Pro shows:
```
CE0463: The definition of this widget has changed. Update this widget...
```

### Required Structure

TextTemplate properties must have a proper `Forms$ClientTemplate` structure:

```json
"TextTemplate": {
  "$ID": "<32-char-guid>",
  "$Type": "Forms$ClientTemplate",
  "Fallback": {
    "$ID": "<32-char-guid>",
    "$Type": "Texts$Text",
    "Items": []
  },
  "Parameters": [],
  "Template": {
    "$ID": "<32-char-guid>",
    "$Type": "Texts$Text",
    "Items": []
  }
}
```

### Important: Empty Arrays

Empty arrays must be `[]`, NOT `[2]`:
```json
// WRONG - serializes as array containing integer 2
"Items": [2]

// CORRECT - truly empty array
"Items": []
```

### How to Identify TextTemplate Properties

1. Search the Type section for `"Type": "TextTemplate"`
2. Note the `$ID` from the parent `ValueType` object
3. Find Object properties where `Value.TypePointer` matches that ID
4. Update those properties' `TextTemplate` from `null` to proper structure

### Affected Widgets

Filter widgets commonly have TextTemplate properties:
- **TextFilter**: `placeholder`, `screenReaderButtonCaption`, `screenReaderInputCaption`
- **DateFilter**: `placeholder`, `screenReaderButtonCaption`, `screenReaderCalendarCaption`, `screenReaderInputCaption`
- **DropdownFilter**: `emptyOptionCaption`, `ariaLabel`, `emptySelectionCaption`, `filterInputPlaceholderCaption`
- **NumberFilter**: `placeholder`, `screenReaderButtonCaption`, `screenReaderInputCaption`
