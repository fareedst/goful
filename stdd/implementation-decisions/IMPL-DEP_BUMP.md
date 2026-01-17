# [IMPL:DEP_BUMP] Dependency Bump

**Cross-References**: [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Refresh direct deps and tidy module graph.

## Rationale

- Pull in security and bug fixes; keep compatible with Go LTS

## Implementation Approach

- Upgraded direct deps to current releases:
  - `github.com/gdamore/tcell/v2 v2.13.5`
  - `github.com/mattn/go-runewidth v0.0.19`
  - `github.com/google/shlex` (latest pseudo-version)
- Upgraded supporting deps:
  - `github.com/lucasb-eyer/go-colorful v1.3.0`
  - `github.com/gdamore/encoding v1.0.1`
  - `github.com/rivo/uniseg v0.4.7`
  - `golang.org/x/sys v0.39.0`
  - `golang.org/x/term v0.38.0`
  - `golang.org/x/text v0.32.0`
  - Added `github.com/clipperhouse/uax29/v2 v2.2.0` via transitive requirements from `tcell`
- Ran `go mod tidy` to sync `go.sum`
- No shims or breaking API adjustments were required after the upgrades

## Code Markers

- `go.mod` entries and related code comments

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] Dependency change commits include `[IMPL:DEP_BUMP] [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]`

Tests that must reference `[REQ:DEPENDENCY_REFRESH]`:
- [ ] N/A - verified via successful test suite

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | All tests pass with updated deps |

## Related Decisions

- Depends on: [IMPL:GO_MOD_UPDATE]
- See also: [ARCH:DEPENDENCY_POLICY], [REQ:DEPENDENCY_REFRESH]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
