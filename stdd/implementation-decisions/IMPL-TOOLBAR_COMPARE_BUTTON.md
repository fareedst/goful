# [IMPL:TOOLBAR_COMPARE_BUTTON] Toolbar Compare Digest Button

**Cross-References**: [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]  
**Status**: Active  
**Created**: 2026-01-17  
**Last Updated**: 2026-01-17

---

## Decision

Add a clickable `[=]` button to the toolbar that triggers digest comparison for all files appearing in multiple directories, with a `SharedFilenames()` method to enumerate comparable files.

## Rationale

- Single-click access to batch digest comparison improves directory sync workflows
- Follows established toolbar pattern from parent and linked buttons
- `SharedFilenames()` enables efficient iteration without exposing internal cache structure

## Implementation Approach

### Module 1: SharedFilenames (`filer/compare.go`)

Add method to enumerate files in the comparison index:

```go
// SharedFilenames returns all filenames that appear in multiple directories.
// [IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]
func (idx *ComparisonIndex) SharedFilenames() []string {
    if idx == nil {
        return nil
    }
    idx.mu.RLock()
    defer idx.mu.RUnlock()
    names := make([]string, 0, len(idx.cache))
    for name := range idx.cache {
        names = append(names, name)
    }
    return names
}
```

### Module 2: Toolbar Callback (`filer/filer.go`)

Add callback infrastructure:

```go
// [IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]
var toolbarCompareDigestFn func()

func SetToolbarCompareDigestFn(fn func()) {
    toolbarCompareDigestFn = fn
}
```

Extend `InvokeToolbarButton()` for "compare" case:

```go
case "compare":
    if toolbarCompareDigestFn != nil {
        toolbarCompareDigestFn()
        return true
    }
```

### Module 3: Button Rendering (`filer/filer.go`)

Extend `drawHeader()` to render `[=]` after `[L]`:

```go
// [IMPL:TOOLBAR_COMPARE_BUTTON] Draw compare digest button - normal style (action button)
compareBtn := "[=]"
compareX1 := x
x = widget.SetCells(x, y, compareBtn, look.Default())
compareX2 := x - 1
toolbarButtons["compare"] = toolbarBounds{x1: compareX1, y: y, x2: compareX2}
x = widget.SetCells(x, y, " ", look.Default())
```

### Module 4: Callback Wiring (`main.go`)

Wire callback that:
1. Gets `SharedFilenames()` from comparison index
2. Iterates and calls `CalculateDigestForFile()` for each
3. Displays summary message

```go
// [IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]
filer.SetToolbarCompareDigestFn(func() {
    idx := g.Workspace().ComparisonIndex()
    if idx == nil {
        message.Info("[REQ:FILE_COMPARISON_COLORS] no files to compare")
        return
    }
    names := idx.SharedFilenames()
    if len(names) == 0 {
        message.Info("[REQ:FILE_COMPARISON_COLORS] no shared filenames")
        return
    }
    totalFiles := 0
    for _, name := range names {
        count := g.Workspace().CalculateDigestForFile(name)
        totalFiles += count
    }
    message.Infof("[REQ:FILE_COMPARISON_COLORS] calculated digests for %d files across %d shared filenames", totalFiles, len(names))
})
```

## Code Markers

Files that must carry `[IMPL:TOOLBAR_COMPARE_BUTTON]` annotations:
- `filer/compare.go` - `SharedFilenames()`
- `filer/filer.go` - `drawHeader()`, `SetToolbarCompareDigestFn()`, `InvokeToolbarButton()`
- `main.go` - callback wiring

Tests that must reference `[REQ:TOOLBAR_COMPARE_BUTTON]`:
- `filer/toolbar_test.go` - compare button hit-testing and invocation tests
- `filer/compare_test.go` - `SharedFilenames()` unit tests

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `filer/compare.go` `SharedFilenames()` - `[IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]`
- [x] `filer/filer.go` `drawHeader()` - `[IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]`
- [x] `filer/filer.go` `SetToolbarCompareDigestFn()` - `[IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]`
- [x] `filer/filer.go` `InvokeToolbarButton()` - `[IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]`
- [x] `main.go` callback wiring - `[IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]`

Tests that must reference `[REQ:TOOLBAR_COMPARE_BUTTON]`:
- [x] `TestToolbarCompareButtonHit_REQ_TOOLBAR_COMPARE_BUTTON`
- [x] `TestInvokeToolbarCompareButton_REQ_TOOLBAR_COMPARE_BUTTON`
- [x] `TestInvokeToolbarCompareButtonWithNilCallback_REQ_TOOLBAR_COMPARE_BUTTON`
- [x] `TestSharedFilenames_NilIndex_REQ_TOOLBAR_COMPARE_BUTTON`
- [x] `TestSharedFilenames_EmptyIndex_REQ_TOOLBAR_COMPARE_BUTTON`
- [x] `TestSharedFilenames_WithFiles_REQ_TOOLBAR_COMPARE_BUTTON`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-17 | — | ✅ Pass | Token validation: 1468 refs across 78 files |

## Related Decisions

- Depends on: [IMPL:TOOLBAR_LINKED_TOGGLE], [IMPL:DIGEST_COMPARISON], [IMPL:FILE_COMPARISON_INDEX]
- See also: [ARCH:FILE_COMPARISON_ENGINE]

---

*Created on 2026-01-17*
