# [IMPL:FILE_COMPARISON_INDEX] File Comparison Index

**Cross-References**: [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]  
**Status**: Active  
**Created**: 2026-01-09  
**Last Updated**: 2026-01-17

---

## Decision

Build a cached index of cross-directory file comparison states for O(1) draw-time lookup.

## Rationale

- Comparison must not block directory reading or initial display
- Index built once after all directories load, cached until invalidation events
- Pure function design enables independent module validation

## Implementation Approach

Add `filer/compare.go` with:

- `type CompareState struct { NamePresent bool; SizeState SizeCompare; TimeState TimeCompare }` where `SizeCompare` and `TimeCompare` are enums (Equal, Smallest, Largest / Equal, Earliest, Latest)
- `type ComparisonIndex struct` with `cache map[string]map[int]CompareState` keyed by filename then dirIndex
- `func BuildIndex(dirs []*Directory) *ComparisonIndex` that:
  - Collects all files by name across directories
  - For files in multiple directories: marks NamePresent=true, computes size/time comparisons
  - For single-directory files: no entry (returns nil on lookup)
- `func (idx *ComparisonIndex) Get(dirIndex int, filename string) *CompareState` for draw-time lookup
- `var comparisonEnabled bool` and `func ToggleComparisonColors() (enabled bool)` for runtime toggle

Workspace tracks `*ComparisonIndex` and rebuilds on invalidation events (`Chdir`, `reload`, `ReloadAll`, `CreateDir`, `CloseDir`).

Index building happens after `ReloadAll` completes, before next draw cycle.

## Code Markers

- `filer/compare.go` includes `[IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/compare.go` - index building and lookup
- [ ] `filer/compare_test.go` - tests

Tests that must reference `[REQ:FILE_COMPARISON_COLORS]`:
- [ ] Index building tests with various window/file combinations
- [ ] Size/time edge case tests

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-09 | — | ✅ Pass | `go test ./filer` |

## Related Decisions

- Depends on: [IMPL:COMPARE_COLOR_CONFIG]
- See also: [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
