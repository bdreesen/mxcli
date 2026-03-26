# Session Management

The TUI supports session persistence, allowing you to restore your previous browsing state when relaunching the tool.

## Restoring a Session

Use the `-c` or `--continue` flag to restore the previous session:

```bash
mxcli tui -c
```

This restores:
- Open project tabs and their file paths
- Navigation state (selected module, document, scroll positions)
- Panel focus and layout configuration
- Preview mode settings

If the session file cannot be loaded (e.g., first launch or corrupted state), a warning is printed to stderr and the TUI starts fresh.

## Combining with Project Flag

When using `-c` with `-p`, the explicit project path takes precedence:

```bash
# Restore session but open a specific project
mxcli tui -c -p app.mpr
```

When using `-c` alone, the project path is loaded from the saved session. If the session contains multiple tabs, the first tab's project path is used.

## Session Storage

Session state is stored locally in the user's configuration directory. The TUI automatically saves session state when exiting normally (via `q`).

## Project File Picker

When launching without either `-p` or `-c`, the TUI presents an interactive file picker:

```bash
mxcli tui
```

The picker shows:
- Recently opened projects (from session history)
- Filesystem browser filtered to `.mpr` files

Select a project and press `Enter` to launch the TUI. The selected project is saved to the history for future sessions.

## Session History

The TUI maintains a history of opened projects. This history is separate from the full session state and persists across sessions. It powers:

- The file picker's "recent projects" list
- Session restore when using `-c`

Each time a project is opened, it is added to the history via `tui.SaveHistory(projectPath)`.
