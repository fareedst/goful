# STDD Processes

**STDD Methodology Version**: 1.1.0

Process documentation is the missing link that keeps tooling, rituals, and expectations traceable back to requirements. This guide defines how to record repeatable processes with semantic tokens so that every operational step you take is measurable, auditable, and associated with the intent that drove it.

## Process Tokens

Introduce `[PROC:*]` tokens whenever you describe how work happens.
Each token declares the process, its scope, and the requirements it serves. Because processes often span multiple artifacts, each entry should refer to:

- **Requirements** (`[REQ:*]`) to show whose intent the process satisfies
- **Architecture** (`[ARCH:*]`) or **Implementation** (`[IMPL:*]`) decisions that depend on the process outcome
- **Tests** (`[TEST:*]`) or other validation steps triggered by the process

Process entries become first-class trace nodes that explain **how** to survey, build, test, deploy, and otherwise steward the requirements themselves.

## Process Entry Template

Use the structure below for every process you document. Each entry should be kept current, reference the controlling requirements, and mention the deliverables or artifacts it produces.

### `[PROC:PROCESS_NAME]`
- **Purpose** — Describe the problem or requirement this process satisfies, ideally referencing a `[REQ:*]` token.
- **Scope** — Describe the boundaries of the process (teams, code areas, environments, or lifecycle phases).
- **Token references** — List `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`, or `[TEST:*]` tokens that the process continuously touches.
- **Status** — Active, deprecated, or scheduled for automation.

#### Core Activities
1. **Survey the Project**
   - Identify the existing intent (documentation, tokens, diagrams) tied to the requirement.
   - Capture discovery artifacts (notes, system maps, dependency lists) labeled with `[PROC:PROJECT_SURVEY]` or a more specific process token.
2. **Build Work**
   - Describe how to prepare the build environment, dependencies, and packages.
   - Reference architecture or implementation tokens that the process must observe before running the build.
3. **Test Work**
   - List the mandatory validation suites, acceptance tests, or checkpoints.
   - Include examples of test names that reference the requirement token (e.g., `TestFoo_REQ_BAR`).
4. **Deploy Work**
   - Outline the deployment targets, release artifacts, and approvals required.
   - Mention any CI/CD pipelines or configuration tokens that guarantee traceability.
5. **Requirements Stewardship**
   - State how the process collects feedback, updates requirements, and revalidates tokens.
   - Explain how this process keeps the `[REQ:*]` token fresh (review cadence, stakeholders, reporting).

#### Artifacts & Metrics
- **Artifacts** — Document the files, checklists, or dashboards produced during the process.
- **Success Metrics** — Name how you know the process satisfied the requirement (e.g., updated token table, green builds, automated audits).

### Example: `[PROC:PROJECT_SURVEY_AND_SETUP]`
- **Purpose** — Capture the context for `[REQ:STDD_SETUP]` before any new feature work.
- **Scope** — Applied to every new module or team onboarding cycle.
- **Token references** — `[REQ:STDD_SETUP]`, `[ARCH:STDD_STRUCTURE]`, `[IMPL:STDD_FILES]`.
- **Status** — Active.

#### Core Activities
1. **Survey**
   - Read `STDD.md`, `semantic-tokens.md`, and recent requirements to understand intent.
   - Tag findings with `[PROC:PROJECT_SURVEY_AND_SETUP]` and record them in the project knowledge base.
2. **Build**
   - Confirm required toolchains (language runtime, STDD tooling) are installed and share the list on the onboarding checklist.
   - Validate any `[ARCH:*]` constraints (folder layout, manifests) before manipulating files.
3. **Test**
   - Run smoke tests that include `[REQ:MODULE_VALIDATION]` to prove tracing works for a new module.
   - Check that tokens surfaced during survey show up in test names and code comments.
4. **Deploy**
   - Ensure deployment documentation references the same requirement tokens and that automated jobs run at least once to prove the configuration.
5. **Requirements Stewardship**
   - Record missing `[REQ:*]` tokens discovered during the survey and assign owners to author them.
   - Tag conclusions in the knowledge base with the `[PROC:PROJECT_SURVEY_AND_SETUP]` token so future reviews can trace the reasoning.

#### Artifacts & Metrics
- **Artifacts** — Onboarding checklist, environment matrix, token discovery log.
- **Success Metrics** — Every new module has `[REQ:*]` tokens defined, token registry updated, and build/test/deploy pipelines run at least once.

---

### `[PROC:TOKEN_AUDIT]`
- **Purpose** — Guarantee every change preserves the requirement → architecture → implementation trace mandated by `[REQ:STDD_SETUP]`.
- **Scope** — Mandatory for all code, test, and documentation edits that reference semantic tokens.
- **Token references** — `[REQ:STDD_SETUP]`, `[ARCH:STDD_STRUCTURE]`, `[ARCH:MODULE_VALIDATION]`, `[ARCH:TOKEN_VALIDATION_AUTOMATION]`, `[IMPL:TOKEN_VALIDATION_SCRIPT]`.
- **Status** — Active.

#### Core Activities
1. **Survey the Project**
   - Identify the requirement token(s) driving the change.
   - Locate related `[ARCH:*]` / `[IMPL:*]` entries and confirm they exist or create them before coding.
2. **Build Work**
   - Annotate every touched source file with the appropriate `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]` trio.
   - Update `semantic-tokens.md` when introducing new identifiers.
3. **Test Work**
   - Ensure each new/updated test name and comment references the requirement token.
   - Cross-check module validation evidence when tests fulfill `[REQ:MODULE_VALIDATION]`.
4. **Deploy Work**
   - Capture audit notes inside the relevant task in `stdd/tasks.md`, referencing files touched.
5. **Requirements Stewardship**
   - Record gaps or drift (e.g., missing architecture decisions) and assign owners to close them before code review.

#### Artifacts & Metrics
- **Artifacts** — Task log entries citing `[PROC:TOKEN_AUDIT]`, updated documentation sections, commit annotations.
- **Success Metrics** — Zero files lacking the required token breadcrumbs; every task lists the audit date and affected tokens.

---

### `[PROC:TOKEN_VALIDATION]`
- **Purpose** — Provide an automated, reproducible proof that every token used in source files exists in `stdd/semantic-tokens.md`.
- **Scope** — Run for every task prior to completion and whenever CI needs to assert registry integrity.
- **Token references** — `[REQ:STDD_SETUP]`, `[ARCH:TOKEN_VALIDATION_AUTOMATION]`, `[IMPL:TOKEN_VALIDATION_SCRIPT]`.
- **Status** — Active (automated via `scripts/validate_tokens.sh`).

#### Core Activities
1. **Survey the Project**
   - Confirm the token registry is updated with any new identifiers discovered during `[PROC:TOKEN_AUDIT]`.
   - Verify `ripgrep` (`rg`) and `git` are installed—the script depends on both.
2. **Build Work**
   - From the repo root, run `./scripts/validate_tokens.sh` (optionally pass explicit paths to scan docs/templates).
   - For PR automation, add a CI step invoking the same script so failures block merges.
3. **Test Work**
   - Inspect the script output; successful runs emit `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified <count> token references...`.
   - On failure, follow the printed list of missing tokens, update the registry, and rerun until the script succeeds.
4. **Deploy Work**
   - Capture the command output (or CI job link) inside `stdd/implementation-decisions.md` and the active task to provide traceable evidence.
5. **Requirements Stewardship**
   - When new directories begin carrying semantic tokens (e.g., docs), expand the default file globs in the script or pass targeted paths to keep coverage comprehensive.

#### Artifacts & Metrics
- **Artifacts** — Script output recorded in tasks and implementation decisions, CI job logs referencing `[PROC:TOKEN_VALIDATION]`.
- **Success Metrics** — Latest run shows zero missing tokens; CI fails fast if drift occurs; every completed task links to a successful validation run.

### `[PROC:TERMINAL_VALIDATION]`
- **Purpose** — Provide a reproducible manual checklist that proves `[REQ:TERMINAL_PORTABILITY]` and `[REQ:TERMINAL_CWD]` remain satisfied after adapter changes.
- **Scope** — Applies to macOS Terminal.app, Linux desktop sessions (gnome-terminal or equivalent), and tmux environments touched by `[ARCH:TERMINAL_LAUNCHER]` and `[IMPL:TERMINAL_ADAPTER]`.
- **Token references** — `[REQ:TERMINAL_PORTABILITY]`, `[REQ:TERMINAL_CWD]`, `[ARCH:TERMINAL_LAUNCHER]`, `[IMPL:TERMINAL_ADAPTER]`.
- **Status** — Active.

#### Core Activities
1. **Prepare Environment**
   - Build or install the current goful binary and ensure `GOFUL_DEBUG_TERMINAL=1` so `DEBUG: [IMPL:TERMINAL_ADAPTER]` logs record branch decisions.
   - Stage a workspace directory whose path contains spaces (e.g., `~/Projects/Terminal Demo`) to verify quoting behaviour.
2. **macOS Terminal.app Validation**
   - Launch goful outside tmux/screen, focus the staged directory, press `:` to open the shell prompt, and enter `echo mac`.
   - Expect Terminal.app to open a new tab/window, log the AppleScript branch, run `cd "<focused dir>"; echo mac;read -p "HIT ENTER KEY"`, and keep the window open until Enter is pressed, after which it exits cleanly.
   - Repeat inside a tmux session; tmux should take precedence and open a new window without invoking Terminal.app.
   - Repeat with `GOFUL_TERMINAL_APP="iTerm2"` and/or `GOFUL_TERMINAL_SHELL="zsh"` set to confirm `[REQ:TERMINAL_PORTABILITY]` applies the new runtime parameters (logs should call out the selected app/shell).
3. **Linux Desktop Validation**
   - On a Linux desktop (no tmux), repeat the shell prompt action and confirm gnome-terminal launches with the title escape (`echo -n '\033]0;cmd\007'`) before running the payload and pause prompt.
   - Set `GOFUL_TERMINAL_CMD="alacritty -e"`, rerun the step, and confirm the override command receives the payload and pause tail while inheriting the focused-directory working dir when GOOS=darwin.
4. **tmux-only Validation**
   - With `TERM`/`TERM_PROGRAM` indicating tmux or screen, trigger the terminal command and confirm `tmux new-window -n <cmd>` is used regardless of OS.
5. **Documentation & Stewardship**
   - Record observed behaviour, attached logs, and any discrepancies inside `stdd/tasks.md` under the active terminal launcher task.
   - Update `[REQ:TERMINAL_PORTABILITY]`/`[REQ:TERMINAL_CWD]` validation evidence if regressions or new findings surface.

#### Artifacts & Metrics
- **Artifacts** — Saved `DEBUG: [IMPL:TERMINAL_ADAPTER]` logs from macOS/Linux runs, operator notes linked in `stdd/tasks.md`.
- **Success Metrics** — macOS Terminal opens in the focused directory with the pause tail, Linux overrides behave predictably, and tmux detection always routes through `tmux new-window` when active.
