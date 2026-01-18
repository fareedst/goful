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

### `[PROC:DEMO_GENERATION]`
- **Purpose** — Provide a reproducible workflow for creating terminal demo recordings (GIFs) that showcase goful features for documentation and README.
- **Scope** — Applies whenever new features need visual documentation or existing demos need updating.
- **Token references** — `[REQ:BATCH_DIFF_REPORT]`, `[REQ:FILE_COMPARISON_COLORS]`, `[REQ:LINKED_NAVIGATION]`, `[IMPL:DIGEST_COMPARISON]`.
- **Status** — Active.

#### Core Activities
1. **Prepare Environment**
   - Install prerequisites: `asciinema` (terminal recording) and `agg` (GIF conversion).
   - Build the goful binary: `go build -o bin/goful .`
   - Set terminal size to 120x35 columns/rows (matching existing demos).
   - Create demo directories with appropriate test files that demonstrate the feature.

2. **Create Demo Scripts**
   - For CLI-only demos (e.g., `--diff-report`): Create a bash script that runs the commands with appropriate pauses (`sleep`) for readability.
   - For interactive TUI demos: Create a bash wrapper script using tmux (see Working Configuration below).

3. **Record Demo** (see Working Configuration for the correct approach)

4. **Convert to GIF**
   - Use agg with asciinema theme to preserve original terminal colors:
     ```bash
     agg --theme asciinema /tmp/demo_name.cast .github/demo_name.gif
     ```

5. **Update Documentation**
   - Add GIF embed to README.md in the relevant feature section: `![demo_name](.github/demo_name.gif)`
   - Verify the GIF renders correctly in markdown preview.

6. **Requirements Stewardship**
   - Record demo creation in `stdd/tasks.md` if part of a larger feature task.
   - Update this process if new tools or techniques are discovered.

#### Artifacts & Metrics
- **Artifacts** — `.cast` recording files (temporary), `.gif` demo files in `.github/`, updated README.md with embeds.
- **Success Metrics** — GIF plays correctly, shows the intended feature clearly, file size is reasonable (<1MB), and README embeds display properly.

---

#### Working Configuration (Colors + Timing)

The following approach successfully captures both **colors** and **intermediate timing** for interactive TUI demos:

**Requirements:**
- tmux (for proper keystroke timing via `send-keys`)
- asciinema (for terminal recording)
- agg (for GIF conversion)

**Critical Environment Settings:**
```bash
unset NO_COLOR                    # CRITICAL: tcell checks this and disables all colors if set
export TERM=xterm-256color        # Required for tcell color detection
export COLORTERM=truecolor        # Enables 24-bit color support
```

**Recording Script Template:**
```bash
#!/bin/bash
SESSION="demo_name"
CAST_FILE="/tmp/demo_name.cast"

# Kill any existing session
tmux kill-session -t $SESSION 2>/dev/null

# Create new tmux session with specific dimensions
tmux new-session -d -s $SESSION -x 120 -y 35

# Start asciinema INSIDE tmux with proper environment
# CRITICAL: unset NO_COLOR and set TERM/COLORTERM before asciinema starts
tmux send-keys -t $SESSION "unset NO_COLOR && export TERM=xterm-256color COLORTERM=truecolor && asciinema rec --overwrite --cols 120 --rows 35 $CAST_FILE" Enter
sleep 3

# Start goful
tmux send-keys -t $SESSION "/path/to/goful /tmp/demo/dir1 /tmp/demo/dir2" Enter
sleep 4

# Send keystrokes with pauses (timing is preserved by tmux send-keys)
tmux send-keys -t $SESSION "j"
sleep 2
tmux send-keys -t $SESSION "="
sleep 3
# ... more keystrokes ...

# Quit goful
tmux send-keys -t $SESSION "q"
sleep 1
tmux send-keys -t $SESSION "y"
sleep 2

# Exit asciinema recording
tmux send-keys -t $SESSION "exit" Enter
sleep 2

# Cleanup
tmux kill-session -t $SESSION 2>/dev/null
```

**Why This Works:**
1. **tmux send-keys** sends keystrokes to the PTY with real-time execution, so `sleep` commands in the wrapper script create actual pauses that asciinema captures.
2. **asciinema runs inside tmux** where the environment variables are properly set before goful starts.
3. **NO_COLOR is unset** inside the tmux session, not in the outer shell, ensuring tcell sees the correct environment.

---

#### Approaches That Did NOT Work

The following configurations were tested and failed to capture both colors and timing correctly:

**1. asciinema + expect (no tmux)**
```bash
# Configuration:
export TERM=xterm-256color
export COLORTERM=truecolor
asciinema rec -c "expect demo.exp" output.cast
```
- **Result:** Colors worked (after unsetting NO_COLOR), but **timing was lost**.
- **Cause:** expect's `sleep` commands pause the script but don't flush output. All UI changes are buffered and appear at once in the recording.
- **Symptoms:** Recording shows goful UI appearing instantly with quit prompt visible immediately; no intermediate states.

**2. expect with output waiting**
```bash
# In expect script:
send "j"
expect -re "."   # Wait for any output
sleep 2
```
- **Result:** Same as above—timing still lost.
- **Cause:** `expect -re "."` matches output but doesn't force synchronous display.

**3. asciinema + stdbuf**
```bash
stdbuf -o0 expect demo.exp  # Disable output buffering
```
- **Result:** Did not improve timing capture.
- **Cause:** stdbuf affects stdio buffering but not PTY buffering.

**4. tmux without unsetting NO_COLOR inside the session**
```bash
# In wrapper script (outer shell):
unset NO_COLOR
export TERM=xterm-256color
# Then start tmux session...
tmux send-keys "asciinema rec..." Enter
```
- **Result:** **Timing worked but colors were black/white**.
- **Cause:** tmux sessions don't inherit the parent shell's unset variables. NO_COLOR remained set inside the tmux session (from user's shell profile).
- **Symptoms:** tcell's `Colors()` method returns 0; only bold/reverse video codes in recording, no color codes.

**5. Setting TERM before tmux session creation**
```bash
TERM=xterm-256color tmux new-session -d -s $SESSION
tmux set-environment -t $SESSION TERM xterm-256color
```
- **Result:** Colors still not working.
- **Cause:** Environment variables set via `tmux set-environment` affect new windows but asciinema was already running. Also, NO_COLOR was still the root cause.

**6. TCELL_TRUECOLOR=enable**
```bash
export TCELL_TRUECOLOR=enable
```
- **Result:** Did not help when NO_COLOR was set.
- **Cause:** tcell checks NO_COLOR first in its `Colors()` method. If NO_COLOR is set (to any value), it returns 0 immediately, before checking TCELL_TRUECOLOR or TERM.

---

#### Root Cause Analysis: NO_COLOR Environment Variable

The primary blocker for color recording was the `NO_COLOR` environment variable. This is a [standard](https://no-color.org/) that many CLI tools honor to disable colored output.

**tcell's behavior** (from `tscreen.go`):
```go
func (t *tScreen) Colors() int {
    if os.Getenv("NO_COLOR") != "" {
        return 0  // Returns 0 if NO_COLOR is set to ANY value
    }
    // ... rest of color detection
}
```

**Diagnosis steps used:**
1. Created a test program to check `screen.Colors()` inside asciinema—returned 0.
2. Checked tcell's terminfo lookup—correctly returned 256 colors.
3. Examined tcell source code and found the NO_COLOR check.
4. Verified `echo $NO_COLOR` showed `NO_COLOR=1` in the environment.
5. Unset NO_COLOR and re-tested—`Colors()` returned 16777216 (truecolor).

---

#### Troubleshooting Quick Reference

| Symptom | Cause | Fix |
|---------|-------|-----|
| Black/white output, no colors | `NO_COLOR` environment variable is set | `unset NO_COLOR` inside tmux session before starting asciinema |
| Colors work but timing is lost | Using expect directly with asciinema | Use tmux with `send-keys` instead of expect |
| "terminal not cursor addressable" panic | TERM not set or set to "dumb" | Set `TERM=xterm-256color` before launching goful |
| tcell `Colors()` returns 0 | NO_COLOR set, or TERM not recognized | Check and unset NO_COLOR; verify TERM has valid terminfo |
| agg converts to wrong colors | Using themed palette (e.g., monokai) | Use `--theme asciinema` to preserve original colors |

---

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
