# Proposal: Concurrent Access Safety for mxcli

## Problem

When Claude Code spawns multiple subagents, each may execute `mxcli` commands against the same Mendix project simultaneously. The current architecture assumes single-process access to all project files, which creates data corruption risks under concurrent use.

### Affected Resources

| Resource | Path | Read Safe | Write Safe | Failure Mode |
|----------|------|-----------|------------|--------------|
| MPR database | `app.mpr` | Yes (SQLite) | No | SQLite lock contention, potential corruption if two writers commit |
| mxunit files | `mprcontents/**/*.mxunit` | Yes | No | Silent data corruption, last writer wins |
| Catalog | `.mxcli/catalog.db` | Yes (SQLite) | No | "database is locked" errors during rebuild |
| Session logs | `~/.mxcli/logs/*.log` | N/A | Yes | Safe (O_APPEND, atomic small writes) |
| History file | `~/.mxcli_history` | N/A | Yes | Benign (readline append) |

### Concrete Scenario

Claude Code runs three subagents in parallel:

```
Agent 1: mxcli -p app.mpr -c "CREATE ENTITY Shop.Product (...); COMMIT;"
Agent 2: mxcli -p app.mpr -c "CREATE MICROFLOW Shop.ACT_Process (...); COMMIT;"
Agent 3: mxcli -p app.mpr -c "DESCRIBE ENTITY Shop.Customer"
```

Agent 3 (read-only) is fine. Agents 1 and 2 both open the MPR for writing, create objects, and commit. For MPR v2, they write to separate `.mxunit` files (different document types), so they *might* succeed. But both also update the SQLite `_units` table simultaneously, risking lock contention or inconsistent state. For MPR v1 (single SQLite file), concurrent writes are guaranteed to conflict.

## Analysis

### What's Actually Safe Today

**Concurrent reads** are safe:
- `SHOW`, `DESCRIBE`, `SEARCH`, `SELECT` commands
- Multiple `mxcli describe` or `mxcli show` processes
- SQLite handles concurrent readers natively

**Sequential writes** are safe:
- One process writes, finishes, then another starts
- Claude Code running commands one at a time

### What's Unsafe

**Concurrent writes** to the same project:
- Two `CREATE`/`DROP`/`COMMIT` commands running simultaneously
- `REFRESH CATALOG` running while another process also rebuilds
- Any write during an active `COMMIT` from another process

**Mixed read-write** during active writes:
- A reader may see partially committed state
- Not dangerous (no corruption) but potentially confusing results

## Proposed Solution: Project-Level File Lock

### Design

Use an advisory file lock (`flock()`) on a lock file next to the MPR to coordinate access between concurrent mxcli processes.

```
app.mpr
app.mpr.mxcli-lock    <-- advisory lock file
mprcontents/
.mxcli/
  catalog.db
  catalog.db.mxcli-lock  <-- separate lock for catalog
```

### Lock Modes

| Operation | Lock Type | Behavior |
|-----------|-----------|----------|
| Read commands (SHOW, DESCRIBE, SELECT) | Shared (LOCK_SH) | Multiple readers allowed concurrently |
| Write commands (CREATE, DROP, COMMIT) | Exclusive (LOCK_EX) | Blocks until other readers/writers finish |
| Catalog rebuild (REFRESH CATALOG) | Exclusive on catalog lock | Doesn't block MPR reads |
| Catalog query (SELECT FROM CATALOG) | Shared on catalog lock | Multiple concurrent queries OK |

### Lock Granularity Options

**Option A: Coarse-grained (recommended for v1)**

Single lock per project, acquired at `CONNECT` time based on read/write mode:

```go
// In executor.go, when opening project
func (e *Executor) execConnect(stmt *ast.ConnectStmt) error {
    // ... existing connect logic ...

    lockPath := mprPath + ".mxcli-lock"
    if e.needsWrite {
        e.lock = acquireExclusive(lockPath)  // blocks until available
    } else {
        e.lock = acquireShared(lockPath)     // concurrent with other readers
    }
    // lock released in Close() or DISCONNECT
}
```

Pros: Simple, correct, prevents all concurrent write issues.
Cons: Write commands block all other access (even reads) for the duration of the session.

**Option B: Fine-grained (future)**

Lock acquired per-statement, only around actual write operations:

```go
func (e *Executor) Execute(stmt ast.Statement) error {
    if isWriteStatement(stmt) {
        e.lock.Upgrade()       // shared -> exclusive
        defer e.lock.Downgrade() // back to shared
    }
    return e.executeInner(stmt)
}
```

Pros: Reads never blocked, writes only block briefly.
Cons: More complex, risk of deadlocks if upgrade fails, requires careful handling of multi-statement scripts.

### Lock Timeout

Waiting indefinitely for a lock is a poor UX. Add a configurable timeout:

```go
const defaultLockTimeout = 30 * time.Second

func acquireExclusive(lockPath string, timeout time.Duration) (*Lock, error) {
    // Try non-blocking first
    if tryFlock(fd, LOCK_EX|LOCK_NB) {
        return &Lock{fd: fd}, nil
    }
    // Log that we're waiting
    log.Info("waiting for project lock", "lock_path", lockPath)
    // Retry with timeout
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        if tryFlock(fd, LOCK_EX|LOCK_NB) {
            return &Lock{fd: fd}, nil
        }
        time.Sleep(100 * time.Millisecond)
    }
    return nil, fmt.Errorf("timeout waiting for project lock (another mxcli may be writing to this project)")
}
```

### Implementation

#### Files to Create

**`sdk/mpr/lock.go`** (~80 lines)

```go
package mpr

type Lock struct {
    fd   *os.File
    path string
}

func AcquireShared(lockPath string, timeout time.Duration) (*Lock, error)
func AcquireExclusive(lockPath string, timeout time.Duration) (*Lock, error)
func (l *Lock) Release() error
func (l *Lock) Upgrade() error    // shared -> exclusive (Option B only)
func (l *Lock) Downgrade() error  // exclusive -> shared (Option B only)
```

Uses `syscall.Flock()` on Linux/macOS. On Windows, uses `LockFileEx()`.

#### Files to Modify

**`sdk/mpr/reader.go`** — Acquire shared lock on Open
**`sdk/mpr/writer.go`** — Acquire exclusive lock on OpenForWriting
**`mdl/executor/executor.go`** — Release lock on Close/Disconnect

### Catalog Lock

Separate lock file for the catalog database:

```
.mxcli/catalog.db.mxcli-lock
```

- `REFRESH CATALOG` acquires exclusive
- `SELECT FROM CATALOG.*` acquires shared
- Independent from MPR lock (catalog reads shouldn't block project writes)

## Alternative Approaches Considered

### A. Document the Constraint

Simply document that concurrent write access is unsupported and rely on Claude Code to serialize write operations.

```markdown
## Concurrent Access
mxcli does not support concurrent write access to the same project.
When using with Claude Code, ensure write commands run sequentially.
Read commands (SHOW, DESCRIBE, SEARCH) can safely run in parallel.
```

Pros: Zero implementation effort.
Cons: Users will hit this, error messages will be confusing, potential data corruption.

### B. Per-Document Locks (mxunit level)

Lock individual `.mxunit` files instead of the whole project.

Pros: Maximum concurrency for MPR v2.
Cons: Complex, doesn't help MPR v1, `_units` table still needs project-level lock.

### C. Single mxcli Server Process

Run mxcli as a persistent server (like the LSP) that serializes all requests internally.

```bash
# Start server
mxcli serve -p app.mpr --socket /tmp/mxcli.sock

# Clients send commands to server
mxcli -S /tmp/mxcli.sock -c "SHOW ENTITIES"
```

Pros: Perfect serialization, shared catalog, single connection overhead.
Cons: Major architectural change, process lifecycle management, more failure modes.

## Recommendation

1. **Immediate**: Document the constraint in CLAUDE.md and skill files. Add a warning when `MXCLI_LOG` shows concurrent sessions on the same project.

2. **Short-term (Option A)**: Implement coarse-grained file locking. Simple, correct, covers the common case. A write session blocks other writers but the lock duration is short (mxcli commands complete in seconds).

3. **Long-term (Option C)**: Consider the server model if Claude Code usage patterns show frequent concurrent access needs. The LSP server (`mxcli lsp`) already demonstrates this architecture.

## Verification

```bash
# Test 1: Concurrent reads (should both succeed)
mxcli -p app.mpr -c "SHOW ENTITIES" &
mxcli -p app.mpr -c "SHOW MICROFLOWS" &
wait

# Test 2: Concurrent writes (second should wait for first)
mxcli -p app.mpr -c "CREATE MODULE Test1; COMMIT;" &
mxcli -p app.mpr -c "CREATE MODULE Test2; COMMIT;" &
wait
# Both should succeed (serialized by lock)

# Test 3: Lock timeout
mxcli -p app.mpr  # start REPL (holds lock)
# In another terminal:
mxcli -p app.mpr -c "CREATE MODULE Test3; COMMIT;"
# Should show "waiting for project lock..." then succeed or timeout

# Test 4: Read during write (should succeed)
mxcli -p app.mpr  # start REPL, run CREATE
# In another terminal:
mxcli -p app.mpr -c "SHOW ENTITIES"
# Should succeed (shared read lock)
```
