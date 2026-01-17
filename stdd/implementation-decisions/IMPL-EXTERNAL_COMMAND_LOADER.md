# [IMPL:EXTERNAL_COMMAND_LOADER] External Command Loader

**Cross-References**: [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]  
**Status**: Active  
**Created**: 2026-01-02  
**Last Updated**: 2026-01-17

---

## Decision

Provide a dedicated loader that resolves the config file path (flag/env/default), parses JSON or YAML, validates schema, filters by platform, and falls back to embedded defaults.

## Rationale

- Keeps customization outside of `main.go`, enabling dotfile repos or tooling scripts to ship command sets without recompiling
- Ensures `[REQ:MODULE_VALIDATION]` can be satisfied with pure unit tests (no widget/app dependencies)
- Preserves historic behavior (platform-specific defaults, cursor offsets) whenever the config file is absent or invalid

## Implementation Approach

- Extend `configpaths.Paths` with `Commands` + `CommandsSource` so `emitPathDebug` reports all three precedence outcomes. CLI flag `-commands` and env var `GOFUL_COMMANDS_FILE` feed the resolver
- New package `externalcmd` exposes:
  - `type Entry` with `Menu`, `Key`, `Label`, `Command`, `Offset`, `Platforms`, `Disabled`
  - `func Defaults(goos string) []Entry` returning the old hard-coded Windows/POSIX menus expressed with `%` macros instead of inline `g.File()` references (e.g., rename defaults to `mv -vi %f %~f`)
  - `func Load(Options) ([]Entry, error)` where `Options` carries `Path`, `GOOS`, `ReadFile`, `Debug`. Loader expands `~`, reads JSON or YAML (supporting raw arrays or `{ commands: [] }`), validates unique `menu/key`, enforces required fields, filters by `Platforms`, skips disabled entries, logs diagnostics tagged with `[IMPL:EXTERNAL_COMMAND_LOADER]`, and **prepends file entries ahead of `Defaults` unless the file sets an explicit opt-out (e.g., `inheritDefaults: false`)**
- Errors reading/parsing the config file bubble up so callers can emit `message.Errorf` but still fall back to defaults
- Dependency note: Introduced `gopkg.in/yaml.v3` to parse YAML files without writing a bespoke parser; the package is already widely used and compatible with Go 1.24

## Code Markers

- `externalcmd/defaults.go` & `externalcmd/loader.go` include `[IMPL:EXTERNAL_COMMAND_LOADER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]`
- `configpaths/resolver.go` references `[IMPL:STATE_PATH_RESOLVER]` while documenting the new `Commands` path fields

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `externalcmd/defaults.go` - `[IMPL:EXTERNAL_COMMAND_LOADER]`
- [ ] `externalcmd/loader.go` - `[IMPL:EXTERNAL_COMMAND_LOADER]`
- [ ] `externalcmd/loader_test.go` - `[REQ:EXTERNAL_COMMAND_CONFIG]`

Tests that must reference `[REQ:EXTERNAL_COMMAND_CONFIG]`:
- [ ] `TestLoadCommands_REQ_EXTERNAL_COMMAND_CONFIG`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-02 | — | ✅ Pass | `go test ./externalcmd` covers loader in isolation |
| 2026-01-02 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 245 token references across 52 files |

## Related Decisions

- Depends on: [IMPL:STATE_PATH_RESOLVER]
- See also: [ARCH:EXTERNAL_COMMAND_REGISTRY], [REQ:EXTERNAL_COMMAND_CONFIG], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
