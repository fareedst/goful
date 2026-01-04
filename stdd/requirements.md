# Requirements Directory

**STDD Methodology Version**: 1.1.0

## Overview
This document serves as the **central listing/registry** for all requirements in this project. Each requirement has a unique semantic token `[REQ:IDENTIFIER]` for traceability.

**For detailed information about how requirements are fulfilled, see:**
- **Architecture decisions**: See `architecture-decisions.md` for high-level design choices that fulfill requirements
- **Implementation decisions**: See `implementation-decisions.md` for detailed implementation approaches, APIs, and algorithms
- **Semantic tokens**: See `semantic-tokens.md` for the complete token registry

### Requirement Structure

Each requirement includes:
- **Description**: What the requirement specifies (WHAT)
- **Rationale**: Why the requirement exists (WHY)
- **Satisfaction Criteria**: How we know the requirement is satisfied (acceptance criteria, success conditions)
- **Validation Criteria**: How we verify/validate the requirement is met (testing approach, verification methods, success metrics)

**Note**: 
- Satisfaction and validation criteria that involve architectural or implementation details reference the appropriate layers
- Architecture decisions in `architecture-decisions.md` explain HOW requirements are fulfilled at a high level
- Implementation decisions in `implementation-decisions.md` explain HOW requirements are fulfilled at a detailed level

## Requirements Registry

### Functional Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|

### Non-Functional Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|

### Immutable Requirements (Major Version Change Required)

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|

### Incomplete Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|

---

## Detailed Requirements

### Core Functionality

### [REQ:STDD_SETUP] STDD Methodology Setup

**Priority: P0 (Critical)**

- **Description**: The project must follow the Semantic Token-Driven Development (STDD) methodology, including a specific directory structure (`stdd/`) and documentation files (`requirements.md`, `architecture-decisions.md`, etc.).
- **Rationale**: To ensure traceability of intent from requirements to code and to maintain a consistent development process.
- **Satisfaction Criteria**:
  - `stdd/` directory exists.
  - All required documentation files exist and are populated from templates.
  - `.cursorrules` contains the STDD rules.
- **Validation Criteria**:
  - Manual verification of file existence.
  - AI agent acknowledgment of principles.
- **Architecture**: See `architecture-decisions.md` § STDD Project Structure [ARCH:STDD_STRUCTURE]
- **Implementation**: See `implementation-decisions.md` § STDD File Creation [IMPL:STDD_FILES]

**Status**: ✅ Implemented

### [REQ:MODULE_VALIDATION] Independent Module Validation Before Integration

**Priority: P0 (Critical)**

- **Description**: Logical modules must be developed and validated independently before integration into code satisfying specific requirements. Each module must have clear boundaries, interfaces, and validation criteria defined before development begins.
- **Rationale**: To eliminate bugs related to code complexity by ensuring each module works correctly in isolation before combining with other modules. Independent validation reduces integration complexity, catches bugs early in the development cycle, and ensures each module meets its defined contract before integration.
- **Satisfaction Criteria**:
  - All logical modules are identified and documented with clear boundaries before development.
  - Module interfaces and contracts are defined and documented.
  - Module validation criteria are specified (what "validated" means for each module).
  - Each module is developed independently.
  - Each module passes independent validation (unit tests with mocks, integration tests with test doubles, contract validation, edge case testing, error handling validation) before integration.
  - Module validation results are documented.
  - Integration tasks are separate from module development and validation tasks.
  - Integration only occurs after module validation passes.
- **Validation Criteria**:
  - Manual verification that modules are identified and documented before development.
  - Automated verification that module validation tests exist and pass before integration.
  - Code review verification that integration code references validated modules.
  - Verification that module validation results are documented.
  - Verification that integration tests validate the combined behavior of validated modules.
- **Architecture**: See `architecture-decisions.md` § Module Validation Strategy [ARCH:MODULE_VALIDATION]
- **Implementation**: See `implementation-decisions.md` § Module Validation Implementation [IMPL:MODULE_VALIDATION]

**Status**: ✅ Implemented

### [REQ:QUIT_DIALOG_DEFAULT] Quit Dialog Default Confirmation

**Priority: P0 (Critical)**

- **Description**: The quit confirmation dialog must accept the Return/Enter key with no additional input to select the default affirmative option so users can exit quickly without typing.
- **Rationale**: Terminal users expect Return to submit dialogs; regressions break muscle memory and can trap the user in the application.
- **Satisfaction Criteria**:
  - Pressing Return while the quit dialog is focused and empty exits the application immediately.
  - The key translation layer produces the same command invocation for Return as for `Ctrl-M` so existing keymaps remain valid.
  - Behavior is consistent across platforms supported by tcell.
- **Validation Criteria**:
  - Automated tests cover the key translation path to ensure Return continues to map to the command execution trigger.
  - Manual verification confirms the quit dialog exits when Return is pressed without entering text.
- **Architecture**: See `architecture-decisions.md` § Quit Dialog Key Translation [ARCH:QUIT_DIALOG_KEYS]
- **Implementation**: See `implementation-decisions.md` § Quit Dialog Return Handling [IMPL:QUIT_DIALOG_ENTER]

**Status**: ✅ Implemented

### [REQ:GO_TOOLCHAIN_LTS] Modern Go Toolchain Baseline

**Priority: P0 (Critical)**

- **Description**: The project must target a current Go LTS release in `go.mod` and CI to ensure modern language features, security fixes, and ecosystem compatibility.
- **Rationale**: Outdated Go versions block security patches and ecosystem tooling; aligning to LTS keeps builds reproducible and supported.
- **Satisfaction Criteria**:
  - `go.mod` declares the agreed LTS Go version; local builds use it.
  - CI matrix pins the same Go version(s) and caches modules.
  - `go fmt` / `gofmt -w` and `go vet` succeed with the updated toolchain.
  - `go mod tidy` produces a clean module graph.
- **Validation Criteria**:
  - CI run on Go LTS passes fmt/vet/test.
  - Manual `go version` in CI logs matches `go.mod`.
  - Token audit shows `[IMPL:GO_MOD_UPDATE]` annotations in go.mod-related changes.
- **Architecture**: See `architecture-decisions.md` § Go Runtime Strategy [ARCH:GO_RUNTIME_STRATEGY]
- **Implementation**: See `implementation-decisions.md` § Go Mod Update [IMPL:GO_MOD_UPDATE]

**Status**: ✅ Implemented

### [REQ:DEPENDENCY_REFRESH] Secure Dependency Updates

**Priority: P0 (Critical)**

- **Description**: Refresh dependencies (e.g., `tcell`, `golang.org/x/*`) to current compatible releases to pick up security and bug fixes.
- **Rationale**: Old transitive versions carry CVEs and incompatibilities; refreshing keeps the runtime safe and supported.
- **Satisfaction Criteria**:
  - `go.mod`/`go.sum` pinned to current stable releases for direct deps.
  - `go mod tidy` leaves no unused or missing entries.
  - Document any breaking changes or shims needed.
- **Validation Criteria**:
  - CI succeeds with updated deps on all targets.
  - Static analysis and tests pass without new regressions.
  - Token audit shows `[IMPL:DEP_BUMP]` references in dependency changes.
- **Architecture**: See `architecture-decisions.md` § Dependency Policy [ARCH:DEPENDENCY_POLICY]
- **Implementation**: See `implementation-decisions.md` § Dependency Bump [IMPL:DEP_BUMP]

**Status**: ✅ Implemented

### [REQ:CI_PIPELINE_CORE] CI Coverage for fmt/vet/tests

**Priority: P0 (Critical)**

- **Description**: Establish GitHub Actions CI that runs formatting, vetting, and unit tests on every push/PR.
- **Rationale**: Prevents regressions and enforces consistency before merging.
- **Satisfaction Criteria**:
  - Workflow triggers on PR and main branch pushes.
  - Steps: `go fmt`/`gofmt -w`, `go vet`, `go test ./...`.
  - Caches Go modules for performance.
- **Validation Criteria**:
  - CI badge shows passing runs.
  - Workflow logs include `DEBUG:`/`DIAGNOSTIC:` markers where applicable.
  - Token audit: `[IMPL:CI_WORKFLOW]` in workflow file comments.
- **Architecture**: See `architecture-decisions.md` § CI Pipeline [ARCH:CI_PIPELINE]
- **Implementation**: See `implementation-decisions.md` § CI Workflow [IMPL:CI_WORKFLOW]

**Status**: ✅ Implemented

### [REQ:STATIC_ANALYSIS] Static Analysis Gate

**Priority: P1 (Important)**

- **Description**: Add staticcheck and (optionally) golangci-lint to CI to catch correctness and style issues early.
- **Rationale**: Static analysis surfaces defects and API misuse not caught by tests.
- **Satisfaction Criteria**:
  - CI job runs staticcheck (and golangci-lint if configured) with project-appropriate config.
  - Baseline excludes are documented and minimized.
- **Validation Criteria**:
  - CI static analysis job passes or fails builds on findings.
  - Token audit shows `[IMPL:STATICCHECK_SETUP]` annotations.
- **Architecture**: See `architecture-decisions.md` § Static Analysis Policy [ARCH:STATIC_ANALYSIS_POLICY]
- **Implementation**: See `implementation-decisions.md` § Staticcheck Setup [IMPL:STATICCHECK_SETUP]

**Status**: ✅ Implemented

### [REQ:RACE_TESTING] Race-Enabled Tests

**Priority: P1 (Important)**

- **Description**: Provide a CI job that runs `go test -race ./...` to detect concurrency issues.
- **Rationale**: Race detection is critical for terminal UI and concurrency-heavy code paths.
- **Satisfaction Criteria**:
  - Dedicated workflow job running `go test -race` against supported OS/arch.
  - Resource/timeouts tuned to complete reliably.
- **Validation Criteria**:
  - Passing race job in CI.
  - Documented flakes and mitigations if any.
- **Architecture**: See `architecture-decisions.md` § Race Testing Pipeline [ARCH:RACE_TESTING_PIPELINE]
- **Implementation**: See `implementation-decisions.md` § Race Job [IMPL:RACE_JOB]

**Status**: ⏳ Planned

### [REQ:UI_PRIMITIVE_TESTS] UI Widget Coverage

**Priority: P0 (Critical)**

- **Description**: Increase test coverage for UI primitives (`widget/`, `filer/`) including rendering and event handling.
- **Rationale**: Widgets drive core UX; regressions are high impact.
- **Satisfaction Criteria**:
  - Unit tests cover widget behaviors and edge cases.
  - Snapshot/render or event tests validate expected output.
  - Module boundaries and interfaces documented.
- **Validation Criteria**:
  - Tests include `[REQ:UI_PRIMITIVE_TESTS]` in names/comments.
  - Module validation evidence recorded before integration.
- **Architecture**: See `architecture-decisions.md` § UI Test Strategy [ARCH:TEST_STRATEGY_UI]
- **Implementation**: See `implementation-decisions.md` § Widget Tests [IMPL:TEST_WIDGETS]

**Status**: ⏳ Planned

### [REQ:CMD_HANDLER_TESTS] Command Handling Coverage

**Priority: P0 (Critical)**

- **Description**: Validate command/app mode handling (`app/`, `cmdline/`) with focused unit tests.
- **Rationale**: Command parsing and mode transitions are correctness-critical.
- **Satisfaction Criteria**:
  - Tests cover parsing, validation, and mode changes.
  - Error paths and edge cases exercised.
- **Validation Criteria**:
  - Tests reference `[REQ:CMD_HANDLER_TESTS]`.
  - Module validation logged prior to integration.
- **Architecture**: See `architecture-decisions.md` § Command Test Strategy [ARCH:TEST_STRATEGY_CMD]
- **Implementation**: See `implementation-decisions.md` § Command Tests [IMPL:TEST_CMDLINE]

**Status**: ⏳ Planned

### [REQ:INTEGRATION_FLOWS] Integration Flows for File Ops

**Priority: P0 (Critical)**

- **Description**: Add integration/snapshot tests for core flows: open directory, navigate, rename, delete.
- **Rationale**: Ensures end-to-end behavior remains stable before refactors.
- **Satisfaction Criteria**:
  - Integration tests exercise the listed flows with fixtures.
  - Baseline outputs (messages, UI state) captured for regression.
- **Validation Criteria**:
  - Tests reference `[REQ:INTEGRATION_FLOWS]` and assert expected outcomes.
  - Logs capture `DEBUG:`/`DIAGNOSTIC:` markers for traceability.
- **Architecture**: See `architecture-decisions.md` § Integration Test Strategy [ARCH:TEST_STRATEGY_INTEGRATION]
- **Implementation**: See `implementation-decisions.md` § Integration Flow Tests [IMPL:TEST_INTEGRATION_FLOWS]

**Status**: ⏳ Planned

### [REQ:CONFIGURABLE_STATE_PATHS] Configurable State & History Persistence

**Priority: P0 (Critical)**

- **Description**: Goful must allow operators to redirect the persisted state (`state.json`) and cmdline history files through CLI flags or environment variables so multiple instances (CI, sandboxes, shared machines) can run without clobbering the default `~/.goful` data.
- **Rationale**: Hard-coded persistence paths prevent isolated test runs and multi-instance workflows; providing overrides enables deterministic automation and safer experimentation.
- **Satisfaction Criteria**:
  - CLI flags (e.g., `-state`, `-history`) override all other sources for the respective files.
  - Environment variables (e.g., `GOFUL_STATE_PATH`, `GOFUL_HISTORY_PATH`) are honored when flags are unset.
  - Default paths remain `~/.goful/state.json` and `~/.goful/history/shell` when neither flags nor environment overrides are provided.
  - Override values expand `~`, create necessary directories on first save, and are logged with `DEBUG:` output referencing `[IMPL:STATE_PATH_RESOLVER]`.
- **Validation Criteria**:
  - Unit tests exercise flag/env/default precedence plus path expansion edge cases.
  - Integration tests verify resolved paths are passed to `filer.NewFromState`, `filer.SaveState`, `cmdline.LoadHistory`, and `cmdline.SaveHistory`.
  - README and developer docs describe the flags and environment variables with their precedence order.
- **Architecture**: See `architecture-decisions.md` § State Path Selection [ARCH:STATE_PATH_SELECTION]
- **Implementation**: See `implementation-decisions.md` § State Path Resolver [IMPL:STATE_PATH_RESOLVER]

**Status**: ⏳ Planned
### [REQ:WINDOW_MACRO_ENUMERATION] External Command Window Enumeration

**Priority: P1 (Important)**

- **Description**: External command macros must expose the full set of visible directories so operators can pass every window to shell scripts. `%D@` appends the other window paths (relative order starting from the next window and wrapping). `%~D@` accepts the historical “non-quote” modifier for symmetry with other macros but each path is still escaped so multi-word directories remain safe when injected into the target command line. `echo %D %D@` therefore prints all window directories with the focused window first.
- **Rationale**: Bulk copy/move workflows depend on knowing all workspace paths. Today `%D2` exposes only the next window, so automation that needs >2 windows requires manual re-entry. Enumerating the remaining windows keeps macros self-contained and removes repetitive typing.
- **Satisfaction Criteria**:
  - `%D@` expands to a space-separated list of every other directory path in deterministic order (start with next window, then wrap through the rest). When only one window is open, the expansion is empty.
  - `%~D@` emits the same ordering and, despite the `~` modifier, individually quotes each path just like `%D@` so directories with spaces remain safe.
  - `%D@` respects the same macro parser features as `%D` (supports escaping, `%~~` safeguards, etc.) and can be combined with other text inside commands.
  - Both the requirement and the README document the new placeholder so `external-command` users can discover it.
  - A regression test proves `expandMacro("echo %D %D@")` covers all window paths with the focused directory appearing only once at the beginning.
- **Validation Criteria**:
  - Pure helper tests validate that the path enumeration logic handles 1–4 windows, wrapping order, and quoting rules.
  - `app/spawn_test.go` exercises `%D@` and `%~D@` end-to-end, including escaping/quoting scenarios that prove every path is shell-safe.
  - Token validation confirms `[REQ:WINDOW_MACRO_ENUMERATION]`, `[ARCH:WINDOW_MACRO_ENUMERATION]`, and `[IMPL:WINDOW_MACRO_ENUMERATION]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Window Macro Enumeration [ARCH:WINDOW_MACRO_ENUMERATION]
- **Implementation**: See `implementation-decisions.md` § Window Macro Enumeration [IMPL:WINDOW_MACRO_ENUMERATION]

**Status**: ⏳ Planned

### [REQ:EXTERNAL_COMMAND_CONFIG] Configurable External Commands

**Priority: P1 (Important)**

- **Description**: Operators must be able to customize the `external-command` menu keys, labels, and shell payloads by editing a configuration file (flag/env/default path) instead of recompiling goful. File-provided commands are **prepended** ahead of the compiled Windows/POSIX defaults by default so custom shortcuts appear at the top of the menu while legacy entries remain available, and the file may optionally opt out of inheritance when operators need a clean slate.
- **Rationale**: The current hard-coded menu makes it impossible to keep personal automation or team-standard tooling in sync without editing Go sources. Moving the definitions into a JSON or YAML file unblocks scripted provisioning, keeps `%D@`-style macros discoverable, allows secure environments to remove dangerous defaults, and preserves historical ergonomics by inheriting built-in commands unless explicitly disabled.
- **Satisfaction Criteria**:
  - CLI flag `-commands` overrides the configuration path; environment variable `GOFUL_COMMANDS_FILE` is honored when the flag is unset, and defaults fall back to `~/.goful/external_commands.yaml`.
  - When a configuration file is present, the loader prepends its entries before the compiled defaults so customized commands appear first while legacy shortcuts remain unless the file explicitly signals “replace defaults” (e.g., `inheritDefaults: false`).
  - The inheritance option defaults to “prepend defaults after custom entries” and can be toggled per file to drop the compiled defaults entirely for locked-down environments.
  - The loader falls back to shipped defaults (POSIX/Windows) when the config file is missing or invalid so first-run behavior matches the legacy menu.
  - Each command entry supports `menu`, `key`, `label`, and either a shell `command` string or a `runMenu` target, plus optional `offset`, optional `platforms` (GOOS list), and a `disabled` flag.
    - Duplicate `menu/key` combinations, missing required fields, or unsupported platforms are rejected with `message.Errorf`, and diagnostics mention `[IMPL:EXTERNAL_COMMAND_LOADER]`.
    - Menu binding reuses resolved entries and exposes the same `g.Shell` offsets so caret placement matches historical commands.
- **Validation Criteria**:
  - Unit tests cover path precedence for the config file (flag/env/default) and emit `[REQ:MODULE_VALIDATION]` evidence for the resolver.
  - Loader tests cover JSON/YAML decoding, default fallback, duplicate/invalid entry rejection, platform filtering, disabled entries, **and the default prepend vs. optional replace behavior**.
  - Binder tests prove menu arguments are generated deterministically and cursor offsets reach `g.Shell` correctly.
  - README/CONTRIBUTING describe the file format, macros, prepend-by-default behavior, and override steps with `[ARCH:EXTERNAL_COMMAND_REGISTRY]` references.
- **Architecture**: See `architecture-decisions.md` § External Command Registry [ARCH:EXTERNAL_COMMAND_REGISTRY]
- **Implementation**: See `implementation-decisions.md` § External Command Loader/Binding [IMPL:EXTERNAL_COMMAND_LOADER], [IMPL:EXTERNAL_COMMAND_BINDER]

**Status**: ⏳ Planned

### [REQ:TERMINAL_PORTABILITY] Cross-Platform Terminal Launcher

**Priority: P0 (Critical)**

- **Description**: Goful must launch foreground commands in an OS-appropriate terminal so macOS, Linux desktops, and tmux users all see a usable window without editing Go code. The launcher must expose an override hook (`GOFUL_TERMINAL_CMD`) and retain the historical pause tail so commands remain visible until acknowledged.
- **Rationale**: Recent dependency upgrades broke macOS Terminal invocation entirely, and Linux users increasingly prefer alternate emulators. Centralizing selection logic keeps behaviour testable and prevents regressions when upstream terminals change.
- **Satisfaction Criteria**:
  - `tmux`/`screen` sessions always use `tmux new-window -n <cmd>` regardless of OS.
  - macOS launches via AppleScript, reusing the historical payload plus pause tail while injecting the focused directory ahead of the command.
  - Linux desktops default to gnome-terminal with the legacy title escape, and overrides (e.g., `GOFUL_TERMINAL_CMD="alacritty -e"`) insert before the `bash -c` payload.
  - `GOFUL_DEBUG_TERMINAL=1` emits `DEBUG: [IMPL:TERMINAL_ADAPTER]` logs describing the branch taken.
  - README and CONTRIBUTING include guidance for macOS behaviour, overrides, and troubleshooting, all tagged with `[REQ:TERMINAL_PORTABILITY]`.
- **Validation Criteria**:
  - Unit tests cover override parsing plus Linux, macOS, and tmux branches with `[REQ:TERMINAL_PORTABILITY]` suffixes.
  - `terminalcmd.Apply` tests prove `g.ConfigTerminal` receives the expected command slices and re-reads the focused directory each time.
  - Manual verification follows `[PROC:TERMINAL_VALIDATION]` on real macOS Terminal.app, Linux desktops, and tmux sessions.
  - README/CONTRIBUTING references are updated whenever behaviour changes.
- **Architecture**: See `architecture-decisions.md` § Terminal Launcher Abstraction [ARCH:TERMINAL_LAUNCHER]
- **Implementation**: See `implementation-decisions.md` § Terminal Adapter Module [IMPL:TERMINAL_ADAPTER]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-04)**:
- `go test ./terminalcmd` (darwin/arm64, Go 1.24.3) covering `TestCommandFactory*`, `TestParseOverride`, and `TestApply*` for all selection branches.
- Manual checklist documented via `[PROC:TERMINAL_VALIDATION]` (macOS Terminal, Linux desktop, and tmux) referenced from `stdd/tasks.md` for operator execution.

### [REQ:TERMINAL_CWD] macOS Terminal Working Directory

**Priority: P0 (Critical)**

- **Description**: When goful opens macOS Terminal.app for a foreground command, the launched session must start inside the focused filer directory (the `%D` macro) so relative commands align with the window the user targeted, without needing an explicit `cd`.
- **Rationale**: `%D` already represents “the current window path” during macro expansion, but Terminal.app historically opened in `$HOME`, forcing users to prepend `cd %D`. Automating the `cd` step preserves intent, reduces mistakes, and keeps macOS parity with Linux terminals that inherit the current process directory.
- **Satisfaction Criteria**:
  - The Terminal adapter injects `cd <focusedDir>;` ahead of every macOS foreground command, with paths safely quoted (spaces, unicode).
  - Overrides via `GOFUL_TERMINAL_CMD` retain the same preamble so alternative macOS terminals behave consistently.
  - Behavior is transparent—users can still override by issuing their own `cd`, but the default always matches `%D`.
  - README/CONTRIBUTING describe the macOS working-directory guarantee.
- **Validation Criteria**:
  - Unit tests assert macOS command payloads include the `cd` preamble while Linux/tmux branches remain unchanged.
  - Configurator/integration tests prove the adapter re-evaluates the focused directory each time `g.ConfigTerminal` closure runs.
  - Manual validation confirms Terminal.app opens in the correct directory for different panes.
  - Token validation shows `[IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_CWD]` across code/docs/tests.
- **Architecture**: See `architecture-decisions.md` § Terminal Launcher Abstraction [ARCH:TERMINAL_LAUNCHER]
- **Implementation**: See `implementation-decisions.md` § Terminal Adapter Module [IMPL:TERMINAL_ADAPTER]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-04)**:
- `TestCommandFactoryDarwin_REQ_TERMINAL_PORTABILITY` and `TestApplyDarwinCwd_REQ_TERMINAL_CWD` (`terminalcmd/factory_test.go`) confirm the `cd` preamble is injected and re-computed per invocation.
- Manual checklist `[PROC:TERMINAL_VALIDATION]` documents macOS Terminal steps (outside and inside tmux) so operators can confirm the working-directory guarantee on physical hardware.

### [REQ:ARCH_DOCUMENTATION] Architecture Guide

**Priority: P1 (Important)**

- **Description**: Provide `ARCHITECTURE.md` explaining main packages (UI widgets, file ops, app/mode) and data flow.
- **Rationale**: New contributors need a concise mental model before changes.
- **Satisfaction Criteria**:
  - Document outlines package responsibilities and data flow diagrams/text.
  - Cross-references semantic tokens and key modules.
  - `ARCHITECTURE.md` published (2026-01-01) and linked from README/CONTRIBUTING so contributors can find it quickly.
- **Validation Criteria**:
  - Document reviewed for completeness; tokens audited.
  - README/CONTRIBUTING cross-links verified as part of `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh`.
- **Architecture**: See `architecture-decisions.md` § Docs Structure [ARCH:DOCS_STRUCTURE]
- **Implementation**: See `implementation-decisions.md` § Architecture Guide [IMPL:DOC_ARCH_GUIDE]

**Status**: ⏳ Planned

### [REQ:CONTRIBUTING_GUIDE] Contributor Standards

**Priority: P1 (Important)**

- **Description**: Add `CONTRIBUTING.md` covering coding standards, branching, review expectations, and token usage.
- **Rationale**: Align contributors on process and quality gates.
- **Satisfaction Criteria**:
  - Document includes coding standards, branch/review flow, test/lint expectations, token guidance.
  - References CI and Makefile targets.
  - `CONTRIBUTING.md` (2026-01-01) published with workflow checklist and debug/logging policy, and linked from README.
- **Validation Criteria**:
  - Document linked from README and tokens audited.
  - Maintainers can follow the checklist end-to-end (fmt → vet → test → token validation) without missing steps.
- **Architecture**: See `architecture-decisions.md` § Contribution Process [ARCH:CONTRIBUTION_PROCESS]
- **Implementation**: See `implementation-decisions.md` § CONTRIBUTING Guide [IMPL:DOC_CONTRIBUTING]

**Status**: ⏳ Planned

### [REQ:RELEASE_BUILD_MATRIX] Reproducible Builds

**Priority: P1 (Important)**

- **Description**: Provide Makefile targets and CI matrix to build static binaries across GOOS/GOARCH.
- **Rationale**: Ensures reproducible releases and cross-platform coverage.
- **Satisfaction Criteria**:
  - Makefile exposes `lint`, `test`, `release`, and `clean-release` targets; release target emits CGO-disabled binaries + `.sha256` digests into `dist/`.
  - CI workflow contains `release-matrix` job covering at least linux/amd64, linux/arm64, and darwin/arm64 using the Makefile target.
  - Tag-triggered `release.yml` workflow reuses the same Makefile target, uploads artifacts/checksums to GitHub Releases, and logs deterministic filenames (`goful_${GOOS}_${GOARCH}`) and SHA256 digests.
- **Validation Criteria**:
  - Local `make release PLATFORM=$(go env GOOS)/$(go env GOARCH)` succeeds and only produces expected artifacts.
  - CI + release workflows finish successfully with uploaded assets and logged checksums.
  - Token audit shows `[IMPL:MAKE_RELEASE_TARGETS]` references across Makefile + workflows, and `./scripts/validate_tokens.sh` passes after changes.
- **Validation Evidence (2026-01-01)**:
  - `make release PLATFORM=darwin/arm64` → `DIAGNOSTIC: [IMPL:MAKE_RELEASE_TARGETS] ... sha256 ad7db0a0... dist/goful_darwin_arm64`
  - `./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 130 token references across 44 files.`
- **Architecture**: See `architecture-decisions.md` § Build Matrix [ARCH:BUILD_MATRIX]
- **Implementation**: See `implementation-decisions.md` § Release Targets [IMPL:MAKE_RELEASE_TARGETS]

**Status**: ✅ Implemented

### [REQ:BEHAVIOR_BASELINE] Baseline Behavior Capture

**Priority: P1 (Important)**

- **Description**: Capture key interactions (keyboard mappings/modes) in automated tests or scripts as a pre-refactor baseline.
- **Rationale**: Protects against behavior drift during major changes.
- **Satisfaction Criteria**:
  - Tests or scripts record current keyboard mappings/modes and expected outputs.
  - Stored fixtures serve as comparison points.
  - `main_keymap_test.go` (`KeymapBaselineSuite`) enumerates canonical filer/cmdline/finder/completion/menu chords with `[TEST:KEYMAP_BASELINE]`.
- **Validation Criteria**:
  - Baseline tests run in CI and gate changes.
  - Documentation lists captured interactions.
  - `go test ./...` (2026-01-01) proves the baseline suite passes before integration.
- **Architecture**: See `architecture-decisions.md` § Baseline Capture [ARCH:BASELINE_CAPTURE]
- **Implementation**: See `implementation-decisions.md` § Baseline Snapshots [IMPL:BASELINE_SNAPSHOTS]

**Status**: ⏳ Planned

### [REQ:DEBT_TRIAGE] Debt and Risk Tracking

**Priority: P1 (Important)**

- **Description**: Triage known pain points (error handling, cross-platform quirks), open issues, and annotate risky areas with TODOs and owners.
- **Rationale**: Makes risk visible before refactors and guides prioritization.
- **Satisfaction Criteria**:
  - Issues/TODOs documented with owners.
  - Hotspots identified with inline `[IMPL:DEBT_TRACKING]` comments.
- **Validation Criteria**:
  - Debt list reviewed and linked from tasks.
  - Token audit confirms references in code/docs.
- **Architecture**: See `architecture-decisions.md` § Debt Management [ARCH:DEBT_MANAGEMENT]
- **Implementation**: See `implementation-decisions.md` § Debt Tracking [IMPL:DEBT_TRACKING]

**Status**: ⏳ Planned

### [REQ:EVENT_LOOP_SHUTDOWN] Event Poller Shutdown Control

**Priority: P0 (Critical)**

- **Description**: Ensure the UI event poller terminates cleanly whenever `app.Goful.Run` exits so goroutines and channels are not left running after the UI closes.
- **Rationale**: The current poller spins forever (`widget.PollEvent` loop) and continues sending to `g.event` even after shutdown, leaking goroutines and wasting CPU; long-running sessions eventually stall or panic.
- **Satisfaction Criteria**:
  - The event loop observes an explicit stop signal (context cancellation or quit channel) and terminates within a bounded timeout when `Run` exits.
  - Pending events are drained or discarded safely without writing to closed channels.
  - Shutdown emits `DEBUG: [IMPL:EVENT_LOOP_SHUTDOWN]` traces so operators can confirm the branch executed.
  - Manual flows described in `[PROC:TERMINAL_VALIDATION]` remain unaffected—the poller shutdown must not regress terminal launching behavior.
- **Validation Criteria**:
  - Unit tests cover the poller with a fake `widget.Poller` to assert it stops when signaled and does not panic when events arrive after stop.
  - Integration-style tests simulate `g.Run` start/stop and verify no goroutines/leaked channels remain (using `testing`'s goroutine leak detection or instrumentation hooks).
  - Debt log item D1 is updated with mitigation notes referencing this requirement once validation passes.
- **Architecture**: See `architecture-decisions.md` § Event Loop Shutdown [ARCH:EVENT_LOOP_SHUTDOWN]
- **Implementation**: See `implementation-decisions.md` § Event Loop Shutdown Controller [IMPL:EVENT_LOOP_SHUTDOWN]

**Status**: ⏳ Planned

### [REQ:IDENTIFIER] Requirement Name

**Priority: P0 (Critical) | P1 (Important) | P2 (Nice-to-have) | P3 (Future)**

- **Description**: What the requirement specifies
- **Rationale**: Why the requirement exists
- **Satisfaction Criteria**:
  - How we know the requirement is satisfied
  - Acceptance criteria
  - Success conditions
- **Validation Criteria**: 
  - How we verify/validate the requirement is met
  - Testing approach
  - Verification methods
  - Success metrics
- **Architecture**: See `architecture-decisions.md` § Decision Name [ARCH:IDENTIFIER]
- **Implementation**: See `implementation-decisions.md` § Implementation Name [IMPL:IDENTIFIER]

**Status**: ✅ Implemented | ⏳ Planned
```

## Notes

- All requirements MUST be documented here with `[REQ:*]` tokens
- Requirements describe WHAT the system should do and WHY, not HOW
- Requirements MUST NOT describe bugs or implementation details
- **Language-Agnostic Requirements**: Requirements MUST be language-agnostic. Language selection, runtime choices, and language-specific implementation details belong in architecture decisions (`architecture-decisions.md`) or implementation decisions (`implementation-decisions.md`), NOT in requirements. The ONLY exception is when language selection is itself a specific requirement (e.g., `[REQ:USE_PYTHON]` for a Python-specific project requirement). When documenting requirements, focus on behavior and capabilities, not on how they are implemented in a specific language.

## Future Enhancements (Out of Scope)

The following features are documented but marked as future enhancements:
- Each requirement should cross-reference architecture and implementation decisions

---

### Core Functionality

### [REQ:IDENTIFIER] Requirement Name

**Priority: P0 (Critical) | P1 (Important) | P2 (Nice-to-have) | P3 (Future)**

- **Description**: What the requirement specifies
- **Rationale**: Why the requirement exists
- **Satisfaction Criteria**:
  - How we know the requirement is satisfied
  - Acceptance criteria
  - Success conditions
- **Validation Criteria**: 
  - How we verify/validate the requirement is met
  - Testing approach
  - Verification methods
  - Success metrics
- **Architecture**: See `architecture-decisions.md` § Decision Name [ARCH:IDENTIFIER]
- **Implementation**: See `implementation-decisions.md` § Implementation Name [IMPL:IDENTIFIER]

**Status**: ✅ Implemented | ⏳ Planned
```

### 2. [REQ:ANOTHER_FEATURE] Another Feature Name

**Priority: P0 (Critical)**

- **Description**: Description of the feature
- **Rationale**: Why it's needed
- **Satisfaction Criteria** (How we know the requirement is satisfied):
  - Criterion 1
  - Criterion 2
- **Validation Criteria** (How we verify/validate the requirement is met):
  - Validation method 1
  - Validation method 2

**Status**: ⏳ Planned

## Non-Functional Requirements

### 1. Performance [REQ:PERFORMANCE]
- Requirement description
- Metrics or targets

### 2. Reliability [REQ:RELIABILITY]
- Requirement description
- Availability targets

### 3. Maintainability [REQ:MAINTAINABILITY]
- Requirement description
- Code quality standards

### 4. Usability [REQ:USABILITY]
- Requirement description
- User experience goals

## Edge Cases to Handle

1. **Edge Case 1**
   - Description
   - Expected behavior

2. **Edge Case 2**
   - Description
   - Expected behavior

## Future Enhancements (Out of Scope)

The following features are documented but marked as future enhancements:
- Feature 1
- Feature 2
- Feature 3

These may be considered for future iterations but are not required for the initial implementation.

