# [IMPL:TOOLBAR_LINKED_TOGGLE] Toolbar Linked Toggle Button

**Cross-References**: [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Add a clickable `[L]` button to the toolbar in the filer header row, immediately after the `[^]` parent button. The button displays the current linked navigation mode state (reverse style when ON, normal style when OFF) and toggles the mode when clicked. This button replaces the existing conditional `[LINKED]` indicator.

## Rationale

- The existing `[LINKED]` indicator only appears when linked mode is ON, providing no visual cue when OFF
- A toggle button serves dual purpose: state display AND mode toggle
- Consistent with the toolbar pattern established by the parent button
- Improves mouse-first user experience by providing clickable access to linked mode toggle

## Implementation Approach

### 1. Add Callback Infrastructure (`filer/filer.go`)

Add new callback variable and setter alongside the existing `toolbarParentNavFn`:

```go
// [IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]
var toolbarLinkedToggleFn func()

func SetToolbarLinkedToggleFn(fn func()) {
    toolbarLinkedToggleFn = fn
}
```

### 2. Extend Button Invocation (`filer/filer.go`)

Update `InvokeToolbarButton()` to handle the "linked" button:

```go
case "linked":
    if toolbarLinkedToggleFn != nil {
        toolbarLinkedToggleFn()
        return true
    }
```

### 3. Render Linked Button (`filer/filer.go`)

Modify `drawHeader()` to render `[L]` immediately after `[^]`:

```go
// [IMPL:TOOLBAR_LINKED_TOGGLE] Draw linked mode toggle button
linkedBtn := "[L]"
linkedX1 := x
linkedStyle := look.Default()
if linkedNavIndicator != nil && linkedNavIndicator() {
    linkedStyle = linkedStyle.Reverse(true)
}
x = widget.SetCells(x, y, linkedBtn, linkedStyle)
linkedX2 := x - 1
toolbarButtons["linked"] = toolbarBounds{x1: linkedX1, y: y, x2: linkedX2}
x = widget.SetCells(x, y, " ", look.Default())
```

### 4. Remove Conditional Indicator (`filer/filer.go`)

Remove the existing `[LINKED]` indicator code that currently appears after the workspace tabs:

```go
// REMOVE THIS BLOCK:
if linkedNavIndicator != nil && linkedNavIndicator() {
    x = widget.SetCells(x, y, "[LINKED]", look.Default().Reverse(true))
    x = widget.SetCells(x, y, " ", look.Default())
}
```

### 5. Wire Callback (`main.go`)

Add callback wiring in the `config()` function:

```go
// [IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]
filer.SetToolbarLinkedToggleFn(func() {
    enabled := g.ToggleLinkedNav()
    state := "disabled"
    if enabled {
        state = "enabled"
    }
    message.Infof("[REQ:LINKED_NAVIGATION] linked navigation %s", state)
})
```

## Visual Behavior

| Linked Mode State | Button Appearance |
|-------------------|-------------------|
| ON | `[L]` with reverse style (highlighted) |
| OFF | `[L]` with normal style (dimmed) |

## Header Layout After Implementation

```
[^] [L] 1  2  3  | [DIFF STATUS] [1] /path... [2] /path...
```

## Code Markers

Files that must carry `[IMPL:TOOLBAR_LINKED_TOGGLE]` annotations:
- `filer/filer.go` - `drawHeader()`, `SetToolbarLinkedToggleFn()`, `InvokeToolbarButton()`
- `main.go` - callback wiring in `config()`

Tests that must reference `[REQ:TOOLBAR_LINKED_TOGGLE]`:
- `filer/toolbar_test.go` - linked button hit-testing and invocation tests

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `filer/filer.go` `drawHeader()` - `[IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]`
- [x] `filer/filer.go` `SetToolbarLinkedToggleFn()` - `[IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]`
- [x] `filer/filer.go` `InvokeToolbarButton()` - `[IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]`
- [x] `main.go` callback wiring - `[IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]`

Tests that must reference `[REQ:TOOLBAR_LINKED_TOGGLE]`:
- [x] `TestToolbarLinkedButtonHit_REQ_TOOLBAR_LINKED_TOGGLE`
- [x] `TestInvokeToolbarLinkedButton_REQ_TOOLBAR_LINKED_TOGGLE`
- [x] `TestInvokeToolbarLinkedButtonWithNilCallback_REQ_TOOLBAR_LINKED_TOGGLE`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ⏳ Pending | Initial plan documented |
| 2026-01-18 | — | ✅ Pass | Token validation: 1360 refs across 78 files |

## Related Decisions

- Depends on: [IMPL:TOOLBAR_PARENT_BUTTON] (toolbar infrastructure)
- Depends on: [IMPL:LINKED_NAVIGATION] (linked mode state management)
- See also: [ARCH:TOOLBAR_LAYOUT]

---

*Created on 2026-01-18*
