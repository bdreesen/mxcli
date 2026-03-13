# Migrating Non-Mendix Applications with mxcli & MDL

This guide describes how to use **mxcli** and **MDL** (Mendix Definition Language) to migrate existing non-Mendix applications to the Mendix platform. The process follows five phases: assess, propose, generate, test, and finish.

## Overview

**mxcli** is a command-line tool that reads and modifies Mendix projects (`.mpr` files) using **MDL**, a SQL-like domain-specific language. Together with an AI assistant (such as Claude), mxcli enables rapid, repeatable migration of applications from any technology stack to Mendix.

Key capabilities for migration:

| Capability | Description |
|------------|-------------|
| **MDL scripting** | SQL-like syntax for creating entities, microflows, pages, security, and more |
| **Migration skills** | Platform-specific guides for OutSystems, Oracle Forms, K2/Nintex, and generic assessments |
| **Automated validation** | Syntax checking, linting (41 rules), and best-practices scoring |
| **Docker testing** | Build, run, and test Mendix apps without Studio Pro |
| **Quality reports** | Scored assessment across naming, security, architecture, performance, and design |

## The Five Phases

```
  Phase 1          Phase 2           Phase 3           Phase 4         Phase 5
 ┌──────────┐   ┌───────────┐   ┌──────────────┐   ┌──────────┐   ┌──────────┐
 │  Assess  │──▶│  Propose  │──▶│   Generate   │──▶│   Test   │──▶│  Finish  │
 │  Source   │   │ Transform │   │  Mendix App  │   │  & Lint  │   │ in Studio│
 └──────────┘   └───────────┘   └──────────────┘   └──────────┘   │   Pro    │
                                                                    └──────────┘
```

---

## Phase 1: Assess the Existing Application

The first step is a thorough investigation of the source application to produce a structured migration inventory. The **assess-migration** skill provides a 6-step framework that works with any technology stack.

### What to Extract

| Category | What to Document | Examples |
|----------|-----------------|----------|
| **Data model** | Entities, attributes, relationships, constraints | JPA entities, Django models, DB schema |
| **Business logic** | Validation rules, calculations, workflows | Service classes, stored procedures, triggers |
| **Pages / UI** | Screens, forms, dashboards, navigation | React components, Razor views, JSP templates |
| **Integrations** | APIs consumed/exposed, file feeds, message queues | REST clients, SOAP services, Kafka topics |
| **Security** | Authentication methods, roles, data access rules | Spring Security, RBAC policies, row-level security |
| **Scheduled jobs** | Background tasks, timers, batch processing | Cron jobs, Quartz schedulers, Celery tasks |

### Assessment Output

The assessment produces a structured report with:

- **Executive summary** — technology stack, application size, complexity rating
- **Categorized inventory** — counts and details for each category above
- **Mendix mapping** — how each source element maps to Mendix concepts
- **Migration risks** — complex stored procedures, custom UI components, real-time integrations
- **Recommended phases** — suggested order of migration work

### Platform-Specific Skills

For common source platforms, dedicated migration skills provide deeper guidance:

| Source Platform | Skill | Key Mappings |
|----------------|-------|-------------|
| **OutSystems** | `migrate-outsystems` | eSpace → Module, Entity → Entity, Action → Microflow, Screen → Page, Static Entity → Enumeration |
| **Oracle Forms** | `migrate-oracle-forms` | Form → Page, Block → Snippet, PL/SQL → Microflow, LOV → Enumeration |
| **K2 / Nintex** | `migrate-k2-nintex` | SmartForm → Page, SmartObject → Entity, Workflow → Microflow chain |
| **Any stack** | `assess-migration` | Generic framework for Java, .NET, Python, Node.js, PHP, Ruby, etc. |

### Example: Starting an Assessment

```bash
# Point the AI assistant at the source codebase
# The assess-migration skill guides investigation of:
# - Build files (pom.xml, package.json, requirements.txt)
# - ORM models (@Entity, DbContext, models.py)
# - Service classes and business logic
# - UI templates and route definitions
# - Security configuration
# - API clients and integrations
```

The assessment report becomes the input for Phase 2.

---

## Phase 2: Create a Transformation Proposal

With the assessment complete, create a detailed proposal mapping every source element to its Mendix equivalent.

### Module Structure

Plan the Mendix module structure. A common pattern is to consolidate many small source modules into fewer Mendix modules:

```
Source (OutSystems)          Mendix
─────────────────           ─────────
CRM_Core, CRM_UI    ──▶    CRM
Cases_Core, Cases_UI ──▶    Cases
Auth_Core, Auth_SSO  ──▶    Administration
Shared_Utils         ──▶    Commons
```

### Transformation Mapping

For each element in the assessment, document the target Mendix implementation:

| Source Element | Source Location | Mendix Target | MDL Statement |
|---------------|----------------|---------------|---------------|
| Customer table | `schema.sql` | `CRM.Customer` entity | `CREATE PERSISTENT ENTITY` |
| OrderStatus enum | `enums.py` | `Sales.OrderStatus` enumeration | `CREATE ENUMERATION` |
| calculateTotal() | `OrderService.java` | `Sales.ACT_Order_CalculateTotal` microflow | `CREATE MICROFLOW` |
| Customer list page | `customers.html` | `CRM.Customer_Overview` page | `CREATE PAGE` |
| Admin role | `security.xml` | `Administration.Admin` module role | `GRANT` statements |
| Stripe integration | `PaymentClient.js` | REST consumed service | `CREATE ODATA CLIENT` or REST call |

### Prioritization

Order the migration work to maximize early value:

1. **Enumerations** — no dependencies, used by entities
2. **Domain model** — entities, attributes, associations
3. **Security** — module roles, user roles, access rules
4. **Core business logic** — validation microflows, calculation microflows
5. **Pages** — overview pages, edit forms, dashboards
6. **Integrations** — REST clients, OData services, file handling
7. **Navigation** — menu structure, home pages
8. **Advanced features** — scheduled events, workflows, business events

---

## Phase 3: Generate the Mendix Application

Use mxcli and MDL scripts to generate the Mendix application. The AI assistant uses **skills** — structured guides for best practices — to produce correct, idiomatic MDL.

### Key Skills for Generation

| Skill | Purpose |
|-------|---------|
| `generate-domain-model` | Entity, association, and enumeration syntax with naming conventions |
| `write-microflows` | Microflow syntax, 60+ activity types, common patterns |
| `create-page` | Page and widget syntax for 50+ widget types |
| `overview-pages` | CRUD page patterns (list + detail) |
| `master-detail-pages` | Master-detail page layouts |
| `manage-security` | Module roles, user roles, GRANT/REVOKE, demo users |
| `manage-navigation` | Navigation profiles, menu items, home pages |
| `organize-project` | Folder structure, MOVE command, project conventions |
| `demo-data` | Seed data via external database import |
| `database-connections` | External database connectivity from microflows |

### Generation Workflow

```bash
# 1. Create a new Mendix project in Studio Pro (or use an existing one)
# 2. Connect mxcli to the project
mxcli -p /path/to/app.mpr

# 3. Execute MDL scripts
mxcli exec domain-model.mdl -p app.mpr
mxcli exec microflows.mdl -p app.mpr
mxcli exec pages.mdl -p app.mpr
mxcli exec security.mdl -p app.mpr

# 4. Or run interactively in the REPL
mxcli -p app.mpr
```

### Example: Domain Model Generation

```sql
-- domain-model.mdl
CONNECT LOCAL './MyApp.mpr';

-- Enumerations first (referenced by entities)
CREATE ENUMERATION Sales.OrderStatus (
  Draft 'Draft',
  Pending 'Pending',
  Confirmed 'Confirmed',
  Shipped 'Shipped',
  Delivered 'Delivered',
  Cancelled 'Cancelled'
);

-- Entities
/** Customer master data */
@Position(100, 100)
CREATE PERSISTENT ENTITY CRM.Customer (
  Name: String(200) NOT NULL ERROR 'Customer name is required',
  Email: String(200) UNIQUE ERROR 'Email already exists',
  Phone: String(50),
  IsActive: Boolean DEFAULT TRUE
)
INDEX (Name)
INDEX (Email);
/

/** Sales order */
@Position(300, 100)
CREATE PERSISTENT ENTITY Sales.Order (
  OrderNumber: String(50) NOT NULL UNIQUE,
  OrderDate: DateTime NOT NULL,
  TotalAmount: Decimal DEFAULT 0,
  Status: Enumeration(Sales.OrderStatus) DEFAULT 'Draft'
)
INDEX (OrderNumber)
INDEX (OrderDate DESC);
/

-- Associations
CREATE ASSOCIATION Sales.Order_Customer
  FROM CRM.Customer
  TO Sales.Order
  TYPE Reference
  OWNER Default;
/

-- Security
GRANT Sales.Admin ON Sales.Order (CREATE, DELETE, READ *, WRITE *);
GRANT Sales.User ON Sales.Order (CREATE, READ *, WRITE *)
  WHERE '[Status != ''Cancelled'']';

COMMIT MESSAGE 'Generated Sales domain model';
```

### Example: Microflow Generation

```sql
-- microflows.mdl
CONNECT LOCAL './MyApp.mpr';

CREATE MICROFLOW Sales.ACT_Order_CalculateTotal
BEGIN
  DECLARE $Order Sales.Order;
  RETRIEVE $Lines FROM Sales.OrderLine WHERE [Sales.OrderLine_Order = $Order];
  DECLARE $Total Decimal = 0;
  LOOP $Line IN $Lines
  BEGIN
    SET $Total = $Total + $Line/Price * $Line/Quantity;
  END;
  CHANGE $Order (TotalAmount = $Total);
  COMMIT $Order;
END;
/

COMMIT MESSAGE 'Generated order calculation microflow';
```

### Example: Page Generation

```sql
-- pages.mdl
CONNECT LOCAL './MyApp.mpr';

CREATE PAGE CRM.Customer_Overview (
  Title: 'Customers',
  Layout: Atlas_Core.Atlas_Default
) {
  DATAGRID2 ON CRM.Customer (
    COLUMN Name { Caption: 'Name' }
    COLUMN Email { Caption: 'Email' }
    COLUMN Phone { Caption: 'Phone' }
    COLUMN IsActive { Caption: 'Active' }
    SEARCH ON Name, Email
    BUTTON 'New' CALL CRM.Customer_NewEdit
    BUTTON 'Edit' CALL CRM.Customer_NewEdit
    BUTTON 'Delete' CALL CONFIRM DELETE
  )
};
/

COMMIT MESSAGE 'Generated customer overview page';
```

---

## Phase 4: Test and Validate

mxcli provides a comprehensive validation pipeline — from syntax checking to running the application in Docker.

### Validation Steps

```bash
# Step 1: Syntax check (fast, no project needed)
mxcli check script.mdl

# Step 2: Reference validation (checks entity/microflow names exist)
mxcli check script.mdl -p app.mpr --references

# Step 3: Lint the project (41 built-in rules + 27 Starlark rules)
mxcli lint -p app.mpr

# Step 4: Quality report (scored 0-100 per category)
mxcli report -p app.mpr --format markdown

# Step 5: Mendix compiler check (requires Docker)
mxcli docker check -p app.mpr

# Step 6: Build and run (requires Docker)
mxcli docker build -p app.mpr
mxcli docker run -p app.mpr
```

### Lint Categories

The linter covers 6 categories with rules in the MDL, SEC, QUAL, ARCH, DESIGN, and CONV series:

| Category | Focus | Example Rules |
|----------|-------|---------------|
| **Naming** (CONV) | Naming conventions | Entity naming, microflow prefixes |
| **Security** (SEC) | Access control | Missing access rules, open security |
| **Quality** (QUAL) | Code quality | Unused variables, empty microflows |
| **Architecture** (ARCH) | Structure | Module dependencies, circular references |
| **Performance** (DESIGN) | Efficiency | Missing indexes, large retrieve-all |
| **Design** (MDL) | Best practices | Entity design, association patterns |

### Automated Testing

Write `.test.mdl` or `.test.md` test files to validate microflow behavior:

```sql
-- tests/order_tests.test.mdl
-- @test: Order total is calculated correctly
CONNECT LOCAL './MyApp.mpr';

$Customer = CREATE CRM.Customer (Name = 'Test');
COMMIT $Customer;

$Order = CREATE Sales.Order (OrderNumber = 'ORD-001', OrderDate = now());
CHANGE $Order (Sales.Order_Customer = $Customer);
COMMIT $Order;

$Line = CREATE Sales.OrderLine (Price = 10.00, Quantity = 3);
CHANGE $Line (Sales.OrderLine_Order = $Order);
COMMIT $Line;

CALL MICROFLOW Sales.ACT_Order_CalculateTotal ($Order = $Order);

-- @assert: $Order/TotalAmount = 30.00
```

Run tests with:

```bash
# Run all tests (requires Docker)
mxcli test tests/ -p app.mpr

# Output JUnit XML for CI/CD integration
mxcli test tests/ -p app.mpr --format junit
```

### Script Impact Preview

Before applying changes, preview what a script will modify:

```bash
# Compare script against current project state
mxcli diff -p app.mpr changes.mdl
```

---

## Phase 5: Review and Finish in Studio Pro

The final phase transitions from mxcli/MDL to Mendix Studio Pro for visual refinement and features that require the IDE.

### What to Do in Studio Pro

| Task | Why Studio Pro |
|------|---------------|
| **Page layout tuning** | Visual drag-and-drop for pixel-perfect layouts |
| **Styling and theming** | CSS/SCSS editing with live preview |
| **Complex workflows** | Workflow editor for multi-step approval processes |
| **Custom widgets** | Pluggable widget development and configuration |
| **Marketplace modules** | Install and configure marketplace modules |
| **Integration testing** | Test REST/OData services with real endpoints |
| **Performance profiling** | Runtime profiling and query optimization |
| **Deployment** | Configure deployment pipelines and environments |

### Handoff Checklist

Before opening in Studio Pro, verify:

```bash
# Final validation
mxcli docker check -p app.mpr          # No compiler errors
mxcli lint -p app.mpr                  # No critical lint issues
mxcli report -p app.mpr               # Review quality scores
mxcli test tests/ -p app.mpr          # All tests pass

# Review project structure
mxcli -p app.mpr -c "SHOW STRUCTURE DEPTH 2"
```

### Iterative Workflow

The mxcli/Studio Pro workflow is iterative — you can switch between them:

```
  mxcli/MDL                    Studio Pro
 ┌──────────┐               ┌──────────────┐
 │ Generate  │──── open ───▶│ Visual edit   │
 │ entities, │               │ styling,     │
 │ microflows│◀── save ─────│ test, deploy  │
 │ pages     │               │              │
 └──────────┘               └──────────────┘
```

Changes made in either tool are persisted in the `.mpr` file. Always close mxcli before opening in Studio Pro, and vice versa.

---

## Migration Workflow Summary

| Phase | Tools | Skills Used | Output |
|-------|-------|-------------|--------|
| 1. Assess | AI assistant + source code | `assess-migration`, platform-specific skills | Assessment report |
| 2. Propose | AI assistant | Assessment report as input | Transformation mapping |
| 3. Generate | mxcli, MDL scripts | `generate-domain-model`, `write-microflows`, `create-page`, `manage-security`, `organize-project` | Working Mendix project |
| 4. Test | mxcli, Docker | `check-syntax`, `assess-quality` | Validated, lint-clean project |
| 5. Finish | Studio Pro | — | Production-ready application |

## Further Reading

- [mxcli Overview](./mxcli-overview.md) — comprehensive mxcli feature reference
- [MDL Quick Reference](../MDL_QUICK_REFERENCE.md) — complete MDL syntax tables
- [MDL Domain Model](../05-mdl-specification/03-domain-model.md) — entity, attribute, and association specification
- [MDL Language Reference](../05-mdl-specification/01-language-reference.md) — full language specification
