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

### Decision: Makefile + CI matrix for GOOS/GOARCH static builds
**Rationale:**
- Ensures reproducible release artifacts across platforms.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Makefile targets and workflow jobs annotated with `[IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]`.

**Cross-References**: [REQ:RELEASE_BUILD_MATRIX], [IMPL:MAKE_RELEASE_TARGETS]

## 22. Baseline Capture [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE]

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

## 23. Debt Management [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]

### Decision: Systematically log and annotate risky areas
**Rationale:**
- Makes hidden risks visible for scheduling and refactors.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Issues/TODOs carry `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`.

**Cross-References**: [REQ:DEBT_TRIAGE], [IMPL:DEBT_TRACKING]

## 24. Token Validation Automation [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]

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

