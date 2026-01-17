# [IMPL:EXTERNAL_COMMAND_BINDER] External Command Binder

**Cross-References**: [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]  
**Status**: Active  
**Created**: 2026-01-02  
**Last Updated**: 2026-01-17

---

## Decision

Extract a pure helper that converts loader entries into `menu.Add` triplets with closures that call `g.Shell`, preserving cursor offsets and surfacing placeholder entries when configs are empty.

## Rationale

- Keeps `main.go` readable and makes binder behavior testable without spinning up the widget stack
- Guarantees deterministic registration order (file order) and simple hooks for future menu destinations beyond `external-command`
- Provides user-facing feedback when no commands remain (placeholder entry says "no commands configured" and logs a `DEBUG:` line)

## Implementation Approach

- New file `main_external_commands.go` in package `main` defines:
  - `type shellInvoker func(cmd string, offset ...int)` to abstract `g.Shell`
  - `func buildExternalMenuSpecs(entries []externalcmd.Entry) []menuSpec` which normalizes menu names, drops entries missing required fields (defensive), and ensures file order is preserved
  - `func registerExternalCommandMenu(g *app.Goful, entries []externalcmd.Entry)` which calls `buildExternalMenuSpecs`, adds a placeholder entry if specs are empty, and feeds `menu.Add` arguments with closures capturing the right offset
  - Placeholder callbacks call `message.Info` so pressing `X` explains that no commands are configured instead of crashing
- Tests in `main_external_commands_test.go` inject fake `shellInvoker` functions and assert commands/offsets propagate correctly and placeholder behavior triggers when expected

## Code Markers

- `main_external_commands.go` and its tests include `[IMPL:EXTERNAL_COMMAND_BINDER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]`
- `main.go` references `[IMPL:EXTERNAL_COMMAND_BINDER]` when wiring the menu after loading definitions

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `main_external_commands.go` - `[IMPL:EXTERNAL_COMMAND_BINDER]`
- [ ] `main_external_commands_test.go` - `[REQ:EXTERNAL_COMMAND_CONFIG]`
- [ ] `main.go` - menu wiring reference

Tests that must reference `[REQ:EXTERNAL_COMMAND_CONFIG]`:
- [ ] `TestBuildExternalMenuSpecs_REQ_EXTERNAL_COMMAND_CONFIG`
- [ ] `TestRegisterExternalCommandsPlaceholder_REQ_EXTERNAL_COMMAND_CONFIG`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-02 | — | ✅ Pass | `go test ./...` covers binder tests |

## Related Decisions

- Depends on: [IMPL:EXTERNAL_COMMAND_LOADER]
- See also: [ARCH:EXTERNAL_COMMAND_REGISTRY], [REQ:EXTERNAL_COMMAND_CONFIG], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
