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

## 2b. External Command Loader [IMPL:EXTERNAL_COMMAND_LOADER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]

### Decision: Provide a dedicated loader that resolves the config file path (flag/env/default), parses JSON or YAML, validates schema, filters by platform, and falls back to embedded defaults.
**Rationale:**
- Keeps customization outside of `main.go`, enabling dotfile repos or tooling scripts to ship command sets without recompiling.
- Ensures `[REQ:MODULE_VALIDATION]` can be satisfied with pure unit tests (no widget/app dependencies).
- Preserves historic behavior (platform-specific defaults, cursor offsets) whenever the config file is absent or invalid.

### Implementation Approach:
- Extend `configpaths.Paths` with `Commands` + `CommandsSource` so `emitPathDebug` reports all three precedence outcomes. CLI flag `-commands` and env var `GOFUL_COMMANDS_FILE` feed the resolver.
- New package `externalcmd` exposes:
  - `type Entry` with `Menu`, `Key`, `Label`, `Command`, `Offset`, `Platforms`, `Disabled`.
  - `func Defaults(goos string) []Entry` returning the old hard-coded Windows/POSIX menus expressed with `%` macros instead of inline `g.File()` references (e.g., rename defaults to `mv -vi %f %~f`).
  - `func Load(Options) ([]Entry, error)` where `Options` carries `Path`, `GOOS`, `ReadFile`, `Debug`. Loader expands `~`, reads JSON or YAML (supporting raw arrays or `{ commands: [] }`), validates unique `menu/key`, enforces required fields, filters by `Platforms`, skips disabled entries, logs diagnostics tagged with `[IMPL:EXTERNAL_COMMAND_LOADER]`, and **prepends file entries ahead of `Defaults` unless the file sets an explicit opt-out (e.g., `inheritDefaults: false`)**.
- Errors reading/parsing the config file bubble up so callers can emit `message.Errorf` but still fall back to defaults.
- Dependency note: Introduced `gopkg.in/yaml.v3` to parse YAML files without writing a bespoke parser; the package is already widely used and compatible with Go 1.24.

**Code Markers**:
- `externalcmd/defaults.go` & `externalcmd/loader.go` include `[IMPL:EXTERNAL_COMMAND_LOADER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]`.
- `configpaths/resolver.go` references `[IMPL:STATE_PATH_RESOLVER]` while documenting the new `Commands` path fields.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Tests in `externalcmd/loader_test.go` named `TestLoadCommands_REQ_EXTERNAL_COMMAND_CONFIG` and friends ensure schema validation, platform filtering, duplicate detection, and fallback behavior.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./externalcmd` (darwin/arm64, Go 1.24.3) covers the loader in isolation, including JSON and YAML fixtures.
- `./scripts/validate_tokens.sh` (2026-01-02) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 245 token references across 52 files.`

**Cross-References**: [ARCH:EXTERNAL_COMMAND_REGISTRY], [REQ:EXTERNAL_COMMAND_CONFIG], [REQ:MODULE_VALIDATION]

## 2c. External Command Binder [IMPL:EXTERNAL_COMMAND_BINDER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]

### Decision: Extract a pure helper that converts loader entries into `menu.Add` triplets with closures that call `g.Shell`, preserving cursor offsets and surfacing placeholder entries when configs are empty.
**Rationale:**
- Keeps `main.go` readable and makes binder behavior testable without spinning up the widget stack.
- Guarantees deterministic registration order (file order) and simple hooks for future menu destinations beyond `external-command`.
- Provides user-facing feedback when no commands remain (placeholder entry says “no commands configured” and logs a `DEBUG:` line).

### Implementation Approach:
- New file `main_external_commands.go` in package `main` defines:
  - `type shellInvoker func(cmd string, offset ...int)` to abstract `g.Shell`.
  - `func buildExternalMenuSpecs(entries []externalcmd.Entry) []menuSpec` which normalizes menu names, drops entries missing required fields (defensive), and ensures file order is preserved.
  - `func registerExternalCommandMenu(g *app.Goful, entries []externalcmd.Entry)` which calls `buildExternalMenuSpecs`, adds a placeholder entry if specs are empty, and feeds `menu.Add` arguments with closures capturing the right offset.
  - Placeholder callbacks call `message.Info` so pressing `X` explains that no commands are configured instead of crashing.
- Tests in `main_external_commands_test.go` inject fake `shellInvoker` functions and assert commands/offsets propagate correctly and placeholder behavior triggers when expected.

**Code Markers**:
- `main_external_commands.go` and its tests include `[IMPL:EXTERNAL_COMMAND_BINDER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]`.
- `main.go` references `[IMPL:EXTERNAL_COMMAND_BINDER]` when wiring the menu after loading definitions.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Tests named `TestBuildExternalMenuSpecs_REQ_EXTERNAL_COMMAND_CONFIG` and `TestRegisterExternalCommandsPlaceholder_REQ_EXTERNAL_COMMAND_CONFIG` cover success and edge cases.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./...` after integration covers binder tests; `./scripts/validate_tokens.sh` ensures the new tokens are registered.

**Cross-References**: [ARCH:EXTERNAL_COMMAND_REGISTRY], [REQ:EXTERNAL_COMMAND_CONFIG], [REQ:MODULE_VALIDATION]

## 2d. External Command Append Toggle [IMPL:EXTERNAL_COMMAND_APPEND] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]

### Decision: Preserve built-in Windows/POSIX menu entries whenever a commands file is present, **prepending** file-defined entries by default and only replacing defaults when the file opts out via `inheritDefaults: false`.
**Rationale:**
- Operators expect historical shortcuts (`cp`, `mv`, etc.) to remain available unless they explicitly suppress them; this keeps onboarding friction low.
- Some environments must ship a clean slate for security reasons, so the same config file needs a deterministic switch to drop defaults entirely.
- Encoding the behavior in semantic tokens makes the inheritance contract testable and discoverable beyond requirements prose.

### Implementation Approach:
- Extend the loader parser to recognize either an array of entries (prepends defaults implicitly) or an object wrapper containing `commands` and `inheritDefaults` (JSON or YAML). Missing flags default to `true` so existing configs pick up prepend semantics automatically.
- After sanitizing file entries, merge them with `externalcmd.Defaults` when inheritance is enabled (custom entries first) or return only the sanitized entries when disabled. Emit `DEBUG: [IMPL:EXTERNAL_COMMAND_APPEND]` logs describing whether defaults were included.
- Surface the new `[IMPL:EXTERNAL_COMMAND_APPEND]` token in loader code comments, docs, and tests so audits can trace the behavior end-to-end.

**Code Markers**:
- `externalcmd/loader.go` merge logic and debug output include `[IMPL:EXTERNAL_COMMAND_APPEND] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]`.
- README + STDD docs document the `inheritDefaults` flag and reference this token.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Tests named `TestLoadAppendsDefaultsByDefault_REQ_EXTERNAL_COMMAND_CONFIG` and `TestLoadCanDisableDefaults_REQ_EXTERNAL_COMMAND_CONFIG` include `[IMPL:EXTERNAL_COMMAND_APPEND]` in comments.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./externalcmd` (darwin/arm64, Go 1.24.3) covers prepend vs. replace behaviors.
- `./scripts/validate_tokens.sh` run after implementation adds `[IMPL:EXTERNAL_COMMAND_APPEND]` to the registry and verifies references.

**Cross-References**: [ARCH:EXTERNAL_COMMAND_REGISTRY], [REQ:EXTERNAL_COMMAND_CONFIG], [REQ:MODULE_VALIDATION]

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

## 19. Window Macro Enumeration [IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]

### Decision: Introduce `%D@`/`%~D@` expansions via helper modules
**Rationale:**
- Extends external-command macros so scripts can see every workspace directory without hard-coding pane order.
- Keeps the macro parser maintainable by isolating the new behavior into two helpers instead of embedding list construction inside the switch.
- Preserves compatibility with existing quoting (`%D`/`%~D`) and escape semantics while offering deterministic ordering for automation.

### Implementation Approach:
- Added `otherWindowDirPaths(ws *filer.Workspace) []string` (Module 1 `WindowSequenceBuilder`) that iterates from `Focus+1` through all directories, wrapping via modulo arithmetic. Returns an empty slice if there is only one directory. A companion `otherWindowDirNames` helper derives the same deterministic ordering but returns `Directory.Base()` so `%d@` can reuse the sequence logic without duplicating basename handling across call sites.
- Added `formatDirListForMacro(paths []string, quote bool) string` (Module 2 `MacroListFormatter`) that applies `util.Quote` per entry when `quote=true`, leaves entries untouched when `quote=false`, and joins with single spaces. `%D@` invokes the quoted branch so each path is escaped, while `%~D@` deliberately uses the raw branch to honor the tilde modifier's non-quote semantics. `%d@` shares the same formatter for directory names so the quoting guarantees (and `%~` override) behave identically whether scripts need paths or basenames. Returns an empty string when no paths are provided.
- Updated `expandMacro` to detect `%D@`, `%~D@`, `%d@`, and `%~d@` by looking ahead for the `macroAllOtherDirs` sentinel (`'@'`). The dispatcher calls the helpers instead of reusing the single-path logic so quoted vs. raw behavior stays localized and tied to whether the tilde modifier was present.

### Code Markers:
- `app/spawn.go` helper functions and `%D@`/`%d@` branches include `// [IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]`.
- `README.md` macro table entry references the same tokens so documentation remains searchable.

### Validation Evidence `[REQ:MODULE_VALIDATION]`:
- `TestOtherWindowDirPaths_REQ_WINDOW_MACRO_ENUMERATION` (helper test) covers 1–4 directory workspaces, wrap-around behavior, and focus movement, and `TestOtherWindowDirNames_REQ_WINDOW_MACRO_ENUMERATION` mirrors that coverage for basenames.
- `TestMacroListFormatting_REQ_WINDOW_MACRO_ENUMERATION` confirms quoting vs. raw output and empty inputs.
- `TestExpandMacro` gained `%d@`/`%~d@` assertions so the integration path is covered with real macro parsing, asserting that `%D@` remains quoted while `%~D@` returns raw paths (including entries with spaces) and that `%d@`/`%~d@` only emit directory names.
- `./scripts/validate_tokens.sh` (2026-01-05) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 260 token references across 55 files.` (captured in this decision and the active task log).

**Cross-References**: [ARCH:WINDOW_MACRO_ENUMERATION], [REQ:WINDOW_MACRO_ENUMERATION], [REQ:MODULE_VALIDATION]

## 20. Debt Tracking [IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]

### Decision: Track debt via issues and TODO annotations
**Rationale:**
- Makes risk visible and assignable.
- Keeps inline breadcrumbs synchronized with the central backlog so `[PROC:TOKEN_AUDIT]` can verify coverage.

### Implementation Approach:
- Capture the backlog in `stdd/debt-log.md` (D1–D4 for the initial pass) with owner, risk, TODO reference, and next action columns.
- Annotate each hotspot (`app/goful.go`, `main.go`, `cmdline/cmdline.go`, `filer/filer.go`) with `TODO(goful-maintainers)` comments that describe the risk and include `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`.
- Update `stdd/tasks.md` to link the P1 debt triage task back to the backlog and record `[PROC:TOKEN_VALIDATION]` output after audits.

**Code Markers**:
- TODOs carry `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`.
- Documentation links refer to the backlog file so future contributors know where to extend the list.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `2026-01-01`: `./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 148 token references across 44 files.`

**Cross-References**: [ARCH:DEBT_MANAGEMENT], [REQ:DEBT_TRIAGE]

## 21. Token Validation Script [IMPL:TOKEN_VALIDATION_SCRIPT] [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]

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

## 21. Quit Dialog Return Handling [IMPL:QUIT_DIALOG_ENTER] [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]

### Decision: Map `tcell.KeyEnter` → `C-m` and guard with regression tests
**Rationale:**
- Keeps the cmdline submission shortcut stable even when upstream terminal libraries change raw key codes.
- Fixes the regression where Return no longer exited the quit dialog after dependency upgrades.

### Implementation Approach:
- Extend `keyToSting` (sic) in `widget.EventToString` to treat `tcell.KeyEnter` identically to `tcell.KeyCtrlM`.
- Add focused unit tests in `widget/widget_test.go` asserting both `KeyEnter` and `KeyCtrlM` emit the canonical `C-m` string.
- Annotate the mapping and tests with `[IMPL:QUIT_DIALOG_ENTER] [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]` comments for traceability.

**Code Markers**:
- `widget/widget.go` mapping comment at the new dictionary entry.
- `widget/widget_test.go` test names/comments referencing `[REQ:QUIT_DIALOG_DEFAULT]`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Requires verifying the translator module and accompanying tests include the `[IMPL:QUIT_DIALOG_ENTER] [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]` markers.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `2026-01-02`: `go test ./...` (darwin/arm64, Go 1.24.3 toolchain) exercising `widget` translator tests.
- `2026-01-02`: `./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 25 token references across 36 files.`

**Cross-References**: [ARCH:QUIT_DIALOG_KEYS], [REQ:QUIT_DIALOG_DEFAULT]

## 21b. Backspace Key Translation [IMPL:BACKSPACE_TRANSLATION] [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]

### Decision: Normalize both tcell backspace key codes so they emit the canonical `backspace` string consumed by filer and prompt keymaps
**Rationale:**
- macOS terminals often report Backspace as `tcell.KeyBackspace` while Linux/tmux sessions use `tcell.KeyBackspace2`. Only one entry existed in the translator map, so Backspace failed silently on half the platforms.
- Mapping both key codes to the same string preserves historical behavior (Backspace opens parent directory, deletes the previous rune) without requiring duplicate keymap entries or user-specific configuration.
- Keeping the normalization inside `widget.EventToString` satisfies `[REQ:MODULE_VALIDATION]` by addressing the issue within the pure translator module that already underpins cmdline/filer behavior.

### Implementation Approach:
- Extend `keyToSting` in `widget/widget.go` with a `tcell.KeyBackspace` entry pointing to `"backspace"` and annotate both entries with `[IMPL:BACKSPACE_TRANSLATION] [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]`.
- Add table-driven unit test `TestEventToStringBackspace_REQ_BACKSPACE_BEHAVIOR` in `widget/widget_test.go` that creates events for both `tcell.KeyBackspace` and `tcell.KeyBackspace2` and asserts `EventToString` returns `backspace`.
- Retain existing `main_keymap_test.go` baseline coverage so the `backspace` binding remains required in filer/cmdline/finder/completion keymaps.

**Code Markers**:
- `widget/widget.go` map entries for both backspace keys include `[IMPL:BACKSPACE_TRANSLATION] [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]`.
- `widget/widget_test.go` test includes the same triplet and references `[REQ:BACKSPACE_BEHAVIOR]` in the function name.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Files: `widget/widget.go`, `widget/widget_test.go`, `main_keymap_test.go`.
- Tests: `TestEventToStringBackspace_REQ_BACKSPACE_BEHAVIOR` plus existing keymap baseline cases referencing `[REQ:BACKSPACE_BEHAVIOR]` ensure cross-layer coverage.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./...` (darwin/arm64, Go 1.24.3) — 2026-01-09.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` (2026-01-09) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 520 token references across 66 files.`

**Cross-References**: [ARCH:BACKSPACE_TRANSLATION], [REQ:BACKSPACE_BEHAVIOR], [REQ:MODULE_VALIDATION]

## 22. Terminal Adapter Module [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD]

### Decision: Implement a `terminalcmd` package that encapsulates platform-aware terminal command creation plus the glue that registers it with `g.ConfigTerminal`.
**Rationale:**
- Keeps OS detection, tmux handling, and macOS automation scripts isolated from UI wiring so we can satisfy [REQ:TERMINAL_PORTABILITY] without scattering conditionals.
- Enables fast module validation—`CommandFactory` can be unit-tested with pure inputs, and `Configurator` can be exercised via fakes that capture the configured command slices.

### Implementation Approach:
- **Module 1: `CommandFactory`**
  - Signature: `func NewFactory(opts Options) Factory`.
  - `Options` include `GOOS`, `IsTmux`, `Override []string`, `Tail string`, plus macOS-specific fields `TerminalApp string` (default `Terminal`) and `TerminalShell string` (default `bash`) so AppleScript launches can be customized without editing Go code.
  - `Factory.CommandWithCwd(cmd string, cwd string) []string` returns:
    - Override path: `Override + []string{"bash", "-c", payload}` where `payload` already prefixes macOS commands with `cd "<cwd>";` to satisfy `[REQ:TERMINAL_CWD]`.
    - Tmux path: `[]string{"tmux", "new-window", "-n", title(cmd), cmd + tail}`.
    - macOS path: `[]string{"osascript", "-e", fmt.Sprintf("tell application \"%s\" to do script \"%s\"", terminalApp, script), "-e", fmt.Sprintf("tell application \"%s\" to activate", terminalApp)}` where `script` embeds the configured shell (`<terminalShell> -c "cd \"<cwd>\"; <cmd + tail>"; exit`).
    - Linux default: maintain current gnome-terminal invocation with title-setting escape.
  - Emits `DEBUG: [IMPL:TERMINAL_ADAPTER] ...` logs describing the branch taken and any overrides, guarded by `GOFUL_DEBUG_TERMINAL=1`.

- **Module 2: `Configurator`**
  - Accepts a `Factory` and returns the closure passed to `g.ConfigTerminal`.
  - Responsible for injecting the “HIT ENTER KEY” tail, escaping titles, ensuring `bash -c` semantics remain unchanged, and providing a live `cwd` callback so macOS launches always reflect the focused directory.
  - Surface helper `Apply(cfg Configurator, factory Factory, cwd func() string)` that wires both shell and terminal commands where appropriate.

- **Environment & Overrides**
  - Parse `GOFUL_TERMINAL_CMD` (string) or `-terminal` flag (future) into the override slice.
  - Read `GOFUL_TERMINAL_APP` and `GOFUL_TERMINAL_SHELL` (with defaults baked into `NewFactory`) so AppleScript launches can target another application or shell without modifying Go code.
  - Document how to supply fallback commands (e.g., `iTerm2`).
- **macOS Shell Invocation Safeguard**
  - [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] Switching the AppleScript payload from `bash -lc` to `<shell> -c` prevents login-shell initialization from hanging Terminal windows that source interactive profiles while respecting the configured shell binary.
  - [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] The payload is still quoted via `strconv.Quote` so `<shell> -c` receives the entire `cd "<cwd>"; <cmd><tail>` sequence intact while avoiding `.bash_profile` prompts.

**Code Markers**:
- `terminalcmd/*.go` and `main.go` wiring include `[IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD]`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Tests named `TestCommandFactoryDarwin_REQ_TERMINAL_PORTABILITY`, etc., assert branch outputs.
- Manual validation checklist is documented as `[PROC:TERMINAL_VALIDATION]` in `stdd/processes.md` and linked from `stdd/tasks.md`.

**Validation Evidence (2026-01-04)** `[PROC:TOKEN_VALIDATION]`:
- `go test ./terminalcmd` (darwin/arm64, Go 1.24.3) exercises `TestCommandFactory*`, `TestParseOverride`, and `TestApply*`, covering every selection branch plus the configurator glue.
- `./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 245 token references across 52 files.`
- Manual execution follows `[PROC:TERMINAL_VALIDATION]`; operators must run the macOS/Linux checklist on physical hardware prior to releases and record findings in `stdd/tasks.md`.

**Cross-References**: [ARCH:TERMINAL_LAUNCHER], [REQ:TERMINAL_PORTABILITY], [REQ:TERMINAL_CWD], [REQ:MODULE_VALIDATION]

## 23. Event Loop Shutdown Controller [IMPL:EVENT_LOOP_SHUTDOWN] [ARCH:EVENT_LOOP_SHUTDOWN] [REQ:EVENT_LOOP_SHUTDOWN]

### Decision: Add explicit stop control to the `app.Goful` event poller so goroutines exit immediately when the UI shuts down.
**Rationale:**
- Current implementations defer to process exit; embedded workflows keep the process alive, so leaked poller goroutines chew CPU and keep writing to `g.event`.
- Providing a controllable shutdown path lets us validate the poller independently and clear Debt Log item D1.

### Implementation Approach:
1. **Poller Abstraction**
   - Introduce an interface (or pure helper) `type Poller interface { Poll(stop <-chan struct{}, out chan<- tcell.Event) }` that wraps `widget.PollEvent` and listens for a stop channel.
   - Use `tcell.Screen` mocks in tests to simulate events and blocked reads.

2. **Shutdown Controller**
   - Extend `app.Goful` with a `pollStop chan struct{}` and `sync.WaitGroup` to track poller goroutines.
   - When `Run` exits (or when a fatal error occurs), close `pollStop`, wait for the poller to return with a timeout, then close `g.event`.
   - Emit `DEBUG: [IMPL:EVENT_LOOP_SHUTDOWN] stop signal sent/received` logs gated by `GOFUL_DEBUG_EVENTLOOP` (new env var) for troubleshooting.

3. **Timeout & Error Handling**
   - If the poller fails to stop within the timeout, log `message.Errorf` with instructions to file a bug referencing `[REQ:EVENT_LOOP_SHUTDOWN]` and continue teardown safely.
   - Ensure repeated shutdown attempts are idempotent (closing an already-closed channel must not panic).

4. **Integration with Existing Flow**
   - Wire the poller start/stop into existing `Run` lifecycle, ensuring other modules (cmdline, filer) still receive events while the UI is running.
   - Update debt log entry D1 to reflect mitigation once manual validation is recorded.

**Code Markers**:
- `app/goful.go`, any new helper file, and associated tests include `[IMPL:EVENT_LOOP_SHUTDOWN] [ARCH:EVENT_LOOP_SHUTDOWN] [REQ:EVENT_LOOP_SHUTDOWN]` comments.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Tests named `TestEventPollerStops_REQ_EVENT_LOOP_SHUTDOWN` (unit) and `TestRunStopsPoller_REQ_EVENT_LOOP_SHUTDOWN` (integration) prove both modules.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]` (to be captured when implementation lands):
- `go test ./app` covering poller + shutdown controller unit tests.
- `./scripts/validate_tokens.sh` updated counts recorded in `stdd/tasks.md` and this section once work completes.

**Cross-References**: [ARCH:EVENT_LOOP_SHUTDOWN], [REQ:EVENT_LOOP_SHUTDOWN], [REQ:MODULE_VALIDATION]

## 24. Xform CLI Script [IMPL:XFORM_CLI_SCRIPT] [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]

### Decision: Implement `scripts/xform.sh` as a portable Bash helper that can be executed directly or sourced, exposing a `xform` function that inserts `--to` before every destination argument while offering a dry-run preview.
**Rationale:**
- Provides a single, well-tested place to handle argument parsing and quoting for workflows that need the `--to <target>` pattern repeated for multiple paths.
- Keeps the helper compatible with macOS `/bin/bash` 3.2 by avoiding Bash 4+ features (no associative arrays or `local -n`) so contributors can run it without Homebrew Bash.

### Implementation Approach:
- `scripts/xform.sh` starts with `set -euo pipefail` and defines:
  - `xform::usage` — prints help text and exits with status 64 when invoked incorrectly.
  - `xform::parse` — consumes `-p/--prefix`, `-k/--keep`, `-n/--dry-run`, `-h/--help`, and `--` flags, ensures at least `keep + 1` positional arguments remain, and exports globals (`XFORM_PREFIX`, `XFORM_KEEP`, `XFORM_DRY_RUN`, `XFORM_ARGS`) for downstream logic. Defaults: prefix `--to`, keep `2`.
  - `xform::run` — builds the transformed argv array by preserving the first `keep` positional arguments and interleaving `<prefix>` between the remaining entries. When `dry-run` is enabled it prints the `%q`-formatted command; otherwise it executes the new argv and propagates the exit code.
  - `xform` — public wrapper that calls `xform::parse` followed by `xform::run`.
- The script checks `[[ "${BASH_SOURCE[0]}" == "$0" ]]` to decide whether to execute immediately or just define the function for callers who `source` it.
- `scripts/xform_test.sh` sources the helper and runs two module-validation suites:
  - Parser tests feed different flag combinations (custom prefix/keep along with error paths) and assert correct exit codes/messages for insufficient arguments.
  - Builder tests run `xform -n ...` and check the printed, quoted command to ensure interleaving/logging semantics work for arguments with spaces and custom prefixes/keep windows.

**Code Markers**:
- `scripts/xform.sh` contains `# [IMPL:XFORM_CLI_SCRIPT] [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]` comments near the parser and builder functions.
- `scripts/xform_test.sh` comments reference `[REQ:CLI_TO_CHAINING]` when asserting dry-run output and error handling.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Files: `scripts/xform.sh`, `scripts/xform_test.sh`, `scripts/xform.bats`.
- Tests: shell harness cases (`test_dry_run_inserts_targets_REQ_CLI_TO_CHAINING`, etc.) and Bats specs (`dry-run uses default prefix [REQ:CLI_TO_CHAINING]`, `invalid keep value fails with guidance [REQ:CLI_TO_CHAINING]`).

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `bash scripts/xform_test.sh` (2026-01-06, macOS 15.1 arm64) validates parser and builder modules independently.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` (2026-01-06) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 269 token references across 55 files.` (recorded in `stdd/tasks.md`).

**Cross-References**: [ARCH:XFORM_CLI_PIPELINE], [REQ:CLI_TO_CHAINING], [REQ:MODULE_VALIDATION]

## 25. Startup Directory Parser & Seeder [IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]

### Decision: Implement helpers that convert positional CLI arguments into deterministic workspace layouts before the UI loop runs.
**Rationale:**
- Keeps startup customization encapsulated, making it easy to validate and maintain without scattering logic across `main.go` and filer internals.
- Provides clear debug output (`GOFUL_DEBUG_WORKSPACE=1`) so operators and CI scripts can diagnose mismatched directories quickly.
- Honors `[REQ:MODULE_VALIDATION]` by keeping the parser pure and the seeder isolated from runtime event handling.

### Implementation Approach:
- `ParseStartupDirs(args []string) ([]string, []string)` (in `app/startup_dirs.go`):
  - Trims whitespace, expands `~`, resolves absolute clean paths, and checks existence + directory-ness via `os.Stat`.
  - Returns ordered directories (duplicates allowed intentionally) plus warnings describing invalid entries; warnings are surfaced through `message.Errorf`.
- `SeedStartupWorkspaces(g *app.Goful, dirs []string, debug bool) bool`:
  - Early-exits when no directories are provided to preserve historical state restoration.
  - Adds or removes workspaces to match the requested count by calling existing `CreateWorkspace` / `CloseWorkspace` helpers.
  - For each pane, focuses the first directory, calls `Dir().Chdir()` + `ReloadAll()`, retitles the workspace, and optionally logs `DEBUG:` entries tagged with `[IMPL:WORKSPACE_START_DIRS]`.
  - Returns a boolean indicating whether seeding occurred so callers can decide whether additional fallback work is needed.
- `main.go` integration:
  - After parsing flags and loading history, `flag.Args()` are fed into the parser, warnings produce `message.Errorf` output, and seeding runs with debug mode tied to `GOFUL_DEBUG_WORKSPACE`.

**Code Markers**:
- `app/startup_dirs.go`, `app/startup_dirs_test.go`, and the new block in `main.go` include `[IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]` comments.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `main.go`, `app/startup_dirs.go`, and associated tests/reference docs.
- Tests: `TestParseStartupDirs_REQ_WORKSPACE_START_DIRS`, `TestSeedStartupWorkspaces_REQ_WORKSPACE_START_DIRS`.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./...` (darwin/arm64, Go 1.24.3) on 2026-01-07 validates parser/seeder helpers plus `main.go` wiring.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 288 token references across 58 files.` (2026-01-07)

**Cross-References**: [ARCH:WORKSPACE_BOOTSTRAP], [REQ:WORKSPACE_START_DIRS], [REQ:MODULE_VALIDATION]

## 25. Startup Directory Parser & Seeder [IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]

### Decision: Parse trailing CLI arguments into deterministic workspace targets and seed filer windows before the UI loop begins.
**Rationale:**
- Keeps positional directory support encapsulated so tests can validate parsing and workspace mutation independently per `[REQ:MODULE_VALIDATION]`.
- Preserves backward compatibility by falling back to the historical workspace flow when no arguments are supplied or when every supplied path fails validation.
- Provides structured diagnostics (`GOFUL_DEBUG_WORKSPACE=1`) that help operators troubleshoot mismatched layouts without attaching a debugger.

### Implementation Approach:
- Added `applyStartupDirs(g *app.Goful, cwd func() string, args []string)` in `main.go` which:
  - Invokes `parseStartupDirs(args)` to normalize positional arguments (trim whitespace, expand `~`, preserve order—including intentional duplicates) and collect validation warnings for nonexistent paths.
  - Emits `DEBUG: [IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS] parsed=<dirs> warnings=<count>` when `GOFUL_DEBUG_WORKSPACE` is set.
  - Calls `seedWorkspace(g, dirs)` only when at least one valid directory remains.
- Introduced `parseStartupDirs` + `seedWorkspace` helpers (new file under `app/`), where:
  - Parser ensures each directory exists (using `os.Stat`), records warnings for missing items, and returns the sanitized slice.
  - Seeder inspects `g.Workspace()` to determine current pane count, then:
    - Calls `g.CreateWorkspace()` until panes reach the requested length, setting each window via `workspace.MoveFocus` + `g.Dir().Chdir(path)`.
    - When more windows exist than directories, closes excess panes via `g.CloseWorkspace()` while retaining at least one window.
    - Ensures window ordering matches the CLI order and focuses the first requested directory.
  - Each mutation path logs `DEBUG: ... seeding window` lines under the same token triplet when debug mode is enabled.
- Nonexistent directories trigger `message.Errorf` with actionable text but do not abort seeding the remaining entries.
- Default (no positional args) path is unchanged; helper returns early without modifying workspaces.

**Code Markers**:
- `main.go`, `app/startup_dirs.go`, and supporting tests include `// [IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]`.
- Debug output and error messages reference `[REQ:WORKSPACE_START_DIRS]` so operators can trace the behavior back to the requirement.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `main.go`, `app/startup_dirs.go`.
- Tests: `app/startup_dirs_test.go` (parser + seeder unit tests) and `filer/integration_test.go` (launch integration) include `[REQ:WORKSPACE_START_DIRS]` markers.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./app ./filer` (darwin/arm64, Go 1.24.3) covering parser/seeder units plus integration updates (`2026-01-07`).
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` (`2026-01-07`) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 273 token references across 55 files.`

**Cross-References**: [ARCH:WORKSPACE_BOOTSTRAP], [REQ:WORKSPACE_START_DIRS], [REQ:MODULE_VALIDATION]

## 26. Filename Exclude Rules [IMPL:FILER_EXCLUDE_RULES] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]

### Decision: Centralize basename filtering inside `filer` so every reader (default, glob, finder) automatically skips excluded entries.
**Rationale:**
- Keeps the filter logic deterministic, testable, and shareable across `Directory` instances without duplicating conditionals.
- Supports `[REQ:MODULE_VALIDATION]` by isolating the rule store from UI wiring, enabling pure unit tests for toggle/state transitions.
- Ensures mark, finder, and macro flows inherit the same behaviour because they all append via `Directory.read`.

### Implementation Approach:
- Add `filer/exclude.go` with:
  - `type excludeSet map[string]struct{}` stored alongside `excludedNamesMu sync.RWMutex`, `excludedNames excludeSet`, and `excludeEnabled bool`.
  - `func ConfigureExcludedNames(names []string, activate bool)` that trims whitespace, lowercases entries, replaces the set, and toggles `excludeEnabled = activate && len(set) > 0`.
  - `func ToggleExcludedNames() (enabled bool, hasRules bool)` plus helpers `ExcludedNamesEnabled()` and `ExcludedNameCount()` for diagnostics/UI integration.
  - `func shouldExclude(name string) bool` used by `Directory.read` (skips once `excludeEnabled` is true and the lowercase basename exists in the set).
- Guard mark insertion: the callback inside `Directory.read` checks `shouldExclude(fs.Name())` before `AppendList`, so `defaultReader`, `globPattern`, `globDirPattern`, and finder flows automatically inherit the filter.
- Emit `DEBUG: [IMPL:FILER_EXCLUDE_RULES] ...` logs when `ConfigureExcludedNames` replaces the set or when toggles occur with `GOFUL_DEBUG_PATHS` to ease troubleshooting.

**Code Markers**:
- `filer/exclude.go`, `filer/directory.go`, and filer tests include `[IMPL:FILER_EXCLUDE_RULES] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `filer/integration_test.go` gains `TestExcludedNamesHidden_REQ_FILER_EXCLUDE_NAMES` plus helper tests to prove toggling reintroduces entries and that finder/glob reuse the same guard.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./filer` (darwin/arm64, Go 1.24.3) on 2026-01-07 covering the new unit + integration cases referenced above.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` (2026-01-07) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 302 token references across 58 files.` (captured in `stdd/tasks.md` once implementation lands).

**Cross-References**: [ARCH:FILER_EXCLUDE_FILTER], [REQ:FILER_EXCLUDE_NAMES], [REQ:MODULE_VALIDATION]

## 27. Filename Exclude Loader & Toggle [IMPL:FILER_EXCLUDE_LOADER] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]

### Decision: Extend the existing path resolver with `-exclude-names` / `GOFUL_EXCLUDES_FILE` and wire a loader + toggle UI hook into `main`.
**Rationale:**
- Reuses the proven precedence model so operators immediately understand how to override the list location.
- Keeps parsing logic (trim, comment skip, case normalization) pure and testable.
- Provides a discoverable runtime toggle via both the View menu and a dedicated keystroke so users can quickly inspect hidden files when necessary.

### Implementation Approach:
- `configpaths/resolver.go` adds `DefaultExcludesPath`, `EnvExcludesKey`, and `Excludes`/`ExcludesSource` fields on `Paths`, plus a new `flagExcludesSourceLabel`. `Resolver.Resolve` now accepts `flagExcludes` and returns the resolved path + provenance so `emitPathDebug` can log it.
- `main.go` defines `excludeNamesFlag`, calls `pathsResolver.Resolve(*stateFlag, *historyFlag, *commandsFlag, *excludeNamesFlag)`, and invokes `loadExcludedNames(paths.Excludes)` before `app.SeedStartupWorkspaces`.
- `loadExcludedNames` (new helper in `main.go`) opens the file (tolerates `os.ErrNotExist`), reads newline-delimited basenames, strips comments (`#` prefix) and whitespace, lowercases entries, and calls `filer.ConfigureExcludedNames(parsed, true)`. Errors use `message.Errorf` referencing `[REQ:FILER_EXCLUDE_NAMES]`; success paths log `message.Infof` counts.
- Toggle helper `toggleExcludedNames(g *app.Goful)` wraps `filer.ToggleExcludedNames`, reports the new state/count via `message.Infof`, and calls `g.Workspace().ReloadAll()`. Bound to `g.AddKeymap("E", toggle)` and added as `view` menu entry (e.g., `n` for "toggle excludes") so the action is reachable via mouse/keyboard menus.
- `emitPathDebug` gains the excludes tuple so `GOFUL_DEBUG_PATHS=1` prints provenance for the new file.

**Code Markers**:
- `main.go`, `configpaths/resolver.go`, and `configpaths/resolver_test.go` reference `[IMPL:FILER_EXCLUDE_LOADER] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]`.
- Toggle handlers include `[REQ:FILER_EXCLUDE_NAMES]` in logged messages for auditability.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `main/exclude_loader_test.go` (new) validates parser behaviour and log emissions with `[REQ:FILER_EXCLUDE_NAMES]`.
- `configpaths` tests updated to ensure the resolver reports the excludes path precedence.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./configpaths ./main` (darwin/arm64, Go 1.24.3) on 2026-01-07 covering loader + resolver updates.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` (2026-01-07) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 302 token references across 58 files.` (same run logged under task + decision).

**Cross-References**: [ARCH:FILER_EXCLUDE_FILTER], [REQ:FILER_EXCLUDE_NAMES], [REQ:MODULE_VALIDATION]

## 28. Comparison Color Configuration [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]

### Decision: Load comparison color scheme from YAML with sensible defaults for missing/invalid configs.
**Rationale:**
- Allows users to customize colors to match their terminal themes without code changes.
- Provides consistent defaults that work on both light and dark terminals.
- Reuses the existing path resolver pattern for flag/env/default precedence.

### Implementation Approach:
- Add `filer/comparecolors/` package with:
  - `type Config struct` containing color definitions for each comparison state (NamePresent, SizeEqual, SizeSmallest, SizeLargest, TimeEqual, TimeEarliest, TimeLatest).
  - `func Load(path string) (*Config, error)` that reads YAML, validates color names, and returns parsed config.
  - `func DefaultConfig() *Config` providing sensible defaults when file is missing.
  - Color names map to `tcell.Color` values via a lookup table supporting named colors ("red", "green", etc.) and hex codes.
- Extend `configpaths.Resolver` with `-compare-colors` flag, `GOFUL_COMPARE_COLORS` env, and default `~/.goful/compare_colors.yaml`.
- `main.go` loads config at startup, passes to `look.ConfigureComparisonColors()`.

**Code Markers**:
- `filer/comparecolors/config.go` includes `[IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `filer/comparecolors/config_test.go` validates YAML parsing, defaults, and invalid input handling with `[REQ:FILE_COMPARISON_COLORS]`.

**Cross-References**: [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS], [REQ:MODULE_VALIDATION]

## 29. File Comparison Index [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]

### Decision: Build a cached index of cross-directory file comparison states for O(1) draw-time lookup.
**Rationale:**
- Comparison must not block directory reading or initial display.
- Index built once after all directories load, cached until invalidation events.
- Pure function design enables independent module validation.

### Implementation Approach:
- Add `filer/compare.go` with:
  - `type CompareState struct { NamePresent bool; SizeState SizeCompare; TimeState TimeCompare }` where `SizeCompare` and `TimeCompare` are enums (Equal, Smallest, Largest / Equal, Earliest, Latest).
  - `type ComparisonIndex struct` with `cache map[string]map[int]CompareState` keyed by filename then dirIndex.
  - `func BuildIndex(dirs []*Directory) *ComparisonIndex` that:
    - Collects all files by name across directories.
    - For files in multiple directories: marks NamePresent=true, computes size/time comparisons.
    - For single-directory files: no entry (returns nil on lookup).
  - `func (idx *ComparisonIndex) Get(dirIndex int, filename string) *CompareState` for draw-time lookup.
  - `var comparisonEnabled bool` and `func ToggleComparisonColors() (enabled bool)` for runtime toggle.
- Workspace tracks `*ComparisonIndex` and rebuilds on invalidation events (`Chdir`, `reload`, `ReloadAll`, `CreateDir`, `CloseDir`).
- Index building happens after `ReloadAll` completes, before next draw cycle.

**Code Markers**:
- `filer/compare.go` includes `[IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- `filer/compare_test.go` tests index building with various window/file combinations, size/time edge cases with `[REQ:FILE_COMPARISON_COLORS]`.

**Cross-References**: [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS], [REQ:MODULE_VALIDATION]

## 30. Comparison Draw Integration [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]

### Decision: Extend `FileStat.Draw()` to accept comparison context and apply colors independently to name, size, and time fields.
**Rationale:**
- Minimal change to existing draw path—comparison is optional context.
- Independent color application means name, size, and time can each show different comparison states.
- Respects existing file-type colors when comparison is disabled or file is unique.
- **Default State (2026-01-11)**: Comparison coloring is **enabled by default** (`comparisonEnabled = true` in `look/comparison.go`) so users immediately see color-coded file listings without manual toggle. Press `` ` `` (backtick) to disable if desired.

### Implementation Approach:
- Add `look/comparison.go` with:
  - Thread-safe style storage for each comparison state.
  - `func CompareNamePresent() tcell.Style`, `func CompareSizeEqual() tcell.Style`, etc.
  - `func ConfigureComparisonColors(cfg *comparecolors.Config)` to apply loaded config.
- Modify `FileStat.Draw(x, y, width int, focus bool)` signature to accept optional `*CompareState`:
  - New signature: `Draw(x, y, width int, focus bool, cmp *CompareState)`.
  - When `cmp != nil` and `cmp.NamePresent`: use comparison name color.
  - Size field uses `cmp.SizeState` to select Equal/Smallest/Largest color.
  - Time field uses `cmp.TimeState` to select Equal/Earliest/Latest color.
- `Directory.drawFiles()` passes comparison state from workspace index to each `FileStat.Draw()`.
- `Workspace.Draw()` ensures index is available before drawing.

**Code Markers**:
- `filer/file.go`, `filer/directory.go`, `filer/workspace.go` include `[IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`.
- `look/comparison.go` includes `[IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Integration tests in `filer/` validate that comparison colors apply correctly with `[REQ:FILE_COMPARISON_COLORS]`.

**Cross-References**: [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS], [IMPL:FILE_COMPARISON_INDEX], [IMPL:COMPARE_COLOR_CONFIG]

## 32. Linked Navigation Implementation [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]

### Decision: Implement linked navigation with minimal state in `app.Goful` and pure navigation helpers in `filer.Workspace`.
**Rationale:**
- Keeps the feature self-contained and easy to test independently per `[REQ:MODULE_VALIDATION]`.
- Uses existing `Directory.Chdir` infrastructure rather than duplicating navigation logic.
- Provides clear toggle feedback via message and header indicator.

### Implementation Approach:
- **State Management (`app/goful.go`)**:
  - Add `linkedNav bool` field to `Goful` struct (default `true`).
  - Add `func (g *Goful) ToggleLinkedNav() bool` that flips the state and returns new value.
  - Add `func (g *Goful) IsLinkedNav() bool` getter.
  - Export `LinkedNavEnabled` callback type for header rendering.

- **Navigation Helpers (`filer/workspace.go`)**:
  - Add `func (w *Workspace) ChdirAllToSubdir(name string)` that iterates all non-focused directories, checks if `name` exists as a subdirectory, and calls `Chdir(name)` if so.
  - Add `func (w *Workspace) ChdirAllToParent()` that iterates all directories and calls `Chdir("..")`.
  - Add `func (w *Workspace) SortAllBy(typ SortType)` that applies the given sort type to all directories.
  - All methods rebuild comparison index after changes.

- **Exported Sort Types (`filer/directory.go`)**:
  - Export `SortType` and constants (`SortName`, `SortNameRev`, `SortSize`, `SortSizeRev`, `SortMtime`, `SortMtimeRev`, `SortExt`, `SortExtRev`) for linked sort synchronization.
  - Export `SortBy(typ SortType)` method to enable workspace-level sorting.

- **Header Indicator (`filer/filer.go`)**:
  - Add `var linkedNavIndicatorFunc func() bool` package variable.
  - Add `func SetLinkedNavIndicator(fn func() bool)` to wire the callback from main.
  - Modify `drawHeader()` to show `[LINKED]` with reverse style when the callback returns true.

- **Keymap Integration (`main.go`)**:
  - Replace direct navigation callbacks with wrappers that check `g.IsLinkedNav()`.
  - For `backspace`/`C-h`/`u`: if linked, call `g.Workspace().ChdirAllToParent()` then `g.Dir().Chdir("..")`.
  - For enter-dir (extmap `.dir`): if linked, call `g.Workspace().ChdirAllToSubdir(name)` then `g.Dir().EnterDir()`.
  - For sort menu: if linked, call `g.Workspace().SortAllBy(typ)` instead of `g.Dir().Sort*()`.
  - Add `L` (uppercase, macOS-compatible) and `M-l` bindings to toggle with `message.Infof` feedback.
  - Wire `filer.SetLinkedNavIndicator(g.IsLinkedNav)` at startup.

**Code Markers**:
- `app/goful.go`: `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `filer/workspace.go`: `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `filer/filer.go`: `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `main.go`: `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `app/goful.go`, `filer/workspace.go`, `filer/filer.go`, `filer/directory.go`, `main.go`.
- Tests: `filer/integration_test.go` tests named `TestChdirAllToSubdir_REQ_LINKED_NAVIGATION`, `TestChdirAllToParent_REQ_LINKED_NAVIGATION`, `TestSortAllBy_REQ_LINKED_NAVIGATION`.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./...` (darwin/arm64, Go 1.24.3) on 2026-01-09 validates linked navigation and sort helpers.
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` (2026-01-09) → token validation passed.

**Cross-References**: [ARCH:LINKED_NAVIGATION], [REQ:LINKED_NAVIGATION], [REQ:MODULE_VALIDATION]

## 31. Digest Comparison [IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]

### Decision: On-demand xxHash64 digest calculation for files with equal sizes across directories
**Rationale:**
- Files may have identical sizes but different content; digest comparison provides content verification.
- On-demand calculation avoids I/O overhead for files the user doesn't need to compare.
- xxHash64 offers excellent speed for non-cryptographic hashing, suitable for file comparison.
- Terminal attributes (underline/strikethrough) provide visual distinction without adding new color configuration.

### Implementation Approach:
- Add `DigestCompare` enum to `filer/compare.go` with states: `DigestUnknown`, `DigestEqual`, `DigestDifferent`, `DigestNA`.
- Add `DigestState` field to `CompareState` struct.
- Implement `CalculateFileDigest(path string) (uint64, error)` using `github.com/cespare/xxhash/v2` with streaming I/O.
- Add `UpdateDigestStates(filename string, dirs []*Directory) int` method to `ComparisonIndex`:
  - Only processes files with `SizeState == SizeEqual`.
  - Calculates digests for all matching files across directories.
  - Sets `DigestState` to `DigestEqual` if all digests match, `DigestDifferent` otherwise.
- Add `CalculateDigestForFile(filename string) int` method to `Workspace` as public API.
- Modify `FileStat.DrawWithComparison()` to apply terminal attributes to size field:
  - `DigestEqual`: `tcell.AttrUnderline`
  - `DigestDifferent`: `tcell.AttrStrikeThrough`
- Bind `=` keystroke to trigger digest calculation for the file under cursor.
- Add "calculate file digest" entry to View menu for discoverability.

**Code Markers**:
- `filer/compare.go`: `[IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`
- `filer/workspace.go`: `[IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`
- `filer/file.go`: `[IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`
- `main.go`: `[IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Unit tests in `filer/compare_test.go` validate digest calculation and state propagation with `[REQ:FILE_COMPARISON_COLORS]`.

**Cross-References**: [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS], [IMPL:FILE_COMPARISON_INDEX], [IMPL:COMPARISON_DRAW]

## 33. Difference Search Implementation [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]

### Decision: Implement a two-command difference search with state in Workspace, pure comparison logic in a dedicated module, and a persistent status display.
**Rationale:**
- Separating state (initial dirs, active flag, status fields) from comparison logic keeps the code testable per `[REQ:MODULE_VALIDATION]`.
- Using cursor position as the implicit bookmark for "Continue" simplifies state and integrates naturally with existing navigation.
- Alphabetic iteration (case-sensitive) matches user expectations for file sorting.
- A dedicated persistent status line provides continuous feedback during long searches without auto-dismissing like regular messages.

### Implementation Approach:
- **State Management (`filer/diffsearch.go`)**:
  - Add `DiffSearchState` struct with `InitialDirs []string`, `Active bool`, and status fields:
    - `LastDiffName string` - Name of last found difference
    - `LastDiffReason string` - Reason for last difference
    - `CurrentPath string` - Current directory being searched
    - `FilesChecked int` - Count of files checked
    - `Searching bool` - Whether actively searching vs paused
  - Add setter methods: `SetSearching()`, `SetCurrentPath()`, `IncrementFilesChecked()`, `SetLastDiff()`.
  - Add `StatusText() string` that returns formatted status for display.
  - Store state in `Workspace` struct as `diffSearch *DiffSearchState`.

- **Core Comparison Logic (`filer/diffsearch.go`)**:
  - Add `func CollectAllNames(dirs []*Directory) []string` that returns the union of all file/directory names (excluding `..`), sorted alphabetically.
  - Add `func CheckDifference(name string, dirs []*Directory) (isDiff bool, reason string)`:
    - Check presence in each directory.
    - If missing from any, return `true, "missing in window N"`.
    - If present in all, compare sizes. If sizes differ, return `true, "size mismatch"`.
    - Otherwise return `false, ""`.
  - Add `func FindNextDifference(dirs []*Directory, startAfter string) (name string, reason string, found bool)`:
    - Call `CollectAllNames`, iterate from after `startAfter`.
    - Return first difference found.

- **Subdirectory Descent (`filer/diffsearch.go`)**:
  - Add `func FindNextSubdir(dirs []*Directory, startAfter string) (name string, existsInAll bool, found bool)`:
    - Collect subdirectory names only.
    - Find next subdir after `startAfter`.
    - Check if it exists in all directories.
  - Add `func FindNextSubdirInAll(dirs []*Directory, startAfter string) (name string, found bool)`:
    - Like `FirstSubdirInAll` but respects `startAfter` parameter.
    - Iterates through subdirectories in alphabetical order starting after `startAfter`.
    - Returns the first subdirectory that exists in ALL directories.
    - Critical for maintaining search position when user manually navigates into subdirectories.
  - Descent logic in command wrapper: if no file differences found, check subdirs. If subdir missing in any, treat as difference. If subdir exists in all, navigate all windows into it and repeat.
  - **Bug Fix (2026-01-10)**: The original implementation used `FirstSubdirInAll` which always returned the first common subdirectory, ignoring the `startAfter` position. This caused the search to loop back to already-searched directories when the user manually navigated into a subdirectory. Fixed by using `FindNextSubdirInAll` which respects the current search position.

- **Cursor Movement (`filer/workspace.go`)**:
  - Add `func (w *Workspace) SetCursorByNameAll(name string)` that moves cursor to `name` in all directories where it exists.
  - Reuses existing `Directory.SetCursorByName(name)` method.

- **Persistent Status Display (`diffstatus/diffstatus.go`)**:
  - New package providing a dedicated status line that persists while diff search is active.
  - `Init()` creates the status window at height-3.
  - `SetStatusFn(fn func() string)` and `SetActiveFn(fn func() bool)` wire up callbacks.
  - `SetMessage(text string)` sets a custom status message that takes priority over the `statusFn` callback. Use this for persistent status updates like "Different: X - Y" that should not auto-dismiss.
  - `ClearMessage()` clears the custom status message, reverting to `statusFn` callback.
  - `IsActive() bool` checks if status line should be shown.
  - `Draw()` renders the status line with reverse video styling, prioritizing `customMessage` over `statusFn()` result.
  - `Resize()` adjusts position during window resize.
  - Integrated into `app/goful.go` draw cycle and resize handling.
  - **Message Routing**: Status messages like "Different: X - Y" are routed to `diffstatus.SetMessage()` instead of `message.Infof()` to ensure they persist in the dedicated row rather than auto-dismissing. Error messages and search completion messages still use the ephemeral `message` package since the diffstatus row disappears when the search ends.

- **Periodic UI Refresh (`app/goful.go`)**:
  - The `findNextDiff()` function starts a goroutine with a 1-second ticker.
  - The ticker calls `g.Draw()` and `widget.Show()` to refresh the UI.
  - The goroutine is stopped via `defer close(quit)` when search pauses or completes.

- **Command Wrappers (`app/goful.go`)**:
  - Add `func (g *Goful) StartDiffSearch()`:
    - Record initial directories, set active and searching state.
    - Call engine to find first difference.
    - Move cursors, update status.
  - Add `func (g *Goful) ContinueDiffSearch()`:
    - Read cursor filename from active window.
    - Set searching state.
    - Call engine to find next difference after that name.
    - Handle subdirectory descent if needed.
    - Move cursors, update status or "No differences found".
  - Add `func (g *Goful) DiffSearchStatus() string` for status callback.
  - Add `func (g *Goful) IsDiffSearchActive() bool` for active callback.

- **Keymap Integration (`main.go`)**:
  - Bind `[` for start diff search, `]` for continue diff search.
  - Wire `diffstatus.SetStatusFn(g.DiffSearchStatus)` and `diffstatus.SetActiveFn(g.IsDiffSearchActive)`.

**Code Markers**:
- `filer/diffsearch.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`
- `filer/workspace.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`
- `app/goful.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`
- `diffstatus/diffstatus.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`
- `main.go`: `[IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]`

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `filer/diffsearch.go`, `filer/workspace.go`, `app/goful.go`, `diffstatus/diffstatus.go`, `main.go`.
- Tests: `filer/diffsearch_test.go` with tests named `Test*_REQ_DIFF_SEARCH`.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- `go test ./filer/... -run "REQ_DIFF_SEARCH"` (darwin/arm64, Go 1.24.3) on 2026-01-10 validates diff search state and comparison logic (16 tests passing).
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` (2026-01-10) → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 798 token references across 69 files.`
- Message routing alignment: Status messages now route to dedicated `diffstatus` row per `[REQ:DIFF_SEARCH]` specification.
- Bug fix (2026-01-10): Added `FindNextSubdirInAll` to respect `startAfter` during subdirectory descent, fixing search state loss when user manually navigates into subdirectories. Added 4 new unit tests: `TestFindNextSubdirInAll_REQ_DIFF_SEARCH`, `TestFindNextSubdirInAllSkipsNonCommon_REQ_DIFF_SEARCH`, `TestFindNextSubdirInAllNoCommon_REQ_DIFF_SEARCH`.

**Cross-References**: [ARCH:DIFF_SEARCH], [REQ:DIFF_SEARCH], [REQ:MODULE_VALIDATION]

## 34. nsync Observer Adapter [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]

### Decision: Implement an observer adapter that bridges nsync progress events to goful's progress widget.
**Rationale:**
- nsync uses the Observer pattern for progress notifications with callbacks like `OnStart`, `OnProgress`, `OnFinish`.
- goful has an existing `progress` package with `Start()`, `Update()`, `Finish()` functions that render a progress bar.
- An adapter bridges these two systems without modifying either one.

### Implementation Approach:
- Create `type gofulObserver struct` implementing `nsync.Observer` interface in `app/nsync.go`.
- `OnStart(plan)`: Call `progress.Start(float64(plan.TotalBytes))` and `progress.StartTaskCount(plan.TotalItems)`.
- `OnProgress(stats)`: Call `progress.Update(float64(stats.BytesCopied - lastBytes))` with delta tracking.
- `OnItemComplete(item, result)`: Call `progress.FinishTask()` per item, emit `message.Infof` for errors.
- `OnFinish(result)`: Call `progress.Finish()`, emit summary message.
- Observer must be thread-safe; use mutex for byte tracking since nsync calls from multiple goroutines.

**Code Markers**:
- `app/nsync.go` includes `// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]` in the observer struct and methods.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `app/nsync.go`.
- Tests: `app/nsync_test.go` with tests named `TestGofulObserver_REQ_NSYNC_MULTI_TARGET`.

**Cross-References**: [ARCH:NSYNC_INTEGRATION], [REQ:NSYNC_MULTI_TARGET], [REQ:MODULE_VALIDATION]

## 35. nsync Copy/Move Wrappers [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]

### Decision: Provide wrapper functions that configure and execute nsync sync operations within goful's async file control pattern.
**Rationale:**
- nsync's `Syncer.Sync()` is synchronous and blocks until complete; goful needs to run it in a background goroutine.
- The `asyncFilectrl` pattern handles UI resizing, progress widget space, and workspace reload after completion.
- Wrapper functions encapsulate nsync configuration (sources, destinations, recursive, move mode) for clean call sites.

### Implementation Approach:
- `func (g *Goful) syncCopy(sources []string, destinations []string)`:
  - Resolves absolute paths for all sources.
  - Configures `nsync.Config{Sources, Destinations, Recursive: true, Move: false, Jobs: 4}`.
  - Creates syncer with `nsync.WithObserver(gofulObserver)`.
  - Calls within `asyncFilectrl` goroutine pattern.
  - Reports result via `message.Infof`/`message.Error`.
- `func (g *Goful) syncMove(sources []string, destinations []string)`:
  - Same as `syncCopy` but with `Move: true` in config.
  - nsync handles source deletion after successful sync to all destinations.
- Context cancellation: Create `context.WithCancel` that listens for user interrupt (future enhancement).

**Code Markers**:
- `app/nsync.go` includes `// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]` in wrapper functions.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `app/nsync.go`.
- Tests: `app/nsync_test.go` with tests named `TestSyncCopy_REQ_NSYNC_MULTI_TARGET`, `TestSyncMove_REQ_NSYNC_MULTI_TARGET`.

**Cross-References**: [ARCH:NSYNC_INTEGRATION], [REQ:NSYNC_MULTI_TARGET], [IMPL:NSYNC_OBSERVER], [REQ:MODULE_VALIDATION]

## 36. CopyAll/MoveAll Functions [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]

### Decision: Add new cmdline modes that collect all visible workspace directories as destinations and delegate to nsync wrappers.
**Rationale:**
- Users expect symmetry with existing `Copy`/`Move` commands.
- Using `otherWindowDirPaths()` (existing helper for `%D@` macro) provides consistent destination enumeration.
- Fallback to builtin single-target operation when only one pane is visible prevents confusing behavior.

### Implementation Approach:
- `func (g *Goful) CopyAll()`:
  - If only one directory visible: delegate to `g.Copy()` with message explaining fallback.
  - Collect sources: if marks exist, use `g.Dir().MarkfilePaths()`; else use cursor file `g.File().Path()`.
  - Collect destinations: `otherWindowDirPaths(g.Workspace())`.
  - Call `g.syncCopy(sources, destinations)`.
- `func (g *Goful) MoveAll()`:
  - Same pattern as `CopyAll` but calls `g.syncMove`.
- No cmdline text input needed—operation is immediate after command invocation.
- Add to `main.go` keybindings: `C` for CopyAll, `M` for MoveAll.
- Add command menu entries for discoverability.
- Note: `` ` `` (backtick) is used for toggle comparison colors (changed from `C` to avoid conflict).

**Code Markers**:
- `app/nsync.go` includes `// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]` in `CopyAll`/`MoveAll` functions.
- `main.go` keybindings include `// [IMPL:NSYNC_COPY_MOVE] [REQ:NSYNC_MULTI_TARGET]`.

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `app/mode.go`, `main.go`.
- Tests: Integration tests in `app/` verifying destination enumeration and fallback behavior with `[REQ:NSYNC_MULTI_TARGET]`.

**Cross-References**: [ARCH:NSYNC_INTEGRATION], [REQ:NSYNC_MULTI_TARGET], [IMPL:NSYNC_OBSERVER], [REQ:MODULE_VALIDATION]

## 37. nsync Confirmation Modes [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]

### Decision: Implement confirmation prompts for CopyAll/MoveAll using cmdline modes similar to quitMode/removeMode.
**Rationale:**
- Multi-target operations are high-risk and users expect confirmation before files are synced to multiple destinations.
- Reusing the existing cmdline mode pattern keeps implementation simple and UX consistent with other confirmation dialogs.
- The confirmation displays source count and destination count so users understand the scope of the operation.

### Implementation Approach:
- **Add `copyAllMode` struct to `app/mode.go`**:
  - Fields: `*Goful`, `sources []string`, `destinations []string`
  - `String()`: returns `"copyall"`
  - `Prompt()`: returns `fmt.Sprintf("Copy %d file(s) to %d destinations? [Y/n] ", len(sources), len(destinations))`
  - `Draw()`: calls `c.DrawLine()`
  - `Run()`: on `Y`/`y`/empty calls `m.doCopyAll(sources, destinations)` and `c.Exit()`; on `n`/`N` calls `c.Exit()`; else clears text

- **Add `moveAllMode` struct to `app/mode.go`**:
  - Same pattern as `copyAllMode` but with "Move" label and calls `m.doMoveAll()`

- **Refactor `app/nsync.go`**:
  - Rename current `syncCopy`/`syncMove` internals to `doCopyAll()`/`doMoveAll()` (private execution methods)
  - New public `CopyAll()` collects sources/destinations, then starts `copyAllMode` if valid
  - New public `MoveAll()` collects sources/destinations, then starts `moveAllMode` if valid
  - Single-pane fallback logic (`g.Copy()`/`g.Move()`) remains in the public methods before mode creation

**Code Markers**:
- `app/mode.go` confirmation mode structs include `// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]`
- `app/nsync.go` refactored public/private methods include `// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]`

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `app/mode.go`, `app/nsync.go`
- Tests: `app/nsync_test.go` with tests named `TestCopyAllConfirmation_REQ_NSYNC_CONFIRMATION`, `TestMoveAllConfirmation_REQ_NSYNC_CONFIRMATION`

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- To be captured after implementation with `./scripts/validate_tokens.sh` output

**Cross-References**: [ARCH:NSYNC_CONFIRMATION], [REQ:NSYNC_CONFIRMATION], [REQ:NSYNC_MULTI_TARGET], [REQ:MODULE_VALIDATION]

## 37. Linked Navigation Comparison Index Timing Fix [IMPL:LINKED_NAVIGATION] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]

### Issue: Digest comparison decoration missing from focused window after linked navigation

**Rationale:**
When navigating into a subdirectory with linked navigation enabled, pressing `=` to calculate file digests would only decorate OTHER windows - the SOURCE (focused) window was not decorated correctly. This caused the message to report "calculated digest for 2 files" instead of 3, and the focused window showed no underline decoration.

**Root Cause Analysis (via runtime instrumentation):**
The comparison index was being rebuilt inside `ChdirAllToSubdir()` BEFORE the focused directory navigated via `EnterDir()`. When the index was built, the focused directory (`dirIdx:0`) still pointed to the parent directory contents, so files only present in the subdirectory were not indexed for the focused window.

**Log Evidence (2026-01-11):**
```
{"message":"BEFORE RebuildComparisonIndex","data":{"paths":[".../notes_git", ".../notes_git/dev", ".../notes_git/dev"]}}
{"message":"scanning dir","data":{"dirIdx":0,"path":".../notes_git","listLen":13}}  // ← OLD PATH!
{"message":"cache lookup","data":{"filename":"tror.txt","cachedDirIndices":[1,2]}}  // ← NO 0!
```

### Fix Approach:
1. Added `ChdirAllToSubdirNoRebuild()` method to `filer.Workspace` that navigates non-focused directories WITHOUT rebuilding the comparison index.
2. Modified `linkedEnterDir` in `main.go` to:
   - Call `ChdirAllToSubdirNoRebuild()` to navigate other directories
   - Call `EnterDir()` to navigate the focused directory
   - Call `RebuildComparisonIndex()` AFTER all directories have navigated

This ensures the comparison index is built with the correct file lists from all directories after they've all navigated to the new location.

**Code Markers**:
- `filer/workspace.go`: `ChdirAllToSubdirNoRebuild()` with `// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `main.go`: `linkedEnterDir` function with comments documenting the navigation sequence

**Token Coverage** `[PROC:TOKEN_AUDIT]`:
- Source: `filer/workspace.go`, `main.go`
- Existing tests cover linked navigation behavior; manual verification confirmed fix

**Validation Evidence** `[PROC:TOKEN_VALIDATION]` (2026-01-11):
- `go test ./...` passes on darwin/arm64, Go 1.24.3
- Manual verification: pressing `=` after linked navigation now reports "calculated digest for 3 files" and all windows show underline decoration
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 905 token references across 71 files.`

**Cross-References**: [ARCH:LINKED_NAVIGATION], [REQ:LINKED_NAVIGATION], [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS], [IMPL:DIGEST_COMPARISON]
