# Fix Issue Skill

A fast-path workflow for diagnosing and fixing bugs in mxcli. Each fix appends
to the symptom table below, so the next similar issue costs fewer reads.

## How to Use

1. Match the issue symptom to a row in the table — go straight to that file.
2. Follow the fix pattern for that row.
3. Write a failing test first, then implement.
4. After the fix: **add a new row** to the table if the symptom is not already covered.

---

## Symptom → Layer → File Table

| Symptom | Root cause layer | First file to open | Fix pattern |
|---------|-----------------|-------------------|-------------|
| `DESCRIBE` shows `$var = LIST OPERATION ...;` | Missing parser case | `sdk/mpr/parser_microflow.go` → `parseListOperation()` | Add `case "Microflows$XxxType":` returning the correct struct |
| `DESCRIBE` shows `$var = ACTION ...;` | Missing formatter case | `mdl/executor/cmd_microflows_format_action.go` → `formatActionStatement()` | Add `case *microflows.XxxAction:` with `fmt.Sprintf` output |
| `DESCRIBE` shows `$var = LIST OPERATION %T;` (with type name) | Missing formatter case | `mdl/executor/cmd_microflows_format_action.go` → `formatListOperation()` | Add `case *microflows.XxxOperation:` before the `default` |
| Compile error: `undefined: microflows.XxxOperation` | Missing SDK struct | `sdk/microflows/microflows_actions.go` | Add struct + `func (XxxOperation) isListOperation() {}` marker |
| `TypeCacheUnknownTypeException` in Studio Pro | Wrong `$Type` storage name in BSON write | `sdk/mpr/writer_microflow.go` | Check the storage name table in CLAUDE.md; verify against `reference/mendixmodellib/reflection-data/` |
| CE0066 "Entity access is out of date" | MemberAccess added to wrong entity | `sdk/mpr/writer_domainmodel.go` | MemberAccess must only be on the FROM entity (`ParentPointer`), not the TO entity — see CLAUDE.md association semantics |
| CE0463 "widget definition changed" | Object property structure doesn't match Type PropertyTypes | `sdk/widgets/templates/` | Re-extract template from Studio Pro; see `sdk/widgets/templates/README.md` |
| Parser returns `nil` for a known BSON type | Unhandled `default` in a `parseXxx()` switch | `sdk/mpr/parser_microflow.go` or `parser_page.go` | Find the switch by grepping for `default: return nil`; add the missing case |
| MDL check gives "unexpected token" on valid-looking syntax | Grammar missing rule or token | `mdl/grammar/MDLParser.g4` + `MDLLexer.g4` | Add rule/token, run `make grammar` |

---

## TDD Protocol

Always follow this order — never implement before the test exists:

```
Step 1: Write a failing unit test (parser test or formatter test)
Step 2: Confirm it fails to compile or fails at runtime
Step 3: Implement the minimum code to make it pass
Step 4: Run: /c/Users/Ylber.Sadiku/go/go/bin/go test ./mdl/executor/... ./sdk/mpr/...
Step 5: Add the symptom row to the table above if not already present
```

Parser tests go in `sdk/mpr/parser_<domain>_test.go`.
Formatter tests go in `mdl/executor/cmd_<domain>_format_<area>_test.go`.

---

## Issue #212 — Reference Fix (seeding example)

**Symptom:** `DESCRIBE MICROFLOW` showed `$var = LIST OPERATION ...;` for
`Microflows$Find`, `Microflows$Filter`, `Microflows$ListRange`.

**Root cause:** `parseListOperation()` in `sdk/mpr/parser_microflow.go` had no
cases for these three BSON types — they fell to `default: return nil`.

**Files changed:**
| File | Change |
|------|--------|
| `sdk/microflows/microflows_actions.go` | Added `FindByAttributeOperation`, `FilterByAttributeOperation`, `RangeOperation` |
| `sdk/mpr/parser_microflow.go` | Added 3 parser cases |
| `mdl/executor/cmd_microflows_format_action.go` | Added 3 formatter cases |
| `mdl/executor/cmd_microflows_format_listop_test.go` | Added 4 formatter tests |
| `sdk/mpr/parser_listoperation_test.go` | New file, 4 parser tests |

**Key insight:** `Microflows$ListRange` stores offset/limit inside a nested
`CustomRange` map — must cast `raw["CustomRange"].(map[string]any)` before
extracting `OffsetExpression`/`LimitExpression`.

---

## After Every Fix — Checklist

- [ ] Failing test written before implementation
- [ ] `go test ./mdl/executor/... ./sdk/mpr/...` passes
- [ ] New symptom row added to the table above (if not already covered)
- [ ] PR title: `fix: <one-line description matching the symptom>`
