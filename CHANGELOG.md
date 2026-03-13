# Changelog

All notable changes to mxcli will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.2.0] - 2026-03-10

### Added

- **Page Variables** — `Variables: { $name: Type = 'expression' }` in page/snippet headers for column visibility and conditional logic (`Forms$LocalVariable`)
- **SHOW DATABASE CONNECTIONS** — List all database connections in the project
- **DESCRIBE MODULE** — Now includes database connections
- **Catalog: role_mappings** — New `CATALOG.role_mappings` table for querying role assignments
- **Starlark: permissions()** — New `permissions()` function for all element types in lint rules
- **Starlark: module_roles** — Expose module roles to Starlark API
- **Short attribute names** — Gallery TEXTFILTER `Attributes:` now resolves short attribute names from entity context
- **Folder in DESCRIBE** — DESCRIBE PAGE/SNIPPET/MICROFLOW now includes Folder in output
- **MOVE examples** — MOVE to folder, module, and cross-module examples in page examples

### Fixed

- **Microflow/nanoflow datasource** — Fix serialization for pluggable widgets (DATAGRID, GALLERY)
- **DESCRIBE PAGE attributes** — Output short attribute names instead of fully qualified names
- **Null AttributeRef crash** — Resolve microflow return types from exec cache for ContentParams
- **DataGrid2 ShowNumberOfRows** — Removed unsupported property that caused CE0463 widget definition error
- **DataGrid column Visible** — Document that column visibility requires page variables, not `$currentObject`
- **DynamicCellClass syntax** — Fix expression syntax to use `if(...) then ... else ...` (not function-call style)
- **DROP MODULE cleanup** — Fully clean up mxunit files, themesource dirs, and all document types
- **Business event nil panic** — Fix nil pointer panic in create/drop when hierarchy fails
- **SHOW CATALOG TABLES** — Fix parsing for catalog table queries
- **Project explorer** — Show integration documents in folder hierarchy; add missing document types

### Changed

- **Snippet help** — Updated from V2 to V3 syntax with Variables example
- **Page help** — Added Variables syntax and usage notes
- **Documentation** — Updated architecture docs, feature matrix, quick reference, and migration guide

## [0.1.0] - 2026-03-06

First public release.

### Added

- **MDL Language** — SQL-like syntax (Mendix Definition Language) for querying and modifying Mendix projects
- **Domain Model** — CREATE/ALTER/DROP ENTITY, CREATE ASSOCIATION, attribute types, indexes, validation rules
- **Microflows & Nanoflows** — 60+ activity types, loops, error handling, expressions, parameters
- **Pages** — 50+ widget types, CREATE/ALTER PAGE/SNIPPET, DataGrid, DataView, ListView, pluggable widgets
- **Security** — Module roles, entity access rules, GRANT/REVOKE, UPDATE SECURITY reconciliation
- **Navigation** — Navigation profiles, menu items, home pages, login pages
- **Enumerations** — CREATE/ALTER/DROP ENUMERATION with localized values
- **Business Events** — CREATE/DROP business event services
- **Project Settings** — SHOW/DESCRIBE/ALTER for runtime, language, and theme settings
- **Database Connections** — CREATE/DESCRIBE DATABASE CONNECTION for Database Connector module
- **Full-text Search** — SEARCH across all strings, messages, captions, labels, and MDL source
- **Code Navigation** — SHOW CALLERS/CALLEES/REFERENCES/IMPACT/CONTEXT for cross-reference analysis
- **Catalog Queries** — SQL-based querying of project metadata via CATALOG tables
- **Linting** — 14 built-in rules + 27 Starlark rules across MDL, SEC, QUAL, ARCH, DESIGN, CONV categories
- **Report** — Scored best practices report with category breakdown (`mxcli report`)
- **Testing** — `.test.mdl` / `.test.md` test files with Docker-based runtime validation
- **Diff** — Compare MDL scripts against project state, git diff for MPR v2 projects
- **External SQL** — Direct queries against PostgreSQL, Oracle, SQL Server with credential isolation
- **Data Import** — IMPORT FROM external DB into Mendix app PostgreSQL with batch insert and ID generation
- **Connector Generation** — Auto-generate Database Connector MDL from external schema discovery
- **OQL** — Query running Mendix runtime via admin API
- **Docker Build** — `mxcli docker build` with PAD patching
- **VS Code Extension** — Syntax highlighting, diagnostics, completion, hover, go-to-definition, symbols, folding
- **LSP Server** — `mxcli lsp --stdio` for editor integration
- **Multi-tool Init** — `mxcli init` with support for Claude Code, Cursor, Continue.dev, Windsurf, Aider
- **Dev Container** — `mxcli init` generates `.devcontainer/` configuration for sandboxed AI agent development
- **MPR v1/v2** — Automatic format detection, read/write support for both formats
- **Fluent API** — High-level Go API (`api/` package) for programmatic model manipulation
