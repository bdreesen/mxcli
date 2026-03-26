# Playwright Testing

mxcli integrates with [playwright-cli](https://github.com/microsoft/playwright-cli) for automated UI testing of Mendix applications. This enables verification that generated pages render correctly, widgets appear in the DOM, and navigation works as expected.

## Why Browser Testing?

Mendix apps are React single-page applications. The server returns a JavaScript shell, and the actual UI is rendered client-side. An HTTP 200 from `/p/Customer_Overview` tells you nothing about whether the page's widgets actually rendered. A button defined in MDL might be missing from the DOM due to conditional visibility, incorrect container nesting, or a BSON serialization issue -- none of which are detectable without executing JavaScript in a real browser.

## Prerequisites

- **Node.js** -- Required for playwright-cli
- **playwright-cli** -- Install globally: `npm install -g @playwright/cli@latest`
- **Chromium** -- Install via: `playwright-cli install --with-deps chromium`
- **Running Mendix app** -- Start with `mxcli docker run -p app.mpr --wait`

### Devcontainer Setup

If using a devcontainer, add the following to `postCreateCommand` in your `devcontainer.json`:

```bash
npm install -g @playwright/cli@latest && playwright-cli install --with-deps chromium
```

The `mxcli init` command generates a `.playwright/cli.config.json` with sensible defaults for headless Chromium and a `PLAYWRIGHT_CLI_SESSION` environment variable for session management.

## mxcli playwright verify

The `mxcli playwright verify` command runs `.test.sh` scripts against a running Mendix application and collects results.

```bash
# Run all test scripts in a directory
mxcli playwright verify tests/ -p app.mpr

# Run a specific script
mxcli playwright verify tests/verify-customers.test.sh

# List discovered scripts without executing
mxcli playwright verify tests/ --list

# Output JUnit XML for CI integration
mxcli playwright verify tests/ -p app.mpr --junit results.xml

# Verbose output (show script stdout/stderr)
mxcli playwright verify tests/ -p app.mpr --verbose

# Custom app URL (auto-detected from .docker/.env by default)
mxcli playwright verify tests/ --base-url http://localhost:9090

# Skip the app health check
mxcli playwright verify tests/ --skip-health-check
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `--list`, `-l` | `false` | List test scripts without executing |
| `--junit`, `-j` | | Write JUnit XML results to file |
| `--verbose`, `-v` | `false` | Show script stdout/stderr during execution |
| `--color` | `false` | Use colored terminal output |
| `--timeout`, `-t` | `2m` | Timeout per script execution |
| `--base-url` | `http://localhost:8080` | Mendix app base URL |
| `--skip-health-check` | `false` | Skip app reachability check before running |
| `-p` | | Path to the `.mpr` project file |

### How It Works

The verify runner performs these steps:

1. **Discovers** all `.test.sh` files in the provided paths
2. **Checks** that `playwright-cli` is available in PATH
3. **Health-checks** the app at the base URL (unless `--skip-health-check`)
4. **Opens** a playwright-cli browser session (Chromium by default)
5. **Runs** each `.test.sh` script sequentially via `bash`
6. **Captures** a screenshot on failure for debugging
7. **Closes** the browser session
8. **Reports** pass/fail per script with timing
9. **Writes** JUnit XML if `--junit` is specified
10. **Exits** with non-zero status if any script failed

## Test Script Format

Test scripts are plain bash files using playwright-cli commands. The naming convention is `tests/verify-<name>.test.sh`. Scripts should use `set -euo pipefail` so that any failing command causes the script to exit with a non-zero code.

### Example Test Script

```bash
#!/usr/bin/env bash
# tests/verify-customers.test.sh -- Customer module smoke test
set -euo pipefail

# --- Login ---
playwright-cli open http://localhost:8080
playwright-cli fill e12 "MxAdmin"
playwright-cli fill e15 "AdminPassword1!"
playwright-cli click e18
playwright-cli state-save mendix-auth

# --- Customer Overview page ---
playwright-cli goto http://localhost:8080/p/Customer_Overview
playwright-cli run-code "document.querySelector('.mx-name-dgCustomers') !== null"
playwright-cli run-code "document.querySelector('.mx-name-btnNew') !== null"
playwright-cli run-code "document.querySelector('.mx-name-btnEdit') !== null"

# --- Create a customer ---
playwright-cli click btnNew
playwright-cli snapshot
playwright-cli fill txtName "CI Test Customer"
playwright-cli fill txtEmail "ci@test.com"
playwright-cli click btnSave

# --- Verify persistence (via OQL, not direct DB query) ---
mxcli oql -p app.mpr --json \
  "SELECT Name FROM MyModule.Customer WHERE Name = 'CI Test Customer'" \
  | grep -q "CI Test Customer"

# --- Cleanup ---
playwright-cli close
echo "PASS: verify-customers"
```

## Widget Selectors

Mendix renders each widget's name property as a CSS class on the corresponding DOM element:

```html
<div class="mx-name-submitButton form-group">
```

These `.mx-name-*` classes are stable, predictable, and directly derived from MDL widget names. This makes them ideal for test assertions.

### Using Selectors in Scripts

The recommended approach for CI scripts is to use `run-code` with CSS selectors, since dynamic element refs (`e12`, `e15`) change between page loads:

```bash
# Check widget presence
playwright-cli run-code "document.querySelector('.mx-name-btnSave') !== null"

# Click a widget
playwright-cli run-code "document.querySelector('.mx-name-btnSave').click()"

# Check text content
playwright-cli run-code "document.querySelector('.mx-name-lblTitle').textContent.includes('Customers')"
```

For assertions that should cause failures, use `throw` so the script exits non-zero under `set -e`:

```bash
playwright-cli run-code "if (!document.querySelector('.mx-name-btnSave')) throw new Error('btnSave not found')"
```

## Three Test Layers

### Layer 1: Smoke Tests

Fast checks that pages are reachable and the app starts without errors. These run first as a gate before heavier tests.

```bash
# HTTP reachability
playwright-cli open http://localhost:8080
playwright-cli goto http://localhost:8080/p/Customer_Overview
playwright-cli goto http://localhost:8080/p/Customer_Edit
```

### Layer 2: UI Widget Tests

Verify that every widget generated in MDL is present and interactive in the DOM.

```bash
# Widget presence
playwright-cli run-code "document.querySelector('.mx-name-dgCustomers') !== null"
playwright-cli run-code "document.querySelector('.mx-name-btnNew') !== null"

# Form interaction
playwright-cli run-code "document.querySelector('.mx-name-txtName input').value = 'Test'"
playwright-cli run-code "document.querySelector('.mx-name-btnSave').click()"
```

### Layer 3: Data Assertions

Verify that UI interactions persist the correct data. Use `mxcli oql` for data validation instead of direct database queries.

```bash
# Submit a form via the UI, then verify persistence
mxcli oql -p app.mpr --json \
  "SELECT Name, Email FROM MyModule.Customer WHERE Name = 'Test'" \
  | grep -q "Test"
```

## Example Output

```
Found 3 test script(s)
Checking app at http://localhost:8080...
  App is reachable
Opening browser session (chromium)...
  [1/3] verify-login... PASS (1.2s)
  [2/3] verify-customers... PASS (3.4s)
  [3/3] verify-orders... FAIL (2.1s)
         Screenshot saved: verify-orders-failure.png
         Error: btnSubmit not found

Playwright Verify Results
============================================================
  PASS  verify-login (1.2s)
  PASS  verify-customers (3.4s)
  FAIL  verify-orders (2.1s)
         Error: btnSubmit not found
------------------------------------------------------------
Total: 3  Passed: 2  Failed: 1  Time: 6.7s
Some scripts failed.
```

## CI/CD Integration

For continuous integration, combine `mxcli docker run --wait` with `mxcli playwright verify`:

```bash
# Build and start the app
mxcli docker run -p app.mpr --wait

# Run all verification scripts
mxcli playwright verify tests/ -p app.mpr --junit results.xml

# JUnit XML output is consumed by CI systems (GitHub Actions, Jenkins, etc.)
```

## Related Pages

- [Testing](testing.md) -- MDL test framework (`mxcli test`)
- [Docker Run](docker-run.md) -- Building and running the Mendix app in Docker
- [Docker Check](docker-check.md) -- Validating projects without building
