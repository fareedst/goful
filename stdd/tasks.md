# Tasks and Incomplete Subtasks

**STDD Methodology Version**: 1.1.0

## Overview
This document tracks all tasks and subtasks for implementing this project. Tasks are organized by priority and implementation phase.

## Priority Levels

- **P0 (Critical)**: Must have - Core functionality, blocks other work
- **P1 (Important)**: Should have - Enhanced functionality, better error handling
- **P2 (Nice-to-Have)**: Could have - UI/UX improvements, convenience features
- **P3 (Future)**: Won't have now - Deferred features, experimental ideas

## Task Format

```markdown
## P0: Task Name [REQ:IDENTIFIER] [ARCH:IDENTIFIER] [IMPL:IDENTIFIER]

**Status**: 🟡 In Progress | ✅ Complete | ⏸️ Blocked | ⏳ Pending

**Description**: Brief description of what this task accomplishes.

**Dependencies**: List of other tasks/tokens this depends on.

**Subtasks**:
- [ ] Subtask 1 [REQ:X] [IMPL:Y]
- [ ] Subtask 2 [REQ:X] [IMPL:Z]
- [ ] Subtask 3 [TEST:X]
- [ ] Token audit & validation [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [ ] All subtasks complete
- [ ] Code implements requirement
- [ ] Tests pass with semantic token references
- [ ] Documentation updated
- [ ] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Priority Rationale**: Why this is P0/P1/P2/P3
```

## Task Management Rules

1. **Subtasks are Temporary**
   - Subtasks exist only while the parent task is in progress
   - Remove subtasks when parent task completes

2. **Priority Must Be Justified**
   - Each task must have a priority rationale
   - Priorities follow: Tests/Code/Functions > DX > Infrastructure > Security

3. **Semantic Token References Required**
   - Every task MUST reference at least one semantic token
   - Cross-reference to related tokens

4. **Token Audits & Validation Required**
   - Every task must include a `[PROC:TOKEN_AUDIT]` subtask and capture its result
   - `./scripts/validate_tokens.sh` (or repo-specific equivalent) must run before closing the task, with results logged under `[PROC:TOKEN_VALIDATION]`

5. **Completion Criteria Must Be Met**
   - All criteria must be checked before marking complete
   - Documentation must be updated

## Task Status Icons

- 🟡 **In Progress**: Actively being worked on
- ✅ **Complete**: All criteria met, subtasks removed
- ⏸️ **Blocked**: Waiting on dependency
- ⏳ **Pending**: Not yet started

## Active Tasks

## P0: Setup STDD Methodology [REQ:STDD_SETUP] [ARCH:STDD_STRUCTURE] [IMPL:STDD_FILES]

**Status**: ✅ Complete

**Description**: Initialize the project with the STDD directory structure and documentation files.

**Dependencies**: None

**Subtasks**:
- [x] Create `stdd/` directory
- [x] Instantiate documentation files from templates
- [x] Update `.cursorrules`
- [x] Register semantic tokens

**Completion Criteria**:
- [x] All subtasks complete
- [x] Code implements requirement
- [x] Documentation updated

**Priority Rationale**: P0 because this is the foundation for all future work.

## P0: Promote Processes into Core Methodology [REQ:STDD_SETUP] [ARCH:STDD_STRUCTURE] [IMPL:STDD_FILES]

**Status**: ✅ Complete

**Description**: Align every methodology reference (docs, templates, registry files) to STDD v1.1.0 after elevating Processes into the primary STDD workflow.

**Dependencies**: None

**Subtasks**:
- [x] Update STDD version references across methodology docs and guides
- [x] Update all template files and project copies with the new version marker
- [x] Refresh `VERSION`, `CHANGELOG.md`, and supporting metadata to announce v1.1.0

**Completion Criteria**:
- [x] All semantic references cite STDD v1.1.0
- [x] VERSION file, changelog, and documentation agree on the new version
- [x] Tasks and supporting docs reflect completion of this work

**Priority Rationale**: Processes are now a primary STDD concern; all consumers must see the v1.1.0 upgrade immediately to maintain alignment.

## Phase 2: Core Components

### Task 2.1: Configurable State & History Paths
**Status:** ✅ Complete  
**Priority:** P0 (Critical)  
**Semantic Tokens:** `[REQ:CONFIGURABLE_STATE_PATHS]`, `[ARCH:STATE_PATH_SELECTION]`, `[IMPL:STATE_PATH_RESOLVER]`

**Description**: Allow goful to honor CLI flags and environment overrides for the persisted state (`state.json`) and shell history so multiple instances or sandboxes do not clobber the default `~/.goful` files. This task wires the resolver into startup, documents the behavior, and proves it via unit + integration tests.

**Dependencies**: [ARCH:MODULE_VALIDATION], [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Expand requirement + architecture + implementation docs with the new tokens before coding [REQ:CONFIGURABLE_STATE_PATHS]
- [x] Identify and document module boundaries (`PathResolver`, `BootstrapPaths`) plus validation criteria [REQ:MODULE_VALIDATION]
- [x] Develop Module 1 `PathResolver` (pure precedence logic) independently [IMPL:STATE_PATH_RESOLVER]
- [x] Validate Module 1 with unit tests covering flag/env/default precedence + edge cases [REQ:MODULE_VALIDATION]
- [x] Develop Module 2 `BootstrapPaths` (flag parsing + env wiring into goful startup) independently
- [x] Validate Module 2 with integration-style tests exercising combined struct behavior [REQ:MODULE_VALIDATION]
- [x] Integrate validated modules into `main.go`, update README, and ensure debug output references the new tokens
- [x] Write end-to-end verification (resolver + bootstrap) to assert both paths flow into filer/cmdline entry points [REQ:CONFIGURABLE_STATE_PATHS]
- [x] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` and record outcomes (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 70 token references across 40 files.`) [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] Modules documented with interfaces + validation evidence
- [x] Module 1 & 2 validation suites pass independently before integration
- [x] Integration wiring + documentation merged with semantic tokens
- [x] README and developer docs describe flags/env overrides
- [x] Token audit + validation logged with command output

**Priority Rationale**: P0 because without configurable paths, multi-instance workflows overwrite global files and block tests/sandboxes that must isolate state.

## P0: Modernize Toolchain and Dependencies [REQ:GO_TOOLCHAIN_LTS] [REQ:DEPENDENCY_REFRESH] [ARCH:GO_RUNTIME_STRATEGY] [ARCH:DEPENDENCY_POLICY] [IMPL:GO_MOD_UPDATE] [IMPL:DEP_BUMP]

**Status**: ✅ Complete

**Description**: Move to current Go LTS, refresh deps (tcell, x/*), and tidy module graph.

**Dependencies**: None

**Subtasks**:
- [x] Decide target Go LTS and document in `go.mod` [REQ:GO_TOOLCHAIN_LTS] [IMPL:GO_MOD_UPDATE]
- [x] Update deps and run `go mod tidy` [REQ:DEPENDENCY_REFRESH] [IMPL:DEP_BUMP]
- [x] Smoke test fmt/vet/test locally with new toolchain [REQ:GO_TOOLCHAIN_LTS]
- [x] Document any shims/breakages discovered (none required) [REQ:DEPENDENCY_REFRESH]
- [x] Run `[PROC:TOKEN_AUDIT]`
- [x] Run `[PROC:TOKEN_VALIDATION]` (`./scripts/validate_tokens.sh`)

**Completion Criteria**:
- [x] `go.mod`/`go.sum` updated to LTS and refreshed deps
- [x] Tests/vet pass on new toolchain
- [x] Docs updated with decisions
- [x] Token audit + validation recorded (`./scripts/validate_tokens.sh`)

**Priority Rationale**: P0 to unblock all downstream CI/tests and security fixes.

## P0: Restore Quit Dialog Return Behavior [REQ:QUIT_DIALOG_DEFAULT] [ARCH:QUIT_DIALOG_KEYS] [IMPL:QUIT_DIALOG_ENTER]

**Status**: ✅ Complete

**Description**: Ensure the quit dialog (and other cmdline modes) accept the Return/Enter key as the default confirmation after recent tcell changes.

**Dependencies**: None (regression fix, but ties into `[REQ:MODULE_VALIDATION]` for translator module)

**Subtasks**:
- [x] Identify affected modules and document translator contract [REQ:QUIT_DIALOG_DEFAULT] [REQ:MODULE_VALIDATION]
- [x] Update `widget.EventToString` mapping for Return/Enter [ARCH:QUIT_DIALOG_KEYS] [IMPL:QUIT_DIALOG_ENTER]
- [x] Add regression tests for `EventToString` handling [REQ:QUIT_DIALOG_DEFAULT] [IMPL:QUIT_DIALOG_ENTER]
- [x] Document manual Return-validation guidance for operators (interactive TUI verification required on a real terminal) [REQ:QUIT_DIALOG_DEFAULT]
- [x] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` and record outcome [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] Translator + cmdline modules validated independently
- [x] Tests cover Return → `C-m` mapping
- [x] Manual verification guidance logged (operator to confirm on-device)
- [x] Documentation and tokens updated end-to-end
- [x] `[PROC:TOKEN_AUDIT]` / `[PROC:TOKEN_VALIDATION]` recorded

**Priority Rationale**: P0 because users cannot exit the application using standard key flow, effectively trapping sessions.

## P0: CI & Static Analysis Foundation [REQ:CI_PIPELINE_CORE] [REQ:STATIC_ANALYSIS] [REQ:RACE_TESTING] [ARCH:CI_PIPELINE] [ARCH:STATIC_ANALYSIS_POLICY] [ARCH:RACE_TESTING_PIPELINE] [IMPL:CI_WORKFLOW] [IMPL:STATICCHECK_SETUP] [IMPL:RACE_JOB]

**Status**: ✅ Complete

**Description**: Establish GitHub Actions for fmt/vet/tests, static analysis, and race job.

**Dependencies**: Modernize Toolchain and Dependencies

**Subtasks**:
- [x] Add fmt/vet/test workflow with cache [REQ:CI_PIPELINE_CORE] [IMPL:CI_WORKFLOW]
- [x] Add staticcheck (and optional golangci-lint) job [REQ:STATIC_ANALYSIS] [IMPL:STATICCHECK_SETUP]
- [x] Add race-enabled test job [REQ:RACE_TESTING] [IMPL:RACE_JOB]
- [x] Ensure jobs reference matching Go version [REQ:GO_TOOLCHAIN_LTS]
- [x] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]`

**Completion Criteria**:
- [x] CI runs fmt/vet/test/staticcheck/race on PRs
- [x] Jobs pass on target branches (workflow validated locally; run on PRs)
- [x] Token audit + validation recorded (`./scripts/validate_tokens.sh`)

**Priority Rationale**: P0 to gate all future changes with automated checks.

## P0: Test Coverage for UI/Commands/Flows [REQ:UI_PRIMITIVE_TESTS] [REQ:CMD_HANDLER_TESTS] [REQ:INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_UI] [ARCH:TEST_STRATEGY_CMD] [ARCH:TEST_STRATEGY_INTEGRATION] [IMPL:TEST_WIDGETS] [IMPL:TEST_CMDLINE] [IMPL:TEST_INTEGRATION_FLOWS]

**Status**: ✅ Complete

**Description**: Add coverage for widgets/filer, command handling, and integration flows (open/navigate/rename/delete).

**Dependencies**: CI & Static Analysis Foundation

**Subtasks**:
- [x] Identify modules and validation criteria per area [REQ:MODULE_VALIDATION]
- [x] Add widget/filer unit/snapshot tests [REQ:UI_PRIMITIVE_TESTS] [IMPL:TEST_WIDGETS]
- [x] Add command/app mode tests [REQ:CMD_HANDLER_TESTS] [IMPL:TEST_CMDLINE]
- [x] Add integration flow tests with fixtures [REQ:INTEGRATION_FLOWS] [IMPL:TEST_INTEGRATION_FLOWS]
- [x] Document validation results before integration [REQ:MODULE_VALIDATION]
- [x] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]`

**Completion Criteria**:
- [x] Module validation evidence recorded
- [x] Tests cover listed areas with tokens
- [x] CI green with new coverage
- [x] Token audit + validation recorded (`./scripts/validate_tokens.sh`)

**Priority Rationale**: P0 to secure behavior before major refactors.

## P1: Docs & Baselines [REQ:ARCH_DOCUMENTATION] [REQ:CONTRIBUTING_GUIDE] [REQ:BEHAVIOR_BASELINE] [ARCH:DOCS_STRUCTURE] [ARCH:CONTRIBUTION_PROCESS] [ARCH:BASELINE_CAPTURE] [IMPL:DOC_ARCH_GUIDE] [IMPL:DOC_CONTRIBUTING] [IMPL:BASELINE_SNAPSHOTS]

**Status**: ✅ Complete

**Description**: Write `ARCHITECTURE.md`, `CONTRIBUTING.md`, and capture baseline keybindings/modes.

**Dependencies**: CI & Static Analysis Foundation

**Module Boundaries**:
- `DocArchitecture` – curates package/data-flow overview before code changes; validated by doc review + cross-references. [REQ:ARCH_DOCUMENTATION] [REQ:MODULE_VALIDATION]
- `DocContributing` – contributor workflow contract describing tooling, semantic tokens, and debug policy. [REQ:CONTRIBUTING_GUIDE] [REQ:MODULE_VALIDATION]
- `KeymapBaselineSuite` – pure Go tests snapshotting filer/cmdline/finder/menu/completion bindings without touching terminal IO. [REQ:BEHAVIOR_BASELINE] [ARCH:BASELINE_CAPTURE] [IMPL:BASELINE_SNAPSHOTS] [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Draft architecture overview with package/data flow [REQ:ARCH_DOCUMENTATION] [IMPL:DOC_ARCH_GUIDE]
- [x] Draft contributing guide with standards/review expectations [REQ:CONTRIBUTING_GUIDE] [IMPL:DOC_CONTRIBUTING]
- [x] Capture baseline interactions/keymaps in tests/scripts [REQ:BEHAVIOR_BASELINE] [IMPL:BASELINE_SNAPSHOTS]
- [x] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]` (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 89 token references across 43 files.`)

**Completion Criteria**:
- [x] Docs published and cross-linked
- [x] Baseline tests run in CI
- [x] Token audit + validation recorded

**Priority Rationale**: P1 to enable contributors and guard behavior.

## P1: Release Build Hygiene [REQ:RELEASE_BUILD_MATRIX] [ARCH:BUILD_MATRIX] [IMPL:MAKE_RELEASE_TARGETS]

**Status**: ✅ Complete

**Description**: Add Makefile targets and CI + tag-triggered workflows for reproducible static builds across GOOS/GOARCH, with deterministic filenames/digests and artifact publication.

**Dependencies**: Modernize Toolchain and Dependencies

**Module Boundaries**:
- `MakeReleaseTargets` – Makefile targets for lint/test/build plus hermetic release bundles (dist directory, CGO disabled) validated via `make release` on host OS and checksum generation. [REQ:RELEASE_BUILD_MATRIX] [IMPL:MAKE_RELEASE_TARGETS] [REQ:MODULE_VALIDATION]
- `ReleaseMatrixWorkflow` – GitHub Actions matrix job (linux/amd64, linux/arm64, darwin/arm64) that reuses Makefile targets and publishes artifacts/checksums. [REQ:RELEASE_BUILD_MATRIX] [ARCH:BUILD_MATRIX] [IMPL:MAKE_RELEASE_TARGETS] [REQ:MODULE_VALIDATION]
- `ArtifactDeterminismAudit` – checksum/filename verification ensuring deterministic naming and digests, implemented as part of Makefile + CI steps with logged SHA256 outputs. [REQ:RELEASE_BUILD_MATRIX] [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Finalize Makefile release targets + local validation run (`make release PLATFORM=darwin/arm64` → `DIAGNOSTIC: [IMPL:MAKE_RELEASE_TARGETS] ... sha256 ad7db0a0...`) [REQ:RELEASE_BUILD_MATRIX] [IMPL:MAKE_RELEASE_TARGETS]
- [x] Implement GitHub Actions release workflow (tag-triggered matrix + asset upload) [REQ:RELEASE_BUILD_MATRIX] [ARCH:BUILD_MATRIX]
- [x] Document release process in README/CONTRIBUTING + update requirements/decisions [REQ:RELEASE_BUILD_MATRIX]
- [x] Verify deterministic artifacts/checksums + run `[PROC:TOKEN_AUDIT]` / `[PROC:TOKEN_VALIDATION]` (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 130 token references across 44 files.`)

**Completion Criteria**:
- [x] Makefile + workflow committed with tokens
- [x] Matrix + release workflows pass and upload artifacts
- [x] Token audit + validation recorded

**Priority Rationale**: P1 to prepare for reproducible releases.

## P1: Debt Triage [REQ:DEBT_TRIAGE] [ARCH:DEBT_MANAGEMENT] [IMPL:DEBT_TRACKING]

**Status**: ✅ Complete

**Description**: Log known pain points (error handling, cross-platform quirks) and annotate risky areas with TODOs/owners. Snapshot recorded in `stdd/debt-log.md`.

**Dependencies**: None

**Subtasks**:
- [x] Create issue list/backlog of known risks (see `stdd/debt-log.md`) [REQ:DEBT_TRIAGE] [IMPL:DEBT_TRACKING]
- [x] Add TODOs with owners in hotspot files [REQ:DEBT_TRIAGE] [IMPL:DEBT_TRACKING]
- [x] Link debt list into docs/tasks (`architecture-decisions.md`, `implementation-decisions.md`) [REQ:DEBT_TRIAGE]
- [x] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]` (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 148 token references across 44 files.`)

**Completion Criteria**:
- [ ] Debt backlog documented and linked
- [ ] TODOs annotated in code
- [ ] Token audit + validation recorded

**Priority Rationale**: P1 to surface risks before refactors.
## P1: Window Macro Enumeration [REQ:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [IMPL:WINDOW_MACRO_ENUMERATION]

**Status**: ✅ Complete

**Description**: Add `%D@` / `%~D@` external-command macros so scripts can enumerate every other workspace directory with deterministic ordering and quoting rules.

**Dependencies**: [REQ:MODULE_VALIDATION]

**Completion Criteria**:
- [x] `%D@` / `%~D@` helper modules implemented with deterministic ordering and quoting.
- [x] Unit + integration tests document module validation evidence.
- [x] README reflects the new macro.
- [x] Tokens registered + validation logs captured.

**Validation Evidence**: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 154 token references across 44 files.`

**Priority Rationale**: P1 because cross-window automation is a frequently requested workflow improvement, but it does not block basic navigation.

## P1: `%~D@` Escaping Regression [REQ:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [IMPL:WINDOW_MACRO_ENUMERATION]

**Status**: ✅ Complete

**Description**: Ensure `%~D@` emits shell-escaped directory paths so multi-window commands remain safe even when the non-quote modifier is used. Aligns `%~D@` behavior with `%D` escaping rules to prevent spaces from splitting arguments.

**Dependencies**: [REQ:WINDOW_MACRO_ENUMERATION], [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Update requirements + architecture + implementation docs to describe the escaped `%~D@` contract [REQ:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [IMPL:WINDOW_MACRO_ENUMERATION]
- [x] Adjust `app/spawn.go` enumeration helpers to quote `%~D@` outputs and add diagnostics if needed [REQ:WINDOW_MACRO_ENUMERATION] [IMPL:WINDOW_MACRO_ENUMERATION]
- [x] Refresh unit/integration tests covering `%~D@` to assert escaped results [REQ:WINDOW_MACRO_ENUMERATION] [IMPL:WINDOW_MACRO_ENUMERATION]
- [x] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` and log outcomes (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 185 token references across 46 files.`) [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] Documentation updated with the escaped `%~D@` expectation
- [x] `%~D@` emits individually quoted paths in all workspaces
- [x] Tests cover the new escaping guarantee
- [x] Token audit + validation recorded

**Priority Rationale**: P1 because broken escaping corrupts automation workflows but does not stop the application from launching.
## P0: Cross-Platform Terminal Launcher [REQ:TERMINAL_PORTABILITY] [ARCH:TERMINAL_LAUNCHER] [IMPL:TERMINAL_ADAPTER]

**Status**: ⏳ Pending

**Description**: Provide a portable terminal launcher so executing commands from goful works on macOS (Terminal.app), Linux desktops, and tmux/screen sessions without editing source. Centralize OS detection and overrides in a testable module that feeds `g.ConfigTerminal`.

**Dependencies**: [REQ:MODULE_VALIDATION], [REQ:STDD_SETUP]

**Module Boundaries**:
- `CommandFactory` – pure function that decides which terminal command slice to run based on GOOS, tmux detection, and overrides.
- `Configurator` – wires the factory output (plus tail suffix) into `g.ConfigTerminal` and emits diagnostics.
- `ManualValidationChecklist` – documents macOS + Linux verification steps after automated tests pass.

**Subtasks**:
- [x] Document requirement + architecture/implementation tokens for terminal portability [REQ:TERMINAL_PORTABILITY] [ARCH:TERMINAL_LAUNCHER] [IMPL:TERMINAL_ADAPTER]
- [x] Implement `CommandFactory` with tmux detection, overrides, Linux default, and macOS `osascript` path [IMPL:TERMINAL_ADAPTER]
- [x] Add unit tests covering Linux, macOS, tmux, and override branches (module validation) [REQ:MODULE_VALIDATION]
- [x] Implement `Configurator` glue + logging, and wire it into `g.ConfigTerminal` [IMPL:TERMINAL_ADAPTER]
- [x] Add integration tests or fakes proving `g.ConfigTerminal` receives the correct command slices [REQ:TERMINAL_PORTABILITY]
- [x] Update README/CONTRIBUTING with macOS guidance and override instructions [REQ:TERMINAL_PORTABILITY]
- [x] Inject `%D` working directory into every macOS Terminal launch (docs + tests) [REQ:TERMINAL_CWD]
- [x] Replace mac login shell flags with non-login `bash -c` to avoid hangs [REQ:TERMINAL_PORTABILITY] [IMPL:TERMINAL_ADAPTER]
- [ ] Manual validation checklist for macOS Terminal + Linux desktop runs [REQ:TERMINAL_PORTABILITY]
- [x] `[PROC:TOKEN_AUDIT]` and `./scripts/validate_tokens.sh` (`[PROC:TOKEN_VALIDATION]`) recorded after code/tests land (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 185 token references across 46 files.` on 2026-01-02)

**Completion Criteria**:
- [ ] Factory + configurator modules validated independently before integration
- [ ] `main.go` uses the new adapter; legacy gnome-terminal call removed
- [ ] Tests cover selection matrix and keep-open tail behavior
- [ ] Documentation updated with macOS instructions and overrides
- [ ] `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]` logs captured

**Priority Rationale**: P0 because macOS users currently cannot execute external commands in a terminal, blocking a core workflow.


