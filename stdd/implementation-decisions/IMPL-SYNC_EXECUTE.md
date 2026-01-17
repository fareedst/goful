# [IMPL:SYNC_EXECUTE] Sync Execute

**Cross-References**: [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]  
**Status**: Active  
**Created**: 2026-01-14  
**Last Updated**: 2026-01-17

---

## Decision

Implement sync command mode with prefix key pattern and sequential execution.

## Rationale

- Implements the requirement for synchronized batch operations across workspace panes per [REQ:SYNC_COMMANDS]
- Uses the established cmdline mode pattern for prompts, matching existing copy/delete/rename UX
- Separates the mode state, execution engine, and operation modes for independent validation per [REQ:MODULE_VALIDATION]

## Implementation Approach

### Add `FindFileByName` to `filer/directory.go`

- Method on `*Directory` that scans the file list for an exact name match
- Returns `*FileStat` if found, `nil` otherwise
- Used by execution engine to locate same-named files in each pane

### Create `app/window_wide.go`

- `syncMode` struct implementing `cmdline.Mode` that captures `ignoreFailures bool`
- `SyncMode(ignoreFailures bool)` entry point on `*Goful` that activates the mode
- Prompt displays "Sync: [c]opy [d]elete [r]ename" (or with "!" suffix for ignore mode)
- Input handler maps `c`/`d`/`r` to operation-specific modes, `!` toggles ignore mode, `C-g`/`C-[` exits
- `SyncResult` struct to aggregate succeeded/skipped/failed counts with error details
- `executeSyncCopy`, `executeSyncDelete`, `executeSyncRename` functions that iterate panes

### Add operation modes to `app/mode.go`

- `syncCopyMode`: Prompts for new filename, validates user enters a different name, stores `ignoreFailures`, copies file to new name in each pane's directory
- `syncDeleteMode`: Confirms deletion with y/n, stores `ignoreFailures`, executes on confirm
- `syncRenameMode`: Prompts for new name, stores `ignoreFailures`, executes on confirm
- All modes follow the existing pattern (implement `cmdline.Mode` interface)

### Wire keybindings in `main.go`

- `"S"`: `func() { g.SyncMode(false) }` – normal mode (abort on failure)

### Execution Order

1. Start with focused pane (`ws.Focus`)
2. Proceed through remaining panes in order (`(ws.Focus + i) % len(ws.Dirs)`)
3. For each pane, find file by exact name match
4. If not found, increment skip count and continue
5. If found, execute operation:
   - Copy: copy source file to new filename in the same directory
   - Rename: rename file to new name
   - Delete: remove file
6. On error, record failure and either abort (default) or continue (ignore mode, toggled with `!`)
7. After completion, reload all directories and report summary

## Code Markers

- `filer/directory.go`: `FindFileByName` with `// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]`
- `app/window_wide.go`: All types and functions with same tokens
- `app/mode.go`: Operation mode structs with same tokens
- `main.go`: Keybindings with `// [IMPL:SYNC_EXECUTE] [REQ:SYNC_COMMANDS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/directory.go` - FindFileByName
- [ ] `app/window_wide.go` - all functions
- [ ] `app/mode.go` - sync operation modes
- [ ] `main.go` - keybindings

Tests that must reference `[REQ:SYNC_COMMANDS]`:
- [ ] `TestFindFileByName_REQ_SYNC_COMMANDS`
- [ ] `TestSyncResult_REQ_SYNC_COMMANDS`
- [ ] `TestCopyFileSimple_REQ_SYNC_COMMANDS`
- [ ] `TestCopyDirRecursive_REQ_SYNC_COMMANDS`
- [ ] `TestSyncMode_REQ_SYNC_COMMANDS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-14 | — | ✅ Pass | `go test ./...` passes on darwin/arm64, Go 1.24.3 |

## Related Decisions

- Depends on: —
- See also: [ARCH:SYNC_MODE], [REQ:SYNC_COMMANDS], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
