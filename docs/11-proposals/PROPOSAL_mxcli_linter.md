# mxcli Linter - Implementation Proposal (Revised)

## Overview

Add an extensible linting framework to mxcli that leverages the existing SQLite-based catalog system. Rules can be written in Go (built-in) or Starlark (user-defined).

**Dependency:** This proposal builds on top of the [Code Search Implementation](./code-search-implementation.md) proposal, which adds the `references` table required for unused element detection and cross-reference analysis.

## Current State Analysis

### Existing Infrastructure

| Component | Location | What It Provides |
|-----------|----------|------------------|
| **Catalog** | `mdl/catalog/` | In-memory SQLite with modules, entities, microflows, pages, etc. |
| **Catalog Builder** | `mdl/catalog/builder*.go` | Populates catalog from MPR via `REFRESH CATALOG [FULL]` |
| **CLI Structure** | `cmd/mxcli/main.go` | Cobra-based commands with `-p` project flag |
| **Init System** | `cmd/mxcli/init.go` | Creates `.claude/` folder with skills |
| **Executor** | `mdl/executor/` | Executes MDL statements, has access to catalog |

### Existing Catalog Schema (from `mdl/catalog/tables.go`)

```
modules      - Id, Name, QualifiedName, ModuleName, Description, IsSystemModule, ...
entities     - Id, Name, QualifiedName, ModuleName, EntityType, Description, Generalization, AttributeCount, ...
microflows   - Id, Name, QualifiedName, ModuleName, MicroflowType, Description, ReturnType, ParameterCount, ActivityCount, ...
pages        - Id, Name, QualifiedName, ModuleName, Title, URL, LayoutRef, Description, WidgetCount, ...
snippets     - Id, Name, QualifiedName, ModuleName, Description, WidgetCount, ...
enumerations - Id, Name, QualifiedName, ModuleName, Description, ValueCount, ...
activities   - Id, Name, Caption, ActivityType, MicroflowId, EntityRef, ActionType, ... (FULL mode only)
widgets      - Id, Name, WidgetType, ContainerId, EntityRef, AttributeRef, ... (FULL mode only)
objects      - VIEW: union of all types above
```

### What's Missing (Required from code-search-implementation.md)

```sql
-- References table (from code-search-implementation.md proposal)
CREATE TABLE references (
    Id TEXT PRIMARY KEY,
    SourceType TEXT NOT NULL,
    SourceId TEXT NOT NULL,
    SourceName TEXT NOT NULL,
    TargetType TEXT NOT NULL,
    TargetId TEXT,
    TargetName TEXT NOT NULL,
    RefKind TEXT NOT NULL,  -- 'call', 'create', 'retrieve', 'change', 'show_page', etc.
    ...
);
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           CLI Layer                                  │
│  mxcli lint [-p app.mpr] [--format text|json|sarif] [--config path] │
│  MDL: LINT [MODULE.* | *] [FORMAT text|json|sarif]                  │
└─────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          Linter Core                                 │
│  mdl/linter/linter.go                                               │
│  - Rule registration & orchestration                                 │
│  - Parallel rule execution (goroutines + semaphore)                 │
│  - Configuration management (.claude/lint-config.yaml)              │
│  - Output formatting (text, JSON, SARIF)                            │
└─────────────────────────────────────────────────────────────────────┘
                │                               │
                ▼                               ▼
┌───────────────────────────┐     ┌─────────────────────────────┐
│   Built-in Go Rules       │     │   Starlark Rule Engine      │
│   mdl/linter/rules/       │     │   mdl/linter/starlark.go    │
│   - MDL001: NamingConv    │     │   - Script loading          │
│   - MDL002: EmptyMicroflow│     │   - API bindings            │
│   - MDL003: UnusedEntity  │     │   - Sandboxed execution     │
│   - MDL004: CircularDeps  │     │   (uses go.starlark.net)    │
└───────────────────────────┘     └─────────────────────────────┘
                │                               │
                └───────────────┬───────────────┘
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          LintContext                                 │
│  mdl/linter/context.go                                              │
│  Wraps *catalog.Catalog and provides rule-friendly API:             │
│  - Entities(), Microflows(), Pages() - iterators                    │
│  - FindReferences(id) - requires references table                   │
│  - FindUnused(kind) - requires references table                     │
│  - ModuleDependencies() - derived from references                   │
│  - Query(sql) - raw SQL access                                      │
└─────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    SQLite Catalog (In-Memory)                        │
│  Built by: REFRESH CATALOG FULL                                      │
│  Tables: modules, entities, microflows, pages, activities, widgets   │
│  Required: references table (from code-search proposal)              │
└─────────────────────────────────────────────────────────────────────┘
```

## File Structure

### New Package: `mdl/linter/`

```
mdl/linter/
├── linter.go          # Core orchestration, Linter struct, Run()
├── context.go         # LintContext wrapping catalog, query methods
├── config.go          # Config loading from .claude/lint-config.yaml
├── output.go          # Formatters: Text, JSON, SARIF
├── starlark.go        # Starlark rule engine, API bindings
├── rules/             # Built-in Go rules
│   ├── naming.go      # MDL001: NamingConvention
│   ├── empty.go       # MDL002: EmptyMicroflow
│   ├── unused.go      # MDL003: UnusedEntity (requires references)
│   ├── circular.go    # MDL004: CircularDependency (requires references)
│   └── security.go    # MDL005: SecurityCheck
└── testdata/          # Test fixtures
```

### User Custom Rules Location

Following the existing `.claude/` pattern from `mxcli init`:

```
<mendix-project>/
├── .claude/
│   ├── settings.json          # Existing Claude Code settings
│   ├── skills/                # Existing skills
│   ├── commands/              # Existing commands
│   ├── lint-config.yaml       # NEW: Linter configuration
│   └── lint-rules/            # NEW: Custom Starlark rules
│       ├── naming.star
│       └── architecture.star
├── CLAUDE.md
└── MyApp.mpr
```

## Usage

### CLI Commands

```bash
# Lint with default rules (requires catalog)
mxcli lint -p app.mpr

# Output as JSON
mxcli lint -p app.mpr --format json

# Output as SARIF (for GitHub/IDE integration)
mxcli lint -p app.mpr --format sarif > results.sarif

# Use custom config file
mxcli lint -p app.mpr --config ./my-lint-config.yaml

# List available rules
mxcli lint -p app.mpr --list-rules

# Lint specific module(s)
mxcli lint -p app.mpr --module Sales --module Customers
```

### MDL REPL Commands

```sql
-- Build catalog first (required for linting)
REFRESH CATALOG FULL;

-- Lint all modules
LINT;

-- Lint specific module
LINT Sales.*;

-- Lint with format
LINT FORMAT json;
LINT FORMAT sarif;

-- Show available rules
SHOW LINT RULES;
```

## Configuration

**File:** `.claude/lint-config.yaml`

```yaml
# mxcli Linter Configuration
linter:
  # Directory containing custom Starlark rules (relative to project)
  rules_dir: ".claude/lint-rules"

  # Output format: text, json, sarif
  output_format: text

  # Modules/patterns to exclude from linting
  exclude:
    - "System"
    - "Administration"
    - "Atlas_*"
    - "*_Generated"

  # Per-rule configuration
  rules:
    MDL001:  # NamingConvention
      enabled: true
      severity: warning
      options:
        entity_pattern: "^[A-Z][a-zA-Z0-9]*$"
        microflow_pattern: "^(ACT_|SUB_|DS_|VAL_|SCH_)?[A-Z][a-zA-Z0-9_]*$"

    MDL002:  # EmptyMicroflow
      enabled: true
      severity: warning

    MDL003:  # UnusedEntity
      enabled: true
      severity: info

    MDL004:  # CircularDependency
      enabled: true
      severity: error

    MDL005:  # SecurityCheck
      enabled: false  # Disabled - requires security model support

    # Custom rule
    CUSTOM_001:
      enabled: true
      severity: warning
```

## Implementation Details

### 1. Core Types (`mdl/linter/linter.go`)

```go
package linter

type Severity int

const (
    SeverityHint Severity = iota
    SeverityInfo
    SeverityWarning
    SeverityError
)

type Violation struct {
    RuleID      string
    Severity    Severity
    Message     string
    Location    Location
    Suggestion  string
}

type Location struct {
    Module       string  // e.g., "Sales"
    DocumentType string  // "entity", "microflow", "page"
    DocumentName string  // e.g., "Customer"
    DocumentID   string  // UUID
}

type Rule interface {
    ID() string
    Name() string
    Description() string
    Severity() Severity
    Category() string
    Check(ctx *LintContext) []Violation
}

type Linter struct {
    ctx       *LintContext
    rules     []Rule
    config    *Config
}

func (l *Linter) Run(ctx context.Context) ([]Violation, error) {
    // Parallel rule execution with semaphore
}
```

### 2. LintContext (`mdl/linter/context.go`)

Wraps the existing `*catalog.Catalog` and provides a rule-friendly API:

```go
package linter

import (
    "database/sql"
    "iter"
    "github.com/mendixlabs/mxcli/mdl/catalog"
)

type LintContext struct {
    catalog   *catalog.Catalog
    db        *sql.DB  // Direct access for complex queries
    config    *Config
    excluded  map[string]bool  // Excluded modules
}

// NewLintContext creates context from existing catalog
func NewLintContext(cat *catalog.Catalog, cfg *Config) *LintContext {
    return &LintContext{
        catalog: cat,
        db:      cat.DB(),
        config:  cfg,
        excluded: buildExcludeMap(cfg.Exclude),
    }
}

// Entity represents a lintable entity
type Entity struct {
    ID             string
    Name           string
    QualifiedName  string
    ModuleName     string
    EntityType     string  // "Persistent", "NonPersistent", "View"
    Description    string
    Generalization string
    AttributeCount int
}

// Entities returns an iterator over all entities (excluding system modules)
func (ctx *LintContext) Entities() iter.Seq[Entity] {
    return func(yield func(Entity) bool) {
        rows, err := ctx.db.Query(`
            SELECT Id, Name, QualifiedName, ModuleName, EntityType,
                   Description, Generalization, AttributeCount
            FROM entities
            WHERE ModuleName NOT IN (SELECT Name FROM modules WHERE IsSystemModule = 1)
            ORDER BY ModuleName, Name
        `)
        if err != nil {
            return
        }
        defer rows.Close()

        for rows.Next() {
            var e Entity
            var desc, gen sql.NullString
            rows.Scan(&e.ID, &e.Name, &e.QualifiedName, &e.ModuleName,
                      &e.EntityType, &desc, &gen, &e.AttributeCount)
            e.Description = desc.String
            e.Generalization = gen.String

            if ctx.excluded[e.ModuleName] {
                continue
            }

            if !yield(e) {
                return
            }
        }
    }
}

// Similar iterators for Microflows(), Pages(), etc.

// Query executes raw SQL (for advanced rules)
func (ctx *LintContext) Query(query string, args ...any) (*sql.Rows, error) {
    return ctx.db.Query(query, args...)
}

// FindReferences finds all references to a given element
// Requires: references table from code-search-implementation.md
func (ctx *LintContext) FindReferences(targetID string) []Reference {
    // ...
}

// FindUnused finds elements with no incoming references
// Requires: references table from code-search-implementation.md
func (ctx *LintContext) FindUnused(kind string) []Symbol {
    // ...
}

// ModuleDependencies returns module dependency graph
// Requires: references table from code-search-implementation.md
func (ctx *LintContext) ModuleDependencies() map[string][]string {
    // ...
}
```

### 3. Built-in Rules (`mdl/linter/rules/`)

**MDL001: NamingConvention** - Works with current schema
```go
func (r *NamingConventionRule) Check(ctx *LintContext) []Violation {
    var violations []Violation

    for entity := range ctx.Entities() {
        if !isPascalCase(entity.Name) {
            violations = append(violations, Violation{
                RuleID:  "MDL001",
                Message: fmt.Sprintf("Entity name '%s' should use PascalCase", entity.Name),
                Location: Location{
                    Module:       entity.ModuleName,
                    DocumentType: "entity",
                    DocumentName: entity.Name,
                    DocumentID:   entity.ID,
                },
                Suggestion: toPascalCase(entity.Name),
            })
        }
    }

    return violations
}
```

**MDL002: EmptyMicroflow** - Works with current schema
```go
func (r *EmptyMicroflowRule) Check(ctx *LintContext) []Violation {
    var violations []Violation

    for mf := range ctx.Microflows() {
        if mf.ActivityCount == 0 {
            violations = append(violations, Violation{
                RuleID:  "MDL002",
                Message: fmt.Sprintf("Microflow '%s' has no activities", mf.Name),
                Location: Location{
                    Module:       mf.ModuleName,
                    DocumentType: "microflow",
                    DocumentName: mf.Name,
                    DocumentID:   mf.ID,
                },
                Suggestion: "Add activities or remove unused microflow",
            })
        }
    }

    return violations
}
```

**MDL003: UnusedEntity** - Requires references table
```go
func (r *UnusedEntityRule) Check(ctx *LintContext) []Violation {
    // This rule requires the references table from code-search-implementation.md
    if !ctx.HasReferencesTable() {
        return nil  // Skip if references not available
    }

    unused := ctx.FindUnused("entity")
    // ...
}
```

### 4. Starlark Integration (`mdl/linter/starlark.go`)

```go
package linter

import (
    "go.starlark.net/starlark"
    "go.starlark.net/starlarkstruct"
)

type StarlarkRule struct {
    meta    RuleMeta
    path    string
    ctx     *LintContext
    checkFn starlark.Value
}

func LoadStarlarkRule(path string, ctx *LintContext) (*StarlarkRule, error) {
    // Load and parse .star file
    // Extract metadata: RULE_ID, RULE_NAME, DESCRIPTION, SEVERITY, CATEGORY
    // Get check() function
}

func (r *StarlarkRule) buildPredeclared() starlark.StringDict {
    return starlark.StringDict{
        // Query functions - map to LintContext methods
        "entities":     starlark.NewBuiltin("entities", r.builtinEntities),
        "microflows":   starlark.NewBuiltin("microflows", r.builtinMicroflows),
        "pages":        starlark.NewBuiltin("pages", r.builtinPages),
        "query":        starlark.NewBuiltin("query", r.builtinQuery),

        // Violation helpers
        "violation":    starlark.NewBuiltin("violation", r.builtinViolation),
        "location":     starlark.NewBuiltin("location", r.builtinLocation),

        // String utilities
        "is_pascal_case": starlark.NewBuiltin("is_pascal_case", builtinIsPascalCase),
        "is_camel_case":  starlark.NewBuiltin("is_camel_case", builtinIsCamelCase),
        "matches":        starlark.NewBuiltin("matches", builtinMatches),
    }
}
```

### 5. CLI Integration (`cmd/mxcli/lint.go`)

```go
var lintCmd = &cobra.Command{
    Use:   "lint",
    Short: "Lint a Mendix project for issues",
    Long: `Run linting rules against a Mendix project to find potential issues.

Built-in rules check for:
  - Naming conventions (MDL001)
  - Empty microflows (MDL002)
  - Unused entities (MDL003) - requires REFRESH CATALOG FULL
  - Circular dependencies (MDL004) - requires references table

Custom rules can be added as Starlark scripts in .claude/lint-rules/

Examples:
  mxcli lint -p app.mpr
  mxcli lint -p app.mpr --format sarif > results.sarif
  mxcli lint -p app.mpr --module Sales
`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    lintCmd.Flags().StringP("format", "f", "text", "Output format: text, json, sarif")
    lintCmd.Flags().StringP("config", "c", "", "Config file path")
    lintCmd.Flags().StringSliceP("module", "m", nil, "Modules to lint (default: all)")
    lintCmd.Flags().Bool("list-rules", false, "List available rules")

    rootCmd.AddCommand(lintCmd)
}
```

### 6. MDL Grammar Extensions (`mdl/grammar/MDLParser.g4`)

```antlr
// Add to statement rule:
statement
    : ...
    | lintStatement
    ;

lintStatement
    : LINT (qualifiedNamePattern)? (FORMAT lintFormat)?
    | SHOW LINT RULES
    ;

qualifiedNamePattern
    : qualifiedName              // Specific element
    | moduleName DOT STAR        // All in module
    | STAR                       // All
    ;

lintFormat
    : TEXT | JSON | SARIF
    ;

// New keywords
LINT: L I N T;
SARIF: S A R I F;
```

### 7. Executor Integration (`mdl/executor/cmd_lint.go`)

```go
func (e *Executor) execLint(s *ast.LintStmt) error {
    // Ensure catalog is built (FULL mode for activities)
    if e.catalog == nil || !e.catalog.IsBuilt() {
        if err := e.buildCatalog(true); err != nil {
            return err
        }
    }

    // Load config
    cfg, err := linter.LoadConfig(e.findLintConfig())
    if err != nil {
        return err
    }

    // Create linter
    ctx := linter.NewLintContext(e.catalog, cfg)
    lint := linter.New(ctx)

    // Load rules
    lint.LoadBuiltinRules()
    if err := lint.LoadCustomRules(cfg.RulesDir); err != nil {
        return fmt.Errorf("loading custom rules: %w", err)
    }

    // Run
    violations, err := lint.Run(context.Background())
    if err != nil {
        return err
    }

    // Output
    formatter := linter.GetFormatter(s.Format, true)
    return formatter.Format(violations, e.output)
}
```

## Output Examples

### Text Format (Default)

```
Sales
-----
  ⚠ Entity name 'customer_info' should use PascalCase [MDL001]
      at Sales.customer_info
      → CustomerInfo

  ⚠ Microflow 'test' has no activities [MDL002]
      at Sales.test
      → Add activities or remove unused microflow

MyModule
--------
  ℹ Entity 'TempData' is not referenced anywhere [MDL003]
      at MyModule.TempData
      → Remove entity or add references

3 issues: 0 errors, 2 warnings, 1 info
```

### SARIF Format (for CI/GitHub)

```json
{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "version": "2.1.0",
  "runs": [{
    "tool": {
      "driver": {
        "name": "mxcli-lint",
        "version": "0.1.0",
        "rules": [
          {"id": "MDL001", "shortDescription": {"text": "NamingConvention"}, ...},
          {"id": "MDL002", "shortDescription": {"text": "EmptyMicroflow"}, ...}
        ]
      }
    },
    "results": [
      {
        "ruleId": "MDL001",
        "level": "warning",
        "message": {"text": "Entity name 'customer_info' should use PascalCase"},
        "locations": [{"physicalLocation": {"artifactLocation": {"uri": "Sales/customer_info"}}}]
      }
    ]
  }]
}
```

## Custom Rule Example

**File:** `.claude/lint-rules/entity_documentation.star`

```python
# Entity Documentation Rule
# Checks that all entities have documentation

RULE_ID = "CUSTOM_001"
RULE_NAME = "EntityDocumentation"
DESCRIPTION = "Entities should have documentation explaining their purpose"
SEVERITY = "info"
CATEGORY = "documentation"

def check():
    """Ensure all entities have documentation."""
    violations = []

    for entity in entities():
        # Check if documentation is empty
        if not entity.description or not entity.description.strip():
            violations.append(violation(
                message = "Entity '{}' has no documentation".format(entity.name),
                location = location(
                    module = entity.module_name,
                    document_type = "entity",
                    document_name = entity.name,
                    id = entity.id
                ),
                suggestion = "Add a description using /** ... */ comment"
            ))

    return violations
```

## Implementation Plan

### Phase 1: Core Framework (No Dependencies)

1. **Create `mdl/linter/` package** with basic types
2. **Implement LintContext** wrapping catalog (queries that work with current schema)
3. **Implement built-in rules** that don't require references:
   - MDL001: NamingConvention
   - MDL002: EmptyMicroflow
4. **Add text output formatter**
5. **Add `lint` CLI command**

### Phase 2: MDL Integration

6. **Add MDL grammar** for LINT statement
7. **Add executor command** `cmd_lint.go`
8. **Add configuration** loading from `.claude/lint-config.yaml`

### Phase 3: Starlark Rules

9. **Add Starlark engine** using `go.starlark.net`
10. **Implement API bindings** for custom rules
11. **Add rule loader** for `.claude/lint-rules/`

### Phase 4: Advanced Rules (Requires code-search proposal)

12. **Implement MDL003: UnusedEntity** (requires references table)
13. **Implement MDL004: CircularDependency** (requires references table)
14. **Add JSON and SARIF formatters**

### Phase 5: CI/CD Integration

15. **Add exit codes** (1 for errors)
16. **Add `--list-rules` option**
17. **Document GitHub Actions integration**

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `mdl/linter/linter.go` | Create | Core linter orchestration |
| `mdl/linter/context.go` | Create | LintContext wrapping catalog |
| `mdl/linter/config.go` | Create | Configuration loading |
| `mdl/linter/output.go` | Create | Text/JSON/SARIF formatters |
| `mdl/linter/starlark.go` | Create | Starlark rule engine |
| `mdl/linter/rules/naming.go` | Create | MDL001 rule |
| `mdl/linter/rules/empty.go` | Create | MDL002 rule |
| `mdl/linter/rules/unused.go` | Create | MDL003 rule |
| `mdl/linter/rules/circular.go` | Create | MDL004 rule |
| `cmd/mxcli/lint.go` | Create | CLI lint command |
| `cmd/mxcli/main.go` | Modify | Register lint command |
| `mdl/grammar/MDLParser.g4` | Modify | Add LINT statement |
| `mdl/ast/ast_lint.go` | Create | AST types for LINT |
| `mdl/visitor/visitor_lint.go` | Create | Parse LINT statements |
| `mdl/executor/cmd_lint.go` | Create | Execute LINT statements |
| `go.mod` | Modify | Add `go.starlark.net` dependency |

## Testing Strategy

1. **Unit tests** for each rule
2. **Integration tests** with sample MPR files
3. **Starlark rule tests** with fixture scripts
4. **CLI tests** for output formats

## Open Questions

1. **Should `mxcli init` be extended** to also create `.claude/lint-config.yaml` and example custom rules?

2. **Should we add a `lint init` subcommand** to create the config and rules directory?

3. **What's the priority** for SARIF output vs getting core rules working first?

4. **Should rules have access to the raw MPR** (via reader) for deeper analysis, or just the catalog?
