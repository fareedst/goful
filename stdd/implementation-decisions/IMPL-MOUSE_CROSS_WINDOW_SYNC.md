# [IMPL:MOUSE_CROSS_WINDOW_SYNC] Mouse Cross-Window Cursor Synchronization

**Cross-References**: [ARCH:MOUSE_CROSS_WINDOW_SYNC] [REQ:MOUSE_CROSS_WINDOW_SYNC]  
**Status**: Active  
**Created**: 2026-01-17  
**Last Updated**: 2026-01-18

---

## Decision

Add cross-window cursor synchronization to `handleLeftClick` in `app/goful.go`, and fix the rendering logic in `filer/directory.go` to display cursor highlights in unfocused windows.

## Rationale

- Reuses existing `SetCursorByNameAll()` infrastructure from diff search feature
- Single-line addition to existing mouse handling code
- No new dependencies or complex logic required
- Rendering fix ensures cursor is visually highlighted in all windows, not just the focused one

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

**Fix:** Decouple cursor highlighting from window focus:
```go
isFocused := i == d.Cursor()  // FIX: Highlight cursor in all windows
```

This allows the cursor position to be visually highlighted in ALL windows, enabling users to see the synchronized cursor across all directory panes.

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

## Related Decisions

- Depends on: [IMPL:MOUSE_FILE_SELECT], [IMPL:MOUSE_HIT_TEST]
- See also: [ARCH:MOUSE_CROSS_WINDOW_SYNC], [REQ:MOUSE_CROSS_WINDOW_SYNC], [REQ:MODULE_VALIDATION]

---

*Created as part of mouse cross-window cursor sync feature on 2026-01-17*
