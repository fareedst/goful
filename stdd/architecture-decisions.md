# Architecture Decisions

**STDD Methodology Version**: 1.1.0

## Overview
This document captures the high-level architectural decisions for this project. All decisions are cross-referenced with requirements using semantic tokens `[REQ:*]` and assigned architecture tokens `[ARCH:*]` for traceability.

## Template Structure

When documenting architecture decisions, use this format:

```markdown
## N. Decision Title [ARCH:IDENTIFIER] [REQ:RELATED_REQUIREMENT]

### Decision: Brief description of the architectural decision
**Rationale:**
- Why this decision was made
- What problems it solves
- What benefits it provides

**Alternatives Considered:**
- Alternative 1: Why it was rejected
- Alternative 2: Why it was rejected

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code files expected to carry `[IMPL:*] [ARCH:*] [REQ:*]` comments
- Tests expected to reference `[REQ:*]` / `[TEST:*]` tokens that validate this decision

**Cross-References**: [REQ:RELATED_REQUIREMENT], [ARCH:OTHER_DECISION]
```

## Notes

- All architecture decisions MUST be recorded here IMMEDIATELY when made
- Each decision MUST include `[ARCH:*]` token and cross-reference `[REQ:*]` tokens
- Architecture decisions are dependent on requirements
- DO NOT defer architecture documentation - record decisions as they are made
- Document the expected code + test touchpoints so `[PROC:TOKEN_AUDIT]` has concrete files/functions to verify.
- Capture the intended validation tooling (e.g., references to `./scripts/validate_tokens.sh`) so `[PROC:TOKEN_VALIDATION]` remains reproducible.
- **Language Selection**: Language selection, runtime choices, and language-specific architectural patterns belong in architecture decisions. Document language choice with `[ARCH:LANGUAGE_SELECTION]` token when it's an architectural decision (not a requirement). Language-specific patterns (e.g., async/await, goroutines, callbacks) should be documented here. Requirements should remain language-agnostic unless language selection is itself a specific requirement.

---

**Rationale:**
- Clear separation of concerns
- Standard project layout
- Testable components

## 3. STDD Project Structure [ARCH:STDD_STRUCTURE] [REQ:STDD_SETUP]

### Decision: Centralized `stdd/` Directory
**Rationale:**
- Keeps documentation close to code but organized in a dedicated namespace.
- Ensures the AI agent can easily find all context in one place.
- Separates meta-documentation from project source code.

**Alternatives Considered:**
- Root-level files: Clutters the root directory.
- `.github` or `.docs` folder: `stdd` is more specific to the methodology.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: `// [IMPL:STDD_FILES] [ARCH:STDD_STRUCTURE] [REQ:STDD_SETUP]` comments in bootstrap scripts.
- Tests: `TestSTDDSetup_REQ_STDD_SETUP` ensures docs + registry exist.

**Cross-References**: [REQ:STDD_SETUP]

## 3. State Path Selection [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]

### Decision: Deterministic path precedence for state/history persistence
**Rationale:**
- Satisfies [REQ:CONFIGURABLE_STATE_PATHS] by letting operators override persistence files without editing source.
- Supports hermetic tests and sandboxes that must run multiple goful instances simultaneously.
- Keeps behavior explicit: CLI flags override environment settings, environment overrides defaults, and defaults remain backwards compatible.

**Alternatives Considered:**
- **Environment-only overrides**: rejected because CI and scripted invocations need per-run control without mutating process-wide env.
- **Config file**: rejected for this iteration; adds file parsing complexity without immediate requirement coverage.

**Implementation:**
- Introduce `configpaths.Resolver` (Module 1) that accepts CLI flag inputs + `LookupEnv` hook and returns expanded paths plus provenance metadata.
- Provide `BootstrapPaths` helper (Module 2) that applies the resolved paths to `app.NewGoful`, `filer.SaveState`, `cmdline.LoadHistory`, and `cmdline.SaveHistory`, emitting `DEBUG:` lines tagged with `[IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]`.
- Surface CLI flags `-state` and `-history`, along with environment variables `GOFUL_STATE_PATH`/`GOFUL_HISTORY_PATH`, to feed the resolver.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: Annotate resolver and bootstrap helpers with `[IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]`.
- Tests: `TestResolvePaths_REQ_CONFIGURABLE_STATE_PATHS` style names validate precedence logic and module validation evidence.

## 3b. External Command Registry [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]

### Decision: External command bindings load from a JSON or YAML registry with flag/env/default path precedence and deterministic menu wiring.
**Rationale:**
- Removes the need to patch `main.go` to customize the `external-command` menu; teams can distribute curated configs per environment.
- Mirrors the existing precedence contract so `-commands` (CLI) overrides `GOFUL_COMMANDS_FILE`, which overrides the default `~/.goful/external_commands.yaml`.
- Preserves current behavior (Windows vs. POSIX defaults, cursor offsets, macros) when no file is present or parsing fails, reducing regression risk.
- By default, inherits the compiled defaults even when a file is present, then **prepends** file entries so customized shortcuts appear first while legacy entries remain unless the file opts into a "replace defaults" mode.
- Surfaces `DEBUG:` diagnostics and `message.Errorf` warnings when entries are skipped (duplicate keys, missing fields, unsupported platforms).

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `CommandConfigPathResolver` extends `configpaths.Resolver` with `Commands` + provenance metadata and emits `[IMPL:STATE_PATH_RESOLVER]` debug output alongside state/history.
- `externalcmd.Loader` (`[IMPL:EXTERNAL_COMMAND_LOADER]` & `[IMPL:EXTERNAL_COMMAND_APPEND]`) expands `~`, reads JSON or YAML, validates schema, enforces unique `menu/key`, filters by GOOS, honors `disabled`, _merges defaults + file entries with prepend-as-default semantics_, and falls back to embedded defaults packaged per platform.
- `externalcmd.Defaults` enumerates the historical Windows/POSIX menus so regression tests can compare lists directly and gives the loader something to inherit from when files only supply deltas.
- `main.registerExternalCommandMenu` (`[IMPL:EXTERNAL_COMMAND_BINDER]`) converts validated entries into `menu.Add` triplets, captures cursor offsets when invoking `g.Shell`, and injects a placeholder item if configuration disables every command to keep `X` from erroring.

**Alternatives Considered:**
- **Ad-hoc parsing inside `main.go`:** rejected to keep loader logic testable and aligned with `[REQ:MODULE_VALIDATION]`.
- **Environment-only overrides:** rejected; operators asked for file-based registries they can check into dotfile repos.
- **TOML-only support:** rejected; JSON+YAML cover the main ecosystems with a single extra dependency (`gopkg.in/yaml.v3`).

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `configpaths/resolver.go` comments mention `[IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]` plus the new commands path metadata.
- `externalcmd/loader.go` and tests include `[IMPL:EXTERNAL_COMMAND_LOADER] [IMPL:EXTERNAL_COMMAND_APPEND] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]`.
- `main_external_commands.go` helper + tests include `[IMPL:EXTERNAL_COMMAND_BINDER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]`.

**Cross-References**: [REQ:EXTERNAL_COMMAND_CONFIG], [IMPL:EXTERNAL_COMMAND_LOADER], [IMPL:EXTERNAL_COMMAND_BINDER]

## 4. Data Management [ARCH:DATA_MANAGEMENT] [REQ:DATA_REQUIREMENT]

### Decision: [Your Data Management Approach]
**Rationale:**
- Reason 1
- Reason 2

**Implementation:**
- Storage approach
- Data access patterns
- Consistency model

## 5. Error Handling Strategy [ARCH:ERROR_HANDLING] [REQ:ERROR_HANDLING]

### Decision: [Your Error Handling Approach]
**Rationale:**
- Idiomatic for chosen language/framework
- Clear error propagation
- Easy to test

**Pattern:**
- Error types
- Error propagation
- Error reporting

## 6. Testing Strategy [ARCH:TESTING_STRATEGY]

### Decision: [Your Testing Approach]
**Rationale:**
- Comprehensive test coverage
- Fast unit tests
- Integration tests for end-to-end scenarios
- Aligns with validation criteria defined in requirements [REQ:*]

**Structure:**
- Unit test organization
- Integration test organization
- Test fixtures and utilities

**Note**: This testing strategy implements the validation criteria specified in `requirements.md`. Each requirement's validation criteria informs what types of tests are needed (unit, integration, manual verification, etc.).

## 7. Dependency Management [ARCH:DEPENDENCY_MANAGEMENT]

### Decision: [Your Dependency Management Approach]
**Rationale:**
- Reduce external dependencies
- Faster builds
- Fewer security concerns

**Allowed Dependencies:**
- Standard library only (or minimal external dependencies)
- Consider external packages only if standard library is insufficient

## 8. Build and Distribution [ARCH:BUILD_DISTRIBUTION]

### Decision: [Your Build and Distribution Approach]
**Rationale:**
- Easy deployment
- No runtime dependencies
- Cross-platform support

**Build Targets:**
- Platform 1
- Platform 2
- Platform 3

## 9. Code Organization Principles [ARCH:CODE_ORGANIZATION]

### Decision: [Your Code Organization Approach]
**Rationale:**
- Testable components
- Clear responsibilities
- Easy to extend
- Maintainable codebase

**Principles:**
- Each module has a single, clear responsibility
- Functions are small and focused
- Interfaces where appropriate for testability
- Avoid global state where possible

## 10. Module Validation Strategy [ARCH:MODULE_VALIDATION] [REQ:MODULE_VALIDATION]

### Decision: Independent Module Validation Before Integration
**Rationale:**
- Eliminates bugs related to code complexity by ensuring each module works correctly in isolation
- Reduces integration complexity by validating modules independently before combining them
- Catches bugs early in the development cycle, before integration issues compound
- Ensures each module meets its defined contract before integration
- Makes debugging easier by isolating issues to specific modules
- Enables parallel development of modules when dependencies are properly mocked

**Module Identification Requirements:**
- Modules must be identified and documented before development begins
- Each module must have clear boundaries and responsibilities
- Module interfaces and contracts must be defined and documented
- Module dependencies must be identified and documented
- Module validation criteria must be specified (what "validated" means for each module)

**Validation Approach:**
- **Unit Testing**: Each module must have comprehensive unit tests with mocked dependencies
- **Integration Testing with Test Doubles**: Modules must be tested with mocks, stubs, or fakes for dependencies
- **Contract Validation**: Input/output validation to ensure modules meet their defined contracts
- **Edge Case Testing**: Modules must be tested with edge cases and boundary conditions
- **Error Handling Validation**: Modules must be tested for proper error handling and error propagation

**Integration Requirements:**
- Integration tasks must be separate from module development and validation tasks
- Integration only occurs after module validation passes
- Integration tests validate the combined behavior of validated modules
- Module validation results must be documented before integration

**Alternatives Considered:**
- **Big Bang Integration**: Integrating all modules at once without independent validation
  - Rejected: Too complex, makes debugging difficult, bugs compound
- **Minimal Validation**: Only basic unit tests before integration
  - Rejected: Insufficient to catch complexity-related bugs, doesn't validate contracts properly
- **Post-Integration Validation Only**: Validating only after integration
  - Rejected: Doesn't catch module-level bugs early, increases debugging complexity

**Cross-References**: [REQ:MODULE_VALIDATION], [IMPL:MODULE_VALIDATION]

## 11. Go Runtime Strategy [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS]

### Decision: Track current Go LTS in `go.mod` and CI
**Rationale:**
- Ensures access to security patches and modern stdlib APIs.
- Aligns local and CI builds to avoid divergence.
- Provides predictable compiler behavior for race/static analysis.

**Alternatives Considered:**
- Staying on Go 1.16: rejected due to lack of support and security fixes.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: `go.mod` comment `// [IMPL:GO_MOD_UPDATE] [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS]`.
- CI: workflow step comments include same tokens.

**Cross-References**: [REQ:GO_TOOLCHAIN_LTS], [IMPL:GO_MOD_UPDATE]

## 12. Dependency Policy [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]

### Decision: Refresh direct deps to current stable versions with tidy
**Rationale:**
- Reduce CVE exposure and improve terminal compatibility.
- Keep transitive `x/*` libraries aligned with Go LTS.

**Alternatives Considered:**
- Pinning to legacy versions: rejected due to security/compatibility risk.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: dependency bumps carry `[IMPL:DEP_BUMP] [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]` comments.
- Tests: regression tests tagged with `[REQ:DEPENDENCY_REFRESH]` when shims added.

**Cross-References**: [REQ:DEPENDENCY_REFRESH], [IMPL:DEP_BUMP]

## 13. CI Pipeline [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]

### Decision: GitHub Actions workflow for fmt/vet/test with caching
**Rationale:**
- Enforces formatting and vetting before merge.
- Runs unit tests on pushes/PRs.
- Uses Go cache to keep runtimes fast.

**Alternatives Considered:**
- Local-only checks: rejected; lacks enforcement.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `.github/workflows/ci.yml` steps annotated with `[IMPL:CI_WORKFLOW] [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]`.
- Tests referenced in workflow include requirement tokens.

**Cross-References**: [REQ:CI_PIPELINE_CORE], [IMPL:CI_WORKFLOW]

## 14. Static Analysis Policy [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]

### Decision: Add staticcheck (and optional golangci-lint) job
**Rationale:**
- Catches API misuse and nil/loop issues early.
- Keeps exclusions explicit and minimal.

**Alternatives Considered:**
- Relying on vet only: rejected; weaker coverage.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Workflow job comments include `[IMPL:STATICCHECK_SETUP] [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]`.
- Config file (if added) carries same tokens.

**Cross-References**: [REQ:STATIC_ANALYSIS], [IMPL:STATICCHECK_SETUP]

## 15. Race Testing Pipeline [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]

### Decision: Dedicated `go test -race ./...` job
**Rationale:**
- Detects concurrency issues in app/widget event handling.
- Keeps runtime separate to manage resource needs.

**Alternatives Considered:**
- Folding into main job: rejected to keep runtimes predictable.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Workflow job annotated with `[IMPL:RACE_JOB] [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]`.

**Cross-References**: [REQ:RACE_TESTING], [IMPL:RACE_JOB]

## 16. UI Test Strategy [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]

### Decision: Unit/snapshot tests for widgets and filer primitives
**Rationale:**
- Protects rendering and event handling from regressions.
- Supports module validation before integration.

**Alternatives Considered:**
- Manual verification only: rejected; not repeatable.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Tests include `[IMPL:TEST_WIDGETS] [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]`.

**Cross-References**: [REQ:UI_PRIMITIVE_TESTS], [IMPL:TEST_WIDGETS]

## 17. Command Test Strategy [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]

### Decision: Focused tests for command parsing and mode transitions
**Rationale:**
- Ensures command line handling remains stable pre-refactor.

**Alternatives Considered:**
- Only integration tests: rejected; misses fast feedback.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Tests annotated with `[IMPL:TEST_CMDLINE] [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]`.

**Cross-References**: [REQ:CMD_HANDLER_TESTS], [IMPL:TEST_CMDLINE]

## 18. Integration Test Strategy [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]

### Decision: Snapshot/integration flows for open/navigate/rename/delete
**Rationale:**
- Validates core file operations end-to-end.

**Alternatives Considered:**
- Unit tests only: rejected; miss cross-module behavior.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Integration tests include `[IMPL:TEST_INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]`.

**Cross-References**: [REQ:INTEGRATION_FLOWS], [IMPL:TEST_INTEGRATION_FLOWS]

## 19. Docs Structure [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]

### Decision: Add `ARCHITECTURE.md` describing packages, event flow, and validation seams
**Rationale:**
- Provides concise onboarding and change impact map.
- Captures how `main` → `config` → `app.Goful` compose widgets so refactors can reason about ripple effects.
- Documents validation seams per [REQ:MODULE_VALIDATION] so test authors know which modules can be mocked.

**Structure & Scope:**
- **Overview** outlines user-facing goals and the relationship between CLI shells, filer panes, and widgets.  
- **Runtime Flow** traces `main.go` flag parsing, `configpaths.Resolver`, and the event loop between `widget.PollEvent`, `app.Goful`, and subordinate widgets.  
- **Module Deep Dives** cover `app`, `filer`, `widget`, `cmdline`, `menu`, `configpaths`, `message/info/progress`, and persistence helpers with inter-module contracts.  
- **Validation Hooks** cite existing test suites (`widget`, `cmdline`, `filer`, keymap baselines) and note extension points for future modules.

**Validation Plan**:
- `DocArchitecture` module is validated via doc review checklist ensuring every section links `[REQ:*] → [ARCH:*] → [IMPL:*]`.
- Cross-links in README + CONTRIBUTING confirm discoverability.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Document includes `[IMPL:DOC_ARCH_GUIDE] [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]` plus backlinks to related tokens (`[REQ:CONFIGURABLE_STATE_PATHS]`, `[ARCH:STATE_PATH_SELECTION]`, etc.).

**Cross-References**: [REQ:ARCH_DOCUMENTATION], [IMPL:DOC_ARCH_GUIDE], [REQ:MODULE_VALIDATION]

## 20. Contribution Process [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE]

### Decision: Contributor guide with coding standards and review flow
**Rationale:**
- Aligns expectations on formatting, testing, and tokens.
- Documents the enforced debug/logging policy so contributors do not strip diagnostic output required by STDD.
- Highlights module-validation steps per [REQ:MODULE_VALIDATION] to keep integration safe.

**Structure & Scope:**
- **Environment & Tooling** (Go 1.24.x, `make fmt/test`, scripts).
- **Workflow Checklist** (format → lint → test → `./scripts/validate_tokens.sh` with `[PROC:TOKEN_AUDIT]` references).
- **Semantic Token Discipline** (how to register new `[REQ:*]/[ARCH:*]/[IMPL:*]/[TEST:*]` entries).
- **Module Validation + Debug Expectations** spanning unit, integration, and baseline tests along with required `DEBUG:`/`DIAGNOSTIC:` prefixes.
- **Review Readiness** signals (CI, token validation logs, documentation cross-links).

**Validation Plan**:
- `DocContributing` module validated via doc review + cross-check with README/CI instructions.
- Test section references `KeymapBaselineSuite` + other modules to prove coverage mapping.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `CONTRIBUTING.md` carries `[IMPL:DOC_CONTRIBUTING] [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE]` plus `[PROC:TOKEN_AUDIT]`, `[PROC:TOKEN_VALIDATION]`, `[REQ:MODULE_VALIDATION]` references.

**Cross-References**: [REQ:CONTRIBUTING_GUIDE], [IMPL:DOC_CONTRIBUTING], [REQ:MODULE_VALIDATION]

## 21. Build Matrix [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]

### Decision: Makefile + CI/release matrix for reproducible CGO-disabled builds
**Rationale:**
- Ensures deterministically named binaries (goful_${GOOS}_${GOARCH}) plus SHA256 digests for linux/amd64, linux/arm64, and darwin/arm64 – the primary release targets today.
- Keeps the release process scriptable (single `make release` locally, `release-matrix` job in CI) so artifacts are identical regardless of where they are produced.
- Provides auditable logs (`DIAGNOSTIC: [IMPL:MAKE_RELEASE_TARGETS] ...`) that capture which platform was built and which checksum was emitted.

**Structure:**
- **Makefile** adds reusable targets (`lint`, `test`, `release`, `clean-release`) inside module `MakeReleaseTargets`. The release target loops over `RELEASE_PLATFORMS` or a supplied `PLATFORM`, sets `GOOS/GOARCH`, enforces `CGO_ENABLED=0`, uses `-trimpath -ldflags "-s -w"`, writes binaries to `dist/`, and immediately generates `.sha256` files.  
- **GitHub Actions CI** contains job `release-matrix` (module `ReleaseMatrixWorkflow`) with strategy include set:
  - linux/amd64
  - linux/arm64
  - darwin/arm64
  The job runs `make release PLATFORM=${{matrix.goos}}/${{matrix.goarch}}`, emits the checksum to logs, and uploads the binary + digest as artifacts.
- **Release workflow** (`.github/workflows/release.yml`) reuses the same matrix + Makefile target when tags `v*` are pushed, then aggregates the uploaded artifacts and publishes them to GitHub Releases via `softprops/action-gh-release`, guaranteeing the downloadable assets match CI output.
- **Artifact verification** (`ArtifactDeterminismAudit`) is satisfied by having Makefile, CI, and release workflows regenerate digests straight from the compiled binary so humans can diff `.sha256` outputs (and rerun `make release` to reproduce).

**Validation Plan [REQ:MODULE_VALIDATION]:**
- `MakeReleaseTargets` validated locally by running `make release PLATFORM=$(go env GOOS)/$(go env GOARCH)` and confirming `dist/` only contains the expected files plus `.sha256`.
- `ReleaseMatrixWorkflow` validation occurs through CI: matrix job must succeed and upload all artifacts; failure indicates platform-specific regression.
- `ArtifactDeterminismAudit` validation is the SHA output diff – rerunning `make release` should yield identical digests (documented when recording release notes).

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Makefile targets, CI job shell blocks, and README/CONTRIBUTING references include `[IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]`.

**Cross-References**: [REQ:RELEASE_BUILD_MATRIX], [IMPL:MAKE_RELEASE_TARGETS], [REQ:MODULE_VALIDATION]

## 22. Window Macro Enumeration [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]

### Decision: Deterministically enumerate workspace directories for `%D@`
**Rationale:**
- External command macros currently expose only the focused (`%D`) and next (`%D2`) directories, forcing multi-window workflows to retype remaining paths.
- Shell automations often move/copy between all open panes; providing a placeholder for “all other windows” lets operators script against the current workspace layout.
- Deterministic ordering (focused window followed by the next windows wrapping around) keeps automation predictable even as panes are added or rotated.

**Module Boundaries & Contracts** `[REQ:MODULE_VALIDATION]`:
- `WindowSequenceBuilder` (Module 1) – Pure helper that inspects `filer.Workspace` state and returns the ordered list of other directory paths. The function must not mutate focus or layout and must gracefully handle 1-window workspaces by returning an empty slice. A companion helper derives just the directory names (`Directory.Base()`) from the same deterministic ordering so `%d@` can reuse the sequence logic without copying path manipulation everywhere.
- `MacroListFormatter` (Module 2) – Formats the sequence for macro insertion, applying quoting rules (`util.Quote`) per entry when requested. `%D@` calls the quoted branch so every emitted path is shell safe, while `%~D@` explicitly opts into the raw branch (no escaping) to preserve the tilde modifier's "non-quote" semantics. `%d@` reuses the same formatter on the basename list so the quoting guarantees stay identical whether scripts need full paths or names. It joins entries with single spaces and returns an empty string for empty input.
- Integration point: `expandMacro` routes `%D@`, `%~D@`, `%d@`, and `%~d@` to these modules, so existing macro parsing (escapes, `%~~` guardrails, `%&`) continues to behave identically for other placeholders. The dispatcher now selects the quoted vs. raw formatter based on whether the tilde modifier was present.

**Pseudo-Code Sketch:**
```text
// [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]
func otherWindowPaths(ws *filer.Workspace) []string {
    paths := []string{}
    for offset := 1; offset < len(ws.Dirs); offset++ {
        idx := (ws.Focus + offset) % len(ws.Dirs)
        paths = append(paths, ws.Dirs[idx].Path)
    }
    return paths
}

func formatDirs(paths []string, quote bool) string {
    if len(paths) == 0 {
        return ""
    }
    parts := make([]string, len(paths))
    for i, p := range paths {
        parts[i] = chooseQuote(p, quote)
    }
    return strings.Join(parts, " ")
}
```

**Validation Plan:**
- Module tests instantiate lightweight workspaces with synthetic directories to verify ordering, wrap-around, and 1-window behavior.
- Integration tests extend `app/spawn_test.go` to assert `%D@` expansions stay quoted while `%~D@` returns raw paths (including cases with spaces) and to prove `%d@`/`%~d@` reuse the same ordering using only directory names. This keeps regression coverage tied back to the requirement even when the tilde modifier is present.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `app/spawn.go` comments for the new helper(s) and `%D@`/`%d@` branches carry `[IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]`.
- `app/spawn_test.go` test names/comments reference `[REQ:WINDOW_MACRO_ENUMERATION]`.
- `README.md` macro table documents `%D@` with matching tokens for discoverability.

**Cross-References**: [REQ:WINDOW_MACRO_ENUMERATION], [IMPL:WINDOW_MACRO_ENUMERATION], [REQ:MODULE_VALIDATION]

## 23. Baseline Capture [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE]

### Decision: Capture key interactions/keymaps as executable baselines
**Rationale:**
- Protects current behavior before refactors.
- Provides changelog-friendly evidence whenever a keybinding or mode is re-mapped.
- Supports `[REQ:MODULE_VALIDATION]` by giving pure tests that fail fast before wiring into terminal I/O.

**Module Boundaries & Coverage:**
- `KeymapBaselineSuite` (tests) snapshots default bindings for:
  - `filerKeymap` navigation/selection/command chords (`j/k`, `space`, `q/Q`, `:`, menu launches, etc.).
  - `cmdlineKeymap` editing, history, run/exit commands.
  - `finderKeymap`, `completionKeymap`, and `menuKeymap` cursor + exit bindings.
- Suite references `[TEST:KEYMAP_BASELINE]` plus `[REQ:BEHAVIOR_BASELINE]` to keep traceability.
- Baselines are intentionally pure (no widget drawing) so they run in CI without tcell initialization.

**Validation Plan:**
- Table-driven tests confirm required key strings exist in returned `widget.Keymap` maps and emit `DEBUG:` logs enumerating coverage.
- Future modules can extend baseline coverage by appending more required chords without touching runtime logic.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Baseline tests/scripts include `[IMPL:BASELINE_SNAPSHOTS] [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE] [TEST:KEYMAP_BASELINE]`.

**Cross-References**: [REQ:BEHAVIOR_BASELINE], [IMPL:BASELINE_SNAPSHOTS], [TEST:KEYMAP_BASELINE], [REQ:MODULE_VALIDATION]

## 24. Debt Management [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]

### Decision: Systematically log and annotate risky areas
**Rationale:**
- Makes hidden risks visible for scheduling and refactors.
- Ensures inline TODOs, backlog entries, and STDD docs share the same semantic tokens for auditability.
- Provides owners plus next steps so the debt backlog can be triaged like any other requirement-driven artifact.

**Structure:**
- `stdd/debt-log.md` is the canonical backlog that maps each risk (D1, D2, …) to owners, affected files, and mitigation plans with `[REQ:DEBT_TRIAGE] [IMPL:DEBT_TRACKING]` breadcrumbs.
- Runtime hotspots are annotated with `TODO(goful-maintainers)` comments that repeat `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`.
- `stdd/tasks.md` links back to the backlog so completion criteria reference concrete evidence.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Issues/TODOs carry `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`.
- Backlog doc references `[ARCH:DEBT_MANAGEMENT]` and links to specific files (e.g., `app/goful.go`, `cmdline/cmdline.go`) for each item.

**Cross-References**: [REQ:DEBT_TRIAGE], [IMPL:DEBT_TRACKING]

## 25. Token Validation Automation [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]

### Decision: Provide a helper script that enforces `[PROC:TOKEN_VALIDATION]`
**Rationale:**
- Ensures code/test references only use tokens registered in `stdd/semantic-tokens.md`.
- Gives contributors a repeatable command (`scripts/validate_tokens.sh`) to gate PRs.
- Unblocks tasks that require `[PROC:TOKEN_VALIDATION]` evidence (e.g., modernization work).

**Implementation Notes:**
- Default scope targets tracked source files (`*.go`, Go modules, shell scripts, workflows, Makefile) to avoid template placeholders.
- Additional paths can be passed to the script once templates are scrubbed of placeholder tokens.
- Script depends on `git` and `ripgrep`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `scripts/validate_tokens.sh` header/comment includes `[IMPL:TOKEN_VALIDATION_SCRIPT] [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]`.

**Cross-References**: [REQ:STDD_SETUP], [IMPL:TOKEN_VALIDATION_SCRIPT], [PROC:TOKEN_VALIDATION]

## 25. Quit Dialog Key Translation [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]

### Decision: Normalize Return/Enter events to the historical `C-m` submission path
**Rationale:**
- Recent tcell upgrades emit `KeyEnter` instead of `KeyCtrlM`, which broke the implicit default-confirm behavior in the quit dialog.
- Mapping Return/Enter to the same symbol used by the cmdline keymap keeps behavior stable without duplicating key bindings.
- Centralizing the mapping in `widget.EventToString` ensures every cmdline-based mode benefits, not just quit.

**Alternatives Considered:**
- Update every cmdline keymap to handle a new `enter` symbol: rejected because it is error-prone and duplicates logic across widgets.
- Special-case quit mode to detect empty text before submission: rejected because it bypasses the shared cmdline infrastructure and would not fix other dialogs.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: `widget/widget.go` `EventToString` includes comments `// [IMPL:QUIT_DIALOG_ENTER] [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]`.
- Tests: `widget/widget_test.go` adds `TestEventToStringReturnKey_REQ_QUIT_DIALOG_DEFAULT` (or equivalent) covering both `KeyEnter` and `KeyCtrlM`.

**Module Boundaries & Validation `[REQ:MODULE_VALIDATION]`:**
- `InputEventTranslator` module (`widget.EventToString`) translates tcell events → semantic strings; validated via new unit test.
- `CmdlineSubmit` module (`cmdline.Keymap` handlers for Run) already validated; integration relies on translator emitting `C-m`.

**Cross-References**: [REQ:QUIT_DIALOG_DEFAULT], [IMPL:QUIT_DIALOG_ENTER]

## 25b. Backspace Key Translation [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]

### Decision: Canonicalize Backspace key codes so every widget observes the same `backspace` symbol
**Rationale:**
- macOS keyboards label the Backspace key as “delete” and commonly emit `tcell.KeyBackspace`, while many Linux terminals emit `tcell.KeyBackspace2`. Without normalization, only one platform sees the `backspace` string, breaking filer navigation and prompt editing.
- Centralizing the mapping in `widget.EventToString` mirrors the `KeyEnter` fix and guarantees every widget that depends on canonical key strings (filer, finder, completion, menu) inherits the behavior.
- Keeping the canonical string aligned with the existing keymap entries preserves documentation (`README`) and `KeymapBaseline` coverage while satisfying `[REQ:MODULE_VALIDATION]`.

**Module Boundaries & Contracts** `[REQ:MODULE_VALIDATION]`:
- `InputEventTranslator` (existing `widget.EventToString`) MUST translate both `tcell.KeyBackspace` and `tcell.KeyBackspace2` into `backspace` while leaving other keys untouched. This ensures `filerKeymap`, `cmdlineKeymap`, `finderKeymap`, and `completionKeymap` receive the same symbol they already bind to parent navigation and `DeleteBackwardChar`.
- `KeymapBaselineSuite` (tests in `main_keymap_test.go`) documents that `backspace` remains a required chord, providing regression coverage beyond the translator-specific tests.

**Validation Plan:**
- Add table-driven unit tests under `widget/widget_test.go` that create events for both `tcell.KeyBackspace` variants and assert `EventToString` returns `backspace`.
- Continue running `KeymapBaselineSuite` so any accidental removal of the binding fails CI immediately.
- Manual smoke test on macOS Terminal and Linux/tmux sessions to prove Backspace opens the parent directory in filer views and deletes characters in prompts.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `widget/widget.go` map entries for `tcell.KeyBackspace`/`tcell.KeyBackspace2` include `// [IMPL:BACKSPACE_TRANSLATION] [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]`.
- `widget/widget_test.go` introduces `TestEventToStringBackspace_REQ_BACKSPACE_BEHAVIOR` with matching token references.

**Cross-References**: [REQ:BACKSPACE_BEHAVIOR], [IMPL:BACKSPACE_TRANSLATION], [REQ:MODULE_VALIDATION]

## 26. Terminal Launcher Abstraction [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD]

### Decision: Create a dedicated terminal launcher module that selects the correct command sequence for tmux, Linux desktops, and macOS Terminal.
**Rationale:**
- Restores “execute in terminal” workflows on macOS by invoking Terminal.app via `osascript` while preserving the Linux/tmux experience.
- Centralizes platform detection and overrides so future terminal additions (iTerm2, alacritty, Kitty) only touch one module.
- Keeps `main.go` declarative by wiring `g.ConfigTerminal` to a pure factory that can be unit-tested per [REQ:MODULE_VALIDATION].

**Architecture Outline:**
- Introduce package `terminalcmd` with two modules:
  1. `CommandFactory` (pure) returns the `[]string` invocation for the active environment. Inputs: requested shell command, detected tmux/screen, runtime GOOS, optional `GOFUL_TERMINAL_CMD` override.
  2. `Configurator` wires the factory into `g.ConfigTerminal`, logs `DEBUG: [IMPL:TERMINAL_ADAPTER] ...` selections, applies the existing “HIT ENTER KEY” tail, and injects a working-directory provider so macOS launches can `cd` into `%D` before running the command.
- Selection order:
  - If `GOFUL_TERMINAL_CMD` (or config flag) is set, split and use it verbatim.
  - Else if `is_tmux`, run `tmux new-window -n <title> <cmd+tail>`.
  - Else if `runtime.GOOS == "darwin"`, run `osascript -e 'tell application "<configured app>" to do script "%s" & activate'` so the selected macOS terminal (default Terminal.app, override via `GOFUL_TERMINAL_APP`) begins inside the focused directory while `GOFUL_TERMINAL_SHELL` controls the inline shell command (default `bash`).
  - Else default to current Linux behavior: gnome-terminal (with title escape) running bash.
- Provide extension points for future emulators by returning structured data rather than building strings inline. macOS-specific extension points include:
  - `GOFUL_TERMINAL_APP` (default `Terminal`) so operators can direct the AppleScript payload to another application (iTerm2, Warp, etc.) without writing their own `osascript`.
  - `GOFUL_TERMINAL_SHELL` (default `bash`) so the in-window command can switch to `zsh`, `fish`, or another shell without edits.


**Module Validation [REQ:MODULE_VALIDATION]:**
- `CommandFactory` validated via table-driven unit tests (no external processes).
- `Configurator` validated via integration-style tests that stub `g.ConfigTerminal` and assert it receives factory output plus the tail suffix.
- Manual validation checklist documents macOS Terminal run plus Linux gnome-terminal run.

**Alternatives Considered:**
- Hard-coding another case inside `main.go`: rejected because it scales poorly and is hard to test.
- Shelling out to `open -a Terminal ...`: rejected in favor of `osascript` so we can run the command immediately and keep the window open.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `terminalcmd/factory.go`, `terminalcmd/factory_test.go`, and `main.go` wiring must include `[IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY]`.
- README/CONTRIBUTING updates describe overrides and reference `[REQ:TERMINAL_PORTABILITY]`.

**Cross-References**: [REQ:TERMINAL_PORTABILITY], [REQ:TERMINAL_CWD], [IMPL:TERMINAL_ADAPTER], [REQ:MODULE_VALIDATION]

## 27. Event Loop Shutdown [ARCH:EVENT_LOOP_SHUTDOWN] [REQ:EVENT_LOOP_SHUTDOWN]

### Decision: Introduce a coordinated shutdown signal for the UI event poller so `app.Goful.Run` can stop `widget.PollEvent` without leaking goroutines or writing to closed channels.
**Rationale:**
- Today the poller runs an infinite loop that never observes `g.exit`; once the UI exits the goroutine keeps spinning, causing CPU spikes and stray writes to `g.event`.
- Explicit stop control aligns with `[REQ:MODULE_VALIDATION]` by giving the poller a contract that can be validated in isolation.

**Architecture Outline:**
- Wrap `widget.PollEvent` invocations inside a `Poller` module that accepts a context or `stop <-chan struct{}` derived from `app.Goful` lifecycle hooks.
- Provide a `ShutdownController` that closes the stop signal when `Run` is unwinding, waits for the poller to exit (with timeout), then closes `g.event` safely.
- Emit `DEBUG: [IMPL:EVENT_LOOP_SHUTDOWN]` logs when shutdown begins, completes, or times out.

**Alternatives Considered:**
- Leave the infinite loop in place and rely on process exit: rejected because CLI integrations embed goful and expect clean teardown.
- Use `runtime.Goexit` or `os.Exit` to kill goroutines: rejected because it bypasses cleanup and breaks terminal restoration.

**Module Validation [REQ:MODULE_VALIDATION]:**
- `Poller` module validated with fakes that simulate tcell events and assert termination once the stop channel closes.
- `ShutdownController` validated with unit tests that model successful shutdown and timeout paths.
- Integration tests for `app.Goful.Run` ensure no events are delivered after shutdown and goroutine counts return to baseline.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: `app/goful.go` poller, shutdown controller, and any new helper files include `[IMPL:EVENT_LOOP_SHUTDOWN] [ARCH:EVENT_LOOP_SHUTDOWN] [REQ:EVENT_LOOP_SHUTDOWN]` comments.
- Tests: `app/goful_shutdown_test.go` (or equivalent) names incorporate `REQ_EVENT_LOOP_SHUTDOWN` to prove validation coverage.

**Cross-References**: [REQ:EVENT_LOOP_SHUTDOWN], [IMPL:EVENT_LOOP_SHUTDOWN], [REQ:MODULE_VALIDATION]

## 28. Xform CLI Pipeline [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]

### Decision: Split the `scripts/xform.sh` helper into a reusable argument parser and command builder so multi-target commands can be rewritten (inserting a configurable prefix after a configurable number of untouched arguments) without duplicating fragile shell quoting.
**Rationale:**
- External command recipes and developer workflows need the same transformation—keep the first *N* arguments, then interleave `<prefix> target` pairs—so the helper must be callable as either a script or a sourced function.
- Providing configurable `--prefix` and `--keep` options keeps the helper generic enough for different consumers while preserving the simple dry-run/exec flow.
- Keeping the parser and builder as pure Bash functions enables module validation via shell-based tests despite Bash 3.2 limitations on macOS.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `XformArgs` parses flags (`-p/--prefix`, `-k/--keep`, `-n/--dry-run`, `-h/--help`, `--`) and ensures at least `keep + 1` positional arguments remain. It normalizes defaults (`prefix="--to"`, `keep=2`) and returns structured state without invoking external commands. Errors exit with status 64 after printing help text.
- `TargetInterleaver` constructs the transformed argv array by copying the leading `keep` arguments verbatim and pairing each remaining argument with the selected prefix. It emits `%q`-formatted output in dry-run mode while execution mode runs the transformed command and propagates its exit status.
- Integration glue detects whether the script is executed directly (`BASH_SOURCE[0] == $0`) or sourced, defining `xform()` for reuse across scripts.

**Alternatives Considered:**
- Inline argument manipulation inside every external command: rejected because it duplicates quoting logic and increases maintenance overhead.
- A Go-based helper: rejected for now to avoid adding a build step for quick shell automation (script must run anywhere Git checkout exists).

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `scripts/xform.sh` comments annotate both modules with `[IMPL:XFORM_CLI_SCRIPT] [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]`.
- `scripts/xform_test.sh` test cases include `[REQ:CLI_TO_CHAINING]` references proving module validation for parser and interleaver behavior.

**Cross-References**: [REQ:CLI_TO_CHAINING], [IMPL:XFORM_CLI_SCRIPT], [REQ:MODULE_VALIDATION]

## 29. Workspace Bootstrap from Positional Directories [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]

### Decision: Treat trailing CLI arguments as ordered workspace directories and seed filer windows before the UI starts.
**Rationale:**
- Keeps invocation parity with other TUIs by letting scripts dictate deterministic layouts without mutating persisted state files.
- Centralizes the behavior into two testable modules so `[REQ:MODULE_VALIDATION]` can be satisfied independently of the rest of the startup logic.
- Ensures fallback behavior remains unchanged when no positional arguments are provided, protecting long-standing workflows.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `StartupDirParser` – Pure helper that consumes `flag.Args()` after `flag.Parse()`, trims whitespace, expands `~`, resolves absolute paths, and collects warnings for invalid entries (missing, non-directory, permission errors). Returns the ordered slice that mirrors user input so duplicates produce multiple panes intentionally.
- `WorkspaceSeeder` – Receives the parsed slice plus the `*app.Goful` instance, resizing the workspace list to match (creating new panes or closing extras) and calling `Dir().Chdir()` / `ReloadAll()` so each pane shows the desired directory. Emits `DEBUG: [IMPL:WORKSPACE_START_DIRS] ...` lines when `GOFUL_DEBUG_WORKSPACE` is set to document the mutations.
- Error Handling – `message.Errorf` surfaces warnings before the UI loop begins; invalid entries are skipped so remaining valid directories still open. If every argument fails, the helper aborts and defaults to the persisted workspace layout.

**Alternatives Considered:**
- **Repeated `-workspace` flags** – Rejected because it duplicates shell quoting requirements and diverges from common CLI expectations where positional directories control startup state.
- **Post-startup seeding** – Rejected to avoid visible flicker and to keep module validation simple by running before `goful.Run()`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: `main.go` and `app/startup_dirs.go` include `[IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]` annotations.
- Tests: `app/startup_dirs_test.go` (parser/seeder unit coverage) and future integration tests reference `[REQ:WORKSPACE_START_DIRS]`.

**Cross-References**: [REQ:WORKSPACE_START_DIRS], [IMPL:WORKSPACE_START_DIRS], [REQ:MODULE_VALIDATION]

## 29. Workspace Bootstrap from Positional Directories [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]

### Decision: Treat trailing CLI arguments as ordered workspace directories and seed filer windows accordingly before the UI starts.
**Rationale:**
- Power users often launch goful via scripts and expect deterministic multi-pane layouts without manual navigation; positional args keep invocation parity with other TUIs.
- Centralizing parsing + seeding logic maintains testability and ensures `[REQ:MODULE_VALIDATION]` coverage by isolating pure helpers from widget wiring.
- Enables automation-friendly diagnostics (`GOFUL_DEBUG_WORKSPACE`) describing how startup directories were interpreted without polluting runtime code paths elsewhere.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `StartupDirParser` (Module 1) – Consumes `flag.Args()` after `flag.Parse()`, normalizes paths via `util.ExpandPath`, and returns `(orderedDirs []string, warnings []error)` so callers can log non-fatal issues (e.g., missing paths) while proceeding with valid ones in the exact order entered. Empty inputs yield `nil` and indicate fallback to legacy behavior.
- `WorkspaceSeeder` (Module 2) – Accepts the parsed directory list plus `*app.Goful` and mutates the underlying `filer.Workspace` so the window count and ordering match the supplied arguments. Responsibilities:
  - Create additional directories via `g.CreateWorkspace()` / `g.Dir().Chdir()` style helpers when more args exist than windows.
  - Reuse existing windows when counts align by calling `Workspace().Chdir(idx, dir)` without resetting focus unnecessarily.
  - Close surplus windows (beyond one) when fewer args arrive than currently open panes, preserving the first window for the first directory.
  - Emit `DEBUG: [IMPL:WORKSPACE_START_DIRS] ...` lines when `GOFUL_DEBUG_WORKSPACE=1`, documenting before/after state.
- Error Handling Contract – Missing directories trigger `message.Errorf` but do not abort seeding; invalid arguments are skipped so valid windows still open. When every argument fails, the seeder aborts and falls back to the default workspace.

**Alternatives Considered:**
- **New `-workspace` flag per directory**: rejected to avoid duplicating shell quoting rules and to keep parity with existing TUIs that use positional args.
- **Lazy seeding post UI launch**: rejected because it complicates module validation and introduces flicker as panes are added after widget initialization.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: `main.go` helper `applyStartupDirs` plus new parser/seeder files include `[IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]`.
- Tests: `app/startup_dirs_test.go` (unit) and `filer/integration_test.go` additions include `[REQ:WORKSPACE_START_DIRS]` to prove module validation before integration.

**Cross-References**: [REQ:WORKSPACE_START_DIRS], [IMPL:WORKSPACE_START_DIRS], [REQ:MODULE_VALIDATION]

## 30. File Comparison Engine [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]

### Decision: Implement progressive comparison with cached indexing for cross-directory file color-coding
**Rationale:**
- Provides instant visual feedback for file relationships across workspace directories without manual comparison.
- Progressive rendering ensures directories display immediately while comparison analysis runs in background.
- Caching comparison results avoids recomputation on every draw, improving performance for large directories.
- Configurable color scheme via YAML allows users to customize the visual appearance to their preferences.

**Architecture Outline:**
- **ComparisonColorConfig** (`filer/comparecolors/`): Loads YAML configuration defining color roles (name_present, size_equal/smallest/largest, time_equal/earliest/latest). Falls back to sensible defaults when config is missing or invalid. Path resolution via `configpaths.Resolver` extension.
- **FileComparisonIndex** (`filer/compare.go`): Pure module that builds an index of file names across workspace directories. For each file appearing in multiple directories, computes comparison states for size and time. Results are cached per workspace and invalidated on directory changes.
- **ComparisonLook** (`look/comparison.go`): Exposes `tcell.Style` getters for each comparison state. Manages runtime toggle state with thread-safe access.
- **Draw Integration**: `FileStat.Draw()` accepts optional comparison context and applies appropriate colors per field (name, size, time) independently.

**Progressive Rendering Strategy:**
1. Initial draw renders file lists with standard colors (no blocking).
2. After all directories complete loading, comparison index is built.
3. Subsequent draws use cached comparison results.
4. Cache invalidated on: `Chdir`, `reload`, `ReloadAll`, `CreateDir`, `CloseDir`, toggle on.

**Default State**: Comparison coloring is **enabled by default** so users immediately benefit from color-coded file comparisons without manual activation. The backtick toggle (`` ` ``) disables the feature when standard colors are preferred.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `CompareColorLoader` (Module 1): Parses YAML config, validates color names, provides defaults. Independently testable with mock file readers.
- `FileComparisonIndex` (Module 2): Pure function that takes workspace directories and returns comparison state map. No side effects, independently testable.
- `ComparisonCache` (Module 3): Thread-safe storage keyed by (dirIndex, filename). Provides O(1) lookup during draw.
- `ComparisonLook` (Module 4): Style providers that consume loaded config. Tested with mock styles.

**Alternatives Considered:**
- **Synchronous comparison on every draw**: Rejected due to performance impact on large directories.
- **Background goroutine for comparison**: Rejected to avoid concurrency complexity; progressive rendering with caching is simpler and sufficient.
- **Hardcoded colors only**: Rejected; users have diverse terminal themes and color preferences.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `filer/compare.go`, `filer/comparecolors/` include `[IMPL:FILE_COMPARISON_INDEX] [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`.
- `look/comparison.go` includes `[IMPL:COMPARISON_LOOK] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`.
- `filer/file.go` draw integration includes `[IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`.
- Tests reference `[REQ:FILE_COMPARISON_COLORS]` to prove validation coverage.

**Cross-References**: [REQ:FILE_COMPARISON_COLORS], [IMPL:COMPARE_COLOR_CONFIG], [IMPL:FILE_COMPARISON_INDEX], [IMPL:COMPARISON_DRAW], [REQ:MODULE_VALIDATION]

## 32. Linked Navigation Mode [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]

### Decision: Implement linked navigation as a toggleable mode that propagates directory changes across workspace windows.
**Rationale:**
- Operators comparing mirrored directory structures need synchronized pane navigation to reduce repetitive manual work.
- Keeping the mode as a simple boolean toggle keeps the implementation minimal and avoids complex state management.
- Propagating navigation to "matching subdirectories" (by name) allows partial synchronization when structures differ slightly.

**Architecture Outline:**
- **State Management**: Add `linkedNav bool` field to `app.Goful` struct with getter/toggle methods.
- **Navigation Helpers**: Add `ChdirAllToSubdir(name string)` and `ChdirAllToParent()` to `filer.Workspace` that iterate non-focused directories and attempt navigation.
- **Header Indicator**: Extend `filer.Filer.drawHeader()` to show `[LINKED]` when the mode is enabled, using a callback or exported flag.
- **Keymap Integration**: Wrap existing navigation callbacks (backspace, enter-dir) to check linked state and invoke the appropriate helper.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `LinkedNavState` (Module 1 in `app/goful.go`) – Pure toggle and query methods; no side effects beyond flipping the boolean.
- `LinkedNavigationHelpers` (Module 2 in `filer/workspace.go`) – Pure workspace methods that iterate directories and call `Chdir`; do not mutate linked state.
- `LinkedNavIndicator` (Module 3 in `filer/filer.go`) – Header rendering that consumes an external flag/callback; no business logic.

**Alternatives Considered:**
- **Per-workspace linked state**: Rejected because operators typically want global linked mode, and per-workspace adds UI complexity.
- **Automatic linking based on directory similarity**: Rejected as too implicit; explicit toggle gives operators control.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Code: `app/goful.go`, `filer/workspace.go`, `filer/filer.go`, `main.go` include `[IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`.
- Tests: `filer/workspace_test.go` and integration tests include `[REQ:LINKED_NAVIGATION]` in names/comments.

**Cross-References**: [REQ:LINKED_NAVIGATION], [IMPL:LINKED_NAVIGATION], [REQ:MODULE_VALIDATION]

## 33. Difference Search Engine [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]

### Decision: Implement a two-command difference search that iterates through workspace files alphabetically and stops at entries that differ across windows.
**Rationale:**
- Users comparing similar directory structures need an efficient way to find differences without manually inspecting each file.
- Separating "Start" and "Continue" commands allows users to inspect differences at their own pace while maintaining search state.
- Using cursor position as the implicit bookmark simplifies state management and integrates naturally with existing navigation.

**Architecture Outline:**
- **State Management**: Add `DiffSearchState` struct to `filer.Workspace` containing `initialDirs []string` (one per window) and `active bool`.
- **Core Comparison**: Pure module that collects the union of all filenames across directories, sorts them alphabetically (case-sensitive), and iterates to find differences.
- **Difference Detection**: A file is "different" if it's missing from any window OR has different sizes across windows. Directories follow the same rule.
- **Cursor Movement**: When a difference is found, move cursors to that entry in all windows where it exists.
- **Files First, Then Directories**: At each level, ALL files are processed in alphabetic order before ANY directories are processed. This ensures a consistent, predictable search order. When `startAfter` is a subdirectory name, file check is skipped (all files at that level have already been processed).
- **Subdirectory Descent**: If no file differences are found at a level, process subdirectories alphabetically at that level. If a subdir exists in all windows, descend into it and repeat. If a subdir is missing in any window, treat it as a difference. The descent function must respect the `startAfter` position to avoid re-searching already-visited subdirectories (see `FindNextSubdirInAll`).
- **Continue Logic**: Command 2 uses the recorded difference name from state (with "/" suffix removed for directories) to skip the last found difference and continue from the next entry. This ensures correct continuation even when the cursor cannot be set (e.g., when a subdirectory is missing in some windows). When descending into subdirectories, the algorithm uses `FindNextSubdirInAll(dirs, subdirStartAfterForDescent)` where `subdirStartAfterForDescent` is determined from `startAfter` (empty if `startAfter` is a filename, or the subdirectory name if `startAfter` is a subdirectory name).
- **Search Completion**: When ascending from a subdirectory to root level, if there are no more subdirectories after the one ascended from, the search is complete and announces "Difference search complete - all differences found".

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `DiffSearchState` (Module 1 in `filer/diffsearch.go`): Pure struct holding initial directories, active state, and status fields (LastDiffName, LastDiffReason, CurrentPath, FilesChecked, Searching). Methods for starting, checking active, clearing, and updating status.
- `DiffSearchEngine` (Module 2 in `filer/diffsearch.go`): Pure functions for collecting file names, detecting differences, and finding the next different entry. No side effects, independently testable.
- `Navigator` (Module 3 in `filer/diffsearch.go`): Interface abstracting directory operations for tree traversal (`GetDirs`, `ChdirAll`, `ChdirParentAll`, `CurrentPath`, `RebuildComparisonIndex`). Allows TreeWalker to be tested with mock implementations.
- `TreeWalker` (Module 4 in `filer/diffsearch.go`): Pure traversal algorithm that uses Navigator interface. `Run(progressFn)` method executes the traversal loop and returns a `Step` result (`StepFoundDiff` or `StepComplete`). Independently testable via MockNavigator.
- `WorkspaceNavigator` (Module 5 in `filer/workspace.go`): Adapter implementing Navigator interface for Workspace. Bridges TreeWalker to real filesystem operations.
- `DiffSearchNavigation` (Module 6 in `filer/workspace.go`): Workspace methods that move cursors to a named entry across all directories.
- `DiffSearchCommands` (Module 7 in `app/goful.go`): Command wrappers that wire state, engine, and navigation together. Creates WorkspaceNavigator, TreeWalker, handles Step result, and manages TUI concerns (periodic UI refresh via goroutine/ticker, status messages, resize).
- `DiffStatusDisplay` (Module 8 in `diffstatus/diffstatus.go`): Dedicated status line display that persists while diff search is active. Shows current search progress or last found difference. Integrated into app draw cycle and resize handling.

**Navigator Interface Pattern:**
The Navigator interface decouples the tree traversal algorithm from TUI concerns:
```go
type Navigator interface {
    GetDirs() []*Directory           // Current directories being compared
    ChdirAll(name string)            // Descend into subdirectory
    ChdirParentAll()                 // Ascend to parent
    CurrentPath() string             // Path of first directory
    RebuildComparisonIndex()         // Refresh after directory changes
}

type StepType int
const (
    StepFoundDiff StepType = iota    // Difference found, search pauses
    StepComplete                      // Search complete
)

type Step struct {
    Type   StepType
    Name   string   // Difference name (with "/" suffix for directories)
    Reason string   // Why it's different
    IsDir  bool     // Whether it's a directory
}
```
This pattern enables:
- Unit testing with MockNavigator (no real filesystem needed for algorithm tests)
- Clear separation between traversal logic and TUI side effects
- Easy addition of new Navigator implementations (e.g., for remote filesystems)

**Algorithm Sketch:**
```text
// [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func findNextDifference(dirs []*Directory, startAfter string) (name string, reason string, found bool) {
    names := collectAllNames(dirs)  // Union of all file/dir names
    sort.Strings(names)              // Case-sensitive alphabetic
    
    for _, name := range names {
        // Skip entries that come before or equal to startAfter alphabetically.
        // Using alphabetical comparison (not exact match) allows startAfter to be
        // a name not in the current list (e.g., a filename when searching subdirs).
        if startAfter != "" && name <= startAfter {
            continue
        }
        if isDifferent, reason := checkDifference(name, dirs); isDifferent {
            return name, reason, true
        }
    }
    return "", "", false  // No more differences at this level
}
```

**Important**: The alphabetical comparison (`name <= startAfter`) instead of exact match (`name == startAfter`) is critical. When continuing a search from a file position (e.g., "date.key"), the subdirectory search needs to find directories that come AFTER that file alphabetically (e.g., "dev", "org", "temp"). Since "date.key" is not in the subdirectory list, exact match would fail to set `started=true`, skipping all subdirectories.

**Alternatives Considered:**
- **Skip set for tracking inspected files**: Rejected; cursor position provides implicit tracking with no extra state.
- **Background comparison thread**: Rejected; synchronous iteration is simpler and search is bounded by user interaction pace.
- **Modal "diff mode" with separate keybindings**: Rejected; background state allows normal navigation between continue commands.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `filer/diffsearch.go` includes `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`.
- `filer/workspace.go` navigation helpers include `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`.
- `app/goful.go` command wiring includes `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`.
- `diffstatus/diffstatus.go` includes `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`.
- `main.go` wiring includes `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`.
- Tests reference `[REQ:DIFF_SEARCH]` in names/comments.

**Cross-References**: [REQ:DIFF_SEARCH], [IMPL:DIFF_SEARCH], [REQ:MODULE_VALIDATION]

## 31. Filename Exclude Filter [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]

### Decision: Load newline-delimited basename filters via flag/env/default and enforce them inside the filer pipeline with a runtime toggle.
**Rationale:**
- Keeps nuisance files (e.g., `.DS_Store`, `Thumbs.db`, build outputs) out of every pane without hand-maintained per-directory ignores.
- Aligns with existing precedence rules (`flag` > `env` > default) so operators can share curated lists without code edits.
- Provides a runtime toggle so users can temporarily show everything, inspect results, and return to the filtered view without restarting goful.

**Module Boundaries & Contracts** `[REQ:MODULE_VALIDATION]`:
- `ExcludeRules` (Module 1 – package `filer`): Pure helper that stores a case-insensitive set of excluded basenames, exposes `ConfigureExcludedNames`, `ToggleExcludedNames`, and `ShouldExclude`. The `Directory.read` callback consults the helper before appending `FileStat` entries so every reader (default, glob, finder) honors the filter automatically. Excludes remain disabled when no names are configured to preserve historical behaviour.
- `ExcludeListLoader` (Module 2 – package `main`): Reads the newline-delimited file resolved by `configpaths.Resolver` (`-exclude-names`, `GOFUL_EXCLUDES_FILE`, default `~/.goful/excludes`), trims whitespace, skips blank lines and `#` comments, normalizes to lower-case, and calls `filer.ConfigureExcludedNames`. Emits `DEBUG:` / `message.Infof` diagnostics tagged with `[IMPL:FILER_EXCLUDE_LOADER]` that list how many entries were loaded or why the file was skipped. Surfaces a View menu entry plus a dedicated `E` key binding that toggles the filter (`Workspace.ReloadAll` afterward) so users can observe changes immediately.

**Alternatives Considered:**
- **Inline glob filters per workspace**: Rejected because it introduces per-pane state and complicates persistence.
- **Embedding excludes in the existing JSON state**: Rejected to keep the block list shareable outside the Go binary and to avoid migrations.
- **Directory-level filter only**: Rejected because finder/glob flows would still leak excluded entries.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `configpaths/resolver.go`, `main.go`, and loader tests include `[IMPL:FILER_EXCLUDE_LOADER] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]`.
- `filer/exclude.go`, `filer/directory.go`, and filer tests include `[IMPL:FILER_EXCLUDE_RULES] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]`.
- Toggle handlers log `[REQ:FILER_EXCLUDE_NAMES]` so operators can trace runtime state changes.

**Cross-References**: [REQ:FILER_EXCLUDE_NAMES], [IMPL:FILER_EXCLUDE_RULES], [IMPL:FILER_EXCLUDE_LOADER], [REQ:MODULE_VALIDATION]

## 34. nsync SDK Integration [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]

### Decision: Integrate the external nsync SDK to provide parallel multi-destination file synchronization as an alternative to the builtin single-target copy/move.
**Rationale:**
- The builtin `copy`/`move` functions in `app/filectrl.go` work with a single destination, requiring users to repeat operations for each target pane.
- The nsync SDK (`github.com/fareedst/nsync/pkg/nsync`) provides production-ready parallel multi-destination sync with progress monitoring, content verification, and move semantics.
- A hybrid approach keeps the existing single-target operations unchanged (muscle memory preserved) while adding explicit "Copy All"/"Move All" commands for multi-destination workflows.
- Using all visible workspace panes as implicit targets aligns with goful's visual paradigm—users see all destinations on screen before invoking the command.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `NsyncObserver` (Module 1 – `app/nsync.go`): Implements `nsync.Observer` interface to bridge nsync progress events to goful's `progress` widget and `message` package. Translates `OnStart`/`OnProgress`/`OnFinish` callbacks to `progress.Start()`/`progress.Update()`/`progress.Finish()` calls. Thread-safe for concurrent goroutine callbacks.
- `SyncCopy`/`SyncMove` (Module 2 – `app/nsync.go`): Wrapper functions that configure `nsync.Config` with sources/destinations, create a `Syncer` with the observer, and execute `Sync()` within `asyncFilectrl` for UI integration. Handle context cancellation for user interrupts and aggregate errors per destination.
- `CopyAll`/`MoveAll` (Module 3 – `app/nsync.go`): Functions that enumerate destination directories from `otherWindowDirPaths()` (reusing the existing `%D@` macro helper), collect source files from marks or cursor, and delegate to `syncCopy`/`syncMove`. Fall back to builtin operations when only one pane exists.

**Dependency Management:**
- nsync is added as a dependency from the public repository `github.com/fareedst/nsync`.
- Transitive dependencies (xxhash, blake3, uuid, x/sync) are already compatible with goful's Go 1.24 toolchain.

**Alternatives Considered:**
- **Extend builtin walker for multiple destinations**: Rejected because it duplicates the parallelism, verification, and progress machinery that nsync already provides.
- **External shell command with `%D@`**: Rejected because it bypasses goful's progress display and error handling.
- **Always use nsync for all copy/move**: Rejected to preserve backward compatibility and avoid regression risk for single-target operations.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `app/nsync.go` includes `[IMPL:NSYNC_OBSERVER] [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]`.
- `app/nsync.go` `CopyAll`/`MoveAll` functions include `[IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]`.
- `main.go` keybinding wiring includes `[IMPL:NSYNC_COPY_MOVE] [REQ:NSYNC_MULTI_TARGET]`.
- Tests reference `[REQ:NSYNC_MULTI_TARGET]` in names/comments.

**Cross-References**: [REQ:NSYNC_MULTI_TARGET], [IMPL:NSYNC_OBSERVER], [IMPL:NSYNC_COPY_MOVE], [REQ:MODULE_VALIDATION]

## 35. nsync Confirmation Mode [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]

### Decision: Add confirmation prompts before multi-target copy/move operations using existing cmdline mode pattern.
**Rationale:**
- Multi-target operations (`CopyAll`/`MoveAll`) affect multiple directories simultaneously and are not easily reversible.
- Users need to verify their intent before files are copied or moved to all visible panes.
- Move operations are particularly risky since source files are deleted after successful sync.
- Reusing the existing cmdline mode pattern (`quitMode`/`removeMode`) keeps implementation simple and UX consistent.

**Architecture Outline:**
- Create `copyAllMode` and `moveAllMode` structs in `app/mode.go` following the `quitMode`/`removeMode` pattern.
- Each mode implements `cmdline.Mode` interface with `String()`, `Prompt()`, `Draw()`, and `Run()` methods.
- The `Prompt()` method displays source count and destination count (e.g., `Copy 3 file(s) to 2 destinations? [Y/n]`).
- The `Run()` method handles user input: `Y`/`y`/empty confirms, `n`/`N` cancels, other input clears.
- Confirmation mode stores sources, destinations, and a reference to `*Goful` to execute the operation.
- Public `CopyAll()`/`MoveAll()` functions collect sources/destinations first, then start the confirmation mode.
- Actual sync execution is deferred to private `doCopyAll()`/`doMoveAll()` methods called only after confirmation.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `copyAllMode`/`moveAllMode` (Module 1 – `app/mode.go`): Confirmation modes that display prompts and handle user input. No side effects until confirmed.
- `doCopyAll()`/`doMoveAll()` (Module 2 – `app/nsync.go`): Private execution methods containing the actual nsync sync logic, called only after user confirmation.

**Alternatives Considered:**
- **No confirmation**: Rejected because multi-target operations are high-risk and users expect confirmation for destructive operations.
- **Synchronous dialog()**: Rejected because it blocks the event loop; cmdline modes integrate better with the existing UI.
- **Confirmation only for moves**: Rejected because copy operations to multiple destinations are also significant and merit confirmation.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `app/mode.go` confirmation modes include `[IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]`.
- `app/nsync.go` refactored methods include `[IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]`.
- Tests reference `[REQ:NSYNC_CONFIRMATION]` in names/comments.

**Cross-References**: [REQ:NSYNC_CONFIRMATION], [IMPL:NSYNC_CONFIRMATION], [REQ:NSYNC_MULTI_TARGET], [REQ:MODULE_VALIDATION]

## 36. Help Widget [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]

### Decision: Create a Help widget based on the existing ListBox/Menu pattern to display a scrollable keystroke catalog.
**Rationale:**
- Users need quick access to keystroke documentation without leaving the application.
- Reusing the existing `widget.ListBox` pattern ensures consistent UI behavior and reduces implementation complexity.
- The popup model (connect via `g.next`, disconnect on exit) integrates naturally with goful's widget chain.

**Architecture Outline:**
- Create package `help` with a `Help` struct that embeds `*widget.ListBox`.
- `New(filer widget.Widget)` constructor creates a popup sized to ~80% of screen dimensions, centered visually.
- The keystroke catalog is maintained as a Go slice within the help package for easy updates.
- `Input(key string)` handles navigation (`C-n`, `C-p`, `up`, `down`, `pgup`, `pgdn`) and exit (`?`, `q`, `C-g`, `C-[`).
- `Exit()` disconnects from the widget chain via `filer.Disconnect()`.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `HelpWidget` (Module 1 – `help/help.go`): Encapsulates the ListBox, keystroke data, and input handling. No external dependencies beyond `widget` package.
- `HelpCommand` (Module 2 – `app/goful.go`): `Help()` method that creates the widget and connects it to the widget chain.
- `HelpKeyBinding` (Module 3 – `main.go`): Wires `?` keystroke to `g.Help()`.

**Alternatives Considered:**
- **External documentation only**: Rejected because users must leave the application to find help.
- **Modal overlay with raw text**: Rejected because ListBox provides scrolling and consistent styling for free.
- **Separate help command/menu**: Rejected because `?` is the universal help convention.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `help/help.go` includes `[IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]`.
- `app/goful.go` `Help()` method includes `[IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]`.
- `main.go` keystroke binding includes `[IMPL:HELP_POPUP] [REQ:HELP_POPUP]`.

**Cross-References**: [REQ:HELP_POPUP], [IMPL:HELP_POPUP], [REQ:MODULE_VALIDATION]

## 38. Mouse Event Routing [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]

### Decision: Enable tcell mouse events and extend the main event loop to dispatch mouse events to the appropriate widget via hit-testing.
**Rationale:**
- Mouse input is expected in modern terminal applications and provides an alternative navigation method.
- Centralizing mouse event handling in `app.Goful.eventHandler` maintains consistency with existing keyboard and resize event routing.
- Hit-testing (determining which widget contains the click coordinates) enables clean separation between event detection and widget response.
- The existing `widget.Widget` interface can be extended with coordinate-based methods without breaking backward compatibility.

**Architecture Outline:**
- **Mouse Initialization** (`widget/widget.go`): Call `screen.EnableMouse()` after `screen.Init()` in the `Init()` function. This enables all mouse events (clicks, drags, wheel) at the tcell level.
- **Event Dispatch** (`app/goful.go`): Extend `eventHandler` to handle `*tcell.EventMouse` events. Extract position via `ev.Position()` and button via `ev.Buttons()`.
- **Hit-Testing Framework** (`widget/widget.go`, `filer/workspace.go`):
  - Add `Contains(x, y int) bool` method to `widget.Window` that returns true if coordinates are within the window bounds.
  - Add `DirectoryAt(x, y int) (*Directory, int)` method to `filer.Workspace` that finds which directory window contains the coordinates.
- **Mouse Handler** (`app/goful.go`): New `mouseHandler(ev *tcell.EventMouse)` method that:
  1. Checks if a modal (`g.Next()`) is active and handles modal-specific mouse events.
  2. Otherwise, delegates to filer-level mouse handling for directory selection.

**Multi-Stage Implementation:**
1. **Stage 1 - Infrastructure**: Enable mouse, add `EventMouse` case to handler, create hit-test framework.
2. **Stage 2 - File Selection**: Implement `FileIndexAtY`, wire click to `SetCursor`.
3. **Stage 3 - Window Focus**: Clicking unfocused window switches focus.
4. **Stage 4 - Scrolling**: Mouse wheel support.
5. **Stage 5 - Modals**: Menu and modal mouse support (future).

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `MouseEventTranslator` (Module 1 in `widget/widget.go`) – Enables mouse at init, provides helpers for button detection.
- `HitTestFramework` (Module 2 in `widget/widget.go`, `filer/workspace.go`) – Pure coordinate-to-widget mapping, independently testable.
- `MouseDispatcher` (Module 3 in `app/goful.go`) – Orchestrates hit-testing and widget method calls.

**Alternatives Considered:**
- **Per-widget mouse polling**: Rejected because centralized dispatch is simpler and matches existing keyboard handling.
- **Full GUI-style event bubbling**: Rejected as over-engineered for a terminal file manager.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `widget/widget.go` includes `[IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]`.
- `filer/workspace.go` and `filer/directory.go` include same tokens for hit-testing methods.
- `app/goful.go` includes `[IMPL:MOUSE_FILE_SELECT] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]`.
- Tests reference `[REQ:MOUSE_FILE_SELECT]` in names/comments.

**Cross-References**: [REQ:MOUSE_FILE_SELECT], [IMPL:MOUSE_HIT_TEST], [IMPL:MOUSE_FILE_SELECT], [REQ:MODULE_VALIDATION]

## 39. Mouse Double-Click Detection [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]

### Decision: Implement time-based double-click detection in the mouse handler with action dispatch for directories and files, respecting the Linked navigation mode.
**Rationale:**
- Double-click is the standard "open/enter" gesture in file managers and users expect this behavior.
- Detection requires tracking the last click time and position to determine if two clicks occur within a threshold (400ms) at the same location.
- For directories, double-click should reuse the existing `linkedEnterDir` pattern to respect Linked navigation mode.
- For files, double-click should trigger the open action, and when Linked mode is enabled, should open same-named files in all windows.

**Architecture Outline:**
- **Click State Tracking** (`app/goful.go`): Add `lastClickTime time.Time`, `lastClickX int`, `lastClickY int` fields to `Goful` struct.
- **Double-Click Detection** (`app/goful.go`): Add `isDoubleClick(x, y int) bool` helper that checks if the current click is within threshold of the last click at the same position, then updates the tracking state.
- **Directory Double-Click Handler**: When double-clicking a directory:
  - If Linked mode is ON: Call `ChdirAllToSubdirNoRebuild()` for other windows, then `EnterDir()` for the focused window, then `RebuildComparisonIndex()`.
  - If Linked mode is OFF: Call only `EnterDir()` for the clicked directory.
- **File Double-Click Handler**: When double-clicking a file:
  - If Linked mode is ON: Collect all same-named file paths from all windows and open each with a **separate command** (one `open` invocation per file). This ensures consistent behavior across different opener applications.
  - If Linked mode is OFF: Trigger open action only for the clicked file via the keymap.
- **Integration with `handleLeftClick`**: After single-click selection, check `isDoubleClick()` and dispatch to appropriate handler.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `DoubleClickDetector` (Module 1 in `app/goful.go`) – Pure timing/position logic, independently testable with mock time.
- `DirectoryDoubleClickHandler` (Module 2 in `app/goful.go`) – Reuses linked navigation pattern, calls workspace and directory methods.
- `FileDoubleClickHandler` (Module 3 in `app/goful.go`) – Iterates windows when Linked, collects matching file paths, and opens each with a separate platform-appropriate command (xdg-open/open/explorer).

**Alternatives Considered:**
- **Separate double-click event**: tcell does not provide built-in double-click events, so manual detection is required.
- **Fixed threshold only**: Chose configurable threshold (default 400ms) for future flexibility.
- **Separate "open all windows" option**: Rejected in favor of reusing the existing Linked navigation setting for consistency.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `app/goful.go` includes `[IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]`.
- Tests reference `[REQ:MOUSE_DOUBLE_CLICK]` in names/comments.

**Cross-References**: [REQ:MOUSE_DOUBLE_CLICK], [IMPL:MOUSE_DOUBLE_CLICK], [ARCH:MOUSE_EVENT_ROUTING], [REQ:LINKED_NAVIGATION], [REQ:MODULE_VALIDATION]

## 37. Sync Mode [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]

### Decision: Implement a two-stage prefix mode for executing synchronized file operations across all workspace panes.
**Rationale:**
- Users managing synchronized directory structures need batch operations on same-named files across all panes.
- A prefix key pattern (similar to `X` for external commands) keeps the keymap organized and discoverable.
- Separating the mode activation (`S`) from the operation key (`c`/`d`/`r`) allows a single prompt to apply to all panes.
- "Ignore failures" mode (toggled with `!`) accommodates divergent directory structures where not all operations will succeed.

**Architecture Outline:**
- **Sync Mode Activation**: `S` keypress enters a transient mode that waits for an operation key.
- **Operation Dispatch**: After `S`, pressing `c`, `d`, or `r` starts the corresponding sync operation mode.
- **Prompt Phase**: Each operation mode prompts once for input/confirmation using the existing cmdline mode pattern.
- **Execution Phase**: Sequential iteration through panes, finding files by name, executing the operation.
- **Failure Handling**: Two modes controlled by a boolean flag toggled with `!`.

**Module Boundaries & Contracts `[REQ:MODULE_VALIDATION]`:**
- `SyncMode` (Module 1 – `app/window_wide.go`): Transient prefix mode implementing `cmdline.Mode`. Captures `ignoreFailures bool` from toggle and waits for operation key.
- `SyncCopyMode`/`SyncDeleteMode`/`SyncRenameMode` (Module 2 – `app/mode.go`): Operation-specific modes that prompt for input and call the execution engine. Copy and rename prompt for a new filename (copy defaults to original name, user must change it); delete prompts for y/n confirmation.
- `SyncExecutor` (Module 3 – `app/window_wide.go`): Pure execution engine that iterates panes, finds files by name, and applies operations with configurable failure handling. Returns structured results (success/failure counts, error details).
- `FindFileByName` (Module 4 – `filer/directory.go`): Helper method on `Directory` that searches the file list for an exact name match.
- `Keymap Integration` (Module 5 – `main.go`): Wires `S` key to `g.SyncMode(false)`.

**Algorithm Sketch:**
```text
// [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func executeSync(ws *Workspace, filename string, newName string, op Operation, ignoreFailures bool) Results {
    results := Results{}
    // Start with focused pane, then wrap through others
    for i := 0; i < len(ws.Dirs); i++ {
        idx := (ws.Focus + i) % len(ws.Dirs)
        dir := ws.Dirs[idx]
        
        file := dir.FindFileByName(filename)
        if file == nil {
            results.Skipped++
            continue
        }
        
        // For copy: copy to newName in same directory
        // For rename: rename to newName
        // For delete: remove file
        err := op.Execute(dir.Path, file, newName)
        if err != nil {
            results.Failures = append(results.Failures, SyncFailure{idx, err})
            if !ignoreFailures {
                return results  // Abort on first failure
            }
        } else {
            results.Succeeded++
        }
    }
    return results
}
```

**Alternatives Considered:**
- **Menu-based approach**: Rejected because it adds an extra interaction step; prefix key is more efficient.
- **Automatic all-pane detection**: Rejected; explicit `S` prefix makes the scope clear to the user.
- **Parallel execution**: Rejected to avoid race conditions and to provide predictable, debuggable ordering.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `app/window_wide.go` includes `[IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]`.
- `app/mode.go` operation modes include `[IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]`.
- `filer/directory.go` helper includes `[IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]`.
- `main.go` keybindings include `[IMPL:SYNC_EXECUTE] [REQ:SYNC_COMMANDS]`.
- Tests reference `[REQ:SYNC_COMMANDS]` in names/comments.

**Cross-References**: [REQ:SYNC_COMMANDS], [IMPL:SYNC_EXECUTE], [REQ:MODULE_VALIDATION]
