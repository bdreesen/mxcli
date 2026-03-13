# Page Syntax V2 - Implementation Reference

## Status: Superseded by V3 ⚠️

> **Note:** V3 syntax is now the recommended approach. See `proposal_pages_v3.md`.
> V2 syntax is still supported for backward compatibility.

This document describes the Page Syntax V2 that has been implemented in MDL.

## Overview

Page Syntax V2 provides a more consistent, readable, and flexible syntax for creating pages and widgets in MDL:

- **`{ }` blocks** instead of `BEGIN/END` for cleaner nesting
- **`->` binding operator** for attributes, variables, actions, and datasources
- **`(Name: value)` property syntax** matching Studio Pro property names
- **Consistent widget pattern**: `WIDGET [id] ['label'] [-> binding] [(properties)] [{ children }]`

## Syntax Pattern

```
WIDGETTYPE [id] ['label'] [-> binding] [(properties)] [{ children }]
```

| Part | Required | Description |
|------|----------|-------------|
| `WIDGETTYPE` | Yes | Widget type keyword (TEXTBOX, DATAVIEW, etc.) |
| `id` | No | Widget identifier for referencing |
| `'label'` | No | Display label (positional string) |
| `-> binding` | No | Binding target (attribute, variable, action, datasource) |
| `(properties)` | No | Additional properties in parentheses |
| `{ children }` | No | Nested widgets in braces |

## Binding Operator `->`

The `->` operator provides a clear, unified way to express bindings:

```mdl
-- Attribute binding (form widgets)
TEXTBOX 'Name' -> Name
CHECKBOX 'Active' -> IsActive

-- Variable binding (containers)
DATAVIEW dvProduct -> $Product { ... }

-- Action binding (buttons)
ACTIONBUTTON 'Save' -> SAVE_CHANGES
ACTIONBUTTON 'Process' -> MICROFLOW MyModule.ProcessOrder

-- Database source with query
DATAGRID -> DATABASE MyModule.Product
  WHERE [IsActive = true]
  ORDER BY Name ASC

-- Selection binding (master-detail)
DATAVIEW -> SELECTION galleryName
```

## Complete Example

```mdl
CREATE PAGE MyModule.ProductDetail (
  $Product: MyModule.Product
)
  TITLE 'Product Details'
  LAYOUT Atlas_Core.Atlas_Default
{
  LAYOUTGRID {
    ROW {
      COLUMN (DesktopWidth: 12) {
        DYNAMICTEXT 'Product: {1}' WITH ({1} = $Product/Name) (RenderMode: H3)
      }
    }

    ROW {
      COLUMN (DesktopWidth: 6) {
        DATAVIEW dvProduct -> $Product {
          TEXTBOX 'Name' -> Name
          TEXTBOX 'Code' -> Code
          TEXTAREA 'Description' -> Description
          DATEPICKER 'Created' -> CreatedDate
          CHECKBOX 'Active' -> IsActive
          DROPDOWN 'Status' -> Status

          FOOTER {
            ACTIONBUTTON 'Save' -> SAVE_CHANGES (ButtonStyle: Primary)
            ACTIONBUTTON 'Cancel' -> CLOSE_PAGE
            ACTIONBUTTON 'Process' -> MICROFLOW MyModule.ACT_ProcessProduct (
              Product: $Product
            ) (ButtonStyle: Success)
          }
        }
      }

      COLUMN (DesktopWidth: 6) {
        DATAGRID dgRelated -> DATABASE MyModule.RelatedItem
          WHERE [MyModule.RelatedItem_Product = $Product]
          ORDER BY Name ASC
        {
          ControlBar: {
            ACTIONBUTTON 'New' -> CREATE_OBJECT MyModule.RelatedItem
              THEN SHOW_PAGE MyModule.RelatedItem_Edit
              (ButtonStyle: Primary)
          }
          Columns: {
            COLUMN 'Item Name' -> Name
            COLUMN 'Category' -> Category
            COLUMN 'Price' -> Price
          }
        }
      }
    }
  }
}
```

## Property Syntax

Properties use `Name: value` syntax with colons:

```mdl
(
  PropertyName: value,
  PropertyName: 'string value',
  PropertyName: 123,
  PropertyName: true
)
```

### Body Properties

Complex widgets can have property groups in their body:

```mdl
DATAGRID dgProducts -> DATABASE MyModule.Product {
  ControlBar: {
    ACTIONBUTTON 'New' -> CREATE_OBJECT MyModule.Product
      THEN SHOW_PAGE MyModule.Product_Edit
      (ButtonStyle: Primary)
  }
  Columns: {
    COLUMN 'Name' -> Name
    COLUMN 'Price' -> Price
    COLUMN 'Actions' {
      ACTIONBUTTON 'Edit' -> SHOW_PAGE MyModule.Product_Edit
        (Product: $currentObject) (ButtonStyle: Default)
    }
  }
}
```

## Implemented Binding Types

| Binding | Syntax | Use Case |
|---------|--------|----------|
| Attribute | `-> Name` | Form inputs bound to entity attribute |
| Variable | `-> $Product` | DataView/Container bound to page parameter |
| Save | `-> SAVE_CHANGES [CLOSE_PAGE]` | Save button action |
| Cancel | `-> CANCEL_CHANGES [CLOSE_PAGE]` | Cancel button action |
| Close | `-> CLOSE_PAGE` | Close page action |
| Delete | `-> DELETE` | Delete object action |
| Show Page | `-> SHOW_PAGE Module.Page (params)` | Navigate to page |
| Microflow | `-> MICROFLOW Module.MF (params)` | Call microflow |
| Create+Show | `-> CREATE_OBJECT Entity THEN SHOW_PAGE Page` | Create and navigate |
| Database | `-> DATABASE Entity WHERE [...] ORDER BY` | Grid/Gallery data source |
| Selection | `-> SELECTION widgetName` | Master-detail binding |

## Implemented Widgets

### Container Widgets
- `LAYOUTGRID` with `ROW` and `COLUMN`
- `CONTAINER`
- `NAVIGATIONLIST` with `ITEM`
- `FOOTER`

### Data Widgets
- `DATAVIEW` - Single object form
- `DATAGRID` - Data table (DataGrid2 widget)
- `GALLERY` - Card gallery with selection
- `LISTVIEW` - Simple list

### Input Widgets
- `TEXTBOX` - Single-line text input
- `TEXTAREA` - Multi-line text input
- `CHECKBOX` - Boolean checkbox
- `RADIOBUTTONS` - Radio button group
- `DATEPICKER` - Date/time picker
- `DROPDOWN` - Dropdown (deprecated, use COMBOBOX)
- `COMBOBOX` - Combo box (pluggable widget)

### Display Widgets
- `DYNAMICTEXT` - Dynamic text with optional template
- `TITLE` - Page heading
- `TEXT` - Static text

### Action Widgets
- `ACTIONBUTTON` - Button with action
- `LINKBUTTON` - Link-styled button

### Special Widgets
- `SNIPPETCALL` - Embed snippet
- `TEMPLATE` - Template for Gallery/ListView items
- `FILTER` with `TEXTFILTER` - Gallery filter

## Property Name Mapping

Property names match Studio Pro for familiarity:

| MDL Property | Widget Types |
|--------------|--------------|
| `Content` | Text widgets |
| `RenderMode` | Text widgets (H1-H6, Paragraph, Text) |
| `ButtonStyle` | Buttons (Default, Primary, Success, etc.) |
| `DesktopWidth` | Columns (1-12, AutoFill, AutoFit) |
| `Selection` | Gallery (Single, Multiple, None) |
| `Class` | All widgets |

## Implementation Files

| File | Purpose |
|------|---------|
| `mdl/grammar/MDLLexer.g4` | Tokens for `{`, `}`, `->`, `:`, keywords |
| `mdl/grammar/MDLParser.g4` | Widget and binding grammar rules |
| `mdl/ast/ast_page_v2.go` | AST types for V2 widgets |
| `mdl/visitor/visitor_page_v2.go` | Parser to AST conversion |
| `mdl/executor/cmd_pages_builder.go` | Page building with V2 support |
| `mdl/executor/cmd_pages_builder_widgets_v2.go` | V2 widget builders |

## Backward Compatibility

Both syntaxes are supported:
- **V1 (BEGIN/END)**: Still works, primarily used by DESCRIBE output
- **V2 ({ } and ->)**: New recommended syntax for writing pages

The parser accepts both syntaxes, and the executor handles both AST formats.

## Example Files

See `mdl-examples/doctype-tests/03-page-examples-v2.mdl` for comprehensive examples including:
- Empty pages
- Layout grids with dynamic text
- DataViews with form inputs
- DataGrids with control bars and columns
- Galleries with filters and templates
- Master-detail patterns
- Snippets with parameters
- String templates with WITH syntax

## Verification

```bash
# Parse V2 syntax
./bin/mxcli check mdl-examples/doctype-tests/03-page-examples-v2.mdl

# Execute against project
./bin/mxcli -p app.mpr -c "execute script 'mdl-examples/doctype-tests/03-page-examples-v2.mdl'"

# Verify in Studio Pro
reference/mxbuild/modeler/mx check app.mpr
```
