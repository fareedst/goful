# [IMPL:DIGEST_COMPARISON] Digest Comparison

**Cross-References**: [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]  
**Status**: Active  
**Created**: 2026-01-11  
**Last Updated**: 2026-01-17

---

## Decision

On-demand xxHash64 digest calculation for files with equal sizes across directories.

## Rationale

- Files may have identical sizes but different content; digest comparison provides content verification
- On-demand calculation avoids I/O overhead for files the user doesn't need to compare
- xxHash64 offers excellent speed for non-cryptographic hashing, suitable for file comparison
- Terminal attributes (underline/strikethrough) provide visual distinction without adding new color configuration

## Implementation Approach

- Add `DigestCompare` enum to `filer/compare.go` with states: `DigestUnknown`, `DigestEqual`, `DigestDifferent`, `DigestNA`
- Add `DigestState` field to `CompareState` struct
- Implement `CalculateFileDigest(path string) (uint64, error)` using `github.com/cespare/xxhash/v2` with streaming I/O
- Add `UpdateDigestStates(filename string, dirs []*Directory) int` method to `ComparisonIndex`:
  - Only processes files with `SizeState == SizeEqual`
  - Calculates digests for all matching files across directories
  - Sets `DigestState` to `DigestEqual` if all digests match, `DigestDifferent` otherwise
- Add `CalculateDigestForFile(filename string) int` method to `Workspace` as public API
- Modify `FileStat.DrawWithComparison()` to apply terminal attributes to size field:
  - `DigestEqual`: `tcell.AttrUnderline`
  - `DigestDifferent`: `tcell.AttrStrikeThrough`
- Bind `=` keystroke to trigger digest calculation for the file under cursor
- Add "calculate file digest" entry to View menu for discoverability

## Code Markers

- `filer/compare.go`: `[IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`
- `filer/workspace.go`: `[IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`
- `filer/file.go`: `[IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`
- `main.go`: `[IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/compare.go` - digest enum and calculation
- [ ] `filer/workspace.go` - public API
- [ ] `filer/file.go` - draw attributes
- [ ] `main.go` - keybinding

Tests that must reference `[REQ:FILE_COMPARISON_COLORS]`:
- [ ] Digest calculation and state propagation tests

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-11 | — | ✅ Pass | Digest comparison working |

## Related Decisions

- Depends on: [IMPL:FILE_COMPARISON_INDEX], [IMPL:COMPARISON_DRAW]
- See also: [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
