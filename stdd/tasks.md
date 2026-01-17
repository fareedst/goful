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

## P0: Backspace Navigation & Editing [REQ:BACKSPACE_BEHAVIOR] [ARCH:BACKSPACE_TRANSLATION] [IMPL:BACKSPACE_TRANSLATION]

**Status**: ✅ Complete

**Description**: Normalize Backspace/Delete key codes so filer panes continue to open the parent directory and prompted input widgets delete the prior character regardless of whether tcell delivers `KeyBackspace` or `KeyBackspace2`.

**Dependencies**: [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Update requirements/architecture/implementation docs plus semantic token registry for the new behavior [REQ:BACKSPACE_BEHAVIOR] [ARCH:BACKSPACE_TRANSLATION] [IMPL:BACKSPACE_TRANSLATION]
- [x] Normalize `tcell.KeyBackspace` and `tcell.KeyBackspace2` to the canonical `backspace` string inside `widget.EventToString` with debug-ready annotations [IMPL:BACKSPACE_TRANSLATION]
- [x] Add translator unit tests proving both key codes emit `backspace` and ensure baseline keymap coverage references the requirement [REQ:BACKSPACE_BEHAVIOR]
- [x] Run `[PROC:TOKEN_AUDIT]` and `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` capturing the diagnostic output (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 520 token references across 66 files.` on 2026-01-09) [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] Docs/tokens updated with cross-references
- [x] Translator change validated by unit tests and manual smoke checks
- [x] Baseline keymaps retain the `backspace` chord, preventing regressions
- [x] `[PROC:TOKEN_AUDIT]` / `[PROC:TOKEN_VALIDATION]` results recorded

**Priority Rationale**: P0 because Backspace is a primary navigation/editing key; without normalization macOS users cannot ascend directories or edit prompts, blocking fundamental workflows.

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

## P1: `%~D@` Raw Output Parity [REQ:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [IMPL:WINDOW_MACRO_ENUMERATION]

**Status**: ✅ Complete

**Description**: Re-align `%~D@` with the historical `%~` non-quote semantics so automation can opt into raw directory arguments without shell escaping, while keeping `%D@` quoted for safety by default.

**Dependencies**: [REQ:WINDOW_MACRO_ENUMERATION], [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Update requirements + architecture + implementation docs + README to describe the quoted vs. raw behavior split [REQ:WINDOW_MACRO_ENUMERATION]
- [x] Adjust `app/spawn.go` helpers so `%~D@` uses the raw formatter and `%D@` stays quoted, with debug annotations preserved [IMPL:WINDOW_MACRO_ENUMERATION]
- [x] Refresh macro tests (helpers + integration) to assert `%~D@` returns raw paths, including directories containing spaces [REQ:WINDOW_MACRO_ENUMERATION]
- [x] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` and capture the latest output in this task + implementation decisions (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 254 token references across 55 files.` on 2026-01-04) [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] Docs, README, and requirements call out the quoting difference explicitly
- [x] `%D@` expansions remain quoted while `%~D@` emits raw paths in all workspaces
- [x] Tests prove raw vs. quoted behavior (space-containing paths) and `[REQ:WINDOW_MACRO_ENUMERATION]` coverage remains intact
- [x] Token audit + validation logs recorded and linked

**Priority Rationale**: P1 because automation workflows that expect raw `%~` semantics cannot consume the quoted variant, but core navigation continues to function.

## P1: `%d@` Directory Name Enumeration [REQ:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [IMPL:WINDOW_MACRO_ENUMERATION]

**Status**: ✅ Complete

**Description**: Extend the window-enumeration macros so `%d@`/`%~d@` append the other directory **names** (basename only) with the same deterministic ordering and tilde-driven quoting rules as `%D@`.

**Dependencies**: [REQ:WINDOW_MACRO_ENUMERATION], [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Update requirements/architecture/implementation docs, README, and macro tables to describe `%d@`/`%~d@` behavior. [REQ:WINDOW_MACRO_ENUMERATION]
- [x] Add helper coverage for directory names (companions to `otherWindowDirPaths`) and annotate with semantic tokens. [IMPL:WINDOW_MACRO_ENUMERATION]
- [x] Teach `expandMacro` + tests to recognize `%d@`/`%~d@`, including quoting vs. raw name guarantees. [IMPL:WINDOW_MACRO_ENUMERATION]
- [x] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` and log the latest diagnostic output (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 260 token references across 55 files.` on 2026-01-05). [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] Documentation + README clearly explain `%d@`/`%~d@` quoting/name semantics alongside `%D@`.
- [x] `%d@` returns quoted basenames and `%~d@` returns raw basenames in deterministic order; single-window workspaces emit empty strings.
- [x] Helper + integration tests cover the new macros and reference `[REQ:WINDOW_MACRO_ENUMERATION]`.
- [x] Token audit + validation outputs recorded in this task and supporting docs.

**Priority Rationale**: P1 because automation workflows that only need directory names must currently post-process `%D@` output, adding brittle shell logic.

## P1: Filename Exclude Filter [REQ:FILER_EXCLUDE_NAMES] [ARCH:FILER_EXCLUDE_FILTER] [IMPL:FILER_EXCLUDE_RULES] [IMPL:FILER_EXCLUDE_LOADER]

**Status**: ✅ Complete

**Description**: Load a newline-delimited basename block list (flag/env/default path), hide matching entries across all filer views, and provide a runtime keystroke/menu toggle so operators can temporarily reveal excluded files without restarting goful.

**Dependencies**: [REQ:MODULE_VALIDATION], [REQ:CONFIGURABLE_STATE_PATHS] (resolver precedence reuse)

**Subtasks**:
- [x] Update `requirements.md`, `architecture-decisions.md`, `implementation-decisions.md`, and `semantic-tokens.md` with `[REQ:FILER_EXCLUDE_NAMES]`, `[ARCH:FILER_EXCLUDE_FILTER]`, `[IMPL:FILER_EXCLUDE_RULES]`, `[IMPL:FILER_EXCLUDE_LOADER]`.
- [x] Extend `configpaths.Resolver` with `-exclude-names` / `GOFUL_EXCLUDES_FILE` precedence and log provenance via `emitPathDebug`.
- [x] Implement loader/parser + diagnostics in `main.go`, wiring the View menu + dedicated keystroke toggle. [REQ:FILER_EXCLUDE_NAMES] [IMPL:FILER_EXCLUDE_LOADER]
- [x] Implement `filer`-level rule store and hook into `Directory.read`, including unit/integration tests that validate hiding + toggling. [REQ:FILER_EXCLUDE_NAMES] [IMPL:FILER_EXCLUDE_RULES]
- [x] Add tests for resolver + loader parsing (comments, blank lines, case-insensitivity) and for filer filtering/toggle behaviour. [REQ:MODULE_VALIDATION]
- [x] Run `[PROC:TOKEN_AUDIT]` and `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh`, recording diagnostic output.

**Completion Criteria**:
- [x] Docs and token registry capture the new requirement, architecture, implementation, and process notes.
- [x] Resolver, loader, runtime toggle, and filer logic are implemented with semantic token comments and debug output.
- [x] Unit + integration tests validate both modules independently before integration; README/ARCHITECTURE mention the new feature.
- [x] `[PROC:TOKEN_AUDIT]` / `[PROC:TOKEN_VALIDATION]` results recorded with command output.

**Validation Evidence**:
- `go test ./...` (darwin/arm64, Go 1.24.3) — 2026-01-07.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 302 token references across 58 files.` (2026-01-07).

**Priority Rationale**: P1 because hiding noisy files dramatically improves navigation ergonomics while remaining optional and discoverable via runtime toggle.

## P1: External Command Configuration [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER] [IMPL:EXTERNAL_COMMAND_BINDER]

**Status**: ✅ Complete

**Description**: Load `external-command` menu entries from a JSON config (flag/env/default) so operators can customize bindings without editing Go code while retaining historical defaults per platform.

**Dependencies**: [REQ:MODULE_VALIDATION], [ARCH:STATE_PATH_SELECTION]

**Subtasks**:
- [x] Update requirements + docs with `[REQ:EXTERNAL_COMMAND_CONFIG]`, `[ARCH:EXTERNAL_COMMAND_REGISTRY]`, `[IMPL:EXTERNAL_COMMAND_LOADER]`, `[IMPL:EXTERNAL_COMMAND_BINDER]`.
- [x] Extend `configpaths.Resolver` with commands path + provenance debug output.
- [x] Implement `externalcmd.Defaults/Load` with schema validation, GOOS filters, disabled entries, and fallback behavior.
- [x] Add loader unit tests covering precedence, invalid files, duplicate keys, platform gating, and disabled entries. [REQ:MODULE_VALIDATION]
- [x] Extract binder helper in `main` + binder unit tests to prove menu args + placeholder behavior. [REQ:MODULE_VALIDATION]
- [x] Wire `main.go` to new loader/binder, replace hard-coded menu, and document file format in README.
- [x] Token audit & validation logs (`./scripts/validate_tokens.sh`). [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] Requirement/architecture/implementation docs updated with cross-references.
- [x] Loader + binder modules validated independently.
- [x] `main.go` uses loader output; defaults preserved for Windows/POSIX.
- [x] README/CONTRIBUTING describe usage + overrides.
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` results recorded (`./scripts/validate_tokens.sh` output pasted into implementation decisions on 2026-01-02).

**Priority Rationale**: P1 because customization improves automation workflows and unblocks secure environments, but it does not block baseline navigation.

## P1: External Command YAML Support [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]

**Status**: ✅ Complete

**Description**: Teach the loader/binder stack to accept YAML configs in addition to JSON so teams with existing YAML workflows can manage the `external-command` menu without format conversions.

**Dependencies**: [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Add `gopkg.in/yaml.v3` dependency and extend `externalcmd.Entry` tags for YAML decoding.
- [x] Update `externalcmd.Load` to parse JSON or YAML arrays/wrapper objects plus new unit tests. [REQ:MODULE_VALIDATION]
- [x] Refresh docs (`README.md`, `stdd/*.md`, `ARCHITECTURE.md`) to describe JSON/YAML support and record the change in `tasks.md`.
- [x] Run `go test ./...` and `./scripts/validate_tokens.sh` to capture the latest `[PROC:TOKEN_VALIDATION]` evidence.

**Completion Criteria**:
- [x] Loader parses YAML files (array + `commands:`) and falls back to defaults on errors.
- [x] Documentation + requirements mention JSON/YAML parity.
- [x] Token audit/validation logs updated with the latest run (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 245 token references across 52 files.`).

**Priority Rationale**: P1 because it materially improves automation ergonomics without blocking core navigation workflows.

## P1: External Command Requirement Update [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER] [IMPL:EXTERNAL_COMMAND_BINDER]

**Status**: ✅ Complete

**Description**: Clarify that file-based external command definitions **prepend** to the built-in defaults by default, with an option to replace or suppress defaults, so requirements stay aligned with loader behavior.

**Dependencies**: [REQ:MODULE_VALIDATION] (documentation change only, but still references loader/binder modules)

**Subtasks**:
- [x] Update `[REQ:EXTERNAL_COMMAND_CONFIG]` in `stdd/requirements.md` to describe prepend semantics and optional override switch. [REQ:EXTERNAL_COMMAND_CONFIG]
- [x] Verify architecture/implementation decisions remain accurate (prepend behavior already covered) or note follow-up edits if discrepancies remain. [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]
- [x] `[PROC:TOKEN_AUDIT]` note: Documentation-only change; ensure semantic tokens already exist and no registry updates required.

**Completion Criteria**:
- [x] Requirement text published with prepend default + optional override description.
- [x] Cross-references reviewed; action items noted if further edits needed.
- [x] Task log updated with results (`DIAGNOSTIC: requirement update complete`).

**Priority Rationale**: P1 to keep documentation accurate for automation workflows relying on external command registries.

## P1: External Command Append Semantics [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_APPEND]

**Status**: ✅ Complete

**Description**: Implement the prepend-by-default behavior (with optional replacement) in code, tests, and docs while introducing dedicated semantic tokens so inheritance is auditable beyond requirements text.

**Dependencies**: [REQ:MODULE_VALIDATION] (loader/binder modules already validated; change reuses same seams)

**Subtasks**:
- [x] Update architecture/implementation docs, README, and `stdd/semantic-tokens.md` with `[IMPL:EXTERNAL_COMMAND_APPEND]` covering the inheritance toggle (prepends by default). [REQ:EXTERNAL_COMMAND_CONFIG]
- [x] Teach `externalcmd.Load` to parse `inheritDefaults` (JSON/YAML) and prepend defaults unless explicitly disabled, with inline tokens/comments. [IMPL:EXTERNAL_COMMAND_APPEND]
- [x] Add tests proving prepend vs. replace behavior, tagged with the new token. [REQ:EXTERNAL_COMMAND_CONFIG]
- [x] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` and paste diagnostic output here. (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 254 token references across 55 files.`)

**Completion Criteria**:
- [x] Docs + registry publish the new token and describe prepend vs. replace flows.
- [x] Code prepends defaults when `inheritDefaults` is omitted/true and replaces them when false, with regression tests covering both cases.
- [x] Token audit + validation logged with the latest script output.

**Priority Rationale**: P1 because automation workflows rely on defaults staying available; implementation and tokens must match the documented contract.

## P1: Startup Workspace Directories [REQ:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [IMPL:WORKSPACE_START_DIRS]

**Status**: ✅ Complete

**Description**: Allow operators to pass positional directories after CLI flags so goful opens one window per argument (ordered), expanding or shrinking the workspace before the UI loop starts while falling back to current behavior when no directories are specified.

**Dependencies**: [REQ:MODULE_VALIDATION], [ARCH:STATE_PATH_SELECTION] (workspace helper access)

**Subtasks**:
- [x] Document requirement/architecture/implementation updates + token registry entries. [REQ:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [IMPL:WORKSPACE_START_DIRS]
- [x] Develop `StartupDirParser` module with unit tests covering empty input, duplicates, missing paths, and tilde expansion. [REQ:WORKSPACE_START_DIRS] [REQ:MODULE_VALIDATION]
- [x] Develop `WorkspaceSeeder` module + tests proving windows are created/closed/reused to match the directory list with deterministic ordering. [REQ:WORKSPACE_START_DIRS] [REQ:MODULE_VALIDATION]
- [x] Integrate helpers into `main.go`, add debug logging/env flag, and ensure errors surface via `message.Errorf`. [REQ:WORKSPACE_START_DIRS]
- [x] Add integration-style coverage (e.g., filer/app tests) showing CLI args seed windows while default behavior remains unchanged. [REQ:WORKSPACE_START_DIRS]
- [x] README/ARCHITECTURE updates highlighting positional syntax and examples. [REQ:WORKSPACE_START_DIRS]
- [x] Token audit & validation (`./scripts/validate_tokens.sh`) recorded with diagnostics. [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] Requirements, architecture, implementation docs, and semantic tokens updated with cross-references.
- [x] Parser + seeder modules independently validated before integration, evidence recorded.
- [x] `main.go` honors positional directories with debug logging and safe fallbacks.
- [x] Tests (unit + integration) cover ordering, error handling, and default behavior.
- [x] README/ARCHITECTURE document the new CLI behavior.
- [x] Latest `[PROC:TOKEN_AUDIT]`/`[PROC:TOKEN_VALIDATION]` results logged.

**Priority Rationale**: P1 because deterministic workspace layouts significantly improve automation and startup ergonomics but do not block core navigation when absent.

**Validation Evidence**:
- `go test ./...` (darwin/arm64, Go 1.24.3) on 2026-01-07 covering parser/seeder helpers plus regressions in `main`.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 288 token references across 58 files.` (2026-01-07)
## P0: Cross-Platform Terminal Launcher [REQ:TERMINAL_PORTABILITY] [ARCH:TERMINAL_LAUNCHER] [IMPL:TERMINAL_ADAPTER]

**Status**: ✅ Complete

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
- [x] Manual validation checklist for macOS Terminal + Linux desktop runs [REQ:TERMINAL_PORTABILITY] — documented in `stdd/processes.md` as `[PROC:TERMINAL_VALIDATION]`; operators should run the checklist on physical macOS/Linux hardware before releases.
- [x] `[PROC:TOKEN_AUDIT]` and `./scripts/validate_tokens.sh` (`[PROC:TOKEN_VALIDATION]`) recorded after code/tests land (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 185 token references across 46 files.` on 2026-01-02)

**Completion Criteria**:
- [x] Factory + configurator modules validated independently before integration (`go test ./terminalcmd` on 2026-01-04).
- [x] `main.go` uses the new adapter; legacy gnome-terminal-only wiring is replaced by `terminalcmd.Apply`.
- [x] Tests cover selection matrix and keep-open tail behavior (`TestCommandFactory*`, `TestApply*`).
- [x] Documentation updated with macOS instructions and override guidance (README/CONTRIBUTING sections tagged with `[REQ:TERMINAL_PORTABILITY]`).
- [x] `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]` logs captured (`./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 245 token references across 52 files.` on 2026-01-04).

**Priority Rationale**: P0 because macOS users currently cannot execute external commands in a terminal, blocking a core workflow.

## P0: macOS AppleScript terminal parameters [REQ:TERMINAL_PORTABILITY] [ARCH:TERMINAL_LAUNCHER] [IMPL:TERMINAL_ADAPTER]

**Status**: ✅ Complete

**Description**: Allow the AppleScript branch to change the target application name and inline shell (e.g., iTerm2 + zsh) via runtime parameters so macOS workflows stay portable without editing Go code.

**Dependencies**: Cross-Platform Terminal Launcher, [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Update `[REQ:TERMINAL_PORTABILITY]`, `[ARCH:TERMINAL_LAUNCHER]`, and `[IMPL:TERMINAL_ADAPTER]` documentation with the new env vars and acceptance criteria.
- [x] Extend `terminalcmd.Options`/`Factory` + `main.go` wiring to honor `GOFUL_TERMINAL_APP` and `GOFUL_TERMINAL_SHELL`, including debug output and defaults.
- [x] Add unit tests covering custom app/shell values and ensure module validation evidence references `[REQ:TERMINAL_PORTABILITY]`.
- [x] Refresh README/CONTRIBUTING/process docs so operators know how to use the new knobs.
- [x] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` once the change lands and record the diagnostic output (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 260 token references across 55 files.` on 2026-01-05).

**Completion Criteria**:
- [x] Docs cross-reference new env vars and satisfaction/validation criteria.
- [x] Code + tests pass for default and customized macOS app/shell settings.
- [x] README/CONTRIBUTING mention the new env vars with `[REQ:TERMINAL_PORTABILITY]` tags.
- [x] `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]` runs captured in this task.

**Priority Rationale**: P0 follow-up to unblock macOS teams that prefer alternate terminal apps or zsh while keeping the AppleScript automation path.

## P2: CLI xform helper [REQ:CLI_TO_CHAINING] [ARCH:XFORM_CLI_PIPELINE] [IMPL:XFORM_CLI_SCRIPT]

**Status**: ✅ Complete

**Description**: Deliver a portable Bash helper that rewrites commands by inserting `--to` before every destination argument, with dry-run output for previewing quoting so external command recipes stay readable.

**Dependencies**: [REQ:MODULE_VALIDATION]

**Module Boundaries**:
- `XformArgs` — parses `-n/--dry-run`, `-h/--help`, and `--` while ensuring the minimum argument count before running transformations.
- `TargetInterleaver` — builds the transformed argv array, prints `%q`-formatted dry-run output, or executes the rewritten command and propagates its exit status.

**Subtasks**:
- [x] Document requirement/architecture/implementation updates before coding [REQ:CLI_TO_CHAINING]
- [x] Implement parser + builder modules with inline debug comments [IMPL:XFORM_CLI_SCRIPT]
- [x] Add shell-based validation suite covering dry-run output + error handling [REQ:MODULE_VALIDATION]
- [x] Run `[PROC:TOKEN_AUDIT]` and `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` `[PROC:TOKEN_VALIDATION]`, recording diagnostic output (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 269 token references across 55 files.` on 2026-01-06)

**Completion Criteria**:
- [x] Helper works when executed directly or sourced, including dry-run/help paths
- [x] Tests validate parser and builder modules independently before integration
- [x] Documentation + token registry updated with new tokens
- [x] `[PROC:TOKEN_AUDIT]` / `[PROC:TOKEN_VALIDATION]` logs stored in this entry and `implementation-decisions.md`

**Priority Rationale**: P2 because this helper improves developer ergonomics and external command authoring but does not block core app flows.

## P2: CLI xform configurability [REQ:CLI_TO_CHAINING] [ARCH:XFORM_CLI_PIPELINE] [IMPL:XFORM_CLI_SCRIPT]

**Status**: ✅ Complete

**Description**: Expand the helper so callers can choose which arguments stay untouched and which prefix is inserted (e.g., `--to`, `--dest`, etc.) while preserving the existing dry-run/exec flow.

**Dependencies**: [REQ:MODULE_VALIDATION]

**Module Boundaries**:
- `XformArgs` — extended parser that captures `--prefix` (string) and `--keep` (number of untouched arguments) along with existing dry-run/help handling.
- `TargetInterleaver` — uses the configured prefix/keep values when constructing the rewritten argv and ensures validation errors surface clearly when insufficient arguments remain.

**Subtasks**:
- [x] Update requirements/architecture/implementation docs with the new configurability expectations [REQ:CLI_TO_CHAINING]
- [x] Implement parser + builder changes supporting `--prefix`/`--keep` defaults and validation [IMPL:XFORM_CLI_SCRIPT]
- [x] Add/extend shell tests covering custom prefix/keep cases and failure paths [REQ:MODULE_VALIDATION]
- [x] Run `[PROC:TOKEN_AUDIT]` and `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` `[PROC:TOKEN_VALIDATION]` once changes land (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 269 token references across 55 files.` on 2026-01-06)

**Completion Criteria**:
- [x] CLI usage/README-like help documents the new options
- [x] Tests validate default behavior plus new prefix/keep paths independently before integration
- [x] Documentation and registries updated to describe configurability
- [x] Token audit + validation logs recorded in this entry and `implementation-decisions.md`

**Priority Rationale**: P2 ergonomics improvement—important for automation flexibility but not blocking core flows.

## P0: Event Loop Shutdown [REQ:EVENT_LOOP_SHUTDOWN] [ARCH:EVENT_LOOP_SHUTDOWN] [IMPL:EVENT_LOOP_SHUTDOWN]

**Status**: ⏳ Pending

**Description**: Stop leaking goroutines by giving `app.Goful`’s event poller an explicit shutdown path that observes application exit, drains pending events safely, and closes channels without writing after close.

**Dependencies**: [REQ:MODULE_VALIDATION], [REQ:DEBT_TRIAGE]

**Module Boundaries**:
- `PollerAdapter` – wraps `widget.PollEvent` to accept a `stop <-chan struct{}` and forward events into `g.event`.
- `ShutdownController` – coordinates stop-signal fan-out, wait groups/timeouts, and channel closure when `Run` is exiting.
- `Diagnostics` – debug logging + counters for leak detection, tied to `DEBUG: [IMPL:EVENT_LOOP_SHUTDOWN]` output.

**Subtasks**:
- [x] Publish requirement + architecture/implementation decisions and register semantic tokens [REQ:EVENT_LOOP_SHUTDOWN] [ARCH:EVENT_LOOP_SHUTDOWN] [IMPL:EVENT_LOOP_SHUTDOWN]
- [ ] Extract/implement `PollerAdapter` with stop signal handling [REQ:EVENT_LOOP_SHUTDOWN]
- [ ] Add unit tests for the adapter using fake poll sources (module validation) [REQ:MODULE_VALIDATION]
- [ ] Implement `ShutdownController` in `app/goful.go` with timeout + logging [IMPL:EVENT_LOOP_SHUTDOWN]
- [ ] Add integration tests proving `Run` stops the poller and no events fire post-shutdown [REQ:EVENT_LOOP_SHUTDOWN]
- [ ] Update `stdd/debt-log.md` D1 entry with mitigation notes once validation passes [REQ:DEBT_TRIAGE]
- [ ] Run `[PROC:TOKEN_AUDIT]` and `./scripts/validate_tokens.sh` (`[PROC:TOKEN_VALIDATION]`) after code/tests land

**Completion Criteria**:
- [ ] Poller and shutdown modules validated independently before integration (unit tests pass).
- [ ] Integration coverage ensures no writes occur after `g.event` closes and goroutine count returns to baseline.
- [ ] Debug logging documents shutdown start/stop (+ timeout) for operators.
- [ ] Debt item D1 updated with resolved status referencing this requirement.
- [ ] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` logs captured with new token references.

**Priority Rationale**: P0 because the leaking poller burns CPU and risks crashes for every exit, making the UI unreliable.

## P1: Cross-Window Difference Search [REQ:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [IMPL:DIFF_SEARCH]

**Status**: ✅ Complete

**Description**: Implement a two-command difference search that iterates through files/directories across workspace windows, finds entries that differ (missing or different size), and highlights them with cursor movement.

**Dependencies**: [REQ:MODULE_VALIDATION]

**Subtasks**:
- [x] Document requirement/architecture/implementation tokens [REQ:DIFF_SEARCH]
- [x] Implement DiffSearchState struct with initial dirs tracking [IMPL:DIFF_SEARCH]
- [x] Implement core comparison logic: union of names, alphabetic sort, difference detection [IMPL:DIFF_SEARCH]
- [x] Implement cursor movement to different file across all windows [IMPL:DIFF_SEARCH]
- [x] Implement subdirectory descent when no file differences found [IMPL:DIFF_SEARCH]
- [x] Wire StartDiffSearch and ContinueDiffSearch commands in app/goful.go [IMPL:DIFF_SEARCH]
- [x] Add keybindings (`[` start, `]` continue) and View menu entries [IMPL:DIFF_SEARCH]
- [x] Add unit tests for DiffSearchState and comparison logic [REQ:DIFF_SEARCH]
- [x] Token audit & validation [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]
- [x] Fix: Subdirectory descent respects `startAfter` position via `FindNextSubdirInAll` [IMPL:DIFF_SEARCH]

**Completion Criteria**:
- [x] Modules documented with interfaces + validation evidence
- [x] Unit tests pass independently before integration
- [x] Integration wiring + documentation merged with semantic tokens
- [x] Token audit + validation logged (`DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 669 token references across 66 files.`)

**Validation Evidence (2026-01-10)**:
- `go test ./filer/... -run "REQ_DIFF_SEARCH"` (darwin/arm64, Go 1.24.3) - 16 tests passing.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 798 token references across 69 files.`
- Message routing alignment (2026-01-10): Added `SetMessage()`/`ClearMessage()` to `diffstatus` package; status messages now route to dedicated row instead of ephemeral `message` window per `[REQ:DIFF_SEARCH]` specification.
- Bug fix (2026-01-10): Added `FindNextSubdirInAll` to respect `startAfter` during subdirectory descent, fixing search state loss when user manually navigates into subdirectories.
- Bug fix (2026-01-12): Fixed `startAfter` assignment after ascending from child directories. After `Chdir("..")`, code was using `ws.Dir().Base()` which returned the PARENT directory name instead of the CHILD we exited. This caused `FindNextSubdirInAll` to skip all siblings because the parent name wasn't in the subdirectory list. Fix: Save child name before ascending, use it for `startAfter`.
- Bug fix (2026-01-13): Fixed alphabetical comparison in `FindNextSubdirInAll`, `FindNextSubdir`, and `FindNextDifference`. Original used exact match (`name == startAfter`) which failed when `startAfter` was a filename not in the subdirs list. Changed to alphabetical comparison (`name <= startAfter`) so subdirectory search correctly finds directories after a file position (e.g., "dev" after "date.key").

**Priority Rationale**: P1 because this feature significantly improves directory comparison workflows but does not block core navigation.

## P1: Multi-Target Copy/Move via nsync [REQ:NSYNC_MULTI_TARGET] [ARCH:NSYNC_INTEGRATION] [IMPL:NSYNC_OBSERVER] [IMPL:NSYNC_COPY_MOVE]

**Status**: ⏳ Pending

**Description**: Integrate the nsync SDK from `github.com/fareedst/nsync` to provide multi-target copy/move operations where files sync to all visible workspace panes simultaneously with parallel execution and progress monitoring.

**Dependencies**: [REQ:MODULE_VALIDATION], [REQ:WINDOW_MACRO_ENUMERATION] (for `otherWindowDirPaths` helper)

**Module Boundaries**:
- `NsyncObserver` (Module 1 – `app/nsync.go`): Adapter implementing `nsync.Observer` to bridge progress events to goful's `progress` widget.
- `SyncCopy/SyncMove` (Module 2 – `app/nsync.go`): Wrapper functions configuring nsync and executing within `asyncFilectrl` pattern.
- `CopyAll/MoveAll` (Module 3 – `app/nsync.go`): Functions collecting destinations from workspace and delegating to nsync wrappers.

**Subtasks**:
- [x] Update requirements with `[REQ:NSYNC_MULTI_TARGET]` [REQ:NSYNC_MULTI_TARGET]
- [x] Update architecture decisions with `[ARCH:NSYNC_INTEGRATION]` [ARCH:NSYNC_INTEGRATION]
- [x] Update implementation decisions with `[IMPL:NSYNC_OBSERVER]` and `[IMPL:NSYNC_COPY_MOVE]` [IMPL:NSYNC_OBSERVER] [IMPL:NSYNC_COPY_MOVE]
- [x] Register tokens in `semantic-tokens.md`
- [x] Add nsync dependency from public repo `github.com/fareedst/nsync` in `go.mod`
- [ ] Implement `NsyncObserver` adapter in `app/nsync.go` [IMPL:NSYNC_OBSERVER]
- [ ] Implement `syncCopy`/`syncMove` wrappers in `app/nsync.go` [IMPL:NSYNC_COPY_MOVE]
- [x] Add `CopyAll`/`MoveAll` functions in `app/nsync.go` [IMPL:NSYNC_COPY_MOVE]
- [ ] Wire keybindings (`C`/`M`) and View menu entries in `main.go` [IMPL:NSYNC_COPY_MOVE]
- [ ] Add unit tests for observer adapter [REQ:MODULE_VALIDATION]
- [ ] Add integration tests for multi-target sync [REQ:MODULE_VALIDATION]
- [ ] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [ ] NsyncObserver module validated independently before integration
- [ ] SyncCopy/SyncMove wrappers validated with temp directory tests
- [ ] CopyAll/MoveAll modes work with 2+ panes and fall back gracefully with 1 pane
- [ ] Progress display updates during multi-file operations
- [ ] Token audit + validation logged

**Priority Rationale**: P1 because multi-target copy/move significantly improves file distribution workflows but does not block core single-target operations.

## P1: Mouse Support for File Selection [REQ:MOUSE_FILE_SELECT] [ARCH:MOUSE_EVENT_ROUTING] [IMPL:MOUSE_HIT_TEST] [IMPL:MOUSE_FILE_SELECT]

**Status**: ✅ Complete

**Description**: Enable mouse input for file selection in directory windows. Left-clicking on a file moves the cursor to that file, clicking in an unfocused window switches focus first, double-clicking enters directories, and mouse wheel scrolls the file list. This is a multi-stage effort starting with file selection and progressing through focus switching, scrolling, and eventually modal support.

**Dependencies**: [REQ:MODULE_VALIDATION]

**Module Boundaries**:
- `MouseEventTranslator` (Module 1 – `widget/widget.go`): Enables mouse at init, exports `EnableMouse`/`DisableMouse` functions. [IMPL:MOUSE_HIT_TEST]
- `HitTestFramework` (Module 2 – `widget/widget.go`, `filer/workspace.go`, `filer/directory.go`): Pure coordinate-to-widget mapping with `Contains`, `DirectoryAt`, `FileIndexAtY`. [IMPL:MOUSE_HIT_TEST]
- `MouseDispatcher` (Module 3 – `app/goful.go`): Orchestrates hit-testing and widget method calls via `mouseHandler`. [IMPL:MOUSE_FILE_SELECT]

**Completion Criteria**:
- [x] Mouse events enabled in tcell
- [x] Hit-testing modules validated independently before integration
- [x] Left-click selects files in directory windows
- [x] Clicking unfocused window switches focus
- [x] Mouse wheel scrolls file list
- [x] Token audit + validation logged

**Validation Evidence (2026-01-17)**:
- `go test ./widget/... ./filer/... -run "MOUSE"` (darwin/arm64, Go 1.24.3) - all hit-testing tests pass
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1137 token references across 74 files.`
- Implemented modules: `widget.Init()` enables mouse, `widget.Window.Contains()`, `filer.Workspace.DirectoryAt()`, `filer.Directory.FileIndexAtY()`, `app.Goful.mouseHandler()` with left-click and wheel support

**Priority Rationale**: P1 because mouse support significantly improves accessibility and user experience for GUI-oriented users, but keyboard navigation remains fully functional without it.

## P0: Linked Navigation Comparison Index Timing Fix [REQ:LINKED_NAVIGATION] [REQ:FILE_COMPARISON_COLORS] [IMPL:LINKED_NAVIGATION]

**Status**: ✅ Complete

**Description**: Fix digest comparison decoration not applying to the focused window after navigating to a subdirectory with linked navigation enabled. The comparison index was being rebuilt before the focused directory finished navigating, causing stale index entries.

**Dependencies**: None (bug fix)

**Completion Criteria**:
- [x] Root cause identified via runtime instrumentation
- [x] `ChdirAllToSubdirNoRebuild()` method added to defer index rebuild
- [x] `linkedEnterDir` sequence corrected: navigate all directories THEN rebuild index
- [x] All tests pass
- [x] Manual verification confirms 3-window digest comparison works correctly

**Validation Evidence** (2026-01-11):
- `go test ./...` passes on darwin/arm64, Go 1.24.3
- Manual verification: `=` key now calculates digests for all windows including the focused window

**Priority Rationale**: P0 because the bug broke a core feature (file digest comparison) when using linked navigation, making the comparison results incomplete and misleading.
