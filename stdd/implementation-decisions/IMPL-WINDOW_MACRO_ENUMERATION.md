# [IMPL:WINDOW_MACRO_ENUMERATION] Window Macro Enumeration

**Cross-References**: [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]  
**Status**: Active  
**Created**: 2026-01-05  
**Last Updated**: 2026-01-17

---

## Decision

Introduce `%D@`/`%~D@` expansions via helper modules.

## Rationale

- Extends external-command macros so scripts can see every workspace directory without hard-coding pane order
- Keeps the macro parser maintainable by isolating the new behavior into two helpers instead of embedding list construction inside the switch
- Preserves compatibility with existing quoting (`%D`/`%~D`) and escape semantics while offering deterministic ordering for automation

## Implementation Approach

- Added `otherWindowDirPaths(ws *filer.Workspace) []string` (Module 1 `WindowSequenceBuilder`) that iterates from `Focus+1` through all directories, wrapping via modulo arithmetic. Returns an empty slice if there is only one directory. A companion `otherWindowDirNames` helper derives the same deterministic ordering but returns `Directory.Base()` so `%d@` can reuse the sequence logic without duplicating basename handling across call sites

- Added `formatDirListForMacro(paths []string, quote bool) string` (Module 2 `MacroListFormatter`) that applies `util.Quote` per entry when `quote=true`, leaves entries untouched when `quote=false`, and joins with single spaces. `%D@` invokes the quoted branch so each path is escaped, while `%~D@` deliberately uses the raw branch to honor the tilde modifier's non-quote semantics. `%d@` shares the same formatter for directory names so the quoting guarantees (and `%~` override) behave identically whether scripts need paths or basenames. Returns an empty string when no paths are provided

- Updated `expandMacro` to detect `%D@`, `%~D@`, `%d@`, and `%~d@` by looking ahead for the `macroAllOtherDirs` sentinel (`'@'`). The dispatcher calls the helpers instead of reusing the single-path logic so quoted vs. raw behavior stays localized and tied to whether the tilde modifier was present

## Code Markers

- `app/spawn.go` helper functions and `%D@`/`%d@` branches include `// [IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]`
- `README.md` macro table entry references the same tokens so documentation remains searchable

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `app/spawn.go` - helper functions
- [ ] `app/spawn_test.go` - test functions

Tests that must reference `[REQ:WINDOW_MACRO_ENUMERATION]`:
- [ ] `TestOtherWindowDirPaths_REQ_WINDOW_MACRO_ENUMERATION`
- [ ] `TestOtherWindowDirNames_REQ_WINDOW_MACRO_ENUMERATION`
- [ ] `TestMacroListFormatting_REQ_WINDOW_MACRO_ENUMERATION`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-05 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 260 token references across 55 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:WINDOW_MACRO_ENUMERATION], [REQ:WINDOW_MACRO_ENUMERATION], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
