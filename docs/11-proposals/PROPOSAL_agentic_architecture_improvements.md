# Proposal: Architecture Improvements for Agentic Development and Reduced Merge Conflicts

## Problem Statement

Two related friction points have been identified:

1. **PR merge conflicts on grammar files** — Every MDL feature that adds syntax touches `MDLParser.g4` (4,082 lines) and triggers a full regeneration of `mdl/grammar/parser/mdl_parser.go` (119,737 lines). When two branches both run `make grammar`, the generated file produces an unresolvable conflict. Multiple recent PRs have been blocked or delayed by this.

2. **High agent cognitive load when adding features** — Adding a new MDL command requires coordinated edits across up to six layers (grammar, AST, visitor, executor, backend interface, backend implementation) with no structural guidance on where to look or what to change. Mistakes are caught only by the PR checklist, not by the compiler.

These are separate problems with separate fixes, but both stem from the same root: the codebase has grown to a point where the implicit conventions that work for a small human team become friction for parallel development — whether that parallel work comes from multiple contributors or from AI-assisted generation.

## Current Architecture (What We Have)

The backend abstraction in `mdl/backend/` is already a sound hexagonal boundary:

```
Grammar (.g4)  →  Generated Parser  →  Visitor (AST)  →  Executor  →  Backend Interface  →  MPR Implementation
                  [119k lines, git]                      [241 files]   [mdl/backend/]       [mdl/backend/mpr/]
                                                                        [mock/]
```

The backend layer is correct: executor handlers depend only on the interface, never on the MPR implementation directly. The existing PR checklist enforces this as a human review step.

The problems are:

- The generated parser is in git, so grammar changes always produce a massive conflict
- `MDLParser.g4` is monolithic — any two PRs adding syntax to different domains conflict on the same file
- The executor has 241 files with no structural guide for where to add a new command
- The mock (`mdl/backend/mock/`) is hand-written and must be kept in sync with the interface manually
- Boundary rules (e.g. "no sdk/mpr imports in executor") are enforced by checklist, not by the compiler

## Proposed Changes

### Change 1: Remove Generated Parser Files from Git

**What:** Add `mdl/grammar/parser/*.go` to `.gitignore` and stop tracking them. The `.g4` source files remain in git. Generated files are produced locally and in CI by running `make grammar`.

**Why this eliminates the conflict:** The 119,737-line generated file is the source of every grammar-related merge conflict. The file is deterministically derived from the `.g4` source — there is no information in it that isn't in the grammar. Removing it from git means two branches that both add syntax no longer produce a file-level conflict; only the human-readable `.g4` source can conflict.

**Prerequisites:** Contributors already have ANTLR4 installed (this has been verified — all grammar contributors have been able to run `make grammar` without issue). The barrier that would normally exist (requiring Java for ANTLR4) is not a practical concern here.

**CI impact:** The CI pipeline must install ANTLR4 and run `make grammar` before `make build` or `make test`. The `push-test.yml` and `release.yml` workflows currently have no grammar step — both need a setup step added.

**Implementation steps:**
1. Add to `.gitignore`:
   ```
   # Generated ANTLR4 parser (regenerate with: make grammar)
   mdl/grammar/parser/*.go
   mdl/grammar/parser/*.interp
   mdl/grammar/parser/*.tokens
   ```
2. Run `git rm --cached mdl/grammar/parser/*.go mdl/grammar/parser/*.interp mdl/grammar/parser/*.tokens`
3. Add to CI workflows before the build step:
   ```yaml
   - name: Install ANTLR4
     run: pip install antlr4-tools
   - name: Generate parser
     run: make grammar
   ```
4. Add `make grammar-check` target that regenerates and diffs — to catch cases where grammar source changed but `make grammar` was not run before committing

---

### Change 2: Split the Grammar Source by Domain

**What:** Use ANTLR4's `import` directive to split `MDLParser.g4` into domain-specific grammar files, imported by a thin master grammar.

**Why this reduces the remaining conflict surface:** After Change 1, grammar conflicts can only occur on the `.g4` source files. A monolithic 4,082-line grammar means two PRs touching different domains (e.g. microflows and pages) still conflict at the source level. Splitting by domain means those PRs are fully independent.

**Proposed split:**

```
mdl/grammar/
  MDLParser.g4          # master grammar: imports only, top-level statement rule
  MDLLexer.g4           # unchanged — tokens are shared across domains
  domains/
    MDLDomainModel.g4   # entity, attribute, association rules
    MDLMicroflow.g4     # microflow, nanoflow, activity rules
    MDLPage.g4          # page, snippet, widget rules
    MDLSecurity.g4      # roles, access rules, grant/revoke
    MDLNavigation.g4    # navigation profiles, menus
    MDLWorkflow.g4      # workflow, user task, decision rules
    MDLService.g4       # REST, OData, published services
    MDLCatalog.g4       # catalog queries, show/describe
    MDLSettings.g4      # project settings, constants, enumerations
    MDLSession.g4       # SET, USE, session-level commands
```

The master `MDLParser.g4` becomes:
```antlr
parser grammar MDLParser;
import MDLDomainModel, MDLMicroflow, MDLPage, MDLSecurity,
       MDLNavigation, MDLWorkflow, MDLService, MDLCatalog,
       MDLSettings, MDLSession;
options { tokenVocab=MDLLexer; }
statement : domainStatement | microflowStatement | pageStatement | ... ;
```

**Note:** ANTLR4's `import` merges rules into the main grammar at generation time. The generated parser is still a single file — the split is a source-level convenience only. All existing listener/visitor code continues to work unchanged.

**Implementation steps:**
1. Identify rule boundaries in current `MDLParser.g4` by domain
2. Extract each domain's rules into a `domains/MDL*.g4` file
3. Replace extracted rules in `MDLParser.g4` with `import` directives
4. Run `make grammar` and verify generated parser is functionally identical (test suite must pass)
5. Update `mdl/grammar/Makefile` to list the domain files as dependencies

---

### Change 3: Command Self-Registration in the Executor

**What:** Replace the central dispatch switch in the executor with a registration pattern where each `cmd_*.go` file self-registers its handler on package init.

**Why this helps agentic development:** Currently, adding a new command requires: (a) creating a `cmd_*.go` file, (b) finding and editing the dispatch switch in `executor.go`, and (c) knowing the right context type. An agent has to read across multiple files to understand the pattern. With self-registration, a new command is entirely contained in its own file — no other file needs to change.

**Proposed pattern:**

```go
// In each cmd_*.go file:
func init() {
    executor.Register(ast.KindCreateEntity, handleCreateEntity)
    executor.Register(ast.KindAlterEntity, handleAlterEntity)
}

func handleCreateEntity(ctx *executor.ExecContext, stmt ast.Statement) error {
    s := stmt.(*ast.CreateEntityStatement)
    return ctx.Backend.DomainModel().CreateEntity(s.Name, ...)
}
```

```go
// In executor/registry.go:
var handlers = map[ast.StatementKind]HandlerFunc{}

func Register(kind ast.StatementKind, fn HandlerFunc) {
    handlers[kind] = fn
}
```

**Implementation steps:**
1. Add `executor/registry.go` with `Register()` and the handler map
2. Replace the central dispatch switch with a lookup into the registry
3. Migrate existing `cmd_*.go` files to register via `init()` — one file at a time, verifiable by running tests after each
4. Remove the dispatch switch from `executor.go` once all commands are migrated

---

### Change 4: Code-Generate the Backend Mock

**What:** Replace the hand-written `mdl/backend/mock/` with a generated mock produced by `mockgen` (or equivalent) from the backend interfaces.

**Why this helps:** Adding a new backend method currently requires four coordinated edits: interface definition, MPR implementation, mock stub, and compile-time check. The mock is the most error-prone because it must be kept in sync manually. An agent can easily miss adding the mock stub, causing compilation failures in unrelated tests. Code generation makes the mock always correct by construction.

**Current state:** The mock has 17 files for 26 sub-interfaces — the counts don't align, indicating the mock is already drifting from the interface.

**Implementation steps:**
1. Add `mockgen` (or `moq`) as a Go tool dependency in `go.mod`
2. Add a `//go:generate mockgen ...` directive to each backend interface file in `mdl/backend/`
3. Add `make mocks` target that runs `go generate ./mdl/backend/...`
4. Add generated mock files to `.gitignore` (same rationale as the parser: derived from source)
5. Add `make mocks` step to CI before `make test`
6. Delete the hand-written mock files in `mdl/backend/mock/`

---

### Change 5: Compiler-Enforced Backend Boundary

**What:** Move the `sdk/mpr` write types behind Go's `internal/` package visibility rules so that importing them from `mdl/executor/` is a compile error, not a checklist item.

**Why this helps:** The PR checklist currently has: *"No sdk/mpr write imports in executor."* This rule exists because executor handlers must go through the backend interface, not call the MPR writer directly. Today this is caught only in code review. An agent that bypasses the backend will produce code that compiles but violates the architecture — and the mistake ships unless a reviewer catches it.

**Proposed structure:**

```
sdk/
  mpr/
    internal/
      writer/     # write types moved here — only sdk/mpr can import
      parser/     # parser types moved here
    reader.go     # public read API (no change)
    writer.go     # delegates to internal/writer
```

Because Go's `internal/` rule allows only the parent package and its children to import, any attempt by `mdl/executor/` to import `sdk/mpr/internal/writer` will be a compile error.

**Implementation steps:**
1. Create `sdk/mpr/internal/writer/` and `sdk/mpr/internal/parser/`
2. Move write-path types into `internal/writer/`
3. Keep `sdk/mpr/writer.go` as the public facade (re-exports what `mdl/backend/mpr/` needs)
4. Verify `mdl/backend/mpr/` still compiles (it is a child of `sdk/mpr` — wait, it is not; it's under `mdl/`)

**Note:** This requires careful analysis. `mdl/backend/mpr/` is the legitimate consumer of the MPR write API, but it lives under `mdl/`, not `sdk/mpr/`. Go's `internal/` rule would block it too. An alternative is to enforce the boundary via a linting rule (`depguard` or `gomodguard`) rather than the compiler — flagging any import of specific `sdk/mpr` write symbols from `mdl/executor/` in CI.

---

## Summary and Priority

| Change | Problem solved | Risk | Effort |
|---|---|---|---|
| 1. Gitignore generated parser | Eliminates 119k-line conflict entirely | Low | Low |
| 2. Split grammar by domain | Reduces source-level grammar conflicts | Medium (grammar refactor) | Medium |
| 3. Command self-registration | Reduces agent discovery cost | Medium (executor refactor) | Medium |
| 4. Code-generate mock | Eliminates mock drift, reduces sync errors | Low | Low |
| 5. Compiler-enforced boundary | Converts checklist rule to compile error | Medium (needs design) | Medium |

**Recommended order:** Changes 1 and 4 first — both are low risk, high value, and immediately unblock parallel development. Changes 2 and 3 next, ideally in separate PRs. Change 5 needs a design decision on the `internal/` boundary approach before implementation.
