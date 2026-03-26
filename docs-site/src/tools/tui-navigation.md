# Navigation and Layout

The TUI uses a Miller column layout for hierarchical browsing of Mendix projects. This design, inspired by file managers like ranger and yazi, shows the current context alongside parent and child items for spatial orientation.

## Three-Column Layout

The default layout divides the terminal into three panels:

| Panel | Content | Width |
|-------|---------|-------|
| Left | Modules and folders | ~20% |
| Center | Documents within the selected module (entities, microflows, pages, etc.) | ~30% |
| Right | Preview of the selected document (MDL/NDSL rendering with syntax highlighting) | ~50% |

The layout adapts dynamically as you navigate:

- **One panel**: Only the module list is visible when no module is selected.
- **Two panels**: Selecting a module expands the document list alongside it (35% + 65% split).
- **Three panels**: Selecting a document adds the preview panel (20% + 30% + 50% split).

## Breadcrumb Navigation

Each panel displays a breadcrumb trail at the top showing the current navigation path. For example:

```
MyModule > DomainModel > Customer
```

Breadcrumbs are clickable with the mouse. Clicking a segment navigates back to that level in the hierarchy. The navigation stack remembers your path so you can drill into nested folders and return.

## Keyboard Navigation

The TUI uses vim-style keybindings for movement:

| Key | Action |
|-----|--------|
| `j` / Down | Move cursor down in the current panel |
| `k` / Up | Move cursor up in the current panel |
| `h` / Left | Move focus to the left panel (or go back in navigation stack) |
| `l` / Right / Enter | Move focus to the right panel (or drill into the selected item) |
| `Tab` | Cycle focus between panels (left, center, right) |
| `/` | Activate filter mode in the current panel |
| `Esc` | Exit filter mode, close overlay, or cancel current action |

### Filtering

Press `/` to activate the filter input in the current panel. Type a substring to filter the list in real time. Only items matching the filter text are shown. Press `Esc` to clear the filter and return to the full list.

## Mouse Support

The TUI enables mouse interaction for all panels:

| Action | Effect |
|--------|--------|
| **Click** on an item | Selects it and focuses the panel |
| **Scroll wheel** | Scrolls the list up or down |
| **Click** on a breadcrumb segment | Navigates back to that level |

Mouse coordinates are translated to panel-local positions using the layout geometry, so clicks always target the correct panel regardless of terminal size.

## Fullscreen Overlay

Press **Enter** or **Z** on a selected document to open a fullscreen overlay showing the complete document content. The overlay includes:

- A title bar with the document name and type
- Scrollable content with syntax highlighting (MDL, SQL, or NDSL)
- Hint bar at the bottom showing available keys

Press `Esc` or `q` to close the overlay and return to the three-column browser.

### Zoom Mode

Press `z` to zoom the currently focused panel to fill the entire terminal. This is useful for reviewing long module lists or document trees. Press `z` again or `Esc` to return to the previous layout. The TUI remembers which panel was focused and which layout was active before zooming.

## Tab Support

The TUI supports multiple project tabs within a single session. Each tab maintains its own:

- Project file path
- Navigation state (selected module, document, preview)
- Panel focus and scroll positions

Tabs allow you to compare or cross-reference elements across different projects or different areas of the same project.

## Contextual Help

Press `?` at any time to open the help overlay. The help content is contextual -- it shows keybindings relevant to the current mode (browser, overlay, command palette, or execution view).
