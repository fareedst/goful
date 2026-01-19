# [IMPL:TESTING] Testing Implementation

**Cross-References**: [ARCH:TESTING_STRATEGY] [REQ:*]  
**Status**: Active  
**Created**: 2025-12-18  
**Last Updated**: 2026-01-18

---

## Decision

Implement a comprehensive testing strategy that realizes the validation criteria specified in `requirements.md` and follows the testing strategy defined in `architecture-decisions.md`.

## Rationale

- Each test validates specific satisfaction criteria from requirements
- Provides traceability from tests back to requirements via `[REQ:*]` tokens
- Ensures module validation before integration per `[REQ:MODULE_VALIDATION]`

## Implementation Approach

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

### Test Reliability: Parallel Execution Compatibility

Tests that change the working directory must handle cases where `os.Getwd()` fails during parallel test execution. This can occur on Ubuntu when the current working directory is deleted by another test's cleanup before `os.Getwd()` is called.

**Fix Pattern** (2026-01-18):
- Use `os.TempDir()` as a fallback when `os.Getwd()` fails
- Ensures tests always have a valid directory to restore to, even if the original directory no longer exists
- Applied to `newTestDirectory()` helper in `filer/integration_test.go` and `TestFindFileByName_REQ_SYNC_COMMANDS` in `app/window_wide_test.go`

```go
origDir, err := os.Getwd()
if err != nil {
    // Fallback to system temp directory if current directory is deleted
    // (can happen during parallel test execution on Ubuntu)
    origDir = os.TempDir()
}
```

## Code Markers

- Test files include `[IMPL:TESTING] [ARCH:TESTING_STRATEGY]` comments
- Test function names include `_REQ_*` suffix for traceability

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] All test files referencing requirements

Tests that must reference `[REQ:*]`:
- [ ] Every test function must include the relevant `[REQ:*]` token

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2025-12-18 | — | ⏳ Pending | Template placeholder |
| 2026-01-18 | — | ✅ Pass | Fixed Ubuntu test failures from `os.Getwd()` errors during parallel execution |

## Related Decisions

- Depends on: —
- See also: [ARCH:TESTING_STRATEGY], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
