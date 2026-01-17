# [IMPL:NSYNC_COPY_MOVE] nsync Copy/Move Wrappers and CopyAll/MoveAll Functions

**Cross-References**: [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]  
**Status**: Active  
**Created**: 2026-01-11  
**Last Updated**: 2026-01-17

---

## Decision

Provide wrapper functions that configure and execute nsync sync operations within goful's async file control pattern, and add new cmdline modes that collect all visible workspace directories as destinations.

## Rationale

- nsync's `Syncer.Sync()` is synchronous and blocks until complete; goful needs to run it in a background goroutine
- The `asyncFilectrl` pattern handles UI resizing, progress widget space, and workspace reload after completion
- Wrapper functions encapsulate nsync configuration for clean call sites
- Users expect symmetry with existing `Copy`/`Move` commands
- Using `otherWindowDirPaths()` provides consistent destination enumeration

## Implementation Approach

### syncCopy and syncMove Wrappers

`func (g *Goful) syncCopy(sources []string, destinations []string)`:
- Resolves absolute paths for all sources
- Configures `nsync.Config{Sources, Destinations, Recursive: true, Move: false, Jobs: 4}`
- Creates syncer with `nsync.WithObserver(gofulObserver)`
- Calls within `asyncFilectrl` goroutine pattern
- Reports result via `message.Infof`/`message.Error`

`func (g *Goful) syncMove(sources []string, destinations []string)`:
- Same as `syncCopy` but with `Move: true` in config
- nsync handles source deletion after successful sync to all destinations

Context cancellation: Create `context.WithCancel` that listens for user interrupt (future enhancement)

### CopyAll and MoveAll Functions

`func (g *Goful) CopyAll()`:
- If only one directory visible: delegate to `g.Copy()` with message explaining fallback
- Collect sources: if marks exist, use `g.Dir().MarkfilePaths()`; else use cursor file `g.File().Path()`
- Collect destinations: `otherWindowDirPaths(g.Workspace())`
- Call `g.syncCopy(sources, destinations)`

`func (g *Goful) MoveAll()`:
- Same pattern as `CopyAll` but calls `g.syncMove`

No cmdline text input needed—operation is immediate after command invocation.

### Keybindings

- `C` for CopyAll, `M` for MoveAll
- `` ` `` (backtick) is used for toggle comparison colors

## Code Markers

- `app/nsync.go` includes `// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]` in wrapper functions and CopyAll/MoveAll
- `main.go` keybindings include `// [IMPL:NSYNC_COPY_MOVE] [REQ:NSYNC_MULTI_TARGET]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `app/nsync.go` - syncCopy, syncMove, CopyAll, MoveAll
- [ ] `main.go` - keybindings

Tests that must reference `[REQ:NSYNC_MULTI_TARGET]`:
- [ ] `TestSyncCopy_REQ_NSYNC_MULTI_TARGET`
- [ ] `TestSyncMove_REQ_NSYNC_MULTI_TARGET`
- [ ] Integration tests verifying destination enumeration and fallback behavior

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-11 | — | ✅ Pass | Wrapper tests pass |

## Related Decisions

- Depends on: [IMPL:NSYNC_OBSERVER]
- See also: [ARCH:NSYNC_INTEGRATION], [REQ:NSYNC_MULTI_TARGET], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
