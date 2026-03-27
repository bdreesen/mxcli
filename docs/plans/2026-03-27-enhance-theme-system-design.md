# Enhanced Theme System Design

**Date:** 2026-03-27
**Status:** Draft
**Branch:** feat/enhance-theme-system

## Goal

Enhance mxcli's theme system to give AI stronger theme customization capability and efficiency. Four feature areas, each building on existing infrastructure.

## Current State

### Existing Commands
- `SHOW DESIGN PROPERTIES [FOR type]` — browse design-properties.json entries
- `DESCRIBE STYLING ON PAGE/SNIPPET ... [WIDGET name]` — inspect widget styling
- `ALTER STYLING ON PAGE/SNIPPET ... WIDGET name SET ...` — modify single widget styling
- Inline styling during page creation (`Class:`, `Style:`, `DesignProperties:`)

### Existing Infrastructure
- `ThemeRegistry` in `executor/theme_reader.go` — loads design-properties.json
- `applyStylingAssignments()` in `executor/cmd_styling.go` — applies Class/Style/DesignProperties
- `SettingsThemeModuleEntry` / `ThemeModuleOrder` in generated metamodel types
- `IsThemeModule` flag on `ProjectsModule`
- Widget type mappings: MDL keyword <-> design-properties.json key <-> BSON $Type

---

## Feature 1: Batch Styling Operations

### Syntax

#### Page/Snippet Scope

```sql
-- By widget type (all ACTIONBUTTONs in the page)
ALTER STYLING ON PAGE Mod.Page WIDGETS OF TYPE ACTIONBUTTON
  SET 'Full width' = ON, Class = 'btn-bordered';

-- By name pattern (SQL LIKE wildcards: % = any chars, _ = single char)
ALTER STYLING ON PAGE Mod.Page WIDGETS LIKE 'btn%'
  SET 'Size' = 'Large';

-- All widgets in a page
ALTER STYLING ON PAGE Mod.Page ALL WIDGETS
  SET 'Spacing top' = 'Large';

-- Same syntax for snippets
ALTER STYLING ON SNIPPET Mod.Snip WIDGETS OF TYPE CONTAINER
  SET 'Background color' = 'Brand Primary';
```

#### Cross-Page / Global Scope

```sql
-- All pages in the project
ALTER STYLING ON ALL PAGES WIDGETS OF TYPE ACTIONBUTTON
  SET 'Full width' = ON;

-- All pages in a specific module
ALTER STYLING IN MODULE MyModule WIDGETS OF TYPE CONTAINER
  SET 'Background color' = 'Brand Primary';

-- All snippets
ALTER STYLING ON ALL SNIPPETS WIDGETS OF TYPE ACTIONBUTTON
  SET Class = 'btn-lg';
```

### Execution Semantics

- Returns affected count: `Updated styling on 23 widgets across 8 pages`
- `LIKE` uses SQL-style wildcards (`%` = any chars, `_` = single char)
- `OF TYPE` matches MDL widget type keywords (CONTAINER, ACTIONBUTTON, TEXTBOX, etc.)
- Reuses existing `applyStylingAssignments()` for each matched widget
- All modified pages/snippets are saved via `writer.UpdatePage()` / `writer.UpdateSnippet()`

### AST Nodes

```go
// AlterStylingBatchStmt represents batch styling operations.
type AlterStylingBatchStmt struct {
    // Target scope
    Scope         string        // "PAGE", "SNIPPET", "ALL_PAGES", "ALL_SNIPPETS", "MODULE"
    ContainerName QualifiedName // For PAGE/SNIPPET scope
    ModuleName    string        // For MODULE scope

    // Widget selection
    Selector     string // "TYPE", "LIKE", "ALL"
    WidgetType   string // For TYPE selector (e.g., "ACTIONBUTTON")
    NamePattern  string // For LIKE selector (e.g., "btn%")

    // Assignments (reuse existing)
    Assignments      []StylingAssignment
    ClearDesignProps bool
}
```

### Grammar Changes

Extend `alterStatement` in `MDLParser.g4`:

```antlr
alterStatement
    : ...existing rules...
    | ALTER STYLING ON (PAGE | SNIPPET) qualifiedName WIDGET IDENTIFIER alterStylingAction+  // existing
    | ALTER STYLING ON (PAGE | SNIPPET) qualifiedName batchWidgetSelector alterStylingAction+ // new: page batch
    | ALTER STYLING ON ALL (PAGES | SNIPPETS) batchWidgetSelector alterStylingAction+         // new: global batch
    | ALTER STYLING IN MODULE (qualifiedName | IDENTIFIER) batchWidgetSelector alterStylingAction+ // new: module batch
    ;

batchWidgetSelector
    : WIDGETS OF TYPE widgetTypeKeyword
    | WIDGETS LIKE STRING_LITERAL
    | ALL WIDGETS
    ;
```

New keywords needed: `WIDGETS`, `PAGES`, `SNIPPETS` (check if already defined as tokens).

### Implementation Files

| File | Change |
|------|--------|
| `mdl/grammar/MDLParser.g4` | Add batch alter rules + `batchWidgetSelector` |
| `mdl/grammar/MDLLexer.g4` | Add `WIDGETS`, `PAGES`, `SNIPPETS` tokens if missing |
| `mdl/ast/ast_styling.go` | Add `AlterStylingBatchStmt` |
| `mdl/visitor/visitor_styling.go` | Add `exitAlterStylingBatchStatement()` |
| `mdl/executor/cmd_styling.go` | Add `execAlterStylingBatch()` with page/snippet iteration |

---

## Feature 2: Structured SCSS Variable Management

### Syntax

#### View Variables

```sql
-- List all custom variables (from theme/web/custom-variables.scss)
SHOW THEME VARIABLES;
-- Output:
-- Variable                    Value                      Source
-- $brand-primary              #264AE5                    custom-variables.scss
-- $brand-secondary            #1FC5A8                    custom-variables.scss
-- $font-family-base           "Open Sans", sans-serif    custom-variables.scss

-- List atlas_core default variables (read-only reference)
SHOW THEME VARIABLES DEFAULT;

-- Search variables by pattern
SHOW THEME VARIABLES LIKE '%brand%';
```

#### Modify Variables

```sql
-- Set a single variable (writes to theme/web/custom-variables.scss)
ALTER THEME VARIABLE '$brand-primary' = '#FF6B35';

-- Batch set multiple variables
ALTER THEME VARIABLES
  '$brand-primary' = '#FF6B35',
  '$brand-secondary' = '#2D3748',
  '$font-size-default' = '15px';

-- Reset to atlas_core default (remove line from custom-variables.scss)
ALTER THEME VARIABLE '$brand-primary' RESET;
```

### Implementation Strategy

#### SCSS Variable Parser

Simple line-based parser — no full SCSS AST needed:

```go
// ScssVariable represents a parsed $variable: value; line.
type ScssVariable struct {
    Name      string // e.g., "$brand-primary"
    Value     string // e.g., "#264AE5"
    Comment   string // trailing comment if any
    LineNum   int    // line number in file
    IsDefault bool   // has !default flag
}

// parseScssVariables reads a .scss file and extracts variable declarations.
// Matches: $var-name: value;  or  $var-name: value !default;
// Preserves non-variable lines (comments, blank lines) for write-back.
func parseScssVariables(content string) ([]ScssVariable, error)
```

Regex pattern: `^\s*(\$[\w-]+)\s*:\s*(.+?)\s*(!default)?\s*;\s*(//.*)?$`

#### Write-Back Strategy

- Parse file into lines
- Find existing variable line and update value in-place
- If variable doesn't exist, append at end before closing comments
- Preserve all non-variable lines (comments, blank lines, imports)
- Never rewrite the entire file — surgical line edits only

#### Validation

- `ALTER THEME VARIABLE` only accepts variable names that exist in atlas_core defaults
- Prevents typos like `$bran-primary` (unknown variable warning)
- Override with `FORCE` keyword if user truly wants a new custom variable

### AST Nodes

```go
type ShowThemeVariablesStmt struct {
    ShowDefault bool   // SHOW THEME VARIABLES DEFAULT
    Pattern     string // LIKE '%brand%' filter
}

type AlterThemeVariableStmt struct {
    Variables []ThemeVariableAssignment // One or more assignments
    Reset     bool                      // RESET (remove custom override)
}

type ThemeVariableAssignment struct {
    Name  string // "$brand-primary"
    Value string // "#FF6B35"
}
```

### Implementation Files

| File | Change |
|------|--------|
| `mdl/grammar/MDLParser.g4` | Add SHOW/ALTER THEME VARIABLE(S) rules |
| `mdl/grammar/MDLLexer.g4` | Add `THEME` token (or reuse existing) |
| `mdl/ast/ast_styling.go` | Add `ShowThemeVariablesStmt`, `AlterThemeVariableStmt` |
| `mdl/visitor/visitor_styling.go` | Add visitor methods |
| `mdl/executor/theme_variables.go` | **New file**: SCSS parser, variable read/write |
| `mdl/executor/cmd_styling.go` | Add `execShowThemeVariables()`, `execAlterThemeVariable()` |

---

## Feature 3: Theme Module Management

### Syntax

```sql
-- List theme modules with priority and metadata
SHOW THEME MODULES;
-- Output:
-- #  Module               IsTheme  Priority  HasDesignProps  Version
-- 1  atlas_core           true     1         true            3.14.3
-- 2  atlas_web_content    true     2         false           3.6.0
-- 3  datawidgets          true     3         true            2.24.0
-- 4  MyModule             false    -         false           -
-- 5  MyCustomTheme        true     4         true            1.0.0

-- Reorder theme module priority (controls SCSS compilation order)
ALTER THEME MODULE ORDER ('atlas_core', 'MyCustomTheme', 'atlas_web_content', 'datawidgets');

-- Mark a module as a theme module
ALTER MODULE MyCustomTheme SET IsThemeModule = true;

-- Unmark
ALTER MODULE MyCustomTheme SET IsThemeModule = false;

-- Detailed view of a theme module's files
DESCRIBE THEME MODULE atlas_core;
-- Output:
-- Module: atlas_core (v3.14.3)
-- Path: themesource/atlas_core/
-- Files:
--   web/design-properties.json (27 widget types, 42 properties)
--   web/variables.scss (186 variables)
--   web/main.scss
--   native/design-properties.json
--   public/resources/fonts/ (6 files)
```

### Implementation Strategy

#### SHOW THEME MODULES

1. Read `ProjectSettings` → `WebUIProjectSettingsPart.ThemeModuleOrder`
2. List all modules, check `IsThemeModule` flag on each
3. Scan `themesource/<module>/` to check for `design-properties.json`
4. Read `.version` file if present
5. Merge and display as table

#### ALTER THEME MODULE ORDER

1. Parse ordered list of module names
2. Validate all names are actual modules with `IsThemeModule = true`
3. Build `[]SettingsThemeModuleEntry` with correct IDs
4. Update `ProjectSettings` via writer

#### ALTER MODULE ... SET IsThemeModule

Extend existing `ALTER MODULE` executor to handle the `IsThemeModule` property. This modifies the module's BSON in the project settings.

#### DESCRIBE THEME MODULE

1. Locate `themesource/<name>/` directory
2. Recursively list files with sizes
3. If `design-properties.json` exists, parse and count widget types + property count
4. If `variables.scss` exists, count variable declarations

### AST Nodes

```go
type ShowThemeModulesStmt struct{}

type AlterThemeModuleOrderStmt struct {
    ModuleNames []string // Ordered list of module names
}

type DescribeThemeModuleStmt struct {
    ModuleName string
}
```

### Implementation Files

| File | Change |
|------|--------|
| `mdl/grammar/MDLParser.g4` | Add SHOW/ALTER/DESCRIBE THEME MODULE rules |
| `mdl/ast/ast_styling.go` | Add AST nodes |
| `mdl/visitor/visitor_styling.go` | Add visitor methods |
| `mdl/executor/cmd_styling.go` | Add executor methods |
| `mdl/executor/cmd_modules.go` | Extend ALTER MODULE for IsThemeModule |

---

## Feature 4: Theme Presets via MDL Script Templates

### Syntax

```sql
-- List available presets (built-in + project-local)
SHOW THEME PRESETS;
-- Output:
-- Name           Description                          Source
-- dark           Dark color scheme with light text     built-in
-- ocean          Blue-toned professional palette       built-in
-- warm           Warm earthy tones                     built-in
-- high-contrast  Accessibility-focused high contrast   built-in
-- corporate      Company brand colors                  project

-- Preview what a preset will do (dry run)
DESCRIBE THEME PRESET 'dark';
-- Output: the MDL statements that will be executed

-- Apply a preset (executes the MDL script)
APPLY THEME PRESET 'dark';
-- Output: Applied theme preset 'dark' (5 variables updated)
```

### Preset File Format

Presets are standard `.mdl` files with metadata comment headers:

```sql
-- preset: dark
-- description: Dark color scheme with light text on dark backgrounds
-- author: mxcli
-- version: 1.0

ALTER THEME VARIABLES
  '$brand-default' = '#1a1a2e',
  '$brand-primary' = '#6c63ff',
  '$brand-success' = '#00b894',
  '$brand-warning' = '#fdcb6e',
  '$brand-danger' = '#e17055',
  '$bg-color' = '#16213e',
  '$bg-color-secondary' = '#0f3460',
  '$font-color-default' = '#edf2f7',
  '$font-color-secondary' = '#a0aec0';
```

Presets can use any MDL command, including batch styling:

```sql
-- preset: professional
-- description: Clean professional look with subtle styling

ALTER THEME VARIABLES
  '$brand-primary' = '#2563EB',
  '$brand-secondary' = '#475569';

-- Apply consistent button styling across all pages
ALTER STYLING ON ALL PAGES WIDGETS OF TYPE ACTIONBUTTON
  SET 'Border' = ON, 'Size' = 'Small';

-- Add shadow to all containers
ALTER STYLING ON ALL PAGES WIDGETS OF TYPE CONTAINER
  SET 'Shadow' = 'Small';
```

### Storage and Discovery

| Location | Purpose | Priority |
|----------|---------|----------|
| `cmd/mxcli/presets/*.mdl` | Built-in presets (via `go:embed`) | Lowest (fallback) |
| `theme/presets/*.mdl` | Project-specific presets | Highest (overrides) |

When names conflict, project presets take priority over built-in.

### Implementation Strategy

1. **Embed built-in presets** in `cmd/mxcli/presets/` using `go:embed`
2. **Parse metadata** from comment headers (`-- preset:`, `-- description:`)
3. **SHOW** scans both locations, merges, displays table
4. **DESCRIBE** reads and outputs the script content
5. **APPLY** feeds the script to the existing MDL executor — zero new execution logic needed
6. **Minimal code**: the preset system is essentially a thin wrapper around script execution

### Built-in Preset Ideas

| Preset | Variables Changed | Description |
|--------|-------------------|-------------|
| `dark` | 8-10 color vars | Dark backgrounds, light text |
| `ocean` | 6-8 color vars | Blue-toned professional palette |
| `warm` | 6-8 color vars | Earthy warm tones |
| `high-contrast` | 10+ vars | WCAG AA compliant contrast ratios |
| `minimal` | 4-5 vars | Reduced visual noise, neutral colors |

### AST Nodes

```go
type ShowThemePresetsStmt struct{}

type DescribeThemePresetStmt struct {
    PresetName string
}

type ApplyThemePresetStmt struct {
    PresetName string
}
```

### Implementation Files

| File | Change |
|------|--------|
| `cmd/mxcli/presets/*.mdl` | **New**: Built-in preset MDL scripts |
| `mdl/grammar/MDLParser.g4` | Add SHOW/DESCRIBE/APPLY THEME PRESET rules |
| `mdl/ast/ast_styling.go` | Add AST nodes |
| `mdl/visitor/visitor_styling.go` | Add visitor methods |
| `mdl/executor/cmd_styling.go` | Add `execApplyThemePreset()` (loads + executes script) |
| `mdl/executor/theme_presets.go` | **New file**: Preset discovery, metadata parsing |

---

## Implementation Priority

| Phase | Features | Dependency | Estimated Scope |
|-------|----------|------------|-----------------|
| Phase 1 | Batch Styling + SCSS Variables | Independent | Grammar + 2 executor files |
| Phase 2 | Theme Module Management | Needs settings reader/writer | Extend existing module code |
| Phase 3 | Theme Presets | Depends on Phase 1 (presets use batch + variables) | Thin wrapper + embedded files |

Phase 1 features are independent and can be developed in parallel.
Phase 3 depends on Phase 1 (presets use `ALTER THEME VARIABLES` and batch styling commands).

## New Keywords Summary

| Keyword | Used In | Already Exists? |
|---------|---------|-----------------|
| `THEME` | All theme commands | Check lexer |
| `VARIABLE` / `VARIABLES` | SCSS variable management | Check lexer |
| `WIDGETS` | Batch styling | Check lexer |
| `PAGES` / `SNIPPETS` | Global batch | Check lexer |
| `PRESET` / `PRESETS` | Theme presets | New |
| `APPLY` | Apply preset | Check lexer |
| `RESET` | Reset variable | Check lexer |
| `DEFAULT` | Show defaults | Check lexer |
| `ORDER` | Module order | Check lexer |

## Skill File Updates

After implementation, update `.claude/skills/mendix/theme-styling.md` to document all new commands with examples.
