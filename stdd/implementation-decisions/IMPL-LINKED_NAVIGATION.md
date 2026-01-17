# [IMPL:LINKED_NAVIGATION] Linked Navigation Implementation

**Cross-References**: [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]  
**Status**: Active  
**Created**: 2026-01-09  
**Last Updated**: 2026-01-17

---

## Decision

Implement linked navigation with minimal state in `app.Goful` and pure navigation helpers in `filer.Workspace`.

## Rationale

- Keeps the feature self-contained and easy to test independently per `[REQ:MODULE_VALIDATION]`
- Uses existing `Directory.Chdir` infrastructure rather than duplicating navigation logic
- Provides clear toggle feedback via message and header indicator

## Implementation Approach

### State Management (`app/goful.go`)

- Add `linkedNav bool` field to `Goful` struct (default `true`)
- Add `func (g *Goful) ToggleLinkedNav() bool` that flips the state and returns new value
- Add `func (g *Goful) IsLinkedNav() bool` getter
- Export `LinkedNavEnabled` callback type for header rendering

### Navigation Helpers (`filer/workspace.go`)

- Add `func (w *Workspace) ChdirAllToSubdir(name string)` that iterates all non-focused directories, checks if `name` exists as a subdirectory, and calls `Chdir(name)` if so
- Add `func (w *Workspace) ChdirAllToParent()` that iterates all directories and calls `Chdir("..")`
- Add `func (w *Workspace) SortAllBy(typ SortType)` that applies the given sort type to all directories
- All methods rebuild comparison index after changes

### Exported Sort Types (`filer/directory.go`)

- Export `SortType` and constants (`SortName`, `SortNameRev`, `SortSize`, `SortSizeRev`, `SortMtime`, `SortMtimeRev`, `SortExt`, `SortExtRev`) for linked sort synchronization
- Export `SortBy(typ SortType)` method to enable workspace-level sorting

### Header Indicator (`filer/filer.go`)

- Add `var linkedNavIndicatorFunc func() bool` package variable
- Add `func SetLinkedNavIndicator(fn func() bool)` to wire the callback from main
- Modify `drawHeader()` to show `[LINKED]` with reverse style when the callback returns true

### Keymap Integration (`main.go`)

- Replace direct navigation callbacks with wrappers that check `g.IsLinkedNav()`
- For `backspace`/`C-h`/`u`: if linked, call `g.Workspace().ChdirAllToParent()` then `g.Dir().Chdir("..")`
- For enter-dir (extmap `.dir`): if linked, call `g.Workspace().ChdirAllToSubdir(name)` then `g.Dir().EnterDir()`
- For sort menu: if linked, call `g.Workspace().SortAllBy(typ)` instead of `g.Dir().Sort*()`
- Add `L` (uppercase, macOS-compatible) and `M-l` bindings to toggle with `message.Infof` feedback
- Wire `filer.SetLinkedNavIndicator(g.IsLinkedNav)` at startup

## Code Markers

- `app/goful.go`: `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `filer/workspace.go`: `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `filer/filer.go`: `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `main.go`: `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `app/goful.go` - state and toggle
- [ ] `filer/workspace.go` - navigation helpers
- [ ] `filer/filer.go` - header indicator
- [ ] `filer/directory.go` - sort types
- [ ] `main.go` - keymap wiring

Tests that must reference `[REQ:LINKED_NAVIGATION]`:
- [ ] `TestChdirAllToSubdir_REQ_LINKED_NAVIGATION`
- [ ] `TestChdirAllToParent_REQ_LINKED_NAVIGATION`
- [ ] `TestSortAllBy_REQ_LINKED_NAVIGATION`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-09 | — | ✅ Pass | `go test ./...` (darwin/arm64, Go 1.24.3) |
| 2026-01-09 | — | ✅ Pass | token validation passed |

## Related Decisions

- Depends on: —
- See also: [ARCH:LINKED_NAVIGATION], [REQ:LINKED_NAVIGATION], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
