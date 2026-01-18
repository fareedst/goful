# [IMPL:ESCAPE_TRANSLATION] Escape Key Translation

**Cross-References**: [ARCH:ESCAPE_TRANSLATION] [REQ:ESCAPE_KEY_BEHAVIOR]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Map `tcell.KeyEscape` to the canonical `"C-["` string so that modal widgets (help popup, menus, etc.) can close with the Escape key using the same case statement that handles Ctrl+[.

## Rationale

- tcell distinguishes `KeyEscape` (key code 27) from `KeyCtrlLeftSq` (key code 91), even though both historically represent the same escape sequence in terminals
- Only `KeyCtrlLeftSq` was mapped to `"C-["` in `keyToSting`, so pressing the physical Escape key caused `EventToString` to fall back to `string(ev.Rune())` which returned `"\u0000"` (null)
- Modal widgets like `help.Help` use `case "C-["` for exit, which never matched the null string
- Mapping `KeyEscape` to `"C-["` ensures consistent behavior: Escape and Ctrl+[ both close modals

## Implementation Approach

- Add `tcell.KeyEscape: "C-["` entry to `keyToSting` map in `widget/widget.go`, annotated with `[IMPL:HELP_POPUP]`
- No keymap changes required—existing `case "C-["` statements automatically work for Escape

## Code Markers

- `widget/widget.go`: `tcell.KeyEscape: "C-[", // [IMPL:HELP_POPUP] Map Escape key same as Ctrl+[ for consistent exit behavior`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `widget/widget.go` - KeyEscape entry

Tests that should reference this behavior:
- [ ] Future: `TestEventToStringEscape_REQ_ESCAPE_KEY_BEHAVIOR` (optional, follows backspace pattern)

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | `go build ./...` and `go test ./...` passed |
| 2026-01-18 | — | ✅ Pass | Runtime debug logs confirmed: `KeyEscape` now returns `"C-["`, exit branch triggers |

**Debug Evidence:**
- Before fix: `"isEscape":true,"inKeyMap":false` → `"result":"\u0000"` → no exit
- After fix: `"isEscape":true,"inKeyMap":true` → `"result":"C-["` → exit triggered

## Related Decisions

- Pattern from: [IMPL:BACKSPACE_TRANSLATION] - same dual-key normalization approach
- See also: [REQ:HELP_POPUP_STYLING], [IMPL:HELP_POPUP], [ARCH:HELP_STYLING]

---

*Created on 2026-01-18*
