# [IMPL:BASELINE_SNAPSHOTS] Baseline Snapshots

**Cross-References**: [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Capture current keybindings/modes as automated baselines.

## Rationale

- Preserve behavior ahead of refactors
- Provide guardrails for future keymap cleanups or menu consolidation work

## Implementation Approach

- Implement `KeymapBaselineSuite` unit tests under `main_keymap_test.go` that:
  - Instantiate maps via `filerKeymap(nil)`, `finderKeymap(nil)`, `cmdlineKeymap(new(cmdline.Cmdline))`, `completionKeymap(new(cmdline.Completion))`, `menuKeymap(new(menu.Menu))`
  - Assert presence of representative key chords for navigation, selection, shell execution, finder/completion movement, and exit behaviors
  - Emit `DEBUG:` logs enumerating the verified chords for traceability
- Introduce helper `assertKeyCoverage` to keep tests declarative and make future updates additive
- Tag tests with `[TEST:KEYMAP_BASELINE]` alongside `[REQ:BEHAVIOR_BASELINE]` tokens
- Keep suite pure (no widget initialization) so it runs instantly in CI

## Code Markers

- Tests include `[IMPL:BASELINE_SNAPSHOTS] [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE] [TEST:KEYMAP_BASELINE]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `main_keymap_test.go` - all baseline tests

Tests that must reference `[REQ:BEHAVIOR_BASELINE]`:
- [ ] `KeymapBaselineSuite` tests

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | Baseline suite passes in CI |

- `go test ./...` (module validation) covers the new baseline suite prior to integrating any runtime changes
- `./scripts/validate_tokens.sh` ensures `[TEST:KEYMAP_BASELINE]` and related tokens are registered

## Related Decisions

- Depends on: —
- See also: [ARCH:BASELINE_CAPTURE], [REQ:BEHAVIOR_BASELINE], [TEST:KEYMAP_BASELINE], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
