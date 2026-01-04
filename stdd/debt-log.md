# Debt & Risk Backlog [REQ:DEBT_TRIAGE] [ARCH:DEBT_MANAGEMENT] [IMPL:DEBT_TRACKING]

## Tracking Approach

- Every entry references the code location plus TODO tag to satisfy `[REQ:DEBT_TRIAGE]`.
- Owners use the `goful-maintainers` alias so outstanding debt is assignable for audits.
- Hotspots inline TODOs contain `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`, keeping the trace intact for `[PROC:TOKEN_AUDIT]`.

## Snapshot — 2026-01-01

| ID | Area | Risk | Owner | TODO Reference | Next Action |
|----|------|------|-------|----------------|-------------|
| D1 | `app/goful.Run` event poller | Goroutine & channel leak after exit causes runaway CPU usage and blocks future runs | goful-maintainers | `// TODO(goful-maintainers)` in `app/goful.go` | Introduce context-aware poller and close `g.event` when shutting down |
| D2 | CLI history persistence (`main.go` + `cmdline/history.go`) | Ignored errors hide corrupt history; missing files are treated as fatal and never logged | goful-maintainers | `// TODO(goful-maintainers)` in `main.go` & `cmdline/cmdline.go` | Differentiate `os.IsNotExist`, surface actionable errors via `message.Error`, add unit coverage |
| D3 | History cache growth (`cmdline/cmdline.go`) | `historyMap` never bounds entries leading to unbounded memory per mode | goful-maintainers | `// TODO(goful-maintainers)` near `historyMap` | Add eviction policy (N most recent) and persistence compaction |
| D4 | `filer.AddExtmap` | Nil map panic when invoked before `MergeExtmap` seeds inner map | goful-maintainers | `// TODO(goful-maintainers)` in `filer/filer.go` | Allocate inner map or deprecate API; add regression test |

## Item Details

### D1. Event Poller Stop Control
- **Context**: `app/goful.go` launches `widget.PollEvent` in a tight infinite loop without observing `g.exit`.
- **Impact**: After `Run` returns the goroutine continues to push into `g.event`, leaking goroutines and hammering the channel buffer.
- **Mitigation Outline**: Move the poller behind a context or expose `stop <- struct{}{}` to break the loop while draining pending events.

### D2. History Error Handling
- **Context**: `_ = cmdline.LoadHistory()` and `_ = cmdline.SaveHistory()` swallow IO failures; `LoadHistory` also treats `os.ErrNotExist` as fatal.
- **Impact**: Users lose shell history silently and cannot distinguish between first-run behavior and permission or disk errors.
- **Mitigation Outline**: Teach `LoadHistory` to treat missing files as success, log real errors via `message.Error`, and surface failures to exit handling.

### D3. History Cache Boundaries
- **Context**: `historyMap` stores every past command forever; dedupe removes duplicates but never trims size.
- **Impact**: Long-running sessions accumulate unbounded history, causing memory growth and larger save files that slow down startup/shutdown.
- **Mitigation Outline**: Track a configurable per-mode limit (e.g., 1k entries) and drop the oldest entries before serialization.

### D4. Extmap API Safety
- **Context**: `filer.AddExtmap` writes to `f.extmap[key][ext]` without initializing the inner map, panicking when used by plugins.
- **Impact**: Third-party integrations calling `AddExtmap` crash immediately, so the helper is effectively unusable outside core.
- **Mitigation Outline**: Allocate `f.extmap[key]` when missing and add regression coverage to keep the API reliable.


