# [IMPL:DIFF_SEARCH] Difference Search Implementation

**Cross-References**: [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]  
**Status**: Active  
**Created**: 2026-01-10  
**Last Updated**: 2026-01-17

---

## Decision

Implement a two-command difference search with state in Workspace, pure comparison logic in a dedicated module, and a persistent status display.

## Rationale

- Separating state (initial dirs, active flag, status fields) from comparison logic keeps the code testable per `[REQ:MODULE_VALIDATION]`
- Using cursor position as the implicit bookmark for "Continue" simplifies state and integrates naturally with existing navigation
- Alphabetic iteration (case-sensitive) matches user expectations for file sorting
- A dedicated persistent status line provides continuous feedback during long searches without auto-dismissing like regular messages

## Implementation Approach

### State Management (`filer/diffsearch.go`)

- Add `DiffSearchState` struct with `InitialDirs []string`, `Active bool`, and status fields:
  - `LastDiffName string` - Name of last found difference
  - `LastDiffReason string` - Reason for last difference
  - `CurrentPath string` - Current directory being searched
  - `FilesChecked int` - Count of files checked
  - `Searching bool` - Whether actively searching vs paused
- Add setter methods: `SetSearching()`, `SetCurrentPath()`, `IncrementFilesChecked()`, `SetLastDiff()`
- Add `StatusText() string` that returns formatted status for display
- Store state in `Workspace` struct as `diffSearch *DiffSearchState`

### Core Comparison Logic (`filer/diffsearch.go`)

- Add `func CollectAllNames(dirs []*Directory) []string` that returns the union of all file/directory names (excluding `..`), sorted alphabetically
- Add `func CheckDifference(name string, dirs []*Directory) (isDiff bool, reason string)`:
  - Check presence in each directory
  - If missing from any, return `true, "missing in window N"`
  - If present in all, compare sizes. If sizes differ, return `true, "size mismatch"`
  - Otherwise return `false, ""`
- Add `func FindNextDifference(dirs []*Directory, startAfter string) (name string, reason string, found bool)`:
  - Call `CollectAllNames`, iterate from after `startAfter`
  - Return first difference found

### Subdirectory Descent (`filer/diffsearch.go`)

- Add `func FindNextSubdir(dirs []*Directory, startAfter string) (name string, existsInAll bool, found bool)`
- Add `func FindNextSubdirInAll(dirs []*Directory, startAfter string) (name string, found bool)`:
  - Iterates through subdirectories in alphabetical order starting after `startAfter`
  - Returns the first subdirectory that exists in ALL directories
  - Critical for maintaining search position when user manually navigates into subdirectories

### Navigator Interface (Refactored 2026-01-13)

Interface `Navigator` abstracts directory operations for tree traversal:
- `GetDirs() []*Directory`
- `ChdirAll(name string)`
- `ChdirParentAll()`
- `CurrentPath() string`
- `RebuildComparisonIndex()`

### TreeWalker (Added 2026-01-13)

- `type TreeWalker struct` - Handles tree traversal algorithm
- `NewTreeWalker(nav Navigator, state *DiffSearchState, startAfter string)`
- `Run(progressFn func()) Step` - Executes traversal, returns Step result
- `type Step struct` with Type, Name, Reason, IsDir fields

### Persistent Status Display (`diffstatus/diffstatus.go`)

New package providing a dedicated status line that persists while diff search is active:
- `Init()`, `SetStatusFn()`, `SetActiveFn()`, `SetMessage()`, `ClearMessage()`, `IsActive()`, `Draw()`, `Resize()`

### Command Wrappers (`app/goful.go`)

- `StartDiffSearch()`, `ContinueDiffSearch()`, `DiffSearchStatus()`, `IsDiffSearchActive()`

### Keymap Integration

- `[` for start diff search, `]` for continue diff search

## Code Markers

- `filer/diffsearch.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`
- `filer/workspace.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`
- `app/goful.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`
- `diffstatus/diffstatus.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`
- `main.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/diffsearch.go` - all functions
- [ ] `filer/workspace.go` - cursor helpers
- [ ] `app/goful.go` - command wrappers
- [ ] `diffstatus/diffstatus.go` - all functions
- [ ] `main.go` - keybindings

Tests that must reference `[REQ:DIFF_SEARCH]`:
- [ ] `filer/diffsearch_test.go` with tests named `Test*_REQ_DIFF_SEARCH`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-10 | — | ✅ Pass | 16 tests passing |
| 2026-01-10 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 798 token references across 69 files |

### Bug Fixes

- **2026-01-10**: Added `FindNextSubdirInAll` to respect `startAfter` during subdirectory descent
- **2026-01-12**: Fixed `startAfter` assignment after ascending from child directories
- **2026-01-13**: Fixed alphabetical comparison to use `name <= startAfter` instead of exact match

## Related Decisions

- Depends on: —
- See also: [ARCH:DIFF_SEARCH], [REQ:DIFF_SEARCH], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
