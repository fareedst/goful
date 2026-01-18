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
| [REQ:STDD_SETUP] | STDD methodology setup | P0 | ✅ Implemented | [ARCH:STDD_STRUCTURE] | [IMPL:STDD_FILES] |
| [REQ:MODULE_VALIDATION] | Independent module validation before integration | P0 | ✅ Implemented | [ARCH:MODULE_VALIDATION] | [IMPL:MODULE_VALIDATION] |
| [REQ:CONFIGURABLE_STATE_PATHS] | Configurable state/history persistence paths | P0 | ✅ Implemented | [ARCH:STATE_PATH_SELECTION] | [IMPL:STATE_PATH_RESOLVER] |
| [REQ:WORKSPACE_START_DIRS] | Positional CLI directories seed workspace windows | P1 | ✅ Implemented | [ARCH:WORKSPACE_BOOTSTRAP] | [IMPL:WORKSPACE_START_DIRS] |
| [REQ:EXTERNAL_COMMAND_CONFIG] | External command menu from config files | P1 | ✅ Implemented | [ARCH:EXTERNAL_COMMAND_REGISTRY] | [IMPL:EXTERNAL_COMMAND_LOADER], [IMPL:EXTERNAL_COMMAND_BINDER] |
| [REQ:FILER_EXCLUDE_NAMES] | Configurable filename exclusions | P1 | ✅ Implemented | [ARCH:FILER_EXCLUDE_FILTER] | [IMPL:FILER_EXCLUDE_RULES], [IMPL:FILER_EXCLUDE_LOADER] |
| [REQ:WINDOW_MACRO_ENUMERATION] | %D@/%d@ enumerate workspace directories | P1 | ✅ Implemented | [ARCH:WINDOW_MACRO_ENUMERATION] | [IMPL:WINDOW_MACRO_ENUMERATION] |
| [REQ:FILE_COMPARISON_COLORS] | Multi-directory file comparison coloring | P1 | ✅ Implemented | [ARCH:FILE_COMPARISON_ENGINE] | [IMPL:COMPARE_COLOR_CONFIG], [IMPL:FILE_COMPARISON_INDEX] |
| [REQ:LINKED_NAVIGATION] | Linked navigation across workspace windows | P1 | ✅ Implemented | [ARCH:LINKED_NAVIGATION] | [IMPL:LINKED_NAVIGATION] |
| [REQ:DIFF_SEARCH] | Cross-window difference search | P1 | ✅ Implemented | [ARCH:DIFF_SEARCH] | [IMPL:DIFF_SEARCH] |
| [REQ:HELP_POPUP] | Help popup displays keystroke catalog on ? key | P2 | ✅ Implemented | [ARCH:HELP_WIDGET] | [IMPL:HELP_POPUP] |
| [REQ:SYNC_COMMANDS] | Sync command operations across panes | P1 | ✅ Implemented | [ARCH:SYNC_MODE] | [IMPL:SYNC_EXECUTE] |
| [REQ:MOUSE_FILE_SELECT] | Mouse input for file selection in directory windows | P1 | ✅ Implemented | [ARCH:MOUSE_EVENT_ROUTING] | [IMPL:MOUSE_HIT_TEST], [IMPL:MOUSE_FILE_SELECT] |
| [REQ:MOUSE_DOUBLE_CLICK] | Double-click to open files and navigate directories | P1 | ✅ Implemented | [ARCH:MOUSE_DOUBLE_CLICK] | [IMPL:MOUSE_DOUBLE_CLICK] |
| [REQ:DEBT_TRIAGE] | Technical debt and risk tracking | P1 | ✅ Implemented | [ARCH:DEBT_MANAGEMENT] | [IMPL:DEBT_TRACKING] |
| [REQ:ARCH_DOCUMENTATION] | Architecture documentation | P1 | ✅ Implemented | [ARCH:DOCS_STRUCTURE] | [IMPL:DOC_ARCH_GUIDE] |
| [REQ:CONTRIBUTING_GUIDE] | Contributor standards | P1 | ✅ Implemented | [ARCH:CONTRIBUTION_PROCESS] | [IMPL:DOC_CONTRIBUTING] |
| [REQ:BEHAVIOR_BASELINE] | Baseline behavior capture | P1 | ✅ Implemented | [ARCH:BASELINE_CAPTURE] | [IMPL:BASELINE_SNAPSHOTS] |
| [REQ:EVENT_LOOP_SHUTDOWN] | Event poller shutdown control | P0 | ✅ Implemented | [ARCH:EVENT_LOOP_SHUTDOWN] | [IMPL:EVENT_LOOP_SHUTDOWN] |
| [REQ:NSYNC_MULTI_TARGET] | Multi-target copy/move via nsync SDK | P1 | ✅ Implemented | [ARCH:NSYNC_INTEGRATION] | [IMPL:NSYNC_OBSERVER], [IMPL:NSYNC_COPY_MOVE] |
| [REQ:NSYNC_CONFIRMATION] | Confirmation before multi-target copy/move | P1 | ✅ Implemented | [ARCH:NSYNC_CONFIRMATION] | [IMPL:NSYNC_CONFIRMATION] |
| [REQ:TOOLBAR_PARENT_BUTTON] | Clickable parent navigation button in toolbar | P1 | ✅ Implemented | [ARCH:TOOLBAR_LAYOUT] | [IMPL:TOOLBAR_PARENT_BUTTON] |
| [REQ:TOOLBAR_LINKED_TOGGLE] | Clickable linked mode toggle button in toolbar | P1 | ✅ Implemented | [ARCH:TOOLBAR_LAYOUT] | [IMPL:TOOLBAR_LINKED_TOGGLE] |
| [REQ:TOOLBAR_COMPARE_BUTTON] | Clickable comparison button in toolbar | P1 | ✅ Implemented | [ARCH:TOOLBAR_LAYOUT] | [IMPL:TOOLBAR_COMPARE_BUTTON] |
| [REQ:TOOLBAR_SYNC_BUTTONS] | Toolbar sync operation buttons (C/D/R/!) | P1 | ✅ Implemented | [ARCH:TOOLBAR_LAYOUT] | [IMPL:TOOLBAR_SYNC_COPY] [IMPL:TOOLBAR_SYNC_DELETE] [IMPL:TOOLBAR_SYNC_RENAME] [IMPL:TOOLBAR_IGNORE_FAILURES] |
| [REQ:HELP_POPUP_STYLING] | Help popup styling and mouse scroll support | P2 | ✅ Implemented | [ARCH:HELP_STYLING] | [IMPL:HELP_STYLING] |
| [REQ:ESCAPE_KEY_BEHAVIOR] | Escape key closes modal widgets | P0 | ✅ Implemented | [ARCH:ESCAPE_TRANSLATION] | [IMPL:ESCAPE_TRANSLATION] |
| [REQ:DOCKER_INTERACTIVE_SETUP] | Docker-based interactive Goful execution | P2 | ✅ Implemented | [ARCH:DOCKER_BUILD_STRATEGY] | [IMPL:DOCKERFILE_MULTISTAGE], [IMPL:DOCKER_COMPOSE_CONFIG] |
| [REQ:DOCKER_WINDOWS_CONTAINER] | Windows container support for Goful testing | P2 | ✅ Implemented | [ARCH:DOCKER_WINDOWS_BUILD] | [IMPL:DOCKERFILE_WINDOWS] |

### Non-Functional Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|

### Immutable Requirements (Major Version Change Required)

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|

### Incomplete Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|
| [REQ:BATCH_DIFF_REPORT] | Batch diff report CLI command | P1 | ✅ Implemented | [ARCH:BATCH_DIFF_REPORT] | [IMPL:BATCH_DIFF_REPORT] |
| [REQ:DOCKER_INTERACTIVE_SETUP] | Docker-based interactive Goful execution | P2 | ✅ Implemented | [ARCH:DOCKER_BUILD_STRATEGY] | [IMPL:DOCKERFILE_MULTISTAGE], [IMPL:DOCKER_COMPOSE_CONFIG] |
| [REQ:DOCKER_WINDOWS_CONTAINER] | Windows container support for Goful testing | P2 | ✅ Implemented | [ARCH:DOCKER_WINDOWS_BUILD] | [IMPL:DOCKERFILE_WINDOWS] |

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

### [REQ:BACKSPACE_BEHAVIOR] Backspace Navigation & Editing

**Priority: P0 (Critical)**

- **Description**: The physical Backspace/Delete key (as labeled on macOS keyboards) must consistently trigger the canonical `backspace` action across goful. In directory panes the action opens the parent directory (mirroring `C-h`), and in prompted input modes (cmdline, finder, completion, prompt dialogs) it must delete the character immediately before the cursor so users can edit text naturally.
- **Rationale**: Backspace is a core navigation/editing shortcut. When tcell delivers different key codes per terminal/OS (`KeyBackspace` vs `KeyBackspace2`), failing to normalize them breaks both workspace navigation and prompt editing, forcing users to rely on obscure alternatives.
- **Satisfaction Criteria**:
  - `widget.EventToString` maps both `tcell.KeyBackspace` and `tcell.KeyBackspace2` to the canonical `backspace` symbol so all keymaps observe the same string regardless of terminal quirks.
  - Filer keymaps keep the existing `backspace` binding that calls `Dir().Chdir("..")`, ensuring Backspace opens the parent directory from any pane.
  - Prompted input widgets (cmdline, finder, completion) keep their `backspace` handlers that invoke `DeleteBackwardChar`, guaranteeing the key erases the prior rune.
  - Documentation (`README`, baseline keymap tests) continues to list `backspace` as a required chord so regressions are caught quickly.
- **Validation Criteria**:
  - Unit tests for `widget.EventToString` cover both backspace key codes and assert the translated string equals `backspace`.
  - Existing `main_keymap_test.go` baseline coverage ensures filer/cmdline keymaps expose the `backspace` chord, acting as a regression net for bindings.
  - Manual smoke test confirms Backspace navigates upward in the filer view and deletes characters inside cmdline prompts on macOS and Linux terminals.
- **Architecture**: See `architecture-decisions.md` § Backspace Key Translation [ARCH:BACKSPACE_TRANSLATION]
- **Implementation**: See `implementation-decisions.md` § Backspace Key Translation [IMPL:BACKSPACE_TRANSLATION]

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

**Status**: ✅ Implemented

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

**Status**: ✅ Implemented

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-17)**:
- `cmdline/history_test.go` tests cover history dedup, cursor movement, and error handling
- Tests include `[REQ:CMD_HANDLER_TESTS]` token references

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-17)**:
- `filer/integration_test.go` tests cover open directory, navigate, rename, delete flows
- Tests include `[REQ:INTEGRATION_FLOWS]` token references

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-07)**:
- Task 2.1 complete with module validation evidence
- Unit tests covering flag/env/default precedence
- `./scripts/validate_tokens.sh` → verified token references

### [REQ:WORKSPACE_START_DIRS] Positional Startup Directories

**Priority: P1 (Important)**

- **Description**: After CLI flags are parsed, any trailing positional arguments must be treated as explicit workspace directory targets. Goful should open one filer window per argument—creating, reusing, or closing panes so the visible workspace order exactly matches the provided list—and fall back to the current directory/open-directory behavior when no arguments are supplied.
- **Rationale**: Power users often launch goful from shell scripts or project-specific wrappers and want deterministic, multi-pane layouts without manual navigation. Enabling positional directories eliminates repetitive `b`/`C-m` workflows, improves automation ergonomics, and keeps invocation parity with other TUI file managers that accept startup paths.
- **Satisfaction Criteria**:
  - `flag.Args()` are interpreted as startup directories in the order supplied; when empty, the existing default workspace arrangement is preserved.
  - The workspace is expanded or shrunk to match the number of directories provided, creating new panes or closing extras so every argument maps to exactly one window with matching focus order.
  - Nonexistent directories trigger `message.Errorf` guidance yet do not crash the UI; remaining valid directories continue to open.
  - `GOFUL_DEBUG_WORKSPACE=1` emits diagnostics referencing `[IMPL:WORKSPACE_START_DIRS]` that show the parsed arguments, normalization, and workspace mutations.
  - README/ARCHITECTURE documentation lists the positional syntax with examples and describes how defaults behave when arguments are omitted.
- **Validation Criteria**:
  - Unit tests cover the startup parser (empty args, multiple entries, whitespace, non-existent paths) and the workspace seeding helper (creation, reuse, trimming, ordering).
  - Integration-style tests prove launching with `go run . dirA dirB dirC` yields three windows in order, while launching without arguments preserves historical behavior.
  - Token validation confirms `[REQ:WORKSPACE_START_DIRS]`, `[ARCH:WORKSPACE_BOOTSTRAP]`, and `[IMPL:WORKSPACE_START_DIRS]` are referenced across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Workspace Bootstrap from Positional Directories [ARCH:WORKSPACE_BOOTSTRAP]
- **Implementation**: See `implementation-decisions.md` § Startup Directory Parser & Seeder [IMPL:WORKSPACE_START_DIRS]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-07)**:
- Parser + seeder modules validated independently
- `go test ./...` (darwin/arm64, Go 1.24.3)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 288 token references across 58 files.`

### [REQ:WORKSPACE_START_DIRS] Positional Startup Directories (Duplicate Entry)

**Priority: P1 (Important)**

- **Description**: After CLI flags are parsed, any trailing positional arguments must be interpreted as explicit workspace directory targets. Goful should open one filer window per argument—creating, reusing, or closing panes so the visible workspace order exactly matches the provided list—and fall back to the current directory restoration behavior when no arguments are supplied.
- **Rationale**: Operators often launch goful via shell scripts or project-specific wrappers and want deterministic multi-pane layouts without manual navigation. Allowing positional directories removes repetitive keybind sequences, enables automation to align UI state with the calling context, and keeps parity with other TUIs that accept startup paths.
- **Satisfaction Criteria**:
  - `flag.Args()` are consumed in order and mapped to workspace directories; empty input leaves startup behavior unchanged.
  - The workspace is expanded or shrunk to match the number of directories provided, with each pane displaying the corresponding path and the first argument receiving focus.
  - Non-existent paths trigger actionable `message.Errorf` guidance yet do not crash the UI; the remaining valid directories continue to open.
  - `GOFUL_DEBUG_WORKSPACE=1` emits `DEBUG: [IMPL:WORKSPACE_START_DIRS] ...` diagnostics that show parsed arguments, normalization decisions, and workspace mutations.
  - README/ARCHITECTURE documentation includes examples illustrating positional usage, duplicate handling, and the debug workflow.
- **Validation Criteria**:
  - Unit tests cover the parser (empty args, tilde expansion, duplicates, invalid paths) and the workspace-seeding helper (window creation/removal, ordering, focus).
  - Integration-style tests demonstrate that launching with positional directories produces deterministic panes, while launching without extra arguments preserves historical behavior.
  - Token validation confirms `[REQ:WORKSPACE_START_DIRS]`, `[ARCH:WORKSPACE_BOOTSTRAP]`, and `[IMPL:WORKSPACE_START_DIRS]` references across documents, code comments, and tests.
- **Architecture**: See `architecture-decisions.md` § Workspace Bootstrap from Positional Directories [ARCH:WORKSPACE_BOOTSTRAP]
- **Implementation**: See `implementation-decisions.md` § Startup Directory Parser & Seeder [IMPL:WORKSPACE_START_DIRS]

**Status**: ✅ Implemented

### [REQ:WINDOW_MACRO_ENUMERATION] External Command Window Enumeration

**Priority: P1 (Important)**

- **Description**: External command macros must expose the full set of visible directories so operators can pass every window to shell scripts. `%D@` appends the other window paths (relative order starting from the next window and wrapping) with shell quoting so each entry stays safe. `%~D@` keeps the historical "non-quote" modifier and therefore emits the exact paths without shell escaping so advanced workflows can opt into raw arguments. `%d@` mirrors the same enumeration but emits just the directory names (no parent directories) so scripts that only care about window labels can avoid extra `basename` calls, and `%~d@` keeps the raw/non-quote semantics for that list as well. `echo %D %D@` therefore prints all window directories with the focused window first, while `echo %d@` mirrors the same ordering using only directory names.
- **Rationale**: Bulk copy/move workflows depend on knowing all workspace paths. Today `%D2` exposes only the next window, so automation that needs >2 windows requires manual re-entry. Enumerating the remaining windows keeps macros self-contained and removes repetitive typing while preserving the tilde modifier's "raw" semantics for compatibility with other macros and reducing extra shell processing when only directory names are needed.
- **Satisfaction Criteria**:
  - `%D@` expands to a space-separated list of every other directory path in deterministic order (start with next window, then wrap through the rest). When only one window is open, the expansion is empty and nothing extra is appended.
  - `%~D@` emits the same ordering as `%D@` but leaves each entry unquoted (no escaping) so scripts that intentionally need raw arguments can opt in by using the tilde modifier.
  - `%d@` expands to the same deterministic ordering but returns just the directory names (`Directory.Base()`), quoting each entry unless the caller supplies the tilde modifier.
  - `%~d@` mirrors `%d@` but returns raw names (no quoting) so existing `%~` automation patterns stay consistent.
  - `%D@`/`%d@` respect the same macro parser features as their single-item counterparts (supports escaping, `%~~` safeguards, etc.) and can be combined with other text inside commands.
  - README and developer docs document all four placeholders, clearly stating the quoting and basename differences so operators choose the right macro for their workflow.
  - Regression tests prove `expandMacro("echo %D %D@")` covers all window paths with the focused directory appearing only once at the beginning, that `expandMacro("echo %d@")` emits only the directory names, and that `%~D@`/`%~d@` return raw entries even when directories include spaces.
- **Validation Criteria**:
  - Pure helper tests validate that the path and name enumeration logic handles 1–4 windows, wrapping order, and both quoted and raw formatting modes.
  - `app/spawn_test.go` exercises `%D@`, `%~D@`, `%d@`, and `%~d@` end-to-end, including scenarios with spaces that prove the quoting guarantees hold.
  - Token validation confirms `[REQ:WINDOW_MACRO_ENUMERATION]`, `[ARCH:WINDOW_MACRO_ENUMERATION]`, and `[IMPL:WINDOW_MACRO_ENUMERATION]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Window Macro Enumeration [ARCH:WINDOW_MACRO_ENUMERATION]
- **Implementation**: See `implementation-decisions.md` § Window Macro Enumeration [IMPL:WINDOW_MACRO_ENUMERATION]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-05)**:
- `%D@`, `%~D@`, `%d@`, `%~d@` macros implemented with deterministic ordering
- Unit + integration tests document module validation evidence
- `./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 260 token references across 55 files.`

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-02)**:
- JSON/YAML loader with prepend-by-default semantics
- Module validation for loader + binder
- `./scripts/validate_tokens.sh` → verified token references

### [REQ:FILER_EXCLUDE_NAMES] Configurable Filename Exclusions

**Priority: P1 (Important)**

- **Description**: Goful must honor a user-supplied block list of file *basenames* (stem + extension after the last path separator) so noise files such as `.DS_Store`, `Thumbs.db`, or build artifacts never appear in directory listings, marks, or finder results while the filter is enabled. The exclude list is managed in a newline-delimited text file located by flag/env/default (`-exclude-names`, `GOFUL_EXCLUDES_FILE`, default `~/.goful/excludes`), supports `#` comments and blank lines, and is case-insensitive per entry. The filter is enabled automatically when at least one name is configured and can be toggled on/off at runtime via a dedicated keystroke (and View menu entry) without restarting the UI.
- **Rationale**: Large repositories and shared workspaces often contain generated files that clutter navigation and lead to accidental operations. Centralizing the filter in configuration keeps behaviour consistent across panes, reduces cognitive load, and allows teams to share curated lists without patching the codebase.
- **Satisfaction Criteria**:
  - CLI flag `-exclude-names`, environment variable `GOFUL_EXCLUDES_FILE`, and default path precedence follow the existing resolver contract so operators can override the block list location without code changes.
  - Each non-empty, non-comment line represents a basename to hide regardless of directory depth; matching entries are omitted from directory listings, finder results, macro expansions, and downstream operations while the filter is active.
  - Filtering is case-insensitive to cover typical cross-platform nuisance files, and diagnostics reference `[IMPL:FILER_EXCLUDE_RULES]` when entries are loaded or skipped.
  - A dedicated keystroke (surfaced both as a View menu item and a direct key binding) toggles the filter state at runtime, emits a `message.Infof` summary with `[REQ:FILER_EXCLUDE_NAMES]`, and forces `Workspace.ReloadAll()` so panes immediately reflect the change.
  - When no list is configured, toggle attempts log actionable guidance instead of silently failing, and the UI behaves exactly as today (no entries hidden).
  - README/ARCHITECTURE documentation explains the file format (comments, whitespace), precedence, case-insensitive matching, default path, and the toggle workflow.
- **Validation Criteria**:
  - Unit tests cover the parser (comments, whitespace, duplicates, case normalization), runtime toggle helpers, and the directory-level filter to prove hidden files never reach the list when enabled.
  - Integration-style tests (or augmented filer integration tests) demonstrate that excluded basenames disappear from directory listings and reappear when the filter is toggled off.
  - Token validation confirms `[REQ:FILER_EXCLUDE_NAMES]`, `[ARCH:FILER_EXCLUDE_FILTER]`, `[IMPL:FILER_EXCLUDE_RULES]`, and `[IMPL:FILER_EXCLUDE_LOADER]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Filename Exclude Filter [ARCH:FILER_EXCLUDE_FILTER]
- **Implementation**: See `implementation-decisions.md` § Filename Exclude Rules / Loader [IMPL:FILER_EXCLUDE_RULES], [IMPL:FILER_EXCLUDE_LOADER]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-07)**:
- `go test ./...` (darwin/arm64, Go 1.24.3)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 302 token references across 58 files.`

### [REQ:TERMINAL_PORTABILITY] Cross-Platform Terminal Launcher

**Priority: P0 (Critical)**

- **Description**: Goful must launch foreground commands in an OS-appropriate terminal so macOS, Linux desktops, and tmux users all see a usable window without editing Go code. The launcher must expose an override hook (`GOFUL_TERMINAL_CMD`) and retain the historical pause tail so commands remain visible until acknowledged.
- **Rationale**: Recent dependency upgrades broke macOS Terminal invocation entirely, and Linux users increasingly prefer alternate emulators. Centralizing selection logic keeps behaviour testable and prevents regressions when upstream terminals change.
- **Satisfaction Criteria**:
  - `tmux`/`screen` sessions always use `tmux new-window -n <cmd>` regardless of OS.
  - macOS launches via AppleScript, reusing the historical payload plus pause tail while injecting the focused directory ahead of the command.
  - macOS operators can change the AppleScript application name (`GOFUL_TERMINAL_APP`, default `Terminal`) and the shell used inside that window (`GOFUL_TERMINAL_SHELL`, default `bash`) at runtime so iTerm2 or zsh workflows work without modifying Go code.
  - Linux desktops default to gnome-terminal with the legacy title escape, and overrides (e.g., `GOFUL_TERMINAL_CMD="alacritty -e"`) insert before the `bash -c` payload.
  - `GOFUL_DEBUG_TERMINAL=1` emits `DEBUG: [IMPL:TERMINAL_ADAPTER]` logs describing the branch taken.
  - README and CONTRIBUTING include guidance for macOS behaviour, overrides, and troubleshooting, all tagged with `[REQ:TERMINAL_PORTABILITY]`.
- **Validation Criteria**:
  - Unit tests cover override parsing plus Linux, macOS, and tmux branches with `[REQ:TERMINAL_PORTABILITY]` suffixes.
  - Tests exercise the AppleScript branch with custom `GOFUL_TERMINAL_APP` / `GOFUL_TERMINAL_SHELL` values to prove the new runtime parameters flow into the generated commands.
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

### [REQ:CLI_TO_CHAINING] Command-Line Target Chaining Helper

**Priority: P2 (Nice-to-have)**

- **Description**: Provide a portable Bash helper (`scripts/xform.sh`) that rewrites multi-target commands by inserting a caller-provided prefix (default `--to`) in front of every argument beyond a configurable “keep” window (default 2). The helper must retain its preview mode, continue working when executed directly or sourced, and remain compatible with macOS `/bin/bash` 3.2+.
- **Rationale**: Copy/move automation frequently needs to call tools that expect repeating flags such as `--to <path>` or `--dest <path>`. Manually inserting these strings between each argument is error-prone, especially when dealing with paths that include spaces. A tiny CLI helper keeps command definitions declarative, reduces quoting mistakes, and becomes reusable across workflows because callers can tune both the prefix and how many leading arguments stay untouched.
- **Satisfaction Criteria**:
  - `scripts/xform.sh` accepts at least `keep + 1` positional arguments and rewrites every argument past the keep index into `<prefix> <target>` pairs while executing the resulting command.
  - CLI options include `-p/--prefix <string>` (defaults to `--to`), `-k/--keep <n>` (defaults to 2 and must be ≥1), and the existing `-n/--dry-run` / `-h/--help` flags. Invalid combinations emit actionable errors and exit with code 64.
  - Dry-run output prints the fully quoted command (using `%q` formatting) without executing it so operators can verify quoting in CI or before destructive runs.
  - The helper supports both direct execution (shebang) and sourcing (defines `xform` function) on macOS `/bin/bash` 3.2+ without using Bash 4+ only features.
- **Validation Criteria**:
  - Shell-based tests exercise the parser (including custom prefix/keep and error handling) and the command construction module independently, proving dry-run output matches expectations.
  - Integration tests (or scripted invocations) ensure the helper can be executed directly, supports dry-run output, and propagates exit codes from the invoked command.
  - Token validation confirms `[REQ:CLI_TO_CHAINING]`, `[ARCH:XFORM_CLI_PIPELINE]`, and `[IMPL:XFORM_CLI_SCRIPT]` references exist across documentation, the script, and its tests.
- **Architecture**: See `architecture-decisions.md` § Xform CLI Pipeline [ARCH:XFORM_CLI_PIPELINE]
- **Implementation**: See `implementation-decisions.md` § Xform CLI Script [IMPL:XFORM_CLI_SCRIPT]

**Status**: ✅ Implemented

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-01)**:
- `ARCHITECTURE.md` published and linked from README/CONTRIBUTING
- Package/data flow documented with semantic token cross-references

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-01)**:
- `CONTRIBUTING.md` published with workflow checklist and debug/logging policy
- Linked from README

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-01)**:
- `main_keymap_test.go` (`KeymapBaselineSuite`) enumerates canonical filer/cmdline/finder/completion/menu chords
- Tests include `[TEST:KEYMAP_BASELINE]` token references

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-01)**:
- `stdd/debt-log.md` created with D1-D4 entries
- TODO annotations added to hotspot files with `[IMPL:DEBT_TRACKING]` tokens

### [REQ:FILE_COMPARISON_COLORS] Multi-Directory File Comparison Coloring

**Priority: P1 (Important)**

- **Description**: When multiple directories are displayed in a workspace, files with common names across windows shall be color-coded. The name, size, and modification time are independently color-coded based on comparison results. File name comparison is case-sensitive, and timestamp comparison is precise to the second. The color scheme is configurable via a global YAML file. The feature is **enabled by default** and can be toggled on/off via keystroke (`` ` `` backtick). Additionally, users can trigger on-demand digest calculation (`=` key) for files with equal sizes to verify content identity; results are displayed with terminal attributes (underline for equal digests, strikethrough for different).
- **Rationale**: Users often compare directories to identify matching, newer, or larger files. Visual color-coding enables instant recognition of file relationships across panes without manual comparison, improving workflow efficiency for file synchronization, backup verification, and duplicate detection tasks. Digest comparison provides definitive content verification for files with matching sizes.
- **Satisfaction Criteria**:
  - Comparison coloring is enabled by default on startup.
  - File names appearing in multiple directories are highlighted; files unique to one directory use neutral colors.
  - For files with matching names across directories: sizes are color-coded as equal/smallest/largest; times are color-coded as equal/earliest/latest.
  - Color scheme is configurable via YAML file at flag/env/default path (`-compare-colors`, `GOFUL_COMPARE_COLORS`, default `~/.goful/compare_colors.yaml`).
  - A dedicated keystroke (`` ` ``) and View menu entry toggle comparison coloring on/off at runtime.
  - Progressive rendering: files display immediately with standard colors, then comparison colors apply once all directories finish loading.
  - Single-directory workspaces or disabled mode show standard file-type colors only.
  - Pressing `=` on a file calculates xxHash64 digests for same-named files with equal sizes across directories.
  - Files with equal digests display size with underline attribute; files with different digests display size with strikethrough attribute.
- **Validation Criteria**:
  - Unit tests cover the comparison index builder for various window counts and file combinations.
  - Unit tests cover YAML config loading with defaults for missing/invalid files.
  - Integration tests prove color-coding applies correctly after directory loading completes.
  - Toggle behavior verified via runtime state checks.
  - Unit tests cover digest calculation and state propagation for equal/different content.
  - Token validation confirms `[REQ:FILE_COMPARISON_COLORS]`, `[ARCH:FILE_COMPARISON_ENGINE]`, and `[IMPL:*]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § File Comparison Engine [ARCH:FILE_COMPARISON_ENGINE]
- **Implementation**: See `implementation-decisions.md` § Comparison Color Config [IMPL:COMPARE_COLOR_CONFIG], File Comparison Index [IMPL:FILE_COMPARISON_INDEX], Comparison Draw [IMPL:COMPARISON_DRAW], Digest Comparison [IMPL:DIGEST_COMPARISON]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-09)**:
- Comparison coloring enabled by default with YAML config
- Digest calculation via `=` key with xxHash64
- Unit + integration tests with token references

### [REQ:LINKED_NAVIGATION] Linked Navigation Mode

**Priority: P1 (Important)**

- **Description**: Goful must support a toggleable "linked" navigation mode where directory navigation in the focused window propagates to all other directory windows in the current workspace. When enabled and the user navigates into a subdirectory, all other windows that contain a matching subdirectory also navigate to it. When the user presses backspace (parent directory), all windows navigate to their respective parent directories. When linked mode is ON, cursor movements (via mouse clicks or keyboard navigation) synchronize the highlight position across all windows by filename. When linked mode is OFF, cursor movements only affect the focused window. The mode is on by default.
- **Rationale**: Operators comparing similar directory structures (e.g., syncing folder hierarchies, comparing release versions) benefit from synchronized navigation across panes. Manual navigation in each pane is tedious and error-prone when directory structures mirror each other. Synchronized cursor highlighting provides visual feedback showing matching files across all windows.
- **Satisfaction Criteria**:
  - A toggle mechanism (`L` uppercase or `M-l` Alt+l) enables/disables linked navigation mode.
  - When disabled, navigation in the focused window affects only that window (historical behavior).
  - When enabled, entering a subdirectory attempts to navigate all other workspace windows to a matching subdirectory (by name) if it exists in each window's current path.
  - When enabled, pressing backspace (parent navigation) causes all windows to navigate to their respective parent directories.
  - When enabled, changing sort order applies the same sort type to all windows in the workspace.
  - When enabled, mouse single-click file selection syncs cursor to same filename in all windows.
  - When enabled, keyboard cursor movement (up/down/page/home/end) syncs cursor to same filename in all windows.
  - When disabled, both mouse and keyboard cursor movements only affect the focused window.
  - A visual indicator (`[L]` button in toolbar, reverse style when ON) shows the current mode state.
  - The mode state is per-session and does not persist across restarts.
  - When enabled and subdirectory navigation cannot complete in one or more windows (subdirectory missing) but succeeds in at least one window, linked navigation is automatically disabled with a message informing the user of the divergent directory structures.
- **Validation Criteria**:
  - Unit tests cover the workspace navigation helpers (`ChdirAllToSubdir`, `ChdirAllToParent`) independently with `[REQ:LINKED_NAVIGATION]` references.
  - Integration tests prove the toggle state affects navigation behavior correctly.
  - Integration tests verify auto-disable behavior when subdirectory navigation partially fails.
  - Manual verification confirms the header indicator displays correctly and navigation propagates as expected.
- **Architecture**: See `architecture-decisions.md` § Linked Navigation Mode [ARCH:LINKED_NAVIGATION]
- **Implementation**: See `implementation-decisions.md` § Linked Navigation Implementation [IMPL:LINKED_NAVIGATION]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-18)**:
- `TestChdirAllToSubdir_REQ_LINKED_NAVIGATION`, `TestChdirAllToParent_REQ_LINKED_NAVIGATION`, `TestLinkedNavigationSingleWindow_REQ_LINKED_NAVIGATION`, `TestSortAllBy_REQ_LINKED_NAVIGATION` in `filer/integration_test.go` covering the workspace navigation and sort helpers.
- `TestMoveCursorLinked_REQ_LINKED_NAVIGATION`, `TestMoveCursorLinkedOff_REQ_LINKED_NAVIGATION`, `TestMoveTopLinked_REQ_LINKED_NAVIGATION`, `TestMoveBottomLinked_REQ_LINKED_NAVIGATION`, `TestLinkedCursorSyncMissingFile_REQ_LINKED_NAVIGATION` in `app/linked_cursor_test.go` covering cursor synchronization.
- Toggle keystroke: `L` (uppercase, works on all platforms) or `M-l` (Alt+l, may not work on macOS where Option produces special characters), or click `[L]` toolbar button.
- Header toolbar button `[L]` displayed with reverse style when mode is active.
- Cursor sync: mouse clicks and keyboard cursor movements sync across windows when linked mode is enabled.
- Sort synchronization: all sort menu options apply to all windows when linked mode is enabled.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1491 token references across 78 files.`

### [REQ:DIFF_SEARCH] Cross-Window Difference Search

**Priority: P1 (Important)**

- **Description**: Goful must provide two commands for iteratively searching for file/directory differences across all windows in the current workspace. Command 1 (Start) records the "initial directories" for each window, then iterates through entries in alphabetic order (case-sensitive) with **files first, then directories** at each level. Entries are "different" if missing from any window OR have different sizes across windows. When a difference is found, the search stops, cursors move to that entry in all windows (where present), and a message explains why the search stopped. The difference name and type are recorded in state. If no file differences exist at a level, subdirectories at that level are checked. If a subdirectory exists in all windows, the search descends into it recursively. Command 2 (Continue) uses the recorded difference name from state (not cursor position) to skip the last found difference and continue from the next entry. This ensures correct continuation even when the cursor cannot be set (e.g., when a subdirectory is missing in some windows).
- **Rationale**: Users comparing similar directory structures need an efficient way to find differences without manually inspecting each file. This feature provides guided navigation through differences, allowing users to quickly identify mismatches, missing files, or size discrepancies across workspace windows.
- **Satisfaction Criteria**:
  - Command 1 records initial directory paths for all windows and begins comparison from the first alphabetic entry.
  - **Files first, then directories**: At each level, ALL files are processed in alphabetic order before ANY directories are processed. This ensures a consistent, predictable search order.
  - Files are compared across all windows: "different" means missing from any window OR having different sizes in any window.
  - When a difference is found, cursors move to that entry in all windows where it exists, and a descriptive message appears (e.g., "Different: file.txt missing in window 2" or "Different: file.txt size mismatch"). The difference name (with "/" suffix for directories) and reason are recorded in state.
  - If no file differences are found at a level, subdirectories at that level are processed in alphabetic order; missing subdirs are treated as differences.
  - Descending into a subdir that exists in all windows continues the comparison recursively at the next level.
  - When ascending from a subdirectory, if at root level and there are no more subdirectories after the one ascended from, the search is complete.
  - If all entries are processed back to initial directories, a "Difference search complete - all differences found" message appears.
  - Command 2 uses the recorded difference name from state (with "/" suffix removed for directories) to skip the last found difference and continue from the next entry. This ensures correct continuation even when the cursor cannot be set (e.g., when a subdirectory is missing in some windows).
  - The search state (initial directories, active flag) persists across commands until a new search is started.
  - A dedicated persistent status line appears when diff search is active, showing the current search status (path being searched, files checked, or last difference found).
  - The status line persists until the search is completed or cleared, and does not auto-dismiss like regular messages.
  - The UI updates once per second during active search to show progress.
- **Validation Criteria**:
  - Unit tests cover the difference detection logic (missing files, size mismatches, all-same scenarios).
  - Unit tests cover the alphabetic iteration and subdirectory descent logic.
  - Integration tests prove cursor movement across windows when a difference is found.
  - Manual verification confirms the dedicated status line displays and persists during search.
  - Token validation confirms `[REQ:DIFF_SEARCH]`, `[ARCH:DIFF_SEARCH]`, and `[IMPL:DIFF_SEARCH]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Difference Search Engine [ARCH:DIFF_SEARCH]
- **Implementation**: See `implementation-decisions.md` § Difference Search Implementation [IMPL:DIFF_SEARCH]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-10)**:
- Unit tests in `filer/diffsearch_test.go` covering state management, difference detection, and alphabetic sorting.
- Keybindings: `[` (start diff search), `]` (continue diff search), also in View menu.
- Dedicated `diffstatus` package provides persistent status line during search.
- Periodic UI refresh (1 second ticker) updates status during active search.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 756 token references across 68 files.`

### [REQ:NSYNC_MULTI_TARGET] Multi-Target Copy/Move via nsync SDK

**Priority: P1 (Important)**

- **Description**: Goful must provide a multi-target copy/move capability that leverages the nsync SDK to sync files to all visible workspace panes simultaneously. When the user invokes a "Copy All" or "Move All" command, the selected files (marked or cursor) are copied/moved in parallel to all other directories visible in the workspace. The existing single-target copy/move commands remain unchanged for single-destination operations.
- **Rationale**: Users frequently need to distribute files across multiple directories at once (e.g., deploying to multiple servers, creating backups across drives, distributing assets to multiple project folders). The builtin single-target copy forces repetitive operations. The nsync SDK provides parallel multi-destination sync with progress monitoring, content verification, and move semantics—features that would be complex to implement from scratch.
- **Satisfaction Criteria**:
  - `CopyAll` command copies selected files (marked or cursor) to all other visible workspace directories in parallel using nsync.
  - `MoveAll` command moves selected files to all other visible workspace directories, deleting the source only after successful sync to all destinations.
  - Progress display integrates with goful's existing `progress` widget showing total bytes, per-item updates, and throughput.
  - When only one pane is visible, these commands fall back to the single-target builtin operation with an informative message.
  - Keybindings (`C` for Copy All, `M` for Move All) and command menu entries provide discoverability.
  - Context cancellation (Ctrl-C during operation) stops the nsync operation gracefully.
  - Errors for specific destinations are reported without aborting the entire operation.
- **Validation Criteria**:
  - Unit tests cover the nsync observer adapter with mocked progress calls.
  - Integration tests verify multi-destination sync to temp directories with file verification.
  - Tests verify fallback to builtin when only one pane exists.
  - Manual verification confirms progress display updates during multi-file operations.
  - Token validation confirms `[REQ:NSYNC_MULTI_TARGET]`, `[ARCH:NSYNC_INTEGRATION]`, and `[IMPL:*]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § nsync Integration [ARCH:NSYNC_INTEGRATION]
- **Implementation**: See `implementation-decisions.md` § nsync Observer [IMPL:NSYNC_OBSERVER], nsync Copy/Move [IMPL:NSYNC_COPY_MOVE]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-17)**:
- `go test ./app/... -run "NSYNC"` (darwin/arm64, Go 1.24.3) - 11 tests pass
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1287 token references across 77 files.`
- Implementation: `gofulObserver`, `syncCopy`, `syncMove`, `CopyAll`, `MoveAll` in `app/nsync.go`
- Keybindings: `C` (Copy All), `M` (Move All) in filer keymap and command menu

### [REQ:NSYNC_CONFIRMATION] Multi-Target Copy/Move Confirmation

**Priority: P1 (Important)**

- **Description**: Before executing a multi-target copy (`CopyAll`) or move (`MoveAll`) operation, goful must display a confirmation prompt showing the number of source files and destination directories, and wait for explicit user confirmation before proceeding.
- **Rationale**: Multi-target operations affect multiple directories simultaneously and are not easily reversible. Users need to verify their intent before files are copied or moved to all visible panes, especially for move operations where source files are deleted after sync.
- **Satisfaction Criteria**:
  - `CopyAll` displays a prompt like `Copy N file(s) to M destinations? [Y/n]` before executing.
  - `MoveAll` displays a prompt like `Move N file(s) to M destinations? [Y/n]` before executing.
  - Pressing `Y`, `y`, or Enter (empty input) confirms and proceeds with the operation.
  - Pressing `n` or `N` cancels the operation and returns to normal mode.
  - Any other input clears the text field and awaits valid input.
  - The confirmation prompt follows the existing cmdline mode pattern used by `removeMode` and `quitMode`.
- **Validation Criteria**:
  - Unit tests verify the confirmation mode accepts valid inputs and rejects invalid ones.
  - Integration tests verify the operation only executes after confirmation.
  - Manual verification confirms the prompt displays correct source/destination counts.
  - Token validation confirms `[REQ:NSYNC_CONFIRMATION]`, `[ARCH:NSYNC_CONFIRMATION]`, and `[IMPL:NSYNC_CONFIRMATION]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § nsync Confirmation [ARCH:NSYNC_CONFIRMATION]
- **Implementation**: See `implementation-decisions.md` § nsync Confirmation Modes [IMPL:NSYNC_CONFIRMATION]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-17)**:
- `go test ./app/... -run "NSYNC_CONFIRMATION"` (darwin/arm64, Go 1.24.3) - 4 tests pass
- Implementation: `copyAllMode`, `moveAllMode` in `app/mode.go`
- Tests: `TestCopyAllMode_String_REQ_NSYNC_CONFIRMATION`, `TestMoveAllMode_String_REQ_NSYNC_CONFIRMATION`, prompt count tests in `app/nsync_test.go`

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

**Status**: ✅ Implemented

**Validation Evidence (2026-01-10)**:
- Keybindings: `[` (start), `]` (continue)
- Dedicated `diffstatus` package for persistent status line
- 16 tests covering difference detection and traversal

### [REQ:HELP_POPUP] Help Popup Keystroke Catalog

**Priority: P2 (Nice-to-have)**

- **Description**: Goful must provide a Help popup that displays the full keystroke catalog when the user presses `?`. The popup is a scrollable list showing all available key bindings organized by function. Pressing `?` again (toggle), `q`, `C-g`, or `C-[` dismisses the popup and returns to normal filer operation.
- **Rationale**: New users and occasional users need a quick reference for available keystrokes without consulting external documentation. A built-in help system improves discoverability and reduces the learning curve.
- **Satisfaction Criteria**:
  - Pressing `?` in the filer view opens a scrollable popup listing all keybindings.
  - The popup displays key combinations and their functions in a readable format.
  - Pressing `?` again dismisses the popup (toggle behavior).
  - Standard exit keys (`q`, `C-g`, `C-[`) also dismiss the popup.
  - Navigation keys (`C-n`/`C-p`/`up`/`down`/`pgup`/`pgdn`) scroll through the help content.
  - The popup follows the existing menu widget pattern for consistent UI behavior.
- **Validation Criteria**:
  - Manual verification that `?` opens the popup with all keystrokes visible.
  - Manual verification that scrolling works correctly.
  - Manual verification that all exit methods (`?`, `q`, `C-g`, `C-[`) dismiss the popup.
  - Token validation confirms `[REQ:HELP_POPUP]`, `[ARCH:HELP_WIDGET]`, and `[IMPL:HELP_POPUP]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Help Widget [ARCH:HELP_WIDGET]
- **Implementation**: See `implementation-decisions.md` § Help Popup Implementation [IMPL:HELP_POPUP]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-15)**:
- `?` key opens scrollable help popup
- Standard exit keys work (`?`, `q`, `C-g`, `C-[`)
- Navigation keys scroll through content

### [REQ:SYNC_COMMANDS] Sync Command Operations

**Priority: P1 (Important)**

- **Description**: Goful must provide a sync command mode activated by pressing `S`, which allows executing file operations (copy, delete, rename) sequentially across all workspace panes on files with the same name as the cursor file in the focused pane. The user presses `S` followed by an operation key (`c` for copy, `d` for delete, `r` for rename). A single prompt gathers necessary input (new filename for copy/rename) or confirmation (for delete), then the operation executes in each pane starting with the focused pane, finding files with the same name in each pane's current directory. By default, the operation aborts on first failure. Pressing `!` while in sync mode toggles "ignore failures" mode that continues through all panes, reporting failures at the end.
- **Rationale**: Users frequently need to perform identical operations across synchronized directory structures (e.g., renaming a file that exists in multiple mirrored directories). Without sync commands, users must repeat the same operation manually in each pane, which is tedious and error-prone. This feature complements the existing linked navigation mode by providing batch operations.
- **Satisfaction Criteria**:
  - Pressing `S` activates sync command mode and displays a prompt indicating the mode.
  - After `S`, pressing `c`, `d`, or `r` initiates the corresponding operation. Pressing `!` toggles ignore-failures mode.
  - For copy: user is prompted for a new filename (default: current filename); the file at cursor (and same-named files in other panes) are copied to the new filename within each pane's directory. User must specify a different name than the original.
  - For rename: user is prompted for new name once; the file at cursor (and same-named files in other panes) are renamed to the new name in each pane.
  - For delete: user confirms once; the file at cursor (and same-named files in other panes) are deleted in each pane.
  - Operations execute sequentially starting from the focused pane, then proceeding through other panes.
  - If a file with the target name doesn't exist in a pane, that pane is skipped (not treated as failure).
  - Default behavior: abort on first actual failure, report which pane failed.
  - Ignore-failures mode (toggle with `!`): continue through all panes even on failures, report all failures at the end.
  - Single-pane workspaces: behave exactly like the regular single-file operation.
  - User can cancel the prompt to abort before any changes are made.
- **Validation Criteria**:
  - Unit tests cover the file-by-name lookup helper in directories.
  - Unit tests cover the execution engine for success, failure, and skip scenarios.
  - Integration tests verify sequential execution across multiple panes.
  - Manual verification confirms the UI flow (prefix key → operation key → prompt → execution).
  - Token validation confirms `[REQ:SYNC_COMMANDS]`, `[ARCH:SYNC_MODE]`, and `[IMPL:SYNC_EXECUTE]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Sync Mode [ARCH:SYNC_MODE]
- **Implementation**: See `implementation-decisions.md` § Sync Execute [IMPL:SYNC_EXECUTE]

**Status**: ✅ Implemented

### [REQ:MOUSE_FILE_SELECT] Mouse File Selection

**Priority: P1 (Important)**

- **Description**: Goful must support mouse input for file selection. Left-clicking on a file in a directory window moves the cursor to that file. Clicking in an unfocused window switches focus to that window before selecting the file. Double-clicking on a file or directory enters/opens it. The feature requires enabling tcell mouse events and implementing hit-testing to convert screen coordinates to file list indices.
- **Rationale**: Mouse input is a common expectation for graphical terminal applications. Supporting mouse selection reduces the learning curve for new users familiar with GUI file managers, provides an alternative navigation method for users who prefer the mouse, and improves accessibility for users who have difficulty with keyboard navigation.
- **Satisfaction Criteria**:
  - tcell mouse events are enabled via `screen.EnableMouse()` during initialization.
  - Left-click on a file in any visible directory window moves the cursor to that file.
  - Clicking in an unfocused directory window first switches focus to that window, then selects the file.
  - Double-click on a file or directory opens/enters it (same as pressing Enter).
  - Mouse wheel scrolls the file list in the directory under the cursor.
  - Mouse events are handled in the main event loop alongside keyboard and resize events.
  - The feature can be enabled/disabled via a configuration option or environment variable (future).
- **Validation Criteria**:
  - Unit tests cover hit-testing functions (`FileIndexAtY`, `DirectoryAt`, `Contains`) with various window configurations and scroll positions.
  - Integration tests verify mouse event routing from the event loop to the appropriate widget.
  - Manual verification confirms file selection, focus switching, double-click, and scrolling work correctly on macOS and Linux terminals.
  - Token validation confirms `[REQ:MOUSE_FILE_SELECT]`, `[ARCH:MOUSE_EVENT_ROUTING]`, and `[IMPL:*]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Mouse Event Routing [ARCH:MOUSE_EVENT_ROUTING]
- **Implementation**: See `implementation-decisions.md` § Mouse Hit Testing [IMPL:MOUSE_HIT_TEST], Mouse File Selection [IMPL:MOUSE_FILE_SELECT]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-17)**:
- Unit tests in `widget/widget_test.go` (`TestWindowContains_REQ_MOUSE_FILE_SELECT`) covering hit-testing.
- Unit tests in `filer/integration_test.go` (`TestDirectoryAt_REQ_MOUSE_FILE_SELECT`, `TestFileIndexAtY_REQ_MOUSE_FILE_SELECT`, `TestDirectoryContains_REQ_MOUSE_FILE_SELECT`) covering directory hit-testing and file index conversion.
- Mouse events enabled via `screen.EnableMouse()` in `widget.Init()`.
- Event loop extended with `*tcell.EventMouse` case in `app/goful.go`.
- Left-click file selection, focus switching, and wheel scrolling operational.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1137 token references across 74 files.`

### [REQ:MOUSE_DOUBLE_CLICK] Mouse Double-Click Behavior

**Priority: P1 (Important)**

- **Description**: Goful must support double-click behavior for mouse navigation. Double-clicking on a directory navigates into it (respecting the Linked navigation mode). Double-clicking on a file opens it, and when Linked mode is enabled, opens same-named files in all windows where they exist. The feature requires time-based double-click detection and integration with existing navigation and file-open mechanisms.
- **Rationale**: Double-click is the standard mouse gesture for "open" or "enter" actions in file managers. Supporting double-click provides intuitive interaction for users accustomed to GUI file managers and completes the mouse navigation experience started with single-click selection.
- **Satisfaction Criteria**:
  - Double-click detection uses time-based threshold (300-500ms) and position matching.
  - Double-click on a directory calls `EnterDir()` to navigate into it.
  - When Linked mode is ON, double-clicking a directory triggers `ChdirAllToSubdir()` to propagate navigation to all windows.
  - Double-click on a file triggers the open action (equivalent to pressing Enter).
  - When Linked mode is ON, double-clicking a file opens same-named files in all windows where they exist.
  - When Linked mode is OFF, only the clicked file/directory is affected.
  - State tracking for last click time and position is maintained in `app.Goful`.
- **Validation Criteria**:
  - Unit tests cover double-click timing detection logic with various thresholds.
  - Unit tests cover directory double-click with Linked mode ON and OFF.
  - Unit tests cover file double-click with Linked mode ON and OFF.
  - Manual verification confirms double-click opens directories and files correctly.
  - Token validation confirms `[REQ:MOUSE_DOUBLE_CLICK]`, `[ARCH:MOUSE_DOUBLE_CLICK]`, and `[IMPL:MOUSE_DOUBLE_CLICK]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Mouse Double-Click Detection [ARCH:MOUSE_DOUBLE_CLICK]
- **Implementation**: See `implementation-decisions.md` § Mouse Double-Click Detection [IMPL:MOUSE_DOUBLE_CLICK]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-17)**:
- Unit tests in `app/mouse_test.go` covering double-click detection timing and position logic.
- Implementation in `app/goful.go`: `isDoubleClick()`, `handleDoubleClickDir()`, `handleDoubleClickFile()`.
- Directory double-click respects `linkedNav` via `ChdirAllToSubdirNoRebuild()`.
- File double-click respects `linkedNav` by moving cursor to same-named files in all windows.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1158 token references across 74 files.`

### [REQ:MOUSE_CROSS_WINDOW_SYNC] Mouse Cross-Window Cursor Synchronization

**Priority: P1 (Important)**

- **Description**: When a user clicks on a file in the active (focused) window and Linked navigation mode is ON, the cursor in all other windows should move to the same filename if it exists. When Linked mode is OFF, only the focused window's cursor moves. Focus remains on the active window. This provides visual feedback showing matching files across all workspace panes.
- **Rationale**: Users comparing directories benefit from seeing the same file highlighted across all windows. This complements the comparison color feature and helps users quickly identify matching files without manually navigating each pane. Respecting the Linked toggle provides consistent behavior with keyboard navigation.
- **Satisfaction Criteria**:
  - When Linked mode is ON, left-clicking a file in the active window moves cursors in all other windows to the same filename (if present).
  - When Linked mode is OFF, left-clicking a file only moves the cursor in the clicked window.
  - Focus remains on the clicked window; other windows do not gain focus.
  - Windows where the filename does not exist have their cursor highlight erased (no file highlighted) when Linked mode is ON.
  - Keyboard navigation in any window restores the cursor highlight.
  - Respects Linked Navigation mode: syncs when ON, individual when OFF.
- **Validation Criteria**:
  - Unit tests verify `SetCursorByNameAll` is called after cursor selection.
  - Manual verification confirms cross-window highlighting with 2+ panes.
  - Token validation confirms `[REQ:MOUSE_CROSS_WINDOW_SYNC]` references exist.
- **Architecture**: See `architecture-decisions.md` § Mouse Cross-Window Sync [ARCH:MOUSE_CROSS_WINDOW_SYNC]
- **Implementation**: See `implementation-decisions.md` § Mouse Cross-Window Sync [IMPL:MOUSE_CROSS_WINDOW_SYNC]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-18)**:
- Unit tests in `filer/integration_test.go` (`TestSetCursorByNameAll_REQ_MOUSE_CROSS_WINDOW_SYNC`, `TestSetCursorByNameAllFocusUnchanged_REQ_MOUSE_CROSS_WINDOW_SYNC`) covering cursor sync and focus preservation.
- Implementation in `app/goful.go`: `handleLeftClick` conditionally calls `SetCursorByNameAll` when `g.IsLinkedNav()` is true.
- Linked cursor sync tests in `app/linked_cursor_test.go` verifying behavior respects Linked mode toggle.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1491 token references across 78 files.`

### [REQ:TOOLBAR_PARENT_BUTTON] Toolbar Parent Navigation Button

**Priority: P1 (Important)**

- **Description**: Goful must provide a clickable button in the filer header row that navigates to the parent directory. When clicked, the button behaves identically to the keyboard `backspace`/`C-h`/`u` commands, respecting the current Linked navigation mode setting: if Linked mode is ON, all workspace windows navigate to their respective parent directories; if Linked mode is OFF, only the focused window navigates. This button is the first element of a planned mouse-first toolbar.
- **Rationale**: Mouse-first users need accessible UI controls for common navigation operations without relying on keyboard shortcuts. The parent directory navigation is one of the most frequently used operations and provides a natural starting point for a toolbar. Positioning the button in the header row keeps it visible and accessible while minimizing UI real estate usage.
- **Satisfaction Criteria**:
  - A clickable `[^]` button appears at the left edge of the filer header row, before the workspace tabs.
  - Left-clicking the button navigates to the parent directory following the current Linked mode setting.
  - When Linked mode is ON, all workspace windows navigate to their respective parent directories.
  - When Linked mode is OFF, only the focused window navigates to its parent directory.
  - The comparison index is rebuilt after navigation (same as keyboard parent navigation).
  - The button does not interfere with existing header elements (workspace tabs, LINKED indicator, directory tabs).
- **Validation Criteria**:
  - Unit tests cover the toolbar button bounds calculation and hit-testing.
  - Unit tests verify Linked mode behavior dispatch (navigate all vs. navigate focused).
  - Manual verification confirms button click triggers parent navigation on macOS and Linux terminals.
  - Token validation confirms `[REQ:TOOLBAR_PARENT_BUTTON]`, `[ARCH:TOOLBAR_LAYOUT]`, and `[IMPL:TOOLBAR_PARENT_BUTTON]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Toolbar Layout [ARCH:TOOLBAR_LAYOUT]
- **Implementation**: See `implementation-decisions.md` § Toolbar Parent Button Implementation [IMPL:TOOLBAR_PARENT_BUTTON]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-18)**:
- Unit tests in `filer/toolbar_test.go` (`TestToolbarButtonAt_REQ_TOOLBAR_PARENT_BUTTON`, `TestToolbarButtonAtMultipleButtons_REQ_TOOLBAR_PARENT_BUTTON`, `TestInvokeToolbarButton_REQ_TOOLBAR_PARENT_BUTTON`, `TestInvokeToolbarButtonWithNilCallback_REQ_TOOLBAR_PARENT_BUTTON`) covering hit-testing and callback invocation.
- Implementation in `filer/filer.go`: `drawHeader()` renders `[^]` button and stores bounds, `ToolbarButtonAt()` provides hit-testing, `InvokeToolbarButton()` dispatches actions.
- Implementation in `app/goful.go`: `handleLeftClick()` checks toolbar before directory selection, `HandleParentButtonPress()` implements linked parent navigation.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1326 token references across 77 files.`

### [REQ:TOOLBAR_LINKED_TOGGLE] Toolbar Linked Mode Toggle Button

**Priority: P1 (Important)**

- **Description**: Goful must provide a clickable `[L]` button in the filer header toolbar (next to the parent `[^]` button) that displays and toggles the linked navigation mode. The button is always visible, with its visual style reflecting the current state: reverse style when linked mode is ON, normal style when OFF. Clicking the button toggles the linked mode state and displays a confirmation message. This button replaces the existing conditional `[LINKED]` indicator that only appeared when the mode was enabled.
- **Rationale**: The existing `[LINKED]` indicator only shows when linked mode is ON, making it unclear when the mode is OFF. A toggle button provides both state visibility and quick access to toggle the mode, improving discoverability and mouse-first user experience. Placing it in the toolbar next to the parent button creates a consistent toolbar UI.
- **Satisfaction Criteria**:
  - A clickable `[L]` button appears in the filer header immediately after the `[^]` parent button.
  - The button uses reverse style when linked mode is ON, normal style when OFF.
  - Left-clicking the button toggles the linked navigation mode.
  - A confirmation message is displayed after toggle (e.g., "linked navigation enabled/disabled").
  - The existing conditional `[LINKED]` indicator is removed.
  - The button does not interfere with other header elements.
- **Validation Criteria**:
  - Unit tests cover the linked button bounds calculation and hit-testing.
  - Unit tests verify linked button invocation triggers the toggle callback.
  - Manual verification confirms button click toggles linked mode and displays message.
  - Token validation confirms `[REQ:TOOLBAR_LINKED_TOGGLE]`, `[ARCH:TOOLBAR_LAYOUT]`, and `[IMPL:TOOLBAR_LINKED_TOGGLE]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Toolbar Layout [ARCH:TOOLBAR_LAYOUT]
- **Implementation**: See `implementation-decisions.md` § Toolbar Linked Toggle Implementation [IMPL:TOOLBAR_LINKED_TOGGLE]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-18)**:
- Unit tests in `filer/toolbar_test.go` (`TestToolbarLinkedButtonHit_REQ_TOOLBAR_LINKED_TOGGLE`, `TestInvokeToolbarLinkedButton_REQ_TOOLBAR_LINKED_TOGGLE`, `TestInvokeToolbarLinkedButtonWithNilCallback_REQ_TOOLBAR_LINKED_TOGGLE`) covering hit-testing and callback invocation.
- Implementation in `filer/filer.go`: `drawHeader()` renders `[L]` button with state-based styling, stores bounds, `SetToolbarLinkedToggleFn()` sets callback, `InvokeToolbarButton()` handles "linked" case.
- Implementation in `main.go`: Callback wiring invokes `ToggleLinkedNav()` and displays confirmation message.
- Conditional `[LINKED]` indicator removed from header; replaced by always-visible `[L]` toggle button.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1360 token references across 78 files.`

### [REQ:TOOLBAR_COMPARE_BUTTON] Toolbar Compare Digest Button

**Priority: P1 (Important)**

- **Description**: Goful must provide a clickable `[=]` button in the filer header toolbar (next to the linked `[L]` button) that triggers digest comparison for all files appearing in multiple windows. Clicking the button calculates xxHash64 digests for every filename that exists in more than one directory pane, equivalent to pressing `=` on each shared file individually.
- **Rationale**: Users comparing directories need a single-click way to calculate digests for all shared files without manually pressing `=` on each file. This improves workflow efficiency for backup verification and duplicate detection tasks where content verification is needed across many files.
- **Satisfaction Criteria**:
  - A clickable `[=]` button appears in the filer header immediately after the `[L]` linked toggle button.
  - The button uses normal style (not reverse) since it triggers an action rather than displaying state.
  - Left-clicking the button calculates digests for all files that appear in 2+ directory panes.
  - A summary message is displayed after completion (e.g., "calculated digests for N files across M filenames").
  - The button does not interfere with other header elements.
- **Validation Criteria**:
  - Unit tests cover the compare button bounds calculation and hit-testing.
  - Unit tests verify compare button invocation triggers the callback.
  - Unit tests verify SharedFilenames() method returns correct filenames.
  - Manual verification confirms button click calculates digests and displays message.
  - Token validation confirms `[REQ:TOOLBAR_COMPARE_BUTTON]`, `[ARCH:TOOLBAR_LAYOUT]`, and `[IMPL:TOOLBAR_COMPARE_BUTTON]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Toolbar Layout [ARCH:TOOLBAR_LAYOUT]
- **Implementation**: See `implementation-decisions.md` § Toolbar Compare Button Implementation [IMPL:TOOLBAR_COMPARE_BUTTON]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-17)**:
- Unit tests in `filer/toolbar_test.go` (`TestToolbarCompareButtonHit_REQ_TOOLBAR_COMPARE_BUTTON`, `TestInvokeToolbarCompareButton_REQ_TOOLBAR_COMPARE_BUTTON`, `TestInvokeToolbarCompareButtonWithNilCallback_REQ_TOOLBAR_COMPARE_BUTTON`) covering hit-testing and callback invocation.
- Unit tests in `filer/compare_test.go` (`TestSharedFilenames_NilIndex_REQ_TOOLBAR_COMPARE_BUTTON`, `TestSharedFilenames_EmptyIndex_REQ_TOOLBAR_COMPARE_BUTTON`, `TestSharedFilenames_WithFiles_REQ_TOOLBAR_COMPARE_BUTTON`) covering SharedFilenames() method.
- Implementation in `filer/filer.go`: `drawHeader()` renders `[=]` button, stores bounds, `SetToolbarCompareDigestFn()` sets callback, `InvokeToolbarButton()` handles "compare" case.
- Implementation in `filer/compare.go`: `SharedFilenames()` method returns all filenames in the comparison index.
- Implementation in `main.go`: Callback wiring iterates shared files and calls `CalculateDigestForFile()` for each, displays summary message.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1468 token references across 78 files.`

### [REQ:TOOLBAR_SYNC_BUTTONS] Toolbar Sync Operation Buttons

**Priority: P1 (Important)**

- **Description**: Goful must provide four clickable buttons in the filer header toolbar after the `[=]` compare button:
  - `[C]` - Copy operation button
  - `[D]` - Delete operation button
  - `[R]` - Rename operation button
  - `[!]` - Ignore-failures mode toggle
  
  These buttons obey the Linked navigation mode:
  - **Linked mode ON**: Buttons trigger Sync operations (copy/delete/rename the same-named file across ALL workspace windows)
  - **Linked mode OFF**: Buttons trigger single-window operations (copy/delete/rename only in the focused window)
  
  The `[!]` button toggles a persistent ignore-failures mode that affects all Sync operations. When enabled, Sync operations continue even if some panes fail (file not found, permission denied, etc.). The button displays reverse style when ignore-failures is ON, normal style when OFF.

- **Rationale**: The existing `S` prefix key for Sync operations requires multiple keystrokes and memorization. Toolbar buttons provide:
  - Single-click access to common file operations
  - Visual indication of the current ignore-failures mode state
  - Consistent mouse-first experience matching the `[^]`, `[L]`, `[=]` buttons
  - Context-aware behavior based on Linked mode without separate keyboard shortcuts

- **Satisfaction Criteria**:
  - Four clickable buttons `[C]`, `[D]`, `[R]`, `[!]` appear in the filer header after `[=]`
  - `[C]`, `[D]`, `[R]` use normal style (action buttons)
  - `[!]` uses reverse style when ignore-failures is ON, normal style when OFF
  - When Linked mode is ON:
    - `[C]` triggers `startSyncCopy` (copies file to new name across all panes)
    - `[D]` triggers `startSyncDelete` (deletes same-named file in all panes)
    - `[R]` triggers `startSyncRename` (renames same-named file in all panes)
  - When Linked mode is OFF:
    - `[C]` triggers `Copy()` (copy in current window only)
    - `[D]` triggers `Remove()` (delete in current window only)
    - `[R]` triggers `Rename()` (rename in current window only)
  - `[!]` toggles ignore-failures state and displays confirmation message
  - Buttons do not interfere with other header elements

- **Validation Criteria**:
  - Unit tests cover button bounds calculation and hit-testing for all four buttons
  - Unit tests verify button invocation triggers correct callbacks
  - Manual verification confirms buttons trigger correct operations based on Linked mode
  - Token validation confirms `[REQ:TOOLBAR_SYNC_BUTTONS]`, `[ARCH:TOOLBAR_LAYOUT]`, and `[IMPL:TOOLBAR_SYNC_*]` references exist

- **Architecture**: See `architecture-decisions.md` § Toolbar Layout [ARCH:TOOLBAR_LAYOUT]
- **Implementation**: See `implementation-decisions/IMPL-TOOLBAR_SYNC_COPY.md`, `IMPL-TOOLBAR_SYNC_DELETE.md`, `IMPL-TOOLBAR_SYNC_RENAME.md`, `IMPL-TOOLBAR_IGNORE_FAILURES.md`

**Status**: ✅ Implemented

**Validation Evidence (2026-01-18)**:
- Unit tests in `filer/toolbar_test.go` (`TestToolbarSyncButtonsHit_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarSyncCopyButton_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarSyncDeleteButton_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarSyncRenameButton_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarIgnoreFailuresButton_REQ_TOOLBAR_SYNC_BUTTONS`, `TestInvokeToolbarSyncButtonsWithNilCallback_REQ_TOOLBAR_SYNC_BUTTONS`, `TestIgnoreFailuresIndicator_REQ_TOOLBAR_SYNC_BUTTONS`) covering hit-testing and callback invocation.
- Implementation in `filer/filer.go`: `drawHeader()` renders `[C]`, `[D]`, `[R]`, `[!]` buttons, stores bounds, `SetToolbarSyncCopyFn()`, `SetToolbarSyncDeleteFn()`, `SetToolbarSyncRenameFn()`, `SetToolbarIgnoreFailuresFn()`, `SetSyncIgnoreFailuresIndicator()` set callbacks, `InvokeToolbarButton()` handles new cases.
- Implementation in `app/goful.go`: `syncIgnoreFailures`, `IsSyncIgnoreFailures()`, `ToggleSyncIgnoreFailures()` for state management.
- Implementation in `app/window_wide.go`: `StartSyncCopy()`, `StartSyncDelete()`, `StartSyncRename()` exported for callback wiring.
- Implementation in `main.go`: Callback wiring with Linked mode logic.
- Token validation: `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1621 token references across 79 files.`

### [REQ:HELP_POPUP_STYLING] Help Popup Styling and Mouse Support

**Priority: P2 (Nice-to-have)**

- **Description**: The help popup shall display with a unified color scheme (border, section headers, key names, descriptions) and support mouse wheel scrolling for navigation.
  - Section headers (`=== Navigation ===`) render with distinct highlight color
  - Key names (left column) render with accent color
  - Descriptions (right column) render with default/subtle color
  - Popup border uses accent color distinguishing it from filer borders
  - Mouse wheel up/down scrolls help content

- **Rationale**:
  - Improves visual distinction between help content types (headers vs. key bindings vs. descriptions)
  - Mouse scroll support provides accessibility parity with keyboard navigation
  - Unified styling aligns with goful's existing color themes (default, midnight, black, white)

- **Satisfaction Criteria**:
  - Section headers render with `look.HelpHeader()` style (bold, accent color)
  - Key binding names render with `look.HelpKey()` style (accent color)
  - Description text renders with `look.HelpDesc()` style (default)
  - Popup border uses `look.HelpBorder()` style
  - Mouse wheel events scroll the help popup when displayed
  - Styling is consistent across all four themes

- **Validation Criteria**:
  - Manual visual verification across all four themes (default, midnight, black, white)
  - Mouse wheel scrolls help popup when displayed
  - Unit tests for content drawer style selection logic
  - Token validation confirms `[REQ:HELP_POPUP_STYLING]`, `[ARCH:HELP_STYLING]`, and `[IMPL:HELP_STYLING]` references exist

- **Architecture**: See `architecture-decisions.md` § Help Styling [ARCH:HELP_STYLING]
- **Implementation**: See `implementation-decisions/IMPL-HELP_STYLING.md`

**Status**: ✅ Implemented

**Validation Evidence (2026-01-18)**:
- `go test ./...` - all tests pass (darwin/arm64, Go 1.24.3)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1708 token references across 79 files.`
- Implementation in `look/look.go`: Theme-aware help styles with `HelpBorder()`, `HelpHeader()`, `HelpKey()`, `HelpDesc()` accessors
- Implementation in `help/help.go`: Custom `helpEntry` drawer with styled borders, headers, and key bindings
- Implementation in `app/goful.go`: Mouse wheel forwarding to modal widgets for scroll support

### [REQ:ESCAPE_KEY_BEHAVIOR] Escape Key Closes Modal Widgets

**Priority: P0 (Critical)**

- **Description**: The Escape key shall close modal widgets (help popup, menus, dialogs) with the same behavior as Ctrl+[.
  - `tcell.KeyEscape` must be mapped to `"C-["` in `widget.EventToString`
  - Modal widgets already handle `case "C-["` for exit, so no widget changes required

- **Rationale**:
  - tcell distinguishes `KeyEscape` (code 27) from `KeyCtrlLeftSq` (code 91) even though both represent escape sequences
  - Without this mapping, pressing Escape caused `EventToString` to return `"\u0000"` (null) which never matched `"C-["` in exit cases
  - Users expect Escape to close popups—this is a standard UI convention

- **Satisfaction Criteria**:
  - `tcell.KeyEscape` maps to `"C-["` in `keyToSting` map
  - Help popup closes when Escape is pressed
  - All modal widgets using `case "C-["` automatically work with Escape

- **Validation Criteria**:
  - Manual verification: press Escape in help popup, confirm it closes
  - Runtime debug logs confirm: `KeyEscape` → `"C-["` → exit triggered

- **Architecture**: See `architecture-decisions.md` § Escape Translation [ARCH:ESCAPE_TRANSLATION]
- **Implementation**: See `implementation-decisions/IMPL-ESCAPE_TRANSLATION.md`

**Status**: ✅ Implemented

**Validation Evidence (2026-01-18)**:
- `go test ./...` - all tests pass (darwin/arm64, Go 1.24.3)
- Runtime debug confirmed: Before fix `"inKeyMap":false,"result":"\u0000"` → After fix `"inKeyMap":true,"result":"C-["`
- Implementation in `widget/widget.go`: `tcell.KeyEscape: "C-["` entry in `keyToSting` map

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

### [REQ:BATCH_DIFF_REPORT] Batch Diff Report CLI Command

**Priority: P1 (Important)**

- **Description**: Goful must provide a `--diff-report` CLI flag that performs a complete non-interactive directory tree comparison across 2 or more directories, outputs a structured YAML report to stdout, and exits without launching the interactive TUI. Progress updates are written to stderr periodically (suppressible with `--quiet`). This feature enables scripted/automated directory comparison workflows using the same traversal algorithm as the interactive `[` and `]` diff search commands.
- **Rationale**: Users need programmatic access to directory comparison results for automation, CI pipelines, backup verification scripts, and other batch processing scenarios. The interactive diff search is powerful but requires manual interaction. A CLI command that produces machine-readable output enables integration with other tools while reusing the proven comparison algorithm.
- **Satisfaction Criteria**:
  - `goful --diff-report dir1 dir2 [dir3 ...]` compares directories and outputs YAML to stdout.
  - Requires at least 2 directories; exits with error if fewer provided.
  - `--quiet` flag suppresses progress output to stderr.
  - Output YAML includes: `directories` list, `totalFilesChecked`, `totalDirectoriesTraversed`, `durationSeconds`, and `differences` array.
  - Each difference entry includes: `name`, `path`, `reason`, and `isDir` fields.
  - Progress updates to stderr every 2 seconds (configurable) showing current path and stats.
  - Exit code 0 on success, 1 on error, 2 if differences found (for scripting).
  - No TUI initialization when `--diff-report` is set.
- **Validation Criteria**:
  - Unit tests cover `BatchNavigator` directory loading and traversal without TUI dependencies.
  - Unit tests cover `RunBatchDiffSearch` collecting all differences into the report structure.
  - Integration tests verify YAML output format and exit codes.
  - Manual verification confirms progress output to stderr and quiet mode suppression.
  - Token validation confirms `[REQ:BATCH_DIFF_REPORT]`, `[ARCH:BATCH_DIFF_REPORT]`, and `[IMPL:BATCH_DIFF_REPORT]` references exist across docs, code, and tests.
- **Architecture**: See `architecture-decisions.md` § Batch Diff Report Architecture [ARCH:BATCH_DIFF_REPORT]
- **Implementation**: See `implementation-decisions.md` § Batch Diff Report Implementation [IMPL:BATCH_DIFF_REPORT]

**Status**: ✅ Implemented

**Validation Evidence (2026-01-17)**:
- 11 unit tests in `filer/diffsearch_test.go` covering BatchNavigator, RunBatchDiffSearch, and report structure
- CLI flags `--diff-report` and `--quiet` implemented in `main.go`
- Exit codes: 0 (no differences), 1 (error), 2 (differences found)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1429 token references across 78 files.`

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

---

### [REQ:DOCKER_INTERACTIVE_SETUP] Docker-Based Interactive Goful Execution

**Priority: P2 (Nice-to-Have)**

- **Description**: Provide a Docker-based setup that allows running Goful interactively in a Linux container with volume mounts for file operations. The setup should include a Dockerfile for building the Goful binary, docker-compose.yml for easy management, a helper script for quick launches, and proper terminal environment configuration.
- **Rationale**: Enables cross-platform development and testing workflows without requiring local Go toolchain installation. Allows developers to test Goful in a consistent Linux environment regardless of their host OS. Useful for CI/CD pipelines and reproducible development environments.
- **Satisfaction Criteria**:
  - Dockerfile builds Goful binary successfully in a Linux container
  - Container runs interactively with proper terminal emulation (`-it` flags)
  - Volume mounts allow file operations on host filesystem
  - Terminal environment variables (`TERM`, `COLORTERM`) configured correctly for tcell
  - Helper script simplifies container execution
  - docker-compose.yml provides convenient service management
- **Validation Criteria**:
  - Docker image builds without errors: `docker build -t goful:latest .`
  - Container runs interactively: `docker run -it --rm goful:latest`
  - Volume mounts work: files created/modified in container are visible on host
  - TUI renders correctly with colors (verify `TERM=xterm-256color` and `COLORTERM=truecolor`)
  - Helper script executes and passes arguments correctly
  - Manual verification on macOS/Windows hosts confirms cross-platform compatibility
- **Architecture**: See `architecture-decisions.md` § Docker Build Strategy [ARCH:DOCKER_BUILD_STRATEGY]
- **Implementation**: See `implementation-decisions.md` § Dockerfile Multistage Build [IMPL:DOCKERFILE_MULTISTAGE] and Docker Compose Configuration [IMPL:DOCKER_COMPOSE_CONFIG]

**Status**: ✅ Implemented

### [REQ:DOCKER_WINDOWS_CONTAINER] Windows Container Support for Goful Testing

**Priority: P2 (Nice-to-Have)**

- **Description**: Extend the existing Docker setup to support Windows containers, enabling testing of Goful on Windows Server environments. This requires a separate Dockerfile targeting Windows Server Core base images, Windows-specific docker-compose configuration, and a PowerShell helper script for Windows hosts.
- **Rationale**: Enables testing Goful on Windows Server environments without requiring dedicated Windows hardware. Complements the existing Alpine-based Linux container for cross-platform development validation. Windows containers can only run on Windows hosts (Docker Desktop with Windows containers mode or Windows Server with containers feature).
- **Satisfaction Criteria**:
  - `Dockerfile.windows` builds Goful binary for Windows (`GOOS=windows`)
  - Windows container runs interactively with proper console emulation (`-it` flags)
  - Volume mounts work with Windows paths (`C:\workspace`)
  - tcell terminal initialization works in Windows console environment
  - PowerShell helper script simplifies container execution on Windows hosts
  - docker-compose.windows.yml provides convenient Windows service management
- **Validation Criteria**:
  - Docker image builds without errors: `docker build -f Dockerfile.windows -t goful:windows .`
  - Container runs interactively on Windows host: `docker run -it --rm goful:windows`
  - Volume mounts work: files created/modified in container are visible on Windows host
  - TUI renders correctly in Windows console (basic keyboard and display functionality)
  - PowerShell helper script executes and passes arguments correctly
  - Testing on Windows Server 2022 or Windows 10/11 with Docker Desktop confirms compatibility
- **Limitations**:
  - Windows containers require Windows host (cannot run on Linux/macOS Docker hosts)
  - Windows Server Core base image is larger (~5GB vs Alpine's ~5MB)
  - Some terminal features (mouse input, resize events) may be limited in Windows containers
  - nsync features remain stub-only on Windows (darwin-only implementation)
- **Architecture**: See `architecture-decisions.md` § Docker Windows Build [ARCH:DOCKER_WINDOWS_BUILD]
- **Implementation**: See `implementation-decisions.md` § Windows Dockerfile [IMPL:DOCKERFILE_WINDOWS]

**Status**: ✅ Implemented (Windows host testing deferred)

**Implementation Progress (2026-01-18)**:
- ✅ `Dockerfile.windows` created with multi-stage build
- ✅ `docker-compose.windows.yml` created with Windows paths
- ✅ `docker-run.ps1` PowerShell helper script created
- ✅ Makefile targets added (`docker-build-windows`, `docker-run-windows`, `docker-clean`)
- ✅ README documentation added with Docker section
- ⏳ Windows host testing deferred (requires Windows host with Docker in Windows containers mode)

