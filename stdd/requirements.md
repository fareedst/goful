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
| [REQ:NSYNC_MULTI_TARGET] | Multi-target copy/move via nsync SDK | P1 | ⏳ Planned | [ARCH:NSYNC_INTEGRATION] | [IMPL:NSYNC_OBSERVER], [IMPL:NSYNC_COPY_MOVE] |
| [REQ:NSYNC_CONFIRMATION] | Confirmation before multi-target copy/move | P1 | ⏳ Planned | [ARCH:NSYNC_CONFIRMATION] | [IMPL:NSYNC_CONFIRMATION] |
| [REQ:HELP_POPUP] | Help popup displays keystroke catalog on ? key | P2 | ⏳ Planned | [ARCH:HELP_WIDGET] | [IMPL:HELP_POPUP] |

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

**Status**: ⏳ Planned
### [REQ:WORKSPACE_START_DIRS] Positional Startup Directories

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

**Status**: ⏳ Planned
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

**Status**: ⏳ Planned

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

**Status**: ⏳ Planned

### [REQ:LINKED_NAVIGATION] Linked Navigation Mode

**Priority: P1 (Important)**

- **Description**: Goful must support a toggleable "linked" navigation mode where directory navigation in the focused window propagates to all other directory windows in the current workspace. When enabled and the user navigates into a subdirectory, all other windows that contain a matching subdirectory also navigate to it. When the user presses backspace (parent directory), all windows navigate to their respective parent directories. The mode is on by default.
- **Rationale**: Operators comparing similar directory structures (e.g., syncing folder hierarchies, comparing release versions) benefit from synchronized navigation across panes. Manual navigation in each pane is tedious and error-prone when directory structures mirror each other.
- **Satisfaction Criteria**:
  - A toggle mechanism (`L` uppercase or `M-l` Alt+l) enables/disables linked navigation mode.
  - When disabled, navigation in the focused window affects only that window (historical behavior).
  - When enabled, entering a subdirectory attempts to navigate all other workspace windows to a matching subdirectory (by name) if it exists in each window's current path.
  - When enabled, pressing backspace (parent navigation) causes all windows to navigate to their respective parent directories.
  - When enabled, changing sort order applies the same sort type to all windows in the workspace.
  - A visual indicator (`[LINKED]`) appears in the filer header when the mode is active.
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

**Validation Evidence (2026-01-09)**:
- `TestChdirAllToSubdir_REQ_LINKED_NAVIGATION`, `TestChdirAllToParent_REQ_LINKED_NAVIGATION`, `TestLinkedNavigationSingleWindow_REQ_LINKED_NAVIGATION`, `TestSortAllBy_REQ_LINKED_NAVIGATION` in `filer/integration_test.go` covering the workspace navigation and sort helpers.
- Toggle keystroke: `L` (uppercase, works on all platforms) or `M-l` (Alt+l, may not work on macOS where Option produces special characters).
- Header indicator `[LINKED]` displayed when mode is active.
- Sort synchronization: all sort menu options apply to all windows when linked mode is enabled.

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

**Status**: ⏳ Planned

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

