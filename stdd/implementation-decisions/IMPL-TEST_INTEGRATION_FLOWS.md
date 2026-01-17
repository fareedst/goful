# [IMPL:TEST_INTEGRATION_FLOWS] Integration Flow Tests

**Cross-References**: [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Snapshot/integration tests for file operation flows.

## Rationale

- Validate end-to-end behavior for open/navigate/rename/delete

## Implementation Approach

### Module Identification `[REQ:MODULE_VALIDATION]`

- `app.App` orchestration (mode wiring + widget graph)
- `filer.Workspace`/`Directory` for FS mutations and navigation
- `widget.ListBox`/`Textbox` for active view state

### Validation Criteria

- Integration fixtures create temporary directories to simulate "open directory, navigate, rename/delete" flows without touching user files
- Tests assert against deterministic transcripts (e.g., active path, list contents, status messages)

### Implemented Coverage

- `TestFlowOpenDirectory_REQ_INTEGRATION_FLOWS` exercises `Workspace` + `Directory` when opening a new path
- `TestFlowNavigateRename_REQ_INTEGRATION_FLOWS` navigates into nested directories and validates rename propagation after reload
- `TestFlowDelete_REQ_INTEGRATION_FLOWS` removes files and confirms the directory state refreshes

Future enhancement: capture golden snapshots of widget buffer output once terminal fakes exist.

## Code Markers

- Tests annotated with `[IMPL:TEST_INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/integration_test.go` - `[IMPL:TEST_INTEGRATION_FLOWS]`

Tests that must reference `[REQ:INTEGRATION_FLOWS]`:
- [ ] `TestFlowOpenDirectory_REQ_INTEGRATION_FLOWS`
- [ ] `TestFlowNavigateRename_REQ_INTEGRATION_FLOWS`
- [ ] `TestFlowDelete_REQ_INTEGRATION_FLOWS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | Integration tests pass |

## Related Decisions

- Depends on: —
- See also: [ARCH:TEST_STRATEGY_INTEGRATION], [REQ:INTEGRATION_FLOWS]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
