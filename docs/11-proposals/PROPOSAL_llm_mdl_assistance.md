# Proposal: Improving LLM Assistance for MDL Code Generation

## Problem Statement

MDL (Mendix Definition Language) is a custom DSL that does not exist in LLM training data. When Claude Code or other LLMs attempt to write MDL, they rely entirely on:

- Grammar files (MDL.g4) - difficult to interpret programmatically
- Skill files - primary teaching mechanism but incomplete
- Example files - learning by pattern matching
- Error messages - learning what went wrong (reactive, not proactive)

This results in common mistakes like:
- Using `SET` on undeclared variables
- Wrong syntax for entity type declarations (`Type = empty` vs `AS Type`)
- Missing qualifications on association paths
- Incorrect enumeration comparisons (string literals vs qualified values)

## Goals

1. **Reduce iteration cycles** - LLM writes correct code on first attempt
2. **Self-correcting errors** - When errors occur, provide enough context to fix them
3. **Pattern consistency** - Encourage best practices and standard patterns
4. **Discoverability** - Make MDL capabilities easy to find and understand

## Proposed Improvements

### 1. Enhanced Error Messages with Examples

**Priority: High | Effort: Low**

Current error messages tell what's wrong but not how to fix it:

```
variable 'IsValid' is not declared. Use DECLARE IsValid: <Type> before using SET
```

Proposed format with inline example:

```
variable 'IsValid' is not declared.

Fix: Add a DECLARE statement before using SET:

  DECLARE $IsValid Boolean = true;
  ...
  SET $IsValid = false;
```

**Implementation:**

```go
// In cmd_microflows_builder.go
func (fb *flowBuilder) addErrorWithExample(message, example string) {
    fb.errors = append(fb.errors, fmt.Sprintf("%s\n\nExample:\n%s", message, example))
}

// Usage
fb.addErrorWithExample(
    fmt.Sprintf("variable '%s' is not declared", s.Target),
    fmt.Sprintf("  DECLARE %s Boolean = true;\n  SET %s = false;", s.Target, s.Target),
)
```

**Error categories to enhance:**

| Error | Current Message | Proposed Addition |
|-------|-----------------|-------------------|
| Undeclared variable | "variable X not declared" | Show DECLARE + SET pattern |
| Entity type syntax | "selected type not allowed" | Show `DECLARE $var AS Module.Entity` |
| Association path | "error in expression" | Show `$var/Module.Association/Attr` |
| Enum comparison | "type mismatch" | Show `Module.Enum.Value` syntax |

### 2. Focused Skills by Document Type / Use Case

**Priority: High | Effort: Medium**

Instead of one large reference, create smaller focused skills that can be loaded individually based on the task at hand. This keeps context focused and reduces token usage.

**Proposed skill organization:**

```
.claude/skills/mendix/
├── README.md                      # Index of all skills
│
├── # By Document Type (syntax reference)
├── mdl-entities.md                # Entity, attributes, associations
├── mdl-enumerations.md            # Enumeration syntax
├── mdl-microflows.md              # Microflow syntax (exists: write-microflows.md)
├── mdl-pages.md                   # Page and widget syntax
│
├── # By Use Case (patterns)
├── patterns-validation.md         # Validation patterns (exists: validation-microflows.md)
├── patterns-crud.md               # Create/Read/Update/Delete patterns
├── patterns-data-processing.md    # Loops, aggregates, batch processing
├── patterns-integration.md        # REST, Java actions, external calls
│
├── # Quick References (cheat sheets)
├── cheatsheet-variables.md        # Variable declaration quick ref
├── cheatsheet-expressions.md      # Operators, functions, XPath
├── cheatsheet-errors.md           # Common errors and fixes
│
└── # Existing skills
    ├── write-microflows.md
    ├── validation-microflows.md
    ├── write-oql-queries.md
    └── ...
```

**Skill loading strategy:**

1. **Task-based loading** - Load relevant skill based on user request:
   - "Create validation microflow" → load `patterns-validation.md`
   - "Add entity with attributes" → load `mdl-entities.md`

2. **Error-based loading** - When errors occur, suggest relevant skill:
   - Variable declaration error → reference `cheatsheet-variables.md`
   - XPath error → reference `cheatsheet-expressions.md`

3. **Keep skills small** - Target 100-200 lines per skill file

**Example: cheatsheet-variables.md (focused, ~50 lines)**

```markdown
# MDL Variable Cheatsheet

## Declaration Syntax

| Type | Syntax | Example |
|------|--------|---------|
| String | `DECLARE $name String = 'value';` | `DECLARE $msg String = '';` |
| Integer | `DECLARE $name Integer = 0;` | `DECLARE $count Integer = 0;` |
| Boolean | `DECLARE $name Boolean = true;` | `DECLARE $valid Boolean = true;` |
| Decimal | `DECLARE $name Decimal = 0.0;` | `DECLARE $amount Decimal = 0;` |
| DateTime | `DECLARE $name DateTime = [%CurrentDateTime%];` | |
| Entity | `DECLARE $name AS Module.Entity;` | `DECLARE $cust AS Sales.Customer;` |
| List | `DECLARE $name List of Module.Entity = empty;` | |

## Key Rules

1. **Primitives**: `DECLARE $var Type = value;` (with initialization)
2. **Entities**: `DECLARE $var AS Module.Entity;` (no initialization, use AS)
3. **SET requires DECLARE**: Always declare before using SET
4. **Parameters are pre-declared**: No need to declare microflow parameters

## Common Mistakes

❌ `DECLARE $product Module.Product = empty;` → Missing AS
✅ `DECLARE $product AS Module.Product;`

❌ `SET $isValid = true;` (without prior DECLARE)
✅ `DECLARE $isValid Boolean = true;` then `SET $isValid = false;`
```

**Example: patterns-crud.md (focused, ~150 lines)**

```markdown
# CRUD Action Patterns

Patterns for Create, Read, Update, Delete operations on entities.

## ACT_Entity_Save (Create or Update)

```mdl
/**
 * Save action for Entity NewEdit page
 * Validates, commits, and closes the page
 */
CREATE MICROFLOW Module.ACT_Entity_Save (
  $Entity: Module.Entity
)
RETURNS Boolean
BEGIN
  -- Validate first
  $IsValid = CALL MICROFLOW Module.VAL_Entity_Save($Entity = $Entity);

  IF $IsValid THEN
    COMMIT $Entity;
    CLOSE PAGE;
  END IF;

  RETURN $IsValid;
END;
/
```

## ACT_Entity_Delete

```mdl
/**
 * Delete action with confirmation
 */
CREATE MICROFLOW Module.ACT_Entity_Delete (
  $Entity: Module.Entity
)
RETURNS Boolean
BEGIN
  DELETE $Entity;
  CLOSE PAGE;
  RETURN true;
END;
/
```

## When to Use

- **ACT_Entity_Save**: Save button on NewEdit pages
- **ACT_Entity_Delete**: Delete button with confirmation dialog
- **ACT_Entity_Cancel**: Cancel button (just close page, no commit)
```

### 4. Check Command with Suggestions

**Priority: Medium | Effort: Medium**

Add `--suggest` flag to provide fix suggestions:

```bash
$ mxcli check script.mdl -p app.mpr --references --suggest

Checking: script.mdl
✓ Syntax OK (3 statements)

Reference errors:
  statement 2: microflow 'Module.Test' has validation errors:
    - variable 'IsValid' is not declared

    Suggested fix (line 7):
    + DECLARE $IsValid Boolean = true;
      IF $Entity/Name = empty THEN
        SET $IsValid = false;  -- line 9

✗ 1 error(s) found
```

**Implementation approach:**

1. Track source locations in AST nodes
2. Generate diff-style suggestions
3. Optionally apply fixes with `--fix` flag

### 5. DESCRIBE Enhancement for Learning

**Priority: Low | Effort: Low**

Add comments to DESCRIBE output explaining syntax:

```bash
$ mxcli -p app.mpr -c "DESCRIBE MICROFLOW Module.Example --annotated"

-- Microflow signature: Name, parameters, return type
CREATE MICROFLOW Module.Example (
  $Input: String           -- Parameter: $name: Type
)
RETURNS Boolean AS $Result -- RETURNS Type AS $variableName
BEGIN
  -- Variable declaration: DECLARE $name Type = value
  DECLARE $Result Boolean = true;

  -- Conditional: IF condition THEN ... END IF
  IF $Input = empty THEN
    SET $Result = false;   -- Assignment: SET $var = expression
  END IF;

  RETURN $Result;          -- Must end with RETURN
END;
/
```

### 6. Lint Rules as Teaching Tools

**Priority: Medium | Effort: Low**

Enhance lint rule output with educational content:

```python
# In lint rules
def check_undeclared_variable(node, context):
    if is_set_statement(node) and not is_declared(node.variable, context):
        return {
            "rule": "MDL020",
            "severity": "error",
            "message": f"Variable '{node.variable}' used in SET but not declared",
            "learn_more": "https://docs.example.com/mdl/variables",
            "quick_fix": {
                "description": "Add DECLARE statement",
                "insert_before": node.line,
                "text": f"DECLARE {node.variable} Boolean = true; -- TODO: set correct type"
            }
        }
```

### 7. Interactive Examples in REPL

**Priority: Low | Effort: Medium**

Add `EXAMPLE` command to REPL:

```
mdl> EXAMPLE validation
-- Validation Microflow Pattern
CREATE MICROFLOW Module.VAL_Entity_Action (
  $Entity: Module.Entity
)
RETURNS Boolean AS $IsValid
BEGIN
  DECLARE $IsValid Boolean = true;

  IF $Entity/RequiredField = empty THEN
    SET $IsValid = false;
    VALIDATION FEEDBACK $Entity/RequiredField MESSAGE 'Required';
  END IF;

  RETURN $IsValid;
END;
/

mdl> EXAMPLE loop
-- Loop Pattern
...
```

## Implementation Roadmap

### Phase 1: Quick Wins (1-2 days)
- [ ] Enhanced error messages with examples
- [ ] Create `cheatsheet-variables.md` skill (~50 lines)
- [ ] Create `cheatsheet-errors.md` skill (~50 lines)
- [ ] Update existing skills with more examples

### Phase 2: Focused Skills (3-5 days)
- [ ] Create `patterns-crud.md` skill
- [ ] Create `patterns-data-processing.md` skill
- [ ] Create `mdl-entities.md` syntax reference
- [ ] Update README.md with skill index and loading guidance

### Phase 3: Tooling (1 week)
- [ ] `--suggest` flag for check command
- [ ] Lint rules with educational output
- [ ] EXAMPLE command in REPL

### Phase 4: Advanced (future)
- [ ] Source location tracking in AST
- [ ] Auto-fix capability (`--fix` flag)
- [ ] Interactive tutorial mode

## Success Metrics

1. **First-attempt success rate** - % of LLM-generated MDL that passes check
2. **Iteration count** - Average attempts needed to get valid MDL
3. **Common error reduction** - Track top 10 errors, measure reduction
4. **User feedback** - Qualitative feedback on error message helpfulness

## Appendix: Common LLM Mistakes

Based on observed patterns:

| Mistake | Frequency | Root Cause |
|---------|-----------|------------|
| SET without DECLARE | High | No equivalent in most languages |
| Entity decl syntax | High | Unusual `AS` keyword requirement |
| String enum comparison | Medium | Most languages use strings |
| Missing association qualification | Medium | XPath-style paths unfamiliar |
| Wrong DECLARE syntax (colon) | Medium | Confusion with TypeScript/Python |
| Missing RETURN | Low | Different from void functions |

## References

- [MDL Grammar](../../mdl/grammar/MDL.g4)
- [Existing Skills](../../.claude/skills/mendix/)
- [Example Files](../../mdl-examples/)
- [Write Microflows Skill](../../.claude/skills/mendix/write-microflows.md)
