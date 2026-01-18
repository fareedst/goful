# [IMPL:TOOLBAR_IGNORE_FAILURES] Toolbar Ignore Failures Toggle

**Cross-References**: [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Add a clickable `[!]` button that toggles a persistent ignore-failures mode for Sync operations. Display reverse style when ON, normal when OFF.

## Rationale

- Replaces the transient `!` toggle within SyncMode
- Provides visual indication of current ignore-failures state
- Allows setting the mode before triggering operations

## Implementation Approach

### Module 1: State Management (`app/goful.go`)

```go
type Goful struct {
    // ... existing fields ...
    syncIgnoreFailures bool // [IMPL:TOOLBAR_IGNORE_FAILURES] Persistent ignore-failures mode
}

// [IMPL:TOOLBAR_IGNORE_FAILURES] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
func (g *Goful) IsSyncIgnoreFailures() bool {
    return g.syncIgnoreFailures
}

// [IMPL:TOOLBAR_IGNORE_FAILURES] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
func (g *Goful) ToggleSyncIgnoreFailures() bool {
    g.syncIgnoreFailures = !g.syncIgnoreFailures
    return g.syncIgnoreFailures
}
```

### Module 2: Indicator Callback (`filer/filer.go`)

```go
// [IMPL:TOOLBAR_IGNORE_FAILURES] Returns whether ignore-failures mode is enabled
var syncIgnoreFailuresIndicator func() bool

func SetSyncIgnoreFailuresIndicator(fn func() bool) {
    syncIgnoreFailuresIndicator = fn
}
```

### Module 3: Button Rendering (`filer/filer.go`)

```go
// [IMPL:TOOLBAR_IGNORE_FAILURES] Draw ignore-failures toggle
ignoreBtnLabel := "[!]"
ignoreBtnX1 := x
ignoreStyle := look.Default()
if syncIgnoreFailuresIndicator != nil && syncIgnoreFailuresIndicator() {
    ignoreStyle = ignoreStyle.Reverse(true)
}
x = widget.SetCells(x, y, ignoreBtnLabel, ignoreStyle)
ignoreBtnX2 := x - 1
toolbarButtons["ignorefailures"] = toolbarBounds{x1: ignoreBtnX1, y: y, x2: ignoreBtnX2}
x = widget.SetCells(x, y, " ", look.Default())
```

### Module 4: Callback Wiring (`main.go`)

```go
// [IMPL:TOOLBAR_IGNORE_FAILURES] Set indicator for button styling
filer.SetSyncIgnoreFailuresIndicator(g.IsSyncIgnoreFailures)

// [IMPL:TOOLBAR_IGNORE_FAILURES] Toggle callback
filer.SetToolbarIgnoreFailuresFn(func() {
    enabled := g.ToggleSyncIgnoreFailures()
    if enabled {
        message.Info("Sync ignore-failures mode ON")
    } else {
        message.Info("Sync ignore-failures mode OFF")
    }
})
```

## Code Markers

Files that must carry `[IMPL:TOOLBAR_IGNORE_FAILURES]` annotations:
- `app/goful.go` - `syncIgnoreFailures`, `IsSyncIgnoreFailures()`, `ToggleSyncIgnoreFailures()`
- `filer/filer.go` - `drawHeader()`, `SetToolbarIgnoreFailuresFn()`, `SetSyncIgnoreFailuresIndicator()`, `InvokeToolbarButton()`
- `main.go` - callback wiring

## Token Coverage `[PROC:TOKEN_AUDIT]`

- [ ] `app/goful.go` state and accessors - `[IMPL:TOOLBAR_IGNORE_FAILURES]`
- [ ] `filer/filer.go` rendering and callbacks - `[IMPL:TOOLBAR_IGNORE_FAILURES]`
- [ ] `main.go` wiring - `[IMPL:TOOLBAR_IGNORE_FAILURES]`
- [ ] `TestToolbarIgnoreFailuresButtonHit_REQ_TOOLBAR_SYNC_BUTTONS`
- [ ] `TestInvokeToolbarIgnoreFailuresButton_REQ_TOOLBAR_SYNC_BUTTONS`
- [ ] `TestIgnoreFailuresIndicator_REQ_TOOLBAR_SYNC_BUTTONS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1621 token references across 79 files. |

## Related Decisions

- Depends on: [IMPL:TOOLBAR_COMPARE_BUTTON]
- See also: [IMPL:TOOLBAR_SYNC_COPY], [IMPL:TOOLBAR_SYNC_DELETE], [IMPL:TOOLBAR_SYNC_RENAME]

---

*Created on 2026-01-18*
