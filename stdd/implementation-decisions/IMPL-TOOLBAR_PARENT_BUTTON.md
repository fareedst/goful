# [IMPL:TOOLBAR_PARENT_BUTTON] Toolbar Parent Button

**Cross-References**: [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]  
**Status**: Active  
**Created**: 2026-01-17  
**Last Updated**: 2026-01-18

---

## Decision

Implement a clickable parent navigation button `[^]` integrated into the filer header row at the left edge, with hit-testing infrastructure and action dispatch that respects Linked navigation mode.

## Rationale

- Provides mouse-first access to parent directory navigation without keyboard shortcuts
- Reuses existing `linkedParentNav` logic for consistent behavior with keyboard navigation
- Minimizes UI disruption by using available header space rather than adding a new row
- Establishes reusable toolbar infrastructure for future buttons

## Implementation Approach

### Module 1: ToolbarRenderer (`filer/filer.go`)

Add toolbar button rendering to `drawHeader()`:

```go
// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
// Draw toolbar at left edge before workspace tabs
var toolbarBounds = make(map[string]struct{ x1, y, x2 int })

func (f *Filer) drawHeader() {
    x, y := f.LeftTop()
    
    // Draw parent button
    // [IMPL:TOOLBAR_PARENT_BUTTON] Parent navigation button
    buttonText := "[^]"
    x = widget.SetCells(x, y, buttonText, look.Default())
    toolbarBounds["parent"] = struct{ x1, y, x2 int }{0, y, x - 1}
    x = widget.SetCells(x, y, " ", look.Default())
    
    // Continue with workspace tabs...
}
```

### Module 2: ToolbarHitTest (`filer/filer.go`)

Add hit-testing function:

```go
// ToolbarButtonAt returns the toolbar button identifier at coordinates (x, y).
// Returns empty string if no button is at that position.
// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
func ToolbarButtonAt(x, y int) string {
    for name, bounds := range toolbarBounds {
        if y == bounds.y && x >= bounds.x1 && x <= bounds.x2 {
            return name
        }
    }
    return ""
}
```

### Module 3: ToolbarDispatcher (`app/goful.go`)

Extend `handleLeftClick()` to check toolbar before directory hit-testing:

```go
func (g *Goful) handleLeftClick(x, y int) {
    // [IMPL:TOOLBAR_PARENT_BUTTON] Check toolbar buttons first (header row)
    if y == 0 {
        if button := filer.ToolbarButtonAt(x, y); button != "" {
            g.handleToolbarClick(button)
            return
        }
    }
    
    // Existing directory hit-testing logic...
}

// handleToolbarClick dispatches toolbar button clicks to appropriate actions.
// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
func (g *Goful) handleToolbarClick(button string) {
    switch button {
    case "parent":
        g.toolbarParentNav()
    }
}

// toolbarParentNav navigates to parent directory, respecting Linked mode.
// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
func (g *Goful) toolbarParentNav() {
    if g.IsLinkedNav() {
        g.Workspace().ChdirAllToParent()
    }
    g.Dir().Chdir("..")
    if g.IsLinkedNav() {
        g.Workspace().RebuildComparisonIndex()
    }
}
```

## Code Markers

Files that must carry `[IMPL:TOOLBAR_PARENT_BUTTON]` annotations:
- `filer/filer.go` - `drawHeader()`, `ToolbarButtonAt()`
- `app/goful.go` - `handleLeftClick()`, `handleToolbarClick()`, `toolbarParentNav()`

Tests that must reference `[REQ:TOOLBAR_PARENT_BUTTON]`:
- `filer/filer_test.go` - `TestToolbarButtonAt_REQ_TOOLBAR_PARENT_BUTTON`
- `app/goful_test.go` - `TestToolbarParentNav_REQ_TOOLBAR_PARENT_BUTTON`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `filer/filer.go` `drawHeader()` - `[IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]`
- [x] `filer/filer.go` `ToolbarButtonAt()` - `[IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]`
- [x] `filer/filer.go` `InvokeToolbarButton()` - `[IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]`
- [x] `app/goful.go` `handleLeftClick()` - `[IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]`
- [x] `app/goful.go` `HandleParentButtonPress()` - `[IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]`
- [x] `main.go` callback wiring - `[IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]`

Tests that must reference `[REQ:TOOLBAR_PARENT_BUTTON]`:
- [x] `TestToolbarButtonAt_REQ_TOOLBAR_PARENT_BUTTON`
- [x] `TestToolbarButtonAtMultipleButtons_REQ_TOOLBAR_PARENT_BUTTON`
- [x] `TestInvokeToolbarButton_REQ_TOOLBAR_PARENT_BUTTON`
- [x] `TestInvokeToolbarButtonWithNilCallback_REQ_TOOLBAR_PARENT_BUTTON`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-17 | — | ⏳ Pending | Initial plan documented |
| 2026-01-18 | — | ✅ Pass | Token validation: 1326 refs across 77 files |

## Related Decisions

- Depends on: [IMPL:LINKED_NAVIGATION], [IMPL:MOUSE_FILE_SELECT]
- See also: [ARCH:MOUSE_EVENT_ROUTING]

---

*Created on 2026-01-17*
