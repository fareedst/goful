# [IMPL:FILER_EXCLUDE_LOADER] Filename Exclude Loader & Toggle

**Cross-References**: [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]  
**Status**: Active  
**Created**: 2026-01-07  
**Last Updated**: 2026-01-17

---

## Decision

Extend the existing path resolver with `-exclude-names` / `GOFUL_EXCLUDES_FILE` and wire a loader + toggle UI hook into `main`.

## Rationale

- Reuses the proven precedence model so operators immediately understand how to override the list location
- Keeps parsing logic (trim, comment skip, case normalization) pure and testable
- Provides a discoverable runtime toggle via both the View menu and a dedicated keystroke so users can quickly inspect hidden files when necessary

## Implementation Approach

### Resolver Extension

`configpaths/resolver.go` adds `DefaultExcludesPath`, `EnvExcludesKey`, and `Excludes`/`ExcludesSource` fields on `Paths`, plus a new `flagExcludesSourceLabel`. `Resolver.Resolve` now accepts `flagExcludes` and returns the resolved path + provenance so `emitPathDebug` can log it.

### Loader

`main.go` defines `excludeNamesFlag`, calls `pathsResolver.Resolve(*stateFlag, *historyFlag, *commandsFlag, *excludeNamesFlag)`, and invokes `loadExcludedNames(paths.Excludes)` before `app.SeedStartupWorkspaces`.

`loadExcludedNames` (new helper in `main.go`) opens the file (tolerates `os.ErrNotExist`), reads newline-delimited basenames, strips comments (`#` prefix) and whitespace, lowercases entries, and calls `filer.ConfigureExcludedNames(parsed, true)`. Errors use `message.Errorf` referencing `[REQ:FILER_EXCLUDE_NAMES]`; success paths log `message.Infof` counts.

### Toggle

Toggle helper `toggleExcludedNames(g *app.Goful)` wraps `filer.ToggleExcludedNames`, reports the new state/count via `message.Infof`, and calls `g.Workspace().ReloadAll()`. Bound to `g.AddKeymap("E", toggle)` and added as `view` menu entry (e.g., `n` for "toggle excludes") so the action is reachable via mouse/keyboard menus.

`emitPathDebug` gains the excludes tuple so `GOFUL_DEBUG_PATHS=1` prints provenance for the new file.

## Code Markers

- `main.go`, `configpaths/resolver.go`, and `configpaths/resolver_test.go` reference `[IMPL:FILER_EXCLUDE_LOADER] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]`
- Toggle handlers include `[REQ:FILER_EXCLUDE_NAMES]` in logged messages for auditability

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `main.go` - loader and toggle
- [ ] `configpaths/resolver.go` - excludes path fields
- [ ] `configpaths/resolver_test.go` - precedence tests

Tests that must reference `[REQ:FILER_EXCLUDE_NAMES]`:
- [ ] Loader behavior tests
- [ ] Resolver precedence tests

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-07 | — | ✅ Pass | `go test ./configpaths` (darwin/arm64, Go 1.24.3) |
| 2026-01-07 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 302 token references across 58 files |

## Related Decisions

- Depends on: [IMPL:FILER_EXCLUDE_RULES]
- See also: [ARCH:FILER_EXCLUDE_FILTER], [REQ:FILER_EXCLUDE_NAMES], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
