# [IMPL:WORKSPACE_START_DIRS] Startup Directory Parser & Seeder

**Cross-References**: [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]  
**Status**: Active  
**Created**: 2026-01-07  
**Last Updated**: 2026-01-17

---

## Decision

Implement helpers that convert positional CLI arguments into deterministic workspace layouts before the UI loop runs.

## Rationale

- Keeps startup customization encapsulated, making it easy to validate and maintain without scattering logic across `main.go` and filer internals
- Provides clear debug output (`GOFUL_DEBUG_WORKSPACE=1`) so operators and CI scripts can diagnose mismatched directories quickly
- Honors `[REQ:MODULE_VALIDATION]` by keeping the parser pure and the seeder isolated from runtime event handling

## Implementation Approach

### Parser: `ParseStartupDirs(args []string) ([]string, []string)`

Located in `app/startup_dirs.go`:
- Trims whitespace, expands `~`, resolves absolute clean paths, and checks existence + directory-ness via `os.Stat`
- Returns ordered directories (duplicates allowed intentionally) plus warnings describing invalid entries; warnings are surfaced through `message.Errorf`

### Seeder: `SeedStartupWorkspaces(g *app.Goful, dirs []string, debug bool) bool`

- Early-exits when no directories are provided to preserve historical state restoration
- Adds or removes workspaces to match the requested count by calling existing `CreateWorkspace` / `CloseWorkspace` helpers
- For each pane, focuses the first directory, calls `Dir().Chdir()` + `ReloadAll()`, retitles the workspace, and optionally logs `DEBUG:` entries tagged with `[IMPL:WORKSPACE_START_DIRS]`
- Returns a boolean indicating whether seeding occurred so callers can decide whether additional fallback work is needed

### main.go Integration

- After parsing flags and loading history, `flag.Args()` are fed into the parser, warnings produce `message.Errorf` output, and seeding runs with debug mode tied to `GOFUL_DEBUG_WORKSPACE`

## Code Markers

- `app/startup_dirs.go`, `app/startup_dirs_test.go`, and the new block in `main.go` include `[IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]` comments

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `main.go` - startup dirs wiring
- [ ] `app/startup_dirs.go` - parser and seeder

Tests that must reference `[REQ:WORKSPACE_START_DIRS]`:
- [ ] `TestParseStartupDirs_REQ_WORKSPACE_START_DIRS`
- [ ] `TestSeedStartupWorkspaces_REQ_WORKSPACE_START_DIRS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-07 | — | ✅ Pass | `go test ./...` (darwin/arm64, Go 1.24.3) |
| 2026-01-07 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 288 token references across 58 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:WORKSPACE_BOOTSTRAP], [REQ:WORKSPACE_START_DIRS], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
