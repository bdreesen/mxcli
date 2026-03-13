// ============================================================================
// Professional Document Template
// ============================================================================

// Page setup
#set page(
  paper: "a4",
  margin: (x: 2.5cm, y: 2.5cm),
  header: context {
    if counter(page).get().first() > 1 [
      #set text(size: 9pt, fill: gray)
      #h(1fr)
      _Mendix for Agentic IDEs_
    ]
  },
  footer: context {
    if counter(page).get().first() > 1 [
      #set text(size: 9pt, fill: gray)
      #h(1fr)
      #counter(page).display("1")
    ]
  }
)

// Typography - Sans-serif fonts (with cross-platform fallbacks)
#set text(
  font: ("Helvetica Neue", "Helvetica", "Arial"),
  size: 10.5pt,
  hyphenate: true
)

// Headings
#set heading(numbering: "1.1.1")
#show heading.where(level: 1): it => {
  pagebreak(weak: true)
  v(1.5em)
  text(size: 18pt, weight: "bold", fill: rgb("#1a365d"), it)
  v(1em)
}
#show heading.where(level: 2): it => {
  v(1.2em)
  text(size: 14pt, weight: "bold", fill: rgb("#2c5282"), it)
  v(0.6em)
}
#show heading.where(level: 3): it => {
  v(0.8em)
  text(size: 12pt, weight: "bold", fill: rgb("#2d3748"), it)
  v(0.4em)
}
#show heading.where(level: 4): it => {
  v(0.6em)
  text(size: 11pt, weight: "bold", style: "italic", it)
  v(0.3em)
}

// Code blocks - monospace with nice styling
// Detect ASCII art diagrams by looking for box-drawing characters
#show raw.where(block: true): it => {
  let content = it.text
  let is_diagram = content.contains("┌") or content.contains("└") or content.contains("│") or content.contains("─") or content.contains("╔") or content.contains("║") or content.contains("├") or content.contains("┬")

  if is_diagram {
    // Diagram styling - centered, no background tint, larger font
    set text(font: ("Menlo", "Monaco", "Courier New", "Courier"), size: 7pt)
    block(
      fill: white,
      stroke: 1pt + rgb("#e2e8f0"),
      inset: 15pt,
      radius: 6pt,
      width: 100%,
      align(center, it)
    )
  } else {
    // Regular code block styling
    set text(font: ("Menlo", "Monaco", "Courier New", "Courier"), size: 8.5pt)
    block(
      fill: rgb("#f7fafc"),
      stroke: (left: 3pt + rgb("#4299e1")),
      inset: (left: 12pt, right: 10pt, top: 10pt, bottom: 10pt),
      radius: (right: 4pt),
      width: 100%,
      it
    )
  }
}

// Inline code
#show raw.where(block: false): it => {
  set text(font: ("Menlo", "Monaco", "Courier New", "Courier"), size: 9.5pt)
  box(fill: rgb("#edf2f7"), inset: (x: 4pt, y: 2pt), radius: 3pt, it)
}

// Tables
#set table(
  stroke: (x, y) => (
    top: if y == 0 { 1.5pt + rgb("#2d3748") } else { 0.5pt + rgb("#e2e8f0") },
    bottom: if y == 0 { 1pt + rgb("#4a5568") } else { 0.5pt + rgb("#e2e8f0") },
    left: 0pt,
    right: 0pt,
  ),
  inset: 8pt,
)
#show table.cell.where(y: 0): strong

// Links
#show link: it => {
  set text(fill: rgb("#3182ce"))
  underline(it)
}

// Block quotes
#show quote: it => {
  block(
    fill: rgb("#ebf8ff"),
    stroke: (left: 4pt + rgb("#4299e1")),
    inset: (left: 16pt, right: 12pt, top: 12pt, bottom: 12pt),
    radius: (right: 4pt),
    it
  )
}

// Horizontal rule
#let horizontalrule = {
  v(1em)
  line(length: 100%, stroke: 1pt + rgb("#e2e8f0"))
  v(1em)
}

// ============================================================================
// Title Page
// ============================================================================

#page(header: none, footer: none)[
  #set align(center)

  #v(3cm)

  // Logo placeholder (optional)
  #block(
    width: 100%,
    height: 2cm,
    // Add logo here if available
  )

  #v(2cm)

  // Title
  #text(size: 32pt, weight: "bold", fill: rgb("#1a365d"))[
    Mendix for Agentic IDEs
  ]

  #v(0.5cm)

  #text(size: 18pt, fill: rgb("#4a5568"))[
    Vision & Architecture
  ]

  #v(3cm)

  // Metadata box
  #block(
    width: 70%,
    stroke: 1pt + rgb("#e2e8f0"),
    radius: 8pt,
    inset: 20pt,
  )[
    #set align(left)
    #set text(size: 11pt)

    #grid(
      columns: (1fr, 2fr),
      row-gutter: 12pt,
      [*Status:*], [Vision Document],
      [*Updated:*], [February 2026],
      [*Audience:*], [Product Strategy, Architecture, Engineering Leadership],
    )
  ]

  #v(1fr)

  // Footer
  #set text(size: 10pt, fill: rgb("#718096"))
  _This document outlines how Mendix can become the preferred target for agentic code generation of business applications._

  #v(2cm)
]

// ============================================================================
// Table of Contents
// ============================================================================

#page(header: none)[
  #heading(outlined: false, numbering: none)[Table of Contents]
  #v(1em)
  #outline(
    title: none,
    indent: 1.5em,
    depth: 4,
  )
]

// ============================================================================
// Document Content
// ============================================================================

<mendix-for-agentic-ides-vision-architecture>
#quote(block: true)[
#strong[Status];: Vision Document (Updated February 2026) #strong[Audience];: Product Strategy, Architecture, Engineering Leadership
]

= Executive Summary
<executive-summary>
The rise of AI-powered coding assistants (Claude Code, GitHub Copilot, Cursor, Windsurf) is transforming software development. These "agentic IDEs" can generate, modify, and debug code autonomously. However, they face significant challenges when generating enterprise business applications:

+ #strong[Verbosity];: General-purpose languages require extensive boilerplate
+ #strong[Risk];: AI-generated code may contain security vulnerabilities, bugs, or architectural flaws
+ #strong[Review burden];: Users struggle to validate large volumes of generated code
+ #strong[Governance];: Enterprises need guardrails around AI-generated software

#strong[Mendix is uniquely positioned to address these challenges] through:

- #strong[MDL (Mendix Definition Language)];: A concise DSL that is 5-10x more token-efficient than equivalent TypeScript/Java
- #strong[Platform Guarantees];: Built-in security, scalability, and compliance from the Mendix runtime
- #strong[Model Validation];: Comprehensive checks that catch errors before deployment
- #strong[Visual Review];: Generated applications can be reviewed in Studio Pro's visual interface

== Strategic Differentiators
<strategic-differentiators>
#figure(
  align(center)[#table(
    columns: (45.71%, 54.29%),
    align: (auto,auto,),
    table.header([Differentiator], [Value Proposition],),
    table.hline(),
    [#strong[Open Platform];], [First low-code platform with full AI agent integration---no proprietary lock-in],
    [#strong[Bring Your Own Agent];], [Works with Claude Code, Copilot, Cursor, or any future AI agent],
    [#strong[Human + AI Collaboration];], [Built for true collaboration---agents generate, humans review and refine],
    [#strong[Hours, Not Weeks];], [Complex use cases (migration, bulk updates, monolith decomposition) achievable in hours],
    [#strong[Co-existence Strategy];], [Best of both worlds---visual Studio Pro for design, CLI/agents for automation],
  )]
  , kind: table
  )

This document outlines how Mendix can become the preferred target for agentic code generation of business applications.

== Current Implementation Highlights
<current-implementation-highlights>
The following components have been implemented:

- #strong[mxcli];: Unified CLI with `exec`, `check`, `lint`, `diff`, `search`, `init` commands
- #strong[ANTLR4 Parser];: Cross-language grammar for MDL syntax
- #strong[CATALOG System];: SQL-based metadata queries with REFS table and code search commands
- #strong[Linting Framework];: Extensible rules with SARIF output for CI/CD
- #strong[Claude Code Integration];: `mxcli init` installs skills and commands into Mendix projects
- #strong[Semantic Validation];: Variable scope checking, return value requirements, reference validation

#horizontalrule

= The Opportunity
<the-opportunity>
== Market Context: Agentic IDEs
<market-context-agentic-ides>
A new category of development tools is emerging:

#figure(
  align(center)[#table(
    columns: (19.35%, 41.94%, 38.71%),
    align: (auto,auto,auto,),
    table.header([Tool], [Description], [Capability],),
    table.hline(),
    [#strong[Claude Code];], [Anthropic's CLI agent], [Full codebase understanding, multi-file edits, terminal access],
    [#strong[GitHub Copilot Workspace];], [GitHub's agentic coding], [Issue-to-PR automation, code generation],
    [#strong[Cursor];], [AI-native IDE], [Inline generation, codebase chat, multi-file edits],
    [#strong[Windsurf];], [Codeium's agentic IDE], [Autonomous coding flows, context awareness],
    [#strong[Devin];], [Cognition's AI developer], [Fully autonomous software engineering],
  )]
  , kind: table
  )

These tools are rapidly improving and will soon be capable of generating entire applications from natural language specifications.

== The Problem: AI + Traditional Code = Risk
<the-problem-ai-traditional-code-risk>
When agentic IDEs generate traditional code (TypeScript, Python, Java), enterprises face:

#figure(image("diagrams/mermaid_1.png"),
  caption: [
    Diagram 1
  ]
)

== The Mendix Advantage
<the-mendix-advantage>
Mendix transforms this equation:

#figure(image("diagrams/mermaid_2.png"),
  caption: [
    Diagram 2
  ]
)

#horizontalrule

= MDL: A Token-Efficient DSL for Business Applications
<mdl-a-token-efficient-dsl-for-business-applications>
== What is MDL?
<what-is-mdl>
MDL (Mendix Definition Language) is a textual representation of Mendix models. It provides:

- #strong[Declarative syntax] for entities, microflows, pages, integrations
- #strong[SQL-like familiarity] for developers and AI models
- #strong[Bidirectional mapping] to/from Mendix visual models

== Token Efficiency: MDL vs Traditional Code
<token-efficiency-mdl-vs-traditional-code>
AI models are constrained by context windows and cost per token. MDL is dramatically more efficient:

#strong[Example: Create a Customer entity with CRUD operations]

=== Traditional TypeScript (Prisma + Next.js): \~450 tokens
<traditional-typescript-prisma-next.js-450-tokens>
```typescript
// schema.prisma
model Customer {
  id        Int      @id @default(autoincrement())
  name      String   @db.VarChar(200)
  email     String   @db.VarChar(200)
  balance   Decimal  @default(0)
  isActive  Boolean  @default(true)
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}

// pages/api/customers/index.ts
import { PrismaClient } from '@prisma/client'
const prisma = new PrismaClient()

export default async function handler(req, res) {
  if (req.method === 'GET') {
    const customers = await prisma.customer.findMany()
    return res.json(customers)
  }
  if (req.method === 'POST') {
    const customer = await prisma.customer.create({
      data: req.body
    })
    return res.json(customer)
  }
}

// pages/api/customers/[id].ts
export default async function handler(req, res) {
  const { id } = req.query
  if (req.method === 'GET') {
    const customer = await prisma.customer.findUnique({
      where: { id: Number(id) }
    })
    return res.json(customer)
  }
  if (req.method === 'PUT') {
    const customer = await prisma.customer.update({
      where: { id: Number(id) },
      data: req.body
    })
    return res.json(customer)
  }
  if (req.method === 'DELETE') {
    await prisma.customer.delete({
      where: { id: Number(id) }
    })
    return res.status(204).end()
  }
}

// Plus: validation, error handling, authentication, authorization...
```

=== MDL: \~80 tokens (5-6x more efficient)
<mdl-80-tokens-5-6x-more-efficient>
```sql
CREATE PERSISTENT ENTITY CRM.Customer (
  Name: String(200) NOT NULL,
  Email: String(200),
  Balance: Decimal DEFAULT 0,
  IsActive: Boolean DEFAULT true
);

-- CRUD operations are automatic in Mendix!
-- Security, validation, API all handled by platform
```

== Why Token Efficiency Matters
<why-token-efficiency-matters>
#figure(
  align(center)[#table(
    columns: 2,
    align: (auto,auto,),
    table.header([Factor], [Impact],),
    table.hline(),
    [#strong[Cost];], [Fewer tokens = lower API costs for AI generation],
    [#strong[Speed];], [Smaller context = faster generation],
    [#strong[Accuracy];], [Less code = fewer opportunities for errors],
    [#strong[Review];], [Concise output = easier human validation],
    [#strong[Context];], [More room for specifications and examples],
  )]
  , kind: table
  )

== MDL Coverage
<mdl-coverage>
MDL aims to express complete business applications. The table below reflects the current implementation status:

#figure(
  align(center)[#table(
    columns: (25%, 50%, 25%),
    align: (auto,auto,auto,),
    table.header([Domain], [MDL Constructs], [Status],),
    table.hline(),
    [#strong[Data Model];], [Entity, Association, Enumeration, Index], [✅ Implemented],
    [#strong[Business Logic];], [Microflow (core activities), Nanoflow, Rules], [⚠️ Partial --- microflows have core activities but many activity types missing; nanoflows and rules not yet implemented],
    [#strong[User Interface];], [Page, Snippet, Layout, Widgets (V3 syntax)], [⚠️ Partial --- basic page/widget support; many widgets and widget options missing],
    [#strong[Integration];], [REST Client, Database Connection, OData, Business Events], [🔄 Syntax designed --- REST has initial syntax; OData (external entities/actions), business events, and database connector lack implementation],
    [#strong[Security];], [Module Roles, Entity Access, Project Security], [📋 Planned --- no implementation or detailed design yet],
    [#strong[Configuration];], [Constants, Scheduled Events, Project Settings, Module Settings], [📋 Planned --- not yet implemented],
    [#strong[Code Analysis];], [Linting, Impact Analysis, Cross-references, Search], [🔄 In progress --- linting and catalog refs working; coverage expanding],
    [#strong[Workflows];], [Workflow definitions, User Tasks, Decisions], [📋 Planned --- not yet implemented],
    [#strong[Mobile];], [Native pages, Navigation profiles, Offline sync], [📋 Planned --- not yet implemented],
    [#strong[Styling];], [Atlas theme customization, Design properties], [📋 Planned --- not yet implemented],
    [#strong[Data Importer];], [Import templates, Data mapping], [📋 Planned --- not yet implemented],
    [#strong[Mappings];], [Import/Export mappings, JSON snippets, XML schemas], [📋 Planned --- not yet implemented],
    [#strong[Advanced];], [JavaScript actions, Task queues, Regular expressions], [📋 Planned --- not yet implemented],
  )]
  , kind: table
  )

== Why Declarative DSL? Comparing Approaches
<why-declarative-dsl-comparing-approaches>
When enabling AI agents to manipulate Mendix models, several approaches are possible. Here's why a declarative DSL (MDL) is the optimal choice:

#figure(
  align(center)[#table(
    columns: (11.43%, 18.57%, 22.86%, 24.29%, 22.86%),
    align: (auto,auto,auto,auto,auto,),
    table.header([Aspect], [Direct JSON], [TypeScript API], [Declarative DSL], [Graph (SPARQL)],),
    table.hline(),
    [#strong[LLM token efficiency];], [❌ 20K+ lines], [⚠️ Verbose], [✅ Compact], [✅ Query-focused],
    [#strong[LLM accuracy];], [❌ Poor at JSON editing], [✅ Well-trained on TS], [✅ Simple syntax], [❌ SPARQL errors],
    [#strong[Paradigm match];], [⚠️ Data only], [❌ Procedural for declarative domain], [✅ Declarative ↔ declarative], [✅ Graphs match relations],
    [#strong[Two-language problem];], [❌ JSON + mental model], [❌ TS + Mendix interleaved], [✅ Single language], [⚠️ Query + update separate],
    [#strong[Human review];], [❌ Unreadable], [⚠️ Must trace execution], [✅ Diff is scannable], [❌ Queries not reviewable],
    [#strong[Skill authoring];], [❌ Impractical], [⚠️ Verbose boilerplate], [✅ Readable patterns], [⚠️ SHACL is complex],
    [#strong[Modification precision];], [❌ Error-prone], [✅ Atomic operations], [✅ Structural identity], [⚠️ SPARQL UPDATE complex],
  )]
  , kind: table
  )

#strong[Key insight];: Mendix microflows are already declarative. Using a procedural API to describe declarative structures creates a paradigm mismatch. MDL's declarative syntax matches the declarative nature of Mendix models.

== CLI Tools vs MCP Service
<cli-tools-vs-mcp-service>
A Mendix project is more than just the model---it includes Java actions, JavaScript widgets, SCSS styling, and more:

```
my-app/
├── model/*.json           ← Model (MDL/mxcli domain)
├── javasource/            ← Java actions (file editing)
├── javascriptsource/      ← Custom widgets (file editing)
├── themesource/           ← SCSS styling (file editing)
└── widgets/               ← Widget packages
```

#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Aspect], [Studio Pro MCP], [CLI Tools (mxcli + file ops)],),
    table.hline(),
    [#strong[Model access];], [✅ Full model API], [✅ Via MDL parsing],
    [#strong[Filesystem access];], [❌ Out of scope], [✅ Native file operations],
    [#strong[Java/JS editing];], [❌ Outside model], [✅ Standard file editing],
    [#strong[Atomic commits];], [⚠️ Model only], [✅ Git for everything],
    [#strong[Offline operation];], [❌ Requires Studio Pro], [✅ Works standalone],
    [#strong[CI/CD integration];], [⚠️ Complex], [✅ Natural fit],
    [#strong[Agent autonomy];], [⚠️ Bound to Studio Pro], [✅ Full control],
  )]
  , kind: table
  )

#strong[Recommendation];: CLI-based approach (mxcli) has a fundamental advantage---a Mendix project is a directory with files. An agent building a complete feature (entity + microflow + Java action + page + styling) works more effectively via CLI tools that treat the entire project as a filesystem.

== Hybrid Architecture: Model + Files
<hybrid-architecture-model-files>
#figure(image("diagrams/mermaid_3.png"),
  caption: [
    Diagram 3
  ]
)

#horizontalrule

= Architecture: Enabling Agentic IDEs
<architecture-enabling-agentic-ides>
== System Overview
<system-overview>
#figure(image("diagrams/mermaid_4.png"),
  caption: [
    Diagram 4
  ]
)

== Capabilities Required by Agentic IDEs
<capabilities-required-by-agentic-ides>
AI coding agents need specific capabilities to work effectively. Here's how MDL tooling provides them:

=== 1. Code Generation & Modification
<code-generation-modification>
#figure(
  align(center)[#table(
    columns: (30.77%, 33.33%, 35.9%),
    align: (auto,auto,auto,),
    table.header([Capability], [Traditional], [MDL Approach],),
    table.hline(),
    [#strong[Create files];], [Write to filesystem], [`CREATE ENTITY`, `CREATE MICROFLOW`],
    [#strong[Modify code];], [Text-based edits], [`ALTER ENTITY`, `CREATE OR MODIFY`],
    [#strong[Delete code];], [File deletion], [`DROP ENTITY`, `DROP MICROFLOW`],
    [#strong[Refactor];], [AST manipulation], [Declarative re-generation],
  )]
  , kind: table
  )

#strong[MDL Advantage];: Atomic operations with clear semantics. AI doesn't need to understand file structure.

=== 2. Code Understanding & Search
<code-understanding-search>
#figure(
  align(center)[#table(
    columns: (30.77%, 33.33%, 35.9%),
    align: (auto,auto,auto,),
    table.header([Capability], [Traditional], [MDL Approach],),
    table.hline(),
    [#strong[Search code];], [grep, ripgrep], [`SEARCH 'keyword'` or `SELECT * FROM CATALOG.ENTITIES WHERE ...`],
    [#strong[Find usages];], [LSP references], [`SHOW REFERENCES TO Module.Entity`],
    [#strong[Find callers];], [Manual tracing], [`SHOW CALLERS OF Module.Microflow [TRANSITIVE]`],
    [#strong[Impact analysis];], [Manual review], [`SHOW IMPACT OF Module.Entity`],
    [#strong[Understand structure];], [Parse AST], [`DESCRIBE ENTITY`, `DESCRIBE MICROFLOW`],
    [#strong[Navigate relationships];], [Go to definition], [`SHOW CALLEES OF`, CATALOG JOINs],
    [#strong[Assemble context];], [Manual gathering], [`SHOW CONTEXT OF Module.Microflow DEPTH 3`],
  )]
  , kind: table
  )

#strong[MDL Advantage];: High-level code search commands and SQL-like queries for semantic search. No regex pattern matching on source code.

```sql
-- Full-text search across strings and source
SEARCH 'validation'

-- Find what calls a microflow (direct and transitive)
SHOW CALLERS OF CRM.ACT_Customer_Save TRANSITIVE

-- Find what a microflow calls
SHOW CALLEES OF CRM.ProcessOrder

-- Find all references to an entity
SHOW REFERENCES TO CRM.Customer

-- Analyze impact of changing an element
SHOW IMPACT OF CRM.Customer

-- Assemble context for LLM consumption
SHOW CONTEXT OF CRM.ACT_Customer_Save DEPTH 3

-- SQL-based queries via CATALOG
SELECT SourceName, RefKind, TargetName
FROM CATALOG.REFS
WHERE TargetName = 'CRM.Customer';
```

=== 3. Validation & Testing
<validation-testing>
#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Capability], [Traditional], [MDL Approach],),
    table.hline(),
    [#strong[Syntax check];], [Compiler/linter], [MDL parser + Language Server],
    [#strong[Type check];], [TypeScript/mypy], [Model validation],
    [#strong[Integration test];], [Jest/pytest], [Mendix test suite],
    [#strong[Security scan];], [SAST tools], [Platform security model],
  )]
  , kind: table
  )

#strong[MDL Advantage];: Multi-level validation catches errors early.

#figure(image("diagrams/mermaid_5.png"),
  caption: [
    Diagram 5
  ]
)

=== 4. Execution & Feedback
<execution-feedback>
#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Capability], [Traditional], [MDL Approach],),
    table.hline(),
    [#strong[Run code];], [Node/Python process], [REPL daemon execution],
    [#strong[Get errors];], [Stack traces], [Structured error messages],
    [#strong[Hot reload];], [Dev server], [Model sync],
    [#strong[Debug];], [Debugger protocol], [Microflow debugging],
  )]
  , kind: table
  )

#strong[MDL Advantage];: Immediate execution with structured feedback.

```bash
# AI validates MDL syntax before execution
$ mxcli check script.mdl
✓ Syntax OK

# AI validates with reference checking
$ mxcli check script.mdl -p app.mpr --references
Error: Entity 'CRM.InvalidEntity' not found

# AI executes MDL against project
$ mxcli exec script.mdl -p app.mpr
Created entity: CRM.Lead
Created microflow: CRM.ACT_Lead_Save

# AI runs linting for best practices
$ mxcli lint -p app.mpr --format sarif
{
  "runs": [{ "results": [...] }]
}

# AI can self-correct and retry based on feedback
```

#horizontalrule

= Skills & Prompting: Guiding AI Generation
<skills-prompting-guiding-ai-generation>
== The Challenge: AI Needs Context
<the-challenge-ai-needs-context>
AI models are powerful but generic. Without domain-specific guidance, they may: - Generate syntactically valid but semantically incorrect MDL - Miss Mendix best practices and conventions - Create inefficient or non-idiomatic solutions - Overlook security or validation requirements

#strong[Skills] (structured prompt templates) bridge this gap by providing: - Domain knowledge about Mendix patterns - MDL syntax examples and conventions - Best practices for specific scenarios - Validation rules and constraints

== Skill Architecture
<skill-architecture>
#figure(image("diagrams/mermaid_6.png"),
  caption: [
    Diagram 6
  ]
)

#strong[Example: Skill-guided validation microflow generation]

The AI receives the skill content which provides: - Naming convention: `VAL_{EntityName}` - Required parameter: the entity to validate - Pattern: use `VALIDATION FEEDBACK` for errors, return `Boolean`

The AI then generates a complete validation microflow following the pattern:

```sql
CREATE MICROFLOW CRM.VAL_Customer ($Customer: CRM.Customer)
RETURNS Boolean
BEGIN
  -- Email validation
  IF $Customer/Email = empty THEN
    VALIDATION FEEDBACK $Customer MEMBER Email MESSAGE 'Email is required';
    RETURN false;
  END IF;

  IF NOT contains($Customer/Email, '@') THEN
    VALIDATION FEEDBACK $Customer MEMBER Email MESSAGE 'Invalid email format';
    RETURN false;
  END IF;

  -- Credit limit validation
  IF $Customer/CreditLimit < 0 THEN
    VALIDATION FEEDBACK $Customer MEMBER CreditLimit MESSAGE 'Cannot be negative';
    RETURN false;
  END IF;

  RETURN true;
END;
```

== Skill Categories
<skill-categories>
#figure(
  align(center)[#table(
    columns: (32.26%, 29.03%, 38.71%),
    align: (auto,auto,auto,),
    table.header([Category], [Purpose], [Skill File],),
    table.hline(),
    [#strong[Domain Model];], [Entity/association patterns], [`generate-domain-model.md`],
    [#strong[Microflow];], [Business logic patterns], [`write-microflows.md`],
    [#strong[Pages];], [UI patterns (V3 syntax)], [`create-page.md`, `overview-pages.md`, `master-detail-pages.md`],
    [#strong[Validation];], [Pre-flight checks], [`check-syntax.md`],
    [#strong[Debugging];], [BSON/serialization issues], [`debug-bson.md`],
  )]
  , kind: table
  )

These skills are installed via `mxcli init` into the `.claude/skills/` folder of a Mendix project.

== Example Skills
<example-skills>
=== Validation Microflow Skill
<validation-microflow-skill>
```markdown
# Validation Microflow Generation

## When to Use
Generate validation microflows when the user requests:
- Input validation for entities
- Business rule enforcement
- Data quality checks

## Naming Convention
- Microflow name: `VAL_{EntityName}` or `Validate{EntityName}`
- Place in same module as entity

## MDL Pattern

CREATE MICROFLOW {Module}.VAL_{EntityName} (
  ${EntityName}: {Module}.{EntityName}
)
RETURNS Boolean
BEGIN
  -- Required field validation
  IF ${EntityName}/{RequiredField} = empty THEN
    VALIDATION FEEDBACK ${EntityName} MEMBER {RequiredField}
      MESSAGE '{FieldLabel} is required';
    RETURN false;
  END IF;

  -- Format validation (email, phone, etc.)
  IF NOT {formatCheck} THEN
    VALIDATION FEEDBACK ${EntityName} MEMBER {Field}
      MESSAGE '{ValidationMessage}';
    RETURN false;
  END IF;

  -- Range validation
  IF ${EntityName}/{NumericField} < {Min} OR ${EntityName}/{NumericField} > {Max} THEN
    VALIDATION FEEDBACK ${EntityName} MEMBER {NumericField}
      MESSAGE '{FieldLabel} must be between {Min} and {Max}';
    RETURN false;
  END IF;

  RETURN true;
END;

## Common Validations

| Type | MDL Expression |
|------|----------------|
| Required | `$Entity/Field = empty` |
| Email format | `NOT contains($Entity/Email, '@')` |
| Min length | `length($Entity/Field) < {min}` |
| Max length | `length($Entity/Field) > {max}` |
| Positive number | `$Entity/Amount <= 0` |
| Date in future | `$Entity/Date < [%CurrentDateTime%]` |
| Regex match | `NOT matches($Entity/Field, '{pattern}')` |

## Example: Complete Customer Validation

CREATE MICROFLOW CRM.VAL_Customer ($Customer: CRM.Customer)
RETURNS Boolean
BEGIN
  -- Name is required
  IF $Customer/Name = empty THEN
    VALIDATION FEEDBACK $Customer MEMBER Name
      MESSAGE 'Customer name is required';
    RETURN false;
  END IF;

  -- Email format
  IF $Customer/Email != empty AND NOT contains($Customer/Email, '@') THEN
    VALIDATION FEEDBACK $Customer MEMBER Email
      MESSAGE 'Please enter a valid email address';
    RETURN false;
  END IF;

  -- Credit limit range
  IF $Customer/CreditLimit < 0 THEN
    VALIDATION FEEDBACK $Customer MEMBER CreditLimit
      MESSAGE 'Credit limit cannot be negative';
    RETURN false;
  END IF;

  IF $Customer/CreditLimit > 1000000 THEN
    VALIDATION FEEDBACK $Customer MEMBER CreditLimit
      MESSAGE 'Credit limit cannot exceed 1,000,000';
    RETURN false;
  END IF;

  RETURN true;
END;
```

=== Overview Page Skill
<overview-page-skill>
````markdown
# Overview Page Generation

## When to Use
Generate overview pages for:
- Entity list views
- Master-detail layouts
- Searchable data grids

## Naming Convention
- Page name: `{EntityName}_Overview`
- Place in same module as entity

## MDL Pattern

```sql
CREATE PAGE {Module}.{EntityName}_Overview
(
  Title: '{Entity Display Name} Overview',
  Layout: Atlas_Core.Atlas_Default
)
{
  LAYOUTGRID gridMain {
    ROW row1 {
      COLUMN col1 (Weight: 12) {
        CONTAINER containerHeader {
          DYNAMICTEXT txtTitle (Content: '{Entity Display Name}', RenderMode: H1)
          ACTIONBUTTON btnNew (
            Caption: 'New {EntityName}',
            Action: CREATE_OBJECT {Module}.{EntityName},
            OnClickPage: {Module}.{EntityName}_NewEdit,
            Style: Primary
          )
        }
      }
      COLUMN col2 (Weight: 12) {
        DATAGRID grid{EntityName} (DataSource: DATABASE {Module}.{EntityName}) {
          COLUMN col{Attr} (AttributePath: {AttributeName}, Caption: '{Display Label}')
          -- ... more columns
          CONTROLBAR ctrlBar {
            SEARCH search1
            PAGING paging1
          }
        }
      }
    }
  }
}
````

= Example: Customer Overview
<example-customer-overview>
```sql
CREATE PAGE CRM.Customer_Overview
(
  Title: 'Customer Overview',
  Layout: Atlas_Core.Atlas_Default
)
{
  LAYOUTGRID gridMain {
    ROW row1 {
      COLUMN colHeader (Weight: 12) {
        CONTAINER containerHeader {
          DYNAMICTEXT txtTitle (Content: 'Customers', RenderMode: H1)
          ACTIONBUTTON btnNew (
            Caption: 'New Customer',
            Action: CREATE_OBJECT CRM.Customer,
            OnClickPage: CRM.Customer_NewEdit,
            Style: Primary
          )
        }
      }
      COLUMN colGrid (Weight: 12) {
        DATAGRID gridCustomer (DataSource: DATABASE CRM.Customer) {
          COLUMN colName (AttributePath: Name, Caption: 'Customer Name')
          COLUMN colEmail (AttributePath: Email, Caption: 'Email Address')
          COLUMN colCredit (AttributePath: CreditLimit, Caption: 'Credit Limit')
          COLUMN colActive (AttributePath: IsActive, Caption: 'Active')
          CONTROLBAR ctrlBar {
            SEARCH search1
            PAGING paging1
          }
        }
      }
    }
  }
}
```

=== CRUD Microflow Skill
<crud-microflow-skill>
```markdown
# CRUD Microflow Generation

## Save Microflow Pattern

CREATE MICROFLOW {Module}.ACT_{EntityName}_Save (
  ${EntityName}: {Module}.{EntityName}
)
RETURNS Boolean
BEGIN
  -- Validate
  $IsValid = CALL {Module}.VAL_{EntityName} ($Object = ${EntityName});

  IF NOT $IsValid THEN
    RETURN false;
  END IF;

  -- Save
  COMMIT ${EntityName};

  -- Provide feedback
  LOG MESSAGE 'INFO' '{EntityName} saved successfully';

  RETURN true;
END;

## Delete Microflow Pattern

CREATE MICROFLOW {Module}.ACT_{EntityName}_Delete (
  ${EntityName}: {Module}.{EntityName}
)
BEGIN
  -- Confirm deletion (handled by calling page)

  -- Check for dependencies
  RETRIEVE $RelatedObjects FROM DATABASE
    WHERE {RelatedEntity}/{Association} = ${EntityName};

  IF $RelatedObjects != empty THEN
    LOG MESSAGE 'WARNING' 'Cannot delete: related records exist';
    RETURN;
  END IF;

  -- Delete
  DELETE ${EntityName};

  -- Feedback
  LOG MESSAGE 'INFO' '{EntityName} deleted successfully';

  CLOSE PAGE;
END;
```

== Skills in Practice: End-to-End Example
<skills-in-practice-end-to-end-example>
#strong[User Request];: "Create a complete customer management module with validation"

#strong[AI Response (using multiple skills)];:

```sql
-- 1. Entity (using entity skill)
CREATE PERSISTENT ENTITY CustomerMgmt.Customer (
  Code: String(20) NOT NULL UNIQUE,
  Name: String(200) NOT NULL,
  Email: String(200),
  Phone: String(20),
  CreditLimit: Decimal DEFAULT 0,
  IsActive: Boolean DEFAULT true,
  CreatedDate: DateTime DEFAULT '[%CurrentDateTime%]',
  ModifiedDate: DateTime
);

CREATE ENUMERATION CustomerMgmt.CustomerStatus (
  ACTIVE 'Active',
  INACTIVE 'Inactive',
  SUSPENDED 'Suspended'
);

-- 2. Validation Microflow (using validation skill)
CREATE MICROFLOW CustomerMgmt.VAL_Customer ($Customer: CustomerMgmt.Customer)
RETURNS Boolean
BEGIN
  IF $Customer/Code = empty THEN
    VALIDATION FEEDBACK $Customer MEMBER Code MESSAGE 'Customer code is required';
    RETURN false;
  END IF;

  IF $Customer/Name = empty THEN
    VALIDATION FEEDBACK $Customer MEMBER Name MESSAGE 'Customer name is required';
    RETURN false;
  END IF;

  IF $Customer/Email != empty AND NOT contains($Customer/Email, '@') THEN
    VALIDATION FEEDBACK $Customer MEMBER Email MESSAGE 'Invalid email format';
    RETURN false;
  END IF;

  IF $Customer/CreditLimit < 0 THEN
    VALIDATION FEEDBACK $Customer MEMBER CreditLimit MESSAGE 'Credit limit cannot be negative';
    RETURN false;
  END IF;

  RETURN true;
END;

-- 3. Save Microflow (using CRUD skill)
CREATE MICROFLOW CustomerMgmt.ACT_Customer_Save ($Customer: CustomerMgmt.Customer)
RETURNS Boolean
BEGIN
  $IsValid = CALL CustomerMgmt.VAL_Customer ($Customer);
  IF NOT $IsValid THEN
    RETURN false;
  END IF;

  CHANGE $Customer SET ModifiedDate = [%CurrentDateTime%];
  COMMIT $Customer;
  RETURN true;
END;

-- 4. Overview Page (using page skill)
CREATE PAGE CustomerMgmt.Customer_Overview
(
  Title: 'Customers',
  Layout: Atlas_Core.Atlas_Default
)
{
  LAYOUTGRID gridMain {
    ROW row1 {
      COLUMN colHeader (Weight: 12) {
        DYNAMICTEXT txtHeader (Content: 'Customer Management', RenderMode: H1)
        ACTIONBUTTON btnNew (
          Caption: 'New Customer',
          Action: CREATE_OBJECT CustomerMgmt.Customer,
          OnClickPage: CustomerMgmt.Customer_NewEdit,
          Style: Primary
        )
      }
      COLUMN colGrid (Weight: 12) {
        DATAGRID gridCustomer (DataSource: DATABASE CustomerMgmt.Customer) {
          COLUMN colCode (AttributePath: Code, Caption: 'Code')
          COLUMN colName (AttributePath: Name, Caption: 'Name')
          COLUMN colEmail (AttributePath: Email, Caption: 'Email')
          COLUMN colCredit (AttributePath: CreditLimit, Caption: 'Credit Limit')
          COLUMN colActive (AttributePath: IsActive, Caption: 'Active')
          CONTROLBAR ctrlBar {
            SEARCH search1
            PAGING paging1
          }
        }
      }
    }
  }
}

-- 5. Edit Page (using page skill)
CREATE PAGE CustomerMgmt.Customer_NewEdit
(
  Params: { $Customer: CustomerMgmt.Customer },
  Title: 'Customer Details',
  Layout: Atlas_Core.Atlas_Default
)
{
  DATAVIEW dvCustomer (DataSource: $Customer) {
    TEXTBOX txtCode (Label: 'Code', Binds: Code)
    TEXTBOX txtName (Label: 'Name', Binds: Name)
    TEXTBOX txtEmail (Label: 'Email', Binds: Email)
    TEXTBOX txtPhone (Label: 'Phone', Binds: Phone)
    TEXTBOX txtCreditLimit (Label: 'Credit Limit', Binds: CreditLimit)
    CHECKBOX chkActive (Label: 'Active', Binds: IsActive)

    FOOTER footerButtons {
      ACTIONBUTTON btnSave (
        Caption: 'Save',
        Action: CALL_MICROFLOW CustomerMgmt.ACT_Customer_Save,
        Style: Primary
      )
      ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CLOSE_PAGE)
    }
  }
}
```

== Skill Distribution
<skill-distribution>
Skills are distributed and customized at multiple levels:

#figure(
  align(center)[#table(
    columns: (26.92%, 38.46%, 34.62%),
    align: (auto,auto,auto,),
    table.header([Level], [Location], [Purpose],),
    table.hline(),
    [#strong[Built-in];], [Embedded in `mxcli` binary], [Core patterns, installed via `mxcli init`],
    [#strong[Organization];], [Shared repository], [Company-specific conventions],
    [#strong[Project];], [`.claude/skills/` folder], [Project-specific customizations],
    [#strong[User];], [`~/.claude/` folder], [Personal preferences],
  )]
  , kind: table
  )

The `mxcli init` command copies built-in skills from the embedded templates (sourced from `reference/mendix-repl/templates/.claude/skills/`) into the target project's `.claude/skills/` folder.

== Benefits of Skill-Based Generation
<benefits-of-skill-based-generation>
#figure(
  align(center)[#table(
    columns: (40.91%, 59.09%),
    align: (auto,auto,),
    table.header([Benefit], [Description],),
    table.hline(),
    [#strong[Consistency];], [Same patterns applied across all generated code],
    [#strong[Quality];], [Best practices embedded in skills],
    [#strong[Efficiency];], [AI doesn't need to rediscover patterns],
    [#strong[Customization];], [Organizations can define their own standards],
    [#strong[Maintainability];], [Update skill once, apply everywhere],
    [#strong[Onboarding];], [New team members learn patterns from skills],
    [#strong[Knowledge Compounding];], [Skills capture lessons learned and encode them permanently. Each bug fix, performance optimization, or pattern improvement added to a skill benefits all future generations. Organizations build cumulative expertise that compounds over time rather than being lost to turnover or forgotten in wikis.],
    [#strong[Composability];], [Skills can reference and build upon other skills, enabling complex workflows from simple building blocks. A "create CRUD pages" skill can compose "create entity", "create overview page", and "create edit page" skills. This modularity allows teams to mix and match capabilities while keeping individual skills focused and testable.],
  )]
  , kind: table
  )

#horizontalrule

= Graph-Based Model Access: Beyond SQL Catalogs
<graph-based-model-access-beyond-sql-catalogs>
== The Limitation of SQL Catalogs
<the-limitation-of-sql-catalogs>
The current CATALOG system uses SQL tables to expose Mendix metadata:

```sql
SELECT * FROM CATALOG.ENTITIES WHERE ModuleName = 'CRM';
SELECT * FROM CATALOG.ASSOCIATIONS WHERE ParentEntity = 'CRM.Customer';
```

While familiar and powerful for simple queries, SQL has limitations for complex model navigation:

#figure(
  align(center)[#table(
    columns: (57.14%, 42.86%),
    align: (auto,auto,),
    table.header([Limitation], [Example],),
    table.hline(),
    [#strong[Path queries];], ["Find all entities reachable from Customer in 3 hops"],
    [#strong[Pattern matching];], ["Find circular dependencies in associations"],
    [#strong[Graph algorithms];], ["What's the shortest path between Order and Invoice?"],
    [#strong[Bulk traversal];], ["Get the complete subgraph for module CRM"],
  )]
  , kind: table
  )

== Current Implementation: Search and Reference Tracking
<current-implementation-search-and-reference-tracking>
While full graph query support remains a future goal, mxcli already implements practical code navigation through the #strong[REFS table] and high-level search commands:

=== The REFS Table
<the-refs-table>
The CATALOG.REFS table tracks cross-references between model elements:

```sql
-- Requires: REFRESH CATALOG FULL
SELECT SourceName, RefKind, TargetName
FROM CATALOG.REFS
WHERE TargetName = 'MyModule.Customer';
```

#figure(
  align(center)[#table(
    columns: (29.03%, 41.94%, 29.03%),
    align: (auto,auto,auto,),
    table.header([RefKind], [Description], [Example],),
    table.hline(),
    [`call`], [Microflow calls microflow], [ACT\_Save → SUB\_Validate],
    [`create`], [Microflow creates entity], [ACT\_New → Customer],
    [`retrieve`], [Microflow retrieves entity], [ACT\_List → Customer],
    [`change`], [Microflow changes entity], [ACT\_Update → Customer],
    [`delete`], [Microflow deletes entity], [ACT\_Remove → Customer],
    [`show_page`], [Microflow shows page], [ACT\_Edit → Customer\_Edit],
    [`generalize`], [Entity extends entity], [Employee → Person],
    [`layout`], [Page uses layout], [Customer\_Edit → PopupLayout],
    [`datasource`], [Widget uses entity], [DataGrid → Customer],
    [`parameter`], [Page parameter typed to entity], [Customer\_Edit(\$Customer)],
    [`action`], [Widget calls microflow], [Button → ACT\_Save],
  )]
  , kind: table
  )

=== High-Level Search Commands
<high-level-search-commands>
mxcli provides developer-friendly commands that query the REFS table:

```bash
# Find what calls a microflow (direct callers)
mxcli callers -p app.mpr Module.MyMicroflow

# Find transitive callers (full call chain)
mxcli callers -p app.mpr Module.MyMicroflow --transitive

# Find what a microflow calls
mxcli callees -p app.mpr Module.MyMicroflow

# Find all references to an element
mxcli refs -p app.mpr Module.Customer

# Analyze impact of changing an element
mxcli impact -p app.mpr Module.Customer

# Assemble context for LLM consumption (with depth control)
mxcli context -p app.mpr Module.MyMicroflow --depth 3
```

Or via MDL syntax in the REPL:

```sql
SHOW CALLERS OF Module.MyMicroflow;
SHOW CALLERS OF Module.MyMicroflow TRANSITIVE;
SHOW CALLEES OF Module.MyMicroflow;
SHOW REFERENCES TO Module.Customer;
SHOW IMPACT OF Module.Customer;
SHOW CONTEXT OF Module.MyMicroflow DEPTH 3;
```

=== Full-Text Search
<full-text-search>
Search across all strings and source in the project:

```bash
# Find all occurrences of a term
mxcli search -p app.mpr "validation"

# Output formats for piping
mxcli search -p app.mpr "error" --format names   # type<TAB>name per line
mxcli search -p app.mpr "error" --format json    # JSON array

# Pipe to other commands
mxcli search -p app.mpr "error" -q --format names | head -1 | \
  awk '{print $2}' | xargs mxcli describe -p app.mpr microflow
```

Or via MDL:

```sql
SEARCH 'validation';
```

=== Why This Matters for AI Agents
<why-this-matters-for-ai-agents>
These commands provide AI agents with essential code navigation capabilities:

#figure(
  align(center)[#table(
    columns: (35.29%, 26.47%, 38.24%),
    align: (auto,auto,auto,),
    table.header([Capability], [Command], [AI Use Case],),
    table.hline(),
    [#strong[Impact analysis];], [`mxcli impact`], [Before modifying an entity, understand what will break],
    [#strong[Call graph];], [`mxcli callers/callees`], [Understand microflow dependencies before refactoring],
    [#strong[Reference lookup];], [`mxcli refs`], [Find all usages before renaming or deleting],
    [#strong[Context gathering];], [`mxcli context`], [Build relevant context window for LLM prompts],
    [#strong[Code search];], [`mxcli search`], [Find where specific patterns or values are used],
  )]
  , kind: table
  )

The catalog is cached in `.mxcli/catalog.db` next to the MPR file. Use `REFRESH CATALOG FULL FORCE` to rebuild after external changes.

== Mendix Models as Graphs
<mendix-models-as-graphs>
A Mendix project is fundamentally a #strong[knowledge graph];:

#figure(image("diagrams/mermaid_7.png"),
  caption: [
    Diagram 7
  ]
)

#strong[Node Types];: Entity, Attribute, Association, Page, Microflow, Widget, … #strong[Edge Types];: has\_attribute, association, uses, contains, calls, …

== Option A: SPARQL Interface
<option-a-sparql-interface>
SPARQL (SPARQL Protocol and RDF Query Language) is the W3C standard for querying RDF graphs:

```sparql
PREFIX mx: <http://mendix.com/model#>

# Find all entities with their attributes
SELECT ?entity ?attrName ?attrType
WHERE {
  ?entity a mx:Entity .
  ?entity mx:hasAttribute ?attr .
  ?attr mx:name ?attrName .
  ?attr mx:type ?attrType .
}

# Find entities reachable from Customer within 3 association hops
SELECT ?entity ?path
WHERE {
  mx:CRM.Customer (mx:association)* ?entity .
  FILTER (?entity != mx:CRM.Customer)
}

# Find circular dependencies (entities that reference each other)
SELECT ?entity1 ?entity2
WHERE {
  ?entity1 mx:association ?entity2 .
  ?entity2 mx:association ?entity1 .
  FILTER (?entity1 < ?entity2)  # Avoid duplicates
}

# Find all microflows that could affect Customer data
SELECT ?microflow
WHERE {
  ?microflow a mx:Microflow .
  ?microflow mx:usesEntity mx:CRM.Customer .
}

# CONSTRUCT a subgraph for export
CONSTRUCT {
  ?entity a mx:Entity .
  ?entity mx:hasAttribute ?attr .
  ?entity mx:association ?target .
}
WHERE {
  ?entity mx:inModule mx:CRM .
  OPTIONAL { ?entity mx:hasAttribute ?attr }
  OPTIONAL { ?entity mx:association ?target }
}
```

#strong[SPARQL Advantages:] - W3C standard with mature tooling - Powerful pattern matching - CONSTRUCT for graph extraction - Federated queries across models

== Option B: Cypher Interface
<option-b-cypher-interface>
Cypher (used by Neo4j) offers a more visual, ASCII-art-like syntax:

```cypher
// Find all entities with their attributes
MATCH (e:Entity)-[:HAS_ATTRIBUTE]->(a:Attribute)
RETURN e.name, a.name, a.type

// Find path from Order to Customer (any length)
MATCH path = (o:Entity {name: 'Order'})-[:ASSOCIATION*1..5]->(c:Entity {name: 'Customer'})
RETURN path

// Find circular dependencies
MATCH (e1:Entity)-[:ASSOCIATION]->(e2:Entity)-[:ASSOCIATION]->(e1)
WHERE id(e1) < id(e2)
RETURN e1.name, e2.name

// Find impact of changing Customer entity
MATCH (c:Entity {name: 'Customer'})<-[:USES|ASSOCIATION*1..3]-(dependent)
RETURN DISTINCT dependent.name, labels(dependent)[0] AS type

// Clone a module's structure
MATCH (m:Module {name: 'CRM'})-[*]->(n)
RETURN m, n

// Find orphaned entities (no associations, no page usage)
MATCH (e:Entity)
WHERE NOT (e)-[:ASSOCIATION]-()
  AND NOT ()-[:USES]->(e)
RETURN e.name AS orphaned_entity

// Bulk update: Add audit fields to all entities in module
MATCH (e:Entity)-[:IN_MODULE]->(m:Module {name: 'CRM'})
WHERE NOT (e)-[:HAS_ATTRIBUTE]->(:Attribute {name: 'CreatedDate'})
CREATE (e)-[:HAS_ATTRIBUTE]->(a:Attribute {name: 'CreatedDate', type: 'DateTime'})
RETURN e.name, a.name
```

#strong[Cypher Advantages:] - Visual, intuitive syntax - Excellent for path queries - Read and write in same language - Popular in developer community

== Option C: GraphQL Interface
<option-c-graphql-interface>
GraphQL provides a typed, hierarchical query interface:

```graphql
# Schema
type Entity {
  id: ID!
  name: String!
  module: Module!
  attributes: [Attribute!]!
  associations: [Association!]!
  usedByPages: [Page!]!
  usedByMicroflows: [Microflow!]!
}

type Query {
  entity(name: String!): Entity
  entities(module: String): [Entity!]!
  impactAnalysis(entityName: String!, depth: Int): ImpactResult!
  pathBetween(from: String!, to: String!): [Path!]!
}

type Mutation {
  createEntity(input: CreateEntityInput!): Entity!
  addAttribute(entityName: String!, attribute: AttributeInput!): Attribute!
  bulkAddAttribute(filter: EntityFilter!, attribute: AttributeInput!): [Entity!]!
}

# Query: Get entity with all relationships
query GetCustomerImpact {
  entity(name: "CRM.Customer") {
    name
    attributes {
      name
      type
    }
    associations {
      name
      targetEntity {
        name
      }
    }
    usedByPages {
      name
      widgets {
        type
      }
    }
    usedByMicroflows {
      name
      actions {
        type
      }
    }
  }
}

# Mutation: Add audit fields to multiple entities
mutation AddAuditFields {
  bulkAddAttribute(
    filter: { module: "CRM" }
    attribute: { name: "CreatedDate", type: DATETIME }
  ) {
    name
  }
}
```

#strong[GraphQL Advantages:] - Strongly typed with schema - Hierarchical, matches model structure - Efficient (fetch only needed fields) - Great tooling (GraphiQL, Apollo)

== Comparison Matrix
<comparison-matrix>
#figure(
  align(center)[#table(
    columns: (24%, 26%, 16%, 16%, 18%),
    align: (auto,auto,auto,auto,auto,),
    table.header([Capability], [SQL Catalog], [SPARQL], [Cypher], [GraphQL],),
    table.hline(),
    [Simple queries], [✅ Excellent], [✅ Good], [✅ Good], [✅ Excellent],
    [Path queries], [❌ Recursive CTEs], [✅ Native], [✅ Excellent], [⚠️ Limited],
    [Pattern matching], [⚠️ Limited], [✅ Excellent], [✅ Excellent], [❌ No],
    [Graph algorithms], [❌ No], [⚠️ Extensions], [✅ Built-in], [❌ No],
    [Write operations], [✅ Via MDL], [✅ SPARQL Update], [✅ Native], [✅ Mutations],
    [Type safety], [⚠️ Runtime], [⚠️ Runtime], [⚠️ Runtime], [✅ Schema],
    [Tooling], [✅ Ubiquitous], [✅ Good], [✅ Good], [✅ Excellent],
    [AI familiarity], [✅ Very High], [⚠️ Medium], [⚠️ Medium], [✅ High],
  )]
  , kind: table
  )

== Recommended Approach: Hybrid Architecture
<recommended-approach-hybrid-architecture>
Rather than choosing one, offer multiple interfaces to the same underlying graph:

#figure(image("diagrams/mermaid_8.png"),
  caption: [
    Diagram 8
  ]
)

== Use Cases by Query Language
<use-cases-by-query-language>
#figure(
  align(center)[#table(
    columns: (29.41%, 44.12%, 26.47%),
    align: (auto,auto,auto,),
    table.header([Use Case], [Best Language], [Example],),
    table.hline(),
    [#strong[Simple lookups];], [SQL], [`SELECT * FROM CATALOG.ENTITIES WHERE Name = 'Customer'`],
    [#strong[Impact analysis];], [Cypher], [`MATCH (e)<-[:USES*1..3]-(dep) RETURN dep`],
    [#strong[Pattern detection];], [SPARQL], [Find all entities matching a naming pattern],
    [#strong[API for UI];], [GraphQL], [VS Code extension fetching model structure],
    [#strong[Bulk modifications];], [Cypher], [Add attributes to all entities in a module],
    [#strong[Export subgraph];], [SPARQL], [CONSTRUCT query for module extraction],
    [#strong[AI agent queries];], [GraphQL/SQL], [Structured, typed responses],
  )]
  , kind: table
  )

== AI Agent Benefits
<ai-agent-benefits>
Graph queries enable powerful AI agent capabilities:

#figure(image("diagrams/mermaid_9.png"),
  caption: [
    Diagram 9
  ]
)

== Implementation Considerations
<implementation-considerations>
=== Option 1: Embedded Graph Database
<option-1-embedded-graph-database>
Use an embedded graph database like #strong[DuckDB] (SQL) + custom graph extensions:

```typescript
// Build graph on project load
const graph = new MendixModelGraph();
await graph.loadFromModel(model);

// Query via SQL with graph functions
const result = await graph.query(`
  SELECT * FROM graph_traverse(
    'CRM.Customer',
    'ASSOCIATION',
    3  -- max depth
  )
`);
```

=== Option 2: RDF Triple Store
<option-2-rdf-triple-store>
Export model to RDF, query via SPARQL:

```typescript
// Export model to RDF triples
const triples = await modelToRDF(model);
const store = new N3.Store(triples);

// Query via SPARQL
const results = await store.query(`
  PREFIX mx: <http://mendix.com/model#>
  SELECT ?entity WHERE {
    ?entity a mx:Entity .
    ?entity mx:inModule mx:CRM .
  }
`);
```

=== Option 3: Virtual Graph Layer
<option-3-virtual-graph-layer>
Keep data in Model SDK, translate queries on-the-fly:

```typescript
// Query parsed and translated to Model SDK calls
const query = parseGraphQL(`
  query {
    entity(name: "CRM.Customer") {
      associations {
        targetEntity { name }
      }
    }
  }
`);

// Executed against Model SDK
const result = await executeQuery(query, model);
```

== Recommendation
<recommendation>
#strong[Phase 1 (Current)];: ✅ SQL CATALOG with REFS table and search commands - REFS table tracks cross-references (call, create, retrieve, show\_page, etc.) - High-level commands: `callers`, `callees`, `refs`, `impact`, `context`, `search` - Sufficient for most AI agent code navigation needs - Familiar SQL interface with practical abstractions

#strong[Phase 2 (Near-term)];: Add Cypher support via embedded graph - Best path/relationship queries (multi-hop traversals) - Pattern matching for detecting anti-patterns - Write operations feel natural

#strong[Phase 3 (Future)];: Add GraphQL for VS Code extension and typed API access - Schema-first development - Excellent tooling (GraphiQL, code generation) - Efficient for UI queries

#horizontalrule

= Language Server: Serving Humans and AI
<language-server-serving-humans-and-ai>
== The Dual-Purpose Design
<the-dual-purpose-design>
The MDL Language Server is architected to serve both audiences:

#figure(image("diagrams/mermaid_10.png"),
  caption: [
    Diagram 10
  ]
)

== LSP Features: Human vs AI Usage
<lsp-features-human-vs-ai-usage>
#figure(
  align(center)[#table(
    columns: (25%, 34.62%, 40.38%),
    align: (auto,auto,auto,),
    table.header([LSP Feature], [Human Experience], [AI Agent Experience],),
    table.hline(),
    [#strong[textDocument/diagnostic];], [See red/yellow squiggles inline], [Parse errors to identify what to fix],
    [#strong[textDocument/completion];], [Press Tab to complete], [Get valid options for generation],
    [#strong[textDocument/documentSymbol];], [Navigate via Outline panel], [Understand existing code structure],
    [#strong[textDocument/hover];], [Read documentation tooltip], [Lookup syntax and semantics],
    [#strong[textDocument/definition];], [Ctrl+Click to navigate], [Resolve references],
    [#strong[textDocument/codeAction];], [Click lightbulb for fixes], [Auto-apply suggested fixes],
    [#strong[textDocument/formatting];], [Shift+Alt+F to format], [Ensure consistent output],
  )]
  , kind: table
  )

== AI-Specific Features (Future)
<ai-specific-features-future>
Beyond standard LSP, we can add AI-optimized capabilities:

```typescript
// Custom LSP extension for AI agents
interface AICapabilities {
  // Generate MDL from natural language
  'mdl/generateFromDescription': {
    description: string;
    context?: string[];  // Existing entities for reference
  } => { mdl: string; confidence: number };

  // Suggest fixes with explanations
  'mdl/explainError': {
    diagnostic: Diagnostic;
  } => { explanation: string; suggestedFix: string };

  // Batch validation for efficiency
  'mdl/validateBatch': {
    documents: TextDocument[];
  } => { diagnostics: Map<string, Diagnostic[]> };

  // Semantic search
  'mdl/search': {
    query: string;
    type?: 'entity' | 'microflow' | 'page' | 'all';
  } => { results: SearchResult[] };
}
```

#horizontalrule

= VS Code Extension: The Human Interface
<vs-code-extension-the-human-interface>
== Why VS Code Matters
<why-vs-code-matters>
Even with AI generation, humans need to: - #strong[Review] generated code before deployment - #strong[Understand] what was generated - #strong[Modify] AI output for edge cases - #strong[Debug] issues in generated applications

The VS Code extension provides the human-friendly interface to MDL:

== Feature Matrix
<feature-matrix>
#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Feature], [Human Benefit], [AI Synergy],),
    table.hline(),
    [#strong[Syntax Highlighting];], [Readable code], [\-],
    [#strong[Error Squiggles];], [Spot mistakes quickly], [AI sees same errors],
    [#strong[Outline View];], [Navigate large files], [AI understands structure],
    [#strong[Autocomplete];], [Faster manual coding], [AI gets same suggestions],
    [#strong[Hover Documentation];], [Learn syntax], [AI can lookup too],
    [#strong[Go to Definition];], [Navigate codebase], [AI can follow references],
    [#strong[Format Document];], [Consistent style], [AI output is formatted],
    [#strong[Code Actions];], [Quick fixes], [AI can apply fixes],
  )]
  , kind: table
  )

== The Review Workflow
<the-review-workflow>
#figure(image("diagrams/mermaid_11.png"),
  caption: [
    Diagram 11
  ]
)

#horizontalrule

= Platform Guarantees: Why Enterprises Trust Mendix
<platform-guarantees-why-enterprises-trust-mendix>
== Security by Default
<security-by-default>
Unlike AI-generated traditional code, Mendix applications inherit platform security:

#figure(
  align(center)[#table(
    columns: (31.03%, 37.93%, 31.03%),
    align: (auto,auto,auto,),
    table.header([Security Concern], [Traditional Code Risk], [Mendix Guarantee],),
    table.hline(),
    [#strong[SQL Injection];], [Must sanitize all inputs], [ORM prevents by design],
    [#strong[XSS Attacks];], [Must escape all output], [Template engine escapes],
    [#strong[Authentication];], [Must implement correctly], [Built-in user management],
    [#strong[Authorization];], [Must check every endpoint], [Declarative role-based access],
    [#strong[Data Exposure];], [Must filter API responses], [Entity access rules],
    [#strong[CSRF];], [Must implement tokens], [Platform handles],
    [#strong[Session Management];], [Must implement securely], [Platform handles],
  )]
  , kind: table
  )

== Validation Pipeline
<validation-pipeline>
#figure(image("diagrams/mermaid_12.png"),
  caption: [
    Diagram 12
  ]
)

== Governance & Compliance
<governance-compliance>
#figure(
  align(center)[#table(
    columns: 2,
    align: (auto,auto,),
    table.header([Requirement], [How Mendix Addresses],),
    table.hline(),
    [#strong[Audit Trail];], [Model versioning, commit history],
    [#strong[Change Review];], [Visual diff in Studio Pro],
    [#strong[Separation of Concerns];], [Module security, app roles],
    [#strong[Compliance Certifications];], [SOC 2, ISO 27001, GDPR],
    [#strong[Data Residency];], [Regional cloud deployments],
    [#strong[Backup & Recovery];], [Platform-managed],
  )]
  , kind: table
  )

#horizontalrule

= Legacy Migration: AI-Assisted Modernization
<legacy-migration-ai-assisted-modernization>
== The Migration Opportunity
<the-migration-opportunity>
Enterprises have massive investments in legacy business applications that are increasingly difficult to maintain:

#figure(
  align(center)[#table(
    columns: (32.61%, 21.74%, 45.65%),
    align: (auto,auto,auto,),
    table.header([Platform Type], [Examples], [Migration Challenges],),
    table.hline(),
    [#strong[Legacy Low-Code];], [OutSystems, K2, Appian, Pega], [Proprietary formats, vendor lock-in, limited export],
    [#strong[3GL Business Apps];], [Java/Spring, .NET/C\#, COBOL], [Large codebases, tribal knowledge, documentation gaps],
    [#strong[Database-Centric];], [Oracle Forms, MS Access, FoxPro], [Tight DB coupling, stored procedures, no separation],
    [#strong[Custom Frameworks];], [Internal RAD tools, 4GL systems], [Unique syntax, no community, skills shortage],
  )]
  , kind: table
  )

#strong[AI + MDL enables a new approach];: Use agentic IDEs to understand legacy systems and generate equivalent Mendix applications.

== Migration Architecture
<migration-architecture>
#figure(image("diagrams/mermaid_13.png"),
  caption: [
    Diagram 13
  ]
)

== Migration by Source Platform
<migration-by-source-platform>
=== From OutSystems / Other Low-Code
<from-outsystems-other-low-code>
OutSystems, Appian, K2, and similar platforms have structured models that map well to Mendix:

#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([OutSystems Concept], [Mendix Equivalent], [MDL Syntax],),
    table.hline(),
    [Entity], [Entity], [`CREATE ENTITY`],
    [Entity Attribute], [Attribute], [`Name: String(200)`],
    [Entity Reference], [Association], [`CREATE ASSOCIATION`],
    [Screen], [Page], [`CREATE PAGE`],
    [Server Action], [Microflow], [`CREATE MICROFLOW`],
    [Client Action], [Nanoflow], [`CREATE NANOFLOW`],
    [Web Block], [Snippet], [`CREATE SNIPPET`],
    [Static Entity], [Enumeration], [`CREATE ENUMERATION`],
  )]
  , kind: table
  )

#strong[AI Migration Workflow:]

#figure(image("diagrams/mermaid_14.png"),
  caption: [
    Diagram 14
  ]
)

#strong[Example: OutSystems Entity to MDL]

OutSystems Entity (conceptual):

```
Entity: Customer
  Id: Integer (Auto Number)
  Name: Text (200)
  Email: Text (200)
  IsActive: Boolean (Default: True)
  CreatedOn: DateTime (Default: CurrentDateTime)
```

Generated MDL:

```sql
CREATE PERSISTENT ENTITY Migration.Customer (
  Name: String(200) NOT NULL,
  Email: String(200),
  IsActive: Boolean DEFAULT true,
  CreatedOn: DateTime DEFAULT '[%CurrentDateTime%]'
);
```

=== From Java/.NET Business Applications
<from-java.net-business-applications>
Traditional 3GL applications require deeper analysis but follow patterns:

#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Java/.NET Pattern], [Mendix Equivalent], [MDL Approach],),
    table.hline(),
    [JPA Entity / EF Model], [Entity], [Extract from annotations/attributes],
    [Repository / DAO], [Built-in], [Mendix handles persistence],
    [Service Class], [Microflow], [Convert method → microflow logic],
    [Controller], [Page + Microflow], [Map endpoints to pages],
    [DTO], [Entity or Non-persistent Entity], [Depends on usage],
    [Validation], [Validation Microflow], [Convert validation rules],
  )]
  , kind: table
  )

#strong[AI Migration Workflow:]

#figure(image("diagrams/mermaid_16.png"),
  caption: [
    Diagram 16
  ]
)

#strong[Example: Spring Boot to MDL]

Java Entity:

```java
@Entity
@Table(name = "customers")
public class Customer {
    @Id @GeneratedValue
    private Long id;

    @Column(length = 200, nullable = false)
    private String name;

    @Column(length = 200)
    private String email;

    @OneToMany(mappedBy = "customer")
    private List<Order> orders;

    @Column(precision = 10, scale = 2)
    private BigDecimal creditLimit;
}
```

Generated MDL:

```sql
CREATE PERSISTENT ENTITY CRM.Customer (
  Name: String(200) NOT NULL,
  Email: String(200),
  CreditLimit: Decimal
);

CREATE ASSOCIATION CRM.Order_Customer
  BETWEEN CRM.Order AND CRM.Customer
  TYPE REFERENCE;
```

Java Service:

```java
@Service
public class CustomerService {
    public BigDecimal calculateDiscount(Order order) {
        Customer customer = order.getCustomer();
        if (customer.isPreferred()) {
            return order.getTotal().multiply(new BigDecimal("0.10"));
        }
        return BigDecimal.ZERO;
    }
}
```

Generated MDL:

```sql
CREATE MICROFLOW CRM.CalculateDiscount (
  $Order: CRM.Order
)
RETURNS Decimal
BEGIN
  RETRIEVE $Customer FROM $Order/Order_Customer;
  IF $Customer/IsPreferred THEN
    RETURN $Order/Total * 0.10;
  END IF;
  RETURN 0;
END;
```

=== From Database-Centric Applications
<from-database-centric-applications>
Oracle Forms, MS Access, and similar tools store logic in the database:

#figure(
  align(center)[#table(
    columns: 2,
    align: (auto,auto,),
    table.header([Legacy Pattern], [Migration Approach],),
    table.hline(),
    [Table], [`CREATE ENTITY` from DDL],
    [View], [`CREATE VIEW ENTITY` with OQL],
    [Stored Procedure], [`CREATE MICROFLOW` with equivalent logic],
    [Trigger], [Before/After commit microflow],
    [Form], [`CREATE PAGE` with data widgets],
    [Report], [Page with DataGrid or external reporting],
  )]
  , kind: table
  )

#strong[AI Migration Workflow:]

```sql
-- Input: Oracle DDL + PL/SQL
CREATE TABLE customers (
  customer_id NUMBER PRIMARY KEY,
  name VARCHAR2(200) NOT NULL,
  credit_limit NUMBER(10,2) DEFAULT 0
);

CREATE PROCEDURE apply_discount(p_order_id NUMBER) AS
  v_total NUMBER;
  v_discount NUMBER;
BEGIN
  SELECT total INTO v_total FROM orders WHERE order_id = p_order_id;
  v_discount := v_total * 0.1;
  UPDATE orders SET discount = v_discount WHERE order_id = p_order_id;
END;
```

Generated MDL:

```sql
-- Entity from table
CREATE PERSISTENT ENTITY Legacy.Customer (
  Name: String(200) NOT NULL,
  CreditLimit: Decimal DEFAULT 0
);

-- Microflow from stored procedure
CREATE MICROFLOW Legacy.ApplyDiscount ($Order: Legacy.Order)
BEGIN
  CHANGE $Order SET Discount = $Order/Total * 0.10;
  COMMIT $Order;
END;
```

== Migration Tooling in VS Code
<migration-tooling-in-vs-code>
The VS Code extension enhances the migration experience:

#figure(image("diagrams/mermaid_17.png"),
  caption: [
    Diagram 17
  ]
)

== Migration Benefits with AI + MDL
<migration-benefits-with-ai-mdl>
#figure(
  align(center)[#table(
    columns: (40.91%, 59.09%),
    align: (auto,auto,),
    table.header([Benefit], [Description],),
    table.hline(),
    [#strong[Speed];], [AI can analyze and generate thousands of lines of MDL in minutes],
    [#strong[Consistency];], [Same patterns applied across entire codebase],
    [#strong[Accuracy];], [Language Server validates every generated statement],
    [#strong[Iteration];], [Easy to refine: adjust prompt, regenerate, validate],
    [#strong[Review];], [Concise MDL is easier to review than equivalent Java/.NET],
    [#strong[Traceability];], [Can document source-to-target mapping],
    [#strong[Risk Reduction];], [Mendix platform guarantees reduce security/quality risks],
  )]
  , kind: table
  )

== Migration Phases
<migration-phases>
#figure(image("diagrams/mermaid_18.png"),
  caption: [
    Diagram 18
  ]
)

== ROI of AI-Assisted Migration
<roi-of-ai-assisted-migration>
#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Factor], [Traditional Migration], [AI + MDL Migration],),
    table.hline(),
    [#strong[Analysis Time];], [Weeks of manual review], [Hours with AI scanning],
    [#strong[Code Generation];], [Manual rewrite], [AI generates 70-80%],
    [#strong[Review Effort];], [Review all generated code], [Review concise MDL],
    [#strong[Iteration Speed];], [Days per change], [Minutes per change],
    [#strong[Consistency];], [Depends on team], [AI applies patterns uniformly],
    [#strong[Documentation];], [Often skipped], [AI can generate alongside],
    [#strong[Risk];], [High (manual errors)], [Lower (platform validation)],
    [#strong[Total Duration];], [12-24 months], [4-8 months],
    [#strong[Cost Reduction];], [Baseline], [50-70% reduction],
  )]
  , kind: table
  )

#horizontalrule

= Implementation Roadmap
<implementation-roadmap>
== Phase 1: Foundation
<phase-1-foundation>
#strong[Status];: ✅ Completed

- ☒ MDL Parser (ANTLR4-based, cross-language grammar)
- ☒ REPL with connection management
- ☒ Basic CRUD operations (Entity, Association, Enumeration)
- ☒ CATALOG queries for searching
- ☒ VS Code extension with syntax highlighting
- ☒ Language Server with diagnostics and symbols

== Phase 2: AI Readiness
<phase-2-ai-readiness>
#strong[Status];: ✅ Completed

- ☒ #strong[mxcli];: Single unified CLI for AI agents (Cobra-based)
  - `mxcli exec` - Execute MDL scripts
  - `mxcli check` - Syntax and semantic validation
  - `mxcli diff` - Compare script against project state
  - `mxcli diff-local` - Compare local changes against git (MPR v2)
  - `mxcli lint` - Extensible linting framework
  - `mxcli describe` - Describe entities, microflows, pages
  - `mxcli init` - Initialize Claude Code integration in Mendix projects
- ☒ #strong[Enhanced CATALOG];:
  - REFS table tracking cross-references (call, create, retrieve, show\_page, etc.)
  - Code search commands: `SHOW CALLERS`, `SHOW CALLEES`, `SHOW REFERENCES`, `SHOW IMPACT`, `SHOW CONTEXT`
  - Full-text `SEARCH` command
  - Widget discovery: `SHOW WIDGETS`
- ☒ #strong[Linting Framework];: Extensible rules with multiple output formats
  - MDL001 (NamingConvention) - PascalCase enforcement
  - MDL002 (EmptyMicroflow) - Detect empty microflows
  - MDL003 (DomainModelSize) - Entity count limits
  - MDL004 (ValidationFeedback) - Non-empty message check
  - SARIF output for CI/GitHub integration
- ☒ #strong[Documentation];: AI-friendly skills in `.claude/skills/`
  - `write-microflows.md` - Microflow syntax and patterns
  - `create-page.md` - V3 page/widget syntax
  - `overview-pages.md` - CRUD page patterns
  - `master-detail-pages.md` - Master-detail patterns
  - `generate-domain-model.md` - Entity/Association syntax
  - `check-syntax.md` - Pre-flight validation checklist
  - `debug-bson.md` - BSON debugging workflow
- ☒ #strong[Microflow Body Validation];: AST-level semantic checks
  - Return value requirements
  - Variable scope validation
  - VALIDATION FEEDBACK message checks (CE0091)
- ☐ #strong[Completion Provider];: Context-aware suggestions
- ☐ #strong[Code Actions];: Auto-fix common errors

== Phase 3: Agentic Integration
<phase-3-agentic-integration>
#strong[Status];: 🔄 In Progress

- ☒ #strong[Claude Code Integration];: `mxcli init` creates `.claude/` folder with skills and commands
- ☒ #strong[Semantic Search];: `SEARCH` command for full-text search across project
- ☒ #strong[Batch Operations];: `UPDATE WIDGETS` for bulk widget property updates (experimental)
- ☐ #strong[Claude Code MCP Server];: Native tool integration
- ☐ #strong[Cursor Extension];: IDE-specific features
- ☐ #strong[GitHub Copilot];: Custom instructions support

== Phase 4: Enterprise Features
<phase-4-enterprise-features>
#strong[Status];: 🔄 Planned

- ☐ #strong[Security Model in MDL];: Role-based access definition
- ☐ #strong[Testing Framework];: Automated test generation
- ☐ #strong[CI/CD Integration];: Pipeline tooling (SARIF output available)
- ☐ #strong[Multi-Project Support];: Cross-app references
- ☐ #strong[Governance Dashboard];: AI generation tracking

#horizontalrule

= Success Metrics
<success-metrics>
== Adoption Metrics
<adoption-metrics>
#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Metric], [Target], [Measurement],),
    table.hline(),
    [#strong[AI-Generated Apps];], [1000+ apps/month], [Platform telemetry],
    [#strong[Token Efficiency];], [5x vs TypeScript], [Benchmark tests],
    [#strong[Generation Accuracy];], [95%+ valid syntax], [LSP error rates],
    [#strong[User Satisfaction];], [NPS \> 50], [Survey],
  )]
  , kind: table
  )

== Quality Metrics
<quality-metrics>
#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Metric], [Target], [Measurement],),
    table.hline(),
    [#strong[Validation Coverage];], [100% of AI output], [Platform checks],
    [#strong[Security Incidents];], [0 AI-related], [Security team],
    [#strong[Production Readiness];], [80%+ first-time], [Deployment success],
    [#strong[Review Time];], [\<10 min avg], [User studies],
  )]
  , kind: table
  )

#horizontalrule

= Strategic Positioning: The Open AI-Ready Platform
<strategic-positioning-the-open-ai-ready-platform>
== Competitive Moat
<competitive-moat>
Mendix's approach to AI agent integration creates sustainable competitive advantages:

#figure(image("diagrams/mermaid_19.png"),
  caption: [
    Diagram 19
  ]
)

== 1. Open Low-Code Platform
<open-low-code-platform>
#strong[First mover advantage];: Mendix is the first low-code platform to offer full AI agent integration through open, non-proprietary interfaces.

#figure(
  align(center)[#table(
    columns: (51.52%, 48.48%),
    align: (auto,auto,),
    table.header([Openness Aspect], [Implementation],),
    table.hline(),
    [#strong[Open DSL];], [MDL grammar published, parsers available in multiple languages (ANTLR4)],
    [#strong[Open Tooling];], [mxcli is open-source, works with any terminal-capable agent],
    [#strong[Open Integration];], [Standard file formats, Git-compatible, CI/CD ready],
    [#strong[No Proprietary Protocol];], [No special MCP server required---standard file and CLI operations],
  )]
  , kind: table
  )

#strong[Contrast with competitors];: Other low-code platforms either lack textual representations entirely, or require proprietary integrations that lock customers into specific AI vendors.

== 2. Bring Your Own Agent (BYOA)
<bring-your-own-agent-byoa>
Enterprises can choose their preferred AI agent without being locked into a specific vendor:

#figure(
  align(center)[#table(
    columns: (20.59%, 55.88%, 23.53%),
    align: (auto,auto,auto,),
    table.header([Agent], [Integration Method], [Status],),
    table.hline(),
    [#strong[Claude Code];], [Native CLI + skills via `mxcli init`], [✅ Supported],
    [#strong[GitHub Copilot];], [Workspace + terminal access], [✅ Supported],
    [#strong[Cursor];], [Terminal + file editing], [✅ Supported],
    [#strong[Windsurf];], [Agentic flows + CLI], [✅ Supported],
    [#strong[Custom/Internal Agents];], [CLI + MDL grammar], [✅ Supported],
    [#strong[Future Agents];], [Any agent with file/CLI access], [✅ Ready],
  )]
  , kind: table
  )

#strong[Strategic value];: As AI agents evolve rapidly, enterprises are not locked into today's best option. They can switch agents as the market evolves, or use different agents for different tasks.

== 3. Built for True Collaboration
<built-for-true-collaboration>
The platform is designed for humans and AI to work together, not to replace one with the other:

#figure(image("diagrams/mermaid_20.png"),
  caption: [
    Diagram 20
  ]
)

#strong[Key design decisions];: - MDL diffs are human-scannable (unlike JSON or procedural code) - Validation happens at multiple layers (syntax → semantic → lint → platform) - Visual review in Studio Pro provides final human checkpoint - Atomic commits ensure human approval before changes take effect

== 4. Hours, Not Weeks
<hours-not-weeks>
Complex enterprise use cases become achievable in hours instead of weeks:

#figure(
  align(center)[#table(
    columns: (19.23%, 40.38%, 40.38%),
    align: (auto,auto,auto,),
    table.header([Use Case], [Traditional Approach], [With MDL + AI Agent],),
    table.hline(),
    [#strong[Legacy Migration];], [Months of manual analysis and rewrite], [AI analyzes source, generates MDL, human reviews],
    [#strong[Bulk Application Updates];], [Manual changes across 50+ apps], [Script generates MDL diffs, batch apply with review],
    [#strong[Monolith Decomposition];], [Weeks of refactoring], [AI proposes module boundaries, generates extraction plan],
    [#strong[Compliance Updates];], [Touch every entity/microflow manually], [Pattern-based bulk updates via MDL],
    [#strong[API Versioning];], [Manual endpoint updates], [AI generates new version, preserves old],
  )]
  , kind: table
  )

#strong[Example: Bulk Update Across Applications]

```bash
# Update all ComboBox widgets across 50 applications
for app in apps/*.mpr; do
  mxcli -p "$app" -c "UPDATE WIDGETS SET 'showLabel' = false WHERE WidgetType LIKE '%combobox%'"
done

# Time: ~30 minutes (including review)
# Traditional: ~2 weeks of manual work
```

== 5. Co-existence Strategy
<co-existence-strategy>
Best of both worlds---visual development and AI automation coexist:

#figure(
  align(center)[#table(
    columns: 3,
    align: (auto,auto,auto,),
    table.header([Workflow], [Tool], [User],),
    table.hline(),
    [#strong[Visual design];], [Studio Pro], [Business developers, UX designers],
    [#strong[Bulk operations];], [mxcli + AI], [DevOps, platform teams],
    [#strong[Feature generation];], [AI agent + MDL], [Full-stack developers],
    [#strong[Review & refinement];], [Studio Pro], [All stakeholders],
    [#strong[CI/CD automation];], [mxcli + scripts], [Platform engineering],
  )]
  , kind: table
  )

#strong[No forced migration];: Teams can adopt AI-assisted development gradually. Visual developers continue using Studio Pro. Power users leverage CLI and AI. Both work on the same projects seamlessly.

== Strategic Implications
<strategic-implications>
#figure(
  align(center)[#table(
    columns: (26.67%, 37.78%, 35.56%),
    align: (auto,auto,auto,),
    table.header([For Mendix], [For Enterprises], [For the Market],),
    table.hline(),
    [Capture AI-native developers], [Reduce development costs 50%+], [Set standard for AI + low-code],
    [Differentiate from competitors], [Future-proof platform choice], [Drive low-code evolution],
    [Enable new use cases (migration)], [Accelerate digital transformation], [Expand addressable market],
    [Build ecosystem (agents, skills)], [Maintain governance & control], [Create network effects],
  )]
  , kind: table
  )

#horizontalrule

= Conclusion
<conclusion>
Mendix is uniquely positioned to become the platform of choice for AI-generated business applications:

+ #strong[MDL provides token efficiency];: 5-10x more concise than traditional code
+ #strong[Platform guarantees reduce risk];: Security, validation, and governance built-in
+ #strong[Language Server enables AI];: Real-time feedback for self-correction
+ #strong[VS Code extension aids review];: Humans can validate AI output
+ #strong[Visual tools complete the loop];: Studio Pro for final review and deployment

By investing in MDL tooling for agentic IDEs, Mendix can capture the emerging market of AI-assisted enterprise application development while maintaining the trust and governance that enterprises require.

#horizontalrule

= Appendix: Competitive Analysis
<appendix-competitive-analysis>
== Alternative Approaches
<alternative-approaches>
#figure(
  align(center)[#table(
    columns: (45.45%, 27.27%, 27.27%),
    align: (auto,auto,auto,),
    table.header([Approach], [Pros], [Cons],),
    table.hline(),
    [#strong[Traditional Code (TypeScript/Python)];], [Flexible, large ecosystem], [Verbose, security risks, review burden],
    [#strong[No-Code Platforms];], [Visual, accessible], [Limited AI generation capability],
    [#strong[Other Low-Code (OutSystems, Appian)];], [Similar guarantees], [No textual DSL for AI],
    [#strong[Infrastructure as Code (Terraform)];], [Declarative], [Not application-focused],
  )]
  , kind: table
  )

== Mendix Differentiation
<mendix-differentiation>
+ #strong[Only low-code platform with mature textual DSL]
+ #strong[Bidirectional: Text ↔ Visual]
+ #strong[Enterprise-grade platform guarantees]
+ #strong[Existing ecosystem of developers]
+ #strong[Proven at scale (thousands of apps)]

#horizontalrule

= References
<references>
- #link("https://docs.anthropic.com/claude-code")[Claude Code Documentation]
- #link("https://microsoft.github.io/language-server-protocol/")[Language Server Protocol]
- #link("https://docs.mendix.com/apidocs-mxsdk/mxsdk/")[Mendix Model SDK]
- #link("../02-features/")[MDL Syntax Reference]
- #link("../../packages/vscode-mdl/")[VS Code MDL Extension]
- #link("../../cmd/mxcli/")[mxcli Commands] - CLI implementation
- #link("../../mdl/grammar/")[MDL Grammar (ANTLR4)] - Parser grammar definition
- #link("../../reference/mendix-repl/templates/.claude/skills/")[Skill Templates] - Built-in skill files
- #link("../../mdl/linter/rules/")[Linting Rules] - Built-in lint rules
