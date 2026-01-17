# [IMPL:STATICCHECK_SETUP] Staticcheck Setup

**Cross-References**: [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Add static analysis job.

## Rationale

- Surface correctness issues early

## Implementation Approach

- Added `staticcheck` job to `.github/workflows/ci.yml` that:
  - Reuses the Go `1.24.3` toolchain setup with cached modules
  - Installs `staticcheck` via `go install honnef.co/go/tools/cmd/staticcheck@latest`
  - Runs `staticcheck ./...` with `[IMPL:STATICCHECK_SETUP] [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]` inline comments

## Code Markers

- Workflow job comments carry `[IMPL:STATICCHECK_SETUP] [ARCH:STATIC_ANALYSIS_POLICY] [REQ:STATIC_ANALYSIS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `.github/workflows/ci.yml` - staticcheck job

Tests that must reference `[REQ:STATIC_ANALYSIS]`:
- [ ] N/A - workflow job is the validation

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | staticcheck passes on codebase |

## Related Decisions

- Depends on: [IMPL:CI_WORKFLOW]
- See also: [ARCH:STATIC_ANALYSIS_POLICY], [REQ:STATIC_ANALYSIS]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
