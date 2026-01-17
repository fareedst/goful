# [IMPL:TEST_CMDLINE] Command Tests

**Cross-References**: [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Add tests for command parsing and modes.

## Rationale

- Prevent regressions in command handling and state transitions

## Implementation Approach

### Module Identification `[REQ:MODULE_VALIDATION]`

- `cmdline.Parser` (tokenization + quoting)
- `cmdline.Completion` helpers (word boundary + suggestion generation)
- `app.Mode` transitions (normal, command, prompt)

### Validation Criteria

- Parser tests feed representative command strings (including quoting, globbing, multi-byte) and assert resulting structs
- Mode tests stimulate key-event handlers with table-driven inputs to ensure state-dependent callbacks fire

### Immediate Test Plan

- `TestParseLine_REQ_CMD_HANDLER_TESTS`: ensures parser emits expected argv slices plus error handling for unterminated quotes
- `TestApplyModeTransitions_REQ_CMD_HANDLER_TESTS`: uses lightweight fakes to confirm `mode.GoNormal()` / `mode.GoCommand()` toggles behavior
- `TestCompletionFilters_REQ_CMD_HANDLER_TESTS`: validates completion filter respects prefixes + case sensitivity
- Additional coverage will mock `cmdline.Extmap` to exercise edge bindings before integration tests

## Code Markers

- Tests carry `[IMPL:TEST_CMDLINE] [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `cmdline/*_test.go` - `[IMPL:TEST_CMDLINE]`

Tests that must reference `[REQ:CMD_HANDLER_TESTS]`:
- [ ] `TestParseLine_REQ_CMD_HANDLER_TESTS`
- [ ] `TestApplyModeTransitions_REQ_CMD_HANDLER_TESTS`
- [ ] `TestCompletionFilters_REQ_CMD_HANDLER_TESTS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | Command tests pass |

## Related Decisions

- Depends on: —
- See also: [ARCH:TEST_STRATEGY_CMD], [REQ:CMD_HANDLER_TESTS]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
