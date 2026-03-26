# Terminal UI (TUI)

mxcli includes a terminal-based interactive UI for browsing, inspecting, and modifying Mendix projects. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), it provides a ranger/yazi-style Miller column interface directly in the terminal.

## Quick Start

```bash
# Launch with a project file
mxcli tui -p app.mpr

# Resume a previous session
mxcli tui -c

# Launch without a project (opens file picker)
mxcli tui
```

The TUI opens in full-screen alternate screen mode with mouse support enabled.

## Feature Summary

| Feature | Description |
|---------|-------------|
| **Miller column navigation** | Three-panel layout: modules, documents, preview |
| **Vim-style keybindings** | h/j/k/l navigation, / to filter, : for commands |
| **Syntax highlighting** | Real-time MDL, SQL, and NDSL highlighting via Chroma |
| **Command palette** | VS Code-style `:` command bar with tab completion |
| **MDL execution** | Type or paste MDL scripts and execute them in-place |
| **mx check integration** | Validate projects with grouped errors, navigation, and filtering |
| **Mouse support** | Click to select, scroll wheel, clickable breadcrumbs |
| **Session restore** | `-c` flag restores previous tabs, selections, and navigation state |
| **File watcher** | Auto-detects MPR changes and refreshes the tree |
| **Fullscreen preview** | Press Enter or Z to expand any panel to full screen |
| **Tab support** | Multiple project tabs in a single session |
| **Contextual help** | Press `?` for keybinding reference |

## Architecture

The TUI delegates all heavy lifting to existing mxcli logic. The project tree is built by calling `buildProjectTree()` directly, while commands like describe, callers, check, and exec run mxcli as a subprocess. Diagrams open in the system browser via `xdg-open` or `open`.

## Flags

| Flag | Description |
|------|-------------|
| `-p, --project` | Path to the `.mpr` project file |
| `-c, --continue` | Restore the previous TUI session (tabs, navigation, preview mode) |

When launched without `-p`, the TUI presents a file picker that shows recently opened projects and allows browsing the filesystem for `.mpr` files.

## Keyboard Reference

| Key | Action |
|-----|--------|
| `h` / Left | Move focus left |
| `l` / Right / Enter | Move focus right / open |
| `j` / Down | Move down |
| `k` / Up | Move up |
| `Tab` | Cycle panel focus |
| `/` | Filter current column |
| `:` | Open command palette |
| `x` | Open MDL execution view |
| `z` / `Z` | Toggle fullscreen zoom on current panel |
| `?` | Show contextual help |
| `q` | Quit |
