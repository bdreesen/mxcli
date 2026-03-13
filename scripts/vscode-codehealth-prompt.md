# Project Prompt: VS Code Code Health Extension

## Overview

Create a VS Code extension called **Code Health** (`code-health`) that provides at-a-glance visualization of codebase structure, size, complexity trends, and test coverage health. The primary audience is developers using GenAI coding tools (Claude Code, Copilot, Cursor, etc.) who need to quickly understand the state of a codebase that's changing rapidly through AI-assisted development.

## Core Concept

The extension answers three questions at a glance:
1. **Structure** — What does this codebase look like? (modules, dependencies, file sizes)
2. **Refactoring** — Which files are getting too large and need splitting?
3. **Testing** — Where are the testing gaps?

## Feature 1: Tree View (Sidebar Panel)

Add a custom Tree View in the VS Code Explorer sidebar (or its own activity bar icon) that mirrors the project's folder structure but augments each item with inline sparkline visualizations.

### Folder nodes show:
- Folder name
- Total lines of code (aggregated from all source files recursively)
- A **size sparkline bar** — a compact inline bar (like `▁▂▃▅▇`) proportional to the folder's LOC relative to the largest folder. Color-coded:
  - The color represents the health of the *largest file* in that folder (green if all files <500 lines, yellow if any 500-999, orange if any 1000-1499, red if any 1500+)
- File count badge (e.g., `12 files`)

### File nodes show:
- File name
- Line count
- A **size sparkline bar** — proportional to the file's LOC relative to the largest file in the project. Color-coded by threshold:
  - Green: <500 lines
  - Yellow: 500-999 lines
  - Orange: 1000-1499 lines
  - Red: 1500+ lines
- A **trend sparkline** — a mini line chart (e.g., `▁▂▃▅▇▅▃`) showing how the file's line count changed over the last 10 git commits that touched it. Rising trend = growing file. Flat = stable. Helps spot files that are ballooning.

### Tree View behavior:
- Respects `.gitignore` and excludes `node_modules`, `vendor`, `.git`, generated code, etc.
- Configurable exclude patterns via settings
- Source file types are configurable (default: `*.go`, `*.ts`, `*.tsx`, `*.js`, `*.jsx`, `*.py`, `*.rs`, `*.java`, `*.cs`)
- Click on a file opens it in the editor
- Refresh button to recalculate
- Sort options: by name (default), by size (largest first), by recent change, by health (worst first)
- Collapsible folders with aggregated stats

## Feature 2: Code Health Dashboard (Webview Document)

A command (`Code Health: Open Dashboard`) that opens a webview panel with a rich visualization of the entire project. Think of it as an interactive treemap / block diagram.

### Layout:
- The dashboard shows **one block per top-level module/folder**
- Each block is sized proportionally to its total LOC
- Inside each block, show nested blocks for subfolders (like a treemap)
- Leaf blocks (files) are colored by their health (green/yellow/orange/red thresholds)

### Each module block contains these sparkline indicators:

1. **File count** — Simple number with a small bar: `14 files ▅▅▅▅`
2. **Size health** — Distribution sparkline showing how many files fall into each health bucket. e.g., `██████▓▓░░` where green blocks = healthy files, yellow = medium, red = large. This gives an instant read on whether the module is well-factored.
3. **Change trend** — A sparkline showing commit activity over the last 10 commits (or 2 weeks). Tall bars = lots of churn. Helps identify hotspots where GenAI tools are actively generating/modifying code.
4. **Test health** — A bar or ratio showing test file coverage:
   - Ratio of `_test.go` / `.test.ts` / `.spec.ts` files to source files
   - Color: green if ratio > 0.5, yellow if 0.3-0.5, red if < 0.3
   - Or based on actual coverage data if available (lcov, go coverage)

### Dashboard interactions:
- Click a module block to drill down into its subfolders
- Click a file block to open the file in the editor
- Hover shows detailed stats (exact LOC, last modified, top contributors)
- Zoom controls for large projects
- Export as SVG or PNG for documentation

## Feature 3: Dependency Awareness (Stretch Goal)

For supported languages, overlay dependency information:

### Go projects:
- Use `go list -json ./...` to extract package imports
- Show dependency arrows between module blocks in the dashboard
- Sort the tree view by dependency depth (leaf packages first, like the source_tree.sh script)
- Show a "tier" badge on each folder indicating its dependency depth

### TypeScript/JavaScript projects:
- Parse `import` statements to build a dependency graph
- Use the same dependency-depth sorting

### Generic fallback:
- If language-specific tooling isn't available, just show the size/health/trend visualizations without dependency ordering

## Technical Architecture

### Stack:
- **Language**: TypeScript
- **Build**: Use Bun's built-in bundler (`bun build`) — no need for esbuild
- **Package manager**: bun (use `bun install`, `bun run compile`, etc. — NOT npm/node)
- **Webview**: HTML/CSS/SVG for the dashboard. Consider using a lightweight charting library (e.g., `sparkline-svg` or hand-rolled SVG) for sparklines. Avoid heavy frameworks — keep it fast.
- **Git integration**: Use `simple-git` or shell out to `git log --follow --format='%H' -- <file>` + `git show <hash>:<file> | wc -l` for historical line counts
- **Tree View**: VS Code TreeDataProvider API

### Performance considerations:
- Cache line counts and git history in memory (invalidate on file save)
- Use a file watcher to detect changes
- Compute git history lazily (only when a file is expanded or the dashboard is opened)
- For large repos, limit git history to last 10 commits per file
- Use worker threads or batched processing for initial scan

### Configuration (VS Code settings):
```json
{
  "codeHealth.excludePatterns": ["**/node_modules/**", "**/vendor/**", "**/.git/**"],
  "codeHealth.sourceExtensions": ["go", "ts", "tsx", "js", "jsx", "py", "rs"],
  "codeHealth.thresholds": {
    "green": 500,
    "yellow": 1000,
    "orange": 1500
  },
  "codeHealth.gitHistoryDepth": 10,
  "codeHealth.testPatterns": ["*_test.go", "*.test.ts", "*.test.tsx", "*.spec.ts", "*.spec.tsx"]
}
```

### Extension activation:
- Activate on workspace open (lightweight — just register tree view)
- Lazy-load dashboard webview only when requested
- Background scan on activation with progress indicator

## Project Structure

```
code-health/
├── src/
│   ├── extension.ts          # Extension entry point, command registration
│   ├── treeView/
│   │   ├── provider.ts       # TreeDataProvider implementation
│   │   ├── treeItem.ts       # Custom TreeItem with sparkline decorations
│   │   └── scanner.ts        # File system scanner with caching
│   ├── dashboard/
│   │   ├── panel.ts          # Webview panel management
│   │   ├── data.ts           # Data aggregation for dashboard
│   │   └── webview/          # HTML/CSS/JS for the webview
│   │       ├── index.html
│   │       ├── dashboard.ts  # Dashboard rendering logic
│   │       ├── treemap.ts    # Treemap layout algorithm
│   │       └── sparkline.ts  # SVG sparkline generator
│   ├── analysis/
│   │   ├── lineCounter.ts    # Count lines per file
│   │   ├── gitHistory.ts     # Git log parsing for trends
│   │   ├── testDetector.ts   # Find test files and compute ratios
│   │   └── dependencies.ts   # Language-specific dependency extraction
│   └── utils/
│       ├── config.ts         # Settings management
│       ├── cache.ts          # In-memory cache with invalidation
│       └── sparkline.ts      # Unicode sparkline rendering for tree view
├── media/
│   └── icon.svg              # Activity bar icon
├── package.json
├── tsconfig.json
├── build.ts                  # Bun build script
└── README.md
```

## Sparkline Rendering

### Tree View (inline text sparklines):
Use Unicode block characters for inline sparklines in tree item descriptions:
- Block elements: `▁▂▃▄▅▆▇█` (U+2581 through U+2588) — 8 height levels
- These render in the tree view's `description` field as monospace text
- Color them using the ThemeColor API or TreeItem's `iconPath` with generated SVG icons

### Dashboard (SVG sparklines):
Generate small SVG sparklines (~80x20px) for embedding in the treemap blocks:
- Line sparklines for trends (series of points connected by lines)
- Bar sparklines for distributions (small vertical bars)
- Use HSL color interpolation for smooth green-to-red gradients

## MVP Scope (v0.1)

For the initial release, implement:
1. Tree View with file/folder LOC counts and color-coded size bars
2. Basic sorting (by name, by size)
3. Configurable exclude patterns and source extensions
4. `Code Health: Open Dashboard` command with a simple treemap view colored by file health
5. File count and LOC aggregation per folder

### Post-MVP (v0.2+):
- Git history trend sparklines
- Test health indicators
- Dependency awareness
- Dashboard drill-down and export

## Design Principles

Follow **Edward Tufte's principles** for information design throughout. With every new feature, visualization, or UI decision, ask: **"What would Tufte do?"** If an element doesn't carry data, remove it. If two separate views can be merged into one dense, layered visualization, merge them. Prefer showing the data over describing it.

1. **High data density** — Maximize the data-ink ratio. Every pixel should convey information. A single tree view row should communicate file name, size, health status, and trend in one line. No decorative chrome, no redundant labels. Sparklines are the ideal Tufte element: data-intense, small, word-sized graphics that can be embedded inline.
2. **Multiple aspects in one visualization** — Each visual element should encode several dimensions simultaneously. A sparkline bar conveys size (length), health (color), and trend (shape) in one compact glyph. The dashboard treemap encodes module hierarchy (nesting), size (area), health (color), and activity (sparkline overlays) in a single view. Avoid separate pages for separate metrics — integrate them.
3. **Clean, chartjunk-free design** — No gratuitous gradients, 3D effects, heavy gridlines, or decorative elements. Use whitespace, subtle borders, and typographic hierarchy. Let the data speak. Muted colors for healthy state, saturated colors only for items needing attention (so problems visually pop).
4. **Small multiples** — The dashboard's module blocks act as small multiples: same layout repeated per module, making cross-module comparison effortless. Each block has the same sparkline positions so the eye can scan across modules rapidly.
5. **Fast** — Initial scan should complete in <2 seconds for a 1000-file project.
6. **Lightweight** — No heavy dependencies. Minimal bundle size.
7. **Language-agnostic** — Core features (LOC, health, trends) work for any language. Dependency features are opt-in per language.
8. **Non-intrusive** — Lives in the sidebar. Doesn't modify files or interrupt workflow.
