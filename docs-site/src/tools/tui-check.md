# Project Validation

The TUI integrates with `mx check` to validate Mendix projects directly from the terminal interface. Errors, warnings, and deprecations are displayed in an interactive overlay with grouping, navigation, and filtering.

## Running a Check

Use the `:check` command from the command palette (press `:` then type `check`) to validate the current project. The TUI runs `mx check -j -w -d` in the background and displays results in a fullscreen overlay.

## Error Overlay

The check results overlay groups diagnostics by error code and deduplicates repeated occurrences:

```
mx check Results  [All: 8E 2W 1D]

CE1613 -- The selected association/attribute no longer exists
  MyFirstModule.P_ComboBox_Enum (Page)
    > Property 'Association' of combo box 'cmbPriority'
  MyFirstModule.P_ComboBox_Assoc (Page) (x7)
    > Property 'Attribute' of combo box 'cmbCategory'
```

Each group shows:
- The error code and description
- Affected documents with their type
- Specific element details
- Occurrence count for deduplicated entries (shown as `(xN)`)

## Filtering by Severity

The overlay supports filtering by diagnostic severity. Press `Tab` to cycle through filters:

| Filter | Shows |
|--------|-------|
| **All** | Errors, warnings, and deprecations |
| **Errors** | Only errors |
| **Warnings** | Only warnings |
| **Deprecations** | Only deprecations |

The title bar reflects the active filter: `[All: 8E 2W 1D]` or `[Errors: 8]`.

## Error Navigation

The check overlay is a selectable list. Use `j`/`k` to move the cursor between error locations. Press `Enter` on a location to:

1. Close the overlay
2. Navigate the project tree to the affected document

After closing the overlay, the TUI enters **check navigation mode**. The status bar shows the current error position and navigation hints:

```
[2/5] CE1613: MyFirstModule.P_ComboBox_Enum  ]e next  [e prev
```

### Navigation Keys

| Key | Action |
|-----|--------|
| `]e` | Jump to the next error location in the project tree |
| `[e` | Jump to the previous error location |
| `!` | Reopen the check overlay |
| `Esc` | Exit check navigation mode |

## File Monitoring

The TUI watches the project file for changes. When it detects that the MPR file has been modified (e.g., after saving in Studio Pro or executing an MDL script), it can automatically re-run `mx check` and update the results. This provides a continuous validation feedback loop when making changes.

## LLM Anchor Tags

The check overlay embeds faint structured anchor tags in its output for machine consumption. These tags are nearly invisible in the terminal but are preserved when copying text or taking screenshots:

```
[mxcli:check] errors=8 warnings=2 deprecations=1
[mxcli:check:CE1613] severity=ERROR count=6 doc=MyFirstModule.P_ComboBox_Assoc type=Page
```

These anchors enable AI assistants to parse check results from terminal screenshots or clipboard content, providing structured data for automated error analysis and fix suggestions.

## Status Badge

When check results are available, a compact badge appears in the TUI status bar showing the diagnostic summary:

```
8E 2W 1D
```

This badge updates whenever a new check completes, providing an at-a-glance project health indicator without opening the full overlay.
