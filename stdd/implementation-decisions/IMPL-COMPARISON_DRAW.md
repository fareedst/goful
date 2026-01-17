# [IMPL:COMPARISON_DRAW] Comparison Draw Integration

**Cross-References**: [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]  
**Status**: Active  
**Created**: 2026-01-09  
**Last Updated**: 2026-01-17

---

## Decision

Extend `FileStat.Draw()` to accept comparison context and apply colors independently to name, size, and time fields.

## Rationale

- Minimal change to existing draw path—comparison is optional context
- Independent color application means name, size, and time can each show different comparison states
- Respects existing file-type colors when comparison is disabled or file is unique
- **Default State (2026-01-11)**: Comparison coloring is **enabled by default** (`comparisonEnabled = true` in `look/comparison.go`) so users immediately see color-coded file listings without manual toggle. Press `` ` `` (backtick) to disable if desired

## Implementation Approach

Add `look/comparison.go` with:

- Thread-safe style storage for each comparison state
- `func CompareNamePresent() tcell.Style`, `func CompareSizeEqual() tcell.Style`, etc.
- `func ConfigureComparisonColors(cfg *comparecolors.Config)` to apply loaded config

Modify `FileStat.Draw(x, y, width int, focus bool)` signature to accept optional `*CompareState`:

- New signature: `Draw(x, y, width int, focus bool, cmp *CompareState)`
- When `cmp != nil` and `cmp.NamePresent`: use comparison name color
- Size field uses `cmp.SizeState` to select Equal/Smallest/Largest color
- Time field uses `cmp.TimeState` to select Equal/Earliest/Latest color

`Directory.drawFiles()` passes comparison state from workspace index to each `FileStat.Draw()`.

`Workspace.Draw()` ensures index is available before drawing.

## Code Markers

- `filer/file.go`, `filer/directory.go`, `filer/workspace.go` include `[IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`
- `look/comparison.go` includes `[IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/file.go` - Draw method
- [ ] `filer/directory.go` - drawFiles
- [ ] `filer/workspace.go` - Draw coordination
- [ ] `look/comparison.go` - style functions

Tests that must reference `[REQ:FILE_COMPARISON_COLORS]`:
- [ ] Integration tests validating comparison colors apply correctly

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-11 | — | ✅ Pass | Comparison coloring enabled by default |

## Related Decisions

- Depends on: [IMPL:FILE_COMPARISON_INDEX], [IMPL:COMPARE_COLOR_CONFIG]
- See also: [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
