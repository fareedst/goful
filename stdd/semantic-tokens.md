# Semantic Tokens Directory

**STDD Methodology Version**: 1.1.0

## Overview
This document serves as the **central directory/registry** for all semantic tokens used in the project. Semantic tokens (`[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`) provide a consistent vocabulary and traceability mechanism that ties together all documentation, code, and tests.

**For detailed information about tokens, see:**
- **Requirements tokens**: See `requirements.md` for full descriptions, rationale, satisfaction criteria, and validation criteria
- **Architecture tokens**: See `architecture-decisions.md` for architectural decisions, rationale, and alternatives considered
- **Implementation tokens**: See `implementation-decisions.md` for implementation details, code structures, and algorithms

## AI Assistant Integration Guidelines [REQ:DOC_016]

### Token Usage for AI Assistants

AI assistants should use semantic tokens for:

1. **Code Navigation**: Search for `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]` tokens to find related code
2. **Feature Understanding**: Trace features from requirements through architecture to implementation
3. **Change Impact Analysis**: Use token cross-references to identify affected components
4. **Test Discovery**: Find tests for features using `[REQ:*]` tokens in test names

### Token-Based Code Navigation

```bash
# Find all implementations of a requirement
grep -r "\[REQ:FEATURE_NAME\]" --include="*.go" .

# Find all tests for a requirement
grep -r "REQ_FEATURE_NAME" --include="*_test.go" .

# Find architecture decisions for a feature
grep -r "\[ARCH:FEATURE_NAME\]" --include="*.md" .

# Find implementation details
grep -r "\[IMPL:FEATURE_NAME\]" --include="*.go" .
```

### Token Creation Requirements

When implementing features:
1. **ALWAYS** create `[REQ:*]` token in `requirements.md` first
2. **ALWAYS** create `[ARCH:*]` token in `architecture-decisions.md` for design decisions
3. **ALWAYS** add `[IMPL:*]` tokens to code comments
4. **ALWAYS** reference `[REQ:*]` tokens in test names/comments
5. **ALWAYS** update `semantic-tokens.md` registry when creating new tokens
6. **ALWAYS** document any `[PROC:*]` process tokens in `processes.md` when defining operational workflows

### Token Audit Workflow `[PROC:TOKEN_AUDIT]`

- Map requirement → architecture → implementation tokens before touching code.
- Annotate every code edit with `[IMPL:*] [ARCH:*] [REQ:*]` (same triplet used in documentation).
- Require tests to include the `[REQ:*]` (and optional `[TEST:*]`) identifiers in both the test name and supporting comments.
- Record the audit result inside the relevant task/subtask so future agents can see when the chain was verified.

### Automated Validation `[PROC:TOKEN_VALIDATION]`

- Run `./scripts/validate_tokens.sh` (or repo-specific equivalent) after each audit to ensure every referenced token exists in the registry.
- Treat validation failures as blocking defects until the registry and documents are synchronized.
- Capture validation outputs in `implementation-decisions.md` so audits remain reproducible.

### Token Validation Requirements

Before marking features complete:
1. **ALWAYS** run token validation scripts (e.g., `./scripts/validate_tokens.sh`) and store the `[PROC:TOKEN_VALIDATION]` result in `implementation-decisions.md`.
2. **ALWAYS** ensure token consistency across all layers
3. **ALWAYS** verify token traceability in documentation
4. **ALWAYS** check that all cross-references are valid

## Token Format

```
[TYPE:IDENTIFIER]
```

## Token Types

- `[REQ:*]` - Requirements (functional/non-functional) - **The source of intent**
- `[ARCH:*]` - Architecture decisions - **High-level design choices that preserve intent**
- `[IMPL:*]` - Implementation decisions - **Low-level choices that preserve intent**
- `[TEST:*]` - Test specifications - **Validation of intent**
- `[PROC:*]` - Process definitions for survey/build/test/deploy work that stay linked to `[REQ:*]`

## Token Naming Convention

- Use UPPER_SNAKE_CASE for identifiers
- Be descriptive but concise
- Example: `[REQ:DUPLICATE_PREVENTION]` not `[REQ:DP]`

## Cross-Reference Format

When referencing other tokens:

```markdown
[IMPL:STATE_PATH_RESOLVER] Description [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]
```

## Requirements Tokens Registry

**📖 Full details**: See `requirements.md`

### Immutable Requirements

-### Core Functional Requirements
- `[REQ:STDD_SETUP]` - STDD methodology setup
- `[REQ:MODULE_VALIDATION]` - Independent module validation before integration
- `[REQ:CONFIGURABLE_STATE_PATHS]` - Configurable state/history persistence
- `[REQ:WORKSPACE_START_DIRS]` - Positional CLI directories seed workspace windows
- `[REQ:EXTERNAL_COMMAND_CONFIG]` - External command menu definitions are configurable via files/overrides
- `[REQ:BACKSPACE_BEHAVIOR]` - Backspace opens parent directory in filer views and deletes previous characters inside prompts across platforms
- `[REQ:FILER_EXCLUDE_NAMES]` - Configurable basename exclusions with runtime toggles
- `[REQ:WINDOW_MACRO_ENUMERATION]` - `%D@`/`%d@` enumerate other workspace directories (paths + basenames) for external commands
- `[REQ:GO_TOOLCHAIN_LTS]` - Go toolchain tracks current LTS baseline
- `[REQ:DEPENDENCY_REFRESH]` - Dependencies refreshed for security and compatibility
- `[REQ:CI_PIPELINE_CORE]` - CI pipeline runs fmt/vet/tests on every change
- `[REQ:STATIC_ANALYSIS]` - Static analysis (staticcheck/golangci-lint) guards regressions
- `[REQ:RACE_TESTING]` - Race-enabled test job in automation
- `[REQ:UI_PRIMITIVE_TESTS]` - UI widget primitives are covered by tests
- `[REQ:CMD_HANDLER_TESTS]` - Command handling/app modes are validated by tests
- `[REQ:INTEGRATION_FLOWS]` - Integration flows (open/navigate/rename/delete) are tested
- `[REQ:ARCH_DOCUMENTATION]` - Architecture is documented for the main packages and data flow
- `[REQ:CONTRIBUTING_GUIDE]` - Contributor guide defines standards and review expectations
- `[REQ:RELEASE_BUILD_MATRIX]` - Reproducible release builds via Makefile and matrix targets
- `[REQ:BEHAVIOR_BASELINE]` - Baseline interactions and key mappings are captured in automation/docs
- `[REQ:DEBT_TRIAGE]` - Technical debt and risky areas are triaged with TODOs/issues
- `[REQ:TERMINAL_PORTABILITY]` - Terminal launcher works across Linux, macOS, and tmux contexts
- `[REQ:TERMINAL_CWD]` - macOS terminal sessions start in the active directory
- `[REQ:EVENT_LOOP_SHUTDOWN]` - Event poller must observe shutdown signals and terminate without leaking goroutines
- `[REQ:CLI_TO_CHAINING]` - CLI helper rewrites commands by interleaving `--to` before every target argument with dry-run support
- `[REQ:FILE_COMPARISON_COLORS]` - Cross-directory file comparison with configurable color-coding for names, sizes, and times
- `[REQ:LINKED_NAVIGATION]` - Linked navigation mode that propagates directory changes across workspace windows
- `[REQ:DIFF_SEARCH]` - Cross-window difference search for finding missing or size-mismatched files
- `[REQ:NSYNC_MULTI_TARGET]` - Multi-target copy/move via nsync SDK to all visible workspace panes
- `[REQ:NSYNC_CONFIRMATION]` - Confirmation prompts before multi-target copy/move operations
- `[REQ:HELP_POPUP]` - Help popup displays keystroke catalog on `?` key
- Add your requirements tokens here

### Non-Functional Requirements
- `[REQ:PERFORMANCE]` - Performance requirements
- `[REQ:RELIABILITY]` - Reliability requirements
- `[REQ:MAINTAINABILITY]` - Maintainability requirements
- `[REQ:USABILITY]` - Usability requirements

## Architecture Tokens Registry

**📖 Full details**: See `architecture-decisions.md`

- `[ARCH:LANGUAGE_SELECTION]` - Language and runtime selection
- `[ARCH:PROJECT_STRUCTURE]` - Project structure decision
- `[ARCH:STDD_STRUCTURE]` - STDD project structure [REQ:STDD_SETUP]
- `[ARCH:MODULE_VALIDATION]` - Module validation strategy [REQ:MODULE_VALIDATION]
- `[ARCH:STATE_PATH_SELECTION]` - Path precedence for persistence [REQ:CONFIGURABLE_STATE_PATHS]
- `[ARCH:EXTERNAL_COMMAND_REGISTRY]` - External command config resolution/binding [REQ:EXTERNAL_COMMAND_CONFIG]
- `[ARCH:BACKSPACE_TRANSLATION]` - Backspace key canonicalization in the input translator [REQ:BACKSPACE_BEHAVIOR]
- `[ARCH:FILER_EXCLUDE_FILTER]` - Basename exclude filter + loader/toggle strategy [REQ:FILER_EXCLUDE_NAMES]
- `[ARCH:WINDOW_MACRO_ENUMERATION]` - Workspace directory enumeration for `%D@` [REQ:WINDOW_MACRO_ENUMERATION]
- `[ARCH:GO_RUNTIME_STRATEGY]` - Go LTS/toolchain policy [REQ:GO_TOOLCHAIN_LTS]
- `[ARCH:DEPENDENCY_POLICY]` - Dependency refresh and security approach [REQ:DEPENDENCY_REFRESH]
- `[ARCH:CI_PIPELINE]` - CI workflow design for fmt/vet/tests [REQ:CI_PIPELINE_CORE]
- `[ARCH:STATIC_ANALYSIS_POLICY]` - Static analysis gates [REQ:STATIC_ANALYSIS]
- `[ARCH:RACE_TESTING_PIPELINE]` - Race-enabled testing strategy [REQ:RACE_TESTING]
- `[ARCH:TEST_STRATEGY_UI]` - UI widget testing approach [REQ:UI_PRIMITIVE_TESTS]
- `[ARCH:TEST_STRATEGY_CMD]` - Command/app mode testing approach [REQ:CMD_HANDLER_TESTS]
- `[ARCH:TEST_STRATEGY_INTEGRATION]` - Integration flow testing approach [REQ:INTEGRATION_FLOWS]
- `[ARCH:DOCS_STRUCTURE]` - Architecture documentation structure [REQ:ARCH_DOCUMENTATION]
- `[ARCH:CONTRIBUTION_PROCESS]` - Contribution standards and review flow [REQ:CONTRIBUTING_GUIDE]
- `[ARCH:BUILD_MATRIX]` - Build/release matrix strategy [REQ:RELEASE_BUILD_MATRIX]
- `[ARCH:BASELINE_CAPTURE]` - Baseline behavior and keymap capture [REQ:BEHAVIOR_BASELINE]
- `[ARCH:DEBT_MANAGEMENT]` - Debt triage and TODO/issue tracking [REQ:DEBT_TRIAGE]
- `[ARCH:TOKEN_VALIDATION_AUTOMATION]` - Token validation helper automation [REQ:STDD_SETUP]
- `[ARCH:QUIT_DIALOG_KEYS]` - Quit dialog key translation guarantees [REQ:QUIT_DIALOG_DEFAULT]
- `[ARCH:TERMINAL_LAUNCHER]` - Cross-platform terminal launcher abstraction [REQ:TERMINAL_PORTABILITY]
- `[ARCH:EVENT_LOOP_SHUTDOWN]` - Event poller shutdown coordination [REQ:EVENT_LOOP_SHUTDOWN]
- `[ARCH:XFORM_CLI_PIPELINE]` - Parser/builder split for the xform helper that inserts `--to` between targets [REQ:CLI_TO_CHAINING]
- `[ARCH:WORKSPACE_BOOTSTRAP]` - Workspace seeding from positional startup directories [REQ:WORKSPACE_START_DIRS]
- `[ARCH:FILE_COMPARISON_ENGINE]` - Progressive comparison with cached indexing for cross-directory file color-coding [REQ:FILE_COMPARISON_COLORS]
- `[ARCH:LINKED_NAVIGATION]` - Linked navigation mode architecture [REQ:LINKED_NAVIGATION]
- `[ARCH:DIFF_SEARCH]` - Difference search engine architecture [REQ:DIFF_SEARCH]
- `[ARCH:NSYNC_INTEGRATION]` - nsync SDK integration for multi-target copy/move [REQ:NSYNC_MULTI_TARGET]
- `[ARCH:NSYNC_CONFIRMATION]` - Confirmation prompts before multi-target operations [REQ:NSYNC_CONFIRMATION]
- `[ARCH:HELP_WIDGET]` - Help widget architecture based on ListBox pattern [REQ:HELP_POPUP]
- Add your architecture tokens here

## Implementation Tokens Registry

**📖 Full details**: See `implementation-decisions.md`

- `[IMPL:CONFIG_STRUCT]` - Configuration structure implementation [ARCH:CONFIG_STRUCTURE] [REQ:CONFIGURATION]
- `[IMPL:STDD_FILES]` - STDD file creation [ARCH:STDD_STRUCTURE] [REQ:STDD_SETUP]
- `[IMPL:MODULE_VALIDATION]` - Module validation implementation [ARCH:MODULE_VALIDATION] [REQ:MODULE_VALIDATION]
- `[IMPL:STATE_PATH_RESOLVER]` - Resolver + bootstrap wiring [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]
- `[IMPL:EXTERNAL_COMMAND_LOADER]` - External command config loader/defaults [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]
- `[IMPL:BACKSPACE_TRANSLATION]` - Dual-key backspace normalization inside `widget.EventToString` [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]
- `[IMPL:EXTERNAL_COMMAND_APPEND]` - Default inheritance (prepended custom entries) + replacement toggle for external command configs [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]
- `[IMPL:EXTERNAL_COMMAND_BINDER]` - Menu binding helpers for commands [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]
- `[IMPL:FILER_EXCLUDE_RULES]` - Centralized basename filtering + toggles [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]
- `[IMPL:FILER_EXCLUDE_LOADER]` - Exclude list loader, resolver wiring, and UI toggle [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]
- `[IMPL:WINDOW_MACRO_ENUMERATION]` - `%D@`/`%d@` macro helpers [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]
- `[IMPL:GO_MOD_UPDATE]` - Go version/toolchain update [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS]
- `[IMPL:DEP_BUMP]` - Dependency refresh and tidy [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]
- `[IMPL:CI_WORKFLOW]` - GitHub Actions workflow for fmt/vet/tests [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]
- `[IMPL:STATICCHECK_SETUP]` - Static analysis configuration [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]
- `[IMPL:RACE_JOB]` - Race-enabled test job [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]
- `[IMPL:TEST_WIDGETS]` - Widget/UI primitive tests [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]
- `[IMPL:TEST_CMDLINE]` - Command handling tests [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]
- `[IMPL:TEST_INTEGRATION_FLOWS]` - Integration/snapshot flow tests [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]
- `[IMPL:DOC_ARCH_GUIDE]` - Architecture documentation [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]
- `[IMPL:DOC_CONTRIBUTING]` - CONTRIBUTING guide [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE]
- `[IMPL:MAKE_RELEASE_TARGETS]` - Makefile + matrix build targets [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]
- `[IMPL:BASELINE_SNAPSHOTS]` - Baseline behavior capture/tests [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE]
- `[IMPL:DEBT_TRACKING]` - Debt triage annotations/issues [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]
- `[IMPL:TOKEN_VALIDATION_SCRIPT]` - Shell script enforcing token registry [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]
- `[IMPL:QUIT_DIALOG_ENTER]` - Return/Enter mapping for quit dialog default [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]
- `[IMPL:TERMINAL_ADAPTER]` - Platform-aware terminal command adapter [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD]
- `[IMPL:EVENT_LOOP_SHUTDOWN]` - Event poller stop controller [ARCH:EVENT_LOOP_SHUTDOWN] [REQ:EVENT_LOOP_SHUTDOWN]
- `[IMPL:XFORM_CLI_SCRIPT]` - Bash helper and tests that expose `xform` with dry-run output [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]
- `[IMPL:WORKSPACE_START_DIRS]` - Parser + seeder for positional startup directories [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]
- `[IMPL:COMPARE_COLOR_CONFIG]` - Comparison color YAML configuration loader [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
- `[IMPL:FILE_COMPARISON_INDEX]` - Cached index of cross-directory file comparison states [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
- `[IMPL:COMPARISON_DRAW]` - Draw integration for comparison colors [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
- `[IMPL:DIGEST_COMPARISON]` - On-demand xxHash64 digest calculation for files with equal sizes [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
- `[IMPL:LINKED_NAVIGATION]` - Linked navigation mode implementation [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
- `[IMPL:LINKED_NAVIGATION_AUTO_DISABLE]` - Auto-disable linked navigation on partial subdirectory navigation failure [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
- `[IMPL:DIFF_SEARCH]` - Difference search implementation [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
- `[IMPL:NSYNC_OBSERVER]` - Observer adapter bridging nsync to goful progress widget [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
- `[IMPL:NSYNC_COPY_MOVE]` - nsync copy/move wrappers and CopyAll/MoveAll modes [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
- `[IMPL:NSYNC_CONFIRMATION]` - Confirmation prompts for multi-target copy/move operations [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
- `[IMPL:HELP_POPUP]` - Help popup implementation with keystroke catalog [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
- Add your implementation tokens here

## Test Tokens Registry

- `[TEST:KEYMAP_BASELINE]` - Captures canonical keybindings for filer/cmdline/finder/completion/menu (`main_keymap_test.go`) ensuring `[REQ:BEHAVIOR_BASELINE]` remains satisfied via `[IMPL:BASELINE_SNAPSHOTS]`.

## Token Relationships

### Hierarchical Relationships
- `[REQ:PARENT_FEATURE]` contains `[REQ:SUB_FEATURE_1]`, `[REQ:SUB_FEATURE_2]`
- `[ARCH:FEATURE]` includes `[ARCH:COMPONENT_1]`, `[ARCH:COMPONENT_2]`

### Flow Relationships
- `[REQ:FEATURE]` → `[ARCH:DESIGN]` → `[IMPL:IMPLEMENTATION]` → Code + Tests

### Dependency Relationships
- `[IMPL:FEATURE]` depends on `[ARCH:DESIGN]` and `[REQ:FEATURE]`
- `[ARCH:DESIGN]` depends on `[REQ:FEATURE]`

## Process Tokens Registry

**📖 Full details**: See `processes.md`

- `[PROC:PROJECT_SURVEY_AND_SETUP]` - Survey and readiness process supporting `[REQ:STDD_SETUP]` and `[ARCH:STDD_STRUCTURE]`
- `[PROC:BUILD_PIPELINE_VALIDATION]` - Build/deploy validation tied to `[REQ:MODULE_VALIDATION]`
- `[PROC:TOKEN_AUDIT]` - Mandatory checklist ensuring every requirement → architecture → implementation → code/test path is annotated and documented
- `[PROC:TOKEN_VALIDATION]` - Automated validation workflow (e.g., `./scripts/validate_tokens.sh`) that proves all referenced tokens exist in the registry
- `[PROC:TERMINAL_VALIDATION]` - Manual macOS/Linux terminal checklist safeguarding `[REQ:TERMINAL_PORTABILITY]`, `[REQ:TERMINAL_CWD]`, and `[ARCH:TERMINAL_LAUNCHER]`
- Add your process tokens here

## Usage Examples

### In Code Comments
```[your-language]
// [REQ:CONFIGURABLE_STATE_PATHS] Implementation of configurable state/history paths
// [IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]
function exampleFunction() {
    // ...
}
```
> **NOTE**: Code merged without these annotations is considered incomplete because it fails `[PROC:TOKEN_AUDIT]`.

### In Tests
```[your-language]
// Test validates [REQ:CONFIGURABLE_STATE_PATHS] is met
function testResolvePaths_REQ_CONFIGURABLE_STATE_PATHS() {
    // Test implementation
}
```
> **NOTE**: Tests without `[REQ:*]` markers are rejected during `[PROC:TOKEN_VALIDATION]` because they cannot prove intent.

### In Documentation
```markdown
The feature uses [ARCH:ARCHITECTURE_NAME] to fulfill [REQ:FEATURE_NAME].
Implementation details are documented in [IMPL:IMPLEMENTATION_NAME].
```

## Token Validation Guidelines

### Cross-Layer Token Consistency

Every feature must have proper token coverage across all layers:

1. **Requirements Layer**: Feature must have `[REQ:*]` token in `requirements.md`
2. **Architecture Layer**: Architecture decisions must have `[ARCH:*]` tokens in `architecture-decisions.md`
3. **Implementation Layer**: Implementation must have `[IMPL:*]` tokens in code comments
4. **Test Layer**: Tests must reference `[REQ:*]` tokens in test names/comments
5. **Documentation Layer**: All documentation must cross-reference tokens consistently

### Token Format Validation

1. **Token Format**: Must follow `[TYPE:IDENTIFIER]` pattern exactly
2. **Token Types**: Must use valid types (`REQ`, `ARCH`, `IMPL`, `TEST`, `PROC`)
3. **Identifier Format**: Must use UPPER_SNAKE_CASE
4. **Cross-References**: Implementation tokens must reference architecture and requirement tokens

### Token Traceability Validation

1. Every requirement in `requirements.md` must have corresponding implementation tokens
2. Every architecture decision must have corresponding implementation tokens
3. Every test must link to specific requirements via `[REQ:*]` tokens
4. All tokens must be discoverable through automated validation
## Token Creation Requirements

When implementing features:
1. **ALWAYS** create `[REQ:*]` token in `requirements.md` first
2. **ALWAYS** create `[ARCH:*]` token in `architecture-decisions.md` for design decisions
3. **ALWAYS** add `[IMPL:*]` tokens to code comments
4. **ALWAYS** reference `[REQ:*]` tokens in test names/comments
5. **ALWAYS** update `semantic-tokens.md` registry when creating new tokens
6. **ALWAYS** document any `[PROC:*]` process tokens in `processes.md` when defining operational workflows

