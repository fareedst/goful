# [IMPL:CLICKABLE_WORKSPACE_TABS] Clickable Workspace Tabs

**Cross-References**: [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Implement clickable workspace tabs with pill-style visual appearance, positioned at the right edge of the header toolbar after directory paths.

## Rationale

- macOS terminals often intercept Meta key combinations, making `M-f`, `M-b` unreliable
- Click-to-switch provides accessible alternative that works on all platforms
- Pill styling creates clear visual affordance distinguishing tabs from action buttons
- Reuses existing toolbar infrastructure patterns for consistency

## Implementation Approach

### Module 1: Tab Bounds Tracking (`filer/filer.go`)

Add bounds tracking alongside existing toolbar infrastructure:

```go
// [IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]
// workspaceTabBounds tracks clickable workspace tab regions.
// Key is the workspace index, value is the screen bounds.
type workspaceTabBounds struct {
    x1, y, x2 int
}
var workspaceTabs = make(map[int]workspaceTabBounds)

// WorkspaceTabAt returns the workspace index at (x, y), or -1 if none.
// [IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]
func WorkspaceTabAt(x, y int) int {
    for idx, bounds := range workspaceTabs {
        if y == bounds.y && x >= bounds.x1 && x <= bounds.x2 {
            return idx
        }
    }
    return -1
}
```

### Module 2: Tab Callback Mechanism (`filer/filer.go`)

Add callback pattern similar to toolbar buttons:

```go
// [IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]
var workspaceTabClickFn func(index int)

// SetWorkspaceTabClickFn sets the callback for workspace tab clicks.
// [IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]
func SetWorkspaceTabClickFn(fn func(index int)) {
    workspaceTabClickFn = fn
}

// InvokeWorkspaceTab invokes the workspace tab click callback.
// Returns true if handled, false if no callback set.
// [IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]
func InvokeWorkspaceTab(index int) bool {
    if workspaceTabClickFn != nil {
        workspaceTabClickFn(index)
        return true
    }
    return false
}
```

### Module 3: Header Layout Reorder (`filer/filer.go`)

Modify `drawHeader()` to render in new order with pill styling:

```go
// In drawHeader(), after toolbar buttons and separator:

// Render directory paths first (existing code, moved earlier)
ws := f.Workspace()
width := (f.Width() - x - estimatedTabsWidth) / len(ws.Dirs)
for i := 0; i < len(ws.Dirs); i++ {
    // ... existing directory path rendering ...
}

// [IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]
// Clear workspace tab bounds for this frame
workspaceTabs = make(map[int]workspaceTabBounds)

// Pill style: Aqua background for inactive, Lime for current
pillStyle := look.Default().Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
pillStyleCurrent := look.Default().Foreground(tcell.ColorBlack).Background(tcell.ColorLime).Bold(true)

for i, ws := range f.Workspaces {
    s := fmt.Sprintf("«%s»", ws.Title)  // Guillemet delimiters
    tabX1 := x
    style := pillStyle
    if f.Current == i {
        style = pillStyleCurrent
    }
    x = widget.SetCells(x, y, s, style)
    workspaceTabs[i] = workspaceTabBounds{x1: tabX1, y: y, x2: x - 1}
    x = widget.SetCells(x, y, " ", look.Default())  // Space between tabs
}
```

### Module 4: Mouse Dispatch Extension (`app/goful.go`)

Extend `handleLeftClick()`:

```go
func (g *Goful) handleLeftClick(x, y int) {
    // Check toolbar buttons first
    if btnName := filer.ToolbarButtonAt(x, y); btnName != "" {
        filer.InvokeToolbarButton(btnName)
        return
    }

    // [IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]
    // Check workspace tabs
    if wsIdx := filer.WorkspaceTabAt(x, y); wsIdx >= 0 {
        filer.InvokeWorkspaceTab(wsIdx)
        return
    }

    // ... existing directory handling ...
}
```

### Module 5: Callback Wiring (`main.go`)

Register callback during initialization:

```go
// [IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]
filer.SetWorkspaceTabClickFn(func(index int) {
    if index >= 0 && index < len(g.Filer.Workspaces) && index != g.Filer.Current {
        g.Filer.Workspace().visible(false)
        g.Filer.Current = index
        g.Filer.Workspace().visible(true)
    }
})
```

## Visual Specification

**Layout:**
```
^ L = C D R ! | [1] /Users/foo [2] /Users/bar «1» «2» «3»
```

**Styling:**
- Delimiters: Guillemets `«»` (Unicode U+00AB, U+00BB)
- Current tab: Lime background (`tcell.ColorLime`) + bold - clearly indicates active workspace
- Other tabs: Aqua background (`tcell.ColorAqua`) with black foreground (`tcell.ColorBlack`)
- Spacing: Single space between tabs

## Code Markers

Files that must carry `[IMPL:CLICKABLE_WORKSPACE_TABS]` annotations:
- `filer/filer.go` - `workspaceTabBounds`, `WorkspaceTabAt()`, `SetWorkspaceTabClickFn()`, `InvokeWorkspaceTab()`, `drawHeader()`
- `app/goful.go` - `handleLeftClick()`
- `main.go` - callback wiring

Tests that must reference `[REQ:CLICKABLE_WORKSPACE_TABS]`:
- `filer/toolbar_test.go` - `TestWorkspaceTabAt_REQ_CLICKABLE_WORKSPACE_TABS`
- `filer/toolbar_test.go` - `TestInvokeWorkspaceTab_REQ_CLICKABLE_WORKSPACE_TABS`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `filer/filer.go` `workspaceTabBounds` - `[IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]`
- [x] `filer/filer.go` `WorkspaceTabAt()` - `[IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]`
- [x] `filer/filer.go` `SetWorkspaceTabClickFn()` - `[IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]`
- [x] `filer/filer.go` `InvokeWorkspaceTab()` - `[IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]`
- [x] `filer/filer.go` `SwitchToWorkspace()` - `[IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]`
- [x] `filer/filer.go` `drawHeader()` workspace tabs section - `[IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]`
- [x] `app/goful.go` `handleLeftClick()` workspace tab check - `[IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]`
- [x] `main.go` callback wiring - `[IMPL:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [REQ:CLICKABLE_WORKSPACE_TABS]`

**Validation Result (2026-01-18)**:
`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1823 token references across 85 files.`
