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

## 2. Core Implementation [IMPL:EXAMPLE_IMPLEMENTATION] [ARCH:EXAMPLE_DECISION] [REQ:EXAMPLE_FEATURE]

### Data Structure
```[your-language]
type ExampleStruct struct {
    Field1 string
    Field2 int
}
```

### Implementation Approach
- Approach description
- Key algorithms
- Performance considerations

### Platform-Specific Considerations
- Platform 1: Specific considerations
- Platform 2: Specific considerations

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
function testExampleFeature_REQ_EXAMPLE_FEATURE() {
    // [REQ:EXAMPLE_FEATURE] Validates expected behavior
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
function testIntegrationScenario_REQ_EXAMPLE_FEATURE() {
    // [REQ:EXAMPLE_FEATURE] End-to-end validation comment
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
- Add `.github/workflows/ci.yml` with steps: checkout, setup-go, cache, `go fmt`/`gofmt -w`, `go vet`, `go test ./...`.
- Include `DEBUG:`/`DIAGNOSTIC:` logs for key steps.

**Code Markers**:
- Workflow file comments with `[IMPL:CI_WORKFLOW] [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]`.

**Validation Evidence** `[PROC:TOKEN_VALIDATION]`:
- Link to latest passing CI run and token validation output.

**Cross-References**: [ARCH:CI_PIPELINE], [REQ:CI_PIPELINE_CORE]

## 10. Staticcheck Setup [IMPL:STATICCHECK_SETUP] [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]

### Decision: Add static analysis job
**Rationale:**
- Surface correctness issues early.

### Implementation Approach:
- Add workflow job invoking staticcheck (and golangci-lint if configured).
- Provide config file if needed with minimal excludes.

**Code Markers**:
- Workflow job comments carry `[IMPL:STATICCHECK_SETUP] [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]`.

**Cross-References**: [ARCH:STATIC_ANALYSIS_POLICY], [REQ:STATIC_ANALYSIS]

## 11. Race Job [IMPL:RACE_JOB] [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]

### Decision: Dedicated race-enabled test job
**Rationale:**
- Detects concurrency issues without impacting main job runtime.

### Implementation Approach:
- Add CI job running `go test -race ./...` with appropriate timeouts and cache.

**Code Markers**:
- Workflow job comments carry `[IMPL:RACE_JOB] [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]`.

**Cross-References**: [ARCH:RACE_TESTING_PIPELINE], [REQ:RACE_TESTING]

## 12. Widget Tests [IMPL:TEST_WIDGETS] [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]

### Decision: Add unit/snapshot tests for widget primitives
**Rationale:**
- Protect rendering/event handling behaviors.

### Implementation Approach:
- Add tests under `widget/test/...` and `filer/`.
- Use golden/snapshot helpers for rendering where applicable.

**Code Markers**:
- Test files include `[IMPL:TEST_WIDGETS] [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]`.

**Cross-References**: [ARCH:TEST_STRATEGY_UI], [REQ:UI_PRIMITIVE_TESTS]

## 13. Command Tests [IMPL:TEST_CMDLINE] [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]

### Decision: Add tests for command parsing and modes
**Rationale:**
- Prevent regressions in command handling and state transitions.

### Implementation Approach:
- Unit tests for `cmdline/` parsing and `app/` mode transitions.
- Cover error paths and edge cases.

**Code Markers**:
- Tests carry `[IMPL:TEST_CMDLINE] [ARCH:TEST_STRATEGY_CMD] [REQ:CMD_HANDLER_TESTS]`.

**Cross-References**: [ARCH:TEST_STRATEGY_CMD], [REQ:CMD_HANDLER_TESTS]

## 14. Integration Flow Tests [IMPL:TEST_INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]

### Decision: Snapshot/integration tests for file operation flows
**Rationale:**
- Validate end-to-end behavior for open/navigate/rename/delete.

### Implementation Approach:
- Use fixtures to simulate file operations and capture output snapshots.
- Include diagnostic logging for traceability.

**Code Markers**:
- Tests annotated with `[IMPL:TEST_INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [REQ:INTEGRATION_FLOWS]`.

**Cross-References**: [ARCH:TEST_STRATEGY_INTEGRATION], [REQ:INTEGRATION_FLOWS]

## 15. Architecture Guide [IMPL:DOC_ARCH_GUIDE] [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]

### Decision: Write `ARCHITECTURE.md`
**Rationale:**
- Provides concise understanding of packages and data flow.

### Implementation Approach:
- Summarize UI widgets, file ops, app/mode pipeline, and data flow diagrams/text.
- Cross-link tokens and modules.

**Code Markers**:
- Document contains `[IMPL:DOC_ARCH_GUIDE] [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]`.

**Cross-References**: [ARCH:DOCS_STRUCTURE], [REQ:ARCH_DOCUMENTATION]

## 16. CONTRIBUTING Guide [IMPL:DOC_CONTRIBUTING] [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE]

### Decision: Add contributor standards document
**Rationale:**
- Aligns development workflow and review expectations.

### Implementation Approach:
- Cover coding standards, branching, review flow, token usage, required checks.
- Link to Makefile/CI targets.

**Code Markers**:
- `CONTRIBUTING.md` includes `[IMPL:DOC_CONTRIBUTING] [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE]`.

**Cross-References**: [ARCH:CONTRIBUTION_PROCESS], [REQ:CONTRIBUTING_GUIDE]

## 17. Release Targets [IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]

### Decision: Makefile and CI matrix for reproducible builds
**Rationale:**
- Deliver consistent artifacts across GOOS/GOARCH.

### Implementation Approach:
- Add Makefile targets for lint/test/build.
- CI matrix builds static binaries for target platforms; name artifacts deterministically.

**Code Markers**:
- Makefile and workflow comments carry `[IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]`.

**Cross-References**: [ARCH:BUILD_MATRIX], [REQ:RELEASE_BUILD_MATRIX]

## 18. Baseline Snapshots [IMPL:BASELINE_SNAPSHOTS] [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE]

### Decision: Capture current keybindings/modes as automated baselines
**Rationale:**
- Preserve behavior ahead of refactors.

### Implementation Approach:
- Add tests or scripts that exercise key flows and record expected outputs.
- Store fixtures for comparison.

**Code Markers**:
- Tests/scripts include `[IMPL:BASELINE_SNAPSHOTS] [ARCH:BASELINE_CAPTURE] [REQ:BEHAVIOR_BASELINE]`.

**Cross-References**: [ARCH:BASELINE_CAPTURE], [REQ:BEHAVIOR_BASELINE]

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

