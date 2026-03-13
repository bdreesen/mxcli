# Contributing to mxcli

Thank you for your interest in contributing to mxcli! This guide will help you get started.

## Prerequisites

- **Go 1.26+** — the project uses pure Go with no CGO dependency
- **bun** — used for the VS Code extension (`vscode-mdl/`); do not use npm/node
- **Docker** — required for integration tests and running Mendix apps
- **ANTLR4** — only needed if modifying the MDL grammar (`.g4` files)

## Getting Started

```bash
# Clone the repository
git clone https://github.com/mendixlabs/mxcli.git
cd mxcli

# Build
make build

# Run tests
make test
```

The binary is output to `./bin/mxcli`.

## Project Structure

See [CLAUDE.md](CLAUDE.md) for a detailed architecture overview. Key directories:

| Directory | Description |
|-----------|-------------|
| `cmd/mxcli/` | CLI entry point (Cobra commands) |
| `sdk/` | Core SDK: domain model, microflows, pages, MPR reader/writer |
| `mdl/` | MDL parser (ANTLR4 grammar, AST, visitor, executor) |
| `api/` | High-level fluent API |
| `sql/` | External database connectivity |
| `vscode-mdl/` | VS Code extension for MDL language support |

## Development Workflow

### Building

```bash
make build          # Build mxcli + source_tree (syncs embedded assets automatically)
make release        # Cross-compile for all platforms
make vscode-ext     # Build the VS Code extension (requires bun)
```

### Testing

```bash
make test           # Unit tests
make test-mdl MPR=path/to/app.mpr   # MDL integration tests (requires Docker)
```

### Code Quality

```bash
make fmt            # Format code
make vet            # Run go vet
```

### Grammar Changes

If you modify `mdl/grammar/MDLLexer.g4` or `mdl/grammar/MDLParser.g4`:

```bash
make grammar        # Regenerate the ANTLR4 parser
```

See `docs/03-development/MDL_PARSER_ARCHITECTURE.md` for parser design details.

### VS Code Extension

The extension lives in `vscode-mdl/` and uses **bun** (not npm):

```bash
cd vscode-mdl
bun install
bun run compile
```

Or use the Makefile:

```bash
make vscode-ext       # Build .vsix package
make vscode-install   # Build and install into VS Code
```

## Code Style

- Follow standard Go conventions (`go fmt`, `go vet`)
- Use descriptive names matching Mendix terminology
- Add `// SPDX-License-Identifier: Apache-2.0` as the first line of every new source file
- Keep BSON/JSON tags consistent with Mendix serialization format
- Export types that should be part of the public API

## Key Conventions

### BSON Storage Names

Mendix uses different "storage names" in BSON `$Type` fields than the "qualified names" in SDK documentation. Always verify storage names against:
1. Existing MPR files
2. Reflection data in `reference/mendixmodellib/reflection-data/`
3. Parser cases in `sdk/mpr/parser_microflow.go`

See the table in [CLAUDE.md](CLAUDE.md) for common mappings.

### Pluggable Widget Templates

Templates in `sdk/widgets/templates/` must include both `type` (PropertyTypes schema) and `object` (default WidgetObject) fields. Always extract from Studio Pro, not generated programmatically.

### MDL Scripts

Before writing MDL scripts, read the relevant skill files in `.claude/skills/mendix/` for syntax reference and common pitfalls. Validate with:

```bash
./bin/mxcli check script.mdl                          # Syntax check
./bin/mxcli check script.mdl -p app.mpr --references  # With reference validation
```

## Submitting Changes

1. Create a feature branch from `main`
2. Make your changes, ensuring tests pass (`make test`)
3. Add tests for new functionality
4. Verify the build succeeds (`make build`)
5. Submit a pull request with a clear description of the changes

## Reporting Issues

Please report bugs and feature requests via the project's issue tracker. Include:
- Steps to reproduce
- Expected vs actual behavior
- Mendix version (if relevant)
- Output of `mxcli --version`

## License

By contributing, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).
