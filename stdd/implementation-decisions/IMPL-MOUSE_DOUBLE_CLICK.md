# [IMPL:MOUSE_DOUBLE_CLICK] Mouse Double-Click Detection

**Cross-References**: [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]  
**Status**: Active  
**Created**: 2026-01-17  
**Last Updated**: 2026-01-18

---

## Decision

Implement time-based double-click detection and action handlers for files and directories.

## Rationale

- Implements double-click behavior per [REQ:MOUSE_DOUBLE_CLICK]
- Reuses the linked navigation pattern for directories per [REQ:LINKED_NAVIGATION]
- Extends existing mouse handler infrastructure from [IMPL:MOUSE_FILE_SELECT]

## Implementation Approach

### Add click state fields to `Goful` struct in `app/goful.go`

```go
lastClickTime time.Time  // [IMPL:MOUSE_DOUBLE_CLICK]
lastClickX    int        // [IMPL:MOUSE_DOUBLE_CLICK]
lastClickY    int        // [IMPL:MOUSE_DOUBLE_CLICK]
```

### Add `isDoubleClick` helper

```go
const doubleClickThreshold = 400 * time.Millisecond

// isDoubleClick checks if this click is a double-click based on timing and position.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func (g *Goful) isDoubleClick(x, y int) bool {
    now := time.Now()
    isDouble := now.Sub(g.lastClickTime) < doubleClickThreshold &&
                g.lastClickX == x && g.lastClickY == y
    g.lastClickTime = now
    g.lastClickX = x
    g.lastClickY = y
    return isDouble
}
```

### Add `handleDoubleClickDir` for directory navigation

```go
// handleDoubleClickDir navigates into a directory, respecting linked mode.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func (g *Goful) handleDoubleClickDir(dir *filer.Directory) {
    if g.IsLinkedNav() {
        name := dir.File().Name()
        navigated, skipped := g.Workspace().ChdirAllToSubdirNoRebuild(name)
        if skipped > 0 {
            g.SetLinkedNav(false)
            message.Infof("linked navigation disabled: %d window(s) missing '%s'", skipped, name)
        }
        _ = navigated
    }
    dir.EnterDir()
    if g.IsLinkedNav() {
        g.Workspace().RebuildComparisonIndex()
    }
}
```

### Add `handleDoubleClickFile` for file opening

When Linked mode is enabled, the handler collects all matching file paths from all windows and opens each one with a **separate command**. This ensures each file is opened independently, which provides better compatibility across different opener applications.

```go
// handleDoubleClickFile opens a file, and opens same-named files in all windows when linked.
// When linked mode is ON, opens same-named files from ALL windows where they exist.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func (g *Goful) handleDoubleClickFile(dir *filer.Directory) {
    filename := dir.File().Name()

    if g.IsLinkedNav() {
        // Collect all file paths from windows that have the same-named file
        var filePaths []string
        for _, d := range g.Workspace().Dirs {
            if found := d.FindFileByName(filename); found != nil {
                d.SetCursorByName(filename)
                filePaths = append(filePaths, found.Path())
            }
        }

        // Open each matching file with a separate command
        opener := "xdg-open"  // Linux
        switch runtime.GOOS {
        case "windows":
            opener = "explorer"
        case "darwin":
            opener = "open"
        }
        for _, path := range filePaths {
            g.Spawn(opener + " " + util.Quote(path) + " %&")
        }
        return
    }

    // Non-linked mode: only open the focused file
    g.Input("C-m")
}
```

**Key design choice**: Using separate commands for each file (one `open` invocation per file) rather than a single command with multiple file arguments. This approach:
- Ensures consistent behavior across different opener applications
- Allows each file to be opened by its associated application independently
- Uses `util.Quote()` for proper shell escaping of paths with special characters

### Modify `handleLeftClick` to detect double-click

The modified function checks for double-click after selection and dispatches to the appropriate handler based on whether the clicked item is a file or directory.

## Code Markers

- `app/goful.go`: `doubleClickThreshold` constant, `lastClickTime`/`lastClickX`/`lastClickY` fields, `isDoubleClick`, `handleDoubleClickDir`, `handleDoubleClickFile` with `// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]`

## Test Coverage

- `app/mouse_test.go`:
  - `TestIsDoubleClick_REQ_MOUSE_DOUBLE_CLICK` (table-driven timing/position tests)
  - `TestDoubleClickThreshold_REQ_MOUSE_DOUBLE_CLICK` (threshold sanity check)
  - `TestIsDoubleClickUpdatesState_REQ_MOUSE_DOUBLE_CLICK` (state updates)
  - `TestDoubleClickSequence_REQ_MOUSE_DOUBLE_CLICK` (realistic click sequences)

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `app/goful.go` - all double-click related functions and fields

Tests that must reference `[REQ:MOUSE_DOUBLE_CLICK]`:
- [ ] Unit tests with names referencing `REQ_MOUSE_DOUBLE_CLICK`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-17 | — | ✅ Pass | 4 tests PASS |
| 2026-01-17 | — | ✅ Pass | verified 1158 token references across 74 files |
| 2026-01-18 | — | ✅ Pass | Bug fix: Linked mode file open now opens all matching files via single command |

## Related Decisions

- Depends on: [IMPL:MOUSE_FILE_SELECT], [IMPL:MOUSE_HIT_TEST]
- See also: [ARCH:MOUSE_DOUBLE_CLICK], [REQ:MOUSE_DOUBLE_CLICK], [REQ:LINKED_NAVIGATION], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
