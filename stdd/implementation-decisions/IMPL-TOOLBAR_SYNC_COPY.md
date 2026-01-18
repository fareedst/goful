# [IMPL:TOOLBAR_SYNC_COPY] Toolbar Sync Copy Button

**Cross-References**: [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Add a clickable `[C]` button to the toolbar that triggers copy operations. When Linked mode is ON, triggers Sync copy (across all windows). When OFF, triggers single-window copy.

## Rationale

- Provides single-click access to copy operation with Linked-mode awareness
- Follows established toolbar pattern from existing buttons
- Reduces keystrokes compared to `S c` prefix sequence

## Implementation Approach

### Module 1: Toolbar Callback (`filer/filer.go`)

```go
// [IMPL:TOOLBAR_SYNC_COPY] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
var toolbarSyncCopyFn func()

func SetToolbarSyncCopyFn(fn func()) {
    toolbarSyncCopyFn = fn
}
```

### Module 2: Button Rendering (`filer/filer.go`)

Extend `drawHeader()` to render `[C]` after `[=]`:

```go
// [IMPL:TOOLBAR_SYNC_COPY] Draw sync copy button
syncCopyBtn := "[C]"
syncCopyX1 := x
x = widget.SetCells(x, y, syncCopyBtn, look.Default())
syncCopyX2 := x - 1
toolbarButtons["synccopy"] = toolbarBounds{x1: syncCopyX1, y: y, x2: syncCopyX2}
x = widget.SetCells(x, y, " ", look.Default())
```

### Module 3: Button Invocation (`filer/filer.go`)

Extend `InvokeToolbarButton()`:

```go
case "synccopy":
    if toolbarSyncCopyFn != nil {
        toolbarSyncCopyFn()
        return true
    }
```

### Module 4: Callback Wiring (`main.go`)

```go
// [IMPL:TOOLBAR_SYNC_COPY] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
filer.SetToolbarSyncCopyFn(func() {
    if g.IsLinkedNav() {
        // Sync copy - operates across all windows
        file := g.File()
        if file == nil || file.Name() == ".." {
            message.Errorf("no file selected")
            return
        }
        g.startSyncCopy(file.Name(), g.IsSyncIgnoreFailures())
    } else {
        // Single-window copy
        g.Copy()
    }
})
```

## Code Markers

Files that must carry `[IMPL:TOOLBAR_SYNC_COPY]` annotations:
- `filer/filer.go` - `drawHeader()`, `SetToolbarSyncCopyFn()`, `InvokeToolbarButton()`
- `main.go` - callback wiring

Tests that must reference `[REQ:TOOLBAR_SYNC_BUTTONS]`:
- `filer/toolbar_test.go` - sync copy button hit-testing and invocation tests

## Token Coverage `[PROC:TOKEN_AUDIT]`

- [ ] `filer/filer.go` `drawHeader()` - `[IMPL:TOOLBAR_SYNC_COPY]`
- [ ] `filer/filer.go` `SetToolbarSyncCopyFn()` - `[IMPL:TOOLBAR_SYNC_COPY]`
- [ ] `filer/filer.go` `InvokeToolbarButton()` - `[IMPL:TOOLBAR_SYNC_COPY]`
- [ ] `main.go` callback wiring - `[IMPL:TOOLBAR_SYNC_COPY]`
- [ ] `TestToolbarSyncCopyButtonHit_REQ_TOOLBAR_SYNC_BUTTONS`
- [ ] `TestInvokeToolbarSyncCopyButton_REQ_TOOLBAR_SYNC_BUTTONS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1621 token references across 79 files. |

## Related Decisions

- Depends on: [IMPL:TOOLBAR_COMPARE_BUTTON], [IMPL:SYNC_EXECUTE]
- See also: [IMPL:TOOLBAR_SYNC_DELETE], [IMPL:TOOLBAR_SYNC_RENAME], [IMPL:TOOLBAR_IGNORE_FAILURES]

---

*Created on 2026-01-18*
