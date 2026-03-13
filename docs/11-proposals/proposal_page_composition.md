# Page Composition and Partial Updates

## Status: Proposal

## Problem Statement

Large MDL page scripts become unwieldy to write, read, and maintain. Currently, the only way to create or modify a page is to specify the entire page structure in a single `CREATE PAGE` or `CREATE OR REPLACE PAGE` statement. This creates several issues:

1. **Unsafe editing of pages with unsupported widgets** - `CREATE OR REPLACE PAGE` rebuilds the entire page from MDL. Any widget type not yet supported by the MDL writer is **lost**. `DESCRIBE PAGE` renders unsupported widget types as comments (`-- Pages$SomeType (name)`), so round-tripping a page silently drops content. With partial updates (ALTER PAGE), we can add, change, or remove specific widgets while leaving unsupported parts of the page untouched.
2. **No reusability** - Common widget patterns (save/cancel buttons, form layouts) must be copy-pasted
3. **All-or-nothing updates** - Changing a single property requires rewriting the entire page
4. **Large script files** - Complex pages result in hundreds of lines of MDL
5. **No incremental development** - Can't build pages piece by piece in REPL sessions

## Goals

1. **Safe partial editing** - Modify pages containing unsupported widget types without data loss
2. **Composability** - Break large pages into smaller, reusable MDL fragments
3. **Partial Updates** - Modify specific widgets or properties without replacing entire pages
4. **Incremental Creation** - Build pages step by step

## Design Principles

- Widget names are unique within a page (flat namespace, no nested paths needed)
- Fragments are MDL-level constructs (not Mendix Studio Pro snippets)
- Fragments are script-scoped (transient, not persisted)
- Operations should be atomic and validate against current page state
- Property assignments must validate against widget type capabilities
- Syntax should feel natural alongside existing MDL

---

## Relationship to Existing Features

This proposal complements existing features. `UPDATE WIDGETS` (bulk property updates) and `ALTER STYLING ON PAGE` (partial styling updates) are already implemented and prove that partial page modification works.

| Feature | Status | Scope | Target | Use Case |
|---------|--------|-------|--------|----------|
| `UPDATE WIDGETS SET ... WHERE ...` | **Implemented** | Project-wide or module | Multiple widgets by filter | "Disable labels on all comboboxes" |
| `ALTER STYLING ON PAGE ... WIDGET ...` | **Implemented** | Single page | Widget styling by name | "Change CSS class on this container" |
| `ALTER PAGE SET ... ON widget` | Proposed | Single page | Single widget by name | "Change this button's caption" |
| `ALTER PAGE REPLACE widget WITH {...}` | Proposed | Single page | Widget subtree | "Restructure this form section" |
| `DEFINE/USE FRAGMENT` | Proposed | Script scope | Reusable widget groups | "Standard save/cancel buttons" |

**How they work together:**

```sql
-- Bulk update (existing): Change ALL combobox labels project-wide
UPDATE WIDGETS
SET 'showLabel' = false
WHERE WidgetType LIKE '%combobox%'

-- Targeted update (proposed): Change ONE specific widget's property
ALTER PAGE Module.CustomerEdit {
  SET 'showLabel' = true ON cbStatus  -- Exception to the rule above
}

-- Structural change (proposed): Replace entire section
ALTER PAGE Module.CustomerEdit {
  REPLACE footer1 WITH {
    USE FRAGMENT NewFooterLayout
  }
}
```

---

## Part 1: Fragment Definition and Usage

### DEFINE FRAGMENT

Define reusable widget groups that can be inserted into pages:

```sql
-- Simple fragment
DEFINE FRAGMENT SaveCancelFooter AS {
  FOOTER footer1 {
    ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES, ButtonStyle: Primary)
    ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
  }
}

-- Fragment with layout structure
DEFINE FRAGMENT TwoColumnForm AS {
  LAYOUTGRID formGrid {
    ROW row1 {
      COLUMN colLeft (DesktopWidth: 6) { }
      COLUMN colRight (DesktopWidth: 6) { }
    }
  }
}

-- Fragment for consistent headings
DEFINE FRAGMENT PageHeader AS {
  LAYOUTGRID headerGrid {
    ROW headerRow {
      COLUMN headerCol (DesktopWidth: 12) {
        DYNAMICTEXT pageTitle (Content: 'Page Title', RenderMode: H2)
      }
    }
  }
}
```

### USE FRAGMENT

Insert a defined fragment at the current position:

```sql
CREATE PAGE Module.CustomerEdit
(
  Params: { $Customer: Module.Customer },
  Title: 'Edit Customer',
  Layout: Atlas_Core.PopupLayout
)
{
  DATAVIEW dvCustomer (DataSource: $Customer) {
    TEXTBOX txtName (Label: 'Name', Attribute: Name)
    TEXTBOX txtEmail (Label: 'Email', Attribute: Email)

    -- Insert the reusable footer
    USE FRAGMENT SaveCancelFooter
  }
}
```

### Parameterized Fragments (Future)

Future enhancement - fragments that accept parameters:

```sql
DEFINE FRAGMENT FormField($label, $attr) AS {
  TEXTBOX txt$attr (Label: $label, Attribute: $attr)
}

-- Usage
USE FRAGMENT FormField('Customer Name', 'Name')
USE FRAGMENT FormField('Email Address', 'Email')
```

### Fragment Naming and Prefixes

When using fragments multiple times, use a prefix to avoid name conflicts:

```sql
DEFINE FRAGMENT SaveCancelFooter AS {
  FOOTER footer1 {
    ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES, ButtonStyle: Primary)
    ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
  }
}

-- Use with prefix to create unique names
USE FRAGMENT SaveCancelFooter AS customer_   -- Creates customer_footer1, customer_btnSave, customer_btnCancel
USE FRAGMENT SaveCancelFooter AS order_      -- Creates order_footer1, order_btnSave, order_btnCancel

-- Without prefix (only use once per page)
USE FRAGMENT SaveCancelFooter
```

### Fragment Scope

- Fragments are defined at script scope (available after definition until script ends)
- Fragments are transient (not persisted in MPR, only exist during script execution)
- Fragments can reference other fragments (but no circular references)
- Fragment names must be unique within a script

### SHOW FRAGMENTS

List all defined fragments in the current session:

```sql
DEFINE FRAGMENT SaveCancelFooter AS { ... }
DEFINE FRAGMENT FormHeader AS { ... }

SHOW FRAGMENTS;
-- Output:
-- SaveCancelFooter
-- FormHeader
```

### DESCRIBE FRAGMENT

Show the definition of a script-defined fragment:

```sql
DESCRIBE FRAGMENT SaveCancelFooter;

-- Output:
-- DEFINE FRAGMENT SaveCancelFooter AS {
--   FOOTER footer1 {
--     ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES, ButtonStyle: Primary)
--     ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
--   }
-- }
```

---

## Part 2: DESCRIBE FRAGMENT FROM PAGE

Extract a widget subtree from an existing page as a fragment definition. This enables the **describe → edit → replace** workflow:

### Basic Syntax

```sql
DESCRIBE FRAGMENT FROM PAGE Module.PageName WIDGET widgetName;
```

### Example Workflow

```sql
-- 1. Extract part of an existing page as a fragment
DESCRIBE FRAGMENT FROM PAGE Module.CustomerEdit WIDGET dvCustomer;

-- Output (can be copied, modified, and used in ALTER PAGE):
-- {
--   DATAVIEW dvCustomer (DataSource: $Customer) {
--     TEXTBOX txtName (Label: 'Name', Attribute: Name)
--     TEXTBOX txtEmail (Label: 'Email', Attribute: Email)
--     FOOTER footer1 {
--       ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES)
--       ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
--     }
--   }
-- }

-- 2. Copy the output, modify it, then replace
ALTER PAGE Module.CustomerEdit {
  REPLACE dvCustomer WITH {
    DATAVIEW dvCustomer (DataSource: $Customer) {
      TEXTBOX txtName (Label: 'Name', Attribute: Name)
      TEXTBOX txtEmail (Label: 'Email', Attribute: Email)
      TEXTBOX txtPhone (Label: 'Phone', Attribute: Phone)  -- Added new field
      FOOTER footer1 {
        ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES, ButtonStyle: Primary)  -- Added style
        ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
      }
    }
  }
}
```

### Extract and Save as Reusable Fragment

```sql
-- Extract footer from one page
DESCRIBE FRAGMENT FROM PAGE Module.CustomerEdit WIDGET footer1;

-- Output can be wrapped in DEFINE FRAGMENT for reuse:
DEFINE FRAGMENT StandardFooter AS {
  FOOTER footer1 {
    ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES, ButtonStyle: Primary)
    ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
  }
}

-- Now use in other pages
CREATE PAGE Module.OrderEdit (...) {
  DATAVIEW dvOrder (DataSource: $Order) {
    -- ... fields ...
    USE FRAGMENT StandardFooter
  }
}
```

### Use Cases

1. **Modify complex widgets** - Extract, edit in your editor, replace
2. **Create fragment library** - Extract patterns from Studio Pro-designed pages
3. **Understand page structure** - Inspect specific sections without full DESCRIBE PAGE
4. **Refactoring** - Extract common patterns, convert to fragments

---

## Part 3: ALTER PAGE Statement

Modify existing pages without full replacement:

### Basic Syntax

```sql
ALTER PAGE Module.PageName {
  -- One or more operations
}
```

### SET Property

Update widget properties by name:

```sql
ALTER PAGE Module.CustomerEdit {
  -- Single property
  SET Caption = 'Update' ON btnSave

  -- Multiple properties
  SET (Caption: 'Save Changes', ButtonStyle: Success) ON btnSave

  -- Property on page itself
  SET Title = 'Edit Customer Record'

  -- Nested property for pluggable widgets (dot notation in quotes)
  SET 'showLabel' = false ON cbStatus
  SET 'labelWidth' = 4 ON cbStatus
}
```

### INSERT Widget

Add new widgets to existing containers:

```sql
ALTER PAGE Module.CustomerEdit {
  -- Insert after a named widget
  INSERT AFTER txtName {
    TEXTBOX txtMiddleName (Label: 'Middle Name', Attribute: MiddleName)
  }

  -- Insert before a named widget
  INSERT BEFORE btnSave {
    ACTIONBUTTON btnValidate (Caption: 'Validate', Action: MICROFLOW Module.VAL_Customer)
  }

  -- Insert as first child of a container
  INSERT FIRST IN dvCustomer {
    DYNAMICTEXT formHeader (Content: 'Customer Information', RenderMode: H3)
  }

  -- Insert as last child of a container
  INSERT LAST IN footer1 {
    LINKBUTTON lnkHelp (Caption: 'Help', Action: SHOW_PAGE Module.HelpPage)
  }

  -- Insert fragment
  INSERT AFTER txtEmail {
    USE FRAGMENT AddressFields
  }
}
```

### DROP Widget

Remove widgets from a page:

```sql
ALTER PAGE Module.CustomerEdit {
  -- Remove a single widget
  DROP WIDGET txtMiddleName

  -- Remove multiple widgets
  DROP WIDGET txtFax, txtPager

  -- Remove with IF EXISTS (no error if not found)
  DROP WIDGET IF EXISTS txtLegacyField
}
```

### REPLACE Widget

Swap a widget with new content:

```sql
ALTER PAGE Module.CustomerEdit {
  -- Replace widget entirely
  REPLACE btnCancel WITH {
    LINKBUTTON lnkCancel (Caption: 'Cancel', Action: CLOSE_PAGE)
  }

  -- Replace with fragment
  REPLACE oldFooter WITH {
    USE FRAGMENT SaveCancelFooter
  }
}
```

### MOVE Widget

Relocate a widget within the page:

```sql
ALTER PAGE Module.CustomerEdit {
  -- Move to different position
  MOVE txtPhone AFTER txtEmail

  -- Move to different container
  MOVE btnHelp LAST IN footer1
}
```

---

## Part 4: Incremental Page Building

Build pages step by step, useful for REPL sessions:

### CREATE EMPTY PAGE

```sql
-- Create minimal page structure
CREATE EMPTY PAGE Module.NewPage
(
  Title: 'New Page',
  Layout: Atlas_Core.Atlas_Default
)

-- Now build it up incrementally
ALTER PAGE Module.NewPage {
  INSERT FIRST IN Main {  -- Main is the layout placeholder
    LAYOUTGRID mainGrid { }
  }
}

ALTER PAGE Module.NewPage {
  INSERT FIRST IN mainGrid {
    ROW row1 {
      COLUMN col1 (DesktopWidth: 12) { }
    }
  }
}

ALTER PAGE Module.NewPage {
  INSERT FIRST IN col1 {
    DYNAMICTEXT heading (Content: 'Welcome', RenderMode: H2)
  }
}
```

---

## Part 5: Script Organization

### Multiple Statements in Sequence

```sql
-- fragments.mdl - Define reusable fragments
DEFINE FRAGMENT SaveCancelFooter AS { ... }
DEFINE FRAGMENT FormHeader AS { ... }
DEFINE FRAGMENT AddressFields AS { ... }
/

-- customer-edit.mdl - Create the page
USE SCRIPT 'fragments.mdl'  -- Future: import fragments from other files

CREATE PAGE Module.CustomerEdit (...) {
  USE FRAGMENT FormHeader
  DATAVIEW dvCustomer (DataSource: $Customer) {
    TEXTBOX txtName (Label: 'Name', Attribute: Name)
    USE FRAGMENT AddressFields
    USE FRAGMENT SaveCancelFooter
  }
}
/

-- Later modifications
ALTER PAGE Module.CustomerEdit {
  SET Title = 'Edit Customer v2'
  INSERT AFTER txtName {
    TEXTBOX txtNickname (Label: 'Nickname', Attribute: Nickname)
  }
}
```

---

## Implementation Phases

### Phase 1: Core Fragment System
- `DEFINE FRAGMENT name AS { ... }`
- `USE FRAGMENT name [AS prefix_]` within CREATE PAGE
- `SHOW FRAGMENTS` - list defined fragments
- `DESCRIBE FRAGMENT name` - show fragment definition
- Fragment storage in executor context
- Fragment expansion during page building with prefix support

### Phase 2: DESCRIBE FRAGMENT FROM PAGE
- `DESCRIBE FRAGMENT FROM PAGE Module.Page WIDGET widgetName`
- Extract widget subtree as MDL fragment syntax
- Enables the describe → edit → replace workflow

### Phase 3: ALTER PAGE Basics
- `SET property = value ON widgetName`
- `SET (prop1: val1, prop2: val2) ON widgetName`
- `INSERT AFTER/BEFORE widgetName { ... }`
- `DROP WIDGET widgetName`
- `REPLACE widgetName WITH { ... }`
- Page loading, modification, and saving
- Operations work on raw BSON widget trees, preserving unsupported widget types
- Property validation against widget types

### Phase 4: Advanced ALTER Operations
- `INSERT FIRST/LAST IN containerName`
- `MOVE widgetName AFTER/BEFORE target`
- `DROP WIDGET IF EXISTS`

### Phase 5: Future Enhancements
- Parameterized fragments: `DEFINE FRAGMENT Name($param) AS { ... }`
- `USE SCRIPT 'file.mdl'` for file includes
- Fragment libraries
- Conditional fragments (`USE FRAGMENT X IF condition`)

### Already Implemented (not in scope)
- **Bulk Widget Updates** — `UPDATE WIDGETS SET ... WHERE ...` with module and widget-type filtering, `DRY RUN` support. Fully working in `cmd_widgets.go`.
- **ALTER STYLING ON PAGE** — Partial styling updates on individual widgets. Working in `cmd_styling.go`. Proves the partial page modification pattern.

---

## Grammar Changes

### New Tokens (MDLLexer.g4)

```antlr
DEFINE: D E F I N E;
FRAGMENT: F R A G M E N T;
INSERT: I N S E R T;
BEFORE: B E F O R E;
AFTER: A F T E R;
FIRST: F I R S T;
LAST: L A S T;
```

The following tokens already exist in the lexer: `ALTER`, `MOVE`, `REPLACE`, `WITH`, `EMPTY`.

### New Rules (MDLParser.g4)

```antlr
// Fragment definition
defineFragmentStatement
    : DEFINE FRAGMENT IDENTIFIER AS LBRACE widgetV3* RBRACE
    ;

// Fragment usage (within widget children)
useFragmentStatement
    : USE FRAGMENT IDENTIFIER (AS IDENTIFIER)?  // Optional prefix
    ;

// Show fragments list
showFragmentsStatement
    : SHOW FRAGMENTS
    ;

// Describe fragment (script-defined or from page)
describeFragmentStatement
    : DESCRIBE FRAGMENT IDENTIFIER
    | DESCRIBE FRAGMENT FROM PAGE qualifiedName WIDGET IDENTIFIER
    ;

// ALTER PAGE statement
alterPageStatement
    : ALTER PAGE qualifiedName LBRACE alterOperation+ RBRACE
    ;

alterOperation
    : setPropertyOperation
    | insertOperation
    | dropWidgetOperation
    | replaceWidgetOperation
    | moveWidgetOperation
    ;

setPropertyOperation
    : SET propertyAssignment ON IDENTIFIER
    | SET LPAREN propertyAssignmentList RPAREN ON IDENTIFIER
    | SET propertyAssignment  // Page-level property
    ;

insertOperation
    : INSERT (AFTER | BEFORE) IDENTIFIER LBRACE widgetV3+ RBRACE
    | INSERT (FIRST | LAST) IN IDENTIFIER LBRACE widgetV3+ RBRACE
    ;

dropWidgetOperation
    : DROP WIDGET (IF EXISTS)? identifierList
    ;

replaceWidgetOperation
    : REPLACE IDENTIFIER WITH LBRACE widgetV3+ RBRACE
    ;

moveWidgetOperation
    : MOVE IDENTIFIER (AFTER | BEFORE) IDENTIFIER
    | MOVE IDENTIFIER (FIRST | LAST) IN IDENTIFIER
    ;
```

---

## AST Types

```go
// ast/ast_page_fragments.go

type DefineFragmentStmt struct {
    Name    string
    Widgets []*WidgetV3
}

func (s *DefineFragmentStmt) isStatement() {}

type UseFragmentStmt struct {
    FragmentName string
    Prefix       string // Optional prefix for widget names
}

type ShowFragmentsStmt struct{}

func (s *ShowFragmentsStmt) isStatement() {}

type DescribeFragmentStmt struct {
    FragmentName string        // For script-defined fragments
    PageName     QualifiedName // For DESCRIBE FRAGMENT FROM PAGE
    WidgetName   string        // Widget to extract
    FromPage     bool          // true if FROM PAGE syntax
}

func (s *DescribeFragmentStmt) isStatement() {}

type AlterPageStmt struct {
    PageName   QualifiedName
    Operations []AlterOperation
}

type AlterOperation interface {
    isAlterOperation()
}

type SetPropertyOp struct {
    WidgetName  string // empty for page-level
    Properties  map[string]interface{}
}

type InsertOp struct {
    Position    InsertPosition // AFTER, BEFORE, FIRST, LAST
    TargetName  string
    Widgets     []*WidgetV3
}

type DropWidgetOp struct {
    WidgetNames []string
    IfExists    bool
}

type ReplaceWidgetOp struct {
    WidgetName string
    NewWidgets []*WidgetV3
}

type MoveWidgetOp struct {
    WidgetName  string
    Position    InsertPosition
    TargetName  string
}

type InsertPosition string
const (
    InsertAfter  InsertPosition = "AFTER"
    InsertBefore InsertPosition = "BEFORE"
    InsertFirst  InsertPosition = "FIRST"
    InsertLast   InsertPosition = "LAST"
)
```

---

## Executor Changes

### Fragment Registry

```go
type Executor struct {
    // ... existing fields
    fragments map[string]*ast.DefineFragmentStmt
}

func (e *Executor) execDefineFragment(s *ast.DefineFragmentStmt) error {
    if _, exists := e.fragments[s.Name]; exists {
        return fmt.Errorf("fragment %s already defined", s.Name)
    }
    e.fragments[s.Name] = s
    return nil
}
```

### Fragment Expansion

During page building, when encountering `USE FRAGMENT`:

```go
func (pb *pageBuilder) expandFragment(name string) ([]*ast.WidgetV3, error) {
    fragment, ok := pb.executor.fragments[name]
    if !ok {
        return nil, fmt.Errorf("fragment not found: %s", name)
    }
    // Return copy of widgets to avoid mutation
    return cloneWidgets(fragment.Widgets), nil
}
```

### ALTER PAGE Execution

```go
func (e *Executor) execAlterPage(s *ast.AlterPageStmt) error {
    // 1. Load existing page
    page, err := e.reader.GetPage(pageID)
    if err != nil {
        return err
    }

    // 2. Build widget index by name
    widgetIndex := buildWidgetIndex(page)

    // 3. Apply operations in order
    for _, op := range s.Operations {
        if err := e.applyAlterOperation(page, widgetIndex, op); err != nil {
            return err
        }
    }

    // 4. Save modified page
    return e.writer.UpdatePage(page)
}
```

---

## Example: Complete Workflow

### Workflow A: Creating New Pages with Fragments

```sql
-- Step 1: Define reusable fragments
DEFINE FRAGMENT CrudButtons AS {
  FOOTER formFooter {
    ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES, ButtonStyle: Primary)
    ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
    ACTIONBUTTON btnDelete (Caption: 'Delete', Action: DELETE, ButtonStyle: Danger)
  }
}
/

-- Step 2: Create page using fragments
CREATE PAGE CRM.Customer_Edit
(
  Params: { $Customer: CRM.Customer },
  Title: 'Edit Customer',
  Layout: Atlas_Core.PopupLayout
)
{
  DATAVIEW dvCustomer (DataSource: $Customer) {
    TEXTBOX txtName (Label: 'Name', Attribute: Name)
    TEXTBOX txtEmail (Label: 'Email', Attribute: Email)
    TEXTBOX txtPhone (Label: 'Phone', Attribute: Phone)
    USE FRAGMENT CrudButtons
  }
}
/

-- Step 3: Simple modifications with ALTER
ALTER PAGE CRM.Customer_Edit {
  INSERT AFTER txtEmail {
    TEXTBOX txtWebsite (Label: 'Website', Attribute: Website)
  }
  SET Caption = 'Save Customer' ON btnSave
  DROP WIDGET btnDelete
}
```

### Workflow B: Describe → Edit → Replace (Modifying Existing Pages)

This is the key workflow for modifying complex existing pages:

```sql
-- Step 1: Extract the section you want to modify
DESCRIBE FRAGMENT FROM PAGE CRM.Customer_Edit WIDGET dvCustomer;

-- Output:
-- {
--   DATAVIEW dvCustomer (DataSource: $Customer) {
--     TEXTBOX txtName (Label: 'Name', Attribute: Name)
--     TEXTBOX txtEmail (Label: 'Email', Attribute: Email)
--     TEXTBOX txtWebsite (Label: 'Website', Attribute: Website)
--     TEXTBOX txtPhone (Label: 'Phone', Attribute: Phone)
--     FOOTER formFooter {
--       ACTIONBUTTON btnSave (Caption: 'Save Customer', Action: SAVE_CHANGES, ButtonStyle: Primary)
--       ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
--     }
--   }
-- }

-- Step 2: Copy output, edit in your editor, then replace
ALTER PAGE CRM.Customer_Edit {
  REPLACE dvCustomer WITH {
    DATAVIEW dvCustomer (DataSource: $Customer) {
      -- Reorganized into two columns
      LAYOUTGRID formGrid {
        ROW row1 {
          COLUMN colLeft (DesktopWidth: 6) {
            TEXTBOX txtName (Label: 'Name', Attribute: Name)
            TEXTBOX txtEmail (Label: 'Email', Attribute: Email)
          }
          COLUMN colRight (DesktopWidth: 6) {
            TEXTBOX txtPhone (Label: 'Phone', Attribute: Phone)
            TEXTBOX txtWebsite (Label: 'Website', Attribute: Website)
          }
        }
      }
      FOOTER formFooter {
        ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES, ButtonStyle: Primary)
        ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES)
      }
    }
  }
}
```

### Workflow C: Extract Patterns from Studio Pro Pages

```sql
-- Extract a well-designed footer from a Studio Pro page
DESCRIBE FRAGMENT FROM PAGE Atlas_Core.ExamplePage WIDGET footerActions;

-- Wrap in DEFINE FRAGMENT for reuse
DEFINE FRAGMENT StandardActions AS {
  -- paste extracted content here
}

-- Use in your pages
CREATE PAGE Module.NewPage (...) {
  DATAVIEW dv (...) {
    -- ... fields ...
    USE FRAGMENT StandardActions
  }
}
```

### Workflow D: Safe Editing of Pages with Unsupported Widgets

This is the primary motivation for this proposal. A page may contain widgets that MDL cannot yet describe or round-trip (e.g., newer pluggable widgets, specialized marketplace widgets). Today, `DESCRIBE PAGE` renders these as comments and `CREATE OR REPLACE PAGE` silently drops them. With ALTER PAGE, we can safely modify the known parts:

```sql
-- Page contains a mix of supported and unsupported widgets.
-- DESCRIBE PAGE shows:
--   DATAVIEW dvCustomer (...) {
--     TEXTBOX txtName (Label: 'Name', Attribute: Name)
--     -- CustomWidgets$SomeMarketplaceWidget (mpWidget1)    <-- unsupported, shown as comment
--     FOOTER footer1 { ... }
--   }

-- ALTER PAGE works on the raw BSON widget tree, so unsupported widgets are preserved:
ALTER PAGE Module.CustomerEdit {
  INSERT AFTER txtName {
    TEXTBOX txtEmail (Label: 'Email', Attribute: Email)
  }
  SET Caption = 'Update' ON btnSave
}
-- mpWidget1 is untouched — it stays in the page exactly as it was.
```

---

## Design Decisions

1. **Raw BSON widget tree for ALTER PAGE**: ALTER PAGE operations must work on the raw BSON widget tree (not parsed/reconstructed widgets). This is what makes unsupported widget preservation possible — widgets that MDL cannot parse are kept as opaque BSON documents and passed through unchanged. Only the targeted widgets are modified. This follows the same approach proven by the existing `UPDATE WIDGETS` and `ALTER STYLING` implementations.

2. **Fragment naming conflicts**: Use prefix syntax to avoid conflicts
   - `USE FRAGMENT Name AS prefix_` creates prefixed widget names
   - Without prefix, error if names conflict

3. **Validation**: ALTER operations validate property assignments against widget types
   - `SET Caption` only valid on widgets that have Caption property
   - Prevents creating invalid models

4. **Dry run**: Not needed for ALTER PAGE (keep it simple for now)

5. **Fragment scope**: Script-scoped (transient)
   - Fragments exist only during script execution
   - Use `DESCRIBE FRAGMENT FROM PAGE` to extract and recreate

---

## Open Questions

1. **Undo support**: Should ALTER operations be reversible?
   - Could generate inverse operations for rollback

2. **Nested fragments**: Should fragments be allowed to USE other fragments?
   - Adds complexity but increases reusability

---

## Success Criteria

1. Pages containing unsupported widget types can be safely modified without data loss
2. A 200-line page script can be broken into 5-6 smaller fragments
3. Single property changes don't require page rewrite
4. Common patterns (CRUD buttons, form layouts) defined once, used everywhere
5. REPL users can build pages incrementally
