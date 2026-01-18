# Implementation Decisions

**STDD Methodology Version**: 1.2.0

## Overview

This document serves as the **index** for all implementation decisions in this project. Each decision has been migrated to a dedicated detail file for better organization, navigation, and AI processing.

All decisions are cross-referenced with architecture decisions using `[ARCH:*]` tokens and requirements using `[REQ:*]` tokens for traceability.

## Directory Structure

```
stdd/
  implementation-decisions.md    # This index
  implementation-decisions/      # Detail files
    IMPL-*.md
```

## Filename Convention

Pattern: `IMPL-{TOKEN_NAME}.md`

Example: `[IMPL:CONFIG_STRUCT]` → `implementation-decisions/IMPL-CONFIG_STRUCT.md`

## Notes

- All implementation decisions MUST be recorded here IMMEDIATELY when made
- Each decision MUST include `[IMPL:*]` token and cross-reference both `[ARCH:*]` and `[REQ:*]` tokens
- Implementation decisions are dependent on both architecture decisions and requirements
- DO NOT defer implementation documentation - record decisions as they are made
- Record where code/tests are annotated so `[PROC:TOKEN_AUDIT]` can succeed later
- Include the most recent `[PROC:TOKEN_VALIDATION]` run information so future contributors know the last verified state
- **Language-Specific Implementation**: Language-specific implementation details (APIs, libraries, syntax patterns, idioms) belong in implementation decisions

## How to Add a New Implementation Decision

1. Create `IMPL-{TOKEN}.md` in `implementation-decisions/` using the template below
2. Add entry to the index table in this file
3. Update `semantic-tokens.md` registry with the new token

---

## Implementation Decisions Index

| Token | Title | Status | Cross-References | Detail |
|-------|-------|--------|------------------|--------|
| `[IMPL:CONFIG_STRUCT]` | Configuration Structure | Active | [ARCH:CONFIG_STRUCTURE] [REQ:CONFIGURATION] | [Detail](implementation-decisions/IMPL-CONFIG_STRUCT.md) |
| `[IMPL:STDD_FILES]` | STDD File Creation | Active | [ARCH:STDD_STRUCTURE] [REQ:STDD_SETUP] | [Detail](implementation-decisions/IMPL-STDD_FILES.md) |
| `[IMPL:STATE_PATH_RESOLVER]` | State Path Resolver | Active | [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS] | [Detail](implementation-decisions/IMPL-STATE_PATH_RESOLVER.md) |
| `[IMPL:EXTERNAL_COMMAND_LOADER]` | External Command Loader | Active | [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG] | [Detail](implementation-decisions/IMPL-EXTERNAL_COMMAND_LOADER.md) |
| `[IMPL:EXTERNAL_COMMAND_BINDER]` | External Command Binder | Active | [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG] | [Detail](implementation-decisions/IMPL-EXTERNAL_COMMAND_BINDER.md) |
| `[IMPL:EXTERNAL_COMMAND_APPEND]` | External Command Append Toggle | Active | [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG] | [Detail](implementation-decisions/IMPL-EXTERNAL_COMMAND_APPEND.md) |
| `[IMPL:ERROR_HANDLING]` | Error Handling Implementation | Active | [ARCH:ERROR_HANDLING] [REQ:ERROR_HANDLING] | [Detail](implementation-decisions/IMPL-ERROR_HANDLING.md) |
| `[IMPL:TESTING]` | Testing Implementation | Active | [ARCH:TESTING_STRATEGY] [REQ:*] | [Detail](implementation-decisions/IMPL-TESTING.md) |
| `[IMPL:CODE_STYLE]` | Code Style and Conventions | Active | — | [Detail](implementation-decisions/IMPL-CODE_STYLE.md) |
| `[IMPL:MODULE_VALIDATION]` | Module Validation Implementation | Active | [ARCH:MODULE_VALIDATION] [REQ:MODULE_VALIDATION] | [Detail](implementation-decisions/IMPL-MODULE_VALIDATION.md) |
| `[IMPL:GO_MOD_UPDATE]` | Go Mod Update | Active | [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS] | [Detail](implementation-decisions/IMPL-GO_MOD_UPDATE.md) |
| `[IMPL:DEP_BUMP]` | Dependency Bump | Active | [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH] | [Detail](implementation-decisions/IMPL-DEP_BUMP.md) |
| `[IMPL:CI_WORKFLOW]` | CI Workflow | Active | [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE] | [Detail](implementation-decisions/IMPL-CI_WORKFLOW.md) |
| `[IMPL:STATICCHECK_SETUP]` | Staticcheck Setup | Active | [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS] | [Detail](implementation-decisions/IMPL-STATICCHECK_SETUP.md) |
| `[IMPL:RACE_JOB]` | Race Job | Active | [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING] | [Detail](implementation-decisions/IMPL-RACE_JOB.md) |
| `[IMPL:TEST_WIDGETS]` | Widget Tests | Active | [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS] | [Detail](implementation-decisions/IMPL-TEST_WIDGETS.md) |
| `[IMPL:TEST_CMDLINE]` | Command Tests | Active | [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS] | [Detail](implementation-decisions/IMPL-TEST_CMDLINE.md) |
| `[IMPL:TEST_INTEGRATION_FLOWS]` | Integration Flow Tests | Active | [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS] | [Detail](implementation-decisions/IMPL-TEST_INTEGRATION_FLOWS.md) |
| `[IMPL:DOC_ARCH_GUIDE]` | Architecture Guide | Active | [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION] | [Detail](implementation-decisions/IMPL-DOC_ARCH_GUIDE.md) |
| `[IMPL:DOC_CONTRIBUTING]` | CONTRIBUTING Guide | Active | [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE] | [Detail](implementation-decisions/IMPL-DOC_CONTRIBUTING.md) |
| `[IMPL:MAKE_RELEASE_TARGETS]` | Release Targets | Active | [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX] | [Detail](implementation-decisions/IMPL-MAKE_RELEASE_TARGETS.md) |
| `[IMPL:BASELINE_SNAPSHOTS]` | Baseline Snapshots | Active | [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE] | [Detail](implementation-decisions/IMPL-BASELINE_SNAPSHOTS.md) |
| `[IMPL:WINDOW_MACRO_ENUMERATION]` | Window Macro Enumeration | Active | [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION] | [Detail](implementation-decisions/IMPL-WINDOW_MACRO_ENUMERATION.md) |
| `[IMPL:DEBT_TRACKING]` | Debt Tracking | Active | [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE] | [Detail](implementation-decisions/IMPL-DEBT_TRACKING.md) |
| `[IMPL:TOKEN_VALIDATION_SCRIPT]` | Token Validation Script | Active | [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP] | [Detail](implementation-decisions/IMPL-TOKEN_VALIDATION_SCRIPT.md) |
| `[IMPL:QUIT_DIALOG_ENTER]` | Quit Dialog Return Handling | Active | [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT] | [Detail](implementation-decisions/IMPL-QUIT_DIALOG_ENTER.md) |
| `[IMPL:BACKSPACE_TRANSLATION]` | Backspace Key Translation | Active | [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR] | [Detail](implementation-decisions/IMPL-BACKSPACE_TRANSLATION.md) |
| `[IMPL:TERMINAL_ADAPTER]` | Terminal Adapter Module | Active | [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD] | [Detail](implementation-decisions/IMPL-TERMINAL_ADAPTER.md) |
| `[IMPL:EVENT_LOOP_SHUTDOWN]` | Event Loop Shutdown Controller | Active | [ARCH:EVENT_LOOP_SHUTDOWN] [REQ:EVENT_LOOP_SHUTDOWN] | [Detail](implementation-decisions/IMPL-EVENT_LOOP_SHUTDOWN.md) |
| `[IMPL:XFORM_CLI_SCRIPT]` | Xform CLI Script | Active | [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING] | [Detail](implementation-decisions/IMPL-XFORM_CLI_SCRIPT.md) |
| `[IMPL:WORKSPACE_START_DIRS]` | Startup Directory Parser & Seeder | Active | [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS] | [Detail](implementation-decisions/IMPL-WORKSPACE_START_DIRS.md) |
| `[IMPL:FILER_EXCLUDE_RULES]` | Filename Exclude Rules | Active | [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES] | [Detail](implementation-decisions/IMPL-FILER_EXCLUDE_RULES.md) |
| `[IMPL:FILER_EXCLUDE_LOADER]` | Filename Exclude Loader & Toggle | Active | [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES] | [Detail](implementation-decisions/IMPL-FILER_EXCLUDE_LOADER.md) |
| `[IMPL:COMPARE_COLOR_CONFIG]` | Comparison Color Configuration | Active | [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS] | [Detail](implementation-decisions/IMPL-COMPARE_COLOR_CONFIG.md) |
| `[IMPL:FILE_COMPARISON_INDEX]` | File Comparison Index | Active | [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS] | [Detail](implementation-decisions/IMPL-FILE_COMPARISON_INDEX.md) |
| `[IMPL:COMPARISON_DRAW]` | Comparison Draw Integration | Active | [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS] | [Detail](implementation-decisions/IMPL-COMPARISON_DRAW.md) |
| `[IMPL:DIGEST_COMPARISON]` | Digest Comparison | Active | [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS] | [Detail](implementation-decisions/IMPL-DIGEST_COMPARISON.md) |
| `[IMPL:LINKED_NAVIGATION]` | Linked Navigation Implementation | Active | [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION] | [Detail](implementation-decisions/IMPL-LINKED_NAVIGATION.md) |
| `[IMPL:LINKED_NAVIGATION_AUTO_DISABLE]` | Linked Navigation Auto-Disable | Active | [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION] | [Detail](implementation-decisions/IMPL-LINKED_NAVIGATION_AUTO_DISABLE.md) |
| `[IMPL:DIFF_SEARCH]` | Difference Search Implementation | Active | [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH] | [Detail](implementation-decisions/IMPL-DIFF_SEARCH.md) |
| `[IMPL:NSYNC_OBSERVER]` | nsync Observer Adapter | Active | [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET] | [Detail](implementation-decisions/IMPL-NSYNC_OBSERVER.md) |
| `[IMPL:NSYNC_COPY_MOVE]` | nsync Copy/Move Wrappers | Active | [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET] | [Detail](implementation-decisions/IMPL-NSYNC_COPY_MOVE.md) |
| `[IMPL:NSYNC_CONFIRMATION]` | nsync Confirmation Modes | Active | [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION] | [Detail](implementation-decisions/IMPL-NSYNC_CONFIRMATION.md) |
| `[IMPL:HELP_POPUP]` | Help Popup Implementation | Active | [ARCH:HELP_WIDGET] [REQ:HELP_POPUP] | [Detail](implementation-decisions/IMPL-HELP_POPUP.md) |
| `[IMPL:SYNC_EXECUTE]` | Sync Execute | Active | [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS] | [Detail](implementation-decisions/IMPL-SYNC_EXECUTE.md) |
| `[IMPL:MOUSE_HIT_TEST]` | Mouse Hit Testing | Active | [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT] | [Detail](implementation-decisions/IMPL-MOUSE_HIT_TEST.md) |
| `[IMPL:MOUSE_FILE_SELECT]` | Mouse File Selection | Active | [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT] | [Detail](implementation-decisions/IMPL-MOUSE_FILE_SELECT.md) |
| `[IMPL:MOUSE_DOUBLE_CLICK]` | Mouse Double-Click Detection | Active | [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK] | [Detail](implementation-decisions/IMPL-MOUSE_DOUBLE_CLICK.md) |
| `[IMPL:MOUSE_CROSS_WINDOW_SYNC]` | Mouse Cross-Window Cursor Sync | Active | [ARCH:MOUSE_CROSS_WINDOW_SYNC] [REQ:MOUSE_CROSS_WINDOW_SYNC] | [Detail](implementation-decisions/IMPL-MOUSE_CROSS_WINDOW_SYNC.md) |
| `[IMPL:TOOLBAR_PARENT_BUTTON]` | Toolbar Parent Button | Active | [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON] | [Detail](implementation-decisions/IMPL-TOOLBAR_PARENT_BUTTON.md) |
| `[IMPL:TOOLBAR_LINKED_TOGGLE]` | Toolbar Linked Toggle Button | Active | [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE] | [Detail](implementation-decisions/IMPL-TOOLBAR_LINKED_TOGGLE.md) |
| `[IMPL:TOOLBAR_COMPARE_BUTTON]` | Toolbar Compare Digest Button | Active | [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON] | [Detail](implementation-decisions/IMPL-TOOLBAR_COMPARE_BUTTON.md) |
| `[IMPL:TOOLBAR_SYNC_COPY]` | Toolbar Sync Copy Button | Active | [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS] | [Detail](implementation-decisions/IMPL-TOOLBAR_SYNC_COPY.md) |
| `[IMPL:TOOLBAR_SYNC_DELETE]` | Toolbar Sync Delete Button | Active | [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS] | [Detail](implementation-decisions/IMPL-TOOLBAR_SYNC_DELETE.md) |
| `[IMPL:TOOLBAR_SYNC_RENAME]` | Toolbar Sync Rename Button | Active | [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS] | [Detail](implementation-decisions/IMPL-TOOLBAR_SYNC_RENAME.md) |
| `[IMPL:TOOLBAR_IGNORE_FAILURES]` | Toolbar Ignore Failures Toggle | Active | [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS] | [Detail](implementation-decisions/IMPL-TOOLBAR_IGNORE_FAILURES.md) |
| `[IMPL:LINKED_CURSOR_SYNC]` | Linked Cursor Synchronization | Active | [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION] | [Detail](implementation-decisions/IMPL-LINKED_CURSOR_SYNC.md) |
| `[IMPL:BATCH_DIFF_REPORT]` | Batch Diff Report CLI | Active | [ARCH:BATCH_DIFF_REPORT] [REQ:BATCH_DIFF_REPORT] | [Detail](implementation-decisions/IMPL-BATCH_DIFF_REPORT.md) |
| `[IMPL:HELP_STYLING]` | Help Popup Styling and Mouse Scroll | Active | [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING] | [Detail](implementation-decisions/IMPL-HELP_STYLING.md) |

### Status Values

- **Active**: Current implementation
- **Deprecated**: Superseded by newer decision
- **Proposed**: Under consideration

---

## Detail File Template

When creating a new implementation decision detail file, use this template:

```markdown
# [IMPL:TOKEN_NAME] Title

**Cross-References**: [ARCH:*] [REQ:*]  
**Status**: Active  
**Created**: YYYY-MM-DD  
**Last Updated**: YYYY-MM-DD

---

## Decision

Brief description of what was decided.

## Rationale

- Why this implementation approach was chosen
- What problems it solves
- How it fulfills the architecture decision

## Implementation Approach

- Specific technical details
- Code structure or patterns
- API design decisions

## Code Markers

Specific code locations, function names, or patterns to look for.

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] file/function - token

Tests that must reference `[REQ:*]`:
- [ ] TestName

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| YYYY-MM-DD | hash | ✅ Pass / ⏳ Pending | Details |

## Related Decisions

- Depends on: ...
- See also: ...

---

*Created on YYYY-MM-DD*
```

---

*Migrated to scalable index structure on 2026-01-17*
