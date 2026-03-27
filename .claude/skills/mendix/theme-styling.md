# Theme & Styling System

## When to Use This Skill

Use this skill when the user wants to:
- View or modify widget styling (CSS classes, inline styles, design properties)
- Browse available Atlas UI design properties
- Understand the Mendix theme/styling architecture
- Work with `themesource/` directories or `design-properties.json`
- Apply consistent styling patterns across pages
- Debug CSS or design property issues
- Hot-reload CSS changes during development

## Mendix Theme Architecture

### Directory Structure

```
MyProject/
├── theme/                          # Project-level overrides
│   ├── web/
│   │   ├── main.scss               # SCSS entry point (import chain)
│   │   ├── custom-variables.scss   # Project variable overrides
│   │   ├── exclusion-variables.scss # Exclude unwanted Atlas components
│   │   ├── login.html              # Custom login page
│   │   └── settings.json           # Theme settings
│   └── native/
│       ├── main.js                 # Native entry point
│       └── custom-variables.js     # Native variable overrides
│
├── themesource/                    # Module-level theme definitions
│   ├── atlas_core/                 # Base framework (always present)
│   │   ├── web/
│   │   │   ├── design-properties.json  # Widget design properties
│   │   │   ├── variables.scss          # Color/spacing/font variables
│   │   │   └── ...                     # Component SCSS files
│   │   ├── native/
│   │   │   └── design-properties.json
│   │   └── public/resources/           # Fonts, icons, images
│   ├── datawidgets/                # DataGrid2, Gallery, etc.
│   ├── atlas_web_content/          # Web content styles
│   └── <module_name>/              # Each module can contribute styles
│       └── web/design-properties.json
│
└── theme-cache/web/                # Compiled CSS output (build artifact)
```

### SCSS Compilation Chain

`atlas_core/web/main.scss` imports in order:
1. Default variables (`atlas_core`)
2. Exclusion variables (disable Atlas components)
3. Project custom variables (`theme/web/custom-variables.scss`)
4. Bootstrap framework
5. MXUI components
6. Core styles (base, animations, spacing, flex)
7. Widget-specific styles
8. Module-specific styles from `themesource/*/web/*.scss`

### Design Properties System

Design properties are defined in `themesource/*/web/design-properties.json` files. They provide structured styling options in Studio Pro's properties panel.

**Property types:**

| Type | MDL Syntax | Example |
|------|------------|---------|
| Toggle | `'Name': ON` / `'Name': OFF` | `'Full width': ON` |
| Dropdown | `'Name': 'Option'` | `'Background color': 'Brand Primary'` |
| ColorPicker | `'Name': 'Option'` | `'Background color': 'Brand Success'` |
| ToggleButtonGroup | `'Name': 'Option'` | `'Size': 'Large'` |
| Spacing | `'Spacing top': 'Large'` | Complex margin/padding groups |

**Inheritance:** All widgets inherit properties from the `"Widget"` key (Spacing, Align self, Hide on). Type-specific keys add extra properties.

**atlas_core common properties:**

| Widget Type | Properties |
|-------------|------------|
| Widget (all) | Spacing, Align self, Hide on |
| DivContainer | Align content, Background color, Shade, Shadow |
| Button | Size, Align icon, Full width, Border |
| DataGrid | Hover rows, Striped rows, Bordered |
| ListView | Styled, List style unstyled |
| StaticImageViewer | Shape |
| DynamicImageViewer | Shape |

## MDL Commands

### SHOW DESIGN PROPERTIES

Browse available design properties from the project's `themesource/`.

```sql
-- List ALL design properties for all widget types
SHOW DESIGN PROPERTIES;

-- Properties for a specific widget type
SHOW DESIGN PROPERTIES FOR CONTAINER;
SHOW DESIGN PROPERTIES FOR ACTIONBUTTON;
SHOW DESIGN PROPERTIES FOR DATAGRID;
SHOW DESIGN PROPERTIES FOR LISTVIEW;
```

Output shows property name, type, and available options:
```
=== Widget (inherited) ===
  Spacing                  Spacing     -- margin/padding groups
  Align self               ToggleButtonGroup [Left, Right]
  Hide on                  ToggleButtonGroup [Phone, Tablet, Desktop]

=== DivContainer ===
  Align content            Dropdown    [Left align as a row, Center align as a row, ...]
  Background color         ColorPicker [Background Primary, Brand Primary, ...]
  Shadow                   ToggleButtonGroup [Small, Medium, Large]
```

### DESCRIBE STYLING

Inspect current styling on widgets in a page or snippet.

```sql
-- All styled widgets in a page
DESCRIBE STYLING ON PAGE MyModule.MyPage;

-- Specific widget
DESCRIBE STYLING ON PAGE MyModule.MyPage WIDGET btnSave;

-- Snippet
DESCRIBE STYLING ON SNIPPET MyModule.MySnippet;
DESCRIBE STYLING ON SNIPPET MyModule.MySnippet WIDGET ctnMain;
```

Output example:
```
WIDGET btnSave (ActionButton)
  Class: 'btn-lg btn-bordered'
  DesignProperties: ['Size': 'Large', 'Full width': ON]

WIDGET ctnHero (Container)
  Class: 'card'
  Style: 'border-left: 4px solid #264AE5;'
  DesignProperties: ['Background color': 'Brand Primary', 'Shadow': 'Medium']
```

### ALTER STYLING

Modify widget styling. Requires write mode (`-w` flag or `OPEN ... FOR WRITING`).

```sql
-- Set CSS class
ALTER STYLING ON PAGE MyModule.MyPage WIDGET ctnMain
  SET Class = 'card mx-spacing-top-large';

-- Set inline style
ALTER STYLING ON PAGE MyModule.MyPage WIDGET ctnMain
  SET Style = 'background-color: #f8f9fa; padding: 16px;';

-- Set design property (option)
ALTER STYLING ON PAGE MyModule.MyPage WIDGET ctnMain
  SET 'Background color' = 'Brand Primary';

-- Toggle design property ON
ALTER STYLING ON PAGE MyModule.MyPage WIDGET btnSave
  SET 'Full width' = ON;

-- Toggle design property OFF (removes it)
ALTER STYLING ON PAGE MyModule.MyPage WIDGET btnSave
  SET 'Full width' = OFF;

-- Multiple assignments in one command
ALTER STYLING ON PAGE MyModule.MyPage WIDGET ctnMain
  SET Class = 'card',
      'Background color' = 'Brand Primary',
      'Shadow' = 'Medium',
      'Full width' = ON;

-- Clear all design properties
ALTER STYLING ON PAGE MyModule.MyPage WIDGET ctnMain
  CLEAR DESIGN PROPERTIES;

-- Works on snippets too
ALTER STYLING ON SNIPPET MyModule.MySnippet WIDGET btnCancel
  SET Class = 'btn-bordered', 'Size' = 'Large';
```

### Inline Styling During Page Creation

Apply styling when creating widgets in `CREATE PAGE` or `CREATE SNIPPET`:

```sql
CREATE PAGE MyModule.StyledPage 'Styled Example' (Layout: 'Atlas_Default')
{
  -- CSS class
  CONTAINER ctnCard (Class: 'card mx-2') {
    DYNAMICTEXT txtTitle (Content: 'Hello', RenderMode: H2)
  }

  -- Inline style
  CONTAINER ctnBanner (Style: 'background-color: #264AE5; color: white; padding: 24px;') {
    DYNAMICTEXT txtBanner (Content: 'Welcome', RenderMode: H1)
  }

  -- Design properties
  CONTAINER ctnSpaced (DesignProperties: ['Spacing top': 'Large', 'Background color': 'Brand Primary']) {
    DYNAMICTEXT txtContent (Content: 'Spaced and colored')
  }

  -- All three combined
  CONTAINER ctnFull (
    Class: 'card',
    Style: 'border-left: 4px solid #264AE5;',
    DesignProperties: ['Spacing top': 'Large', 'Shadow': 'Medium']
  ) {
    ACTIONBUTTON btnAction (Caption: 'Go', DesignProperties: ['Size': 'Large', 'Full width': ON])
  }
}
```

### ALTER PAGE with Styling

Set styling via `ALTER PAGE` operations:

```sql
ALTER PAGE MyModule.MyPage {
  -- Set Class on an existing widget
  SET Class = 'card mx-spacing-top-large' ON ctnMain

  -- Set Style
  SET Style = 'padding: 16px;' ON ctnHeader

  -- Insert a styled widget
  INSERT AFTER txtName {
    CONTAINER ctnHighlight (
      Class: 'alert alert-info',
      DesignProperties: ['Spacing top': 'Small']
    ) {
      DYNAMICTEXT txtHint (Content: 'Enter your full name')
    }
  }
}
```

## Widget Type Mapping

MDL keywords map to design-properties.json keys:

| MDL Keyword | design-properties.json Key | Notes |
|-------------|---------------------------|-------|
| `CONTAINER` | `DivContainer` | |
| `ACTIONBUTTON` | `Button` | |
| `LINKBUTTON` | `Button` | Same as ACTIONBUTTON |
| `TEXTBOX` | `TextBox` | |
| `TEXTAREA` | `TextArea` | |
| `DATEPICKER` | `DatePicker` | |
| `CHECKBOX` | `CheckBox` | |
| `RADIOBUTTONS` | `RadioButtons` | |
| `COMBOBOX` | `ReferenceSelector` | |
| `DROPDOWN` | `DropDown` | |
| `DATAGRID` | `DataGrid` | Classic DataGrid |
| `DATAVIEW` | `DataView` | |
| `LISTVIEW` | `ListView` | |
| `GALLERY` | `Gallery` | Pluggable widget |
| `LAYOUTGRID` | `LayoutGrid` | |
| `DYNAMICTEXT` | `DynamicText` | |
| `STATICTEXT` | `Label` | |
| `IMAGE` / `STATICIMAGE` | `StaticImageViewer` | |
| `DYNAMICIMAGE` | `DynamicImageViewer` | |

Pluggable widgets use their full widget ID as key (e.g., `com.mendix.widget.web.datagrid.Datagrid`).

## CSS Hot-Reload Workflow

For theme/styling changes during development with Docker:

```bash
# 1. Compile SCSS into deployment package (~55s)
mxcli docker build -p app.mpr

# 2. Push compiled CSS to browsers (instant, no page reload)
mxcli docker reload -p app.mpr --css
```

The `--css` flag calls the M2EE `update_styling` action, which pushes CSS via WebSocket to all connected browsers. **It does NOT compile SCSS** — always run `docker build` first.

For non-CSS changes (Class, Style, DesignProperties on widgets), use normal reload:
```bash
mxcli docker reload -p app.mpr
```

## Common Styling Patterns

### Card Layout
```sql
CONTAINER ctnCard (Class: 'card', DesignProperties: ['Shadow': 'Small', 'Spacing top': 'Large']) {
  CONTAINER ctnCardBody (Class: 'card-body') {
    DYNAMICTEXT txtTitle (Content: 'Card Title', RenderMode: H3)
    DYNAMICTEXT txtBody (Content: 'Card content here')
  }
}
```

### Alert / Notification
```sql
CONTAINER ctnAlert (Class: 'alert alert-info', DesignProperties: ['Spacing top': 'Small']) {
  DYNAMICTEXT txtMessage (Content: 'This is an informational message')
}
```

### Full-Width Action Button
```sql
ACTIONBUTTON btnSubmit (
  Caption: 'Submit',
  OnClick: MICROFLOW 'MyModule.ACT_Submit',
  ButtonStyle: Primary,
  DesignProperties: ['Full width': ON, 'Size': 'Large']
)
```

### Responsive Hide
```sql
-- Hide sidebar on phone
CONTAINER ctnSidebar (DesignProperties: ['Hide on': 'Phone']) {
  -- sidebar content
}
```

### Centered Content Column
```sql
CONTAINER ctnCenter (DesignProperties: ['Align content': 'Center align as a column']) {
  -- centered content
}
```

## Caveats

### DYNAMICTEXT + Style Crash

**Never** apply `Style` directly to a DYNAMICTEXT widget — it crashes MxBuild with a NullReferenceException. Wrap in a CONTAINER:

```sql
-- WRONG: crashes MxBuild
DYNAMICTEXT txt (Content: 'Hello', Style: 'color: red;')

-- CORRECT: style the container
CONTAINER ctn (Style: 'color: red;') {
  DYNAMICTEXT txt (Content: 'Hello')
}
```

This also applies to `ALTER STYLING` and `ALTER PAGE SET Style`:
```sql
-- WRONG
ALTER STYLING ON PAGE Mod.Page WIDGET txtHeading SET Style = 'color: red;';

-- CORRECT: wrap first, then style the container
ALTER PAGE Mod.Page {
  REPLACE txtHeading WITH {
    CONTAINER ctnHeading (Style: 'color: red;') {
      DYNAMICTEXT txtHeading (Content: 'Heading', RenderMode: H2)
    }
  }
}
```

### Design Property Key Must Match Exactly

Design property keys are **case-sensitive** and must match the `name` field in `design-properties.json` exactly:
```sql
-- CORRECT
DesignProperties: ['Spacing top': 'Large']

-- WRONG (case mismatch)
DesignProperties: ['spacing top': 'Large']
DesignProperties: ['SPACING TOP': 'Large']
```

### Module Cleanup

When deleting a module, its `themesource/<modulename>/` directory is automatically removed by `DeleteModuleWithCleanup()`.

## Checklist

- [ ] Run `SHOW DESIGN PROPERTIES FOR <type>` to check available options before applying
- [ ] Use `DESCRIBE STYLING ON PAGE ...` to inspect current styling before modifying
- [ ] Never apply `Style` directly to DYNAMICTEXT — wrap in a CONTAINER
- [ ] Design property keys are case-sensitive — match `design-properties.json` exactly
- [ ] Toggle properties use `ON`/`OFF`, option properties use quoted strings
- [ ] For CSS changes, run `docker build` then `docker reload --css` (build compiles SCSS, reload pushes to browsers)
- [ ] Verify changes with `DESCRIBE STYLING` after modification
