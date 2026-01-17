# [IMPL:MOUSE_FILE_SELECT] Mouse File Selection

**Cross-References**: [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]  
**Status**: Active  
**Created**: 2026-01-15  
**Last Updated**: 2026-01-17

---

## Decision

Wire mouse click events to file cursor movement and directory focus switching.

## Rationale

- Implements the core mouse selection feature per [REQ:MOUSE_FILE_SELECT]
- Uses the hit-testing framework from [IMPL:MOUSE_HIT_TEST] to find the target directory and file
- Integrates with the existing event loop pattern in `app.Goful`

## Implementation Approach

### Extend `eventHandler` in `app/goful.go`

```go
func (g *Goful) eventHandler(ev tcell.Event) {
    switch ev := ev.(type) {
    case *tcell.EventKey:
        // ... existing code ...
    case *tcell.EventResize:
        // ... existing code ...
    case *tcell.EventMouse:
        g.mouseHandler(ev)
    }
}
```

### Add `mouseHandler` in `app/goful.go`

```go
// mouseHandler handles mouse events for file selection and scrolling.
// [IMPL:MOUSE_FILE_SELECT] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func (g *Goful) mouseHandler(ev *tcell.EventMouse) {
    x, y := ev.Position()
    buttons := ev.Buttons()
    
    // Handle modal widgets first
    if !widget.IsNil(g.Next()) {
        // TODO: Modal mouse handling (future stage)
        return
    }
    
    // Handle left click for file selection
    if buttons&tcell.Button1 != 0 {
        g.handleLeftClick(x, y)
    }
    
    // Handle wheel for scrolling
    if buttons&tcell.WheelUp != 0 {
        g.handleWheelUp(x, y)
    }
    if buttons&tcell.WheelDown != 0 {
        g.handleWheelDown(x, y)
    }
}

func (g *Goful) handleLeftClick(x, y int) {
    ws := g.Workspace()
    dir, idx := ws.DirectoryAt(x, y)
    if dir == nil {
        return
    }
    
    // Switch focus if clicking in unfocused window
    if idx != ws.Focus {
        ws.SetFocus(idx)
    }
    
    // Convert Y to file index and move cursor
    fileIdx := dir.FileIndexAtY(y)
    if fileIdx >= 0 {
        dir.SetCursor(fileIdx)
    }
}
```

## Code Markers

- `app/goful.go`: `mouseHandler`, `handleLeftClick`, `handleWheelUp`, `handleWheelDown` with `// [IMPL:MOUSE_FILE_SELECT] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `app/goful.go` - mouseHandler, handleLeftClick, handleWheelUp, handleWheelDown

Tests that must reference `[REQ:MOUSE_FILE_SELECT]`:
- [ ] Integration tests for mouse file selection

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-15 | — | ✅ Pass | Mouse file selection working |

## Related Decisions

- Depends on: [IMPL:MOUSE_HIT_TEST]
- See also: [ARCH:MOUSE_EVENT_ROUTING], [REQ:MOUSE_FILE_SELECT], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
