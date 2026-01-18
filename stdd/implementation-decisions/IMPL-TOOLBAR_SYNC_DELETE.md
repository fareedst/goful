# [IMPL:TOOLBAR_SYNC_DELETE] Toolbar Sync Delete Button

**Cross-References**: [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Add a clickable `[D]` button to the toolbar that triggers delete operations. When Linked mode is ON, triggers Sync delete (across all windows). When OFF, triggers single-window delete.

## Rationale

- Provides single-click access to delete operation with Linked-mode awareness
- Follows established toolbar pattern from existing buttons
- Reduces keystrokes compared to `S d` prefix sequence

## Implementation Approach

### Module 1: Toolbar Callback (`filer/filer.go`)

```go
// [IMPL:TOOLBAR_SYNC_DELETE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
var toolbarSyncDeleteFn func()

func SetToolbarSyncDeleteFn(fn func()) {
    toolbarSyncDeleteFn = fn
}
```

### Module 2: Button Rendering (`filer/filer.go`)

Extend `drawHeader()` to render `[D]` after `[C]`:

```go
// [IMPL:TOOLBAR_SYNC_DELETE] Draw sync delete button
syncDeleteBtn := "[D]"
syncDeleteX1 := x
x = widget.SetCells(x, y, syncDeleteBtn, look.Default())
syncDeleteX2 := x - 1
toolbarButtons["syncdelete"] = toolbarBounds{x1: syncDeleteX1, y: y, x2: syncDeleteX2}
x = widget.SetCells(x, y, " ", look.Default())
```

### Module 3: Button Invocation (`filer/filer.go`)

Extend `InvokeToolbarButton()`:

```go
case "syncdelete":
    if toolbarSyncDeleteFn != nil {
        toolbarSyncDeleteFn()
        return true
    }
```

### Module 4: Callback Wiring (`main.go`)

```go
// [IMPL:TOOLBAR_SYNC_DELETE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
filer.SetToolbarSyncDeleteFn(func() {
    if g.IsLinkedNav() {
        // Sync delete - operates across all windows
        file := g.File()
        if file == nil || file.Name() == ".." {
            message.Errorf("no file selected")
            return
        }
        g.startSyncDelete(file.Name(), g.IsSyncIgnoreFailures())
    } else {
        // Single-window delete
        g.Remove()
    }
})
```

## Code Markers

Files that must carry `[IMPL:TOOLBAR_SYNC_DELETE]` annotations:
- `filer/filer.go` - `drawHeader()`, `SetToolbarSyncDeleteFn()`, `InvokeToolbarButton()`
- `main.go` - callback wiring

Tests that must reference `[REQ:TOOLBAR_SYNC_BUTTONS]`:
- `filer/toolbar_test.go` - sync delete button hit-testing and invocation tests

## Token Coverage `[PROC:TOKEN_AUDIT]`

- [ ] `filer/filer.go` `drawHeader()` - `[IMPL:TOOLBAR_SYNC_DELETE]`
- [ ] `filer/filer.go` `SetToolbarSyncDeleteFn()` - `[IMPL:TOOLBAR_SYNC_DELETE]`
- [ ] `filer/filer.go` `InvokeToolbarButton()` - `[IMPL:TOOLBAR_SYNC_DELETE]`
- [ ] `main.go` callback wiring - `[IMPL:TOOLBAR_SYNC_DELETE]`
- [ ] `TestToolbarSyncDeleteButtonHit_REQ_TOOLBAR_SYNC_BUTTONS`
- [ ] `TestInvokeToolbarSyncDeleteButton_REQ_TOOLBAR_SYNC_BUTTONS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1621 token references across 79 files. |

## Related Decisions

- Depends on: [IMPL:TOOLBAR_COMPARE_BUTTON], [IMPL:SYNC_EXECUTE]
- See also: [IMPL:TOOLBAR_SYNC_COPY], [IMPL:TOOLBAR_SYNC_RENAME], [IMPL:TOOLBAR_IGNORE_FAILURES]

---

*Created on 2026-01-18*
