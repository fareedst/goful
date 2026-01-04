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

### Task 2.1: Core Feature Implementation
**Status:** ⏳ Pending  
**Priority:** P0 (Critical)  
**Semantic Tokens:** `[REQ:EXAMPLE_FEATURE]`, `[ARCH:EXAMPLE_DECISION]`, `[IMPL:EXAMPLE_IMPLEMENTATION]`

**Description**: Implement the core feature according to requirements and architecture.

**Subtasks**:
- [ ] Identify logical modules and document module boundaries [REQ:MODULE_VALIDATION]
- [ ] Define module interfaces and validation criteria [REQ:MODULE_VALIDATION]
- [ ] Develop Module 1 independently
- [ ] Validate Module 1 independently (unit tests, contract tests, edge cases, error handling) [REQ:MODULE_VALIDATION]
- [ ] Develop Module 2 independently
- [ ] Validate Module 2 independently (unit tests, contract tests, edge cases, error handling) [REQ:MODULE_VALIDATION]
- [ ] Integrate validated modules [REQ:MODULE_VALIDATION]
- [ ] Write integration tests for combined behavior
- [ ] Write end-to-end tests [REQ:EXAMPLE_FEATURE]
- [ ] Run `[PROC:TOKEN_AUDIT]` + `./scripts/validate_tokens.sh` and record outcomes [PROC:TOKEN_VALIDATION]

**Completion Criteria**:
- [ ] All modules identified and documented
- [ ] All modules validated independently before integration
- [ ] Integration tests pass
- [ ] All documentation updated
- [ ] Token audit + validation logged

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

## P0: CI & Static Analysis Foundation [REQ:CI_PIPELINE_CORE] [REQ:STATIC_ANALYSIS] [REQ:RACE_TESTING] [ARCH:CI_PIPELINE] [ARCH:STATIC_ANALYSIS_POLICY] [ARCH:RACE_TESTING_PIPELINE] [IMPL:CI_WORKFLOW] [IMPL:STATICCHECK_SETUP] [IMPL:RACE_JOB]

**Status**: ⏳ Pending

**Description**: Establish GitHub Actions for fmt/vet/tests, static analysis, and race job.

**Dependencies**: Modernize Toolchain and Dependencies

**Subtasks**:
- [ ] Add fmt/vet/test workflow with cache [REQ:CI_PIPELINE_CORE] [IMPL:CI_WORKFLOW]
- [ ] Add staticcheck (and optional golangci-lint) job [REQ:STATIC_ANALYSIS] [IMPL:STATICCHECK_SETUP]
- [ ] Add race-enabled test job [REQ:RACE_TESTING] [IMPL:RACE_JOB]
- [ ] Ensure jobs reference matching Go version [REQ:GO_TOOLCHAIN_LTS]
- [ ] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]`

**Completion Criteria**:
- [ ] CI runs fmt/vet/test/staticcheck/race on PRs
- [ ] Jobs pass on target branches
- [ ] Token audit + validation recorded

**Priority Rationale**: P0 to gate all future changes with automated checks.

## P0: Test Coverage for UI/Commands/Flows [REQ:UI_PRIMITIVE_TESTS] [REQ:CMD_HANDLER_TESTS] [REQ:INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_UI] [ARCH:TEST_STRATEGY_CMD] [ARCH:TEST_STRATEGY_INTEGRATION] [IMPL:TEST_WIDGETS] [IMPL:TEST_CMDLINE] [IMPL:TEST_INTEGRATION_FLOWS]

**Status**: ⏳ Pending

**Description**: Add coverage for widgets/filer, command handling, and integration flows (open/navigate/rename/delete).

**Dependencies**: CI & Static Analysis Foundation

**Subtasks**:
- [ ] Identify modules and validation criteria per area [REQ:MODULE_VALIDATION]
- [ ] Add widget/filer unit/snapshot tests [REQ:UI_PRIMITIVE_TESTS] [IMPL:TEST_WIDGETS]
- [ ] Add command/app mode tests [REQ:CMD_HANDLER_TESTS] [IMPL:TEST_CMDLINE]
- [ ] Add integration flow tests with fixtures [REQ:INTEGRATION_FLOWS] [IMPL:TEST_INTEGRATION_FLOWS]
- [ ] Document validation results before integration [REQ:MODULE_VALIDATION]
- [ ] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]`

**Completion Criteria**:
- [ ] Module validation evidence recorded
- [ ] Tests cover listed areas with tokens
- [ ] CI green with new coverage
- [ ] Token audit + validation recorded

**Priority Rationale**: P0 to secure behavior before major refactors.

## P1: Docs & Baselines [REQ:ARCH_DOCUMENTATION] [REQ:CONTRIBUTING_GUIDE] [REQ:BEHAVIOR_BASELINE] [ARCH:DOCS_STRUCTURE] [ARCH:CONTRIBUTION_PROCESS] [ARCH:BASELINE_CAPTURE] [IMPL:DOC_ARCH_GUIDE] [IMPL:DOC_CONTRIBUTING] [IMPL:BASELINE_SNAPSHOTS]

**Status**: ⏳ Pending

**Description**: Write `ARCHITECTURE.md`, `CONTRIBUTING.md`, and capture baseline keybindings/modes.

**Dependencies**: CI & Static Analysis Foundation

**Subtasks**:
- [ ] Draft architecture overview with package/data flow [REQ:ARCH_DOCUMENTATION] [IMPL:DOC_ARCH_GUIDE]
- [ ] Draft contributing guide with standards/review expectations [REQ:CONTRIBUTING_GUIDE] [IMPL:DOC_CONTRIBUTING]
- [ ] Capture baseline interactions/keymaps in tests/scripts [REQ:BEHAVIOR_BASELINE] [IMPL:BASELINE_SNAPSHOTS]
- [ ] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]`

**Completion Criteria**:
- [ ] Docs published and cross-linked
- [ ] Baseline tests run in CI
- [ ] Token audit + validation recorded

**Priority Rationale**: P1 to enable contributors and guard behavior.

## P1: Release Build Hygiene [REQ:RELEASE_BUILD_MATRIX] [ARCH:BUILD_MATRIX] [IMPL:MAKE_RELEASE_TARGETS]

**Status**: ⏳ Pending

**Description**: Add Makefile targets and CI matrix for reproducible static builds across GOOS/GOARCH.

**Dependencies**: Modernize Toolchain and Dependencies

**Subtasks**:
- [ ] Add lint/test/build targets to Makefile [REQ:RELEASE_BUILD_MATRIX] [IMPL:MAKE_RELEASE_TARGETS]
- [ ] Add CI matrix build job (e.g., linux/amd64, darwin/arm64) [REQ:RELEASE_BUILD_MATRIX]
- [ ] Verify artifacts names/digests deterministic [REQ:RELEASE_BUILD_MATRIX]
- [ ] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]`

**Completion Criteria**:
- [ ] Makefile + workflow committed with tokens
- [ ] Matrix builds succeed
- [ ] Token audit + validation recorded

**Priority Rationale**: P1 to prepare for reproducible releases.

## P1: Debt Triage [REQ:DEBT_TRIAGE] [ARCH:DEBT_MANAGEMENT] [IMPL:DEBT_TRACKING]

**Status**: ⏳ Pending

**Description**: Log known pain points (error handling, cross-platform quirks) and annotate risky areas with TODOs/owners.

**Dependencies**: None

**Subtasks**:
- [ ] Create issue list/backlog of known risks [REQ:DEBT_TRIAGE] [IMPL:DEBT_TRACKING]
- [ ] Add TODOs with owners in hotspot files [REQ:DEBT_TRIAGE] [IMPL:DEBT_TRACKING]
- [ ] Link debt list into docs/tasks [REQ:DEBT_TRIAGE]
- [ ] Run `[PROC:TOKEN_AUDIT]` + `[PROC:TOKEN_VALIDATION]`

**Completion Criteria**:
- [ ] Debt backlog documented and linked
- [ ] TODOs annotated in code
- [ ] Token audit + validation recorded

**Priority Rationale**: P1 to surface risks before refactors.


