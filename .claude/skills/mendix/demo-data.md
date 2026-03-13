# Skill: Connect to Application Database and Generate Demo Data

## Purpose

Connect directly to the Mendix application's PostgreSQL database from a devcontainer
and insert demo data — bypassing the runtime. Covers reading DB settings, understanding
Mendix's internal ID system, and safely inserting rows with correct IDs and association links.

## When to Use This Skill

- User asks to seed or populate the database with test/demo data
- User needs data in the app before the UI is built
- User wants to inspect the database directly (schema, row counts, etc.)
- Bulk data import that is impractical through the Mendix UI

---

## Step 1: Get Database Settings

Use `mxcli` to read the project's configured database connection:

```bash
./mxcli -p <project>.mpr -c "SHOW SETTINGS;"
```

Example output:
```
| Configuration 'Default' | PostgreSql, localhost:5434, db=mxcli2-dev, http=8080 |
```

For full credentials (username, password):

```bash
./mxcli -p <project>.mpr -c "DESCRIBE SETTINGS;"
```

Example output:
```sql
ALTER SETTINGS CONFIGURATION 'Default'
  DatabaseType = 'PostgreSql',
  DatabaseUrl = 'localhost:5434',
  DatabaseName = 'mxcli2-dev',
  DatabaseUserName = 'mendix',
  DatabasePassword = 'mendix',
  HttpPortNumber = 8080;
```

---

## Step 2: Connect to the Database

### From a devcontainer on macOS

The Mendix app's `localhost` in the project settings refers to the **Mac host**, not the
devcontainer. Use `host.docker.internal` to reach it:

```bash
PGPASSWORD=mendix psql -h host.docker.internal -p 5434 -U mendix -d mxcli2-dev
```

### Useful psql commands

```bash
# List all tables
\dt

# Describe a table
\d tasklist$task

# Run a query and exit
PGPASSWORD=mendix psql -h host.docker.internal -p 5434 -U mendix -d mxcli2-dev \
  -c "SELECT * FROM \"tasklist\$task\" LIMIT 5;"
```

---

## Step 3: Understand the Mendix ID System

Every Mendix object has a `bigint` ID composed of two parts:

```
| bits 63–48  | bits 47–0          |
|  entity ID  |  sequence number   |
| (16 bits)   |  (48 bits)         |
```

**Formula:** `id = (short_id::bigint << 48) | sequence_number`

### Look up an entity's short_id and current sequence

```sql
SELECT e.entity_name, e.table_name, ei.short_id, ei.object_sequence,
       (ei.short_id::bigint << 48) AS id_base
FROM mendixsystem$entityidentifier ei
JOIN mendixsystem$entity e ON e.id = ei.id
WHERE e.entity_name = 'TaskList.Task';
```

Example result:
```
 entity_name  | table_name    | short_id | object_sequence | id_base
--------------+---------------+----------+-----------------+-------------------
 TaskList.Task| tasklist$task |       50 |              11 | 14073748835532800
```

### Decode an existing ID

```sql
SELECT id,
       to_hex(id::bigint)              AS hex_id,
       (id::bigint >> 48)              AS entity_short_id,
       id::bigint & x'ffffffffffff'::bigint AS sequence_num
FROM "tasklist$task";
```

### ID generation rules

- `object_sequence` is the **next available** sequence number for that entity
- After inserting N rows, **advance** `object_sequence` by N so the running runtime
  does not reuse those IDs
- IDs are entity-scoped: two entities can have the same sequence number but different
  `short_id`, giving different `id` values

---

## Step 4: Check Association Storage and Optimistic Locking

### Determine association storage mode

Query `mendixsystem$association` to see how each association is stored:

```sql
SELECT association_name, table_name, child_column_name, storage_format
FROM mendixsystem$association
WHERE table_name LIKE 'tasklist%';
```

Mendix stores associations in one of two ways, controlled by the project's
`AssocStorage` convention setting (check with `SHOW SETTINGS`):

#### Mode A — Column storage (`AssocStorage: Column`)

The FK is a regular column in the **owner** entity's table. No junction table exists.

```
tasklist$note
  id                   bigint  PK
  content              varchar
  tasklist$note_task   bigint  FK → tasklist$task.id   ← inline association column
  mxobjectversion      bigint                          ← optimistic lock version
```

Column naming convention: `{module}${associationname}` — all lowercase, `$` separator.

To insert a note linked to a task, simply set the FK column:

```sql
INSERT INTO "tasklist$note" (id, content, author, datecreated, "tasklist$note_task", mxobjectversion)
VALUES (16607023625928722, 'Note text', 'Alice', '2026-02-18 10:00:00', 14073748835532801, 1);
```

#### Mode B — Junction table storage

Mendix creates a separate join table. Both entity IDs are stored there.

```
tasklist$note_task
  tasklist$noteid  bigint  FK → tasklist$note.id   (unique — enforces one task per note)
  tasklist$taskid  bigint  FK → tasklist$task.id
```

Inspect with `\d "tasklist$note_task"`. Insert the entity row first, then the link:

```sql
INSERT INTO "tasklist$note" (id, content, author, datecreated) VALUES
  (16607023625928722, 'Note text', 'Alice', '2026-02-18 10:00:00');

INSERT INTO "tasklist$note_task" ("tasklist$noteid", "tasklist$taskid") VALUES
  (16607023625928722, 14073748835532801);
```

### Optimistic locking — `mxobjectversion`

When the project has optimistic locking enabled, every entity table gets an
`mxobjectversion bigint` column. The runtime:

- Initialises the column to `1` for all existing rows during schema sync
- Increments it by 1 on every `COMMIT`
- Rejects a save if the version in the DB doesn't match what the client loaded

**Always set `mxobjectversion = 1` when inserting rows directly.** Leaving it `NULL`
will cause the runtime to reject the object the first time a user saves it.

Check whether a table has the column:
```sql
SELECT column_name FROM information_schema.columns
WHERE table_name = 'tasklist$task' AND column_name = 'mxobjectversion';
```

---

## Step 5: Insert Demo Data

### Template — entity with column-storage association + optimistic locking

```sql
BEGIN;

-- Insert notes with inline FK and version column
INSERT INTO "tasklist$note" (id, content, author, datecreated, "tasklist$note_task", mxobjectversion)
VALUES
  (16607023625928722, 'First note content',  'Bob',   '2026-02-18 10:00:00', 14073748835532801, 1),
  (16607023625928723, 'Second note content', 'Alice', '2026-02-18 11:00:00', 14073748835532801, 1);

-- Advance Note sequence (was 18, inserted 2, now 20)
UPDATE mendixsystem$entityidentifier ei
SET object_sequence = 20
FROM mendixsystem$entity e
WHERE e.id = ei.id AND e.entity_name = 'TaskList.Note';

COMMIT;
```

### Template — entity with junction-table association (no optimistic locking)

```sql
BEGIN;

INSERT INTO "tasklist$note" (id, content, author, datecreated) VALUES
  (16607023625928722, 'First note content',  'Bob',   '2026-02-18 10:00:00'),
  (16607023625928723, 'Second note content', 'Alice', '2026-02-18 11:00:00');

INSERT INTO "tasklist$note_task" ("tasklist$noteid", "tasklist$taskid") VALUES
  (16607023625928722, 14073748835532801),
  (16607023625928723, 14073748835532801);

UPDATE mendixsystem$entityidentifier ei
SET object_sequence = 20
FROM mendixsystem$entity e
WHERE e.id = ei.id AND e.entity_name = 'TaskList.Note';

COMMIT;
```

### Template — standalone entity (no association)

```sql
BEGIN;

-- short_id=50, object_sequence=11 → id_base=14073748835532800
INSERT INTO "tasklist$task" (id, title, taskstatus, priority, assignedto, duedate, iscompleted, estimatedhours, mxobjectversion)
VALUES
  (14073748835532811, 'My demo task', 'ToDo', 'Medium', 'Alice', '2026-03-01 09:00:00', false, 4.0, 1);

-- Advance sequence (was 11, inserted 1 row)
UPDATE mendixsystem$entityidentifier ei
SET object_sequence = 12
FROM mendixsystem$entity e
WHERE e.id = ei.id AND e.entity_name = 'TaskList.Task';

COMMIT;
```

### Helper query — compute next N IDs for an entity

```sql
SELECT
  entity_name,
  short_id,
  object_sequence                             AS next_seq,
  (short_id::bigint << 48) + object_sequence AS first_new_id,
  (short_id::bigint << 48) + object_sequence + 9 AS last_new_id_if_10_rows
FROM mendixsystem$entityidentifier ei
JOIN mendixsystem$entity e ON e.id = ei.id
WHERE e.entity_name = 'TaskList.Note';
```

---

## Important Caveats

### Reserved attribute names

Mendix automatically adds system attributes to every entity. **Do not use these names**
for custom attributes — they will cause errors when the app tries to sync the schema:

| Reserved name | System meaning |
|---------------|----------------|
| `CreatedDate` | Auto-set on object creation |
| `ChangedDate` | Auto-set on every commit |
| `Owner`       | Reference to creating user |
| `ChangedBy`   | Reference to last user to commit |

If you need a "date created" field, name it `DateCreated`, `NoteDate`, etc.

### New entities need a runtime sync before demo data can be inserted

When you create a new entity with `mxcli exec`, the table and `mendixsystem$entity`
registration only appear **after the Mendix runtime starts and syncs the schema**.
The runtime does this automatically on startup. Until then:
- `\dt` will not show the table
- `mendixsystem$entityidentifier` will not have a row for the entity

Workflow:
1. Create entity with `mxcli exec`
2. Start (or restart) the Mendix runtime
3. Verify the table exists: `\dt *entityname*`
4. Insert demo data

### Sequence safety

Always update `object_sequence` in the same transaction as your inserts. If the runtime
is running concurrently, it may also allocate IDs from the same sequence. To be safe,
insert demo data while the runtime is stopped, or use a sequence value well above the
current `object_sequence` to leave headroom.

---

## Quick Reference

```bash
# Get DB settings
./mxcli -p <project>.mpr -c "DESCRIBE SETTINGS;"

# Connect (devcontainer on macOS)
PGPASSWORD=mendix psql -h host.docker.internal -p 5434 -U mendix -d mxcli2-dev

# Find entity short_id and id_base
SELECT e.entity_name, ei.short_id, ei.object_sequence,
       (ei.short_id::bigint << 48) AS id_base
FROM mendixsystem$entityidentifier ei
JOIN mendixsystem$entity e ON e.id = ei.id
WHERE e.entity_name = 'Module.Entity';

# ID formula
id = (short_id::bigint << 48) | sequence_number

# Check association storage mode
SELECT association_name, table_name, child_column_name
FROM mendixsystem$association
WHERE table_name LIKE 'mymodule%';

# Check if optimistic locking is enabled on a table
SELECT column_name FROM information_schema.columns
WHERE table_name = 'mymodule$myentity' AND column_name = 'mxobjectversion';

# After inserting N rows, advance the sequence
UPDATE mendixsystem$entityidentifier ei
SET object_sequence = <old_value + N>
FROM mendixsystem$entity e
WHERE e.id = ei.id AND e.entity_name = 'Module.Entity';
```

### INSERT column checklist

| Column | Required | Value |
|--------|----------|-------|
| `id` | Always | `(short_id::bigint << 48) \| sequence` |
| `mxobjectversion` | If column exists | `1` |
| `module$assocname` | If column-storage association | FK id of related object |
| Custom attributes | As needed | Your data |

---

## Automated Alternative: IMPORT FROM

For bulk imports from an external database, use the `IMPORT FROM` command instead of
writing manual INSERT statements. It handles ID generation, sequence updates, and
`mxobjectversion` automatically:

```sql
-- Connect to external database
SQL CONNECT postgres 'postgres://user:pass@host:5432/legacydb' AS source;

-- Import rows directly into Mendix app database
IMPORT FROM source QUERY 'SELECT name, email, department FROM employees'
  INTO HRModule.Employee
  MAP (name AS Name, email AS Email, department AS Department)
  BATCH 500;

-- Import with association linking (lookup by natural key)
IMPORT FROM source QUERY 'SELECT name, email, dept_name FROM employees'
  INTO HR.Employee
  MAP (name AS Name, email AS Email)
  LINK (dept_name TO Employee_Department ON Name);

-- Multiple associations
IMPORT FROM source QUERY 'SELECT name, dept, mgr_email FROM employees'
  INTO HR.Employee
  MAP (name AS Name)
  LINK (dept TO Employee_Department ON Name,
        mgr_email TO Employee_Manager ON Email);
```

The `IMPORT` command auto-connects to the Mendix app's PostgreSQL database using
project settings. Override with env vars for devcontainers/Docker:
`MXCLI_DB_TYPE`, `MXCLI_DB_HOST`, `MXCLI_DB_PORT`, `MXCLI_DB_NAME`,
`MXCLI_DB_USER`, `MXCLI_DB_PASSWORD`.

The `LINK` clause maps source columns to Mendix associations:
- `ON ChildAttr` — looks up the child entity by attribute value (builds a cache)
- Without `ON` — treats the source value as a raw Mendix object ID
- Handles both Column storage (inline FK) and Table storage (junction table) automatically
- Only Reference associations supported (not ReferenceSet)

Use manual INSERT (described above) when you need:
- ReferenceSet association linking
- Custom ID allocation or sequence management
- Non-standard data transformations

## Related Skills

- [database-connections.md](./database-connections.md) — Connecting to *external* databases from Mendix microflows
- [project-settings.md](./project-settings.md) — Reading and changing project configuration with `ALTER SETTINGS`
- [generate-domain-model.md](./generate-domain-model.md) — Creating entities before inserting data
