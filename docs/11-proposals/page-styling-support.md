# Proposal: Add Styling Support to MDL Pages

## Current State

Mendix has four styling mechanisms on every widget, all stored in a BSON `Forms$Appearance` object:

| Mechanism | BSON Field | Type | Example |
|---|---|---|---|
| CSS classes | `Class` | string | `"btn-lg mx-spacing-top-large"` |
| Inline CSS | `Style` | string | `"color: red; margin: 10px;"` |
| Dynamic classes | `DynamicClasses` | string (XPath) | `"if $currentObject/IsActive then 'highlight' else ''"` |
| Design properties | `DesignProperties` | array of typed tokens | Atlas UI spacing, colors, toggles |

Additionally, some widgets have their own style enums (e.g. `ButtonStyle: Primary` on buttons — already supported as a separate keyword).

**What exists today:** The serializer (`serializeAppearance()` in `writer_widgets.go`) writes the `Forms$Appearance` BSON structure but always with empty values. `BaseWidget` has `Class` and `Style` fields in the Go struct but they're never populated from MDL. DESCRIBE doesn't output any styling. The `ButtonStyle:` keyword was renamed from `Style:` to free up `Style:` for CSS inline styling.

## Proposal: Three Phases

### Phase 1 — `Class` and `Style` properties (highest value) ✅ DONE

`Class` and `Style` are standard widget properties in the V3 syntax:

```sql
TEXTBOX txtName (Label: 'Name', Attribute: Name, Class: 'form-control-lg mx-spacing-top-large')
CONTAINER ctn1 (Class: 'card', Style: 'padding: 16px; border-radius: 8px;') {
  DYNAMICTEXT txt1 (Content: '{1}', Params: [FullName], Class: 'text-primary h3')
}
ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES, ButtonStyle: Primary, Class: 'btn-block')
```

This is the most common way developers style Mendix apps and maps directly to existing fields on `BaseWidget`. Implementation:

- **Grammar**: `CLASS` and `STYLE` keywords in `widgetPropertyV3` rule (MDLParser.g4)
- **AST**: `GetClass()` and `GetStyle()` helpers on `WidgetV3` (ast_page_v3.go)
- **Visitor**: Extracts Class/Style from parsed string literals (visitor_page_v3.go)
- **Builder**: `applyWidgetAppearance()` sets Class/Style on any widget via `SetAppearance()` (cmd_pages_builder_v3.go)
- **Serializer**: `serializeAppearance(class, style)` passes values through to BSON (writer_widgets.go + all callers)
- **Describe Parse**: Extracts Class/Style from `Appearance` BSON (cmd_pages_describe_parse.go)
- **Describe Output**: Emits `Class:` and `Style:` when non-empty via `appendAppearanceProps()` (cmd_pages_describe_output.go)
- **Wireframe**: Class/Style fields on `wireframeNode` (cmd_page_wireframe.go)

### Phase 2 — `DynamicClasses` (conditional styling)

```sql
CONTAINER ctn1 (
  Class: 'card',
  DynamicClasses: "if $currentObject/Status = 'Error' then 'alert-danger' else 'alert-info'"
) {
  ...
}
```

This is an XPath expression string. Implementation is similar to Phase 1 — just another string property — but it's less commonly used and the XPath expressions can be complex.

### Phase 3 — Design Properties (Atlas UI tokens)

Design properties are structured arrays of typed key-value pairs set via Atlas UI's design system.

#### Where Design Property Definitions Live

Design property **definitions** (what properties are available per widget type) are NOT in the MPR. They live in the project's `themesource` folder:

```
<project-root>/themesource/<Module>/<platform>/design-properties.json
```

The primary file for most Mendix apps:
```
themesource/atlas_core/web/design-properties.json      -- Web platform
themesource/atlas_core/native/design-properties.json    -- Native mobile
```

Any module can add its own design properties at `themesource/<YourModule>/web/design-properties.json`. Multiple modules' properties are merged at load time.

#### Design Property Definition Format

The `design-properties.json` is a JSON object where keys are widget type names and values are arrays of property definitions:

```json
{
    "Widget": [
        { "name": "Spacing top", "type": "Dropdown", "description": "...",
          "options": [
            { "name": "None", "class": "spacing-outer-top-none" },
            { "name": "Small", "class": "spacing-outer-top" },
            { "name": "Large", "class": "spacing-outer-top-large" }
          ]
        },
        { "name": "Hide on phone", "type": "Toggle", "description": "...",
          "class": "hide-phone" }
    ],
    "DivContainer": [
        { "name": "Align content", "type": "Dropdown", "description": "...",
          "options": [
            { "name": "Left align as a row", "class": "row-left" },
            { "name": "Center align as a row", "class": "row-center" }
          ]
        },
        { "name": "Background color", "type": "Dropdown", "description": "...",
          "options": [
            { "name": "Brand Primary", "class": "background-primary" },
            { "name": "Brand Inverse", "class": "background-inverse" }
          ]
        }
    ],
    "Button": [
        { "name": "Size", "type": "Dropdown", ... },
        { "name": "Full width", "type": "Toggle", "class": "btn-block", ... },
        { "name": "Border", "type": "Toggle", "class": "btn-bordered", ... }
    ],
    "com.mendix.widget.web.accordion.Accordion": [ ... ]
}
```

**Widget type keys**: Built-in widgets use Model SDK class names (`DivContainer`, `Button`, `DataGrid`, `ListView`, `DynamicText`). Pluggable widgets use their widget ID (`com.mendix.widget.web.accordion.Accordion`). The `Widget` key defines properties that apply to ALL widgets (inherited by all subtypes).

**Five property types**:

| Type | Description | Value stored in MPR |
|------|-------------|---------------------|
| `Toggle` | On/off CSS class | `ToggleDesignPropertyValue` (presence = on) |
| `Dropdown` | Single-select option list | `OptionDesignPropertyValue` (option name string) |
| `ColorPicker` | Dropdown with color preview | `OptionDesignPropertyValue` (same storage as Dropdown) |
| `ToggleButtonGroup` | Related options as buttons | `OptionDesignPropertyValue` (single) or `CompoundDesignPropertyValue` (multi) |
| `Spacing` | Margin/padding in 4 directions | `CompoundDesignPropertyValue` (nested properties) |

#### BSON Storage Format

Design property **values** are stored per widget in `Appearance.DesignProperties` as an array of `Forms$DesignPropertyValue` objects:

```bson
"Appearance": {
  "$Type": "Forms$Appearance",
  "Class": "",
  "Style": "",
  "DynamicClasses": "",
  "DesignProperties": [
    2,
    {
      "$Type": "Forms$DesignPropertyValue",
      "Key": "Spacing top",
      "Value": {
        "$Type": "Forms$OptionDesignPropertyValue",
        "Option": "Large"
      }
    },
    {
      "$Type": "Forms$DesignPropertyValue",
      "Key": "Full width",
      "Value": {
        "$Type": "Forms$ToggleDesignPropertyValue"
      }
    },
    {
      "$Type": "Forms$DesignPropertyValue",
      "Key": "Spacing",
      "Value": {
        "$Type": "Forms$CompoundDesignPropertyValue",
        "Properties": [
          2,
          {
            "$Type": "Forms$DesignPropertyValue",
            "Key": "Top",
            "Value": {
              "$Type": "Forms$OptionDesignPropertyValue",
              "Option": "M"
            }
          }
        ]
      }
    }
  ]
}
```

**Value type hierarchy**:
- `Forms$ToggleDesignPropertyValue` — empty struct, presence means enabled
- `Forms$OptionDesignPropertyValue` — has `Option` string (the selected option name)
- `Forms$CustomDesignPropertyValue` — has `Value` string (arbitrary text)
- `Forms$CompoundDesignPropertyValue` — has nested `Properties` array of `DesignPropertyValue`

#### Inline Syntax in CREATE PAGE

Design properties can be included in `CREATE PAGE` using an explicit `DesignProperties:` array on any widget. This keeps them clearly separated from built-in widget properties (Class, Style, Label, Binds, etc.):

```sql
CREATE PAGE MyModule.Customer_Edit
(
  Params: { $Customer: MyModule.Customer },
  Title: 'Edit Customer',
  Layout: Atlas_Core.PopupLayout
)
{
  DATAVIEW dvCustomer (DataSource: $Customer) {
    CONTAINER ctn1 (
      Class: 'card',
      Style: 'padding: 16px;',
      DesignProperties: [
        'Spacing top': 'Large',
        'Background color': 'Brand Primary',
        'Hide on phone': ON
      ]
    ) {
      TEXTBOX txtName (Label: 'Name', Attribute: Name)
      TEXTBOX txtEmail (Label: 'Email', Attribute: Email)
    }

    FOOTER footer1 {
      ACTIONBUTTON btnSave (
        Caption: 'Save',
        Action: SAVE_CHANGES,
        ButtonStyle: Primary,
        DesignProperties: ['Full width': ON, 'Size': 'Large']
      )
      ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
    }
  }
}
```

**Syntax rules**:
- `DesignProperties:` is followed by `[` ... `]` (array brackets, same pattern as `Params:`)
- Each entry is `STRING_LITERAL COLON (STRING_LITERAL | ON | OFF)` — quoted key, colon, value
- Toggle properties use `ON` / `OFF` keywords
- Dropdown/ColorPicker/ToggleButtonGroup properties use quoted option name strings
- Entries are comma-separated; single-line for few properties, multi-line for many
- The array is optional — omitting `DesignProperties:` means no design properties (same as today)

**Why explicit `DesignProperties:` wrapper** (not bare quoted keys in the property block):
1. **No ambiguity** — a design property named `Style` or `Label` won't collide with built-in keywords
2. **Cleaner grammar** — one rule for `DESIGNPROPERTIES COLON LBRACKET ... RBRACKET` instead of a catch-all
3. **Explicit intent** — reading the MDL you immediately know these are theme-driven, not built-in settings
4. **Simpler builder logic** — everything inside `DesignProperties:` needs theme registry lookup; everything outside doesn't

**Serialization dependency**: The builder must read the project's `themesource/*/web/design-properties.json` to determine the BSON value type for each property. For example, `'Full width': ON` serializes as `Forms$ToggleDesignPropertyValue` (empty struct), while `'Size': 'Large'` serializes as `Forms$OptionDesignPropertyValue` with `Option: "Large"`. Without the theme registry, the builder cannot distinguish between property types.

**DESCRIBE roundtrip**: `DESCRIBE PAGE` outputs `DesignProperties: [...]` on any widget that has non-empty design properties, making the output re-executable.

#### Proposed Commands — Styling Fragments

In addition to inline syntax in `CREATE PAGE`, design properties can be managed through **fragment-style commands** that operate on individual widgets within existing pages. This is useful for modifying styling without rewriting entire pages.

##### 1. Discover Available Design Properties

Read the `design-properties.json` from the project's themesource and show what properties are available for a given widget type:

```sql
-- Show all design properties available for Container widgets
SHOW DESIGN PROPERTIES FOR CONTAINER;

-- Output:
-- From: Widget (inherited)
--   Spacing top          Dropdown    [None, Small, Medium, Large]
--   Spacing bottom       Dropdown    [None, Small, Medium, Large]
--   Hide on phone        Toggle      class: hide-phone
--   Hide on tablet       Toggle      class: hide-tablet
-- From: DivContainer
--   Align content        Dropdown    [Left align as a row, Center align as a row, ...]
--   Background color     Dropdown    [Brand Default, Brand Primary, Brand Inverse, ...]

-- Show available properties for a pluggable widget
SHOW DESIGN PROPERTIES FOR DATAGRID2;

-- Show available properties for all widget types
SHOW DESIGN PROPERTIES;
```

This requires reading `themesource/*/web/design-properties.json` from the project directory and mapping widget type keys to MDL widget types.

##### 2. Describe Current Styling on a Widget

Show all styling (Class, Style, DynamicClasses, and DesignProperties) for a specific widget on a page:

```sql
-- Show styling for a specific widget on a page
DESCRIBE STYLING ON PAGE MyModule.CustomerEdit WIDGET btnSave;

-- Output:
-- WIDGET btnSave (ActionButton)
--   Class: 'btn-block'
--   Style: 'margin-top: 8px;'
--   Design Properties:
--     Spacing top: Large
--     Full width: ON

-- Show styling for ALL widgets on a page
DESCRIBE STYLING ON PAGE MyModule.CustomerEdit;

-- Output (one section per styled widget):
-- WIDGET ctn1 (Container)
--   Class: 'card'
--   Design Properties:
--     Background color: Brand Primary
-- WIDGET btnSave (ActionButton)
--   Class: 'btn-block'
--   Design Properties:
--     Spacing top: Large
```

##### 3. Set Styling on a Single Widget (Fragment Update)

Change styling properties on a single widget without rewriting the entire page:

```sql
-- Set Class and Style on a widget
ALTER STYLING ON PAGE MyModule.CustomerEdit WIDGET btnSave
  SET Class = 'btn-block btn-lg',
      Style = 'margin-top: 16px;';

-- Set a design property (dropdown/colorpicker selection)
ALTER STYLING ON PAGE MyModule.CustomerEdit WIDGET ctn1
  SET 'Spacing top' = 'Large',
      'Background color' = 'Brand Primary';

-- Toggle a design property on
ALTER STYLING ON PAGE MyModule.CustomerEdit WIDGET btnSave
  SET 'Full width' = ON;

-- Toggle a design property off
ALTER STYLING ON PAGE MyModule.CustomerEdit WIDGET btnSave
  SET 'Full width' = OFF;

-- Clear all design properties on a widget
ALTER STYLING ON PAGE MyModule.CustomerEdit WIDGET btnSave
  CLEAR DESIGN PROPERTIES;

-- Mixed: set Class, Style, and design properties together
ALTER STYLING ON PAGE MyModule.CustomerEdit WIDGET ctn1
  SET Class = 'card custom-card',
      Style = 'border-radius: 12px;',
      'Spacing top' = 'Large',
      'Background color' = 'Brand Primary';
```

The `ALTER STYLING` command:
1. Opens the page from the MPR
2. Finds the widget by name (searching the widget tree)
3. Reads the current `Appearance` BSON
4. Modifies only the specified properties (preserving others)
5. Writes the modified `Appearance` back

This is a **surgical update** — it doesn't require parsing or rewriting the full page MDL, just modifying one widget's Appearance blob.

##### 4. Bulk Styling Updates (Extension of UPDATE WIDGETS)

The existing `UPDATE WIDGETS` command could be extended for styling:

```sql
-- Set a design property on all buttons across a module
UPDATE WIDGETS SET 'Full width' = ON
  WHERE WidgetType LIKE '%Button%' IN MyModule DRY RUN;

-- Set Class on all containers
UPDATE WIDGETS SET Class = 'card'
  WHERE WidgetType = 'DivContainer' IN MyModule;
```

#### Widget Type Mapping

The `design-properties.json` keys must be mapped to BSON `$Type` values used in the MPR:

| design-properties.json key | BSON $Type | MDL keyword |
|---------------------------|------------|-------------|
| `Widget` | *(all)* | *(all)* |
| `DivContainer` | `Forms$DivContainer` | `CONTAINER` |
| `Button` / `ActionButton` | `Forms$ActionButton` | `ACTIONBUTTON` |
| `DataGrid` | `Forms$DataGrid` | `DATAGRID` |
| `ListView` | `Forms$ListView` | `LISTVIEW` |
| `DynamicText` | `Forms$DynamicText` | `DYNAMICTEXT` |
| `StaticImageViewer` | `Forms$StaticImageViewer` | `STATICIMAGE` |
| `Label` | `Forms$Label` | `STATICTEXT` |
| `GroupBox` | `Forms$GroupBox` | `GROUPBOX` |
| `TabContainer` | `Forms$TabContainer` | `TABCONTAINER` |
| Pluggable widget ID | `CustomWidgets$CustomWidget` | Widget-specific |

#### Implementation Approach

**Step 1: Theme reader** — Parse `themesource/*/web/design-properties.json` files, merge by widget type, and build an in-memory registry of available properties per widget type. This is a prerequisite for both inline syntax and fragment commands.

**Step 2: Grammar + AST + Visitor** — Add `DESIGNPROPERTIES COLON LBRACKET designPropertyEntry (COMMA designPropertyEntry)* RBRACKET` rule to `widgetPropertyV3` in MDLParser.g4. AST stores design properties as `map[string]string` on `WidgetV3` (key = property name, value = option name or "ON"/"OFF"). Visitor extracts from parse tree.

**Step 3: Builder + Serializer** — In `buildWidgetV3`, read design properties from AST, look up each in the theme registry to determine the BSON value type, and pass to a new `serializeDesignProperties()` function that builds the `Forms$DesignPropertyValue` array inside `Appearance`.

**Step 4: DESCRIBE Parse + Output** — Parse `DesignProperties` array from Appearance BSON, store on `rawWidget`, and emit `DesignProperties: [...]` in MDL output when non-empty.

**Step 5: DESCRIBE STYLING** — Dedicated command to show all styling on a widget or page, with human-readable output cross-referencing the theme registry.

**Step 6: ALTER STYLING** — Locate widget in page BSON by name, modify Appearance in place, write back to MPR. Validates property names and values against theme registry.

**Step 7: SHOW DESIGN PROPERTIES** — Query the theme registry and format as a table, showing available properties per widget type with their allowed values.

#### Open Questions

1. **Snippet styling** — Should `ALTER STYLING` also work on snippets (`ALTER STYLING ON SNIPPET Module.Name WIDGET ...`)?
2. **Compound properties** — Spacing has nested structure (margin-top, margin-bottom, etc.). Should we flatten to `'Spacing.Margin.Top' = 'M'` or use a nested syntax?
3. **Validation** — Should `ALTER STYLING` reject unknown property names / invalid option values, or allow them (for forward compatibility with newer themes)?
4. **Theme discovery** — The MPR path gives us the project root, but we need to verify `themesource/` exists and handle projects without Atlas Core.

## Recommendation

Phase 1 alone covers ~90% of real-world styling needs. It's a small change (the plumbing already exists in `BaseWidget` and `serializeAppearance`) and fits naturally into the existing property syntax. Phase 2 is a straightforward extension. Phase 3 provides two complementary interfaces: inline `DesignProperties: [...]` in `CREATE PAGE` for new pages, and fragment commands (`DESCRIBE STYLING`, `ALTER STYLING`, `SHOW DESIGN PROPERTIES`) for surgical updates to existing pages. Both depend on a theme reader that parses `design-properties.json` from the project's `themesource` folder.
