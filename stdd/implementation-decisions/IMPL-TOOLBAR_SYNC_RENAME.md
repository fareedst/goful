# [IMPL:TOOLBAR_SYNC_RENAME] Toolbar Sync Rename Button

**Cross-References**: [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Add a clickable `[R]` button to the toolbar that triggers rename operations. When Linked mode is ON, triggers Sync rename (across all windows). When OFF, triggers single-window rename.

## Rationale

- Provides single-click access to rename operation with Linked-mode awareness
- Follows established toolbar pattern from existing buttons
- Reduces keystrokes compared to `S r` prefix sequence

## Implementation Approach

### Module 1: Toolbar Callback (`filer/filer.go`)

```go
// [IMPL:TOOLBAR_SYNC_RENAME] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
var toolbarSyncRenameFn func()

func SetToolbarSyncRenameFn(fn func()) {
    toolbarSyncRenameFn = fn
}
```

### Module 2: Button Rendering (`filer/filer.go`)

Extend `drawHeader()` to render `[R]` after `[D]`:

```go
// [IMPL:TOOLBAR_SYNC_RENAME] Draw sync rename button
syncRenameBtn := "[R]"
syncRenameX1 := x
x = widget.SetCells(x, y, syncRenameBtn, look.Default())
syncRenameX2 := x - 1
toolbarButtons["syncrename"] = toolbarBounds{x1: syncRenameX1, y: y, x2: syncRenameX2}
x = widget.SetCells(x, y, " ", look.Default())
```

### Module 3: Button Invocation (`filer/filer.go`)

Extend `InvokeToolbarButton()`:

```go
case "syncrename":
    if toolbarSyncRenameFn != nil {
        toolbarSyncRenameFn()
        return true
    }
```

### Module 4: Callback Wiring (`main.go`)

```go
// [IMPL:TOOLBAR_SYNC_RENAME] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
filer.SetToolbarSyncRenameFn(func() {
    if g.IsLinkedNav() {
        // Sync rename - operates across all windows
        file := g.File()
        if file == nil || file.Name() == ".." {
            message.Errorf("no file selected")
            return
        }
        g.startSyncRename(file.Name(), g.IsSyncIgnoreFailures())
    } else {
        // Single-window rename
        g.Rename()
    }
})
```

## Code Markers

Files that must carry `[IMPL:TOOLBAR_SYNC_RENAME]` annotations:
- `filer/filer.go` - `drawHeader()`, `SetToolbarSyncRenameFn()`, `InvokeToolbarButton()`
- `main.go` - callback wiring

Tests that must reference `[REQ:TOOLBAR_SYNC_BUTTONS]`:
- `filer/toolbar_test.go` - sync rename button hit-testing and invocation tests

## Token Coverage `[PROC:TOKEN_AUDIT]`

- [ ] `filer/filer.go` `drawHeader()` - `[IMPL:TOOLBAR_SYNC_RENAME]`
- [ ] `filer/filer.go` `SetToolbarSyncRenameFn()` - `[IMPL:TOOLBAR_SYNC_RENAME]`
- [ ] `filer/filer.go` `InvokeToolbarButton()` - `[IMPL:TOOLBAR_SYNC_RENAME]`
- [ ] `main.go` callback wiring - `[IMPL:TOOLBAR_SYNC_RENAME]`
- [ ] `TestToolbarSyncRenameButtonHit_REQ_TOOLBAR_SYNC_BUTTONS`
- [ ] `TestInvokeToolbarSyncRenameButton_REQ_TOOLBAR_SYNC_BUTTONS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1621 token references across 79 files. |

## Related Decisions

- Depends on: [IMPL:TOOLBAR_COMPARE_BUTTON], [IMPL:SYNC_EXECUTE]
- See also: [IMPL:TOOLBAR_SYNC_COPY], [IMPL:TOOLBAR_SYNC_DELETE], [IMPL:TOOLBAR_IGNORE_FAILURES]

---

*Created on 2026-01-18*
