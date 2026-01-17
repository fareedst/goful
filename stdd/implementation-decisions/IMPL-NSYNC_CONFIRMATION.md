# [IMPL:NSYNC_CONFIRMATION] nsync Confirmation Modes

**Cross-References**: [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]  
**Status**: Active  
**Created**: 2026-01-11  
**Last Updated**: 2026-01-17

---

## Decision

Implement confirmation prompts for CopyAll/MoveAll using cmdline modes similar to quitMode/removeMode.

## Rationale

- Multi-target operations are high-risk and users expect confirmation before files are synced to multiple destinations
- Reusing the existing cmdline mode pattern keeps implementation simple and UX consistent with other confirmation dialogs
- The confirmation displays source count and destination count so users understand the scope of the operation

## Implementation Approach

### copyAllMode (`app/mode.go`)

- Fields: `*Goful`, `sources []string`, `destinations []string`
- `String()`: returns `"copyall"`
- `Prompt()`: returns `fmt.Sprintf("Copy %d file(s) to %d destinations? [Y/n] ", len(sources), len(destinations))`
- `Draw()`: calls `c.DrawLine()`
- `Run()`: on `Y`/`y`/empty calls `m.doCopyAll(sources, destinations)` and `c.Exit()`; on `n`/`N` calls `c.Exit()`; else clears text

### moveAllMode (`app/mode.go`)

- Same pattern as `copyAllMode` but with "Move" label and calls `m.doMoveAll()`

### Refactored `app/nsync.go`

- Rename current `syncCopy`/`syncMove` internals to `doCopyAll()`/`doMoveAll()` (private execution methods)
- New public `CopyAll()` collects sources/destinations, then starts `copyAllMode` if valid
- New public `MoveAll()` collects sources/destinations, then starts `moveAllMode` if valid
- Single-pane fallback logic (`g.Copy()`/`g.Move()`) remains in the public methods before mode creation

## Code Markers

- `app/mode.go` confirmation mode structs include `// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]`
- `app/nsync.go` refactored public/private methods include `// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `app/mode.go` - copyAllMode, moveAllMode
- [ ] `app/nsync.go` - public methods and private execution methods

Tests that must reference `[REQ:NSYNC_CONFIRMATION]`:
- [ ] `TestCopyAllConfirmation_REQ_NSYNC_CONFIRMATION`
- [ ] `TestMoveAllConfirmation_REQ_NSYNC_CONFIRMATION`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| — | — | ⏳ Pending | To be captured after implementation |

## Related Decisions

- Depends on: [IMPL:NSYNC_COPY_MOVE]
- See also: [ARCH:NSYNC_CONFIRMATION], [REQ:NSYNC_CONFIRMATION], [REQ:NSYNC_MULTI_TARGET], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
