# Proposal: Comprehensive Nanoflow Support

## Overview

**Status:** In Progress (PR chain: #7 → #8)
**Priority:** High — nanoflows are heavily used (227 across test projects) and CLI parity with microflows is expected.

This proposal covers the full nanoflow feature surface in mxcli: CREATE, DROP, CALL, GRANT/REVOKE, SHOW, DESCRIBE, and validation. It supersedes the earlier `show-describe-nanoflows.md` proposal which focused only on DESCRIBE/DROP (both now implemented).

## Background

Nanoflows execute client-side (browser or native app). In the Mendix metamodel, `Nanoflow` inherits from `MicroflowBase` (not `ServerSideMicroflow`), sharing the same flow structure (`MicroflowObjectCollection`, `SequenceFlow`, action types) but with restricted action set and different properties.

### Nanoflow vs Microflow Model

| Property | Nanoflow | Microflow |
|----------|----------|-----------|
| Inheritance | `MicroflowBase` (direct) | `ServerSideMicroflow` → `MicroflowBase` |
| `AllowedModuleRoles` | Yes (design-time only) | Yes (runtime enforced) |
| `ApplyEntityAccess` | No | Yes |
| `AllowConcurrentExecution` | No | Yes |
| `ConcurrencyErrorMessage` | No | Yes |
| `MicroflowActionInfo` | No | Yes |
| `WorkflowActionInfo` | No | Yes |
| `Url*` / `StableId` | No | Yes |
| `UseListParameterByReference` | Yes (default true) | No |
| `ReturnVariableName` | Inherited from `MicroflowBase` but not exposed in Studio Pro | Yes |
| Allowed return types | No `Binary`, no `Float` | All types |
| `ErrorEvent` | Forbidden | Allowed |
| Expression context | `ClientExpressionContext` | `MicroflowExpressionContext` |
| Predefined variables | `$latestError` (String) | (microflow-specific set) |

### Action Restrictions

**Allowed in nanoflows** (25 actions): ChangeVariable, AggregateList, CreateVariable, Rollback, Retrieve, Delete, CreateChange, Commit, Cast, Change, LogMessage, ListOperation, CreateList, ChangeList, MicroflowCall, ValidationFeedback, ShowPage, ShowMessage, CloseForm, **NanoflowCall**, **JavaScriptActionCall**, **Synchronize**, **CancelSynchronization**, **ClearFromClient**.

**Disallowed** (32+ actions): All Java actions, REST calls, workflow actions, import/export, external object ops, download file, push to client, show home page, email, document generation, metrics, ML model calls.

## Implementation Status

### Completed (PR #7 — `pr4-nanoflows-create-drop`)

| Feature | Layer | Status |
|---------|-------|--------|
| CREATE NANOFLOW | Grammar, AST, Visitor, Executor | Done |
| DROP NANOFLOW | Grammar (existed), AST, Visitor, Executor | Done |
| Nanoflow action validation | Executor (`nanoflow_validation.go`) | Done |
| Return type validation (reject Binary) | Executor | Done |
| Executor cache (created/dropped nanoflows) | Executor | Done |

### Completed (PR #8 — `pr5-nanoflows-call-grant`)

| Feature | Layer | Status |
|---------|-------|--------|
| CALL NANOFLOW (inside flow body) | Grammar, AST, Visitor, Flow builder | Done |
| GRANT EXECUTE ON NANOFLOW | Grammar, AST, Visitor, Executor | Done |
| REVOKE EXECUTE ON NANOFLOW | Grammar, AST, Visitor, Executor | Done |
| `NanoflowCallAction` SDK type | SDK types, BSON parser/writer | Done |
| `AllowedModuleRoles` on Nanoflow | SDK struct, BSON parser/writer | Done |
| Cross-reference validation | `validate.go` | Done |
| Flow body validation | `validate_microflow.go` | Done |
| MDL diff support | `cmd_diff_mdl.go` | Done |
| Statement summary | `stmt_summary.go` | Done |
| Agentic skill | `write-nanoflows.md` | Done |

### Already Working (no changes needed)

| Feature | Notes |
|---------|-------|
| SHOW NANOFLOWS | Lists nanoflows with activity counts |
| DESCRIBE NANOFLOW | Outputs MDL representation |
| RENAME NANOFLOW | Grammar already supports both |
| MOVE NANOFLOW | Grammar already supports both |
| Flow builder reuse | `flowBuilder` + `buildFlowGraph` work for both |
| SDK backend | `ListNanoflows`, `GetNanoflow`, `CreateNanoflow`, `UpdateNanoflow`, `DeleteNanoflow`, `MoveNanoflow` all exist |
| Linter | Iterates both microflows and nanoflows via shared type |

### Not Planned (by design)

| Feature | Reason |
|---------|--------|
| HOME NANOFLOW (navigation) | Home page/microflow is server-side |
| MENU ITEM NANOFLOW | Menu items use server-side navigation |
| Workflow CALL NANOFLOW | Workflow activities are server-side |
| Published REST NANOFLOW handler | REST operations are server-side |

### Future Work (separate PRs)

| Feature | Priority | Notes |
|---------|----------|-------|
| SHOW ACCESS ON NANOFLOW | P2 | Display nanoflow access roles |
| ELK layout for nanoflows | P3 | Visual layout (low priority) |
| Roundtrip tests | P2 | Verify CREATE → DESCRIBE → re-CREATE |
| JavaScriptActionCall in nanoflows | P2 | `call javascript action` syntax |
| SynchronizeAction | P3 | `synchronize` action (offline nanoflows) |
| Web/Native platform mixing check | P3 | CE6051 validation |

## Grammar Changes

### PR #7 — CREATE NANOFLOW

```antlr
createNanoflowStatement
    : NANOFLOW qualifiedName
      LPAREN microflowParameterList? RPAREN
      microflowReturnType?
      microflowOptions?
      BEGIN microflowBody END SEMICOLON? SLASH?
    ;
```

Added to `createStatement` alternatives. Reuses all microflow sub-rules (parameters, return type, options, body).

### PR #8 — CALL NANOFLOW + GRANT/REVOKE

```antlr
callNanoflowStatement
    : (VARIABLE EQUALS)? CALL NANOFLOW qualifiedName
      LPAREN callArgumentList? RPAREN onErrorClause?
    ;

grantNanoflowAccessStatement
    : GRANT EXECUTE ON NANOFLOW qualifiedName TO moduleRoleList
    ;

revokeNanoflowAccessStatement
    : REVOKE EXECUTE ON NANOFLOW qualifiedName FROM moduleRoleList
    ;
```

## Validation Rules

Implemented in `mdl/executor/nanoflow_validation.go`:

1. **Disallowed actions** — Rejects 21 microflow-only action types with descriptive error messages
2. **ErrorEvent forbidden** — Reports `ErrorEvent is not supported in nanoflows`
3. **Binary return type rejected** — Reports `Binary return type is not supported in nanoflows`
4. **Recursive validation** — Checks compound statements (IF/LOOP/WHILE bodies) and error handling blocks
5. **Cross-reference validation** — `validate.go` checks that `call nanoflow Module.Name` targets exist

## Testing

- All existing tests pass (`make build && make test && make lint-go`)
- Registry test updated with all new AST types
- Manual verification: `grep -r 'sdk/mpr' mdl/executor/` confirms no new `sdk/mpr` imports in executor
- Roundtrip tests planned for future PR

## Files Changed

### PR #7 (16 files)
- `mdl/grammar/MDLParser.g4` — `createNanoflowStatement` rule
- `mdl/grammar/` — regenerated ANTLR parser files
- `mdl/ast/ast_microflow.go` — `CreateNanoflowStmt`, `DropNanoflowStmt`
- `mdl/visitor/visitor_microflow.go` — `ExitCreateNanoflowStatement`
- `mdl/visitor/visitor_entity.go` — NANOFLOW branch in `ExitDropStatement`
- `mdl/executor/cmd_nanoflows_create.go` — CREATE handler
- `mdl/executor/cmd_nanoflows_drop.go` — DROP handler
- `mdl/executor/nanoflow_validation.go` — Action/return type validation
- `mdl/executor/executor.go` — Cache fields + helpers
- `mdl/executor/exec_context.go` — `trackCreatedNanoflow`
- `mdl/executor/register_stubs.go` — Handler registration
- `mdl/executor/registry_test.go` — Test update

### PR #8 (21 files)
- `mdl/grammar/MDLParser.g4` — `callNanoflowStatement`, `grantNanoflowAccessStatement`, `revokeNanoflowAccessStatement`
- `mdl/grammar/` — regenerated ANTLR parser files
- `mdl/ast/ast_microflow.go` — `CallNanoflowStmt`
- `mdl/ast/ast_security.go` — `GrantNanoflowAccessStmt`, `RevokeNanoflowAccessStmt`
- `sdk/microflows/microflows.go` — `AllowedModuleRoles` on Nanoflow
- `sdk/microflows/microflows_actions.go` — `NanoflowCallAction`, `NanoflowCall`, `NanoflowCallParameterMapping`
- `sdk/mpr/parser_nanoflow.go` — AllowedModuleRoles parsing
- `sdk/mpr/parser_microflow_actions.go` — `parseNanoflowCallAction`
- `sdk/mpr/parser_microflow.go` — Action type map registration
- `sdk/mpr/writer_microflow.go` — AllowedModuleRoles serialization
- `sdk/mpr/writer_microflow_actions.go` — NanoflowCallAction serialization
- `mdl/visitor/visitor_microflow_actions.go` — `buildCallNanoflowStatement`
- `mdl/visitor/visitor_microflow_statements.go` — Dispatch + annotation
- `mdl/visitor/visitor_security.go` — Grant/revoke visitors
- `mdl/executor/cmd_microflows_builder_calls.go` — `addCallNanoflowAction`
- `mdl/executor/cmd_microflows_builder_graph.go` — Dispatch
- `mdl/executor/cmd_microflows_builder.go` — `lookupNanoflowReturnType`
- `mdl/executor/cmd_microflows_builder_validate.go` — Validation
- `mdl/executor/cmd_security_write.go` — Grant/revoke handlers
- `mdl/executor/nanoflow_validation.go` — CallNanoflowStmt case
- `mdl/executor/validate.go` — Cross-reference validation
- `mdl/executor/validate_microflow.go` — Flow body validation
- `mdl/executor/cmd_diff_mdl.go` — Diff formatting
- `mdl/executor/stmt_summary.go` — Statement summary
- `mdl/executor/register_stubs.go` — Handler registration
- `mdl/executor/registry_test.go` — Test update
- `cmd/mxcli/skills/write-nanoflows.md` — Agentic skill

## Complexity

- PR #7: Medium (16 files, 110 ins / 75 del)
- PR #8: High (21 files, ~600 ins / ~20 del) — touches more layers due to CALL action + BSON + validation
