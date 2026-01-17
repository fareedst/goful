# [IMPL:RACE_JOB] Race Job

**Cross-References**: [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Dedicated race-enabled test job.

## Rationale

- Detects concurrency issues without impacting main job runtime

## Implementation Approach

- Added `race-tests` job to `.github/workflows/ci.yml` that:
  - Sets up Go `1.24.3` with caching and reuses the module cache
  - Executes `go test -race ./...` so concurrency regressions fail CI early

## Code Markers

- Workflow job comments carry `[IMPL:RACE_JOB] [ARCH:RACE_TESTING_PIPELINE] [REQ:RACE_TESTING]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `.github/workflows/ci.yml` - race-tests job

Tests that must reference `[REQ:RACE_TESTING]`:
- [ ] N/A - the race detector is the test

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | Race tests pass without detected races |

## Related Decisions

- Depends on: [IMPL:CI_WORKFLOW]
- See also: [ARCH:RACE_TESTING_PIPELINE], [REQ:RACE_TESTING]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
