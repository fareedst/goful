# [IMPL:MOUSE_HIT_TEST] Mouse Hit Testing

**Cross-References**: [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]  
**Status**: Active  
**Created**: 2026-01-15  
**Last Updated**: 2026-01-17

---

## Decision

Implement coordinate-based hit-testing for mouse event routing.

## Rationale

- Implements the foundational layer for mouse support per [REQ:MOUSE_FILE_SELECT]
- Separates coordinate math from event handling for independent validation per [REQ:MODULE_VALIDATION]
- Extends existing `widget.Window` with a simple bounds-check method

## Implementation Approach

### Enable Mouse in `widget/widget.go`

- Add `screen.EnableMouse()` call after `screen.Init()` in `Init()` function
- Export `EnableMouse()` and `DisableMouse()` functions for runtime control

### Add `Contains` to `widget.Window`

```go
// Contains returns true if (x, y) is within the window bounds.
// [IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func (w *Window) Contains(x, y int) bool {
    rx, ry := w.RightBottom()
    return x >= w.x && x <= rx && y >= w.y && y <= ry
}
```

### Add `DirectoryAt` to `filer.Workspace`

```go
// DirectoryAt returns the directory containing (x, y) and its index, or nil/-1.
// [IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func (w *Workspace) DirectoryAt(x, y int) (*Directory, int) {
    for i, dir := range w.Dirs {
        if dir.Contains(x, y) {
            return dir, i
        }
    }
    return nil, -1
}
```

### Add `FileIndexAtY` to `filer.Directory`

```go
// FileIndexAtY converts a screen Y coordinate to a list index, or -1 if outside.
// [IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func (d *Directory) FileIndexAtY(y int) int {
    _, topY := d.LeftTop()
    contentStart := topY + 1  // Account for header/border
    row := y - contentStart
    if row < 0 || row >= d.Height()-2 {
        return -1
    }
    idx := d.Offset() + row
    if idx >= d.Upper() {
        return -1
    }
    return idx
}
```

## Code Markers

- `widget/widget.go`: `EnableMouse`, `DisableMouse`, `Contains` with `// [IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]`
- `filer/workspace.go`: `DirectoryAt` with same tokens
- `filer/directory.go`: `FileIndexAtY` with same tokens

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `widget/widget.go` - EnableMouse, DisableMouse, Contains
- [ ] `filer/workspace.go` - DirectoryAt
- [ ] `filer/directory.go` - FileIndexAtY

Tests that must reference `[REQ:MOUSE_FILE_SELECT]`:
- [ ] Unit tests for Contains, DirectoryAt, FileIndexAtY

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-15 | — | ✅ Pass | Hit testing functions work correctly |

## Related Decisions

- Depends on: —
- See also: [ARCH:MOUSE_EVENT_ROUTING], [REQ:MOUSE_FILE_SELECT], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
