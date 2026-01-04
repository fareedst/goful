# Implementation Decisions

**STDD Methodology Version**: 1.1.0

## Overview
This document captures detailed implementation decisions for this project, including specific APIs, data structures, and algorithms. All decisions are cross-referenced with architecture decisions using `[ARCH:*]` tokens and requirements using `[REQ:*]` tokens for traceability.

## Template Structure

When documenting implementation decisions, use this format:

```markdown
## N. Implementation Title [IMPL:IDENTIFIER] [ARCH:RELATED_ARCHITECTURE] [REQ:RELATED_REQUIREMENT]

### Decision: Brief description of the implementation decision
**Rationale:**
- Why this implementation approach was chosen
- What problems it solves
- How it fulfills the architecture decision

### Implementation Approach:
- Specific technical details
- Code structure or patterns
- API design decisions

**Code Markers**: Specific code locations, function names, or patterns to look for

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Which files/functions must carry the `[IMPL:*] [ARCH:*] [REQ:*]` annotations
- Which tests (names + locations) must reference the matching `[REQ:*]`

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- Latest `./scripts/validate_tokens.sh` (or repo equivalent) output summary
- Date/commit hash of the validation run

**Cross-References**: [ARCH:RELATED_ARCHITECTURE], [REQ:RELATED_REQUIREMENT]
```

## Notes

- All implementation decisions MUST be recorded here IMMEDIATELY when made
- Each decision MUST include `[IMPL:*]` token and cross-reference both `[ARCH:*]` and `[REQ:*]` tokens
- Implementation decisions are dependent on both architecture decisions and requirements
- DO NOT defer implementation documentation - record decisions as they are made
- Record where code/tests are annotated so `[PROC:TOKEN_AUDIT]` can succeed later.
- Include the most recent `[PROC:TOKEN_VALIDATION]` run information so future contributors know the last verified state.
- **Language-Specific Implementation**: Language-specific implementation details (APIs, libraries, syntax patterns, idioms) belong in implementation decisions. Code examples in documentation should use `[your-language]` placeholders or be language-agnostic pseudo-code unless demonstrating a specific language requirement. Requirements and architecture decisions should remain language-agnostic.

---
## 1. Configuration Structure [IMPL:CONFIG_STRUCT] [ARCH:CONFIG_STRUCTURE] [REQ:CONFIGURATION]

### Config Type
```[your-language]
// [IMPL:CONFIG_STRUCT] [ARCH:CONFIG_STRUCTURE] [REQ:CONFIGURATION]
type Config struct {
    // Add your configuration fields here
    Field1 string
    Field2 int
    Field3 bool
}
```

### Default Values
- Field1: default value
- Field2: default value
- Field3: default value

## 2. STDD File Creation [IMPL:STDD_FILES] [ARCH:STDD_STRUCTURE] [REQ:STDD_SETUP]

### Implementation Approach:
- Created `stdd/` directory.
- Instantiated `requirements.md`, `architecture-decisions.md`, `implementation-decisions.md`, `semantic-tokens.md`, `tasks.md`, and `ai-principles.md` from templates.
- Updated `.cursorrules` to enforce STDD rules.

**Cross-References**: [ARCH:STDD_STRUCTURE], [REQ:STDD_SETUP]

## 2. State Path Resolver [IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]

### Decision: Provide a pure resolver plus bootstrap glue for persistence paths
**Rationale:**
- Implements the precedence + debug contract from [ARCH:STATE_PATH_SELECTION] without tying tests to global process state.
- Makes it trivial to inject fake environments for module validation required by [REQ:MODULE_VALIDATION].

### Implementation Approach:
- Add package `configpaths` with:
  - `const DefaultState = "~/.goful/state.json"` / `DefaultHistory = "~/.goful/history/shell"`.
  - `const EnvStateKey = "GOFUL_STATE_PATH"` / `EnvHistoryKey = "GOFUL_HISTORY_PATH"`.
  - `type Paths struct { State, History, StateSource, HistorySource string }`.
  - `type Resolver struct { LookupEnv func(string) (string, bool) }` with method `Resolve(flagState, flagHistory string) Paths`.
  - Resolver order: CLI flag → env var → default. All outputs pass through `util.ExpandPath`.
  - `func EnsureParent(path string) error` helper to call `os.MkdirAll(filepath.Dir(path), 0o755)` before state/history saves.
- Add `BootstrapPaths` helper (same package or `main.go`) that:
  - Parses CLI flags (`-state`, `-history`).
  - Calls resolver and logs `DEBUG: [IMPL:STATE_PATH_RESOLVER] ...` lines when `GOFUL_DEBUG_PATHS=1`.
  - Applies resolved paths to `app.NewGoful`, `cmdline.LoadHistory`, and the corresponding save paths when exiting.
- Update `filer.SaveState` to create parent directories before writing to satisfy the requirement.

**Code Markers**:
- Resolver + helper functions carry `[IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]` comments.
- `main.go` flag definitions include inline references to the same tokens.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Touched files: `configpaths/*.go`, `configpaths/*._test.go`, `main.go`, `filer/filer.go`, `README.md`.
- Tests named `TestResolvePaths_REQ_CONFIGURABLE_STATE_PATHS` prove module validation.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `2026-01-01`: `./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 70 token references across 40 files.`

**Cross-References**: [ARCH:STATE_PATH_SELECTION], [REQ:CONFIGURABLE_STATE_PATHS]

## 3. Error Handling Implementation [IMPL:ERROR_HANDLING] [ARCH:ERROR_HANDLING] [REQ:ERROR_HANDLING]

### Error Types
```[your-language]
// Define error types/constants for your language
// Example patterns:
// - Error constants or enums
// - Error classes or types
// - Error codes with messages
```

### Error Wrapping
```[your-language]
// Wrap errors with context in your language's idiomatic way
// Example patterns:
// - Error wrapping with context
// - Exception chaining
// - Error propagation with additional information
```

### Error Reporting
- Error logging approach
- Error propagation pattern
- User-facing error messages

## 4. Testing Implementation [IMPL:TESTING] [ARCH:TESTING_STRATEGY] [REQ:*]

**Note**: This implementation realizes the validation criteria specified in `requirements.md` and follows the testing strategy defined in `architecture-decisions.md`. Each test validates specific satisfaction criteria from requirements.

### Unit Test Structure
```[your-language]
// Unit test structure for your language
// Example pattern:
function testResolvePaths_REQ_CONFIGURABLE_STATE_PATHS() {
    // [REQ:CONFIGURABLE_STATE_PATHS] Validates configurable persistence behavior
    testCases = [
        {
            name: "test case 1",
            input: inputValue,
            expected: expectedValue
        }
    ]
    
    // Run test cases
    for each testCase in testCases {
        result = functionUnderTest(testCase.input)
        assert result equals testCase.expected
    }
}
```
> **Remember**: Without the `[REQ:*]` suffix + inline comment, this test fails `[PROC:TOKEN_AUDIT]`.

### Integration Test Structure
```[your-language]
// Integration test structure for your language
function testIntegrationScenario_REQ_CONFIGURABLE_STATE_PATHS() {
    // [REQ:CONFIGURABLE_STATE_PATHS] End-to-end validation comment
    // Setup: Prepare test environment
    // Execute: Run integration scenario
    // Verify: Assert expected outcomes
}
```
> **Log** the execution of these tests alongside your `[PROC:TOKEN_VALIDATION]` run so future audits see when behavior was last verified.

## 5. Code Style and Conventions [IMPL:CODE_STYLE]

### Naming
- Use descriptive names
- Follow language naming conventions
- Exported types/functions: PascalCase (or language equivalent)
- Unexported: camelCase (or language equivalent)

### Documentation
- Package-level documentation
- Exported function documentation
- Inline comments for complex logic
- Examples in test files

### Formatting
- Use standard formatter for chosen language
- Use linter for code quality

## 6. Module Validation Implementation [IMPL:MODULE_VALIDATION] [ARCH:MODULE_VALIDATION] [REQ:MODULE_VALIDATION]

### Decision: Independent Module Validation Before Integration
**Rationale:**
- Implements [REQ:MODULE_VALIDATION] requirement for independent module validation
- Realizes [ARCH:MODULE_VALIDATION] architecture decision
- Ensures modules are validated independently before integration to eliminate complexity-related bugs

### Implementation Approach:

#### Module Identification Phase
1. **Before Development**: Identify logical modules and document:
   - Module boundaries and responsibilities
   - Module interfaces and contracts
   - Module dependencies
   - Module validation criteria

#### Module Development Phase
2. **Develop Module Independently**:
   - Implement module according to defined interface
   - Use dependency injection or interfaces for dependencies
   - Keep module isolated from other modules during development

#### Module Validation Phase
3. **Validate Module Independently** (BEFORE integration):
   ```[your-language]
   // Example: Module validation test structure
   function testModuleName_IndependentValidation() {
       // Setup: Create module with mocked dependencies
       mockDependency = createMockDependency()
       module = createModule(mockDependency)
       
       // Test: Unit tests with mocked dependencies
       test("contract validation", function() {
           result = module.process(input)
           assert result equals expectedOutput
       })
       
       // Test: Edge cases
       test("edge cases", function() {
           // Test boundary conditions
       })
       
       // Test: Error handling
       test("error handling", function() {
           // Test error scenarios
       })
   }
   ```

4. **Validation Requirements**:
   - **Unit Tests**: Comprehensive unit tests with all dependencies mocked
   - **Contract Tests**: Validate input/output contracts
   - **Edge Case Tests**: Test boundary conditions and edge cases
   - **Error Handling Tests**: Test error scenarios and error propagation
   - **Integration Tests with Test Doubles**: Test module with mocks/stubs/fakes for dependencies

5. **Document Validation Results**:
   - Document which validation tests passed
   - Document any known limitations or assumptions
   - Mark module as "validated" only after all validation criteria pass

#### Integration Phase
6. **Integrate Validated Modules** (ONLY after validation passes):
   ```[your-language]
   // Example: Integration after module validation
   // [REQ:MODULE_VALIDATION] Only integrate after module validation passes
   // [IMPL:MODULE_VALIDATION] [ARCH:MODULE_VALIDATION] [REQ:MODULE_VALIDATION]
   function integrateModules(validatedModule1, validatedModule2) {
       // Integration code that combines validated modules
   }
   ```

7. **Integration Testing**:
   - Test combined behavior of validated modules
   - Verify integration points work correctly
   - Test end-to-end scenarios with validated modules

### Task Structure:
- **Separate Tasks**: Module development, module validation, and integration must be separate tasks
- **Task Dependencies**: Integration tasks depend on module validation tasks
- **Task Priorities**: Module validation is typically P0 or P1 priority

### Code Markers:
- Look for module validation test files: `*_module_test.[ext]` or `*_validation_test.[ext]`
- Look for integration test files: `*_integration_test.[ext]`
- Code comments: `// [REQ:MODULE_VALIDATION] Module validated independently before integration`

### Cross-References: [ARCH:MODULE_VALIDATION], [REQ:MODULE_VALIDATION]

## 7. Go Mod Update [IMPL:GO_MOD_UPDATE] [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS]

### Decision: Update `go.mod` to current LTS and align tooling
**Rationale:**
- Required to use supported compiler features and security fixes.
- Aligns local/CI builds and reduces drift.

### Implementation Approach:
- Set `go 1.24.0` plus `toolchain go1.24.3` in `go.mod` with `[IMPL:GO_MOD_UPDATE] [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS]` comments.
- Recomputed the module graph with `go mod tidy` under `go1.24.3`.
- Updated `message.log` to use a constant format string so `go vet` passes on the modern toolchain.
- Verified the upgraded toolchain by running `go fmt ./...`, `go vet ./...`, and `go test ./...` on darwin/arm64 (go1.24.3).

**Code Markers**:
- `go.mod` header and CI workflow `setup-go` steps.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Comments present in `go.mod` and workflow.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `2026-01-01 go1.24.3 darwin/arm64`
  - `go fmt ./...` (touched `info/info_unix.go`, `info/info_windows.go`)
  - `go vet ./...`
  - `go test ./...`
- `./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 12 token references across 35 files.`

**Cross-References**: [ARCH:GO_RUNTIME_STRATEGY], [REQ:GO_TOOLCHAIN_LTS]

## 8. Dependency Bump [IMPL:DEP_BUMP] [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]

### Decision: Refresh direct deps and tidy module graph
**Rationale:**
- Pull in security and bug fixes; keep compatible with Go LTS.

### Implementation Approach:
- Upgraded direct deps to current releases:
  - `github.com/gdamore/tcell/v2 v2.13.5`
  - `github.com/mattn/go-runewidth v0.0.19`
  - `github.com/google/shlex` (latest pseudo-version)
- Upgraded supporting deps:
  - `github.com/lucasb-eyer/go-colorful v1.3.0`, `github.com/gdamore/encoding v1.0.1`, `github.com/rivo/uniseg v0.4.7`
  - `golang.org/x/sys v0.39.0`, `golang.org/x/term v0.38.0`, `golang.org/x/text v0.32.0`
  - Added `github.com/clipperhouse/uax29/v2 v2.2.0` via transitive requirements from `tcell`.
- Ran `go mod tidy` to sync `go.sum`.
- No shims or breaking API adjustments were required after the upgrades.

**Code Markers**:
- `go.mod` entries and related code comments.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Dependency change commits include `[IMPL:DEP_BUMP] [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]`.

**Cross-References**: [ARCH:DEPENDENCY_POLICY], [REQ:DEPENDENCY_REFRESH]

## 9. CI Workflow [IMPL:CI_WORKFLOW] [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]

### Decision: GitHub Actions workflow for fmt/vet/tests
**Rationale:**
- Enforces consistency and prevents regressions on PRs.

### Implementation Approach:
- Added `.github/workflows/ci.yml` with a `format-vet-test` job that:
  - Checks out the repo, sets up Go `1.24.3` via `actions/setup-go@v5`, and caches modules (`go.sum`).
  - Runs `gofmt -w $(git ls-files '*.go')` followed by `git status --short` and `git diff --exit-code` to enforce formatting.
  - Executes `go vet ./...` and `go test ./...` for regression coverage.
  - Runs `./scripts/validate_tokens.sh` so every CI pass records `[PROC:TOKEN_VALIDATION]`.
- Each shell block embeds `[IMPL:CI_WORKFLOW] [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]` (or the token validation equivalents) to keep traceability in the workflow itself.

**Code Markers**:
- Workflow file comments with `[IMPL:CI_WORKFLOW] [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]`.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `2026-01-01`: `./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 19 token references across 36 files.`

**Cross-References**: [ARCH:CI_PIPELINE], [REQ:CI_PIPELINE_CORE]

## 10. Staticcheck Setup [IMPL:STATICCHECK_SETUP] [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]

### Decision: Add static analysis job
**Rationale:**
- Surface correctness issues early.

### Implementation Approach:
- Added `staticcheck` job to `.github/workflows/ci.yml` that:
  - Reuses the Go `1.24.3` toolchain setup with cached modules.
  - Installs `staticcheck` via `go install honnef.co/go/tools/cmd/staticcheck@latest`.
  - Runs `staticcheck ./...` with `[IMPL:STATICCHECK_SETUP] [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]` inline comments.

**Code Markers**:
- Workflow job comments carry `[IMPL:STATICCHECK_SETUP] [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]`.

**Cross-References**: [ARCH:STATIC_ANALYSIS_POLICY], [REQ:STATIC_ANALYSIS]

## 11. Race Job [IMPL:RACE_JOB] [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]

### Decision: Dedicated race-enabled test job
**Rationale:**
- Detects concurrency issues without impacting main job runtime.

### Implementation Approach:
- Added `race-tests` job to `.github/workflows/ci.yml` that:
  - Sets up Go `1.24.3` with caching and reuses the module cache.
  - Executes `go test -race ./...` so concurrency regressions fail CI early.

**Code Markers**:
- Workflow job comments carry `[IMPL:RACE_JOB] [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]`.

**Cross-References**: [ARCH:RACE_TESTING_PIPELINE], [REQ:RACE_TESTING]

## 12. Widget Tests [IMPL:TEST_WIDGETS] [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]

### Decision: Add unit/snapshot tests for widget primitives
**Rationale:**
- Protect rendering/event handling behaviors.

### Implementation Approach:
- **Module identification [REQ:MODULE_VALIDATION]:**
  - `widget.ListBox` (cursor + scrolling state machine)
  - `widget.Gauge` (progress rendering & normalization)
  - `widget.TextBox` (buffer editing helpers such as `InsertBytes`/`DeleteBytes`)
  - Supporting `widget.Window` helpers (column calculations, offsets)
- **Validation criteria:**
  - Pure functions (cursor math, offset adjustments, column sizing) validated with table-driven Go tests.
  - Rendering helpers validated indirectly by inspecting state (e.g., `ScrollRate`, `ColumnAdjustContentsWidth`) to avoid brittle terminal assertions; future work can introduce snapshot harnesses once deterministic screen fakes exist.
- **Immediate test plan:**
  - `TestListBoxCursorClamping_REQ_UI_PRIMITIVE_TESTS`: proves `SetCursor`, `MoveCursor`, and `SetCursorByName` respect `Lower()/Upper()` bounds and fallback semantics.  
  - `TestListBoxScrollRate_REQ_UI_PRIMITIVE_TESTS`: verifies offset math for `ScrollRate` (“Top”/percentage/“Bot”).  
  - `TestListBoxColumnAdjust_REQ_UI_PRIMITIVE_TESTS`: confirms column auto-fit honors widest content within available width.  
  - Follow-on work: add gauge fill-ratio tests and textbox editing regressions; snapshot harness will cover highlight rendering once `SetCells` fakes land.
- Tests live in `widget/*.go` so they can directly access unexported helpers while retaining `[REQ:UI_PRIMITIVE_TESTS] [ARCH:TEST_STRATEGY_UI] [IMPL:TEST_WIDGETS]` breadcrumbs for auditability.

**Code Markers**:
- Test files include `[IMPL:TEST_WIDGETS] [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]`.

**Cross-References**: [ARCH:TEST_STRATEGY_UI], [REQ:UI_PRIMITIVE_TESTS]

## 13. Command Tests [IMPL:TEST_CMDLINE] [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]

### Decision: Add tests for command parsing and modes
**Rationale:**
- Prevent regressions in command handling and state transitions.

### Implementation Approach:
- **Module identification [REQ:MODULE_VALIDATION]:**
  - `cmdline.Parser` (tokenization + quoting)
  - `cmdline.Completion` helpers (word boundary + suggestion generation)
  - `app.Mode` transitions (normal, command, prompt)
- **Validation criteria:**
  - Parser tests feed representative command strings (including quoting, globbing, multi-byte) and assert resulting structs.
  - Mode tests stimulate key-event handlers with table-driven inputs to ensure state-dependent callbacks fire.
- **Immediate test plan:**
  - `TestParseLine_REQ_CMD_HANDLER_TESTS`: ensures parser emits expected argv slices plus error handling for unterminated quotes.
  - `TestApplyModeTransitions_REQ_CMD_HANDLER_TESTS`: uses lightweight fakes to confirm `mode.GoNormal()` / `mode.GoCommand()` toggles behavior.
  - `TestCompletionFilters_REQ_CMD_HANDLER_TESTS`: validates completion filter respects prefixes + case sensitivity.
  - Additional coverage will mock `cmdline.Extmap` to exercise edge bindings before integration tests.

**Code Markers**:
- Tests carry `[IMPL:TEST_CMDLINE] [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]`.

**Cross-References**: [ARCH:TEST_STRATEGY_CMD], [REQ:CMD_HANDLER_TESTS]

## 14. Integration Flow Tests [IMPL:TEST_INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]

### Decision: Snapshot/integration tests for file operation flows
**Rationale:**
- Validate end-to-end behavior for open/navigate/rename/delete.

- **Module identification [REQ:MODULE_VALIDATION]:**
  - `app.App` orchestration (mode wiring + widget graph)
  - `filer.Workspace`/`Directory` for FS mutations and navigation
  - `widget.ListBox`/`Textbox` for active view state
- **Validation criteria:**
  - Integration fixtures create temporary directories to simulate “open directory, navigate, rename/delete” flows without touching user files.
  - Tests assert against deterministic transcripts (e.g., active path, list contents, status messages).
- **Implemented coverage:**
  - `TestFlowOpenDirectory_REQ_INTEGRATION_FLOWS` exercises `Workspace` + `Directory` when opening a new path.
  - `TestFlowNavigateRename_REQ_INTEGRATION_FLOWS` navigates into nested directories and validates rename propagation after reload.
  - `TestFlowDelete_REQ_INTEGRATION_FLOWS` removes files and confirms the directory state refreshes.
- Future enhancement: capture golden snapshots of widget buffer output once terminal fakes exist.

**Code Markers**:
- Tests annotated with `[IMPL:TEST_INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]`.

**Cross-References**: [ARCH:TEST_STRATEGY_INTEGRATION], [REQ:INTEGRATION_FLOWS]

## 15. Architecture Guide [IMPL:DOC_ARCH_GUIDE] [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]

### Decision: Write `ARCHITECTURE.md`
**Rationale:**
- Provides concise understanding of packages and data flow.
- Establishes a stable "map" before larger refactors touch keymap/menu wiring.

### Implementation Approach:
- Structure the document into:
  1. **Overview & Goals** referencing `[REQ:ARCH_DOCUMENTATION]`.
  2. **Runtime Flow** describing `main` → `configpaths.Resolver` → `app.Goful` event loop.
  3. **Module Deep Dives** (`app`, `filer`, `widget`, `cmdline`, `menu`, `look/message/info/progress`, `configpaths`, `util`).
  4. **Validation & Testing Surfaces** listing module-level tests (widgets, cmdline, integration flows, keymap baselines).
- Embed ASCII-style flow diagrams or bullet chains to highlight dependencies.
- Cross-link to requirements/architecture/implementation tokens inline to preserve STDD traceability.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- Document review ensures every section references at least one `[REQ:*]` token and the doc is linked from `README.md`.

**Code Markers**:
- Document contains `[IMPL:DOC_ARCH_GUIDE] [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]` plus relevant cross-references for each section (e.g., `[REQ:CONFIGURABLE_STATE_PATHS]`, `[ARCH:STATE_PATH_SELECTION]`).

**Cross-References**: [ARCH:DOCS_STRUCTURE], [REQ:ARCH_DOCUMENTATION], [REQ:MODULE_VALIDATION]

## 16. CONTRIBUTING Guide [IMPL:DOC_CONTRIBUTING] [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE]

### Decision: Add contributor standards document
**Rationale:**
- Aligns development workflow and review expectations.
- Documents Go / Makefile targets, CI steps, and STDD-specific requirements (semantic tokens, module validation, debug logging).

### Implementation Approach:
- Sections:
  - **Tooling & Setup** (Go LTS, `make` targets, local environment variables).
  - **Workflow Checklist** enumerating fmt → vet → test → race/staticcheck (via CI) plus manual `./scripts/validate_tokens.sh`.
  - **Semantic Token Discipline** linking to registry updates and `[PROC:TOKEN_AUDIT]`.
  - **Module Validation Expectations** referencing `[REQ:MODULE_VALIDATION]` and `KeymapBaselineSuite`.
  - **Debug Logging Policy** (retain `DEBUG:`/`DIAGNOSTIC:` outputs).
  - **Review Gate** referencing required doc/test updates before opening PRs.
- Provide copy/paste friendly command blocks (e.g., `make fmt`, `go test ./...`).
- Link to `ARCHITECTURE.md`, `README.md`, and STDD docs for quick navigation.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- Document is cross-linked from README and includes instructions to run the validation script before requesting review.

**Code Markers**:
- `CONTRIBUTING.md` includes `[IMPL:DOC_CONTRIBUTING] [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE] [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION] [REQ:MODULE_VALIDATION]`.

**Cross-References**: [ARCH:CONTRIBUTION_PROCESS], [REQ:CONTRIBUTING_GUIDE], [REQ:MODULE_VALIDATION]

## 17. Release Targets [IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]

### Decision: Automate reproducible releases via Makefile + CI + tag workflows
**Rationale:**
- Guarantees that every release binary is built with the same flags (CGO disabled, `-trimpath -ldflags "-s -w"`), stored under predictable filenames, and accompanied by SHA256 digests.
- Keeps local, CI, and tag-triggered release flows identical: `make release` locally mirrors both the CI `release-matrix` job and the GitHub Releases workflow.

### Implementation Approach:
- **Makefile Enhancements**
  - New helpers: `vet`, `lint` (fmt + vet), `clean-release`, `release`, with defaults `DIST_DIR=dist`, `RELEASE_PLATFORMS="linux/amd64 linux/arm64 darwin/arm64"`, and `SHASUM=shasum -a 256`.
  - `release` target accepts optional `PLATFORM=os/arch`; otherwise iterates `RELEASE_PLATFORMS`. For each platform it emits `DIAGNOSTIC: [IMPL:MAKE_RELEASE_TARGETS] ...`, builds `dist/goful_${os}_${arch}`, and writes `dist/goful_${os}_${arch}.sha256`.
  - Targets stamp `[IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]` in echoes so logs remain traceable.
- **CI Workflow (`release-matrix` job)**
  - Adds matrix include entries for linux/amd64, linux/arm64, darwin/arm64.
  - Step 1: checkout + setup Go 1.24.3.
  - Step 2: `make release PLATFORM=${{matrix.goos}}/${{matrix.goarch}}`.
  - Step 3: display checksum file and upload both binary + `.sha256` via `actions/upload-artifact`.
- **Release Workflow (`.github/workflows/release.yml`)**
  - Trigger: `push` tags matching `v*`.
  - Job `matrix-build` mirrors the CI matrix and runs the same `make release` command, uploading artifacts per platform.
  - Job `publish` downloads all artifacts (merged) and calls `softprops/action-gh-release` to attach binaries + `.sha256` files to the GitHub release while logging checksum contents for `ArtifactDeterminismAudit`.

### Validation Evidence `[REQ:MODULE_VALIDATION]`:
- Local command: `make release PLATFORM=$(go env GOOS)/$(go env GOARCH)` followed by `ls dist/` ensures host builds succeed.
- CI evidence: successful `release-matrix` job plus artifact uploads recorded in workflow run logs.
- Release evidence: tag push triggers the release workflow; both matrix and publish jobs must succeed and attach artifacts/digests to the release page.

**Code Markers**:
- Makefile release recipes and workflow shell blocks include `[IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]`.

**Cross-References**: [ARCH:BUILD_MATRIX], [REQ:RELEASE_BUILD_MATRIX], [REQ:MODULE_VALIDATION]

## 18. Baseline Snapshots [IMPL:BASELINE_SNAPSHOTS] [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE]

### Decision: Capture current keybindings/modes as automated baselines
**Rationale:**
- Preserve behavior ahead of refactors.
- Provide guardrails for future keymap cleanups or menu consolidation work.

### Implementation Approach:
- Implement `KeymapBaselineSuite` unit tests under `main_keymap_test.go` that:
  - Instantiate maps via `filerKeymap(nil)`, `finderKeymap(nil)`, `cmdlineKeymap(new(cmdline.Cmdline))`, `completionKeymap(new(cmdline.Completion))`, `menuKeymap(new(menu.Menu))`.
  - Assert presence of representative key chords for navigation, selection, shell execution, finder/completion movement, and exit behaviors.
  - Emit `DEBUG:` logs enumerating the verified chords for traceability.
- Introduce helper `assertKeyCoverage` to keep tests declarative and make future updates additive.
- Tag tests with `[TEST:KEYMAP_BASELINE]` alongside `[REQ:BEHAVIOR_BASELINE]` tokens.
- Keep suite pure (no widget initialization) so it runs instantly in CI.

**Validation Evidence**:
- `go test ./...` (module validation) covers the new baseline suite prior to integrating any runtime changes.
- `./scripts/validate_tokens.sh` ensures `[TEST:KEYMAP_BASELINE]` and related tokens are registered.

**Code Markers**:
- Tests include `[IMPL:BASELINE_SNAPSHOTS] [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE] [TEST:KEYMAP_BASELINE]`.

**Cross-References**: [ARCH:BASELINE_CAPTURE], [REQ:BEHAVIOR_BASELINE], [TEST:KEYMAP_BASELINE], [REQ:MODULE_VALIDATION]

## 19. Debt Tracking [IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]

### Decision: Track debt via issues and TODO annotations
**Rationale:**
- Makes risk visible and assignable.

### Implementation Approach:
- Create issue list for known pain points; annotate code with TODO + owner.
- Link debt items into documentation and tasks.

**Code Markers**:
- TODOs carry `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`.

**Cross-References**: [ARCH:DEBT_MANAGEMENT], [REQ:DEBT_TRIAGE]

## 20. Token Validation Script [IMPL:TOKEN_VALIDATION_SCRIPT] [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]

### Decision: Automate `[PROC:TOKEN_VALIDATION]` via `scripts/validate_tokens.sh`
**Rationale:**
- Contributors need a single command to prove token references are registered.
- Satisfies modernization tasks blocked on running `[PROC:TOKEN_VALIDATION]`.

### Implementation Approach:
- Added `scripts/validate_tokens.sh` (Bash, `set -euo pipefail`).
- Script requires `git` and `rg`, builds the token registry from `stdd/semantic-tokens.md`, and scans tracked source files (`*.go`, module files, shell scripts, workflows, Makefile) unless custom paths are supplied.
- Emits diagnostic output and fails if tokens found in source are missing from the registry.
- Produces success message with counts for audit trails.

**Code Markers**:
- Script header includes `[IMPL:TOKEN_VALIDATION_SCRIPT] [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]`.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `2026-01-01`: `./scripts/validate_tokens.sh` (default globs) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 12 token references across 35 files.`

**Cross-References**: [ARCH:TOKEN_VALIDATION_AUTOMATION], [REQ:STDD_SETUP], [PROC:TOKEN_VALIDATION]

