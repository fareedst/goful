# [IMPL:FILER_EXCLUDE_RULES] Filename Exclude Rules

**Cross-References**: [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]  
**Status**: Active  
**Created**: 2026-01-07  
**Last Updated**: 2026-01-17

---

## Decision

Centralize basename filtering inside `filer` so every reader (default, glob, finder) automatically skips excluded entries.

## Rationale

- Keeps the filter logic deterministic, testable, and shareable across `Directory` instances without duplicating conditionals
- Supports `[REQ:MODULE_VALIDATION]` by isolating the rule store from UI wiring, enabling pure unit tests for toggle/state transitions
- Ensures mark, finder, and macro flows inherit the same behaviour because they all append via `Directory.read`

## Implementation Approach

Add `filer/exclude.go` with:

- `type excludeSet map[string]struct{}` stored alongside `excludedNamesMu sync.RWMutex`, `excludedNames excludeSet`, and `excludeEnabled bool`
- `func ConfigureExcludedNames(names []string, activate bool)` that trims whitespace, lowercases entries, replaces the set, and toggles `excludeEnabled = activate && len(set) > 0`
- `func ToggleExcludedNames() (enabled bool, hasRules bool)` plus helpers `ExcludedNamesEnabled()` and `ExcludedNameCount()` for diagnostics/UI integration
- `func shouldExclude(name string) bool` used by `Directory.read` (skips once `excludeEnabled` is true and the lowercase basename exists in the set)

Guard mark insertion: the callback inside `Directory.read` checks `shouldExclude(fs.Name())` before `AppendList`, so `defaultReader`, `globPattern`, `globDirPattern`, and finder flows automatically inherit the filter.

Emit `DEBUG: [IMPL:FILER_EXCLUDE_RULES] ...` logs when `ConfigureExcludedNames` replaces the set or when toggles occur with `GOFUL_DEBUG_PATHS` to ease troubleshooting.

## Code Markers

- `filer/exclude.go`, `filer/directory.go`, and filer tests include `[IMPL:FILER_EXCLUDE_RULES] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/exclude.go` - all functions
- [ ] `filer/directory.go` - read function guards

Tests that must reference `[REQ:FILER_EXCLUDE_NAMES]`:
- [ ] `TestExcludedNamesHidden_REQ_FILER_EXCLUDE_NAMES`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-07 | — | ✅ Pass | `go test ./filer` (darwin/arm64, Go 1.24.3) |
| 2026-01-07 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 302 token references across 58 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:FILER_EXCLUDE_FILTER], [REQ:FILER_EXCLUDE_NAMES], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
