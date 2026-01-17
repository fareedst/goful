# [IMPL:STATE_PATH_RESOLVER] State Path Resolver

**Cross-References**: [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Provide a pure resolver plus bootstrap glue for persistence paths.

## Rationale

- Implements the precedence + debug contract from [ARCH:STATE_PATH_SELECTION] without tying tests to global process state
- Makes it trivial to inject fake environments for module validation required by [REQ:MODULE_VALIDATION]

## Implementation Approach

- Add package `configpaths` with:
  - `const DefaultState = "~/.goful/state.json"` / `DefaultHistory = "~/.goful/history/shell"`
  - `const EnvStateKey = "GOFUL_STATE_PATH"` / `EnvHistoryKey = "GOFUL_HISTORY_PATH"`
  - `type Paths struct { State, History, StateSource, HistorySource string }`
  - `type Resolver struct { LookupEnv func(string) (string, bool) }` with method `Resolve(flagState, flagHistory string) Paths`
  - Resolver order: CLI flag → env var → default. All outputs pass through `util.ExpandPath`
  - `func EnsureParent(path string) error` helper to call `os.MkdirAll(filepath.Dir(path), 0o755)` before state/history saves
- Add `BootstrapPaths` helper (same package or `main.go`) that:
  - Parses CLI flags (`-state`, `-history`)
  - Calls resolver and logs `DEBUG: [IMPL:STATE_PATH_RESOLVER] ...` lines when `GOFUL_DEBUG_PATHS=1`
  - Applies resolved paths to `app.NewGoful`, `cmdline.LoadHistory`, and the corresponding save paths when exiting
- Update `filer.SaveState` to create parent directories before writing to satisfy the requirement

## Code Markers

- Resolver + helper functions carry `[IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]` comments
- `main.go` flag definitions include inline references to the same tokens

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `configpaths/*.go` - `[IMPL:STATE_PATH_RESOLVER]`
- [ ] `configpaths/*_test.go` - `[REQ:CONFIGURABLE_STATE_PATHS]`
- [ ] `main.go` - flag definitions
- [ ] `filer/filer.go` - SaveState updates
- [ ] `README.md` - documentation

Tests that must reference `[REQ:CONFIGURABLE_STATE_PATHS]`:
- [ ] `TestResolvePaths_REQ_CONFIGURABLE_STATE_PATHS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 70 token references across 40 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:STATE_PATH_SELECTION], [REQ:CONFIGURABLE_STATE_PATHS]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
