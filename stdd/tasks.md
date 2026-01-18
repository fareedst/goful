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

**Status**: ✅ Complete

**Description**: Stop leaking goroutines by giving `app.Goful`'s event poller an explicit shutdown path that observes application exit, drains pending events safely, and closes channels without writing after close.

**Dependencies**: [REQ:MODULE_VALIDATION], [REQ:DEBT_TRIAGE]

**Module Boundaries**:
- `PollerAdapter` – wraps `widget.PollEvent` to accept a `stop <-chan struct{}` and forward events into `g.event`.
- `ShutdownController` – coordinates stop-signal fan-out, wait groups/timeouts, and channel closure when `Run` is exiting.
- `Diagnostics` – debug logging + counters for leak detection, tied to `DEBUG: [IMPL:EVENT_LOOP_SHUTDOWN]` output.

**Subtasks**:
- [x] Publish requirement + architecture/implementation decisions and register semantic tokens [REQ:EVENT_LOOP_SHUTDOWN] [ARCH:EVENT_LOOP_SHUTDOWN] [IMPL:EVENT_LOOP_SHUTDOWN]
- [x] Extract/implement `PollerAdapter` with stop signal handling [REQ:EVENT_LOOP_SHUTDOWN]
- [x] Add unit tests for the adapter using fake poll sources (module validation) [REQ:MODULE_VALIDATION]
- [x] Implement `ShutdownController` in `app/goful.go` with timeout + logging [IMPL:EVENT_LOOP_SHUTDOWN]
- [x] Add integration tests proving `Run` stops the poller and no events fire post-shutdown [REQ:EVENT_LOOP_SHUTDOWN]
- [x] Update `stdd/debt-log.md` D1 entry with mitigation notes once validation passes [REQ:DEBT_TRIAGE]
- [x] Run `[PROC:TOKEN_AUDIT]` and `./scripts/validate_tokens.sh` (`[PROC:TOKEN_VALIDATION]`) after code/tests land

**Completion Criteria**:
- [x] Poller and shutdown modules validated independently before integration (unit tests pass).
- [x] Integration coverage ensures no writes occur after `g.event` closes and goroutine count returns to baseline.
- [x] Debug logging documents shutdown start/stop (+ timeout) for operators.
- [x] Debt item D1 updated with resolved status referencing this requirement.
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` logs captured with new token references.

**Validation Evidence (2026-01-17)**:
- `go test ./app/... -run "EVENT_LOOP_SHUTDOWN" -v` - 8 tests pass (darwin/arm64, Go 1.24.3)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1200 token references across 76 files.`
- Implementation: `pollEvents()`, `shutdownPoller()`, `debugLog()` in `app/goful.go`
- Tests: `app/shutdown_test.go` with 8 test cases covering timeout, idempotency, concurrency, channel closure

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

**Status**: ✅ Complete

**Description**: Integrate the nsync SDK from `github.com/fareedst/nsync` to provide multi-target copy/move operations where files sync to all visible workspace panes simultaneously with parallel execution and progress monitoring.

**Dependencies**: [REQ:MODULE_VALIDATION], [REQ:WINDOW_MACRO_ENUMERATION] (for `otherWindowDirPaths` helper)

**Module Boundaries**:
- `NsyncObserver` (Module 1 – `app/nsync.go`): Adapter implementing `nsync.Observer` to bridge progress events to goful's `progress` widget.
- `SyncCopy/SyncMove` (Module 2 – `app/nsync.go`): Wrapper functions configuring nsync and executing within `asyncFilectrl` pattern.
- `CopyAll/MoveAll` (Module 3 – `app/nsync.go`): Functions collecting destinations from workspace and delegating to nsync wrappers.

**Completion Criteria**:
- [x] NsyncObserver module validated independently before integration
- [x] SyncCopy/SyncMove wrappers validated with temp directory tests
- [x] CopyAll/MoveAll modes work with 2+ panes and fall back gracefully with 1 pane
- [x] Progress display updates during multi-file operations
- [x] Token audit + validation logged

**Validation Evidence (2026-01-17)**:
- `go test ./app/... -run "NSYNC"` (darwin/arm64, Go 1.24.3) - 11 tests pass
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1287 token references across 77 files.`
- Implementation: `gofulObserver`, `syncCopy`, `syncMove`, `CopyAll`, `MoveAll` in `app/nsync.go`
- Keybindings: `C` (Copy All), `M` (Move All) in filer keymap and command menu
- Confirmation modes: `copyAllMode`, `moveAllMode` in `app/mode.go`

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

## P1: Mouse Double-Click Behavior [REQ:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [IMPL:MOUSE_DOUBLE_CLICK]

**Status**: ✅ Complete

**Description**: Implement double-click detection and actions for mouse navigation. Double-clicking directories navigates into them (respecting Linked mode), double-clicking files opens them (and opens same-named files in all windows when Linked mode is enabled).

**Dependencies**: [REQ:MOUSE_FILE_SELECT] (complete), [REQ:LINKED_NAVIGATION] (complete)

**Completion Criteria**:
- [x] All subtasks complete
- [x] Double-click on directory enters it, with linked mode propagation
- [x] Double-click on file opens it, with linked mode opening all same-named files
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Validation Evidence** (2026-01-17):
- `go test ./app/... -run MOUSE_DOUBLE_CLICK` - 4 tests pass (darwin/arm64, Go 1.24.3)
- `go test ./...` - all tests pass
- `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1158 token references across 74 files.`
- Implementation in `app/goful.go`: `isDoubleClick()`, `handleDoubleClickDir()`, `handleDoubleClickFile()`
- Unit tests in `app/mouse_test.go`: `TestIsDoubleClick_REQ_MOUSE_DOUBLE_CLICK`, `TestDoubleClickThreshold_REQ_MOUSE_DOUBLE_CLICK`, `TestIsDoubleClickUpdatesState_REQ_MOUSE_DOUBLE_CLICK`, `TestDoubleClickSequence_REQ_MOUSE_DOUBLE_CLICK`
- Bug fix (2026-01-18): Linked mode file double-click now opens ALL matching files from all windows, executing the open command once for each same-named file. Root cause was that `g.Input("C-m")` only expanded `%f` macro to the focused file path.

**Priority Rationale**: P1 because double-click completes the mouse navigation experience but keyboard navigation remains fully functional without it.

---

## Phase 3: Code Quality Improvements (Identified 2026-01-17)

The following tasks were identified during an STDD documentation review to address runtime reliability issues, documented technical debt, and documentation drift.

## P0: Event Loop Shutdown [REQ:EVENT_LOOP_SHUTDOWN] [ARCH:EVENT_LOOP_SHUTDOWN] [IMPL:EVENT_LOOP_SHUTDOWN]

**Status**: ✅ Complete (see main task entry above for validation evidence)

**Priority Rationale**: P0 because the leaking poller burns CPU and risks crashes for every exit, making the UI unreliable. Completed 2026-01-17.

## P1: History Error Handling [REQ:DEBT_TRIAGE] [ARCH:DEBT_MANAGEMENT] [IMPL:HISTORY_ERROR_HANDLING]

**Status**: ✅ Complete

**Description**: Fix silent error swallowing in CLI history persistence. Distinguish first-run missing files (`os.ErrNotExist`) from actual IO failures (permissions, disk full) and surface actionable errors via `message.Error`.

**Dependencies**: None

**Debt Reference**: D2 in `stdd/debt-log.md`

**Completion Criteria**:
- [x] First-run (missing history file) works silently without errors
- [x] Real IO failures surface via `message.Error` with context
- [x] Tests cover error differentiation
- [x] Debt item D2 updated with resolved status

**Validation Evidence (2026-01-17)**:
- `go test ./cmdline/... -run "REQ_DEBT_TRIAGE"` - 6 tests pass (darwin/arm64, Go 1.24.3)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1254 token references across 77 files.`
- Implementation: `HistoryError` struct, `IsFirstRun()` helper, updated `LoadHistory`/`SaveHistory` in `cmdline/cmdline.go`
- `main.go` surfaces load errors via `message.Errorf`, save errors via `stderr` (post-TUI)
- Unit tests in `cmdline/history_test.go`: first-run success, permission errors, valid file load/save

**Priority Rationale**: P1 because silent error swallowing hides corruption and permission issues, but users can still use goful without history.

## P1: History Cache Boundaries [REQ:DEBT_TRIAGE] [ARCH:DEBT_MANAGEMENT] [IMPL:HISTORY_CACHE_LIMIT]

**Status**: ✅ Complete

**Description**: Add a configurable per-mode entry limit to `historyMap` with an eviction policy that drops the oldest entries before serialization, preventing unbounded memory growth in long-running sessions.

**Dependencies**: None

**Debt Reference**: D3 in `stdd/debt-log.md`

**Completion Criteria**:
- [x] `historyMap` never exceeds configured limit per mode
- [x] Eviction preserves most recent entries
- [x] Tests cover boundary conditions
- [x] Debt item D3 updated with resolved status

**Validation Evidence (2026-01-17)**:
- `go test ./cmdline/... -run "REQ_DEBT_TRIAGE"` - 9 tests pass (darwin/arm64, Go 1.24.3)
- Implementation: `HistoryLimit` variable (default 1000), `trimHistory()` helper in `cmdline/cmdline.go`
- Trimming applied in `History.add()` and `SaveHistory()`
- Unit tests: `TestTrimHistory_REQ_DEBT_TRIAGE`, `TestHistoryAddWithLimit_REQ_DEBT_TRIAGE`, `TestSaveHistory_TrimsBeforeSave_REQ_DEBT_TRIAGE`

**Priority Rationale**: P1 because unbounded growth affects long-running sessions, but typical sessions are short enough that this rarely impacts users.

## P1: Extmap API Safety [REQ:DEBT_TRIAGE] [ARCH:DEBT_MANAGEMENT] [IMPL:EXTMAP_API_SAFETY]

**Status**: ✅ Complete

**Description**: Fix nil map panic in `filer.AddExtmap` by allocating the inner map before writing, making the API safe for third-party integrations.

**Dependencies**: None

**Debt Reference**: D4 in `stdd/debt-log.md`

**Completion Criteria**:
- [x] `AddExtmap` works correctly when called before `MergeExtmap`
- [x] Regression test prevents future breakage
- [x] Debt item D4 updated with resolved status

**Validation Evidence (2026-01-17)**:
- `go test ./filer/... -run "REQ_DEBT_TRIAGE"` - 2 tests pass (darwin/arm64, Go 1.24.3)
- Implementation: Check + allocate inner map in `filer/filer.go` `AddExtmap()`
- Unit tests: `TestAddExtmap_NilMapSafe_REQ_DEBT_TRIAGE`, `TestAddExtmap_MultipleEntries_REQ_DEBT_TRIAGE`

**Priority Rationale**: P1 because the panic blocks third-party integrations, but the core goful app never triggers this path.

## P1: Requirements Status Synchronization [REQ:STDD_SETUP] [PROC:TOKEN_AUDIT]

**Status**: ✅ Complete

**Description**: Update `stdd/requirements.md` status flags to reflect completed tasks. Multiple requirements show "⏳ Planned" but their corresponding tasks in `tasks.md` are marked "✅ Complete", indicating documentation drift.

**Dependencies**: None

**Completion Criteria**:
- [x] All requirement statuses match task completion status
- [x] No documentation drift between requirements.md and tasks.md

**Validation Evidence (2026-01-17)**:
- Updated Requirements Registry table with 20+ implemented requirements
- Updated `[REQ:CMD_HANDLER_TESTS]` and `[REQ:INTEGRATION_FLOWS]` statuses to ✅ Implemented
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1275 token references across 77 files.`

**Priority Rationale**: P1 because documentation accuracy is essential for STDD methodology integrity and contributor onboarding.

## P1: Complete nsync Multi-Target Integration [REQ:NSYNC_MULTI_TARGET] [ARCH:NSYNC_INTEGRATION] [IMPL:NSYNC_OBSERVER] [IMPL:NSYNC_COPY_MOVE]

**Status**: ✅ Complete

**Description**: Complete the nsync SDK integration for multi-target copy/move operations. The dependency is added and `CopyAll`/`MoveAll` function stubs exist, but core implementation is incomplete.

**Dependencies**: [REQ:MODULE_VALIDATION], [REQ:WINDOW_MACRO_ENUMERATION] (for `otherWindowDirPaths` helper)

**Module Boundaries**:
- `NsyncObserver` (Module 1 – `app/nsync.go`): Adapter implementing `nsync.Observer` to bridge progress events to goful's `progress` widget.
- `SyncCopy/SyncMove` (Module 2 – `app/nsync.go`): Wrapper functions configuring nsync and executing within `asyncFilectrl` pattern.
- `CopyAll/MoveAll` (Module 3 – `app/nsync.go`): Functions collecting destinations from workspace and delegating to nsync wrappers.

**Completion Criteria**:
- [x] NsyncObserver module validated independently before integration
- [x] SyncCopy/SyncMove wrappers validated with temp directory tests
- [x] CopyAll/MoveAll modes work with 2+ panes and fall back gracefully with 1 pane
- [x] Progress display updates during multi-file operations
- [x] Token audit + validation logged

**Validation Evidence (2026-01-17)**: See main "Multi-Target Copy/Move via nsync" task for validation evidence.

**Priority Rationale**: P1 because multi-target copy/move significantly improves file distribution workflows but does not block core single-target operations.

## P1: nsync Confirmation Prompts [REQ:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [IMPL:NSYNC_CONFIRMATION]

**Status**: ✅ Complete

**Description**: Add confirmation prompts before multi-target copy/move operations using existing cmdline mode pattern.

**Dependencies**: [REQ:NSYNC_MULTI_TARGET]

**Completion Criteria**:
- [x] Confirmation prompts display source/destination counts
- [x] Y/y/Enter confirms, n/N cancels
- [x] Tests cover input handling
- [x] Token audit + validation logged

**Validation Evidence (2026-01-17)**:
- `go test ./app/... -run "NSYNC_CONFIRMATION"` (darwin/arm64, Go 1.24.3) - 4 tests pass
- Implementation: `copyAllMode`, `moveAllMode` in `app/mode.go`
- Tests: `TestCopyAllMode_String_REQ_NSYNC_CONFIRMATION`, `TestMoveAllMode_String_REQ_NSYNC_CONFIRMATION`, prompt count tests in `app/nsync_test.go`
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1287 token references across 77 files.`

**Priority Rationale**: P1 because multi-target operations are high-risk and users expect confirmation for operations affecting multiple destinations.

## P1: Mouse Cross-Window Cursor Sync [REQ:MOUSE_CROSS_WINDOW_SYNC] [ARCH:MOUSE_CROSS_WINDOW_SYNC] [IMPL:MOUSE_CROSS_WINDOW_SYNC]

**Status**: ✅ Complete

**Description**: When clicking a file in the active window, synchronize cursor to the same filename in all other windows. Focus remains on the active window.

**Dependencies**: [REQ:MOUSE_FILE_SELECT] (complete)

**Completion Criteria**:
- [x] All subtasks complete
- [x] Code implements requirement with semantic token annotations
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Validation Evidence (2026-01-17)**:
- `go test ./filer/... -run "REQ_MOUSE_CROSS_WINDOW_SYNC"` (darwin/arm64, Go 1.24.3) - 2 tests pass
- `go test ./...` - all tests pass
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1298 token references across 77 files.`
- Implementation: `handleLeftClick` in `app/goful.go` calls `ws.SetCursorByNameAll(filename)` after cursor selection
- Unit tests: `TestSetCursorByNameAll_REQ_MOUSE_CROSS_WINDOW_SYNC`, `TestSetCursorByNameAllFocusUnchanged_REQ_MOUSE_CROSS_WINDOW_SYNC` in `filer/integration_test.go`

**Priority Rationale**: P1 because this enhances file comparison workflows but does not block core navigation.

## P1: Toolbar Parent Navigation Button [REQ:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_PARENT_BUTTON]

**Status**: ✅ Complete (2026-01-18)

**Description**: Add a clickable `[^]` parent navigation button to the filer header row. The button respects Linked navigation mode: when ON, all windows navigate to parent; when OFF, only the focused window navigates. This is the first element of a planned mouse-first toolbar.

**Dependencies**: [REQ:LINKED_NAVIGATION] (complete), [REQ:MOUSE_FILE_SELECT] (complete)

**Module Boundaries**:
- `ToolbarRenderer` (Module 1 – `filer/filer.go`): Renders button and tracks screen bounds in package-level map.
- `ToolbarHitTest` (Module 2 – `filer/filer.go`): Pure coordinate-to-button-name mapping.
- `ToolbarDispatcher` (Module 3 – `app/goful.go`): Orchestrates hit-testing and action invocation.

**Subtasks**:
- [x] Add `[REQ:TOOLBAR_PARENT_BUTTON]` requirement to stdd/requirements.md [REQ:TOOLBAR_PARENT_BUTTON]
- [x] Add `[ARCH:TOOLBAR_LAYOUT]` architecture decision to stdd/architecture-decisions.md [ARCH:TOOLBAR_LAYOUT]
- [x] Add `[IMPL:TOOLBAR_PARENT_BUTTON]` to implementation-decisions.md [IMPL:TOOLBAR_PARENT_BUTTON]
- [x] Register new tokens in stdd/semantic-tokens.md
- [x] Modify `drawHeader()` in `filer/filer.go` to render `[^]` button with bounds tracking [IMPL:TOOLBAR_PARENT_BUTTON]
- [x] Add `ToolbarButtonAt()` hit-testing method to `filer/filer.go` [IMPL:TOOLBAR_PARENT_BUTTON]
- [x] Extend `handleLeftClick()` in `app/goful.go` to check toolbar hits [IMPL:TOOLBAR_PARENT_BUTTON]
- [x] Wire button click to `linkedParentNav` logic respecting Linked mode [IMPL:TOOLBAR_PARENT_BUTTON]
- [x] Add unit tests for toolbar hit-testing and navigation dispatch [REQ:TOOLBAR_PARENT_BUTTON]
- [x] Token audit & validation [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] All subtasks complete
- [x] `[^]` button appears at left edge of header row
- [x] Clicking button navigates parent (all windows when Linked, focused only when not)
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Priority Rationale**: P1 because mouse-first navigation significantly improves accessibility and user experience for GUI-oriented users, but keyboard navigation remains fully functional without it.

## P1: Toolbar Linked Mode Toggle Button [REQ:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_LINKED_TOGGLE]

**Status**: ✅ Complete (2026-01-18)

**Description**: Add a clickable `[L]` button to the toolbar in the filer header row, immediately after the `[^]` parent button. The button displays the current linked navigation mode state (reverse style when ON, normal style when OFF) and toggles the mode when clicked. This button replaces the existing conditional `[LINKED]` indicator.

**Dependencies**: [REQ:TOOLBAR_PARENT_BUTTON] (complete), [REQ:LINKED_NAVIGATION] (complete)

**Subtasks**:
- [x] Add `toolbarLinkedToggleFn` callback and setter in `filer/filer.go` [IMPL:TOOLBAR_LINKED_TOGGLE]
- [x] Extend `InvokeToolbarButton()` to handle "linked" button [IMPL:TOOLBAR_LINKED_TOGGLE]
- [x] Render `[L]` button in `drawHeader()` after `[^]` with state-based styling [IMPL:TOOLBAR_LINKED_TOGGLE]
- [x] Remove conditional `[LINKED]` indicator from `drawHeader()` [IMPL:TOOLBAR_LINKED_TOGGLE]
- [x] Wire `SetToolbarLinkedToggleFn` callback in `main.go` [IMPL:TOOLBAR_LINKED_TOGGLE]
- [x] Add unit tests for linked button hit-testing and invocation [REQ:TOOLBAR_LINKED_TOGGLE]
- [x] Token audit & validation [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] All subtasks complete
- [x] `[L]` button appears after `[^]` in header row
- [x] Button style reflects linked mode state (reverse when ON, normal when OFF)
- [x] Clicking button toggles linked mode and displays confirmation message
- [x] Existing `[LINKED]` indicator removed
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Priority Rationale**: P1 because it completes the toolbar UI pattern and provides mouse-first access to linked mode toggle, improving discoverability and accessibility.

## P1: Batch Diff Report CLI [REQ:BATCH_DIFF_REPORT] [ARCH:BATCH_DIFF_REPORT] [IMPL:BATCH_DIFF_REPORT]

**Status**: ✅ Complete (2026-01-17)

**Description**: Add a `--diff-report` CLI flag that performs a complete non-interactive directory tree comparison across 2+ directories, outputs a structured YAML report to stdout with periodic progress to stderr, then exits without launching the interactive TUI.

**Dependencies**: [REQ:DIFF_SEARCH] (complete), [REQ:MODULE_VALIDATION]

**Module Boundaries**:
- `DiffReport` (struct in `filer/diffsearch.go`): YAML-serializable report structure with directories, stats, and differences.
- `DiffEntry` (struct in `filer/diffsearch.go`): Individual difference entry with name, path, reason, isDir.
- `BatchNavigator` (struct in `filer/diffsearch.go`): Headless Navigator implementation without TUI dependencies.
- `RunBatchDiffSearch()` (function in `filer/diffsearch.go`): Main entry point that creates navigator, runs TreeWalker in collection mode, returns DiffReport.
- CLI wiring (in `main.go`): Flag parsing, validation, progress goroutine, YAML output, exit code.

**Completion Criteria**:
- [x] All subtasks complete
- [x] `goful --diff-report dir1 dir2` produces valid YAML output
- [x] Progress reporting to stderr works (and `--quiet` suppresses it)
- [x] Exit codes: 0 (no differences), 1 (error), 2 (differences found)
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Validation Evidence (2026-01-17)**:
- `go test ./filer/... -run "BATCH_DIFF_REPORT"` - 11 tests pass (darwin/arm64, Go 1.24.3)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1429 token references across 78 files.`
- Implementation: `DiffReport`, `DiffEntry`, `BatchNavigator`, `RunBatchDiffSearch` in `filer/diffsearch.go`
- CLI: `--diff-report` and `--quiet` flags, `runBatchDiffReport()` in `main.go`

**Priority Rationale**: P1 because batch directory comparison enables automation and scripting workflows but does not block interactive usage.

## P1: Toolbar Compare Digest Button [REQ:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_COMPARE_BUTTON]

**Status**: ✅ Complete

**Description**: Add a clickable `[=]` button to the toolbar that triggers digest comparison for all files appearing in multiple windows, providing single-click batch digest calculation.

**Dependencies**: [REQ:TOOLBAR_LINKED_TOGGLE] (complete), [REQ:FILE_COMPARISON_COLORS] (complete)

**Subtasks**:
- [x] Add requirement/architecture/implementation docs and register tokens
- [x] Add SharedFilenames() method to ComparisonIndex
- [x] Add toolbarCompareDigestFn callback and setter
- [x] Extend InvokeToolbarButton() for "compare" case
- [x] Render [=] button in drawHeader() after [L]
- [x] Wire callback in main.go to iterate and calculate digests
- [x] Add unit tests for hit-testing and invocation
- [x] Token audit & validation [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] All subtasks complete
- [x] [=] button appears after [L] in header row
- [x] Clicking button calculates digests for all shared files
- [x] Summary message displays file count
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] [PROC:TOKEN_AUDIT] and [PROC:TOKEN_VALIDATION] outcomes logged

**Validation Evidence (2026-01-17)**:
- `go test ./filer/... -run "REQ_TOOLBAR_COMPARE_BUTTON"` - 6 tests pass (darwin/arm64, Go 1.24.3)
- `go test ./...` - all tests pass
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1468 token references across 78 files.`
- Implementation: `SharedFilenames()` in `filer/compare.go`, `SetToolbarCompareDigestFn()` and `InvokeToolbarButton()` in `filer/filer.go`, callback wiring in `main.go`
- Unit tests: `TestToolbarCompareButtonHit_REQ_TOOLBAR_COMPARE_BUTTON`, `TestInvokeToolbarCompareButton_REQ_TOOLBAR_COMPARE_BUTTON`, `TestSharedFilenames_*_REQ_TOOLBAR_COMPARE_BUTTON` in `filer/toolbar_test.go` and `filer/compare_test.go`

**Priority Rationale**: P1 because batch digest comparison significantly improves file verification workflows but does not block core navigation.

## P1: Unify Linked Cursor Synchronization [REQ:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [IMPL:LINKED_CURSOR_SYNC]

**Status**: ✅ Complete (2026-01-18)

**Description**: Make both mouse and keyboard cursor movements respect the Linked toggle. When linked is ON, sync highlights across all windows. When OFF, only affect the focused window.

**Dependencies**: [REQ:MOUSE_CROSS_WINDOW_SYNC] (complete), [IMPL:MOUSE_CROSS_WINDOW_SYNC] (complete)

**Completion Criteria**:
- [x] All subtasks complete
- [x] Mouse and keyboard both respect Linked toggle for cursor sync
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Validation Evidence (2026-01-18)**:
- `go test ./app/... -run "REQ_LINKED_NAVIGATION"` (darwin/arm64, Go 1.24.3) - 5 tests pass
- `go test ./...` - all tests pass
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1491 token references across 78 files.`
- Implementation: `MoveCursorLinked()`, `MoveTopLinked()`, `MoveBottomLinked()`, `PageUpLinked()`, `PageDownLinked()` in `app/goful.go`
- Mouse click `SetCursorByNameAll` now conditional on `g.IsLinkedNav()` in `handleLeftClick()`
- Keymap updated in `main.go` to use linked cursor movement methods
- Unit tests: `TestMoveCursorLinked_REQ_LINKED_NAVIGATION`, `TestMoveCursorLinkedOff_REQ_LINKED_NAVIGATION`, `TestMoveTopLinked_REQ_LINKED_NAVIGATION`, `TestMoveBottomLinked_REQ_LINKED_NAVIGATION`, `TestLinkedCursorSyncMissingFile_REQ_LINKED_NAVIGATION` in `app/linked_cursor_test.go`

**Priority Rationale**: P1 because this unifies existing behavior and improves consistency between mouse and keyboard workflows.

## P1: Toolbar Sync Operation Buttons [REQ:TOOLBAR_SYNC_BUTTONS] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_SYNC_COPY] [IMPL:TOOLBAR_SYNC_DELETE] [IMPL:TOOLBAR_SYNC_RENAME] [IMPL:TOOLBAR_IGNORE_FAILURES]

**Status**: ✅ Complete

**Description**: Add four clickable toolbar buttons `[C]`, `[D]`, `[R]`, `[!]` after the `[=]` compare button. These buttons obey Linked navigation mode: when Linked is ON, they trigger Sync operations across all windows; when OFF, they trigger single-window operations. The `[!]` button toggles a persistent ignore-failures mode.

**Dependencies**: [REQ:TOOLBAR_COMPARE_BUTTON] (complete), [REQ:SYNC_COMMANDS] (complete), [REQ:LINKED_NAVIGATION] (complete)

**Subtasks**:
- [x] Register STDD tokens in semantic-tokens.md, requirements.md, architecture-decisions.md [REQ:TOOLBAR_SYNC_BUTTONS]
- [x] Create implementation decision detail files for each button [IMPL:TOOLBAR_SYNC_COPY] [IMPL:TOOLBAR_SYNC_DELETE] [IMPL:TOOLBAR_SYNC_RENAME] [IMPL:TOOLBAR_IGNORE_FAILURES]
- [x] Add `syncIgnoreFailures` state and accessors to `app/goful.go` [IMPL:TOOLBAR_IGNORE_FAILURES]
- [x] Add toolbar button callbacks and setters in `filer/filer.go` [IMPL:TOOLBAR_SYNC_COPY] [IMPL:TOOLBAR_SYNC_DELETE] [IMPL:TOOLBAR_SYNC_RENAME] [IMPL:TOOLBAR_IGNORE_FAILURES]
- [x] Extend `drawHeader()` to render `[C]`, `[D]`, `[R]`, `[!]` buttons [ARCH:TOOLBAR_LAYOUT]
- [x] Extend `InvokeToolbarButton()` with new button cases [ARCH:TOOLBAR_LAYOUT]
- [x] Wire callbacks in `main.go` with Linked mode logic [REQ:TOOLBAR_SYNC_BUTTONS]
- [x] Add unit tests for new buttons in `filer/toolbar_test.go` [REQ:TOOLBAR_SYNC_BUTTONS]
- [x] Token audit & validation [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [x] All subtasks complete
- [x] Four buttons appear after `[=]` in header row
- [x] Buttons trigger correct operations based on Linked mode
- [x] `[!]` button styling reflects ignore-failures state
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Validation Evidence (2026-01-18)**:
- `go test ./...` - all tests pass
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1621 token references across 79 files.`
- Implementation: `syncIgnoreFailures`, `IsSyncIgnoreFailures()`, `ToggleSyncIgnoreFailures()` in `app/goful.go`
- Toolbar rendering and callbacks in `filer/filer.go`: `SetToolbarSyncCopyFn()`, `SetToolbarSyncDeleteFn()`, `SetToolbarSyncRenameFn()`, `SetToolbarIgnoreFailuresFn()`, `SetSyncIgnoreFailuresIndicator()`
- Exported sync methods: `StartSyncCopy()`, `StartSyncDelete()`, `StartSyncRename()` in `app/window_wide.go`
- Unit tests: `TestToolbarSyncButtonsHit_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarSyncCopyButton_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarSyncDeleteButton_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarSyncRenameButton_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarIgnoreFailuresButton_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarSyncButtonsWithNilCallback_REQ_TOOLBAR_SYNC_BUTTONS`, `TestIgnoreFailuresIndicator_REQ_TOOLBAR_SYNC_BUTTONS` in `filer/toolbar_test.go`

**Priority Rationale**: P1 because mouse-first sync operations significantly improve accessibility and workflow efficiency, but keyboard shortcuts remain fully functional without them.

## P2: Help Popup Styling and Mouse Support [REQ:HELP_POPUP_STYLING] [ARCH:HELP_STYLING] [IMPL:HELP_STYLING]

**Status**: ✅ Complete

**Description**: Add unified color scheme to help popup (border, headers, keys, descriptions) and mouse wheel scrolling support.

**Dependencies**: [REQ:HELP_POPUP] (complete), [REQ:MOUSE_FILE_SELECT] (complete)

**Module Boundaries**:
- `HelpStyler` (Module 1 - `look/look.go`): Pure style accessors with theme-aware configuration.
- `HelpContentDrawer` (Module 2 - `help/help.go`): Content type detection and styled rendering via `widget.Drawer` interface.
- `MouseModalForwarder` (Module 3 - `app/goful.go`): Wheel event forwarding to modal widgets.

**Completion Criteria**:
- [x] All subtasks complete
- [x] Visual verification across four themes
- [x] Mouse wheel scrolls help content
- [x] Tests pass with semantic token references
- [x] Documentation updated
- [x] `[PROC:TOKEN_AUDIT]` and `[PROC:TOKEN_VALIDATION]` outcomes logged

**Validation Evidence (2026-01-18)**:
- `go build ./...` - successful
- `go test ./...` - all tests pass (darwin/arm64, Go 1.24.3)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1708 token references across 79 files.`
- Implementation in `look/look.go`: `HelpBorder()`, `HelpHeader()`, `HelpKey()`, `HelpDesc()` accessors with theme configs
- Implementation in `help/help.go`: `helpEntry` custom drawer, `drawColoredBorder()`, `drawColoredHeader()`, styled content rendering
- Implementation in `app/goful.go`: `mouseHandler()` forwards wheel events to modal widgets via `g.Next().Input()`
- README updated with Help Popup section documenting color styling and mouse scroll support

**Priority Rationale**: P2 because visual styling and mouse scroll improve user experience but keyboard navigation remains fully functional without them.

## P2: Three-Window Demo GIFs [PROC:DEMO_GENERATION] [REQ:BATCH_DIFF_REPORT] [REQ:FILE_COMPARISON_COLORS] [REQ:LINKED_NAVIGATION]

**Status**: ✅ Complete (2026-01-17)

**Description**: Create three GIF demos showcasing goful's 3-window comparison features for README documentation: batch diff report CLI, hash comparison with `=` key, and linked navigation mode.

**Dependencies**: asciinema, agg (asciinema-agg)

**Subtasks**:
- [x] Install prerequisites (asciinema 2.4.0, agg 1.6.0)
- [x] Build goful binary
- [x] Create demo directories with test files (alpha/beta/gamma with size variations)
- [x] Record and convert demo_diff_report.gif (CLI batch comparison)
- [x] Record and convert demo_compare.gif (hash comparison with = key)
- [x] Record and convert demo_linked.gif (linked navigation)
- [x] Update README.md to embed new demo GIFs
- [x] Document demo generation process in stdd/processes.md

**Completion Criteria**:
- [x] All subtasks complete
- [x] Three GIF demos created in .github/ directory
- [x] README.md updated with demo embeds
- [x] Demo generation process documented as `[PROC:DEMO_GENERATION]`

**Validation Evidence (2026-01-17)**:
- `demo_diff_report.gif` (82KB) - CLI batch diff report demo
- `demo_compare.gif` (38KB) - Interactive hash comparison demo
- `demo_linked.gif` (35KB) - Interactive linked navigation demo
- Demo directories: `/tmp/demo/{alpha,beta,gamma}` with file1, file2, file3, subdir/
- Recording tools: asciinema with TERM=xterm-256color, expect scripts for TUI automation
- Troubleshooting: Fixed "terminal not cursor addressable" panic by setting TERM in expect scripts

**Priority Rationale**: P2 because visual demos improve documentation and user onboarding but do not affect functionality.
