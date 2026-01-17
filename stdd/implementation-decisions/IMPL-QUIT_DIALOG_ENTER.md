# [IMPL:QUIT_DIALOG_ENTER] Quit Dialog Return Handling

**Cross-References**: [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]  
**Status**: Active  
**Created**: 2026-01-02  
**Last Updated**: 2026-01-17

---

## Decision

Map `tcell.KeyEnter` → `C-m` and guard with regression tests.

## Rationale

- Keeps the cmdline submission shortcut stable even when upstream terminal libraries change raw key codes
- Fixes the regression where Return no longer exited the quit dialog after dependency upgrades

## Implementation Approach

- Extend `keyToSting` (sic) in `widget.EventToString` to treat `tcell.KeyEnter` identically to `tcell.KeyCtrlM`
- Add focused unit tests in `widget/widget_test.go` asserting both `KeyEnter` and `KeyCtrlM` emit the canonical `C-m` string
- Annotate the mapping and tests with `[IMPL:QUIT_DIALOG_ENTER] [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]` comments for traceability

## Code Markers

- `widget/widget.go` mapping comment at the new dictionary entry
- `widget/widget_test.go` test names/comments referencing `[REQ:QUIT_DIALOG_DEFAULT]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `widget/widget.go` - keyToSting map entry
- [ ] `widget/widget_test.go` - translator tests

Tests that must reference `[REQ:QUIT_DIALOG_DEFAULT]`:
- [ ] Tests verifying `KeyEnter` and `KeyCtrlM` emit `C-m`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-02 | — | ✅ Pass | `go test ./...` (darwin/arm64, Go 1.24.3) |
| 2026-01-02 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 25 token references across 36 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:QUIT_DIALOG_KEYS], [REQ:QUIT_DIALOG_DEFAULT]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
