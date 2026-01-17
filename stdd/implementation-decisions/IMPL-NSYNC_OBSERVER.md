# [IMPL:NSYNC_OBSERVER] nsync Observer Adapter

**Cross-References**: [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]  
**Status**: Active  
**Created**: 2026-01-11  
**Last Updated**: 2026-01-17

---

## Decision

Implement an observer adapter that bridges nsync progress events to goful's progress widget.

## Rationale

- nsync uses the Observer pattern for progress notifications with callbacks like `OnStart`, `OnProgress`, `OnFinish`
- goful has an existing `progress` package with `Start()`, `Update()`, `Finish()` functions that render a progress bar
- An adapter bridges these two systems without modifying either one

## Implementation Approach

- Create `type gofulObserver struct` implementing `nsync.Observer` interface in `app/nsync.go`
- `OnStart(plan)`: Call `progress.Start(float64(plan.TotalBytes))` and `progress.StartTaskCount(plan.TotalItems)`
- `OnProgress(stats)`: Call `progress.Update(float64(stats.BytesCopied - lastBytes))` with delta tracking
- `OnItemComplete(item, result)`: Call `progress.FinishTask()` per item, emit `message.Infof` for errors
- `OnFinish(result)`: Call `progress.Finish()`, emit summary message
- Observer must be thread-safe; use mutex for byte tracking since nsync calls from multiple goroutines

## Code Markers

- `app/nsync.go` includes `// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]` in the observer struct and methods

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `app/nsync.go` - observer struct and methods

Tests that must reference `[REQ:NSYNC_MULTI_TARGET]`:
- [ ] `TestGofulObserver_REQ_NSYNC_MULTI_TARGET`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-11 | — | ✅ Pass | Observer tests pass |

## Related Decisions

- Depends on: —
- See also: [ARCH:NSYNC_INTEGRATION], [REQ:NSYNC_MULTI_TARGET], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
