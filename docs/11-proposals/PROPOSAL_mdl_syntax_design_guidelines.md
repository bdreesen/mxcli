# Proposal: MDL Syntax Design Guidelines

**Status:** Draft
**Date:** 2026-04-01
**Author:** AI-assisted design
**Related:** [PROPOSAL_mdl_syntax_improvements.md](PROPOSAL_mdl_syntax_improvements.md), [PROPOSAL_mdl_syntax_improvements_v2.md](PROPOSAL_mdl_syntax_improvements_v2.md)

---

## Summary

As multiple developers and AI agents contribute MDL syntax, we need shared design principles to keep the language coherent. This proposal defines guardrails for new MDL additions, covering readability, token efficiency, consistency, reviewability, and LLM fitness. It also proposes a companion skill file (`.claude/skills/design-mdl-syntax.md`) that Claude and contributors consult before designing new syntax.

## Problem

MDL has grown organically across 50+ statement types, contributed by different developers and AI agents. Without explicit design guidelines, inconsistencies have crept in:

- **Variable assignment uses three syntaxes**: `DECLARE $X Type = val`, `SET $X = val`, `$X = CREATE ...`
- **Block delimiters vary**: `BEGIN...END LOOP`, `THEN...END IF`, `{ }` for pages
- **Property syntax varies**: `Key: value` in some contexts, `Key = value` in others
- **Keyword verbosity varies**: `CALL MICROFLOW` vs direct function calls, `RETRIEVE ... FROM ... WHERE` vs simpler forms

These are documented in the v1 and v2 syntax improvement proposals. But those proposals focus on *what to change* — they don't establish principles for *how to design new syntax going forward*. Every new contributor (human or AI) faces the same questions: Should I use `BEGIN...END` or `{}`? Should the keyword be `CREATE` or `ADD`? How verbose should property lists be?

## Design Principles

The following principles are ordered by priority. When principles conflict, higher-priority ones win.

### 1. Read Like English, Not Code

**MDL's primary audience is citizen developers, business analysts, and non-software engineers.** Syntax should read as close to natural language as possible, following SQL and BASIC traditions rather than C or Go.

| Principle | Good | Bad | Why |
|-----------|------|-----|-----|
| Keywords over symbols | `FROM`, `WHERE`, `IN` | `->`, `=>`, `\|>` | Symbols require learning; words are self-documenting |
| Spell out intent | `RETRIEVE $Order FROM Shop.Order WHERE Status = 'Open'` | `$Order = Shop.Order.find({status: 'Open'})` | The keyword version reads as a sentence |
| Avoid abbreviations | `MICROFLOW`, `ENUMERATION`, `ASSOCIATION` | `MF`, `ENUM`, `ASSOC` | Full words reduce ambiguity; tokens are cheap |
| Use prepositions | `GRANT READ ON Entity TO Role` | `GRANT READ Entity Role` | Prepositions clarify the relationship between arguments |

**Test**: Read the statement aloud. If a business analyst would understand it on first hearing, it passes.

### 2. One Way to Do Each Thing

Every concept should have exactly one syntax. When a new feature overlaps with an existing pattern, reuse the existing pattern rather than inventing a new one.

| Pattern | Use for | Example |
|---------|---------|---------|
| `CREATE <TYPE> <QualifiedName> (...)` | New elements | `CREATE ENTITY`, `CREATE MICROFLOW`, `CREATE PAGE` |
| `ALTER <TYPE> <QualifiedName> <operation>` | Modify existing elements | `ALTER ENTITY ... ADD`, `ALTER PAGE ... SET` |
| `DROP <TYPE> <QualifiedName>` | Remove elements | `DROP ENTITY`, `DROP MICROFLOW` |
| `SHOW <TYPE>S [IN Module]` | List elements | `SHOW ENTITIES`, `SHOW MICROFLOWS` |
| `DESCRIBE <TYPE> <QualifiedName>` | Inspect one element | `DESCRIBE ENTITY`, `DESCRIBE MICROFLOW` |
| `GRANT/REVOKE <permission> ON <target> TO/FROM <role>` | Security | All access control |

**When adding a new document type**, follow the existing CRUD pattern: `CREATE`, `ALTER`, `DROP`, `SHOW`, `DESCRIBE`. Do not invent alternative verbs (e.g., `ADD`, `REMOVE`, `LIST`, `VIEW`).

### 3. Optimize for LLM Generation and Comprehension

MDL is increasingly generated and consumed by LLMs. Syntax choices should account for how language models tokenize and reason about code.

**Token efficiency:**
- Prefer shorter keywords when equally readable: `IN` over `CONTAINED_IN`, `OR MODIFY` over `OR_MODIFY_IF_EXISTS`
- Avoid deeply nested structures — LLMs handle flat statement sequences better than trees
- Keep statements self-contained: each statement should have full context, no implicit state from prior statements

**Predictable patterns:**
- LLMs generate more accurate code when patterns are regular. Irregular exceptions (e.g., `BEGIN...END` for microflows but `{}` for pages) cause errors
- Use the same keyword order across statement types: `<VERB> <TYPE> <NAME> <MODIFIERS> <BODY>`
- Property lists should follow a consistent format regardless of context

**Unambiguous parsing:**
- Avoid context-dependent keywords that mean different things in different positions
- Prefer explicit terminators (`END`, `;`) over implicit block boundaries
- Identifiers that collide with keywords should always be quotable with the same mechanism (double-quotes or backticks)

**Test**: Can an LLM generate the syntax correctly from a single example? If it needs 3+ examples to get the pattern right, the syntax is too irregular.

### 4. Make Diffs Reviewable

MDL scripts are reviewed by humans in pull requests and `mxcli diff` output. Syntax should produce small, meaningful diffs.

- **One property per line** in multi-property constructs — adding a property should be a one-line diff
- **Trailing commas allowed** — adding the last item shouldn't modify the previous line
- **Stable ordering** — `DESCRIBE` output should use a deterministic property order so re-running it doesn't produce false diffs
- **No redundant defaults** — `DESCRIBE` should omit properties set to their default values unless the value is non-obvious

```mdl
-- Good: one property per line, adding Width is a one-line diff
CREATE PERSISTENT ENTITY Shop.Product (
    Name: String(200),
    Price: Decimal,
    Description: String(unlimited),
);

-- Bad: all on one line, any change touches the entire statement
CREATE PERSISTENT ENTITY Shop.Product (Name: String(200), Price: Decimal, Description: String(unlimited));
```

### 5. Token Efficiency Without Sacrificing Clarity

Conciseness matters for LLM context windows and human scanning, but never at the expense of readability (Principle 1).

**Do:**
- Omit noise keywords that add no information: `CREATE ENTITY` not `CREATE A NEW ENTITY`
- Allow shorthand for common patterns: `String(200)` not `String WITH LENGTH 200`
- Support `OR MODIFY` to avoid check-then-create sequences
- Use type inference where unambiguous: `DECLARE $Count = 0` (obviously Integer)

**Don't:**
- Use single-character operators for domain operations: `+>` for "add to list"
- Omit keywords that clarify intent: `FROM`, `WHERE`, `TO` are cheap and essential
- Create aliases (two keywords for the same thing): if the verb is `CREATE`, don't also accept `ADD` or `NEW`

### 6. Consistency Across Document Types

The same concept should use the same syntax regardless of where it appears.

| Concept | Consistent syntax | Not this |
|---------|------------------|----------|
| Qualified names | `Module.Element` everywhere | `Module::Element` or `Module/Element` in some contexts |
| Property assignment | `Key: value` in definitions | `Key = value` in some places, `Key: value` in others |
| Boolean properties | `VISIBLE`, `EDITABLE`, `REQUIRED` | `Visible: true`, `IsVisible`, `VISIBLE = YES` |
| Optional clauses | `[clause]` is always omittable | Some optional clauses that error when omitted |
| Block bodies | Consistent delimiter per context | Mixing `BEGIN...END` and `{}` in the same context |

**Current exception**: Microflow bodies use `BEGIN...END` while page bodies use `{}`. This is a known inconsistency. New features should follow whichever convention their parent context uses — microflow actions use `BEGIN...END`, widget definitions use `{}`.

## Decision Framework for New Syntax

When designing syntax for a new MDL feature, answer these questions in order:

### Step 1: Does an Existing Pattern Cover This?

Check the MDL Quick Reference (`docs/01-project/MDL_QUICK_REFERENCE.md`). If an existing statement type covers the concept, extend it rather than creating new syntax.

```
New concept: "scheduled events"
→ Existing pattern: CREATE/ALTER/DROP/SHOW/DESCRIBE
→ Design: CREATE SCHEDULED EVENT Module.Name (...)
           DESCRIBE SCHEDULED EVENT Module.Name
           SHOW SCHEDULED EVENTS [IN Module]
```

### Step 2: What's the Statement Shape?

All MDL statements follow one of these shapes:

```
DDL:   <VERB> [MODIFIERS] <TYPE> <QualifiedName> [CLAUSES] [BODY];
DML:   <ACTION> <TARGET> [CLAUSES];
DQL:   <QUERY-VERB> <TYPE>S [FILTERS];
```

New statements must fit one of these shapes. If your feature doesn't fit, reconsider whether it belongs in MDL or should be a CLI command instead.

### Step 3: Keyword Selection

1. **Reuse existing keywords** before inventing new ones. Check the reserved words list in the grammar.
2. **Use standard SQL/DDL verbs**: CREATE, ALTER, DROP, SHOW, DESCRIBE, GRANT, REVOKE, SET
3. **Use Mendix terminology** for domain concepts: ENTITY, MICROFLOW, PAGE, ASSOCIATION (not TABLE, FUNCTION, VIEW, RELATION)
4. **Prepositions clarify structure**: FROM, TO, IN, ON, BY, WITH, AS, WHERE, INTO

### Step 4: Property Lists

Use this format for all property-bearing constructs:

```mdl
CREATE <TYPE> Module.Name (
    Property1: value,
    Property2: value,
    Property3: value,
);
```

Rules:
- Parentheses `()` delimit property lists
- Colon `:` separates key from value (not `=`)
- Comma `,` separates properties
- Trailing comma allowed
- Properties on separate lines for readability (single line acceptable for 1-2 properties)
- Default values omitted unless non-obvious

### Step 5: Read It Aloud

Read the proposed syntax as an English sentence. Verify:
- A business analyst understands the intent
- No ambiguous interpretations exist
- The statement is self-contained (doesn't depend on implicit context)

### Step 6: Test LLM Generation

Give an LLM one example of the new syntax and ask it to generate a variant. If the LLM consistently gets it wrong, the pattern is too irregular or too different from established MDL patterns.

### Step 7: Check Diff Impact

Write two versions of a statement (before and after a small change) and verify the diff is minimal and readable.

## Anti-Patterns

These are patterns to **avoid** in new MDL syntax.

### Overloaded Keywords

```mdl
-- Bad: SET means different things in different contexts
SET $Variable = value;           -- variable assignment
ALTER PAGE ... SET Caption = ''; -- property modification
ALTER SETTINGS SET Key = value;  -- settings change
```

When a keyword has an established meaning in one context, avoid repurposing it with a different meaning elsewhere. (The `SET` overload above is a known debt, not a pattern to extend.)

### Implicit Context

```mdl
-- Bad: what does CONNECT refer to? Where is Module set?
USE MODULE Shop;
CREATE ENTITY Customer (...);   -- implicitly Shop.Customer?

-- Good: explicit qualified name, no implicit state
CREATE ENTITY Shop.Customer (...);
```

Every statement should be independently understandable. Implicit module context (like SQL's `USE DATABASE`) saves a few tokens but makes scripts fragile and hard to review in diffs.

### Symbolic Soup

```mdl
-- Bad: requires learning symbol meanings
$items |> filter($.active) |> map($.name) |> join(",")

-- Good: reads as English
FILTER $Items WHERE Active = true
```

Symbols are powerful for experienced programmers but hostile to MDL's target audience. Prefer keyword-based syntax.

### Feature-Specific Verbs

```mdl
-- Bad: unique verbs for each feature
SCHEDULE EVENT Module.Name ...
REGISTER WEBHOOK Module.Name ...
DEPLOY SERVICE Module.Name ...

-- Good: consistent CREATE pattern
CREATE SCHEDULED EVENT Module.Name (...)
CREATE WEBHOOK Module.Name (...)
CREATE PUBLISHED SERVICE Module.Name (...)
```

### Magic Strings and Positional Arguments

```mdl
-- Bad: position-dependent, meaning unclear
CREATE RULE Shop Process Order ACT_ProcessOrder

-- Good: labeled, self-documenting
CREATE RULE Shop.ProcessOrder (
    Type: Validation,
    Microflow: Shop.ACT_ProcessOrder,
);
```

## Checklist for New Syntax

Before merging any PR that adds new MDL syntax:

- [ ] **Follows CREATE/ALTER/DROP/SHOW/DESCRIBE pattern** — no custom verbs for standard CRUD operations
- [ ] **Uses `Module.Element` qualified names** — no bare names, no alternative separators
- [ ] **Property lists use `( Key: value, ... )` format** — consistent delimiters and separators
- [ ] **Keywords are full English words** — no abbreviations, no symbols for domain operations
- [ ] **Statement reads as an English sentence** — a business analyst can understand the intent
- [ ] **One example is sufficient for LLM generation** — tested by giving one example and asking for a variant
- [ ] **Diff is minimal for small changes** — adding one property is a one-line diff
- [ ] **No new keyword overloading** — each keyword means one thing
- [ ] **No implicit context** — every statement is self-contained with qualified names
- [ ] **DESCRIBE output roundtrips** — `DESCRIBE` produces valid MDL that can be re-executed
- [ ] **Grammar updated and regenerated** — `make grammar` runs clean
- [ ] **Quick reference updated** — `docs/01-project/MDL_QUICK_REFERENCE.md` has the new syntax
- [ ] **Skill file consulted** — developer read `.claude/skills/design-mdl-syntax.md` before designing

## Examples: Applying the Guidelines

### Example 1: Adding Workflow Support

```mdl
-- Follows CREATE pattern, uses Mendix terminology, reads as English
CREATE WORKFLOW Shop.ApproveOrder (
    Description: 'Order approval workflow',
    Parameter: $Order Shop.Order,
)
BEGIN
    USER TASK ReviewOrder (
        Assignee: Shop.Manager,
        Page: Shop.OrderReview_Task,
        Description: 'Review the order details',
    );
    DECISION IsApproved (
        Caption: 'Approved?',
    )
        WHEN $Order/Approved = true THEN
            CALL MICROFLOW Shop.ACT_FulfillOrder(Order: $Order);
        WHEN $Order/Approved = false THEN
            CALL MICROFLOW Shop.ACT_RejectOrder(Order: $Order);
    END DECISION;
END;
```

### Example 2: Adding Scheduled Event Support

```mdl
-- Standard pattern: CREATE + DESCRIBE + SHOW + DROP
CREATE SCHEDULED EVENT Shop.DailyCleanup (
    Microflow: Shop.ACT_Cleanup,
    Interval: 'Daily',
    StartTime: '02:00',
    Enabled: true,
);

SHOW SCHEDULED EVENTS IN Shop;
DESCRIBE SCHEDULED EVENT Shop.DailyCleanup;
DROP SCHEDULED EVENT Shop.DailyCleanup;
```

### Example 3: Wrong Way (anti-patterns)

```mdl
-- Anti-pattern 1: Custom verb instead of CREATE
SCHEDULE Shop.DailyCleanup EVERY DAY AT '02:00' RUN Shop.ACT_Cleanup;

-- Anti-pattern 2: Implicit module context
USE MODULE Shop;
SCHEDULE DailyCleanup ...;

-- Anti-pattern 3: Symbolic syntax
Shop::DailyCleanup => Shop::ACT_Cleanup @ "0 2 * * *"

-- Anti-pattern 4: Positional arguments
CREATE SCHEDULED EVENT Shop DailyCleanup Shop.ACT_Cleanup Daily 02:00 true
```
