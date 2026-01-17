# [IMPL:TERMINAL_ADAPTER] Terminal Adapter Module

**Cross-References**: [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD]  
**Status**: Active  
**Created**: 2026-01-04  
**Last Updated**: 2026-01-17

---

## Decision

Implement a `terminalcmd` package that encapsulates platform-aware terminal command creation plus the glue that registers it with `g.ConfigTerminal`.

## Rationale

- Keeps OS detection, tmux handling, and macOS automation scripts isolated from UI wiring so we can satisfy [REQ:TERMINAL_PORTABILITY] without scattering conditionals
- Enables fast module validation—`CommandFactory` can be unit-tested with pure inputs, and `Configurator` can be exercised via fakes that capture the configured command slices

## Implementation Approach

### Module 1: `CommandFactory`

- Signature: `func NewFactory(opts Options) Factory`
- `Options` include `GOOS`, `IsTmux`, `Override []string`, `Tail string`, plus macOS-specific fields `TerminalApp string` (default `Terminal`) and `TerminalShell string` (default `bash`) so AppleScript launches can be customized without editing Go code
- `Factory.CommandWithCwd(cmd string, cwd string) []string` returns:
  - Override path: `Override + []string{"bash", "-c", payload}` where `payload` already prefixes macOS commands with `cd "<cwd>";` to satisfy `[REQ:TERMINAL_CWD]`
  - Tmux path: `[]string{"tmux", "new-window", "-n", title(cmd), cmd + tail}`
  - macOS path: `[]string{"osascript", "-e", fmt.Sprintf("tell application \"%s\" to do script \"%s\"", terminalApp, script), "-e", fmt.Sprintf("tell application \"%s\" to activate", terminalApp)}` where `script` embeds the configured shell (`<terminalShell> -c "cd \"<cwd>\"; <cmd + tail>"; exit`)
  - Linux default: maintain current gnome-terminal invocation with title-setting escape
- Emits `DEBUG: [IMPL:TERMINAL_ADAPTER] ...` logs describing the branch taken and any overrides, guarded by `GOFUL_DEBUG_TERMINAL=1`

### Module 2: `Configurator`

- Accepts a `Factory` and returns the closure passed to `g.ConfigTerminal`
- Responsible for injecting the "HIT ENTER KEY" tail, escaping titles, ensuring `bash -c` semantics remain unchanged, and providing a live `cwd` callback so macOS launches always reflect the focused directory
- Surface helper `Apply(cfg Configurator, factory Factory, cwd func() string)` that wires both shell and terminal commands where appropriate

### Environment & Overrides

- Parse `GOFUL_TERMINAL_CMD` (string) or `-terminal` flag (future) into the override slice
- Read `GOFUL_TERMINAL_APP` and `GOFUL_TERMINAL_SHELL` (with defaults baked into `NewFactory`) so AppleScript launches can target another application or shell without modifying Go code
- Document how to supply fallback commands (e.g., `iTerm2`)

### macOS Shell Invocation Safeguard

- Switching the AppleScript payload from `bash -lc` to `<shell> -c` prevents login-shell initialization from hanging Terminal windows that source interactive profiles while respecting the configured shell binary
- The payload is still quoted via `strconv.Quote` so `<shell> -c` receives the entire `cd "<cwd>"; <cmd><tail>` sequence intact while avoiding `.bash_profile` prompts

## Code Markers

- `terminalcmd/*.go` and `main.go` wiring include `[IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `terminalcmd/factory.go` - CommandFactory
- [ ] `terminalcmd/factory_test.go` - tests
- [ ] `main.go` - wiring

Tests that must reference `[REQ:TERMINAL_PORTABILITY]`:
- [ ] `TestCommandFactoryDarwin_REQ_TERMINAL_PORTABILITY`
- [ ] Other platform tests

Manual validation checklist is documented as `[PROC:TERMINAL_VALIDATION]` in `stdd/processes.md`.

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-04 | — | ✅ Pass | `go test ./terminalcmd` (darwin/arm64, Go 1.24.3) |
| 2026-01-04 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 245 token references across 52 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:TERMINAL_LAUNCHER], [REQ:TERMINAL_PORTABILITY], [REQ:TERMINAL_CWD], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
