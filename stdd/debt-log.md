# Debt & Risk Backlog [REQ:DEBT_TRIAGE] [ARCH:DEBT_MANAGEMENT] [IMPL:DEBT_TRACKING]

## Tracking Approach

- Every entry references the code location plus TODO tag to satisfy `[REQ:DEBT_TRIAGE]`.
- Owners use the `goful-maintainers` alias so outstanding debt is assignable for audits.
- Hotspots inline TODOs contain `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`, keeping the trace intact for `[PROC:TOKEN_AUDIT]`.

## Snapshot — 2026-01-17

| ID | Area | Risk | Owner | TODO Reference | Next Action |
|----|------|------|-------|----------------|-------------|
| D1 | `app/goful.Run` event poller | ✅ RESOLVED — Goroutine & channel leak after exit causes runaway CPU usage and blocks future runs | goful-maintainers | Implemented in `app/goful.go` | — |
| D2 | CLI history persistence (`main.go` + `cmdline/history.go`) | ✅ RESOLVED — First-run missing files return nil; actual IO failures surface via message.Errorf | goful-maintainers | Implemented in `cmdline/cmdline.go` | — |
| D3 | History cache growth (`cmdline/cmdline.go`) | ✅ RESOLVED — `HistoryLimit` bounds entries per mode (default 1000); oldest entries evicted | goful-maintainers | Implemented in `cmdline/cmdline.go` | — |
| D4 | `filer.AddExtmap` | ✅ RESOLVED — Inner map allocated if missing; safe for third-party integrations | goful-maintainers | Implemented in `filer/filer.go` | — |

## Item Details

### D1. Event Poller Stop Control — ✅ RESOLVED

- **Context**: `app/goful.go` launches `widget.PollEvent` in a tight infinite loop without observing `g.exit`.
- **Impact**: After `Run` returns the goroutine continues to push into `g.event`, leaking goroutines and hammering the channel buffer.
- **Mitigation Outline**: Move the poller behind a context or expose `stop <- struct{}{}` to break the loop while draining pending events.
- **Resolution (2026-01-17)**: Implemented `[REQ:EVENT_LOOP_SHUTDOWN]`, `[ARCH:EVENT_LOOP_SHUTDOWN]`, `[IMPL:EVENT_LOOP_SHUTDOWN]`:
  - Added `pollStop` channel and `pollWg` wait group to coordinate shutdown.
  - `pollEvents()` goroutine checks for stop signal before and after `widget.PollEvent` calls.
  - `shutdownPoller()` closes the stop channel with mutex protection (idempotent), waits for poller with timeout.
  - Debug logging gated by `GOFUL_DEBUG_EVENTLOOP` environment variable.
  - 8 unit tests validate shutdown behavior including concurrent shutdown safety.
  - Old `// TODO(goful-maintainers)` removed from `app/goful.go`.

### D2. History Error Handling — ✅ RESOLVED
- **Context**: `_ = cmdline.LoadHistory()` and `_ = cmdline.SaveHistory()` swallow IO failures; `LoadHistory` also treats `os.ErrNotExist` as fatal.
- **Impact**: Users lose shell history silently and cannot distinguish between first-run behavior and permission or disk errors.
- **Mitigation Outline**: Teach `LoadHistory` to treat missing files as success, log real errors via `message.Error`, and surface failures to exit handling.
- **Resolution (2026-01-17)**: Implemented `[IMPL:HISTORY_ERROR_HANDLING]`:
  - Added `HistoryError` struct with `Path`, `Op`, `Err` fields and `IsFirstRun()` helper.
  - `LoadHistory` returns nil for `os.ErrNotExist` (first-run), returns `*HistoryError` for real IO failures.
  - `SaveHistory` returns `*HistoryError` for permission/disk errors.
  - `main.go` surfaces load errors via `message.Errorf` and save errors via `stderr` (TUI finalized at exit).
  - 6 unit tests validate error differentiation (`history_test.go`).
  - Old `// TODO(goful-maintainers)` removed from `main.go` and `cmdline/cmdline.go`.

### D3. History Cache Boundaries — ✅ RESOLVED
- **Context**: `historyMap` stores every past command forever; dedupe removes duplicates but never trims size.
- **Impact**: Long-running sessions accumulate unbounded history, causing memory growth and larger save files that slow down startup/shutdown.
- **Mitigation Outline**: Track a configurable per-mode limit (e.g., 1k entries) and drop the oldest entries before serialization.
- **Resolution (2026-01-17)**: Implemented `[IMPL:HISTORY_CACHE_LIMIT]`:
  - Added `HistoryLimit` variable (default 1000 entries per mode; set to 0 for unlimited).
  - `trimHistory()` helper drops oldest entries when limit exceeded, keeping most recent.
  - Trimming applied during `History.add()` (runtime) and `SaveHistory` (persistence compaction).
  - 3 unit tests validate eviction behavior (`history_test.go`).
  - Old `// TODO(goful-maintainers)` removed from `cmdline/cmdline.go`.

### D4. Extmap API Safety — ✅ RESOLVED
- **Context**: `filer.AddExtmap` writes to `f.extmap[key][ext]` without initializing the inner map, panicking when used by plugins.
- **Impact**: Third-party integrations calling `AddExtmap` crash immediately, so the helper is effectively unusable outside core.
- **Mitigation Outline**: Allocate `f.extmap[key]` when missing and add regression coverage to keep the API reliable.
- **Resolution (2026-01-17)**: Implemented `[IMPL:EXTMAP_API_SAFETY]`:
  - Check if inner map exists before writing; allocate if nil.
  - 2 regression tests validate nil-map safety and multiple-entry scenarios (`integration_test.go`).
  - Old `// TODO(goful-maintainers)` removed from `filer/filer.go`.


