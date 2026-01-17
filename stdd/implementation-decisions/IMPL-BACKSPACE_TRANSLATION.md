# [IMPL:BACKSPACE_TRANSLATION] Backspace Key Translation

**Cross-References**: [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]  
**Status**: Active  
**Created**: 2026-01-09  
**Last Updated**: 2026-01-17

---

## Decision

Normalize both tcell backspace key codes so they emit the canonical `backspace` string consumed by filer and prompt keymaps.

## Rationale

- macOS terminals often report Backspace as `tcell.KeyBackspace` while Linux/tmux sessions use `tcell.KeyBackspace2`. Only one entry existed in the translator map, so Backspace failed silently on half the platforms
- Mapping both key codes to the same string preserves historical behavior (Backspace opens parent directory, deletes the previous rune) without requiring duplicate keymap entries or user-specific configuration
- Keeping the normalization inside `widget.EventToString` satisfies `[REQ:MODULE_VALIDATION]` by addressing the issue within the pure translator module that already underpins cmdline/filer behavior

## Implementation Approach

- Extend `keyToSting` in `widget/widget.go` with a `tcell.KeyBackspace` entry pointing to `"backspace"` and annotate both entries with `[IMPL:BACKSPACE_TRANSLATION] [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]`
- Add table-driven unit test `TestEventToStringBackspace_REQ_BACKSPACE_BEHAVIOR` in `widget/widget_test.go` that creates events for both `tcell.KeyBackspace` and `tcell.KeyBackspace2` and asserts `EventToString` returns `backspace`
- Retain existing `main_keymap_test.go` baseline coverage so the `backspace` binding remains required in filer/cmdline/finder/completion keymaps

## Code Markers

- `widget/widget.go` map entries for both backspace keys include `[IMPL:BACKSPACE_TRANSLATION] [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]`
- `widget/widget_test.go` test includes the same triplet and references `[REQ:BACKSPACE_BEHAVIOR]` in the function name

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `widget/widget.go` - both backspace key entries
- [ ] `widget/widget_test.go` - backspace test
- [ ] `main_keymap_test.go` - baseline coverage

Tests that must reference `[REQ:BACKSPACE_BEHAVIOR]`:
- [ ] `TestEventToStringBackspace_REQ_BACKSPACE_BEHAVIOR`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-09 | — | ✅ Pass | `go test ./...` (darwin/arm64, Go 1.24.3) |
| 2026-01-09 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 520 token references across 66 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:BACKSPACE_TRANSLATION], [REQ:BACKSPACE_BEHAVIOR], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
