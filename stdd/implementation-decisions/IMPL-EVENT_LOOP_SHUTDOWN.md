# [IMPL:EVENT_LOOP_SHUTDOWN] Event Loop Shutdown Controller

**Cross-References**: [ARCH:EVENT_LOOP_SHUTDOWN] [REQ:EVENT_LOOP_SHUTDOWN]  
**Status**: Active  
**Created**: 2026-01-04  
**Last Updated**: 2026-01-17

---

## Decision

Add explicit stop control to the `app.Goful` event poller so goroutines exit immediately when the UI shuts down.

## Rationale

- Current implementations defer to process exit; embedded workflows keep the process alive, so leaked poller goroutines chew CPU and keep writing to `g.event`
- Providing a controllable shutdown path lets us validate the poller independently and clear Debt Log item D1

## Implementation Approach

### 1. Poller Abstraction

- Introduce an interface (or pure helper) `type Poller interface { Poll(stop <-chan struct{}, out chan<- tcell.Event) }` that wraps `widget.PollEvent` and listens for a stop channel
- Use `tcell.Screen` mocks in tests to simulate events and blocked reads

### 2. Shutdown Controller

- Extend `app.Goful` with a `pollStop chan struct{}` and `sync.WaitGroup` to track poller goroutines
- When `Run` exits (or when a fatal error occurs), close `pollStop`, wait for the poller to return with a timeout, then close `g.event`
- Emit `DEBUG: [IMPL:EVENT_LOOP_SHUTDOWN] stop signal sent/received` logs gated by `GOFUL_DEBUG_EVENTLOOP` (new env var) for troubleshooting

### 3. Timeout & Error Handling

- If the poller fails to stop within the timeout, log `message.Errorf` with instructions to file a bug referencing `[REQ:EVENT_LOOP_SHUTDOWN]` and continue teardown safely
- Ensure repeated shutdown attempts are idempotent (closing an already-closed channel must not panic)

### 4. Integration with Existing Flow

- Wire the poller start/stop into existing `Run` lifecycle, ensuring other modules (cmdline, filer) still receive events while the UI is running
- Update debt log entry D1 to reflect mitigation once manual validation is recorded

## Code Markers

- `app/goful.go`, any new helper file, and associated tests include `[IMPL:EVENT_LOOP_SHUTDOWN] [ARCH:EVENT_LOOP_SHUTDOWN] [REQ:EVENT_LOOP_SHUTDOWN]` comments

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `app/goful.go` - poller and shutdown controller (pollEvents, shutdownPoller, debugLog, Goful struct fields)

Tests that must reference `[REQ:EVENT_LOOP_SHUTDOWN]`:
- [x] `TestPollerShutdownTimeout_REQ_EVENT_LOOP_SHUTDOWN`
- [x] `TestShutdownPollerIdempotent_REQ_EVENT_LOOP_SHUTDOWN`
- [x] `TestPollStopChannelClosure_REQ_EVENT_LOOP_SHUTDOWN`
- [x] `TestGofulStructFields_REQ_EVENT_LOOP_SHUTDOWN`
- [x] `TestShutdownPollerSetsClosedFlag_REQ_EVENT_LOOP_SHUTDOWN`
- [x] `TestPollerStopSignalReceived_REQ_EVENT_LOOP_SHUTDOWN`
- [x] `TestConcurrentShutdown_REQ_EVENT_LOOP_SHUTDOWN`
- [x] `TestDebugLogEnvVar_REQ_EVENT_LOOP_SHUTDOWN`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-17 | — | ✅ Complete | 8 tests pass; `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1200 token references across 76 files.` |

## Related Decisions

- Depends on: —
- See also: [ARCH:EVENT_LOOP_SHUTDOWN], [REQ:EVENT_LOOP_SHUTDOWN], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
