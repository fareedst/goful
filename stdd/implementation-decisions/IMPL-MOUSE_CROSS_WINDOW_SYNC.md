# [IMPL:MOUSE_CROSS_WINDOW_SYNC] Mouse Cross-Window Cursor Synchronization

**Cross-References**: [ARCH:MOUSE_CROSS_WINDOW_SYNC] [REQ:MOUSE_CROSS_WINDOW_SYNC]  
**Status**: Active  
**Created**: 2026-01-17  
**Last Updated**: 2026-01-18

---

## Decision

Add cross-window cursor synchronization to `handleLeftClick` in `app/goful.go`, fix the rendering logic in `filer/directory.go` to display cursor highlights in unfocused windows, and implement cursor hiding for windows without matching files.

## Rationale

- Reuses existing `SetCursorByNameAll()` infrastructure from diff search feature
- Single-line addition to existing mouse handling code
- No new dependencies or complex logic required
- Rendering fix ensures cursor is visually highlighted in all windows, not just the focused one
- Cursor hiding prevents misleading highlights in windows that don't contain the clicked file

## Implementation

### 1. Cursor Sync in `handleLeftClick()` (`app/goful.go`)

```go
// Convert Y to file index and move cursor
fileIdx := dir.FileIndexAtY(y)
if fileIdx >= 0 {
    dir.SetCursor(fileIdx)
    // [IMPL:MOUSE_CROSS_WINDOW_SYNC] Sync cursor to same filename in all windows
    filename := dir.File().Name()
    ws.SetCursorByNameAll(filename)
}
```

### 2. Rendering Fix in `drawFilesWithComparison()` (`filer/directory.go`)

**Problem:** The original code only highlighted the cursor when the window was focused:
```go
isFocused := focus && i == d.Cursor()  // BUG: Only highlights in focused window
```

**Fix:** Decouple cursor highlighting from window focus, but respect cursorHidden state:
```go
isFocused := i == d.Cursor() && !d.IsCursorHidden()  // FIX: Highlight cursor in all windows, hide when no match
```

This allows the cursor position to be visually highlighted in ALL windows where the file exists, while hiding the highlight in windows without a matching file.

### 3. Cursor Hidden State in `widget/listbox.go`

**Added `cursorHidden` field to `ListBox` struct:**
```go
type ListBox struct {
    // ... existing fields ...
    cursorHidden bool // [IMPL:MOUSE_CROSS_WINDOW_SYNC] when true, cursor highlight is not shown
}
```

**Fixed `IndexByName` to return `-1` when not found:**
```go
func (b *ListBox) IndexByName(name string) int {
    for i, content := range b.list {
        if name == content.Name() {
            return i
        }
    }
    return -1  // Previously returned b.lower, causing incorrect cursor movement
}
```

**Modified `SetCursorByName` to control cursor visibility:**
```go
func (b *ListBox) SetCursorByName(name string) {
    idx := b.IndexByName(name)
    if idx != -1 {
        b.SetCursor(idx)
        b.cursorHidden = false
    } else {
        b.cursorHidden = true  // Hide cursor when file not found
    }
}
```

**Added `IsCursorHidden()` method for rendering:**
```go
func (b *ListBox) IsCursorHidden() bool {
    return b.cursorHidden
}
```

**Modified `SetCursor` and `MoveCursor` to reset hidden state:**
Direct cursor movements (from keyboard navigation) always show the cursor by resetting `cursorHidden = false`.

## Code Markers

- `app/goful.go`: `handleLeftClick` with `// [IMPL:MOUSE_CROSS_WINDOW_SYNC] [ARCH:MOUSE_CROSS_WINDOW_SYNC] [REQ:MOUSE_CROSS_WINDOW_SYNC]`
- `filer/directory.go`: `drawFilesWithComparison` with rendering fix (line ~556)

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files that must carry annotations:
- [x] `app/goful.go` - handleLeftClick

Tests that must reference `[REQ:MOUSE_CROSS_WINDOW_SYNC]`:
- [x] `filer/integration_test.go` - `TestSetCursorByNameAll_REQ_MOUSE_CROSS_WINDOW_SYNC`
- [x] `filer/integration_test.go` - `TestSetCursorByNameAllFocusUnchanged_REQ_MOUSE_CROSS_WINDOW_SYNC`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-17 | — | ✅ Pass | Initial implementation. `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1298 token references across 77 files.` |
| 2026-01-18 | — | ✅ Pass | Rendering fix applied. `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1299 token references across 77 files.` |

## Bug Fix History

| Date | Issue | Root Cause | Fix |
|------|-------|------------|-----|
| 2026-01-18 | Cursor not visually highlighted in unfocused windows | `drawFilesWithComparison` in `filer/directory.go` conditioned highlighting on `focus && i == d.Cursor()` | Changed to `i == d.Cursor()` to decouple highlighting from window focus |
| 2026-01-18 | Windows without matching file incorrectly highlight first directory | `IndexByName` returned `b.lower` instead of `-1` when file not found, causing `SetCursorByName` to always call `SetCursor()` | Fixed `IndexByName` to return `-1`; added `cursorHidden` flag to hide highlight when file not found; updated rendering to check `!d.IsCursorHidden()` |

## Related Decisions

- Depends on: [IMPL:MOUSE_FILE_SELECT], [IMPL:MOUSE_HIT_TEST]
- See also: [ARCH:MOUSE_CROSS_WINDOW_SYNC], [REQ:MOUSE_CROSS_WINDOW_SYNC], [REQ:MODULE_VALIDATION]

---

*Created as part of mouse cross-window cursor sync feature on 2026-01-17*
*Updated 2026-01-18 with cursor hiding fix for non-matching windows*